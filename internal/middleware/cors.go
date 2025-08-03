package middleware

import (
	"context"
	"net/http"
	"strings"

	"gosuda.org/boilerplate/internal/config"
)

// CORSMiddleware handles Cross-Origin Resource Sharing
type CORSMiddleware struct {
	config *config.CORSConfig
}

// NewCORSMiddleware creates a new CORS middleware
func NewCORSMiddleware(cfg *config.CORSConfig) *CORSMiddleware {
	return &CORSMiddleware{
		config: cfg,
	}
}

// Handler returns the CORS middleware handler
func (m *CORSMiddleware) Handler(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		// Handle preflight requests
		if r.Method == "OPTIONS" {
			m.handlePreflight(w, r)
			return
		}

		// Set CORS headers for actual requests
		m.setCORSHeaders(w, r)

		next.ServeHTTP(w, r)
	})
}

// handlePreflight handles OPTIONS preflight requests
func (m *CORSMiddleware) handlePreflight(w http.ResponseWriter, r *http.Request) {
	// Set CORS headers
	m.setCORSHeaders(w, r)

	// Set allowed methods
	if len(m.config.AllowedMethods) > 0 {
		w.Header().Set("Access-Control-Allow-Methods", strings.Join(m.config.AllowedMethods, ", "))
	}

	// Set allowed headers
	if len(m.config.AllowedHeaders) > 0 {
		w.Header().Set("Access-Control-Allow-Headers", strings.Join(m.config.AllowedHeaders, ", "))
	}

	// Set max age
	if m.config.MaxAge > 0 {
		w.Header().Set("Access-Control-Max-Age", string(rune(m.config.MaxAge)))
	}

	// Return 200 OK for preflight requests
	w.WriteHeader(http.StatusOK)
}

// setCORSHeaders sets CORS headers for requests
func (m *CORSMiddleware) setCORSHeaders(w http.ResponseWriter, r *http.Request) {
	// Set allowed origins
	if len(m.config.AllowedOrigins) > 0 {
		origin := r.Header.Get("Origin")
		if m.isOriginAllowed(origin) {
			w.Header().Set("Access-Control-Allow-Origin", origin)
		} else if m.config.AllowedOrigins[0] == "*" {
			w.Header().Set("Access-Control-Allow-Origin", "*")
		}
	}

	// Set credentials
	w.Header().Set("Access-Control-Allow-Credentials", "true")

	// Set exposed headers
	w.Header().Set("Access-Control-Expose-Headers", "X-Request-ID")
}

// isOriginAllowed checks if the origin is allowed
func (m *CORSMiddleware) isOriginAllowed(origin string) bool {
	if len(m.config.AllowedOrigins) == 0 {
		return false
	}

	// Allow all origins if "*" is specified
	if m.config.AllowedOrigins[0] == "*" {
		return true
	}

	// Check if origin is in the allowed list
	for _, allowedOrigin := range m.config.AllowedOrigins {
		if allowedOrigin == origin {
			return true
		}
	}

	return false
}

// WithContext adds the CORS middleware to a context
func (m *CORSMiddleware) WithContext(ctx context.Context) context.Context {
	return context.WithValue(ctx, "cors_middleware", m)
}

// CORSFromContext retrieves the CORS middleware from a context
func CORSFromContext(ctx context.Context) *CORSMiddleware {
	if middleware, ok := ctx.Value("cors_middleware").(*CORSMiddleware); ok {
		return middleware
	}
	return nil
}