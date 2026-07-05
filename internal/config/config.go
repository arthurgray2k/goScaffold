package config

import (
	"errors"
	"os"
	"path/filepath"

	"gopkg.in/yaml.v3"
)

// Config represents the application configuration.
type Config struct {
	DefaultAuthor       string `yaml:"default_author"`
	DefaultLicense      string `yaml:"default_license"`
	DefaultModulePrefix string `yaml:"default_module_prefix"`
}

// Load reads the configuration from ~/.goscaffold.yaml.
// If the file does not exist, it returns an empty config without an error.
func Load() (*Config, error) {
	homeDir, err := os.UserHomeDir()
	if err != nil {
		return nil, err
	}

	configPath := filepath.Join(homeDir, ".goscaffold.yaml")
	data, err := os.ReadFile(configPath)
	if err != nil {
		if errors.Is(err, os.ErrNotExist) {
			return &Config{}, nil // Return empty config if file is missing
		}
		return nil, err
	}

	var cfg Config
	if err := yaml.Unmarshal(data, &cfg); err != nil {
		return nil, err
	}

	return &cfg, nil
}
