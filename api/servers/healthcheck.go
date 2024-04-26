package servers

import (
	"context"
	"fmt"
	"log"
	"net/http"
	"schwarz/api/handlers"
	"time"

	"github.com/gorilla/mux"
	"github.com/heptiolabs/healthcheck"
)

const (
	// default kubernetes readiness probe path
	defaultReadyProbe = "/.well-known/ready"
	// default kubernetes liveness probe path
	defaultLiveProbe = "/.well-known/live"
	// default (Prometheus) metrics scraping endpoint path
	defaultMetricsPath = "/metrics"
)

type HealthcheckServer struct {
	httpAddr string
	router   *mux.Router
	timeout  time.Duration
	httpSvr  *http.Server
}

func NewHealthcheck(host, port string, timeout time.Duration, metrics handlers.Metrics, handler healthcheck.Handler) HealthcheckServer {
	server := HealthcheckServer{
		httpAddr: fmt.Sprintf("%s:%s", host, port),
		router:   mux.NewRouter(),
		timeout:  timeout,
	}
	server.router.HandleFunc(defaultReadyProbe, handler.ReadyEndpoint).Methods(http.MethodGet)
	server.router.HandleFunc(defaultLiveProbe, handler.LiveEndpoint).Methods(http.MethodGet)
	server.router.HandleFunc(defaultMetricsPath, metrics.MetricsEndpoint).Methods(http.MethodGet)
	server.httpSvr = &http.Server{
		Addr:              server.httpAddr,
		Handler:           server.router,
		ReadHeaderTimeout: server.timeout,
		ReadTimeout:       server.timeout,
		WriteTimeout:      server.timeout,
	}
	return server
}

func (s *HealthcheckServer) Run() {
	log.Println("Server running on", s.httpAddr)
	go func() {
		if err := s.httpSvr.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			log.Fatal("server shut down", err)
		}
	}()
}

func (s *HealthcheckServer) ShutDown() error {
	ctxShutDown, cancel := context.WithTimeout(context.Background(), s.timeout)
	defer cancel()
	return s.httpSvr.Shutdown(ctxShutDown)
}
