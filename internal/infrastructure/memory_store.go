package infrastructure

import (
	"encoding/json"
	"fmt"
	"strings"
	"sync"

	"gosuda.org/boilerplate/internal/domain"
)

// MemoryStore implements the Store interface using an in-memory map
type MemoryStore struct {
	data    map[string][]byte
	mu      sync.RWMutex
	metrics *Metrics
}

// NewMemoryStore creates a new in-memory store instance
func NewMemoryStore(metrics *Metrics) *MemoryStore {
	return &MemoryStore{
		data:    make(map[string][]byte),
		metrics: metrics,
	}
}

// Set stores a value with the given key
func (s *MemoryStore) Set(key string, value any) error {
	if key == "" {
		return fmt.Errorf("key cannot be empty")
	}

	s.mu.Lock()
	defer s.mu.Unlock()

	// Serialize the value to JSON
	data, err := json.Marshal(value)
	if err != nil {
		return fmt.Errorf("failed to marshal value: %w", err)
	}

	s.data[key] = data
	if s.metrics != nil {
		s.metrics.IncStorageOp("set")
		s.metrics.SetStorageItems(len(s.data))
	}
	return nil
}

// Get retrieves a value by key
func (s *MemoryStore) Get(key string) (value any, err error) {
	if key == "" {
		return nil, fmt.Errorf("key cannot be empty")
	}

	s.mu.RLock()
	defer s.mu.RUnlock()

	data, exists := s.data[key]
	if !exists {
		return nil, domain.ErrKeyNotFound
	}

	// Try to unmarshal as a generic interface{} first
	var result any
	if err := json.Unmarshal(data, &result); err != nil {
		return nil, fmt.Errorf("failed to unmarshal value: %w", err)
	}

	if s.metrics != nil {
		s.metrics.IncStorageOp("get")
	}
	return result, nil
}

// GetTyped retrieves a value by key and unmarshals it into the provided type
func (s *MemoryStore) GetTyped(key string, value any) error {
	if key == "" {
		return fmt.Errorf("key cannot be empty")
	}

	s.mu.RLock()
	defer s.mu.RUnlock()

	data, exists := s.data[key]
	if !exists {
		return domain.ErrKeyNotFound
	}

	if err := json.Unmarshal(data, value); err != nil {
		return fmt.Errorf("failed to unmarshal value: %w", err)
	}

	return nil
}

// List retrieves all values with keys that start with the given prefix
func (s *MemoryStore) List(keyPrefix string) (values []any, err error) {
	s.mu.RLock()
	defer s.mu.RUnlock()

	var result []any
	for key, data := range s.data {
		if strings.HasPrefix(key, keyPrefix) {
			var value any
			if err := json.Unmarshal(data, &value); err != nil {
				return nil, fmt.Errorf("failed to unmarshal value for key %s: %w", key, err)
			}
			result = append(result, value)
		}
	}

	if s.metrics != nil {
		s.metrics.IncStorageOp("list")
	}
	return result, nil
}

// ListKeys retrieves all keys that start with the given prefix
func (s *MemoryStore) ListKeys(keyPrefix string) (keys []string, err error) {
	s.mu.RLock()
	defer s.mu.RUnlock()

	var result []string
	for key := range s.data {
		if strings.HasPrefix(key, keyPrefix) {
			result = append(result, key)
		}
	}

	return result, nil
}

// Delete removes a value by key
func (s *MemoryStore) Delete(key string) error {
	if key == "" {
		return fmt.Errorf("key cannot be empty")
	}

	s.mu.Lock()
	defer s.mu.Unlock()

	if _, exists := s.data[key]; !exists {
		return domain.ErrKeyNotFound
	}

	delete(s.data, key)
	if s.metrics != nil {
		s.metrics.IncStorageOp("delete")
		s.metrics.SetStorageItems(len(s.data))
	}
	return nil
}

// Close closes the storage and performs cleanup
func (s *MemoryStore) Close() error {
	s.mu.Lock()
	defer s.mu.Unlock()

	// Clear all data
	s.data = make(map[string][]byte)
	return nil
}

// Size returns the number of items in the store
func (s *MemoryStore) Size() int {
	s.mu.RLock()
	defer s.mu.RUnlock()
	return len(s.data)
}

// Exists checks if a key exists in the store
func (s *MemoryStore) Exists(key string) bool {
	if key == "" {
		return false
	}

	s.mu.RLock()
	defer s.mu.RUnlock()
	_, exists := s.data[key]
	return exists
}

// Clear removes all data from the store
func (s *MemoryStore) Clear() error {
	s.mu.Lock()
	defer s.mu.Unlock()
	s.data = make(map[string][]byte)
	return nil
}
