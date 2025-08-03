package infrastructure

import (
	"bytes"
	"net/http"

	"github.com/prometheus/client_golang/prometheus"
	"github.com/prometheus/client_golang/prometheus/promhttp"
	"github.com/prometheus/common/expfmt"
)

// Metrics provides Prometheus metrics collection and exposition.
type Metrics struct {
	registry          *prometheus.Registry
	requestsTotal     *prometheus.CounterVec
	storageOpsTotal   *prometheus.CounterVec
	storageItemsGauge prometheus.Gauge
}

// NewMetrics creates a new Metrics instance with pre-defined collectors.
func NewMetrics() *Metrics {
	r := prometheus.NewRegistry()
	m := &Metrics{
		registry: r,
		requestsTotal: prometheus.NewCounterVec(
			prometheus.CounterOpts{
				Name: "app_requests_total",
				Help: "Total number of HTTP requests",
			},
			[]string{"method", "path"},
		),
		storageOpsTotal: prometheus.NewCounterVec(
			prometheus.CounterOpts{
				Name: "app_storage_operations_total",
				Help: "Total number of storage operations",
			},
			[]string{"operation"},
		),
		storageItemsGauge: prometheus.NewGauge(
			prometheus.GaugeOpts{
				Name: "app_storage_items_current",
				Help: "Current number of items in storage",
			},
		),
	}
	r.MustRegister(m.requestsTotal, m.storageOpsTotal, m.storageItemsGauge)
	return m
}

// Handler returns an HTTP handler for serving metrics in Prometheus format.
func (m *Metrics) Handler() http.Handler {
	return promhttp.HandlerFor(m.registry, promhttp.HandlerOpts{})
}

// Gather returns the metrics exposition format as a string.
func (m *Metrics) Gather() (string, error) {
	mfs, err := m.registry.Gather()
	if err != nil {
		return "", err
	}
	var buf bytes.Buffer
	for _, mf := range mfs {
		if _, err := expfmt.MetricFamilyToText(&buf, mf); err != nil {
			return "", err
		}
	}
	return buf.String(), nil
}

// IncRequest increments the HTTP request counter.
func (m *Metrics) IncRequest(method, path string) {
	m.requestsTotal.WithLabelValues(method, path).Inc()
}

// IncStorageOp increments the storage operation counter.
func (m *Metrics) IncStorageOp(op string) {
	m.storageOpsTotal.WithLabelValues(op).Inc()
}

// SetStorageItems sets the current number of items in storage.
func (m *Metrics) SetStorageItems(n int) {
	m.storageItemsGauge.Set(float64(n))
}
