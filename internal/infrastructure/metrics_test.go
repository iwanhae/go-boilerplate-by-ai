package infrastructure

import (
	"context"
	"strings"
	"testing"
	"time"

	"github.com/prometheus/client_golang/prometheus"
)

func TestNewMetricsCollector(t *testing.T) {
	reg := prometheus.NewRegistry()
	mc := NewMetricsCollectorWithRegistry(reg)
	if mc == nil {
		t.Fatal("NewMetricsCollector returned nil")
	}

	// Check initial log level
	if mc.GetLogLevel() != "info" {
		t.Errorf("Expected initial log level to be 'info', got '%s'", mc.GetLogLevel())
	}
}

func TestMetricsCollector_SetLogLevel(t *testing.T) {
	reg := prometheus.NewRegistry()
	mc := NewMetricsCollectorWithRegistry(reg)

	testCases := []struct {
		level      string
		expectPass bool
	}{
		{"debug", true},
		{"info", true},
		{"warn", true},
		{"error", true},
		{"invalid", true}, // Should default to info
	}

	for _, tc := range testCases {
		t.Run(tc.level, func(t *testing.T) {
			mc.SetLogLevel(tc.level)
			result := mc.GetLogLevel()
			if tc.expectPass {
				if tc.level == "invalid" {
					if result != "info" {
						t.Errorf("Expected 'invalid' to default to 'info', got '%s'", result)
					}
				} else if result != tc.level {
					t.Errorf("Expected log level '%s', got '%s'", tc.level, result)
				}
			}
		})
	}
}

func TestMetricsCollector_RecordHTTPRequest(t *testing.T) {
	reg := prometheus.NewRegistry()
	mc := NewMetricsCollectorWithRegistry(reg)

	// Test recording HTTP request
	mc.RecordHTTPRequest("GET", "/posts", 200, 100*time.Millisecond)
	mc.RecordHTTPRequest("POST", "/posts", 201, 150*time.Millisecond)

	// The metrics are recorded but we can't easily verify them without exposing internal state
	// In a real implementation, you might want to expose methods to query the metrics
}

func TestMetricsCollector_RecordHTTPRequestStartEnd(t *testing.T) {
	reg := prometheus.NewRegistry()
	mc := NewMetricsCollectorWithRegistry(reg)

	// Test recording request start/end
	mc.RecordHTTPRequestStart("GET", "/posts")
	mc.RecordHTTPRequestEnd("GET", "/posts")

	// The metrics are recorded but we can't easily verify them without exposing internal state
}

func TestMetricsCollector_RecordStorageOperation(t *testing.T) {
	reg := prometheus.NewRegistry()
	mc := NewMetricsCollectorWithRegistry(reg)

	// Test recording storage operations
	mc.RecordStorageOperation("get", 10*time.Millisecond)
	mc.RecordStorageOperation("set", 20*time.Millisecond)
	mc.RecordStorageOperation("delete", 15*time.Millisecond)
	mc.RecordStorageOperation("list", 25*time.Millisecond)

	// The metrics are recorded but we can't easily verify them without exposing internal state
}

func TestMetricsCollector_SetStorageItemsCount(t *testing.T) {
	reg := prometheus.NewRegistry()
	mc := NewMetricsCollectorWithRegistry(reg)

	// Test setting storage items count
	mc.SetStorageItemsCount(10)
	mc.SetStorageItemsCount(25)
	mc.SetStorageItemsCount(0)

	// The metrics are recorded but we can't easily verify them without exposing internal state
}

func TestMetricsCollector_SetPostsCount(t *testing.T) {
	reg := prometheus.NewRegistry()
	mc := NewMetricsCollectorWithRegistry(reg)

	// Test setting posts count
	mc.SetPostsCount(5)
	mc.SetPostsCount(15)
	mc.SetPostsCount(0)

	// The metrics are recorded but we can't easily verify them without exposing internal state
}

func TestMetricsCollector_GetMetrics(t *testing.T) {
	reg := prometheus.NewRegistry()
	mc := NewMetricsCollectorWithRegistry(reg)
	mc.SetLogLevel("debug")

	metrics, err := mc.GetMetrics(context.Background())
	if err != nil {
		t.Fatalf("GetMetrics failed: %v", err)
	}

	if metrics == "" {
		t.Fatal("GetMetrics returned empty string")
	}

	// Check that the metrics contain expected content
	if !strings.Contains(metrics, "http_requests_total") {
		t.Error("Metrics should contain http_requests_total")
	}

	if !strings.Contains(metrics, "storage_operations_total") {
		t.Error("Metrics should contain storage_operations_total")
	}

	if !strings.Contains(metrics, "posts_total") {
		t.Error("Metrics should contain posts_total")
	}

	if !strings.Contains(metrics, "log_level") {
		t.Error("Metrics should contain log_level")
	}

	// Check that the log level is included
	if !strings.Contains(metrics, "debug") {
		t.Error("Metrics should contain the current log level")
	}
}