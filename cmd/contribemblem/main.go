package main

import (
	"encoding/json"
	"fmt"
	"os"
	"time"

	"github.com/castrojo/contribemblem/internal/badge"
	"github.com/castrojo/contribemblem/internal/bungie"
	"github.com/castrojo/contribemblem/internal/config"
	"github.com/castrojo/contribemblem/internal/emblem"
	"github.com/castrojo/contribemblem/internal/github"
	"github.com/castrojo/contribemblem/internal/readme"
)

func main() {
	if len(os.Args) < 2 {
		printUsage()
		os.Exit(1)
	}

	// Load config from contribemblem.yml (fallback to env vars if not found)
	cfg, err := config.Load(config.DefaultConfigPath)
	if err != nil {
		// Config file not found or invalid - this is OK for backwards compatibility
		// Commands will fall back to environment variables (GITHUB_ACTOR)
		fmt.Fprintf(os.Stderr, "Note: Config file not loaded (%v), using environment variables\n", err)
		cfg = nil
	}

	cmd := os.Args[1]
	switch cmd {
	case "fetch-stats":
		username := getUsername(cfg)
		stats, err := github.FetchStats(username, nil)
		if err != nil {
			fmt.Fprintf(os.Stderr, "Error: %v\n", err)
			os.Exit(1)
		}
		data, _ := json.MarshalIndent(stats, "", "  ")
		fmt.Println(string(data))
	case "select-emblem":
		// Try config first, fall back to JSON file
		var selectedEmblem string
		if cfg != nil {
			selectedEmblem, err = emblem.SelectEmblemFromConfig(&cfg.Emblems)
		} else {
			selectedEmblem, err = emblem.SelectEmblem(emblem.DefaultConfigPath)
		}
		if err != nil {
			fmt.Fprintf(os.Stderr, "Error: %v\n", err)
			os.Exit(1)
		}
		fmt.Println(selectedEmblem)
	case "fetch-emblem":
		// Read emblem hash from args or stdin
		var emblemHash string
		if len(os.Args) > 2 {
			emblemHash = os.Args[2]
		} else {
			if _, err := fmt.Scanln(&emblemHash); err != nil {
				fmt.Fprintf(os.Stderr, "Error reading emblem hash: %v\n", err)
				os.Exit(1)
			}
		}
		if err := bungie.FetchEmblem(emblemHash); err != nil {
			fmt.Fprintf(os.Stderr, "Error: %v\n", err)
			os.Exit(1)
		}
	case "generate":
		// Read stats from data/stats.json
		statsData, err := os.ReadFile("data/stats.json")
		if err != nil {
			fmt.Fprintf(os.Stderr, "Error reading stats.json: %v\n", err)
			os.Exit(1)
		}

		var ghStats github.Stats
		if err := json.Unmarshal(statsData, &ghStats); err != nil {
			fmt.Fprintf(os.Stderr, "Error parsing stats.json: %v\n", err)
			os.Exit(1)
		}

		// Convert to badge.Stats
		badgeStats := &badge.Stats{
			Username:     getUsername(cfg),
			Commits:      ghStats.Commits,
			PullRequests: ghStats.PullRequests,
			Issues:       ghStats.Issues,
			Reviews:      ghStats.Reviews,
			Stars:        ghStats.StarsReceived,
		}

		// Generate badge
		if err := badge.Generate("data/emblem.jpg", badgeStats, "badge.png"); err != nil {
			fmt.Fprintf(os.Stderr, "Error generating badge: %v\n", err)
			os.Exit(1)
		}

		fmt.Println("✓ Badge generated: badge.png")
	case "update-readme":
		changed, err := readme.Inject("README.md", "badge.png", time.Now())
		if err != nil {
			fmt.Fprintf(os.Stderr, "Error updating README: %v\n", err)
			os.Exit(1)
		}
		if changed {
			fmt.Println("✓ README updated")
		} else {
			fmt.Println("✓ README already current")
		}
	case "run":
		fmt.Println("Running full ContribEmblem pipeline...")

		// Step 1: Fetch GitHub stats
		fmt.Println("[1/5] Fetching GitHub stats...")
		username := getUsername(cfg)
		stats, err := github.FetchStats(username, nil)
		if err != nil {
			fmt.Fprintf(os.Stderr, "Failed to fetch stats: %v\n", err)
			os.Exit(1)
		}

		// Save stats to data/stats.json
		if err := os.MkdirAll("data", 0755); err != nil {
			fmt.Fprintf(os.Stderr, "Failed to create data directory: %v\n", err)
			os.Exit(1)
		}
		statsJSON, _ := json.MarshalIndent(stats, "", "  ")
		if err := os.WriteFile("data/stats.json", statsJSON, 0644); err != nil {
			fmt.Fprintf(os.Stderr, "Failed to write stats.json: %v\n", err)
			os.Exit(1)
		}
		fmt.Println("✓ Stats saved to data/stats.json")

		// Step 2: Select emblem
		fmt.Println("[2/5] Selecting weekly emblem...")
		var emblemHash string
		if cfg != nil {
			emblemHash, err = emblem.SelectEmblemFromConfig(&cfg.Emblems)
		} else {
			emblemHash, err = emblem.SelectEmblem(emblem.DefaultConfigPath)
		}
		if err != nil {
			fmt.Fprintf(os.Stderr, "Failed to select emblem: %v\n", err)
			os.Exit(1)
		}
		fmt.Printf("✓ Selected emblem: %s\n", emblemHash)

		// Step 3: Fetch emblem from Bungie
		fmt.Println("[3/5] Fetching emblem from Bungie API...")
		if err := bungie.FetchEmblem(emblemHash); err != nil {
			fmt.Fprintf(os.Stderr, "Failed to fetch emblem: %v\n", err)
			os.Exit(1)
		}
		fmt.Println("✓ Emblem downloaded to data/emblem.jpg")

		// Step 4: Generate badge
		fmt.Println("[4/5] Generating badge image...")
		badgeStats := &badge.Stats{
			Username:     getUsername(cfg),
			Commits:      stats.Commits,
			PullRequests: stats.PullRequests,
			Issues:       stats.Issues,
			Reviews:      stats.Reviews,
			Stars:        stats.StarsReceived,
		}
		if err := badge.Generate("data/emblem.jpg", badgeStats, "badge.png"); err != nil {
			fmt.Fprintf(os.Stderr, "Failed to generate badge: %v\n", err)
			os.Exit(1)
		}
		fmt.Println("✓ Badge generated: badge.png")

		// Step 5: Update README
		fmt.Println("[5/5] Updating README...")
		changed, err := readme.Inject("README.md", "badge.png", time.Now())
		if err != nil {
			fmt.Fprintf(os.Stderr, "Failed to update README: %v\n", err)
			os.Exit(1)
		}
		if changed {
			fmt.Println("✓ README updated")
		} else {
			fmt.Println("✓ README already current")
		}

		fmt.Println("\n🎉 Pipeline complete! Badge ready at badge.png")
	case "generate-demos":
		if err := generateDemos(); err != nil {
			fmt.Fprintf(os.Stderr, "Error: %v\n", err)
			os.Exit(1)
		}
	case "help", "--help", "-h":
		printUsage()
		os.Exit(0)
	default:
		fmt.Fprintf(os.Stderr, "Unknown command: %s\n", cmd)
		printUsage()
		os.Exit(1)
	}
}

// getUsername returns the username from config, falling back to GITHUB_ACTOR env var
func getUsername(cfg *config.Config) string {
	if cfg != nil && cfg.Username != "" {
		return cfg.Username
	}
	return os.Getenv("GITHUB_ACTOR")
}

type demoUser struct {
	username   string
	emblemHash string
	commits    int
	prs        int
	issues     int
	reviews    int
	stars      int
}

var demoUsers = []demoUser{
	{"castrojo", "4052831236", 842, 156, 89, 234, 1247},
	{"jeefy", "1901885391", 567, 123, 67, 189, 892},
	{"mrbobbytables", "1661191194", 423, 98, 156, 312, 634},
}

func generateDemos() error {
	// Validate BUNGIE_API_KEY
	if os.Getenv("BUNGIE_API_KEY") == "" {
		return fmt.Errorf("BUNGIE_API_KEY environment variable not set\nGet your API key from https://www.bungie.net/en/Application")
	}

	// Create directories
	os.MkdirAll("examples", 0755)
	os.MkdirAll("data", 0755)

	for _, user := range demoUsers {
		fmt.Printf("\n=== Generating badge for @%s ===\n", user.username)

		// Write stats JSON
		stats := github.Stats{
			Year:          time.Now().Year(),
			UpdatedAt:     time.Now().Format(time.RFC3339),
			Commits:       user.commits,
			PullRequests:  user.prs,
			Issues:        user.issues,
			Reviews:       user.reviews,
			StarsReceived: user.stars,
		}
		statsJSON, _ := json.MarshalIndent(stats, "", "  ")
		if err := os.WriteFile("data/stats.json", statsJSON, 0644); err != nil {
			return fmt.Errorf("writing stats for %s: %w", user.username, err)
		}

		// Delete cached emblem to force fresh fetch
		os.Remove("data/emblem.jpg")

		// Fetch emblem
		fmt.Printf("Fetching emblem %s...\n", user.emblemHash)
		if err := bungie.FetchEmblem(user.emblemHash); err != nil {
			return fmt.Errorf("fetching emblem for %s: %w", user.username, err)
		}

		// Generate badge
		fmt.Println("Generating badge...")
		badgeStats := &badge.Stats{
			Username:     user.username,
			Commits:      user.commits,
			PullRequests: user.prs,
			Issues:       user.issues,
			Reviews:      user.reviews,
			Stars:        user.stars,
		}
		outputPath := fmt.Sprintf("examples/%s.png", user.username)
		if err := badge.Generate("data/emblem.jpg", badgeStats, outputPath); err != nil {
			return fmt.Errorf("generating badge for %s: %w", user.username, err)
		}

		fmt.Printf("✓ Badge saved to %s\n", outputPath)
	}

	fmt.Println("\n✨ All demo badges generated successfully!")
	return nil
}

func printUsage() {
	fmt.Fprintf(os.Stderr, "Usage: contribemblem <command>\n\n")
	fmt.Fprintf(os.Stderr, "Commands:\n")
	fmt.Fprintf(os.Stderr, "  fetch-stats      Fetch GitHub stats via GraphQL\n")
	fmt.Fprintf(os.Stderr, "  select-emblem    Select weekly emblem hash\n")
	fmt.Fprintf(os.Stderr, "  fetch-emblem     Fetch emblem image from Bungie API\n")
	fmt.Fprintf(os.Stderr, "  generate         Generate badge image\n")
	fmt.Fprintf(os.Stderr, "  update-readme    Update README with badge and timestamp\n")
	fmt.Fprintf(os.Stderr, "  run              Run full pipeline\n")
	fmt.Fprintf(os.Stderr, "  generate-demos   Generate example badges for demo users\n")
	fmt.Fprintf(os.Stderr, "  help             Show this help message\n")
}
