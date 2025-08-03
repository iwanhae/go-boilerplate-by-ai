package middleware

import (
	"context"
	"net/http"

	"github.com/google/uuid"
)

// RequestIDMiddleware generates and propagates request IDs
type RequestIDMiddleware struct{}

// NewRequestIDMiddleware creates a new request ID middleware
func NewRequestIDMiddleware() *RequestIDMiddleware {
	return &RequestIDMiddleware{}
}

// Handler returns the request ID middleware handler
func (m *RequestIDMiddleware) Handler(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		// Check if request ID is already present in headers
		requestID := r.Header.Get("X-Request-ID")
		if requestID == "" {
			// Generate a new request ID
			requestID = uuid.New().String()
		}

		// Add request ID to response headers
		w.Header().Set("X-Request-ID", requestID)

		// Add request ID to request context
		ctx := context.WithValue(r.Context(), "request_id", requestID)
		r = r.WithContext(ctx)

		next.ServeHTTP(w, r)
	})
}

// GetRequestID retrieves the request ID from the context
func GetRequestID(ctx context.Context) string {
	if requestID, ok := ctx.Value("request_id").(string); ok {
		return requestID
	}
	return ""
}

// WithRequestID adds a request ID to a context
func WithRequestID(ctx context.Context, requestID string) context.Context {
	return context.WithValue(ctx, "request_id", requestID)
}

// WithContext adds the request ID middleware to a context
func (m *RequestIDMiddleware) WithContext(ctx context.Context) context.Context {
	return context.WithValue(ctx, "request_id_middleware", m)
}

// RequestIDFromContext retrieves the request ID middleware from a context
func RequestIDFromContext(ctx context.Context) *RequestIDMiddleware {
	if middleware, ok := ctx.Value("request_id_middleware").(*RequestIDMiddleware); ok {
		return middleware
	}
	return nil
}