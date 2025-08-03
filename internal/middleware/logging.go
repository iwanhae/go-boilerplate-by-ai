package middleware

import (
	"context"
	"net/http"
	"time"

	"gosuda.org/boilerplate/internal/infrastructure"
)

// LoggingMiddleware provides HTTP request/response logging
type LoggingMiddleware struct {
	logger infrastructure.LoggerInterface
}

// NewLoggingMiddleware creates a new logging middleware
func NewLoggingMiddleware(logger infrastructure.LoggerInterface) *LoggingMiddleware {
	return &LoggingMiddleware{
		logger: logger,
	}
}

// Handler returns the logging middleware handler
func (m *LoggingMiddleware) Handler(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		start := time.Now()

		// Create a response writer wrapper to capture status code
		wrappedWriter := &responseWriter{ResponseWriter: w, statusCode: http.StatusOK}

		// Log the request
		m.logger.LogHTTPRequest(
			r.Context(),
			r.Method,
			r.URL.Path,
			r.RemoteAddr,
			r.UserAgent(),
			0, // Status code will be logged after response
			time.Since(start).Milliseconds(),
		)

		// Call the next handler
		next.ServeHTTP(wrappedWriter, r)

		// Calculate duration
		duration := time.Since(start)

		// Log the response
		m.logger.LogHTTPRequest(
			r.Context(),
			r.Method,
			r.URL.Path,
			r.RemoteAddr,
			r.UserAgent(),
			wrappedWriter.statusCode,
			duration.Milliseconds(),
		)

		// Log errors if status code indicates an error
		if wrappedWriter.statusCode >= 400 {
			m.logger.LogHTTPError(
				r.Context(),
				r.Method,
				r.URL.Path,
				wrappedWriter.statusCode,
				nil, // Error details would be available in a real implementation
			)
		}
	})
}

// responseWriter wraps http.ResponseWriter to capture status code
type responseWriter struct {
	http.ResponseWriter
	statusCode int
}

// WriteHeader captures the status code
func (rw *responseWriter) WriteHeader(code int) {
	rw.statusCode = code
	rw.ResponseWriter.WriteHeader(code)
}

// WithContext adds the logging middleware to a context
func (m *LoggingMiddleware) WithContext(ctx context.Context) context.Context {
	return context.WithValue(ctx, "logging_middleware", m)
}

// LoggingFromContext retrieves the logging middleware from a context
func LoggingFromContext(ctx context.Context) *LoggingMiddleware {
	if middleware, ok := ctx.Value("logging_middleware").(*LoggingMiddleware); ok {
		return middleware
	}
	return nil
}