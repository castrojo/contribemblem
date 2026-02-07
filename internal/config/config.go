package config

import (
	"fmt"
	"os"

	"gopkg.in/yaml.v3"
)

const (
	DefaultConfigPath = "contribemblem.yml"
)

// Config represents the YAML configuration structure
type Config struct {
	// Username to fetch GitHub stats for
	Username string `yaml:"username"`

	// Metrics to display on the badge
	Metrics MetricsConfig `yaml:"metrics"`

	// Emblem rotation list
	Emblems EmblemsConfig `yaml:"emblems"`
}

// MetricsConfig defines which metrics to display
type MetricsConfig struct {
	Commits      bool `yaml:"commits"`
	PullRequests bool `yaml:"pull_requests"`
	Issues       bool `yaml:"issues"`
	Reviews      bool `yaml:"reviews"`
	Stars        bool `yaml:"stars"`
}

// EmblemsConfig defines emblem rotation settings
type EmblemsConfig struct {
	Rotation []string `yaml:"rotation"`
	Fallback string   `yaml:"fallback"`
}

// Load reads and parses the YAML configuration file
func Load(path string) (*Config, error) {
	// Check if file exists
	if _, err := os.Stat(path); os.IsNotExist(err) {
		return nil, fmt.Errorf("config file not found: %s", path)
	}

	data, err := os.ReadFile(path)
	if err != nil {
		return nil, fmt.Errorf("failed to read config file: %w", err)
	}

	var cfg Config
	if err := yaml.Unmarshal(data, &cfg); err != nil {
		return nil, fmt.Errorf("failed to parse YAML config: %w", err)
	}

	// Validate configuration
	if err := cfg.Validate(); err != nil {
		return nil, err
	}

	return &cfg, nil
}

// Validate checks that the configuration is valid
func (c *Config) Validate() error {
	// Username is required
	if c.Username == "" {
		return fmt.Errorf("username is required in config")
	}

	// At least one metric must be enabled
	if !c.Metrics.Commits && !c.Metrics.PullRequests && !c.Metrics.Issues &&
		!c.Metrics.Reviews && !c.Metrics.Stars {
		return fmt.Errorf("at least one metric must be enabled")
	}

	// Emblem rotation must have at least one emblem
	if len(c.Emblems.Rotation) == 0 {
		return fmt.Errorf("emblems.rotation must contain at least one emblem ID")
	}

	// Validate emblem IDs are non-empty
	for i, emblem := range c.Emblems.Rotation {
		if emblem == "" {
			return fmt.Errorf("emblems.rotation[%d] is empty", i)
		}
	}

	// Fallback emblem should be specified
	if c.Emblems.Fallback == "" {
		return fmt.Errorf("emblems.fallback is required")
	}

	return nil
}

// DefaultConfig returns a configuration with sensible defaults
func DefaultConfig() *Config {
	return &Config{
		Username: "",
		Metrics: MetricsConfig{
			Commits:      true,
			PullRequests: true,
			Issues:       true,
			Reviews:      true,
			Stars:        true,
		},
		Emblems: EmblemsConfig{
			Rotation: []string{
				"1538938257", // Seventh Column Projection
			},
			Fallback: "1538938257",
		},
	}
}
