package badge

import (
	"image"
	_ "image/png"
	"os"
	"path/filepath"
	"testing"
)

func TestGenerate(t *testing.T) {
	tests := []struct {
		name        string
		emblemPath  string
		stats       *Stats
		wantErr     bool
		wantWidth   int
		wantHeight  int
		description string
	}{
		{
			name:       "basic generation with sample stats",
			emblemPath: "testdata/test_emblem.jpg",
			stats: &Stats{
				Username:     "testuser",
				Commits:      150,
				PullRequests: 42,
				Issues:       18,
				Reviews:      67,
				Stars:        23,
			},
			wantErr:     false,
			wantWidth:   Width,
			wantHeight:  Height,
			description: "Should generate valid badge with all stats",
		},
		{
			name:       "empty username",
			emblemPath: "testdata/test_emblem.jpg",
			stats: &Stats{
				Username:     "",
				Commits:      100,
				PullRequests: 20,
				Issues:       10,
				Reviews:      30,
				Stars:        5,
			},
			wantErr:     false,
			wantWidth:   Width,
			wantHeight:  Height,
			description: "Should generate badge without username",
		},
		{
			name:       "zero stats",
			emblemPath: "testdata/test_emblem.jpg",
			stats: &Stats{
				Username:     "newbie",
				Commits:      0,
				PullRequests: 0,
				Issues:       0,
				Reviews:      0,
				Stars:        0,
			},
			wantErr:     false,
			wantWidth:   Width,
			wantHeight:  Height,
			description: "Should generate badge with zero power level",
		},
		{
			name:       "high stats values",
			emblemPath: "testdata/test_emblem.jpg",
			stats: &Stats{
				Username:     "veteran",
				Commits:      5000,
				PullRequests: 1200,
				Issues:       800,
				Reviews:      3400,
				Stars:        600,
			},
			wantErr:     false,
			wantWidth:   Width,
			wantHeight:  Height,
			description: "Should handle large stat values",
		},
		{
			name:       "missing emblem file",
			emblemPath: "testdata/nonexistent.jpg",
			stats: &Stats{
				Username: "testuser",
				Commits:  100,
			},
			wantErr:     true,
			description: "Should error when emblem file doesn't exist",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// Create temp output file
			tmpDir := t.TempDir()
			outputPath := filepath.Join(tmpDir, "badge.png")

			// Run Generate
			err := Generate(tt.emblemPath, tt.stats, outputPath)

			// Check error expectation
			if (err != nil) != tt.wantErr {
				t.Errorf("Generate() error = %v, wantErr %v", err, tt.wantErr)
				return
			}

			// If we expected an error, we're done
			if tt.wantErr {
				return
			}

			// Verify output file exists
			if _, err := os.Stat(outputPath); os.IsNotExist(err) {
				t.Errorf("Generate() did not create output file at %s", outputPath)
				return
			}

			// Verify PNG dimensions
			f, err := os.Open(outputPath)
			if err != nil {
				t.Fatalf("Failed to open generated badge: %v", err)
			}
			defer f.Close()

			img, _, err := image.Decode(f)
			if err != nil {
				t.Fatalf("Failed to decode generated badge as PNG: %v", err)
			}

			bounds := img.Bounds()
			if bounds.Dx() != tt.wantWidth {
				t.Errorf("Generated badge width = %d, want %d", bounds.Dx(), tt.wantWidth)
			}
			if bounds.Dy() != tt.wantHeight {
				t.Errorf("Generated badge height = %d, want %d", bounds.Dy(), tt.wantHeight)
			}

			// Verify dimensions match expected 800x162
			if bounds.Dx() != 800 || bounds.Dy() != 162 {
				t.Errorf("Generated badge dimensions = %dx%d, want 800x162", bounds.Dx(), bounds.Dy())
			}
		})
	}
}
