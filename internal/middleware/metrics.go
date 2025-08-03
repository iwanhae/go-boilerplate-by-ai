package middleware

import (
	"net/http"
	"strings"
	"time"

	"gosuda.org/boilerplate/internal/infrastructure"
)

// MetricsMiddleware tracks HTTP request metrics
type MetricsMiddleware struct {
	metrics *infrastructure.MetricsCollector
}

// NewMetricsMiddleware creates a new metrics middleware
func NewMetricsMiddleware(metrics *infrastructure.MetricsCollector) *MetricsMiddleware {
	return &MetricsMiddleware{
		metrics: metrics,
	}
}

// Handler is the middleware function
func (m *MetricsMiddleware) Handler(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		start := time.Now()

		// Extract method and path for metrics
		method := r.Method
		path := m.normalizePath(r.URL.Path)

		// Record request start
		m.metrics.RecordHTTPRequestStart(method, path)

		// Create a response writer wrapper to capture status code
		wrappedWriter := &metricsResponseWriter{ResponseWriter: w, statusCode: http.StatusOK}

		// Call next handler
		next.ServeHTTP(wrappedWriter, r)

		// Record request end
		m.metrics.RecordHTTPRequestEnd(method, path)

		// Record request metrics
		duration := time.Since(start)
		m.metrics.RecordHTTPRequest(method, path, wrappedWriter.statusCode, duration)
	})
}

// normalizePath normalizes the path for metrics (e.g., /posts/123 -> /posts/{id})
func (m *MetricsMiddleware) normalizePath(path string) string {
	// Normalize /posts/{id} patterns
	if strings.HasPrefix(path, "/posts/") && len(path) > 7 {
		// Find the next slash after /posts/
		rest := path[7:]
		nextSlash := strings.Index(rest, "/")
		
		if nextSlash == -1 {
			// No more slashes, this is /posts/{id}
			return "/posts/{id}"
		} else {
			// There are more slashes, this is /posts/{id}/something
			return "/posts/{id}" + rest[nextSlash:]
		}
	}
	return path
}

// metricsResponseWriter wraps http.ResponseWriter to capture status code
type metricsResponseWriter struct {
	http.ResponseWriter
	statusCode int
}

func (rw *metricsResponseWriter) WriteHeader(code int) {
	rw.statusCode = code
	rw.ResponseWriter.WriteHeader(code)
}

func (rw *metricsResponseWriter) Write(b []byte) (int, error) {
	return rw.ResponseWriter.Write(b)
}