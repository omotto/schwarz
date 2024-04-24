package app

import (
	"context"
	"fmt"
	"log"
	"net"
	"os"
	"os/signal"
	"path/filepath"
	"schwarz/api/middlewares"
	"schwarz/api/servers"
	"schwarz/services"
	"syscall"

	grpcLogrus "github.com/grpc-ecosystem/go-grpc-middleware/logging/logrus"
	grpcRecovery "github.com/grpc-ecosystem/go-grpc-middleware/recovery"
	"github.com/sirupsen/logrus"
	"google.golang.org/grpc"
	"k8s.io/client-go/kubernetes"
	"k8s.io/client-go/tools/clientcmd"

	pb "schwarz/api/proto"
)

func Start() {
	// Get config params
	cfg, err := NewConfig()
	if err != nil {
		log.Fatalf("failed to load config params %v", err)
	}

	// Kubernetes Init
	kubeConfigPath, err := filepath.Abs("./configs/config")
	if err != nil {
		log.Fatalf("failed to load kubeConfig path: %v", err)
	}
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

	// Server Context
	sCtx := serverContext(context.Background())

	// HTTP Server Init
	httpServer := servers.NewHealthcheck("", cfg.HealthPort, cfg.HttpTimeout)
	httpServer.Run()

	// server-side TLS
	// TODO: disabled to check with INSOMNIA gRPC client
	// creds, err := credentials.NewServerTLSFromFile(testdata.Path("server1.pem"), testdata.Path("server1.key"))
	// if err != nil {
	//	log.Fatalf("failed to create credentials: %v", err)
	//}

	// GRPC Server Init
	listener, err := net.Listen("tcp", fmt.Sprintf(":%s", cfg.GRPCPort))
	if err != nil {
		log.Fatalf("failed to listen: %v", err)
	}
	log.Println("starting grpc server")
	grpcServer := grpc.NewServer([]grpc.ServerOption{
		//grpc.Creds(creds),
		grpc.ChainUnaryInterceptor(
			middlewares.AuthInterceptor,
			grpcLogrus.UnaryServerInterceptor(logrus.NewEntry(logrus.New()), []grpcLogrus.Option{}...),
			grpcRecovery.UnaryServerInterceptor(),
		),
		grpc.ChainStreamInterceptor(
			grpcRecovery.StreamServerInterceptor(),
		),
	}...)
	pb.RegisterPostgresServiceServer(grpcServer, servers.NewPostgres(validatorService))

	// GRPC Run
	go func() {
		log.Printf("gRPC server listening at %v", listener.Addr())
		if err = grpcServer.Serve(listener); err != nil {
			log.Fatalf("failed to serve: %v", err)
		}
	}()
	<-sCtx.Done()
	grpcServer.GracefulStop()
	err = httpServer.ShutDown()
	if err != nil {
		log.Println("error: ", err)
	}
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
