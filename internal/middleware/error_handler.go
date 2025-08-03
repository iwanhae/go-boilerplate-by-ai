package middleware

import (
	"context"
	"encoding/json"
	"net/http"

	"gosuda.org/boilerplate/internal/domain"
	"gosuda.org/boilerplate/internal/infrastructure"
)

// ErrorResponse represents a standardized error response
type ErrorResponse struct {
	Code      string `json:"code"`
	Message   string `json:"message"`
	RequestID string `json:"requestId,omitempty"`
}

// ErrorHandlerMiddleware provides centralized error handling
type ErrorHandlerMiddleware struct {
	logger infrastructure.LoggerInterface
}

// NewErrorHandlerMiddleware creates a new error handler middleware
func NewErrorHandlerMiddleware(logger infrastructure.LoggerInterface) *ErrorHandlerMiddleware {
	return &ErrorHandlerMiddleware{
		logger: logger,
	}
}

// Handler returns the error handler middleware
func (m *ErrorHandlerMiddleware) Handler(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		// Create a response writer wrapper to capture errors
		wrappedWriter := &errorResponseWriter{
			ResponseWriter: w,
			errorHandler:   m,
			request:        r,
		}

		next.ServeHTTP(wrappedWriter, r)
	})
}

// errorResponseWriter wraps http.ResponseWriter to handle errors
type errorResponseWriter struct {
	http.ResponseWriter
	errorHandler *ErrorHandlerMiddleware
	request      *http.Request
	statusCode   int
	written      bool
}

// WriteHeader captures the status code and handles errors
func (rw *errorResponseWriter) WriteHeader(code int) {
	rw.statusCode = code
	rw.ResponseWriter.WriteHeader(code)
	rw.written = true
}

// Write handles writing response data
func (rw *errorResponseWriter) Write(data []byte) (int, error) {
	if !rw.written {
		rw.WriteHeader(http.StatusOK)
	}
	return rw.ResponseWriter.Write(data)
}

// HandleError handles domain errors and converts them to HTTP responses
func (m *ErrorHandlerMiddleware) HandleError(w http.ResponseWriter, r *http.Request, err error) {
	// Get request ID from context
	requestID := GetRequestID(r.Context())

	// Determine status code and error response based on error type
	statusCode, errorResponse := m.mapErrorToResponse(err, requestID)

	// Log the error
	m.logger.LogHTTPError(
		r.Context(),
		r.Method,
		r.URL.Path,
		statusCode,
		err,
	)

	// Set response headers
	w.Header().Set("Content-Type", "application/json")
	w.Header().Set("X-Request-ID", requestID)
	w.WriteHeader(statusCode)

	// Write error response
	json.NewEncoder(w).Encode(errorResponse)
}

// mapErrorToResponse maps domain errors to HTTP status codes and responses
func (m *ErrorHandlerMiddleware) mapErrorToResponse(err error, requestID string) (int, ErrorResponse) {
	switch e := err.(type) {
	case *domain.PostNotFoundError:
		return http.StatusNotFound, ErrorResponse{
			Code:      domain.ErrorCodePostNotFound,
			Message:   e.Error(),
			RequestID: requestID,
		}
	case *domain.InvalidPostDataError:
		return http.StatusBadRequest, ErrorResponse{
			Code:      domain.ErrorCodeInvalidPostData,
			Message:   e.Error(),
			RequestID: requestID,
		}
	case *domain.ValidationError:
		return http.StatusBadRequest, ErrorResponse{
			Code:      domain.ErrorCodeValidationError,
			Message:   e.Error(),
			RequestID: requestID,
		}
	case *domain.PaginationError:
		return http.StatusBadRequest, ErrorResponse{
			Code:      domain.ErrorCodePaginationError,
			Message:   e.Error(),
			RequestID: requestID,
		}
	case *domain.StorageError:
		return http.StatusInternalServerError, ErrorResponse{
			Code:      domain.ErrorCodeStorageError,
			Message:   e.Error(),
			RequestID: requestID,
		}
	default:
		// Default to internal server error for unknown errors
		return http.StatusInternalServerError, ErrorResponse{
			Code:      domain.ErrorCodeInternalError,
			Message:   "Internal server error",
			RequestID: requestID,
		}
	}
}

// WithContext adds the error handler middleware to a context
func (m *ErrorHandlerMiddleware) WithContext(ctx context.Context) context.Context {
	return context.WithValue(ctx, "error_handler_middleware", m)
}

// ErrorHandlerFromContext retrieves the error handler middleware from a context
func ErrorHandlerFromContext(ctx context.Context) *ErrorHandlerMiddleware {
	if middleware, ok := ctx.Value("error_handler_middleware").(*ErrorHandlerMiddleware); ok {
		return middleware
	}
	return nil
}