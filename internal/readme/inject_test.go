package readme

import (
	"os"
	"strings"
	"testing"
	"time"
)

func TestInjectNoMarkers(t *testing.T) {
	// Create temp file without markers
	tmpFile, err := os.CreateTemp("", "readme_*.md")
	if err != nil {
		t.Fatalf("Failed to create temp file: %v", err)
	}
	defer os.Remove(tmpFile.Name())

	originalContent := "# My Project\n\nSome content here.\n"
	if _, err := tmpFile.WriteString(originalContent); err != nil {
		t.Fatalf("Failed to write temp file: %v", err)
	}
	tmpFile.Close()

	// Inject badge
	testTime := time.Date(2026, 2, 6, 12, 0, 0, 0, time.UTC)
	changed, err := Inject(tmpFile.Name(), "badge.png", testTime)
	if err != nil {
		t.Fatalf("Inject failed: %v", err)
	}

	if !changed {
		t.Error("Expected changed=true when markers don't exist")
	}

	// Read back and verify
	content, err := os.ReadFile(tmpFile.Name())
	if err != nil {
		t.Fatalf("Failed to read file: %v", err)
	}

	result := string(content)
	if !strings.Contains(result, MarkerStart) {
		t.Error("Expected MarkerStart in result")
	}
	if !strings.Contains(result, MarkerEnd) {
		t.Error("Expected MarkerEnd in result")
	}
	if !strings.Contains(result, "![ContribEmblem](badge.png)") {
		t.Error("Expected badge image markdown in result")
	}
	if !strings.Contains(result, "*Last updated: February 6, 2026*") {
		t.Error("Expected timestamp in result")
	}
	if !strings.HasPrefix(result, originalContent) {
		t.Error("Expected original content to be preserved at start")
	}
}

func TestInjectWithMarkers(t *testing.T) {
	// Create temp file with existing markers
	tmpFile, err := os.CreateTemp("", "readme_*.md")
	if err != nil {
		t.Fatalf("Failed to create temp file: %v", err)
	}
	defer os.Remove(tmpFile.Name())

	originalContent := `# My Project

Some content before.

<!-- CONTRIBEMBLEM:START -->
Old badge content here
<!-- CONTRIBEMBLEM:END -->

Some content after.
`
	if _, err := tmpFile.WriteString(originalContent); err != nil {
		t.Fatalf("Failed to write temp file: %v", err)
	}
	tmpFile.Close()

	// Inject new badge
	testTime := time.Date(2026, 2, 6, 12, 0, 0, 0, time.UTC)
	changed, err := Inject(tmpFile.Name(), "badge.png", testTime)
	if err != nil {
		t.Fatalf("Inject failed: %v", err)
	}

	if !changed {
		t.Error("Expected changed=true when content is different")
	}

	// Read back and verify
	content, err := os.ReadFile(tmpFile.Name())
	if err != nil {
		t.Fatalf("Failed to read file: %v", err)
	}

	result := string(content)
	if !strings.Contains(result, "Some content before.") {
		t.Error("Expected content before markers to be preserved")
	}
	if !strings.Contains(result, "Some content after.") {
		t.Error("Expected content after markers to be preserved")
	}
	if strings.Contains(result, "Old badge content here") {
		t.Error("Expected old badge content to be replaced")
	}
	if !strings.Contains(result, "![ContribEmblem](badge.png)") {
		t.Error("Expected new badge image markdown in result")
	}
	if !strings.Contains(result, "*Last updated: February 6, 2026*") {
		t.Error("Expected timestamp in result")
	}
}

func TestInjectUnchanged(t *testing.T) {
	// Create temp file with current content
	tmpFile, err := os.CreateTemp("", "readme_*.md")
	if err != nil {
		t.Fatalf("Failed to create temp file: %v", err)
	}
	defer os.Remove(tmpFile.Name())

	testTime := time.Date(2026, 2, 6, 12, 0, 0, 0, time.UTC)
	injection := buildInjection("badge.png", testTime)
	currentContent := "# My Project\n\n" + MarkerStart + "\n" + injection + "\n" + MarkerEnd + "\n"

	if _, err := tmpFile.WriteString(currentContent); err != nil {
		t.Fatalf("Failed to write temp file: %v", err)
	}
	tmpFile.Close()

	// Inject same content
	changed, err := Inject(tmpFile.Name(), "badge.png", testTime)
	if err != nil {
		t.Fatalf("Inject failed: %v", err)
	}

	if changed {
		t.Error("Expected changed=false when content is identical (idempotent)")
	}

	// Verify file wasn't modified
	content, err := os.ReadFile(tmpFile.Name())
	if err != nil {
		t.Fatalf("Failed to read file: %v", err)
	}

	if string(content) != currentContent {
		t.Error("Expected file content to remain unchanged")
	}
}

func TestBuildInjection(t *testing.T) {
	testTime := time.Date(2026, 2, 6, 12, 30, 45, 0, time.UTC)
	result := buildInjection("badge.png", testTime)

	expected := "![ContribEmblem](badge.png)\n\n*Last updated: February 6, 2026*"
	if result != expected {
		t.Errorf("Expected:\n%s\n\nGot:\n%s", expected, result)
	}
}

func TestBuildInjectionDifferentPath(t *testing.T) {
	testTime := time.Date(2025, 12, 31, 23, 59, 59, 0, time.UTC)
	result := buildInjection("assets/my-badge.png", testTime)

	expected := "![ContribEmblem](assets/my-badge.png)\n\n*Last updated: December 31, 2025*"
	if result != expected {
		t.Errorf("Expected:\n%s\n\nGot:\n%s", expected, result)
	}
}
