package readme

import (
	"fmt"
	"os"
	"strings"
	"time"
)

const (
	MarkerStart = "<!-- CONTRIBEMBLEM:START -->"
	MarkerEnd   = "<!-- CONTRIBEMBLEM:END -->"
)

// Inject updates the README between markers with the emblem badge.
// If markers don't exist, they are appended to the end of the file.
// badgeImagePath is the relative path to the badge image (e.g., "badge.png").
// Returns true if the README content changed, false if unchanged.
func Inject(readmePath string, badgeImagePath string, updatedAt time.Time) (bool, error) {
	content, err := os.ReadFile(readmePath)
	if err != nil {
		return false, fmt.Errorf("failed to read README: %w", err)
	}

	injection := buildInjection(badgeImagePath, updatedAt)
	original := string(content)
	var updated string

	startIdx := strings.Index(original, MarkerStart)
	endIdx := strings.Index(original, MarkerEnd)

	if startIdx >= 0 && endIdx >= 0 && endIdx > startIdx {
		// Replace content between markers (preserve everything outside)
		before := original[:startIdx]
		after := original[endIdx+len(MarkerEnd):]
		updated = before + MarkerStart + "\n" + injection + "\n" + MarkerEnd + after
	} else {
		// No markers found — append to end
		updated = original
		if !strings.HasSuffix(updated, "\n") {
			updated += "\n"
		}
		updated += "\n" + MarkerStart + "\n" + injection + "\n" + MarkerEnd + "\n"
	}

	if updated == original {
		return false, nil
	}

	if err := os.WriteFile(readmePath, []byte(updated), 0644); err != nil {
		return false, fmt.Errorf("failed to write README: %w", err)
	}

	return true, nil
}

// buildInjection creates the markdown content to inject between markers.
func buildInjection(badgeImagePath string, updatedAt time.Time) string {
	timestamp := updatedAt.UTC().Format("January 2, 2006")
	return fmt.Sprintf("![ContribEmblem](%s)\n\n*Last updated: %s*", badgeImagePath, timestamp)
}
