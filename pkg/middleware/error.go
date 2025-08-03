package middleware

import (
	"encoding/json"
	"net/http"
)

// ErrorResponse represents an error message.
type ErrorResponse struct {
	Error string `json:"error"`
}

// Recover responds with JSON on panics.
func Recover(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		defer func() {
			if rec := recover(); rec != nil {
				FromContext(r.Context()).Error("panic", "err", rec)
				w.Header().Set("Content-Type", "application/json")
				w.WriteHeader(http.StatusInternalServerError)
				_ = json.NewEncoder(w).Encode(ErrorResponse{Error: http.StatusText(http.StatusInternalServerError)})
			}
		}()
		next.ServeHTTP(w, r)
	})
}
