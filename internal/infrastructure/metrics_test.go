package infrastructure

import (
	"strings"
	"testing"
)

func TestMetrics(t *testing.T) {
	m := NewMetrics()
	m.IncRequest("GET", "/posts")
	m.IncStorageOp("set")
	m.SetStorageItems(1)

	metrics, err := m.Gather()
	if err != nil {
		t.Fatalf("gather metrics: %v", err)
	}

	if !strings.Contains(metrics, "app_requests_total{method=\"GET\",path=\"/posts\"} 1") {
		t.Errorf("request metric not found: %s", metrics)
	}
	if !strings.Contains(metrics, "app_storage_operations_total{operation=\"set\"} 1") {
		t.Errorf("storage metric not found: %s", metrics)
	}
	if !strings.Contains(metrics, "app_storage_items_current 1") {
		t.Errorf("storage items metric not found: %s", metrics)
	}
}
