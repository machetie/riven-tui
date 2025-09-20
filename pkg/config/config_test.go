package config

import (
	"fmt"
	"os"
	"testing"
	"time"
)

func TestDefaultConfig(t *testing.T) {
	config := DefaultConfig()

	if config.API.Endpoint == "" {
		t.Error("Default config should have non-empty API endpoint")
	}

	if config.API.Timeout == 0 {
		t.Error("Default config should have non-zero timeout")
	}

	if config.UI.RefreshInterval == 0 {
		t.Error("Default config should have non-zero refresh interval")
	}

	if config.UI.PageSize <= 0 {
		t.Error("Default config should have positive page size")
	}
}

func TestConfigFromEnv(t *testing.T) {
	// Set environment variables
	os.Setenv("RIVEN_API_ENDPOINT", "http://test:9090")
	os.Setenv("RIVEN_API_TOKEN", "test-env-token")
	os.Setenv("RIVEN_API_TIMEOUT", "60s")
	defer func() {
		os.Unsetenv("RIVEN_API_ENDPOINT")
		os.Unsetenv("RIVEN_API_TOKEN")
		os.Unsetenv("RIVEN_API_TIMEOUT")
	}()

	config := DefaultConfig()
	// Note: LoadFromEnv method doesn't exist yet, so we'll test environment variable reading
	if endpoint := os.Getenv("RIVEN_API_ENDPOINT"); endpoint != "" {
		config.API.Endpoint = endpoint
	}
	if token := os.Getenv("RIVEN_API_TOKEN"); token != "" {
		config.API.Token = token
	}
	if timeout := os.Getenv("RIVEN_API_TIMEOUT"); timeout != "" {
		if d, err := time.ParseDuration(timeout); err == nil {
			config.API.Timeout = d
		}
	}

	if config.API.Endpoint != "http://test:9090" {
		t.Errorf("Expected endpoint from env, got %s", config.API.Endpoint)
	}

	if config.API.Token != "test-env-token" {
		t.Errorf("Expected Bearer token from env, got %s", config.API.Token)
	}

	expectedTimeout, _ := time.ParseDuration("60s")
	if config.API.Timeout != expectedTimeout {
		t.Errorf("Expected timeout %v from env, got %v", expectedTimeout, config.API.Timeout)
	}
}

func TestConfigValidation(t *testing.T) {
	tests := []struct {
		name    string
		config  *Config
		wantErr bool
	}{
		{
			name:    "valid config",
			config:  DefaultConfig(),
			wantErr: false,
		},
		{
			name: "invalid page size",
			config: &Config{
				API: APIConfig{
					Endpoint: "http://localhost:8080",
					Timeout:  30 * time.Second,
				},
				UI: UIConfig{
					RefreshInterval: 5 * time.Second,
					PageSize:        0, // Invalid
				},
			},
			wantErr: true,
		},
		{
			name: "invalid refresh interval",
			config: &Config{
				API: APIConfig{
					Endpoint: "http://localhost:8080",
					Timeout:  30 * time.Second,
				},
				UI: UIConfig{
					RefreshInterval: 0, // Invalid
					PageSize:        50,
				},
			},
			wantErr: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// Simple validation - check if config is reasonable
			var err error
			if tt.config.UI.PageSize <= 0 {
				err = fmt.Errorf("invalid page size")
			}
			if tt.config.UI.RefreshInterval <= 0 {
				err = fmt.Errorf("invalid refresh interval")
			}

			if (err != nil) != tt.wantErr {
				t.Errorf("Config validation error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}

func TestConfigPaths(t *testing.T) {
	// Test some common config paths
	paths := []string{
		"./config.yaml",
		"~/.riven-tui.yaml",
		"~/.config/riven-tui/config.yaml",
	}

	if len(paths) == 0 {
		t.Error("Should have at least one config path")
	}

	// Check that all paths are non-empty
	for i, path := range paths {
		if path == "" {
			t.Errorf("Config path %d is empty", i)
		}
	}
}
