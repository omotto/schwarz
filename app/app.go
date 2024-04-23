package app

import (
	"context"
	"fmt"
	"os"
	"path/filepath"

	"k8s.io/client-go/kubernetes"
	"k8s.io/client-go/tools/clientcmd"

	"schwarz/models"
	"schwarz/services"
)

func Start() {

	kubeConfigPath := filepath.Join(os.Getenv("HOME"), ".kube", "config")
	kubeConfig, err := clientcmd.BuildConfigFromFlags("", kubeConfigPath)
	if err != nil {
		panic(err)
	}
	kubeClient, err := kubernetes.NewForConfig(kubeConfig)
	if err != nil {
		panic(err)
	}
	fmt.Println(kubeClient)
	service := services.New(kubeClient)
	resp, err := service.Create(context.Background(), models.CreateRequest{
		DBName:     "mottoDB",
		UserName:   "omotto",
		UserPass:   "123456",
		PortNum:    5432,
		Replicas:   2,
		Capacity:   "10Mi",
		AccessMode: "ReadWriteMany",
	})
	fmt.Println(err)
	err = service.Update(context.Background(), models.UpdateRequest{
		Id:       resp.Id,
		Replicas: 4,
	})
	fmt.Println(err)
	err = service.Delete(context.Background(), models.DeleteRequest{
		Id: resp.Id,
	})
	fmt.Println(err)
}
