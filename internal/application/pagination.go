package application

import (
	"encoding/base64"
	"encoding/json"
	"fmt"
	"strconv"

	"gosuda.org/boilerplate/internal/domain"
)

// PaginationParams represents pagination parameters
type PaginationParams struct {
	Cursor string `json:"cursor,omitempty"`
	Limit  int    `json:"limit"`
}

// PaginationResult represents pagination result
type PaginationResult struct {
	Items      []interface{} `json:"items"`
	NextCursor string        `json:"nextCursor,omitempty"`
	HasMore    bool          `json:"hasMore"`
}

// Pagination constants
const (
	DefaultLimit = 20
	MinLimit     = 1
	MaxLimit     = 100
)

// Cursor represents a pagination cursor
type Cursor struct {
	ID    string `json:"id"`
	Limit int    `json:"limit"`
}

// NewPaginationParams creates new pagination parameters with defaults
func NewPaginationParams(cursor string, limit int) *PaginationParams {
	if limit <= 0 {
		limit = DefaultLimit
	}
	if limit > MaxLimit {
		limit = MaxLimit
	}
	if limit < MinLimit {
		limit = MinLimit
	}

	return &PaginationParams{
		Cursor: cursor,
		Limit:  limit,
	}
}

// DecodeCursor decodes a cursor string into a Cursor struct
func DecodeCursor(cursorStr string) (*Cursor, error) {
	if cursorStr == "" {
		return nil, nil
	}

	decoded, err := base64.URLEncoding.DecodeString(cursorStr)
	if err != nil {
		return nil, &domain.PaginationError{Cursor: cursorStr}
	}

	var cursor Cursor
	if err := json.Unmarshal(decoded, &cursor); err != nil {
		return nil, &domain.PaginationError{Cursor: cursorStr}
	}

	return &cursor, nil
}

// EncodeCursor encodes a Cursor struct into a string
func EncodeCursor(cursor *Cursor) (string, error) {
	if cursor == nil {
		return "", nil
	}

	data, err := json.Marshal(cursor)
	if err != nil {
		return "", err
	}

	return base64.URLEncoding.EncodeToString(data), nil
}

// ParseLimit parses and validates a limit string
func ParseLimit(limitStr string) (int, error) {
	if limitStr == "" {
		return DefaultLimit, nil
	}

	limit, err := strconv.Atoi(limitStr)
	if err != nil {
		return 0, fmt.Errorf("invalid limit: %s", limitStr)
	}

	if limit < MinLimit {
		return 0, fmt.Errorf("limit too small: %d", limit)
	}

	if limit > MaxLimit {
		return 0, fmt.Errorf("limit too large: %d", limit)
	}

	return limit, nil
}

// ValidatePaginationParams validates pagination parameters
func ValidatePaginationParams(cursor string, limit int) error {
	if limit < MinLimit || limit > MaxLimit {
		return fmt.Errorf("limit must be between %d and %d", MinLimit, MaxLimit)
	}

	if cursor != "" {
		if _, err := DecodeCursor(cursor); err != nil {
			return err
		}
	}

	return nil
}

// CreateNextCursor creates a next cursor for pagination
func CreateNextCursor(lastID string, limit int) (string, error) {
	if lastID == "" {
		return "", nil
	}

	cursor := &Cursor{
		ID:    lastID,
		Limit: limit,
	}

	return EncodeCursor(cursor)
}

// ExtractIDFromCursor extracts the ID from a cursor string
func ExtractIDFromCursor(cursorStr string) (string, error) {
	cursor, err := DecodeCursor(cursorStr)
	if err != nil {
		return "", err
	}
	if cursor == nil {
		return "", nil
	}
	return cursor.ID, nil
}

// IsValidCursor checks if a cursor string is valid
func IsValidCursor(cursorStr string) bool {
	if cursorStr == "" {
		return true
	}
	_, err := DecodeCursor(cursorStr)
	return err == nil
}