package domain

import "errors"

// Domain errors
var (
	ErrKeyNotFound = errors.New("key not found")
)

// PostNotFoundError represents when a post is not found
type PostNotFoundError struct {
	ID string
}

func (e PostNotFoundError) Error() string {
	return "post not found: " + e.ID
}

// InvalidPostDataError represents invalid post data
type InvalidPostDataError struct {
	Field string
	Value string
}

func (e InvalidPostDataError) Error() string {
	return "invalid post data: " + e.Field + " = " + e.Value
}

// StorageError represents storage operation errors
type StorageError struct {
	Err error
}

func (e StorageError) Error() string {
	return "storage error: " + e.Err.Error()
}

func (e StorageError) Unwrap() error {
	return e.Err
}

// ValidationError represents validation errors
type ValidationError struct {
	Field   string
	Message string
}

func (e ValidationError) Error() string {
	return "validation error: " + e.Field + " - " + e.Message
}

// PaginationError represents pagination errors
type PaginationError struct {
	Cursor string
}

func (e PaginationError) Error() string {
	return "invalid pagination cursor: " + e.Cursor
}

// Error codes for HTTP responses
const (
	ErrorCodePostNotFound     = "POST_NOT_FOUND"
	ErrorCodeInvalidPostData  = "INVALID_POST_DATA"
	ErrorCodeStorageError     = "STORAGE_ERROR"
	ErrorCodeValidationError  = "VALIDATION_ERROR"
	ErrorCodePaginationError  = "PAGINATION_ERROR"
	ErrorCodeInternalError    = "INTERNAL_ERROR"
)