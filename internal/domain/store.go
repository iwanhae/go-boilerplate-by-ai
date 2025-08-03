package domain

// Store defines the interface for persistent data storage
type Store interface {
	// Set stores a value with the given key
	Set(key string, value any) error
	
	// Get retrieves a value by key
	Get(key string) (value any, err error)
	
	// GetTyped retrieves a value by key and unmarshals it into the provided type
	GetTyped(key string, value any) error
	
	// List retrieves all values with keys that start with the given prefix
	List(keyPrefix string) (values []any, err error)
	
	// Delete removes a value by key
	Delete(key string) error
	
	// Close closes the storage and performs cleanup
	Close() error
}