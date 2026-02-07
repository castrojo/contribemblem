package emblem

import (
	"os"
	"path/filepath"
	"testing"
)

func TestSelectEmblem(t *testing.T) {
	// Create temp config file
	tmpDir := t.TempDir()
	configPath := filepath.Join(tmpDir, "emblem-config.json")

	configJSON := `{"rotation":["1538938257","2962058744","2962058745"],"fallback":"1538938257"}`
	if err := os.WriteFile(configPath, []byte(configJSON), 0644); err != nil {
		t.Fatalf("Failed to write test config: %v", err)
	}

	emblem, err := SelectEmblem(configPath)
	if err != nil {
		t.Fatalf("SelectEmblem failed: %v", err)
	}

	// Verify emblem is one of the rotation values
	validEmblems := map[string]bool{
		"1538938257": true,
		"2962058744": true,
		"2962058745": true,
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

	configJSON := `{"rotation":[],"fallback":"1538938257"}`
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

	configJSON := `{"rotation":["1538938257","2962058744","2962058745"],"fallback":"1538938257"}`
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
