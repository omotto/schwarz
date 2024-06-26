package models

// https://www.digitalocean.com/community/tutorials/how-to-deploy-postgres-to-kubernetes-cluster

type CreateRequest struct {
	DBName     string
	UserName   string
	UserPass   string
	PortNum    int32  // Number of port to expose on the pod's IP address
	Replicas   int32  // Number of desired pods.
	Capacity   string // https://kubernetes.io/docs/concepts/storage/persistent-volumes#resources
	AccessMode string // https://kubernetes.io/docs/concepts/storage/persistent-volumes#binding
}

type CreateResponse struct {
	ID string
}

type DeleteRequest struct {
	ID string
}

type UpdateRequest struct {
	ID       string
	Replicas int32 // Number of desired pods.
}
