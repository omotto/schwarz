package app

import (
	"context"
	"fmt"
	"log"
	"net"
	"os"
	"os/signal"
	"path/filepath"
	"syscall"

	"google.golang.org/grpc"
	"k8s.io/client-go/kubernetes"
	"k8s.io/client-go/tools/clientcmd"

	pb "schwarz/api/proto"
	"schwarz/api/servers"
	"schwarz/services"
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
	log.Println("starting grpc server")
	grpcServer := grpc.NewServer([]grpc.ServerOption{}...)
	pb.RegisterPostgresServiceServer(grpcServer, servers.NewPostgres(validatorService))
	sCtx := serverContext(context.Background())
	go func() {
		log.Printf("gRPC server listening at %v", listener.Addr())
		if err = grpcServer.Serve(listener); err != nil {
			log.Fatalf("failed to serve: %v", err)
		}
	}()
	<-sCtx.Done()
	grpcServer.GracefulStop()
	log.Println("clean shutdown")
}

func serverContext(ctx context.Context) context.Context {
	c := make(chan os.Signal, 1)
	signal.Notify(c, os.Interrupt, syscall.SIGTERM)
	ctx, cancel := context.WithCancel(ctx)
	go func() {
		s := <-c
		log.Printf("got signal %v, attempting graceful shutdown", s)
		cancel()
	}()
	return ctx
}
