package prometheus

import (
	"fmt"

	"github.com/prometheus/client_golang/prometheus"
)

type Prometheus struct {
	counters         map[string]*prometheus.CounterVec
	histograms       map[string]*prometheus.HistogramVec
	gauges           map[string]*prometheus.GaugeVec
	numCounterLabels map[string]int
}

type MetricType int

const (
	Counter MetricType = iota
	Histogram
	Gauge
)

type Metric struct {
	Type        MetricType
	Name        string
	Description string
	Buckets     []float64 // only relevant for the Histogram type
	Labels      []string
}

const (
	duplicatedMetricLabelError   = "duplicate metrics collector registration attempted"
	counterMetricNotFoundError   = "counter metric not found: %s"
	histogramMetricNotFoundError = "histogram metric not found: %s"
	gaugeMetricNotFoundError     = "gauge metric not found: %s"
)

// NewPrometheusService create a new prometheus service
func NewPrometheusService(metrics []Metric) (*Prometheus, error) {
	counters := make(map[string]*prometheus.CounterVec)
	histograms := make(map[string]*prometheus.HistogramVec)
	gauges := make(map[string]*prometheus.GaugeVec)
	numCounterLabels := make(map[string]int)
	// Fill data
	for _, m := range metrics {
		var err error
		switch m.Type {
		case Counter:
			numCounterLabels[m.Name] = len(m.Labels)
			counters[m.Name] = prometheus.NewCounterVec(prometheus.CounterOpts{
				Name: m.Name,
				Help: m.Description,
			}, m.Labels)
			err = prometheus.Register(counters[m.Name])
		case Histogram:
			histograms[m.Name] = prometheus.NewHistogramVec(prometheus.HistogramOpts{
				Name:    m.Name,
				Help:    m.Description,
				Buckets: m.Buckets,
			}, m.Labels)
			err = prometheus.Register(histograms[m.Name])
		case Gauge:
			gauges[m.Name] = prometheus.NewGaugeVec(prometheus.GaugeOpts{
				Name: m.Name,
				Help: m.Description,
			}, m.Labels)
			err = prometheus.Register(gauges[m.Name])
		}
		if err != nil && err.Error() != duplicatedMetricLabelError {
			return nil, err
		}
	}
	return &Prometheus{
		counters:         counters,
		histograms:       histograms,
		gauges:           gauges,
		numCounterLabels: numCounterLabels,
	}, nil
}

func (p *Prometheus) IncreaseCounterMetric(metric string, value float64, labels map[string]string) error {
	if _, ok := p.counters[metric]; !ok {
		return fmt.Errorf(counterMetricNotFoundError, metric)
	}
	p.counters[metric].With(labels).Add(value)
	return nil
}

func (p *Prometheus) ObserveHistogramMetric(metric string, value float64, labels map[string]string) error {
	if _, ok := p.histograms[metric]; !ok {
		return fmt.Errorf(histogramMetricNotFoundError, metric)
	}
	p.histograms[metric].With(labels).Observe(value)
	return nil
}

func (p *Prometheus) SetGaugeMetric(metric string, value float64, labels map[string]string) error {
	if _, ok := p.gauges[metric]; !ok {
		return fmt.Errorf(gaugeMetricNotFoundError, metric)
	}
	p.gauges[metric].With(labels).Set(value)
	return nil
}

// Describe sends the super-set of all possible descriptors of metrics
func (p *Prometheus) Describe(descs chan<- *prometheus.Desc) {
	for _, counter := range p.counters {
		counter.Describe(descs)
	}
	for _, histogram := range p.histograms {
		histogram.Describe(descs)
	}
	for _, gauge := range p.gauges {
		gauge.Describe(descs)
	}
}

// Collect is called by the Prometheus registry when collecting metrics
func (p *Prometheus) Collect(metrics chan<- prometheus.Metric) {
	for _, counter := range p.counters {
		counter.Collect(metrics)
	}
	for _, histogram := range p.histograms {
		histogram.Collect(metrics)
	}
	for _, gauge := range p.gauges {
		gauge.Collect(metrics)
	}
}

const (
	MetricDeploymentAccessTotal       = "deployment_access_total"
	MetricDeploymentAccessFailedTotal = "deployment_access_failed_total"

	LabelID        = "id"
	LabelOperation = "operation"
)

func GetMetricsDefinition() []Metric {
	return []Metric{
		{
			Type:        Counter,
			Name:        MetricDeploymentAccessFailedTotal,
			Description: "External deployment operation failed",
			Labels:      []string{LabelID, LabelOperation},
		},
		{
			Type:        Counter,
			Name:        MetricDeploymentAccessTotal,
			Description: "External deployment operation requested",
			Labels:      []string{LabelID, LabelOperation},
		},
	}
}
