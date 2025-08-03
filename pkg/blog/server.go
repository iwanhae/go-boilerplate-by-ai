package blog

import (
	"context"
	"errors"

	"github.com/google/uuid"
	"gosuda.org/boilerplate/api"
	"gosuda.org/boilerplate/pkg/store"
)

// Server implements blog CRUD operations.
type Server struct {
	store store.Store
}

// New creates a Server.
func New(st store.Store) *Server { return &Server{store: st} }

const prefix = "post:"

func key(id string) string { return prefix + id }

func (s *Server) ListPosts(ctx context.Context, request api.ListPostsRequestObject) (api.ListPostsResponseObject, error) {
	items, err := s.store.List(prefix)
	if err != nil {
		return nil, err
	}
	posts := make([]api.Post, 0, len(items))
	for _, v := range items {
		posts = append(posts, v.(api.Post))
	}
	return api.ListPosts200JSONResponse(posts), nil
}

func (s *Server) CreatePost(ctx context.Context, request api.CreatePostRequestObject) (api.CreatePostResponseObject, error) {
	id := uuid.NewString()
	post := api.Post{Id: id, Title: request.Body.Title, Content: request.Body.Content}
	if err := s.store.Set(key(id), post); err != nil {
		return nil, err
	}
	return api.CreatePost201JSONResponse(post), nil
}

func (s *Server) GetPost(ctx context.Context, request api.GetPostRequestObject) (api.GetPostResponseObject, error) {
	v, err := s.store.Get(key(request.Id))
	if err != nil {
		if errors.Is(err, store.ErrNotFound) {
			return api.GetPost404Response{}, nil
		}
		return nil, err
	}
	return api.GetPost200JSONResponse(v.(api.Post)), nil
}

func (s *Server) UpdatePost(ctx context.Context, request api.UpdatePostRequestObject) (api.UpdatePostResponseObject, error) {
	if _, err := s.store.Get(key(request.Id)); err != nil {
		if errors.Is(err, store.ErrNotFound) {
			return api.UpdatePost404Response{}, nil
		}
		return nil, err
	}
	post := api.Post{Id: request.Id, Title: request.Body.Title, Content: request.Body.Content}
	if err := s.store.Set(key(request.Id), post); err != nil {
		return nil, err
	}
	return api.UpdatePost200JSONResponse(post), nil
}

func (s *Server) DeletePost(ctx context.Context, request api.DeletePostRequestObject) (api.DeletePostResponseObject, error) {
	if err := s.store.Delete(key(request.Id)); err != nil {
		if errors.Is(err, store.ErrNotFound) {
			return api.DeletePost404Response{}, nil
		}
		return nil, err
	}
	return api.DeletePost204Response{}, nil
}
