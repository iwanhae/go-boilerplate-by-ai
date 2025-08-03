package infrastructure

import (
	"time"

	"gosuda.org/boilerplate/internal/domain"
)

// MetricsStore wraps a Store implementation to add metrics tracking
type MetricsStore struct {
	store   domain.Store
	metrics *MetricsCollector
}

// NewMetricsStore creates a new metrics-aware store wrapper
func NewMetricsStore(store domain.Store, metrics *MetricsCollector) *MetricsStore {
	return &MetricsStore{
		store:   store,
		metrics: metrics,
	}
}

// Set stores a value with the given key
func (s *MetricsStore) Set(key string, value any) error {
	start := time.Now()
	err := s.store.Set(key, value)
	duration := time.Since(start)

	s.metrics.RecordStorageOperation("set", duration)
	s.updateStorageMetrics()

	return err
}

// Get retrieves a value by key
func (s *MetricsStore) Get(key string) (value any, err error) {
	start := time.Now()
	value, err = s.store.Get(key)
	duration := time.Since(start)

	s.metrics.RecordStorageOperation("get", duration)

	return value, err
}

// GetTyped retrieves a value by key and unmarshals it into the provided type
func (s *MetricsStore) GetTyped(key string, value any) error {
	start := time.Now()
	err := s.store.GetTyped(key, value)
	duration := time.Since(start)

	s.metrics.RecordStorageOperation("get", duration)

	return err
}

// List retrieves all values with keys that start with the given prefix
func (s *MetricsStore) List(keyPrefix string) (values []any, err error) {
	start := time.Now()
	values, err = s.store.List(keyPrefix)
	duration := time.Since(start)

	s.metrics.RecordStorageOperation("list", duration)

	return values, err
}

// Delete removes a value by key
func (s *MetricsStore) Delete(key string) error {
	start := time.Now()
	err := s.store.Delete(key)
	duration := time.Since(start)

	s.metrics.RecordStorageOperation("delete", duration)
	s.updateStorageMetrics()

	return err
}

// Close closes the storage and performs cleanup
func (s *MetricsStore) Close() error {
	return s.store.Close()
}

// updateStorageMetrics updates the storage metrics
func (s *MetricsStore) updateStorageMetrics() {
	// Try to get the size if the store supports it
	if sizeable, ok := s.store.(interface{ Size() int }); ok {
		size := sizeable.Size()
		s.metrics.SetStorageItemsCount(size)
	}
}