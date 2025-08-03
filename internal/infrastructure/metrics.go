package infrastructure

import (
	"context"
	"fmt"
	"sync"
	"time"

	"github.com/prometheus/client_golang/prometheus"
	"github.com/prometheus/client_golang/prometheus/promauto"
)

// MetricsCollector handles application metrics collection
type MetricsCollector struct {
	mu sync.RWMutex

	// HTTP metrics
	httpRequestsTotal   *prometheus.CounterVec
	httpRequestDuration *prometheus.HistogramVec
	httpRequestsInFlight *prometheus.GaugeVec

	// Storage metrics
	storageOperationsTotal   *prometheus.CounterVec
	storageOperationDuration *prometheus.HistogramVec
	storageItemsCurrent      prometheus.Gauge

	// Business metrics
	postsTotal prometheus.Gauge

	// Application metrics
	logLevel prometheus.Gauge

	// Internal state
	currentLogLevel string
}

// NewMetricsCollector creates a new metrics collector
func NewMetricsCollector() *MetricsCollector {
	return NewMetricsCollectorWithRegistry(prometheus.DefaultRegisterer)
}

// NewMetricsCollectorWithRegistry creates a new metrics collector with a custom registry
func NewMetricsCollectorWithRegistry(reg prometheus.Registerer) *MetricsCollector {
	factory := promauto.With(reg)
	
	mc := &MetricsCollector{
		httpRequestsTotal: factory.NewCounterVec(
			prometheus.CounterOpts{
				Name: "http_requests_total",
				Help: "Total number of HTTP requests",
			},
			[]string{"method", "path", "status_code"},
		),
		httpRequestDuration: factory.NewHistogramVec(
			prometheus.HistogramOpts{
				Name:    "http_request_duration_seconds",
				Help:    "HTTP request duration in seconds",
				Buckets: prometheus.DefBuckets,
			},
			[]string{"method", "path"},
		),
		httpRequestsInFlight: factory.NewGaugeVec(
			prometheus.GaugeOpts{
				Name: "http_requests_in_flight",
				Help: "Current number of HTTP requests being processed",
			},
			[]string{"method", "path"},
		),
		storageOperationsTotal: factory.NewCounterVec(
			prometheus.CounterOpts{
				Name: "storage_operations_total",
				Help: "Total number of storage operations",
			},
			[]string{"operation"},
		),
		storageOperationDuration: factory.NewHistogramVec(
			prometheus.HistogramOpts{
				Name:    "storage_operation_duration_seconds",
				Help:    "Storage operation duration in seconds",
				Buckets: prometheus.DefBuckets,
			},
			[]string{"operation"},
		),
		storageItemsCurrent: factory.NewGauge(
			prometheus.GaugeOpts{
				Name: "storage_items_current",
				Help: "Current number of items in storage",
			},
		),
		postsTotal: factory.NewGauge(
			prometheus.GaugeOpts{
				Name: "posts_total",
				Help: "Total number of posts",
			},
		),
		logLevel: factory.NewGauge(
			prometheus.GaugeOpts{
				Name: "log_level",
				Help: "Current log level (0=debug, 1=info, 2=warn, 3=error)",
			},
		),
	}

	// Initialize log level
	mc.SetLogLevel("info")

	return mc
}

// RecordHTTPRequest records an HTTP request
func (mc *MetricsCollector) RecordHTTPRequest(method, path string, statusCode int, duration time.Duration) {
	mc.httpRequestsTotal.WithLabelValues(method, path, fmt.Sprintf("%d", statusCode)).Inc()
	mc.httpRequestDuration.WithLabelValues(method, path).Observe(duration.Seconds())
}

// RecordHTTPRequestStart records the start of an HTTP request
func (mc *MetricsCollector) RecordHTTPRequestStart(method, path string) {
	mc.httpRequestsInFlight.WithLabelValues(method, path).Inc()
}

// RecordHTTPRequestEnd records the end of an HTTP request
func (mc *MetricsCollector) RecordHTTPRequestEnd(method, path string) {
	mc.httpRequestsInFlight.WithLabelValues(method, path).Dec()
}

// RecordStorageOperation records a storage operation
func (mc *MetricsCollector) RecordStorageOperation(operation string, duration time.Duration) {
	mc.storageOperationsTotal.WithLabelValues(operation).Inc()
	mc.storageOperationDuration.WithLabelValues(operation).Observe(duration.Seconds())
}

// SetStorageItemsCount sets the current number of items in storage
func (mc *MetricsCollector) SetStorageItemsCount(count int) {
	mc.storageItemsCurrent.Set(float64(count))
}

// SetPostsCount sets the current number of posts
func (mc *MetricsCollector) SetPostsCount(count int) {
	mc.postsTotal.Set(float64(count))
}

// SetLogLevel sets the current log level
func (mc *MetricsCollector) SetLogLevel(level string) {
	mc.mu.Lock()
	defer mc.mu.Unlock()

	// Validate and normalize log level
	var normalizedLevel string
	var levelValue float64
	switch level {
	case "debug":
		normalizedLevel = "debug"
		levelValue = 0
	case "info":
		normalizedLevel = "info"
		levelValue = 1
	case "warn":
		normalizedLevel = "warn"
		levelValue = 2
	case "error":
		normalizedLevel = "error"
		levelValue = 3
	default:
		normalizedLevel = "info" // default to info
		levelValue = 1
	}

	mc.currentLogLevel = normalizedLevel
	mc.logLevel.Set(levelValue)
}

// GetLogLevel returns the current log level
func (mc *MetricsCollector) GetLogLevel() string {
	mc.mu.RLock()
	defer mc.mu.RUnlock()
	return mc.currentLogLevel
}

// GetMetrics returns the current metrics as a string
func (mc *MetricsCollector) GetMetrics(ctx context.Context) (string, error) {
	// In a real implementation, this would use prometheus.Registry
	// For now, we'll return a simple metrics format
	metrics := fmt.Sprintf(`# HELP http_requests_total Total number of HTTP requests
# TYPE http_requests_total counter
http_requests_total{method="GET",path="/posts",status_code="200"} 0
http_requests_total{method="POST",path="/posts",status_code="201"} 0
http_requests_total{method="GET",path="/posts/{id}",status_code="200"} 0
http_requests_total{method="PUT",path="/posts/{id}",status_code="200"} 0
http_requests_total{method="DELETE",path="/posts/{id}",status_code="204"} 0

# HELP http_request_duration_seconds HTTP request duration in seconds
# TYPE http_request_duration_seconds histogram
http_request_duration_seconds_bucket{method="GET",path="/posts",le="0.1"} 0
http_request_duration_seconds_bucket{method="GET",path="/posts",le="0.5"} 0
http_request_duration_seconds_bucket{method="GET",path="/posts",le="1"} 0
http_request_duration_seconds_bucket{method="GET",path="/posts",le="+Inf"} 0
http_request_duration_seconds_sum{method="GET",path="/posts"} 0
http_request_duration_seconds_count{method="GET",path="/posts"} 0

# HELP http_requests_in_flight Current number of HTTP requests being processed
# TYPE http_requests_in_flight gauge
http_requests_in_flight{method="GET",path="/posts"} 0

# HELP storage_operations_total Total number of storage operations
# TYPE storage_operations_total counter
storage_operations_total{operation="get"} 0
storage_operations_total{operation="set"} 0
storage_operations_total{operation="delete"} 0
storage_operations_total{operation="list"} 0

# HELP storage_operation_duration_seconds Storage operation duration in seconds
# TYPE storage_operation_duration_seconds histogram
storage_operation_duration_seconds_bucket{operation="get",le="0.001"} 0
storage_operation_duration_seconds_bucket{operation="get",le="0.01"} 0
storage_operation_duration_seconds_bucket{operation="get",le="0.1"} 0
storage_operation_duration_seconds_bucket{operation="get",le="+Inf"} 0
storage_operation_duration_seconds_sum{operation="get"} 0
storage_operation_duration_seconds_count{operation="get"} 0

# HELP storage_items_current Current number of items in storage
# TYPE storage_items_current gauge
storage_items_current 0

# HELP posts_total Total number of posts
# TYPE posts_total gauge
posts_total 0

# HELP log_level Current log level
# TYPE log_level gauge
log_level{level="%s"} 1
`, mc.GetLogLevel())

	return metrics, nil
}