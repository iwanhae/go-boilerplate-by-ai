package config

import (
	"bytes"
	_ "embed"
	"errors"
	"os"

	"gopkg.in/yaml.v3"
)

//go:embed defaults.yaml
var defaults []byte

// Config holds application configuration.
type Config struct {
	Server struct {
		Addr string `yaml:"addr"`
	} `yaml:"server"`
	Log struct {
		Level  string `yaml:"level"`
		Format string `yaml:"format"`
	} `yaml:"log"`
}

// Load reads configuration from embedded defaults and applies environment overrides.
func Load() (*Config, error) {
	var cfg Config
	if err := yaml.NewDecoder(bytes.NewReader(defaults)).Decode(&cfg); err != nil {
		return nil, err
	}

	if addr := os.Getenv("SERVER_ADDR"); addr != "" {
		cfg.Server.Addr = addr
	}
	if lvl := os.Getenv("LOG_LEVEL"); lvl != "" {
		cfg.Log.Level = lvl
	}
	if fmt := os.Getenv("LOG_FORMAT"); fmt != "" {
		cfg.Log.Format = fmt
	}

	if cfg.Server.Addr == "" {
		return nil, errors.New("server addr required")
	}
	return &cfg, nil
}
