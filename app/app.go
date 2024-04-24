package app

import (
	"context"
	"fmt"
	"log"
	"net"
	"os"
	"os/signal"
	"path/filepath"
	"schwarz/api/handlers"
	"schwarz/api/middlewares"
	"schwarz/api/servers"
	kubernetesService "schwarz/services/kubernetes"
	prometheusService "schwarz/services/prometheus"
	"syscall"

	grpcLogrus "github.com/grpc-ecosystem/go-grpc-middleware/logging/logrus"
	grpcPrometheus "github.com/grpc-ecosystem/go-grpc-middleware/providers/prometheus"
	grpcRecovery "github.com/grpc-ecosystem/go-grpc-middleware/recovery"
	"github.com/prometheus/client_golang/prometheus"
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

	// Prometheus Metrics Services
	serverMetrics := grpcPrometheus.NewServerMetrics(grpcPrometheus.WithServerHandlingTimeHistogram())
	registry := prometheus.NewRegistry()
	registry.MustRegister(serverMetrics)
	// Prometheus custom metrics
	customMetrics, err := prometheusService.NewPrometheusService(prometheusService.GetMetricsDefinition())
	if err != nil {
		log.Fatalf("failed to load metric definitions: %v", err)
	}
	registry.MustRegister(customMetrics)

	// Postgres Service Init
	postgresService := kubernetesService.NewPostgres(kubeClient, customMetrics)

	// Validator Service Init
	validatorService := kubernetesService.NewValidator(postgresService)

	// Server Context
	sCtx := serverContext(context.Background())

	// HTTP Server Init
	healthcheckHandlers := handlers.NewHealthcheck(registry)
	httpServer := servers.NewHealthcheck("", cfg.HealthPort, cfg.HttpTimeout, healthcheckHandlers)
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
		// TODO: disabled to check with INSOMNIA gRPC client
		// grpc.Creds(creds),
		grpc.ChainUnaryInterceptor(
			middlewares.AuthInterceptor,
			grpcLogrus.UnaryServerInterceptor(logrus.NewEntry(logrus.New()), []grpcLogrus.Option{}...),
			serverMetrics.UnaryServerInterceptor(),
			grpcRecovery.UnaryServerInterceptor(),
		),
		grpc.ChainStreamInterceptor(
			serverMetrics.StreamServerInterceptor(),
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
	registry.Unregister(customMetrics)
	registry.Unregister(serverMetrics)
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
