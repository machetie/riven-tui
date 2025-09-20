package api

import (
	"riven-tui/pkg/config"
	"testing"
	"time"
)

func TestNewClient(t *testing.T) {
	cfg := &config.Config{
		API: config.APIConfig{
			Endpoint: "http://localhost:8080",
			Token:    "test-token",
			Timeout:  30 * time.Second,
		},
	}

	client := NewClient(cfg)
	if client == nil {
		t.Fatal("NewClient returned nil")
	}

	if client.baseURL != cfg.API.Endpoint {
		t.Errorf("Expected baseURL %s, got %s", cfg.API.Endpoint, client.baseURL)
	}

	if client.token != cfg.API.Token {
		t.Errorf("Expected token %s, got %s", cfg.API.Token, client.token)
	}
}

func TestClientWithContext(t *testing.T) {
	cfg := &config.Config{
		API: config.APIConfig{
			Endpoint: "http://localhost:8080",
			Token:    "test-token",
			Timeout:  30 * time.Second,
		},
	}

	client := NewClient(cfg)

	// Test that context is properly handled (this won't make actual requests)
	if client.httpClient.Timeout != cfg.API.Timeout {
		t.Errorf("Expected timeout %v, got %v", cfg.API.Timeout, client.httpClient.Timeout)
	}
}

func TestClientConfigValidation(t *testing.T) {
	tests := []struct {
		name    string
		config  *config.Config
		wantErr bool
	}{
		{
			name: "valid config",
			config: &config.Config{
				API: config.APIConfig{
					Endpoint: "http://localhost:8080",
					Token:    "test-token",
					Timeout:  30 * time.Second,
				},
			},
			wantErr: false,
		},
		{
			name: "empty base URL",
			config: &config.Config{
				API: config.APIConfig{
					Endpoint: "",
					Token:    "test-token",
					Timeout:  30 * time.Second,
				},
			},
			wantErr: false, // Client should handle empty URL gracefully
		},
		{
			name: "zero timeout",
			config: &config.Config{
				API: config.APIConfig{
					Endpoint: "http://localhost:8080",
					Token:    "test-token",
					Timeout:  0,
				},
			},
			wantErr: false, // Should use default timeout
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			client := NewClient(tt.config)
			if client == nil && !tt.wantErr {
				t.Error("NewClient returned nil unexpectedly")
			}
			if client != nil && tt.wantErr {
				t.Error("NewClient should have failed but didn't")
			}
		})
	}
}
