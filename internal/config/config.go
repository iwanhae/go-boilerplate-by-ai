package config

import (
	_ "embed"
	"fmt"
	"os"
	"strings"
	"time"

	"gopkg.in/yaml.v3"
)

//go:embed defaults.yaml
var defaultConfig []byte

// Config represents the application configuration
type Config struct {
	Server  ServerConfig  `yaml:"server"`
	Logging LoggingConfig `yaml:"logging"`
	Storage StorageConfig `yaml:"storage"`
	Debug   DebugConfig   `yaml:"debug"`
	CORS    CORSConfig    `yaml:"cors"`
}

// ServerConfig represents server configuration
type ServerConfig struct {
	Host         string        `yaml:"host"`
	Port         int           `yaml:"port"`
	ReadTimeout  time.Duration `yaml:"readTimeout"`
	WriteTimeout time.Duration `yaml:"writeTimeout"`
	IdleTimeout  time.Duration `yaml:"idleTimeout"`
}

// LoggingConfig represents logging configuration
type LoggingConfig struct {
	Level  string `yaml:"level"`
	Format string `yaml:"format"`
	Output string `yaml:"output"`
}

// StorageConfig represents storage configuration
type StorageConfig struct {
	Type string `yaml:"type"`
}

// DebugConfig represents debug configuration
type DebugConfig struct {
	Metrics MetricsConfig `yaml:"metrics"`
	Pprof   PprofConfig   `yaml:"pprof"`
}

// MetricsConfig represents metrics configuration
type MetricsConfig struct {
	Enabled bool   `yaml:"enabled"`
	Path    string `yaml:"path"`
}

// PprofConfig represents pprof configuration
type PprofConfig struct {
	Enabled bool   `yaml:"enabled"`
	Path    string `yaml:"path"`
}

// CORSConfig represents CORS configuration
type CORSConfig struct {
	AllowedOrigins []string `yaml:"allowedOrigins"`
	AllowedMethods []string `yaml:"allowedMethods"`
	AllowedHeaders []string `yaml:"allowedHeaders"`
	MaxAge         int      `yaml:"maxAge"`
}

// Load loads configuration from defaults and environment variables
func Load() (*Config, error) {
	// Load default configuration
	config := &Config{}
	if err := yaml.Unmarshal(defaultConfig, config); err != nil {
		return nil, fmt.Errorf("failed to unmarshal default config: %w", err)
	}

	// Override with environment variables
	if err := overrideFromEnv(config); err != nil {
		return nil, fmt.Errorf("failed to override from environment: %w", err)
	}

	// Validate configuration
	if err := validate(config); err != nil {
		return nil, fmt.Errorf("configuration validation failed: %w", err)
	}

	return config, nil
}

// overrideFromEnv overrides configuration values from environment variables
func overrideFromEnv(config *Config) error {
	// Server configuration
	if port := os.Getenv("SERVER_PORT"); port != "" {
		if p, err := parseInt(port); err != nil {
			return fmt.Errorf("invalid SERVER_PORT: %w", err)
		} else {
			config.Server.Port = p
		}
	}

	if host := os.Getenv("SERVER_HOST"); host != "" {
		config.Server.Host = host
	}

	if readTimeout := os.Getenv("SERVER_READ_TIMEOUT"); readTimeout != "" {
		if rt, err := time.ParseDuration(readTimeout); err != nil {
			return fmt.Errorf("invalid SERVER_READ_TIMEOUT: %w", err)
		} else {
			config.Server.ReadTimeout = rt
		}
	}

	if writeTimeout := os.Getenv("SERVER_WRITE_TIMEOUT"); writeTimeout != "" {
		if wt, err := time.ParseDuration(writeTimeout); err != nil {
			return fmt.Errorf("invalid SERVER_WRITE_TIMEOUT: %w", err)
		} else {
			config.Server.WriteTimeout = wt
		}
	}

	if idleTimeout := os.Getenv("SERVER_IDLE_TIMEOUT"); idleTimeout != "" {
		if it, err := time.ParseDuration(idleTimeout); err != nil {
			return fmt.Errorf("invalid SERVER_IDLE_TIMEOUT: %w", err)
		} else {
			config.Server.IdleTimeout = it
		}
	}

	// Logging configuration
	if level := os.Getenv("LOGGING_LEVEL"); level != "" {
		config.Logging.Level = level
	}

	if format := os.Getenv("LOGGING_FORMAT"); format != "" {
		config.Logging.Format = format
	}

	if output := os.Getenv("LOGGING_OUTPUT"); output != "" {
		config.Logging.Output = output
	}

	// Storage configuration
	if storageType := os.Getenv("STORAGE_TYPE"); storageType != "" {
		config.Storage.Type = storageType
	}

	// Debug configuration
	if metricsEnabled := os.Getenv("DEBUG_METRICS_ENABLED"); metricsEnabled != "" {
		if enabled, err := parseBool(metricsEnabled); err != nil {
			return fmt.Errorf("invalid DEBUG_METRICS_ENABLED: %w", err)
		} else {
			config.Debug.Metrics.Enabled = enabled
		}
	}

	if metricsPath := os.Getenv("DEBUG_METRICS_PATH"); metricsPath != "" {
		config.Debug.Metrics.Path = metricsPath
	}

	if pprofEnabled := os.Getenv("DEBUG_PPROF_ENABLED"); pprofEnabled != "" {
		if enabled, err := parseBool(pprofEnabled); err != nil {
			return fmt.Errorf("invalid DEBUG_PPROF_ENABLED: %w", err)
		} else {
			config.Debug.Pprof.Enabled = enabled
		}
	}

	if pprofPath := os.Getenv("DEBUG_PPROF_PATH"); pprofPath != "" {
		config.Debug.Pprof.Path = pprofPath
	}

	// CORS configuration
	if allowedOrigins := os.Getenv("CORS_ALLOWED_ORIGINS"); allowedOrigins != "" {
		config.CORS.AllowedOrigins = strings.Split(allowedOrigins, ",")
	}

	if allowedMethods := os.Getenv("CORS_ALLOWED_METHODS"); allowedMethods != "" {
		config.CORS.AllowedMethods = strings.Split(allowedMethods, ",")
	}

	if allowedHeaders := os.Getenv("CORS_ALLOWED_HEADERS"); allowedHeaders != "" {
		config.CORS.AllowedHeaders = strings.Split(allowedHeaders, ",")
	}

	if maxAge := os.Getenv("CORS_MAX_AGE"); maxAge != "" {
		if ma, err := parseInt(maxAge); err != nil {
			return fmt.Errorf("invalid CORS_MAX_AGE: %w", err)
		} else {
			config.CORS.MaxAge = ma
		}
	}

	return nil
}

// validate validates the configuration
func validate(config *Config) error {
	// Server validation
	if config.Server.Port <= 0 || config.Server.Port > 65535 {
		return fmt.Errorf("invalid server port: %d", config.Server.Port)
	}

	if config.Server.ReadTimeout <= 0 {
		return fmt.Errorf("invalid read timeout: %v", config.Server.ReadTimeout)
	}

	if config.Server.WriteTimeout <= 0 {
		return fmt.Errorf("invalid write timeout: %v", config.Server.WriteTimeout)
	}

	if config.Server.IdleTimeout <= 0 {
		return fmt.Errorf("invalid idle timeout: %v", config.Server.IdleTimeout)
	}

	// Logging validation
	validLevels := map[string]bool{"debug": true, "info": true, "warn": true, "error": true}
	if !validLevels[config.Logging.Level] {
		return fmt.Errorf("invalid logging level: %s", config.Logging.Level)
	}

	validFormats := map[string]bool{"text": true, "json": true}
	if !validFormats[config.Logging.Format] {
		return fmt.Errorf("invalid logging format: %s", config.Logging.Format)
	}

	// Storage validation
	validStorageTypes := map[string]bool{"memory": true}
	if !validStorageTypes[config.Storage.Type] {
		return fmt.Errorf("invalid storage type: %s", config.Storage.Type)
	}

	return nil
}

// Helper functions for parsing environment variables
func parseInt(s string) (int, error) {
	var i int
	_, err := fmt.Sscanf(s, "%d", &i)
	return i, err
}

func parseBool(s string) (bool, error) {
	switch strings.ToLower(s) {
	case "true", "1", "yes", "on":
		return true, nil
	case "false", "0", "no", "off":
		return false, nil
	default:
		return false, fmt.Errorf("invalid boolean value: %s", s)
	}
}