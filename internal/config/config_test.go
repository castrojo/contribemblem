package config

import (
	"os"
	"path/filepath"
	"testing"
)

func TestLoadValidConfig(t *testing.T) {
	// Create temp config file
	tmpDir := t.TempDir()
	configPath := filepath.Join(tmpDir, "config.yml")

	configContent := `username: testuser
metrics:
  commits: true
  pull_requests: true
  issues: false
  reviews: true
  stars: false
emblems:
  rotation:
    - "1538938257"
    - "1409726931"
  fallback: "1538938257"
`
	if err := os.WriteFile(configPath, []byte(configContent), 0644); err != nil {
		t.Fatalf("Failed to write test config: %v", err)
	}

	// Load config
	cfg, err := Load(configPath)
	if err != nil {
		t.Fatalf("Failed to load config: %v", err)
	}

	// Verify username
	if cfg.Username != "testuser" {
		t.Errorf("Expected username 'testuser', got '%s'", cfg.Username)
	}

	// Verify metrics
	if !cfg.Metrics.Commits {
		t.Error("Expected commits to be enabled")
	}
	if cfg.Metrics.Issues {
		t.Error("Expected issues to be disabled")
	}

	// Verify emblems
	if len(cfg.Emblems.Rotation) != 2 {
		t.Errorf("Expected 2 emblems in rotation, got %d", len(cfg.Emblems.Rotation))
	}
	if cfg.Emblems.Fallback != "1538938257" {
		t.Errorf("Expected fallback '1538938257', got '%s'", cfg.Emblems.Fallback)
	}
}

func TestLoadMissingFile(t *testing.T) {
	_, err := Load("/nonexistent/config.yml")
	if err == nil {
		t.Error("Expected error for missing config file")
	}
}

func TestLoadInvalidYAML(t *testing.T) {
	tmpDir := t.TempDir()
	configPath := filepath.Join(tmpDir, "config.yml")

	invalidContent := `username: testuser
metrics:
  commits: true
  indentation error here
`
	if err := os.WriteFile(configPath, []byte(invalidContent), 0644); err != nil {
		t.Fatalf("Failed to write test config: %v", err)
	}

	_, err := Load(configPath)
	if err == nil {
		t.Error("Expected error for invalid YAML")
	}
}

func TestValidateMissingUsername(t *testing.T) {
	cfg := DefaultConfig()
	cfg.Username = ""

	err := cfg.Validate()
	if err == nil {
		t.Error("Expected validation error for missing username")
	}
}

func TestValidateNoMetricsEnabled(t *testing.T) {
	cfg := DefaultConfig()
	cfg.Username = "testuser"
	cfg.Metrics = MetricsConfig{
		Commits:      false,
		PullRequests: false,
		Issues:       false,
		Reviews:      false,
		Stars:        false,
	}

	err := cfg.Validate()
	if err == nil {
		t.Error("Expected validation error when no metrics enabled")
	}
}

func TestValidateEmptyEmblemRotation(t *testing.T) {
	cfg := DefaultConfig()
	cfg.Username = "testuser"
	cfg.Emblems.Rotation = []string{}

	err := cfg.Validate()
	if err == nil {
		t.Error("Expected validation error for empty emblem rotation")
	}
}

func TestValidateEmptyEmblemID(t *testing.T) {
	cfg := DefaultConfig()
	cfg.Username = "testuser"
	cfg.Emblems.Rotation = []string{"1538938257", "", "1409726931"}

	err := cfg.Validate()
	if err == nil {
		t.Error("Expected validation error for empty emblem ID")
	}
}

func TestValidateMissingFallback(t *testing.T) {
	cfg := DefaultConfig()
	cfg.Username = "testuser"
	cfg.Emblems.Fallback = ""

	err := cfg.Validate()
	if err == nil {
		t.Error("Expected validation error for missing fallback emblem")
	}
}

func TestDefaultConfig(t *testing.T) {
	cfg := DefaultConfig()

	// Default should have all metrics enabled
	if !cfg.Metrics.Commits || !cfg.Metrics.PullRequests || !cfg.Metrics.Issues ||
		!cfg.Metrics.Reviews || !cfg.Metrics.Stars {
		t.Error("Expected all metrics enabled in default config")
	}

	// Default should have at least one emblem
	if len(cfg.Emblems.Rotation) == 0 {
		t.Error("Expected at least one emblem in default rotation")
	}

	// Default should have fallback
	if cfg.Emblems.Fallback == "" {
		t.Error("Expected fallback emblem in default config")
	}
}
