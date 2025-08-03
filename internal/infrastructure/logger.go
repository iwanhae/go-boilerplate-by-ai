package infrastructure

import (
	"context"
	"io"
	"log/slog"
	"os"
	"sync/atomic"

	"gosuda.org/boilerplate/internal/config"
)

// LoggerInterface defines the interface for logging functionality
type LoggerInterface interface {
	SetLevel(level string) error
	GetLevel() slog.Level
	Debug(msg string, args ...any)
	Info(msg string, args ...any)
	Warn(msg string, args ...any)
	Error(msg string, args ...any)
	WithContext(ctx context.Context) *slog.Logger
	WithRequestID(requestID string) *slog.Logger
	WithFields(fields map[string]any) *slog.Logger
	LogHTTPRequest(ctx context.Context, method, path, remoteAddr, userAgent string, statusCode int, durationMs int64)
	LogHTTPError(ctx context.Context, method, path string, statusCode int, err error)
	LogStorageOperation(ctx context.Context, operation, key string, err error)
	LogStartup(version, commitHash string, config *config.Config)
	LogShutdown(reason string)
	LogGracefulShutdown(phase string, remainingRequests int)
	LogLevelChange(oldLevel, newLevel string)
}

// Logger provides structured logging functionality
type Logger struct {
	logger *slog.Logger
	level  atomic.Value // stores slog.Level
}

// Ensure Logger implements LoggerInterface
var _ LoggerInterface = (*Logger)(nil)

// NewLogger creates a new logger instance
func NewLogger(cfg *config.LoggingConfig) (*Logger, error) {
	var writer io.Writer
	switch cfg.Output {
	case "stdout":
		writer = os.Stdout
	case "stderr":
		writer = os.Stderr
	default:
		writer = os.Stdout
	}

	var level slog.Level
	switch cfg.Level {
	case "debug":
		level = slog.LevelDebug
	case "info":
		level = slog.LevelInfo
	case "warn":
		level = slog.LevelWarn
	case "error":
		level = slog.LevelError
	default:
		level = slog.LevelInfo
	}

	var handler slog.Handler
	switch cfg.Format {
	case "text":
		handler = slog.NewTextHandler(writer, &slog.HandlerOptions{
			Level: level,
		})
	case "json":
		handler = slog.NewJSONHandler(writer, &slog.HandlerOptions{
			Level: level,
		})
	default:
		handler = slog.NewJSONHandler(writer, &slog.HandlerOptions{
			Level: level,
		})
	}

	logger := &Logger{
		logger: slog.New(handler),
	}
	logger.level.Store(level)

	return logger, nil
}

// SetLevel changes the log level at runtime
func (l *Logger) SetLevel(level string) error {
	var slogLevel slog.Level
	switch level {
	case "debug":
		slogLevel = slog.LevelDebug
	case "info":
		slogLevel = slog.LevelInfo
	case "warn":
		slogLevel = slog.LevelWarn
	case "error":
		slogLevel = slog.LevelError
	default:
		return ErrInvalidLogLevel
	}

	l.level.Store(slogLevel)
	return nil
}

// GetLevel returns the current log level
func (l *Logger) GetLevel() slog.Level {
	return l.level.Load().(slog.Level)
}

// Debug logs a debug message
func (l *Logger) Debug(msg string, args ...any) {
	l.logger.Debug(msg, args...)
}

// Info logs an info message
func (l *Logger) Info(msg string, args ...any) {
	l.logger.Info(msg, args...)
}

// Warn logs a warning message
func (l *Logger) Warn(msg string, args ...any) {
	l.logger.Warn(msg, args...)
}

// Error logs an error message
func (l *Logger) Error(msg string, args ...any) {
	l.logger.Error(msg, args...)
}

// WithContext creates a logger with context values
func (l *Logger) WithContext(ctx context.Context) *slog.Logger {
	if ctx == nil {
		return l.logger
	}

	// Extract request ID from context
	if requestID, ok := ctx.Value("request_id").(string); ok && requestID != "" {
		return l.logger.With("request_id", requestID)
	}

	return l.logger
}

// WithRequestID creates a logger with request ID
func (l *Logger) WithRequestID(requestID string) *slog.Logger {
	if requestID == "" {
		return l.logger
	}
	return l.logger.With("request_id", requestID)
}

// WithFields creates a logger with additional fields
func (l *Logger) WithFields(fields map[string]any) *slog.Logger {
	if len(fields) == 0 {
		return l.logger
	}

	args := make([]any, 0, len(fields)*2)
	for key, value := range fields {
		args = append(args, key, value)
	}

	return l.logger.With(args...)
}

// LogHTTPRequest logs HTTP request information
func (l *Logger) LogHTTPRequest(ctx context.Context, method, path, remoteAddr, userAgent string, statusCode int, durationMs int64) {
	logger := l.WithContext(ctx)
	logger.Info("HTTP request completed",
		"method", method,
		"path", path,
		"remote_addr", remoteAddr,
		"user_agent", userAgent,
		"status_code", statusCode,
		"duration_ms", durationMs,
	)
}

// LogHTTPError logs HTTP error information
func (l *Logger) LogHTTPError(ctx context.Context, method, path string, statusCode int, err error) {
	logger := l.WithContext(ctx)
	logger.Error("HTTP request failed",
		"method", method,
		"path", path,
		"status_code", statusCode,
		"error", err.Error(),
	)
}

// LogStorageOperation logs storage operation information
func (l *Logger) LogStorageOperation(ctx context.Context, operation, key string, err error) {
	logger := l.WithContext(ctx)
	if err != nil {
		logger.Error("Storage operation failed",
			"operation", operation,
			"key", key,
			"error", err.Error(),
		)
	} else {
		logger.Debug("Storage operation completed",
			"operation", operation,
			"key", key,
		)
	}
}

// LogStartup logs application startup information
func (l *Logger) LogStartup(version, commitHash string, config *config.Config) {
	l.Info("Application starting",
		"version", version,
		"commit_hash", commitHash,
		"server_host", config.Server.Host,
		"server_port", config.Server.Port,
		"logging_level", config.Logging.Level,
		"logging_format", config.Logging.Format,
		"storage_type", config.Storage.Type,
	)
}

// LogShutdown logs application shutdown information
func (l *Logger) LogShutdown(reason string) {
	l.Info("Application shutting down", "reason", reason)
}

// LogGracefulShutdown logs graceful shutdown progress
func (l *Logger) LogGracefulShutdown(phase string, remainingRequests int) {
	l.Info("Graceful shutdown in progress",
		"phase", phase,
		"remaining_requests", remainingRequests,
	)
}

// LogLevelChange logs log level change
func (l *Logger) LogLevelChange(oldLevel, newLevel string) {
	l.Info("Log level changed",
		"old_level", oldLevel,
		"new_level", newLevel,
	)
}