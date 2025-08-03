package middleware

import (
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/prometheus/client_golang/prometheus"
	"gosuda.org/boilerplate/internal/infrastructure"
)

func TestNewMetricsMiddleware(t *testing.T) {
	reg := prometheus.NewRegistry()
	metrics := infrastructure.NewMetricsCollectorWithRegistry(reg)
	mm := NewMetricsMiddleware(metrics)
	
	if mm == nil {
		t.Fatal("NewMetricsMiddleware returned nil")
	}
}

func TestMetricsMiddleware_Handler(t *testing.T) {
	reg := prometheus.NewRegistry()
	metrics := infrastructure.NewMetricsCollectorWithRegistry(reg)
	mm := NewMetricsMiddleware(metrics)

	// Create a test handler
	handler := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusOK)
		w.Write([]byte("test response"))
	})

	// Create a test request
	req := httptest.NewRequest("GET", "/posts", nil)
	w := httptest.NewRecorder()

	// Test the middleware
	middlewareHandler := mm.Handler(handler)
	middlewareHandler.ServeHTTP(w, req)

	// Check response
	if w.Code != http.StatusOK {
		t.Errorf("Expected status code %d, got %d", http.StatusOK, w.Code)
	}

	if w.Body.String() != "test response" {
		t.Errorf("Expected body 'test response', got '%s'", w.Body.String())
	}
}

func TestMetricsMiddleware_Handler_WithError(t *testing.T) {
	reg := prometheus.NewRegistry()
	metrics := infrastructure.NewMetricsCollectorWithRegistry(reg)
	mm := NewMetricsMiddleware(metrics)

	// Create a test handler that returns an error
	handler := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusInternalServerError)
		w.Write([]byte("error response"))
	})

	// Create a test request
	req := httptest.NewRequest("POST", "/posts", nil)
	w := httptest.NewRecorder()

	// Test the middleware
	middlewareHandler := mm.Handler(handler)
	middlewareHandler.ServeHTTP(w, req)

	// Check response
	if w.Code != http.StatusInternalServerError {
		t.Errorf("Expected status code %d, got %d", http.StatusInternalServerError, w.Code)
	}

	if w.Body.String() != "error response" {
		t.Errorf("Expected body 'error response', got '%s'", w.Body.String())
	}
}

func TestMetricsMiddleware_Handler_DifferentMethods(t *testing.T) {
	reg := prometheus.NewRegistry()
	metrics := infrastructure.NewMetricsCollectorWithRegistry(reg)
	mm := NewMetricsMiddleware(metrics)

	// Create a test handler
	handler := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusOK)
	})

	testCases := []struct {
		method string
		path   string
	}{
		{"GET", "/posts"},
		{"POST", "/posts"},
		{"PUT", "/posts/123"},
		{"DELETE", "/posts/123"},
		{"GET", "/debug/metrics"},
	}

	for _, tc := range testCases {
		t.Run(tc.method+"_"+tc.path, func(t *testing.T) {
			req := httptest.NewRequest(tc.method, tc.path, nil)
			w := httptest.NewRecorder()

			middlewareHandler := mm.Handler(handler)
			middlewareHandler.ServeHTTP(w, req)

			if w.Code != http.StatusOK {
				t.Errorf("Expected status code %d, got %d", http.StatusOK, w.Code)
			}
		})
	}
}

func TestMetricsMiddleware_NormalizePath(t *testing.T) {
	reg := prometheus.NewRegistry()
	metrics := infrastructure.NewMetricsCollectorWithRegistry(reg)
	mm := NewMetricsMiddleware(metrics)

	// Test path normalization
	testCases := []struct {
		input    string
		expected string
	}{
		{"/posts", "/posts"},
		{"/posts/123", "/posts/{id}"},
		{"/posts/abc-def", "/posts/{id}"},
		{"/posts/", "/posts/"},
		{"/debug/metrics", "/debug/metrics"},
		{"/posts/123/comments", "/posts/{id}/comments"}, // Should normalize the ID part
	}

	for _, tc := range testCases {
		t.Run(tc.input, func(t *testing.T) {
			result := mm.normalizePath(tc.input)
			if result != tc.expected {
				t.Errorf("Expected '%s', got '%s'", tc.expected, result)
			}
		})
	}
}

func TestMetricsResponseWriter(t *testing.T) {
	// Create a test response writer
	w := httptest.NewRecorder()
	rw := &metricsResponseWriter{
		ResponseWriter: w,
		statusCode:     http.StatusOK,
	}

	// Test WriteHeader
	rw.WriteHeader(http.StatusNotFound)
	if rw.statusCode != http.StatusNotFound {
		t.Errorf("Expected status code %d, got %d", http.StatusNotFound, rw.statusCode)
	}

	// Test Write
	testData := []byte("test data")
	n, err := rw.Write(testData)
	if err != nil {
		t.Fatalf("Write failed: %v", err)
	}

	if n != len(testData) {
		t.Errorf("Expected to write %d bytes, wrote %d", len(testData), n)
	}

	if w.Body.String() != "test data" {
		t.Errorf("Expected body 'test data', got '%s'", w.Body.String())
	}
}

func TestMetricsMiddleware_Integration(t *testing.T) {
	reg := prometheus.NewRegistry()
	metrics := infrastructure.NewMetricsCollectorWithRegistry(reg)
	mm := NewMetricsMiddleware(metrics)

	// Create a test handler that simulates different scenarios
	handler := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		switch r.URL.Path {
		case "/posts/123":
			w.WriteHeader(http.StatusOK)
			w.Write([]byte("post found"))
		case "/posts/456":
			w.WriteHeader(http.StatusNotFound)
			w.Write([]byte("post not found"))
		default:
			w.WriteHeader(http.StatusOK)
			w.Write([]byte("default response"))
		}
	})

	testCases := []struct {
		method       string
		path         string
		expectedCode int
		expectedBody string
	}{
		{"GET", "/posts/123", http.StatusOK, "post found"},
		{"GET", "/posts/456", http.StatusNotFound, "post not found"},
		{"GET", "/posts", http.StatusOK, "default response"},
		{"POST", "/posts", http.StatusOK, "default response"},
	}

	for _, tc := range testCases {
		t.Run(tc.method+"_"+tc.path, func(t *testing.T) {
			req := httptest.NewRequest(tc.method, tc.path, nil)
			w := httptest.NewRecorder()

			middlewareHandler := mm.Handler(handler)
			middlewareHandler.ServeHTTP(w, req)

			if w.Code != tc.expectedCode {
				t.Errorf("Expected status code %d, got %d", tc.expectedCode, w.Code)
			}

			if w.Body.String() != tc.expectedBody {
				t.Errorf("Expected body '%s', got '%s'", tc.expectedBody, w.Body.String())
			}
		})
	}
}