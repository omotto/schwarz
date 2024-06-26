package kubernetes

import (
	"context"
	"schwarz/models"
	"schwarz/services/prometheus"

	"github.com/google/uuid"
	appsv1 "k8s.io/api/apps/v1"
	apiv1 "k8s.io/api/core/v1"
	"k8s.io/apimachinery/pkg/api/resource"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/client-go/kubernetes"
	"k8s.io/client-go/util/retry"
)

const (
	postgresPrefix            = "postgres-"
	postgresVolumePrefix      = "postgres-volume-"
	postgresVolumeClaimPrefix = "postgres-volume-claim-"
	postgresSecretPrefix      = "postgres-secret-"
)

type Postgres struct {
	kubeClient *kubernetes.Clientset
	metrics    *prometheus.Prometheus
}

func NewPostgres(clientset *kubernetes.Clientset, metrics *prometheus.Prometheus) Service {
	return &Postgres{
		kubeClient: clientset,
		metrics:    metrics,
	}
}

func (s *Postgres) Create(ctx context.Context, request models.CreateRequest) (models.CreateResponse, error) {
	id := uuid.New().String()
	_ = s.metrics.IncreaseCounterMetric(prometheus.MetricDeploymentAccessTotal, 1, map[string]string{prometheus.LabelID: id, prometheus.LabelOperation: "create"})
	configMap := setConfigMap(request.DBName, request.UserName, request.UserPass, id)
	persistentVolume := setPersistentVolume(request.Capacity, []string{request.AccessMode}, id)
	persistentVolumeClaim := setPersistentVolumeClaim(request.Capacity, []string{request.AccessMode}, id)
	deployment := setDeployment(request.Replicas, request.PortNum, id)
	service := setService(request.PortNum, id)
	if _, err := s.kubeClient.CoreV1().ConfigMaps(apiv1.NamespaceDefault).Create(ctx, configMap, metav1.CreateOptions{}); err != nil {
		return models.CreateResponse{}, err
	} else if _, err := s.kubeClient.CoreV1().PersistentVolumes().Create(ctx, persistentVolume, metav1.CreateOptions{}); err != nil {
		return models.CreateResponse{}, err
	} else if _, err := s.kubeClient.CoreV1().PersistentVolumeClaims(apiv1.NamespaceDefault).Create(ctx, persistentVolumeClaim, metav1.CreateOptions{}); err != nil {
		return models.CreateResponse{}, err
	} else if _, err = s.kubeClient.AppsV1().Deployments(apiv1.NamespaceDefault).Create(ctx, deployment, metav1.CreateOptions{}); err != nil {
		_ = s.metrics.IncreaseCounterMetric(prometheus.MetricDeploymentAccessFailedTotal, 1, map[string]string{prometheus.LabelID: id, prometheus.LabelOperation: "create"})
		return models.CreateResponse{}, err
	} else if _, err := s.kubeClient.CoreV1().Services(apiv1.NamespaceDefault).Create(ctx, service, metav1.CreateOptions{}); err != nil {
		return models.CreateResponse{}, err
	}
	return models.CreateResponse{ID: id}, nil
}

func (s *Postgres) Delete(ctx context.Context, request models.DeleteRequest) error {
	_ = s.metrics.IncreaseCounterMetric(prometheus.MetricDeploymentAccessTotal, 1, map[string]string{prometheus.LabelID: request.ID, prometheus.LabelOperation: "delete"})
	deletePolicy := metav1.DeletePropagationForeground
	if err := s.kubeClient.AppsV1().Deployments(apiv1.NamespaceDefault).Delete(ctx, request.ID, metav1.DeleteOptions{
		PropagationPolicy: &deletePolicy,
	}); err != nil {
		_ = s.metrics.IncreaseCounterMetric(prometheus.MetricDeploymentAccessFailedTotal, 1, map[string]string{prometheus.LabelID: request.ID, prometheus.LabelOperation: "delete"})
		return err
	}
	if err := s.kubeClient.CoreV1().Services(apiv1.NamespaceDefault).Delete(ctx, postgresPrefix+request.ID, metav1.DeleteOptions{}); err != nil {
		return err
	}
	if err := s.kubeClient.CoreV1().PersistentVolumeClaims(apiv1.NamespaceDefault).Delete(ctx, postgresVolumeClaimPrefix+request.ID, metav1.DeleteOptions{}); err != nil {
		return err
	}
	if err := s.kubeClient.CoreV1().PersistentVolumes().Delete(ctx, postgresVolumePrefix+request.ID, metav1.DeleteOptions{}); err != nil {
		return err
	}
	if err := s.kubeClient.CoreV1().ConfigMaps(apiv1.NamespaceDefault).Delete(ctx, postgresSecretPrefix+request.ID, metav1.DeleteOptions{}); err != nil {
		return err
	}
	return nil
}

func (s *Postgres) Update(ctx context.Context, request models.UpdateRequest) error {
	_ = s.metrics.IncreaseCounterMetric(prometheus.MetricDeploymentAccessTotal, 1, map[string]string{prometheus.LabelID: request.ID, prometheus.LabelOperation: "update"})
	return retry.RetryOnConflict(retry.DefaultRetry, func() error {
		// Retrieve the latest version of Deployment before attempting update
		// RetryOnConflict uses exponential backoff to avoid exhausting the apiserver
		result, err := s.kubeClient.AppsV1().Deployments(apiv1.NamespaceDefault).Get(ctx, request.ID, metav1.GetOptions{})
		if err != nil {
			_ = s.metrics.IncreaseCounterMetric(prometheus.MetricDeploymentAccessFailedTotal, 1, map[string]string{prometheus.LabelID: request.ID, prometheus.LabelOperation: "read"})
			return err
		}
		result.Spec.Replicas = &request.Replicas
		_, err = s.kubeClient.AppsV1().Deployments(apiv1.NamespaceDefault).Update(ctx, result, metav1.UpdateOptions{})
		if err != nil {
			_ = s.metrics.IncreaseCounterMetric(prometheus.MetricDeploymentAccessFailedTotal, 1, map[string]string{prometheus.LabelID: request.ID, prometheus.LabelOperation: "update"})
		}
		return err
	})
}

func setConfigMap(dbName, user, pass, id string) *apiv1.ConfigMap {
	labelData := map[string]string{
		"app": "postgres",
	}
	postgresData := map[string]string{
		"POSTGRES_DB":       dbName,
		"POSTGRES_USER":     user,
		"POSTGRES_PASSWORD": pass,
	}
	return &apiv1.ConfigMap{
		TypeMeta: metav1.TypeMeta{
			Kind:       "ConfigMap",
			APIVersion: "v1",
		},
		ObjectMeta: metav1.ObjectMeta{
			Name:   postgresSecretPrefix + id,
			Labels: labelData,
		},
		Data: postgresData,
	}
}

func setPersistentVolume(storage string, accessModes []string, id string) *apiv1.PersistentVolume {
	persistentVolumeAccessModes := make([]apiv1.PersistentVolumeAccessMode, len(accessModes))
	for idx, accessMode := range accessModes {
		persistentVolumeAccessModes[idx] = apiv1.PersistentVolumeAccessMode(accessMode)
	}
	capacity := apiv1.ResourceList{apiv1.ResourceStorage: resource.MustParse(storage)}
	labelData := map[string]string{
		"app":  "postgres",
		"type": "local",
	}
	return &apiv1.PersistentVolume{
		TypeMeta: metav1.TypeMeta{
			Kind:       "PersistentVolume",
			APIVersion: "v1",
		},
		ObjectMeta: metav1.ObjectMeta{
			Name:   postgresVolumePrefix + id,
			Labels: labelData,
		},
		Spec: apiv1.PersistentVolumeSpec{
			StorageClassName: "manual",
			AccessModes:      persistentVolumeAccessModes,
			Capacity:         capacity,
			PersistentVolumeSource: apiv1.PersistentVolumeSource{
				HostPath: &apiv1.HostPathVolumeSource{
					Path: "/data/postgresql",
				},
			},
		},
	}
}

func setPersistentVolumeClaim(storage string, accessModes []string, id string) *apiv1.PersistentVolumeClaim {
	persistentVolumeClaimAccessModes := make([]apiv1.PersistentVolumeAccessMode, len(accessModes))
	for idx, accessMode := range accessModes {
		persistentVolumeClaimAccessModes[idx] = apiv1.PersistentVolumeAccessMode(accessMode)
	}
	storageClassName := "manual"
	labelData := map[string]string{
		"app": "postgres",
	}
	capacity := apiv1.ResourceList{apiv1.ResourceStorage: resource.MustParse(storage)}
	return &apiv1.PersistentVolumeClaim{
		TypeMeta: metav1.TypeMeta{
			Kind:       "PersistentVolumeClaim",
			APIVersion: "v1",
		},
		ObjectMeta: metav1.ObjectMeta{
			Name:   postgresVolumeClaimPrefix + id,
			Labels: labelData,
		},
		Spec: apiv1.PersistentVolumeClaimSpec{
			AccessModes:      persistentVolumeClaimAccessModes,
			StorageClassName: &storageClassName,
			Resources: apiv1.VolumeResourceRequirements{
				Requests: capacity,
			},
		},
	}
}

func setService(port int32, id string) *apiv1.Service {
	labelData := map[string]string{
		"app": "postgres",
	}
	selector := map[string]string{
		"app": "postgres",
	}
	return &apiv1.Service{
		TypeMeta: metav1.TypeMeta{
			Kind:       "Service",
			APIVersion: "v1",
		},
		ObjectMeta: metav1.ObjectMeta{
			Name:   postgresPrefix + id,
			Labels: labelData,
		},
		Spec: apiv1.ServiceSpec{
			Ports:    []apiv1.ServicePort{{Port: port}},
			Selector: selector,
			Type:     "NodePort",
		},
	}
}

func setDeployment(replicas, port int32, id string) *appsv1.Deployment {
	matchLabels := map[string]string{"app": "postgres"}
	labels := map[string]string{"app": "postgres"}
	pullPolicy := "IfNotPresent"
	return &appsv1.Deployment{
		TypeMeta: metav1.TypeMeta{
			Kind:       "Deployment",
			APIVersion: "apps/v1",
		},
		ObjectMeta: metav1.ObjectMeta{
			Name: id,
		},
		Spec: appsv1.DeploymentSpec{
			Replicas: &replicas,
			Selector: &metav1.LabelSelector{
				MatchLabels: matchLabels,
			},
			Template: apiv1.PodTemplateSpec{
				ObjectMeta: metav1.ObjectMeta{
					Labels: labels,
				},
				Spec: apiv1.PodSpec{
					Volumes: []apiv1.Volume{{
						Name: "postgresdata",
						VolumeSource: apiv1.VolumeSource{
							PersistentVolumeClaim: &apiv1.PersistentVolumeClaimVolumeSource{
								ClaimName: postgresVolumeClaimPrefix + id,
							},
						},
					}},
					Containers: []apiv1.Container{{
						Name:            "postgres",
						Image:           "postgres:14",
						ImagePullPolicy: apiv1.PullPolicy(pullPolicy),
						Ports: []apiv1.ContainerPort{{
							ContainerPort: port,
						}},
						EnvFrom: []apiv1.EnvFromSource{{
							ConfigMapRef: &apiv1.ConfigMapEnvSource{
								LocalObjectReference: apiv1.LocalObjectReference{
									Name: postgresSecretPrefix + id,
								},
							},
						}},
						VolumeMounts: []apiv1.VolumeMount{{
							Name:      "postgresdata",
							MountPath: "/var/lib/postgresql/data",
						}},
					}},
				},
			},
		},
	}
}
