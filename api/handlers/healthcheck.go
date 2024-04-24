package handlers

import (
	"net/http"

	"github.com/prometheus/client_golang/prometheus"
	"github.com/prometheus/client_golang/prometheus/promhttp"
)

type Healthcheck struct {
	registry *prometheus.Registry
}

func NewHealthcheck(registry *prometheus.Registry) Healthcheck {
	return Healthcheck{
		registry: registry,
	}
}

func (h *Healthcheck) ReadinessProbe(writer http.ResponseWriter, _ *http.Request) {
	// TODO: Check if all external dependencies are ready either
	writer.WriteHeader(http.StatusOK)
}

func (h *Healthcheck) LivenessProbe(writer http.ResponseWriter, _ *http.Request) {
	writer.WriteHeader(http.StatusOK)
}

func (h *Healthcheck) MetricsHandler(writer http.ResponseWriter, request *http.Request) {
	promhttp.HandlerFor(h.registry, promhttp.HandlerOpts{}).ServeHTTP(writer, request)
}
