package infrastructure

import (
	"testing"

	"github.com/prometheus/client_golang/prometheus"
	"gosuda.org/boilerplate/internal/domain"
)

// mockStore is a mock implementation of domain.Store for testing
type mockStore struct {
	data map[string]interface{}
}

func newMockStore() *mockStore {
	return &mockStore{
		data: make(map[string]interface{}),
	}
}

func (m *mockStore) Set(key string, value any) error {
	m.data[key] = value
	return nil
}

func (m *mockStore) Get(key string) (value any, err error) {
	if val, exists := m.data[key]; exists {
		return val, nil
	}
	return nil, domain.ErrKeyNotFound
}

func (m *mockStore) GetTyped(key string, value any) error {
	if val, exists := m.data[key]; exists {
		// Simple mock implementation - in real code this would unmarshal
		if str, ok := val.(string); ok {
			if strValue, ok := value.(*string); ok {
				*strValue = str
				return nil
			}
		}
		return nil
	}
	return domain.ErrKeyNotFound
}

func (m *mockStore) List(keyPrefix string) (values []any, err error) {
	var result []any
	for key, value := range m.data {
		if len(key) >= len(keyPrefix) && key[:len(keyPrefix)] == keyPrefix {
			result = append(result, value)
		}
	}
	return result, nil
}

func (m *mockStore) Delete(key string) error {
	if _, exists := m.data[key]; exists {
		delete(m.data, key)
		return nil
	}
	return domain.ErrKeyNotFound
}

func (m *mockStore) Close() error {
	return nil
}

func (m *mockStore) Size() int {
	return len(m.data)
}

func TestNewMetricsStore(t *testing.T) {
	mock := newMockStore()
	reg := prometheus.NewRegistry()
	metrics := NewMetricsCollectorWithRegistry(reg)
	
	ms := NewMetricsStore(mock, metrics)
	if ms == nil {
		t.Fatal("NewMetricsStore returned nil")
	}
}

func TestMetricsStore_Set(t *testing.T) {
	mock := newMockStore()
	reg := prometheus.NewRegistry()
	metrics := NewMetricsCollectorWithRegistry(reg)
	ms := NewMetricsStore(mock, metrics)

	// Test Set operation
	err := ms.Set("test-key", "test-value")
	if err != nil {
		t.Fatalf("Set failed: %v", err)
	}

	// Verify the value was stored
	if val, exists := mock.data["test-key"]; !exists {
		t.Error("Value was not stored in underlying store")
	} else if val != "test-value" {
		t.Errorf("Expected 'test-value', got '%v'", val)
	}
}

func TestMetricsStore_Get(t *testing.T) {
	mock := newMockStore()
	reg := prometheus.NewRegistry()
	metrics := NewMetricsCollectorWithRegistry(reg)
	ms := NewMetricsStore(mock, metrics)

	// Store a value first
	mock.Set("test-key", "test-value")

	// Test Get operation
	value, err := ms.Get("test-key")
	if err != nil {
		t.Fatalf("Get failed: %v", err)
	}

	if value != "test-value" {
		t.Errorf("Expected 'test-value', got '%v'", value)
	}
}

func TestMetricsStore_Get_NotFound(t *testing.T) {
	mock := newMockStore()
	reg := prometheus.NewRegistry()
	metrics := NewMetricsCollectorWithRegistry(reg)
	ms := NewMetricsStore(mock, metrics)

	// Test Get operation with non-existent key
	_, err := ms.Get("non-existent")
	if err != domain.ErrKeyNotFound {
		t.Errorf("Expected ErrKeyNotFound, got %v", err)
	}
}

func TestMetricsStore_GetTyped(t *testing.T) {
	mock := newMockStore()
	reg := prometheus.NewRegistry()
	metrics := NewMetricsCollectorWithRegistry(reg)
	ms := NewMetricsStore(mock, metrics)

	// Store a value first
	mock.Set("test-key", "test-value")

	// Test GetTyped operation
	var value string
	err := ms.GetTyped("test-key", &value)
	if err != nil {
		t.Fatalf("GetTyped failed: %v", err)
	}

	if value != "test-value" {
		t.Errorf("Expected 'test-value', got '%s'", value)
	}
}

func TestMetricsStore_GetTyped_NotFound(t *testing.T) {
	mock := newMockStore()
	reg := prometheus.NewRegistry()
	metrics := NewMetricsCollectorWithRegistry(reg)
	ms := NewMetricsStore(mock, metrics)

	// Test GetTyped operation with non-existent key
	var value string
	err := ms.GetTyped("non-existent", &value)
	if err != domain.ErrKeyNotFound {
		t.Errorf("Expected ErrKeyNotFound, got %v", err)
	}
}

func TestMetricsStore_List(t *testing.T) {
	mock := newMockStore()
	reg := prometheus.NewRegistry()
	metrics := NewMetricsCollectorWithRegistry(reg)
	ms := NewMetricsStore(mock, metrics)

	// Store some values
	mock.Set("posts:1", "post1")
	mock.Set("posts:2", "post2")
	mock.Set("users:1", "user1")

	// Test List operation
	values, err := ms.List("posts:")
	if err != nil {
		t.Fatalf("List failed: %v", err)
	}

	if len(values) != 2 {
		t.Errorf("Expected 2 values, got %d", len(values))
	}

	// Check that we got the right values
	found := make(map[string]bool)
	for _, value := range values {
		found[value.(string)] = true
	}

	if !found["post1"] || !found["post2"] {
		t.Error("Expected to find 'post1' and 'post2' in results")
	}
}

func TestMetricsStore_Delete(t *testing.T) {
	mock := newMockStore()
	reg := prometheus.NewRegistry()
	metrics := NewMetricsCollectorWithRegistry(reg)
	ms := NewMetricsStore(mock, metrics)

	// Store a value first
	mock.Set("test-key", "test-value")

	// Test Delete operation
	err := ms.Delete("test-key")
	if err != nil {
		t.Fatalf("Delete failed: %v", err)
	}

	// Verify the value was deleted
	if _, exists := mock.data["test-key"]; exists {
		t.Error("Value was not deleted from underlying store")
	}
}

func TestMetricsStore_Delete_NotFound(t *testing.T) {
	mock := newMockStore()
	reg := prometheus.NewRegistry()
	metrics := NewMetricsCollectorWithRegistry(reg)
	ms := NewMetricsStore(mock, metrics)

	// Test Delete operation with non-existent key
	err := ms.Delete("non-existent")
	if err != domain.ErrKeyNotFound {
		t.Errorf("Expected ErrKeyNotFound, got %v", err)
	}
}

func TestMetricsStore_Close(t *testing.T) {
	mock := newMockStore()
	reg := prometheus.NewRegistry()
	metrics := NewMetricsCollectorWithRegistry(reg)
	ms := NewMetricsStore(mock, metrics)

	// Test Close operation
	err := ms.Close()
	if err != nil {
		t.Fatalf("Close failed: %v", err)
	}
}

func TestMetricsStore_UpdateStorageMetrics(t *testing.T) {
	mock := newMockStore()
	reg := prometheus.NewRegistry()
	metrics := NewMetricsCollectorWithRegistry(reg)
	ms := NewMetricsStore(mock, metrics)

	// Store some values
	mock.Set("key1", "value1")
	mock.Set("key2", "value2")

	// The updateStorageMetrics method is called internally by Set and Delete
	// We can't easily test it directly, but we can verify it doesn't cause errors
	err := ms.Set("key3", "value3")
	if err != nil {
		t.Fatalf("Set failed: %v", err)
	}

	err = ms.Delete("key1")
	if err != nil {
		t.Fatalf("Delete failed: %v", err)
	}
}