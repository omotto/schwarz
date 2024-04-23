package app

import (
	"fmt"
	"log"
	"net"
	"os"
	"path/filepath"
	"schwarz/api/servers"
	"schwarz/services"

	"google.golang.org/grpc"
	"k8s.io/client-go/kubernetes"
	"k8s.io/client-go/tools/clientcmd"

	pb "schwarz/api/proto"
)

func Start() {
	// Kubernetes Init
	kubeConfigPath := filepath.Join(os.Getenv("HOME"), ".kube", "config")
	kubeConfig, err := clientcmd.BuildConfigFromFlags("", kubeConfigPath)
	if err != nil {
		log.Fatalf("failed to load kubeConfig path: %v", err)
	}
	kubeClient, err := kubernetes.NewForConfig(kubeConfig)
	if err != nil {
		log.Fatalf("failed to load kubeConfig: %v", err)
	}

	// Postgres Service Init
	postgresService := services.NewPostgres(kubeClient)

	// Validator Service Init
	validatorService := services.NewValidator(postgresService)

	// GRPC Server Init
	port := 5000
	listener, err := net.Listen("tcp", fmt.Sprintf("localhost:%d", port))
	if err != nil {
		log.Fatalf("failed to listen: %v", err)
	}
	grpcServer := grpc.NewServer()
	pb.RegisterPostgresServiceServer(grpcServer, servers.NewPostgres(validatorService))
	if err = grpcServer.Serve(listener); err != nil {
		log.Fatalf("failed to serve: %v", err)
	}
}
