package domain

import (
	"errors"
	"testing"
)

func TestPostNotFoundError(t *testing.T) {
	err := &PostNotFoundError{ID: "test-123"}
	
	if err.Error() == "" {
		t.Error("PostNotFoundError should have a non-empty error message")
	}
	
	if err.ID != "test-123" {
		t.Errorf("Expected ID 'test-123', got '%s'", err.ID)
	}
}

func TestInvalidPostDataError(t *testing.T) {
	err := &InvalidPostDataError{Field: "title", Value: "empty"}
	
	if err.Error() == "" {
		t.Error("InvalidPostDataError should have a non-empty error message")
	}
	
	if err.Field != "title" {
		t.Errorf("Expected Field 'title', got '%s'", err.Field)
	}
	
	if err.Value != "empty" {
		t.Errorf("Expected Value 'empty', got '%s'", err.Value)
	}
}

func TestStorageError(t *testing.T) {
	originalErr := errors.New("database connection failed")
	err := &StorageError{Err: originalErr}
	
	if err.Error() == "" {
		t.Error("StorageError should have a non-empty error message")
	}
	
	if err.Err != originalErr {
		t.Error("StorageError should wrap the original error")
	}
}

func TestValidationError(t *testing.T) {
	err := &ValidationError{Field: "email", Message: "invalid email format"}
	
	if err.Error() == "" {
		t.Error("ValidationError should have a non-empty error message")
	}
	
	if err.Field != "email" {
		t.Errorf("Expected Field 'email', got '%s'", err.Field)
	}
	
	if err.Message != "invalid email format" {
		t.Errorf("Expected Message 'invalid email format', got '%s'", err.Message)
	}
}

func TestPaginationError(t *testing.T) {
	err := &PaginationError{Cursor: "invalid-cursor"}
	
	if err.Error() == "" {
		t.Error("PaginationError should have a non-empty error message")
	}
	
	if err.Cursor != "invalid-cursor" {
		t.Errorf("Expected Cursor 'invalid-cursor', got '%s'", err.Cursor)
	}
}

func TestErrorConstants(t *testing.T) {
	// Test that error constants are defined
	if ErrKeyNotFound == nil {
		t.Error("ErrKeyNotFound should be defined")
	}
	
	// Test error code constants
	expectedCodes := []string{
		ErrorCodePostNotFound,
		ErrorCodeInvalidPostData,
		ErrorCodeValidationError,
		ErrorCodeStorageError,
		ErrorCodePaginationError,
	}
	
	for _, code := range expectedCodes {
		if code == "" {
			t.Errorf("Error code should not be empty")
		}
	}
}