package domain

import (
	"testing"
	"time"
)

func TestNewPost(t *testing.T) {
	id := "test-123"
	title := "Test Post"
	content := "This is a test post content"
	
	post := NewPost(id, title, content)
	
	if post.ID != id {
		t.Errorf("Expected ID '%s', got '%s'", id, post.ID)
	}
	
	if post.Title != title {
		t.Errorf("Expected Title '%s', got '%s'", title, post.Title)
	}
	
	if post.Content != content {
		t.Errorf("Expected Content '%s', got '%s'", content, post.Content)
	}
	
	if post.CreatedAt.IsZero() {
		t.Error("CreatedAt should be set")
	}
	
	if post.UpdatedAt.IsZero() {
		t.Error("UpdatedAt should be set")
	}
	
	if post.CreatedAt != post.UpdatedAt {
		t.Error("CreatedAt and UpdatedAt should be equal for new posts")
	}
}

func TestPost_Update(t *testing.T) {
	post := NewPost("test-123", "Original Title", "Original content")
	originalCreatedAt := post.CreatedAt
	originalUpdatedAt := post.UpdatedAt
	
	// Wait a bit to ensure time difference
	time.Sleep(1 * time.Millisecond)
	
	newTitle := "Updated Title"
	newContent := "Updated content"
	post.Update(newTitle, newContent)
	
	if post.Title != newTitle {
		t.Errorf("Expected Title '%s', got '%s'", newTitle, post.Title)
	}
	
	if post.Content != newContent {
		t.Errorf("Expected Content '%s', got '%s'", newContent, post.Content)
	}
	
	if post.CreatedAt != originalCreatedAt {
		t.Error("CreatedAt should not change when updating")
	}
	
	if post.UpdatedAt == originalUpdatedAt {
		t.Error("UpdatedAt should change when updating")
	}
	
	if post.UpdatedAt.Before(originalUpdatedAt) {
		t.Error("UpdatedAt should be after the original UpdatedAt")
	}
}



func TestCreatePostRequest_Validate(t *testing.T) {
	testCases := []struct {
		name    string
		req     *CreatePostRequest
		wantErr bool
	}{
		{
			name: "valid request",
			req: &CreatePostRequest{
				Title:   "Valid Title",
				Content: "Valid content",
			},
			wantErr: false,
		},
		{
			name: "empty title",
			req: &CreatePostRequest{
				Title:   "",
				Content: "Valid content",
			},
			wantErr: true,
		},
		{
			name: "empty content",
			req: &CreatePostRequest{
				Title:   "Valid Title",
				Content: "",
			},
			wantErr: true,
		},
		{
			name: "title too long",
			req: &CreatePostRequest{
				Title:   string(make([]byte, 201)), // 201 characters
				Content: "Valid content",
			},
			wantErr: true,
		},
		{
			name: "content too long",
			req: &CreatePostRequest{
				Title:   "Valid Title",
				Content: string(make([]byte, 10001)), // 10001 characters
			},
			wantErr: true,
		},
	}
	
	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			err := tc.req.Validate()
			if (err != nil) != tc.wantErr {
				t.Errorf("CreatePostRequest.Validate() error = %v, wantErr %v", err, tc.wantErr)
			}
		})
	}
}

