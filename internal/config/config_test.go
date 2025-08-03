package config

import (
	"os"
	"testing"
	"time"
)

func TestLoad(t *testing.T) {
	// Test loading default configuration
	config, err := Load()
	if err != nil {
		t.Fatalf("Failed to load default config: %v", err)
	}

	// Verify default values
	if config.Server.Port != 8080 {
		t.Errorf("Expected server port 8080, got %d", config.Server.Port)
	}

	if config.Server.Host != "0.0.0.0" {
		t.Errorf("Expected server host 0.0.0.0, got %s", config.Server.Host)
	}

	if config.Logging.Level != "info" {
		t.Errorf("Expected logging level info, got %s", config.Logging.Level)
	}

	if config.Logging.Format != "json" {
		t.Errorf("Expected logging format json, got %s", config.Logging.Format)
	}

	if config.Storage.Type != "memory" {
		t.Errorf("Expected storage type memory, got %s", config.Storage.Type)
	}
}

func TestEnvironmentVariableOverrides(t *testing.T) {
	// Set environment variables
	os.Setenv("SERVER_PORT", "9090")
	os.Setenv("SERVER_HOST", "127.0.0.1")
	os.Setenv("LOGGING_LEVEL", "debug")
	os.Setenv("LOGGING_FORMAT", "text")
	os.Setenv("STORAGE_TYPE", "memory")
	os.Setenv("SERVER_READ_TIMEOUT", "60s")
	os.Setenv("SERVER_WRITE_TIMEOUT", "60s")
	os.Setenv("SERVER_IDLE_TIMEOUT", "120s")
	os.Setenv("DEBUG_METRICS_ENABLED", "true")
	os.Setenv("DEBUG_PPROF_ENABLED", "false")
	os.Setenv("CORS_MAX_AGE", "3600")

	defer func() {
		os.Unsetenv("SERVER_PORT")
		os.Unsetenv("SERVER_HOST")
		os.Unsetenv("LOGGING_LEVEL")
		os.Unsetenv("LOGGING_FORMAT")
		os.Unsetenv("STORAGE_TYPE")
		os.Unsetenv("SERVER_READ_TIMEOUT")
		os.Unsetenv("SERVER_WRITE_TIMEOUT")
		os.Unsetenv("SERVER_IDLE_TIMEOUT")
		os.Unsetenv("DEBUG_METRICS_ENABLED")
		os.Unsetenv("DEBUG_PPROF_ENABLED")
		os.Unsetenv("CORS_MAX_AGE")
	}()

	config, err := Load()
	if err != nil {
		t.Fatalf("Failed to load config with environment overrides: %v", err)
	}

	// Verify overridden values
	if config.Server.Port != 9090 {
		t.Errorf("Expected server port 9090, got %d", config.Server.Port)
	}

	if config.Server.Host != "127.0.0.1" {
		t.Errorf("Expected server host 127.0.0.1, got %s", config.Server.Host)
	}

	if config.Logging.Level != "debug" {
		t.Errorf("Expected logging level debug, got %s", config.Logging.Level)
	}

	if config.Logging.Format != "text" {
		t.Errorf("Expected logging format text, got %s", config.Logging.Format)
	}

	if config.Server.ReadTimeout != 60*time.Second {
		t.Errorf("Expected read timeout 60s, got %v", config.Server.ReadTimeout)
	}

	if config.Server.WriteTimeout != 60*time.Second {
		t.Errorf("Expected write timeout 60s, got %v", config.Server.WriteTimeout)
	}

	if config.Server.IdleTimeout != 120*time.Second {
		t.Errorf("Expected idle timeout 120s, got %v", config.Server.IdleTimeout)
	}

	if !config.Debug.Metrics.Enabled {
		t.Error("Expected metrics enabled true, got false")
	}

	if config.Debug.Pprof.Enabled {
		t.Error("Expected pprof enabled false, got true")
	}

	if config.CORS.MaxAge != 3600 {
		t.Errorf("Expected CORS max age 3600, got %d", config.CORS.MaxAge)
	}
}

func TestInvalidEnvironmentVariables(t *testing.T) {
	testCases := []struct {
		name        string
		envKey      string
		envValue    string
		expectError bool
	}{
		{"invalid port", "SERVER_PORT", "invalid", true},
		{"invalid read timeout", "SERVER_READ_TIMEOUT", "invalid", true},
		{"invalid write timeout", "SERVER_WRITE_TIMEOUT", "invalid", true},
		{"invalid idle timeout", "SERVER_IDLE_TIMEOUT", "invalid", true},
		{"invalid metrics enabled", "DEBUG_METRICS_ENABLED", "invalid", true},
		{"invalid pprof enabled", "DEBUG_PPROF_ENABLED", "invalid", true},
		{"invalid CORS max age", "CORS_MAX_AGE", "invalid", true},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			os.Setenv(tc.envKey, tc.envValue)
			defer os.Unsetenv(tc.envKey)

			_, err := Load()
			if tc.expectError && err == nil {
				t.Error("Expected error but got none")
			}
			if !tc.expectError && err != nil {
				t.Errorf("Unexpected error: %v", err)
			}
		})
	}
}

func TestCORSEnvironmentVariables(t *testing.T) {
	// Test CORS environment variables
	os.Setenv("CORS_ALLOWED_ORIGINS", "http://localhost:3000,https://example.com")
	os.Setenv("CORS_ALLOWED_METHODS", "GET,POST,PUT")
	os.Setenv("CORS_ALLOWED_HEADERS", "Content-Type,Authorization")

	defer func() {
		os.Unsetenv("CORS_ALLOWED_ORIGINS")
		os.Unsetenv("CORS_ALLOWED_METHODS")
		os.Unsetenv("CORS_ALLOWED_HEADERS")
	}()

	config, err := Load()
	if err != nil {
		t.Fatalf("Failed to load config with CORS environment variables: %v", err)
	}

	expectedOrigins := []string{"http://localhost:3000", "https://example.com"}
	if len(config.CORS.AllowedOrigins) != len(expectedOrigins) {
		t.Errorf("Expected %d allowed origins, got %d", len(expectedOrigins), len(config.CORS.AllowedOrigins))
	}

	expectedMethods := []string{"GET", "POST", "PUT"}
	if len(config.CORS.AllowedMethods) != len(expectedMethods) {
		t.Errorf("Expected %d allowed methods, got %d", len(expectedMethods), len(config.CORS.AllowedMethods))
	}

	expectedHeaders := []string{"Content-Type", "Authorization"}
	if len(config.CORS.AllowedHeaders) != len(expectedHeaders) {
		t.Errorf("Expected %d allowed headers, got %d", len(expectedHeaders), len(config.CORS.AllowedHeaders))
	}
}

func TestParseInt(t *testing.T) {
	testCases := []struct {
		input    string
		expected int
		hasError bool
	}{
		{"123", 123, false},
		{"0", 0, false},
		{"-1", -1, false},
		{"invalid", 0, true},
		{"", 0, true},
	}

	for _, tc := range testCases {
		t.Run(tc.input, func(t *testing.T) {
			result, err := parseInt(tc.input)
			if tc.hasError && err == nil {
				t.Error("Expected error but got none")
			}
			if !tc.hasError && err != nil {
				t.Errorf("Unexpected error: %v", err)
			}
			if !tc.hasError && result != tc.expected {
				t.Errorf("Expected %d, got %d", tc.expected, result)
			}
		})
	}
}

func TestParseBool(t *testing.T) {
	testCases := []struct {
		input    string
		expected bool
		hasError bool
	}{
		{"true", true, false},
		{"TRUE", true, false},
		{"True", true, false},
		{"1", true, false},
		{"yes", true, false},
		{"on", true, false},
		{"false", false, false},
		{"FALSE", false, false},
		{"False", false, false},
		{"0", false, false},
		{"no", false, false},
		{"off", false, false},
		{"invalid", false, true},
		{"", false, true},
	}

	for _, tc := range testCases {
		t.Run(tc.input, func(t *testing.T) {
			result, err := parseBool(tc.input)
			if tc.hasError && err == nil {
				t.Error("Expected error but got none")
			}
			if !tc.hasError && err != nil {
				t.Errorf("Unexpected error: %v", err)
			}
			if !tc.hasError && result != tc.expected {
				t.Errorf("Expected %t, got %t", tc.expected, result)
			}
		})
	}
}

func TestValidation(t *testing.T) {
	// Test with valid configuration
	config, err := Load()
	if err != nil {
		t.Fatalf("Failed to load valid config: %v", err)
	}

	// Test validation function directly
	err = validate(config)
	if err != nil {
		t.Errorf("Validation failed for valid config: %v", err)
	}
}

func TestInvalidConfigurations(t *testing.T) {
	// Test invalid port
	os.Setenv("SERVER_PORT", "70000") // Invalid port
	defer os.Unsetenv("SERVER_PORT")

	_, err := Load()
	if err == nil {
		t.Error("Expected error for invalid port but got none")
	}

	// Test invalid logging level
	os.Setenv("SERVER_PORT", "8080") // Reset to valid port
	os.Setenv("LOGGING_LEVEL", "invalid")
	defer os.Unsetenv("LOGGING_LEVEL")

	_, err = Load()
	if err == nil {
		t.Error("Expected error for invalid logging level but got none")
	}

	// Test invalid logging format
	os.Setenv("LOGGING_LEVEL", "info") // Reset to valid level
	os.Setenv("LOGGING_FORMAT", "invalid")
	defer os.Unsetenv("LOGGING_FORMAT")

	_, err = Load()
	if err == nil {
		t.Error("Expected error for invalid logging format but got none")
	}

	// Test invalid storage type
	os.Setenv("LOGGING_FORMAT", "json") // Reset to valid format
	os.Setenv("STORAGE_TYPE", "invalid")
	defer os.Unsetenv("STORAGE_TYPE")

	_, err = Load()
	if err == nil {
		t.Error("Expected error for invalid storage type but got none")
	}
}