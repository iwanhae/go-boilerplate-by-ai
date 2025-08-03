package middleware

import (
	"context"
	"net/http"
	"runtime/debug"

	"gosuda.org/boilerplate/internal/infrastructure"
)

// RecoveryMiddleware provides panic recovery for HTTP handlers
type RecoveryMiddleware struct {
	logger infrastructure.LoggerInterface
}

// NewRecoveryMiddleware creates a new recovery middleware
func NewRecoveryMiddleware(logger infrastructure.LoggerInterface) *RecoveryMiddleware {
	return &RecoveryMiddleware{
		logger: logger,
	}
}

// Handler returns the recovery middleware handler
func (m *RecoveryMiddleware) Handler(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		defer func() {
			if err := recover(); err != nil {
				// Log the panic
				m.logger.Error("panic recovered",
					"error", err,
					"stack", string(debug.Stack()),
					"method", r.Method,
					"path", r.URL.Path,
					"remote_addr", r.RemoteAddr,
					"user_agent", r.UserAgent(),
				)

				// Get request ID from context if available
				requestID := ""
				if ctx := r.Context(); ctx != nil {
					if id, ok := ctx.Value("request_id").(string); ok {
						requestID = id
					}
				}

				// Return 500 Internal Server Error
				w.Header().Set("Content-Type", "application/json")
				w.Header().Set("X-Request-ID", requestID)
				w.WriteHeader(http.StatusInternalServerError)

				// In a real implementation, you'd use proper JSON encoding
				// For now, we'll write a simple JSON string
				jsonResponse := `{"code":"INTERNAL_ERROR","message":"Internal server error","requestId":"` + requestID + `"}`
				w.Write([]byte(jsonResponse))
			}
		}()

		next.ServeHTTP(w, r)
	})
}

// WithContext adds the recovery middleware to a context
func (m *RecoveryMiddleware) WithContext(ctx context.Context) context.Context {
	return context.WithValue(ctx, "recovery_middleware", m)
}

// RecoveryFromContext retrieves the recovery middleware from a context
func RecoveryFromContext(ctx context.Context) *RecoveryMiddleware {
	if middleware, ok := ctx.Value("recovery_middleware").(*RecoveryMiddleware); ok {
		return middleware
	}
	return nil
}