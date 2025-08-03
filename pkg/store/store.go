package store

import "errors"

// Store abstracts persistence operations.
type Store interface {
	Set(key string, value any) error
	Get(key string) (any, error)
	List(prefix string) ([]any, error)
	Delete(key string) error
	Close() error
}

var ErrNotFound = errors.New("not found")
