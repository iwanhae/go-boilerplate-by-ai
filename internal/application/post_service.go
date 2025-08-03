package application

import (
	"context"
	"fmt"
	"sort"
	"time"

	"gosuda.org/boilerplate/internal/domain"
)

// PostService handles business logic for posts
type PostService struct {
	store domain.Store
}

// NewPostService creates a new post service
func NewPostService(store domain.Store) *PostService {
	return &PostService{
		store: store,
	}
}

// CreatePost creates a new post
func (s *PostService) CreatePost(ctx context.Context, req *domain.CreatePostRequest) (*domain.Post, error) {
	// Validate request
	if err := req.Validate(); err != nil {
		return nil, err
	}

	// Generate ID
	id := generatePostID()

	// Create post
	post := domain.NewPost(id, req.Title, req.Content)

	// Store post
	key := postKey(id)
	if err := s.store.Set(key, post); err != nil {
		return nil, &domain.StorageError{Err: err}
	}

	return post, nil
}

// GetPost retrieves a post by ID
func (s *PostService) GetPost(ctx context.Context, id string) (*domain.Post, error) {
	if err := validatePostID(id); err != nil {
		return nil, err
	}

	key := postKey(id)
	var post domain.Post
	if err := s.store.GetTyped(key, &post); err != nil {
		if err == domain.ErrKeyNotFound {
			return nil, &domain.PostNotFoundError{ID: id}
		}
		return nil, &domain.StorageError{Err: err}
	}

	return &post, nil
}

// UpdatePost updates an existing post
func (s *PostService) UpdatePost(ctx context.Context, id string, req *domain.UpdatePostRequest) (*domain.Post, error) {
	if err := validatePostID(id); err != nil {
		return nil, err
	}

	// Validate request
	if err := req.Validate(); err != nil {
		return nil, err
	}

	// Get existing post
	key := postKey(id)
	var post domain.Post
	if err := s.store.GetTyped(key, &post); err != nil {
		if err == domain.ErrKeyNotFound {
			return nil, &domain.PostNotFoundError{ID: id}
		}
		return nil, &domain.StorageError{Err: err}
	}

	// Update post
	post.Update(req.Title, req.Content)

	// Store updated post
	if err := s.store.Set(key, &post); err != nil {
		return nil, &domain.StorageError{Err: err}
	}

	return &post, nil
}

// DeletePost deletes a post by ID
func (s *PostService) DeletePost(ctx context.Context, id string) error {
	if err := validatePostID(id); err != nil {
		return err
	}

	key := postKey(id)
	if err := s.store.Delete(key); err != nil {
		if err == domain.ErrKeyNotFound {
			return &domain.PostNotFoundError{ID: id}
		}
		return &domain.StorageError{Err: err}
	}

	return nil
}

// ListPosts retrieves a paginated list of posts
func (s *PostService) ListPosts(ctx context.Context, cursor string, limit int) (*domain.PostList, error) {
	// Parse and validate pagination parameters
	params := NewPaginationParams(cursor, limit)
	if err := ValidatePaginationParams(params.Cursor, params.Limit); err != nil {
		return nil, &domain.ValidationError{
			Field:   "pagination",
			Message: err.Error(),
		}
	}

	// Get all posts
	values, err := s.store.List("posts:")
	if err != nil {
		return nil, &domain.StorageError{Err: err}
	}

	// Convert to posts
	var posts []domain.Post
	for _, value := range values {
		if post, ok := value.(domain.Post); ok {
			posts = append(posts, post)
		}
	}

	// Sort by creation date (newest first)
	sort.Slice(posts, func(i, j int) bool {
		return posts[i].CreatedAt.After(posts[j].CreatedAt)
	})

	// Apply pagination
	startIndex := 0
	if params.Cursor != "" {
		cursorObj, err := DecodeCursor(params.Cursor)
		if err != nil {
			return nil, &domain.PaginationError{Cursor: params.Cursor}
		}
		if cursorObj != nil {
			// Find the post with the cursor ID
			for i, post := range posts {
				if post.ID == cursorObj.ID {
					startIndex = i + 1
					break
				}
			}
		}
	}

	endIndex := startIndex + params.Limit
	if endIndex > len(posts) {
		endIndex = len(posts)
	}

	// Get the page of posts
	pagePosts := posts[startIndex:endIndex]

	// Create next cursor
	var nextCursor string
	if endIndex < len(posts) {
		nextCursor, err = CreateNextCursor(posts[endIndex-1].ID, params.Limit)
		if err != nil {
			return nil, &domain.StorageError{Err: err}
		}
	}

	return &domain.PostList{
		Posts:     pagePosts,
		NextCursor: nextCursor,
	}, nil
}

// validatePostID validates a post ID
func validatePostID(id string) error {
	if id == "" {
		return &domain.ValidationError{
			Field:   "id",
			Message: "post ID is required",
		}
	}

	// Check if ID matches the pattern from OpenAPI spec
	if !isValidPostID(id) {
		return &domain.ValidationError{
			Field:   "id",
			Message: "invalid post ID format",
		}
	}

	return nil
}

// isValidPostID checks if a post ID is valid according to the OpenAPI pattern
func isValidPostID(id string) bool {
	// Pattern: ^[a-zA-Z0-9-]+$
	for _, r := range id {
		if !((r >= 'a' && r <= 'z') || (r >= 'A' && r <= 'Z') || (r >= '0' && r <= '9') || r == '-') {
			return false
		}
	}
	return true
}

// postKey generates a storage key for a post
func postKey(id string) string {
	return fmt.Sprintf("posts:%s", id)
}

// generatePostID generates a unique post ID
func generatePostID() string {
	// Simple ID generation - in a real app, you might use UUID or a more sophisticated approach
	return fmt.Sprintf("post-%d", time.Now().UnixNano())
}