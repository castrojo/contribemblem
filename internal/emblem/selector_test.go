package emblem

import (
	"os"
	"path/filepath"
	"testing"

	"github.com/castrojo/contribemblem/internal/config"
)

func TestSelectEmblem(t *testing.T) {
	// Create temp config file
	tmpDir := t.TempDir()
	configPath := filepath.Join(tmpDir, "emblem-config.json")

	configJSON := `{"rotation":["4052831236","1901885391","1661191194"],"fallback":"4052831236"}`
	if err := os.WriteFile(configPath, []byte(configJSON), 0644); err != nil {
		t.Fatalf("Failed to write test config: %v", err)
	}

	emblem, err := SelectEmblem(configPath)
	if err != nil {
		t.Fatalf("SelectEmblem failed: %v", err)
	}

	// Verify emblem is one of the rotation values
	validEmblems := map[string]bool{
		"4052831236": true,
		"1901885391": true,
		"1661191194": true,
	}
	if !validEmblems[emblem] {
		t.Errorf("Selected emblem %s not in rotation", emblem)
	}
}

func TestSelectEmblemMissingConfig(t *testing.T) {
	emblem, err := SelectEmblem("/nonexistent/config.json")
	if err != nil {
		t.Fatalf("Expected no error on missing config, got %v", err)
	}
	if emblem != FallbackEmblem {
		t.Errorf("Expected fallback emblem %s, got %s", FallbackEmblem, emblem)
	}
}

func TestSelectEmblemEmptyRotation(t *testing.T) {
	tmpDir := t.TempDir()
	configPath := filepath.Join(tmpDir, "emblem-config.json")

	configJSON := `{"rotation":[],"fallback":"4052831236"}`
	if err := os.WriteFile(configPath, []byte(configJSON), 0644); err != nil {
		t.Fatalf("Failed to write test config: %v", err)
	}

	emblem, err := SelectEmblem(configPath)
	if err != nil {
		t.Fatalf("Expected no error on empty rotation, got %v", err)
	}
	if emblem != FallbackEmblem {
		t.Errorf("Expected fallback emblem %s, got %s", FallbackEmblem, emblem)
	}
}

func TestSelectEmblemDeterministic(t *testing.T) {
	tmpDir := t.TempDir()
	configPath := filepath.Join(tmpDir, "emblem-config.json")

	configJSON := `{"rotation":["4052831236","1901885391","1661191194"],"fallback":"4052831236"}`
	if err := os.WriteFile(configPath, []byte(configJSON), 0644); err != nil {
		t.Fatalf("Failed to write test config: %v", err)
	}

	// Select emblem twice - should be identical
	emblem1, _ := SelectEmblem(configPath)
	emblem2, _ := SelectEmblem(configPath)

	if emblem1 != emblem2 {
		t.Errorf("Selection not deterministic: %s != %s", emblem1, emblem2)
	}
}

func TestSelectEmblemFromConfig(t *testing.T) {
	cfg := &config.EmblemsConfig{
		Rotation: []string{"4052831236", "1901885391", "1661191194"},
		Fallback: "4052831236",
	}

	emblem, err := SelectEmblemFromConfig(cfg)
	if err != nil {
		t.Fatalf("SelectEmblemFromConfig failed: %v", err)
	}

	// Verify emblem is one of the rotation values
	validEmblems := map[string]bool{
		"4052831236": true,
		"1901885391": true,
		"1661191194": true,
	}
	if !validEmblems[emblem] {
		t.Errorf("Selected emblem %s not in rotation", emblem)
	}
}

func TestSelectEmblemFromConfigNilConfig(t *testing.T) {
	emblem, err := SelectEmblemFromConfig(nil)
	if err != nil {
		t.Fatalf("Expected no error on nil config, got %v", err)
	}
	if emblem != FallbackEmblem {
		t.Errorf("Expected fallback emblem %s, got %s", FallbackEmblem, emblem)
	}
}

func TestSelectEmblemFromConfigEmptyRotation(t *testing.T) {
	cfg := &config.EmblemsConfig{
		Rotation: []string{},
		Fallback: "1234567890",
	}

	emblem, err := SelectEmblemFromConfig(cfg)
	if err != nil {
		t.Fatalf("Expected no error on empty rotation, got %v", err)
	}
	if emblem != "1234567890" {
		t.Errorf("Expected custom fallback emblem 1234567890, got %s", emblem)
	}
}

func TestSelectEmblemFromConfigDeterministic(t *testing.T) {
	cfg := &config.EmblemsConfig{
		Rotation: []string{"4052831236", "1901885391", "1661191194"},
		Fallback: "4052831236",
	}

	// Select emblem twice - should be identical
	emblem1, _ := SelectEmblemFromConfig(cfg)
	emblem2, _ := SelectEmblemFromConfig(cfg)

	if emblem1 != emblem2 {
		t.Errorf("Selection not deterministic: %s != %s", emblem1, emblem2)
	}
}
