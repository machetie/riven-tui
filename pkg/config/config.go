package config

import (
	"fmt"
	"os"
	"path/filepath"
	"time"

	"gopkg.in/yaml.v3"
)

// Config represents the application configuration
type Config struct {
	API APIConfig `yaml:"api"`
	UI  UIConfig  `yaml:"ui"`
}

// APIConfig represents API-related configuration
type APIConfig struct {
	Endpoint string        `yaml:"endpoint"`
	Token    string        `yaml:"token" json:"token"`
	Timeout  time.Duration `yaml:"timeout"`
}

// UIConfig represents UI-related configuration
type UIConfig struct {
	RefreshInterval time.Duration `yaml:"refresh_interval"`
	Theme           string        `yaml:"theme"`
	PageSize        int           `yaml:"page_size"`
}

// DefaultConfig returns a configuration with default values
func DefaultConfig() *Config {
	return &Config{
		API: APIConfig{
			Endpoint: "http://localhost:8080",
			Token:    "",
			Timeout:  30 * time.Second,
		},
		UI: UIConfig{
			RefreshInterval: 5 * time.Second,
			Theme:           "default",
			PageSize:        50,
		},
	}
}

// LoadConfig loads configuration from file and environment variables
func LoadConfig(configPath string) (*Config, error) {
	config := DefaultConfig()

	// Load from file if it exists
	if configPath != "" {
		if err := loadFromFile(config, configPath); err != nil {
			return nil, fmt.Errorf("failed to load config from file: %w", err)
		}
	} else {
		// Try to load from default locations
		defaultPaths := []string{
			filepath.Join(os.Getenv("HOME"), ".config", "riven-tui", "config.yaml"),
			filepath.Join(os.Getenv("HOME"), ".riven-tui.yaml"),
			"config.yaml",
		}

		for _, path := range defaultPaths {
			if _, err := os.Stat(path); err == nil {
				if err := loadFromFile(config, path); err != nil {
					return nil, fmt.Errorf("failed to load config from %s: %w", path, err)
				}
				break
			}
		}
	}

	// Override with environment variables
	loadFromEnv(config)

	// Validate configuration
	if err := validateConfig(config); err != nil {
		return nil, fmt.Errorf("invalid configuration: %w", err)
	}

	return config, nil
}

// loadFromFile loads configuration from a YAML file
func loadFromFile(config *Config, path string) error {
	data, err := os.ReadFile(path)
	if err != nil {
		return err
	}

	// First, check for legacy api_key field
	var legacyCheck struct {
		API struct {
			APIKey string `yaml:"api_key"`
		} `yaml:"api"`
	}

	if err := yaml.Unmarshal(data, &legacyCheck); err == nil {
		if legacyCheck.API.APIKey != "" {
			return fmt.Errorf("configuration uses deprecated 'api_key' field. Please update your config file to use 'token' instead of 'api_key'. See README.md for migration instructions")
		}
	}

	return yaml.Unmarshal(data, config)
}

// loadFromEnv loads configuration from environment variables
func loadFromEnv(config *Config) {
	if endpoint := os.Getenv("RIVEN_API_ENDPOINT"); endpoint != "" {
		config.API.Endpoint = endpoint
	}

	// Support both new and legacy environment variable names
	if token := os.Getenv("RIVEN_API_TOKEN"); token != "" {
		config.API.Token = token
	} else if apiKey := os.Getenv("RIVEN_API_KEY"); apiKey != "" {
		// Legacy support with deprecation warning
		fmt.Fprintf(os.Stderr, "Warning: RIVEN_API_KEY is deprecated, please use RIVEN_API_TOKEN instead\n")
		config.API.Token = apiKey
	}

	if timeout := os.Getenv("RIVEN_API_TIMEOUT"); timeout != "" {
		if d, err := time.ParseDuration(timeout); err == nil {
			config.API.Timeout = d
		}
	}

	if refreshInterval := os.Getenv("RIVEN_UI_REFRESH_INTERVAL"); refreshInterval != "" {
		if d, err := time.ParseDuration(refreshInterval); err == nil {
			config.UI.RefreshInterval = d
		}
	}

	if theme := os.Getenv("RIVEN_UI_THEME"); theme != "" {
		config.UI.Theme = theme
	}
}

// validateConfig validates the configuration
func validateConfig(config *Config) error {
	if config.API.Endpoint == "" {
		return fmt.Errorf("API endpoint is required")
	}

	if config.API.Token == "" {
		return fmt.Errorf("Bearer token is required")
	}

	if config.API.Timeout <= 0 {
		return fmt.Errorf("API timeout must be positive")
	}

	if config.UI.RefreshInterval <= 0 {
		return fmt.Errorf("UI refresh interval must be positive")
	}

	if config.UI.PageSize <= 0 {
		config.UI.PageSize = 50 // Set default if invalid
	}

	return nil
}

// SaveConfig saves the configuration to a file
func SaveConfig(config *Config, path string) error {
	// Create directory if it doesn't exist
	dir := filepath.Dir(path)
	if err := os.MkdirAll(dir, 0755); err != nil {
		return fmt.Errorf("failed to create config directory: %w", err)
	}

	data, err := yaml.Marshal(config)
	if err != nil {
		return fmt.Errorf("failed to marshal config: %w", err)
	}

	return os.WriteFile(path, data, 0644)
}

// GetDefaultConfigPath returns the default configuration file path
func GetDefaultConfigPath() string {
	return filepath.Join(os.Getenv("HOME"), ".config", "riven-tui", "config.yaml")
}
