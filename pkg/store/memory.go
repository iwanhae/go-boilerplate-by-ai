package store

import "sync"

// Memory implements Store using an in-memory map.
type Memory struct {
	mu sync.RWMutex
	m  map[string]any
}

// NewMemory creates a new Memory store.
func NewMemory() *Memory {
	return &Memory{m: make(map[string]any)}
}

func (s *Memory) Set(key string, value any) error {
	s.mu.Lock()
	defer s.mu.Unlock()
	s.m[key] = value
	return nil
}

func (s *Memory) Get(key string) (any, error) {
	s.mu.RLock()
	defer s.mu.RUnlock()
	v, ok := s.m[key]
	if !ok {
		return nil, ErrNotFound
	}
	return v, nil
}

func (s *Memory) List(prefix string) ([]any, error) {
	s.mu.RLock()
	defer s.mu.RUnlock()
	out := make([]any, 0)
	for k, v := range s.m {
		if len(prefix) == 0 || (len(k) >= len(prefix) && k[:len(prefix)] == prefix) {
			out = append(out, v)
		}
	}
	return out, nil
}

func (s *Memory) Delete(key string) error {
	s.mu.Lock()
	defer s.mu.Unlock()
	if _, ok := s.m[key]; !ok {
		return ErrNotFound
	}
	delete(s.m, key)
	return nil
}

func (s *Memory) Close() error { return nil }
