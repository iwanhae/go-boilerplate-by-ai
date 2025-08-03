package middleware

import (
	"context"
	"net/http"
	"time"

	chimid "github.com/go-chi/chi/v5/middleware"
	"log/slog"
)

type loggerKey struct{}

// Logger injects a request-scoped slog logger and logs request details.
func Logger(l *slog.Logger) func(http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			id := chimid.GetReqID(r.Context())
			logger := l.With("request_id", id, "method", r.Method, "path", r.URL.Path)
			ctx := context.WithValue(r.Context(), loggerKey{}, logger)
			start := time.Now()
			defer logger.Info("completed", "duration", time.Since(start))
			next.ServeHTTP(w, r.WithContext(ctx))
		})
	}
}

// FromContext retrieves the request logger from context.
func FromContext(ctx context.Context) *slog.Logger {
	if l, ok := ctx.Value(loggerKey{}).(*slog.Logger); ok {
		return l
	}
	return slog.Default()
}
