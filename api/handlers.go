package api

import (
	"encoding/json"
	"net/http"
	"strconv"

	"gosuda.org/boilerplate/internal/application"
	"gosuda.org/boilerplate/internal/domain"
	"gosuda.org/boilerplate/internal/middleware"
)

// Handlers implements the API endpoints
type Handlers struct {
	postService  *application.PostService
	debugService *application.DebugService
	errorHandler *middleware.ErrorHandlerMiddleware
}

// NewHandlers creates new API handlers
func NewHandlers(
	postService *application.PostService,
	debugService *application.DebugService,
	errorHandler *middleware.ErrorHandlerMiddleware,
) *Handlers {
	return &Handlers{
		postService:  postService,
		debugService: debugService,
		errorHandler: errorHandler,
	}
}

// GetMetrics handles GET /debug/metrics
func (h *Handlers) GetMetrics(w http.ResponseWriter, r *http.Request) {
	metrics, err := h.debugService.GetMetrics(r.Context())
	if err != nil {
		h.errorHandler.HandleError(w, r, err)
		return
	}

	w.Header().Set("Content-Type", "text/plain")
	w.WriteHeader(http.StatusOK)
	w.Write([]byte(metrics))
}

// SetLogLevel handles POST /debug/logs
func (h *Handlers) SetLogLevel(w http.ResponseWriter, r *http.Request) {
	level := r.URL.Query().Get("level")
	if level == "" {
		h.errorHandler.HandleError(w, r, &domain.ValidationError{
			Field:   "level",
			Message: "level parameter is required",
		})
		return
	}

	err := h.debugService.SetLogLevel(r.Context(), level)
	if err != nil {
		h.errorHandler.HandleError(w, r, err)
		return
	}

	w.WriteHeader(http.StatusOK)
}

// GetPprofProfile handles GET /debug/pprof/{profile}
func (h *Handlers) GetPprofProfile(w http.ResponseWriter, r *http.Request) {
	// Extract profile from URL path
	profile := r.URL.Path[len("/debug/pprof/"):]
	if profile == "" {
		h.errorHandler.HandleError(w, r, &domain.ValidationError{
			Field:   "profile",
			Message: "profile parameter is required",
		})
		return
	}

	data, err := h.debugService.GetPprofProfile(r.Context(), profile)
	if err != nil {
		h.errorHandler.HandleError(w, r, err)
		return
	}

	w.Header().Set("Content-Type", "application/octet-stream")
	w.WriteHeader(http.StatusOK)
	w.Write(data)
}

// ListPosts handles GET /posts
func (h *Handlers) ListPosts(w http.ResponseWriter, r *http.Request) {
	cursor := r.URL.Query().Get("cursor")
	limitStr := r.URL.Query().Get("limit")

	limit := 20 // default
	if limitStr != "" {
		if parsedLimit, err := strconv.Atoi(limitStr); err == nil {
			limit = parsedLimit
		}
	}

	posts, err := h.postService.ListPosts(r.Context(), cursor, limit)
	if err != nil {
		h.errorHandler.HandleError(w, r, err)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(posts)
}

// CreatePost handles POST /posts
func (h *Handlers) CreatePost(w http.ResponseWriter, r *http.Request) {
	var req domain.CreatePostRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		h.errorHandler.HandleError(w, r, &domain.ValidationError{
			Field:   "body",
			Message: "invalid JSON body",
		})
		return
	}

	post, err := h.postService.CreatePost(r.Context(), &req)
	if err != nil {
		h.errorHandler.HandleError(w, r, err)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusCreated)
	json.NewEncoder(w).Encode(post)
}

// GetPost handles GET /posts/{id}
func (h *Handlers) GetPost(w http.ResponseWriter, r *http.Request) {
	// Extract ID from URL path
	id := r.URL.Path[len("/posts/"):]
	if id == "" {
		h.errorHandler.HandleError(w, r, &domain.ValidationError{
			Field:   "id",
			Message: "post ID is required",
		})
		return
	}

	post, err := h.postService.GetPost(r.Context(), id)
	if err != nil {
		h.errorHandler.HandleError(w, r, err)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(post)
}

// UpdatePost handles PUT /posts/{id}
func (h *Handlers) UpdatePost(w http.ResponseWriter, r *http.Request) {
	// Extract ID from URL path
	id := r.URL.Path[len("/posts/"):]
	if id == "" {
		h.errorHandler.HandleError(w, r, &domain.ValidationError{
			Field:   "id",
			Message: "post ID is required",
		})
		return
	}

	var req domain.UpdatePostRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		h.errorHandler.HandleError(w, r, &domain.ValidationError{
			Field:   "body",
			Message: "invalid JSON body",
		})
		return
	}

	post, err := h.postService.UpdatePost(r.Context(), id, &req)
	if err != nil {
		h.errorHandler.HandleError(w, r, err)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(post)
}

// DeletePost handles DELETE /posts/{id}
func (h *Handlers) DeletePost(w http.ResponseWriter, r *http.Request) {
	// Extract ID from URL path
	id := r.URL.Path[len("/posts/"):]
	if id == "" {
		h.errorHandler.HandleError(w, r, &domain.ValidationError{
			Field:   "id",
			Message: "post ID is required",
		})
		return
	}

	err := h.postService.DeletePost(r.Context(), id)
	if err != nil {
		h.errorHandler.HandleError(w, r, err)
		return
	}

	w.WriteHeader(http.StatusNoContent)
}

// GetHealth handles GET /health
func (h *Handlers) GetHealth(w http.ResponseWriter, r *http.Request) {
	status, err := h.debugService.GetHealthStatus(r.Context())
	if err != nil {
		h.errorHandler.HandleError(w, r, err)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(status)
}