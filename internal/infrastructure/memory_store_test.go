package infrastructure

import (
	"fmt"
	"testing"

	"gosuda.org/boilerplate/internal/domain"
)

func TestNewMemoryStore(t *testing.T) {
	store := NewMemoryStore()
	if store == nil {
		t.Fatal("Expected store but got nil")
	}

	if store.data == nil {
		t.Error("Expected data map to be initialized")
	}

	if store.Size() != 0 {
		t.Errorf("Expected empty store, got size %d", store.Size())
	}
}

func TestMemoryStore_Set(t *testing.T) {
	store := NewMemoryStore()

	// Test setting a simple value
	err := store.Set("key1", "value1")
	if err != nil {
		t.Errorf("Failed to set value: %v", err)
	}

	if !store.Exists("key1") {
		t.Error("Key should exist after setting")
	}

	// Test setting a complex value
	complexValue := map[string]interface{}{
		"name": "test",
		"age":  25,
		"tags": []string{"tag1", "tag2"},
	}

	err = store.Set("key2", complexValue)
	if err != nil {
		t.Errorf("Failed to set complex value: %v", err)
	}

	// Test setting with empty key
	err = store.Set("", "value")
	if err == nil {
		t.Error("Expected error for empty key")
	}
}

func TestMemoryStore_Get(t *testing.T) {
	store := NewMemoryStore()

	// Test getting non-existent key
	_, err := store.Get("nonexistent")
	if err != domain.ErrKeyNotFound {
		t.Errorf("Expected ErrKeyNotFound, got %v", err)
	}

	// Test getting with empty key
	_, err = store.Get("")
	if err == nil {
		t.Error("Expected error for empty key")
	}

	// Test getting existing key
	store.Set("key1", "value1")
	value, err := store.Get("key1")
	if err != nil {
		t.Errorf("Failed to get value: %v", err)
	}

	if value != "value1" {
		t.Errorf("Expected 'value1', got %v", value)
	}

	// Test getting complex value
	complexValue := map[string]interface{}{
		"name": "test",
		"age":  25,
	}
	store.Set("key2", complexValue)

	retrieved, err := store.Get("key2")
	if err != nil {
		t.Errorf("Failed to get complex value: %v", err)
	}

	retrievedMap, ok := retrieved.(map[string]interface{})
	if !ok {
		t.Fatal("Expected map[string]interface{}")
	}

	if retrievedMap["name"] != "test" {
		t.Errorf("Expected 'test', got %v", retrievedMap["name"])
	}
}

func TestMemoryStore_GetTyped(t *testing.T) {
	store := NewMemoryStore()

	// Test getting typed value
	originalValue := map[string]interface{}{
		"name": "test",
		"age":  25,
	}
	store.Set("key1", originalValue)

	var retrievedValue map[string]interface{}
	err := store.GetTyped("key1", &retrievedValue)
	if err != nil {
		t.Errorf("Failed to get typed value: %v", err)
	}

	if retrievedValue["name"] != "test" {
		t.Errorf("Expected 'test', got %v", retrievedValue["name"])
	}

	// Test getting non-existent key
	err = store.GetTyped("nonexistent", &retrievedValue)
	if err != domain.ErrKeyNotFound {
		t.Errorf("Expected ErrKeyNotFound, got %v", err)
	}

	// Test getting with empty key
	err = store.GetTyped("", &retrievedValue)
	if err == nil {
		t.Error("Expected error for empty key")
	}
}

func TestMemoryStore_List(t *testing.T) {
	store := NewMemoryStore()

	// Set up test data
	store.Set("user:1", "user1")
	store.Set("user:2", "user2")
	store.Set("post:1", "post1")
	store.Set("post:2", "post2")

	// Test listing with prefix
	users, err := store.List("user:")
	if err != nil {
		t.Errorf("Failed to list users: %v", err)
	}

	if len(users) != 2 {
		t.Errorf("Expected 2 users, got %d", len(users))
	}

	posts, err := store.List("post:")
	if err != nil {
		t.Errorf("Failed to list posts: %v", err)
	}

	if len(posts) != 2 {
		t.Errorf("Expected 2 posts, got %d", len(posts))
	}

	// Test listing with non-existent prefix
	empty, err := store.List("nonexistent:")
	if err != nil {
		t.Errorf("Failed to list with non-existent prefix: %v", err)
	}

	if len(empty) != 0 {
		t.Errorf("Expected empty list, got %d items", len(empty))
	}

	// Test listing all
	all, err := store.List("")
	if err != nil {
		t.Errorf("Failed to list all: %v", err)
	}

	if len(all) != 4 {
		t.Errorf("Expected 4 items, got %d", len(all))
	}
}

func TestMemoryStore_ListKeys(t *testing.T) {
	store := NewMemoryStore()

	// Set up test data
	store.Set("user:1", "user1")
	store.Set("user:2", "user2")
	store.Set("post:1", "post1")

	// Test listing keys with prefix
	userKeys, err := store.ListKeys("user:")
	if err != nil {
		t.Errorf("Failed to list user keys: %v", err)
	}

	if len(userKeys) != 2 {
		t.Errorf("Expected 2 user keys, got %d", len(userKeys))
	}

	// Verify keys are correct
	expectedKeys := map[string]bool{"user:1": true, "user:2": true}
	for _, key := range userKeys {
		if !expectedKeys[key] {
			t.Errorf("Unexpected key: %s", key)
		}
	}
}

func TestMemoryStore_Delete(t *testing.T) {
	store := NewMemoryStore()

	// Test deleting non-existent key
	err := store.Delete("nonexistent")
	if err != domain.ErrKeyNotFound {
		t.Errorf("Expected ErrKeyNotFound, got %v", err)
	}

	// Test deleting with empty key
	err = store.Delete("")
	if err == nil {
		t.Error("Expected error for empty key")
	}

	// Test deleting existing key
	store.Set("key1", "value1")
	if !store.Exists("key1") {
		t.Error("Key should exist before deletion")
	}

	err = store.Delete("key1")
	if err != nil {
		t.Errorf("Failed to delete key: %v", err)
	}

	if store.Exists("key1") {
		t.Error("Key should not exist after deletion")
	}
}

func TestMemoryStore_Close(t *testing.T) {
	store := NewMemoryStore()

	// Add some data
	store.Set("key1", "value1")
	store.Set("key2", "value2")

	if store.Size() != 2 {
		t.Errorf("Expected size 2, got %d", store.Size())
	}

	// Close the store
	err := store.Close()
	if err != nil {
		t.Errorf("Failed to close store: %v", err)
	}

	// Verify data is cleared
	if store.Size() != 0 {
		t.Errorf("Expected size 0 after close, got %d", store.Size())
	}
}

func TestMemoryStore_Size(t *testing.T) {
	store := NewMemoryStore()

	if store.Size() != 0 {
		t.Errorf("Expected size 0, got %d", store.Size())
	}

	store.Set("key1", "value1")
	if store.Size() != 1 {
		t.Errorf("Expected size 1, got %d", store.Size())
	}

	store.Set("key2", "value2")
	if store.Size() != 2 {
		t.Errorf("Expected size 2, got %d", store.Size())
	}

	store.Delete("key1")
	if store.Size() != 1 {
		t.Errorf("Expected size 1 after deletion, got %d", store.Size())
	}
}

func TestMemoryStore_Exists(t *testing.T) {
	store := NewMemoryStore()

	// Test non-existent key
	if store.Exists("nonexistent") {
		t.Error("Non-existent key should return false")
	}

	// Test empty key
	if store.Exists("") {
		t.Error("Empty key should return false")
	}

	// Test existing key
	store.Set("key1", "value1")
	if !store.Exists("key1") {
		t.Error("Existing key should return true")
	}

	// Test after deletion
	store.Delete("key1")
	if store.Exists("key1") {
		t.Error("Deleted key should return false")
	}
}

func TestMemoryStore_Clear(t *testing.T) {
	store := NewMemoryStore()

	// Add some data
	store.Set("key1", "value1")
	store.Set("key2", "value2")

	if store.Size() != 2 {
		t.Errorf("Expected size 2, got %d", store.Size())
	}

	// Clear the store
	err := store.Clear()
	if err != nil {
		t.Errorf("Failed to clear store: %v", err)
	}

	// Verify data is cleared
	if store.Size() != 0 {
		t.Errorf("Expected size 0 after clear, got %d", store.Size())
	}
}

func TestMemoryStore_Concurrency(t *testing.T) {
	store := NewMemoryStore()
	done := make(chan bool, 10)

	// Start multiple goroutines writing
	for i := 0; i < 5; i++ {
		go func(id int) {
			for j := 0; j < 10; j++ {
				key := fmt.Sprintf("key-%d-%d", id, j)
				store.Set(key, fmt.Sprintf("value-%d-%d", id, j))
			}
			done <- true
		}(i)
	}

	// Start multiple goroutines reading
	for i := 0; i < 5; i++ {
		go func(id int) {
			for j := 0; j < 10; j++ {
				key := fmt.Sprintf("key-%d-%d", id, j)
				store.Get(key)
			}
			done <- true
		}(i)
	}

	// Wait for all goroutines to complete
	for i := 0; i < 10; i++ {
		<-done
	}

	// Verify final state
	expectedSize := 50 // 5 goroutines * 10 keys each
	if store.Size() != expectedSize {
		t.Errorf("Expected size %d, got %d", expectedSize, store.Size())
	}
}

func TestMemoryStore_JSONSerialization(t *testing.T) {
	store := NewMemoryStore()

	// Test storing various data types
	testCases := []struct {
		name  string
		key   string
		value interface{}
	}{
		{"string", "str", "hello world"},
		{"int", "int", 42},
		{"float", "float", 3.14},
		{"bool", "bool", true},
		{"slice", "slice", []string{"a", "b", "c"}},
		{"map", "map", map[string]int{"a": 1, "b": 2}},
		{"struct", "struct", struct {
			Name string `json:"name"`
			Age  int    `json:"age"`
		}{Name: "test", Age: 25}},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			err := store.Set(tc.key, tc.value)
			if err != nil {
				t.Errorf("Failed to set %s: %v", tc.name, err)
			}

			retrieved, err := store.Get(tc.key)
			if err != nil {
				t.Errorf("Failed to get %s: %v", tc.name, err)
			}

			// Basic type checking
			if retrieved == nil {
				t.Errorf("Retrieved value is nil for %s", tc.name)
			}
		})
	}
}