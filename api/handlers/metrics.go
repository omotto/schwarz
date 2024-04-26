package handlers

import (
	"net/http"

	"github.com/prometheus/client_golang/prometheus"
	"github.com/prometheus/client_golang/prometheus/promhttp"
)

type Metrics struct {
	registry *prometheus.Registry
}

func NewMetrics(registry *prometheus.Registry) Metrics {
	return Metrics{
		registry: registry,
	}
}

func (m *Metrics) MetricsEndpoint(writer http.ResponseWriter, request *http.Request) {
	promhttp.HandlerFor(m.registry, promhttp.HandlerOpts{}).ServeHTTP(writer, request)
}
