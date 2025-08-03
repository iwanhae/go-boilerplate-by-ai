package infrastructure

import (
	"context"
	"log/slog"
	"testing"

	"gosuda.org/boilerplate/internal/config"
)

func TestNewLogger(t *testing.T) {
	testCases := []struct {
		name    string
		config  *config.LoggingConfig
		wantErr bool
	}{
		{
			name: "valid json config",
			config: &config.LoggingConfig{
				Level:  "info",
				Format: "json",
				Output: "stdout",
			},
			wantErr: false,
		},
		{
			name: "valid text config",
			config: &config.LoggingConfig{
				Level:  "debug",
				Format: "text",
				Output: "stderr",
			},
			wantErr: false,
		},
		{
			name: "default config",
			config: &config.LoggingConfig{
				Level:  "warn",
				Format: "invalid",
				Output: "invalid",
			},
			wantErr: false,
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			logger, err := NewLogger(tc.config)
			if tc.wantErr && err == nil {
				t.Error("Expected error but got none")
			}
			if !tc.wantErr && err != nil {
				t.Errorf("Unexpected error: %v", err)
			}
			if !tc.wantErr && logger == nil {
				t.Error("Expected logger but got nil")
			}
		})
	}
}

func TestLogger_SetLevel(t *testing.T) {
	cfg := &config.LoggingConfig{
		Level:  "info",
		Format: "json",
		Output: "stdout",
	}

	logger, err := NewLogger(cfg)
	if err != nil {
		t.Fatalf("Failed to create logger: %v", err)
	}

	testCases := []struct {
		name    string
		level   string
		wantErr bool
	}{
		{"debug level", "debug", false},
		{"info level", "info", false},
		{"warn level", "warn", false},
		{"error level", "error", false},
		{"invalid level", "invalid", true},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			err := logger.SetLevel(tc.level)
			if tc.wantErr && err == nil {
				t.Error("Expected error but got none")
			}
			if !tc.wantErr && err != nil {
				t.Errorf("Unexpected error: %v", err)
			}
		})
	}
}

func TestLogger_GetLevel(t *testing.T) {
	cfg := &config.LoggingConfig{
		Level:  "info",
		Format: "json",
		Output: "stdout",
	}

	logger, err := NewLogger(cfg)
	if err != nil {
		t.Fatalf("Failed to create logger: %v", err)
	}

	level := logger.GetLevel()
	if level != slog.LevelInfo {
		t.Errorf("Expected level %v, got %v", slog.LevelInfo, level)
	}
}

func TestLogger_WithContext(t *testing.T) {
	cfg := &config.LoggingConfig{
		Level:  "info",
		Format: "json",
		Output: "stdout",
	}

	logger, err := NewLogger(cfg)
	if err != nil {
		t.Fatalf("Failed to create logger: %v", err)
	}

	// Test with nil context
	result := logger.WithContext(nil)
	if result == nil {
		t.Error("Expected logger but got nil")
	}

	// Test with context without request ID
	ctx := context.Background()
	result = logger.WithContext(ctx)
	if result == nil {
		t.Error("Expected logger but got nil")
	}

	// Test with context with request ID
	ctxWithID := context.WithValue(context.Background(), "request_id", "test-123")
	result = logger.WithContext(ctxWithID)
	if result == nil {
		t.Error("Expected logger but got nil")
	}
}

func TestLogger_WithRequestID(t *testing.T) {
	cfg := &config.LoggingConfig{
		Level:  "info",
		Format: "json",
		Output: "stdout",
	}

	logger, err := NewLogger(cfg)
	if err != nil {
		t.Fatalf("Failed to create logger: %v", err)
	}

	// Test with empty request ID
	result := logger.WithRequestID("")
	if result == nil {
		t.Error("Expected logger but got nil")
	}

	// Test with valid request ID
	result = logger.WithRequestID("test-123")
	if result == nil {
		t.Error("Expected logger but got nil")
	}
}

func TestLogger_WithFields(t *testing.T) {
	cfg := &config.LoggingConfig{
		Level:  "info",
		Format: "json",
		Output: "stdout",
	}

	logger, err := NewLogger(cfg)
	if err != nil {
		t.Fatalf("Failed to create logger: %v", err)
	}

	// Test with empty fields
	result := logger.WithFields(nil)
	if result == nil {
		t.Error("Expected logger but got nil")
	}

	result = logger.WithFields(map[string]any{})
	if result == nil {
		t.Error("Expected logger but got nil")
	}

	// Test with valid fields
	fields := map[string]any{
		"key1": "value1",
		"key2": 123,
	}
	result = logger.WithFields(fields)
	if result == nil {
		t.Error("Expected logger but got nil")
	}
}

func TestLogger_LogHTTPRequest(t *testing.T) {
	cfg := &config.LoggingConfig{
		Level:  "info",
		Format: "json",
		Output: "stdout",
	}

	logger, err := NewLogger(cfg)
	if err != nil {
		t.Fatalf("Failed to create logger: %v", err)
	}

	ctx := context.WithValue(context.Background(), "request_id", "test-123")
	logger.LogHTTPRequest(ctx, "GET", "/test", "127.0.0.1", "test-agent", 200, 15)
}

func TestLogger_LogHTTPError(t *testing.T) {
	cfg := &config.LoggingConfig{
		Level:  "info",
		Format: "json",
		Output: "stdout",
	}

	logger, err := NewLogger(cfg)
	if err != nil {
		t.Fatalf("Failed to create logger: %v", err)
	}

	ctx := context.WithValue(context.Background(), "request_id", "test-123")
	logger.LogHTTPError(ctx, "GET", "/test", 500, ErrInvalidLogLevel)
}

func TestLogger_LogStorageOperation(t *testing.T) {
	cfg := &config.LoggingConfig{
		Level:  "info",
		Format: "json",
		Output: "stdout",
	}

	logger, err := NewLogger(cfg)
	if err != nil {
		t.Fatalf("Failed to create logger: %v", err)
	}

	ctx := context.WithValue(context.Background(), "request_id", "test-123")

	// Test successful operation
	logger.LogStorageOperation(ctx, "GET", "test-key", nil)

	// Test failed operation
	logger.LogStorageOperation(ctx, "SET", "test-key", ErrInvalidLogLevel)
}

func TestLogger_LogStartup(t *testing.T) {
	cfg := &config.LoggingConfig{
		Level:  "info",
		Format: "json",
		Output: "stdout",
	}

	logger, err := NewLogger(cfg)
	if err != nil {
		t.Fatalf("Failed to create logger: %v", err)
	}

	appConfig := &config.Config{
		Server: config.ServerConfig{
			Host: "0.0.0.0",
			Port: 8080,
		},
		Logging: *cfg,
		Storage: config.StorageConfig{
			Type: "memory",
		},
	}

	logger.LogStartup("1.0.0", "abc123", appConfig)
}

func TestLogger_LogShutdown(t *testing.T) {
	cfg := &config.LoggingConfig{
		Level:  "info",
		Format: "json",
		Output: "stdout",
	}

	logger, err := NewLogger(cfg)
	if err != nil {
		t.Fatalf("Failed to create logger: %v", err)
	}

	logger.LogShutdown("SIGTERM received")
}

func TestLogger_LogGracefulShutdown(t *testing.T) {
	cfg := &config.LoggingConfig{
		Level:  "info",
		Format: "json",
		Output: "stdout",
	}

	logger, err := NewLogger(cfg)
	if err != nil {
		t.Fatalf("Failed to create logger: %v", err)
	}

	logger.LogGracefulShutdown("stopping server", 5)
}

func TestLogger_LogLevelChange(t *testing.T) {
	cfg := &config.LoggingConfig{
		Level:  "info",
		Format: "json",
		Output: "stdout",
	}

	logger, err := NewLogger(cfg)
	if err != nil {
		t.Fatalf("Failed to create logger: %v", err)
	}

	logger.LogLevelChange("info", "debug")
}

func TestLogger_BasicLogging(t *testing.T) {
	cfg := &config.LoggingConfig{
		Level:  "debug",
		Format: "json",
		Output: "stdout",
	}

	logger, err := NewLogger(cfg)
	if err != nil {
		t.Fatalf("Failed to create logger: %v", err)
	}

	// Test all log levels
	logger.Debug("debug message", "key", "value")
	logger.Info("info message", "key", "value")
	logger.Warn("warn message", "key", "value")
	logger.Error("error message", "key", "value")
}

func TestNewLoggerWithDifferentLevels(t *testing.T) {
	levels := []string{"debug", "info", "warn", "error"}

	for _, level := range levels {
		t.Run(level, func(t *testing.T) {
			cfg := &config.LoggingConfig{
				Level:  level,
				Format: "json",
				Output: "stdout",
			}

			logger, err := NewLogger(cfg)
			if err != nil {
				t.Fatalf("Failed to create logger with level %s: %v", level, err)
			}

			if logger == nil {
				t.Error("Expected logger but got nil")
			}
		})
	}
}

func TestNewLoggerWithDifferentOutputs(t *testing.T) {
	outputs := []string{"stdout", "stderr", "invalid"}

	for _, output := range outputs {
		t.Run(output, func(t *testing.T) {
			cfg := &config.LoggingConfig{
				Level:  "info",
				Format: "json",
				Output: output,
			}

			logger, err := NewLogger(cfg)
			if err != nil {
				t.Fatalf("Failed to create logger with output %s: %v", output, err)
			}

			if logger == nil {
				t.Error("Expected logger but got nil")
			}
		})
	}
}