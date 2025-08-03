package config

import (
	"os"
	"testing"
)

func TestLoadDefaults(t *testing.T) {
	t.Setenv("SERVER_ADDR", "")
	cfg, err := Load()
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if cfg.Server.Addr != ":8080" {
		t.Fatalf("addr %q", cfg.Server.Addr)
	}
	if cfg.Log.Level != "info" {
		t.Fatalf("level %q", cfg.Log.Level)
	}
	if cfg.Log.Format != "text" {
		t.Fatalf("format %q", cfg.Log.Format)
	}
}

func TestEnvOverride(t *testing.T) {
	t.Setenv("SERVER_ADDR", ":9090")
	t.Setenv("LOG_LEVEL", "debug")
	t.Setenv("LOG_FORMAT", "json")
	cfg, err := Load()
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if cfg.Server.Addr != ":9090" || cfg.Log.Level != "debug" || cfg.Log.Format != "json" {
		t.Fatalf("env override failed: %#v", cfg)
	}
}

func TestMissingAddr(t *testing.T) {
	// Simulate empty default by temporarily clearing defaults variable.
	orig := defaults
	defaults = []byte("server:\n  addr: ''\n")
	t.Cleanup(func() { defaults = orig })
	os.Unsetenv("SERVER_ADDR")
	if _, err := Load(); err == nil {
		t.Fatal("expected error")
	}
}
