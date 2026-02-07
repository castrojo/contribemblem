package emblem

import (
	"crypto/sha256"
	"encoding/binary"
	"encoding/json"
	"fmt"
	"os"
	"time"

	"github.com/castrojo/contribemblem/internal/config"
)

const (
	FallbackEmblem    = "1538938257" // "Seventh Column Projection" emblem
	DefaultConfigPath = "data/emblem-config.json"
)

// Config represents emblem-config.json structure
type Config struct {
	Rotation []string `json:"rotation"`
	Fallback string   `json:"fallback"`
}

// SelectEmblem performs deterministic weekly emblem selection
// Uses ISO week number and SHA256-seeded modulo selection
func SelectEmblem(configPath string) (string, error) {
	// Load config
	config, err := loadConfig(configPath)
	if err != nil || len(config.Rotation) == 0 {
		// Use fallback if config missing or rotation empty
		return FallbackEmblem, nil
	}

	// Calculate ISO week in UTC (format: YYYY-Www, e.g., "2026-W06")
	now := time.Now().UTC()
	year, week := now.ISOWeek()
	isoWeek := fmt.Sprintf("%d-W%02d", year, week)

	// Generate SHA256 hash of ISO week string
	hash := sha256.Sum256([]byte(isoWeek))

	// Convert first 8 bytes of hash to uint64 (deterministic seed)
	hashValue := binary.BigEndian.Uint64(hash[:8])

	// Select emblem using modulo operation
	index := int(hashValue % uint64(len(config.Rotation)))
	selectedEmblem := config.Rotation[index]

	return selectedEmblem, nil
}

func loadConfig(path string) (*Config, error) {
	data, err := os.ReadFile(path)
	if err != nil {
		return nil, fmt.Errorf("failed to read config: %w", err)
	}

	var config Config
	if err := json.Unmarshal(data, &config); err != nil {
		return nil, fmt.Errorf("failed to parse config: %w", err)
	}

	return &config, nil
}

// SelectEmblemFromConfig performs deterministic weekly emblem selection
// from a config.EmblemsConfig structure
func SelectEmblemFromConfig(cfg *config.EmblemsConfig) (string, error) {
	if cfg == nil || len(cfg.Rotation) == 0 {
		// Use fallback if config missing or rotation empty
		if cfg != nil && cfg.Fallback != "" {
			return cfg.Fallback, nil
		}
		return FallbackEmblem, nil
	}

	// Calculate ISO week in UTC (format: YYYY-Www, e.g., "2026-W06")
	now := time.Now().UTC()
	year, week := now.ISOWeek()
	isoWeek := fmt.Sprintf("%d-W%02d", year, week)

	// Generate SHA256 hash of ISO week string
	hash := sha256.Sum256([]byte(isoWeek))

	// Convert first 8 bytes of hash to uint64 (deterministic seed)
	hashValue := binary.BigEndian.Uint64(hash[:8])

	// Select emblem using modulo operation
	index := int(hashValue % uint64(len(cfg.Rotation)))
	selectedEmblem := cfg.Rotation[index]

	return selectedEmblem, nil
}
