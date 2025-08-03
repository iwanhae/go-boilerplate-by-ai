package domain

import (
	"time"
	"unicode/utf8"
)

// Post represents a blog post entity
type Post struct {
	ID        string    `json:"id"`
	Title     string    `json:"title"`
	Content   string    `json:"content"`
	CreatedAt time.Time `json:"createdAt"`
	UpdatedAt time.Time `json:"updatedAt"`
}

// CreatePostRequest represents a request to create a new post
type CreatePostRequest struct {
	Title   string `json:"title"`
	Content string `json:"content"`
}

// UpdatePostRequest represents a request to update an existing post
type UpdatePostRequest struct {
	Title   string `json:"title"`
	Content string `json:"content"`
}

// PostList represents a paginated list of posts
type PostList struct {
	Posts     []Post `json:"posts"`
	NextCursor string `json:"nextCursor,omitempty"`
}

// Validation constants
const (
	MinTitleLength   = 1
	MaxTitleLength   = 200
	MinContentLength = 1
	MaxContentLength = 10000
)

// ValidateCreateRequest validates a create post request
func (r *CreatePostRequest) Validate() error {
	if err := validateTitle(r.Title); err != nil {
		return err
	}
	if err := validateContent(r.Content); err != nil {
		return err
	}
	return nil
}

// ValidateUpdateRequest validates an update post request
func (r *UpdatePostRequest) Validate() error {
	if err := validateTitle(r.Title); err != nil {
		return err
	}
	if err := validateContent(r.Content); err != nil {
		return err
	}
	return nil
}

// validateTitle validates the title field
func validateTitle(title string) error {
	length := utf8.RuneCountInString(title)
	if length < MinTitleLength {
		return &ValidationError{
			Field:   "title",
			Message: "title is required",
		}
	}
	if length > MaxTitleLength {
		return &ValidationError{
			Field:   "title",
			Message: "title is too long",
		}
	}
	return nil
}

// validateContent validates the content field
func validateContent(content string) error {
	length := utf8.RuneCountInString(content)
	if length < MinContentLength {
		return &ValidationError{
			Field:   "content",
			Message: "content is required",
		}
	}
	if length > MaxContentLength {
		return &ValidationError{
			Field:   "content",
			Message: "content is too long",
		}
	}
	return nil
}

// NewPost creates a new post with the given data
func NewPost(id, title, content string) *Post {
	now := time.Now()
	return &Post{
		ID:        id,
		Title:     title,
		Content:   content,
		CreatedAt: now,
		UpdatedAt: now,
	}
}

// Update updates the post with new data
func (p *Post) Update(title, content string) {
	p.Title = title
	p.Content = content
	p.UpdatedAt = time.Now()
}