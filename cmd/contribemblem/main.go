package main

import (
	"encoding/json"
	"fmt"
	"os"

	"github.com/castrojo/contribemblem/internal/badge"
	"github.com/castrojo/contribemblem/internal/bungie"
	"github.com/castrojo/contribemblem/internal/emblem"
	"github.com/castrojo/contribemblem/internal/github"
)

func main() {
	if len(os.Args) < 2 {
		printUsage()
		os.Exit(1)
	}

	cmd := os.Args[1]
	switch cmd {
	case "fetch-stats":
		stats, err := github.FetchStats()
		if err != nil {
			fmt.Fprintf(os.Stderr, "Error: %v\n", err)
			os.Exit(1)
		}
		data, _ := json.MarshalIndent(stats, "", "  ")
		fmt.Println(string(data))
	case "select-emblem":
		selectedEmblem, err := emblem.SelectEmblem(emblem.DefaultConfigPath)
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
	case "run":
		fmt.Println("Running full ContribEmblem pipeline...")

		// Step 1: Fetch GitHub stats
		fmt.Println("[1/4] Fetching GitHub stats...")
		stats, err := github.FetchStats()
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
		fmt.Println("[2/4] Selecting weekly emblem...")
		emblemHash, err := emblem.SelectEmblem(emblem.DefaultConfigPath)
		if err != nil {
			fmt.Fprintf(os.Stderr, "Failed to select emblem: %v\n", err)
			os.Exit(1)
		}
		fmt.Printf("✓ Selected emblem: %s\n", emblemHash)

		// Step 3: Fetch emblem from Bungie
		fmt.Println("[3/4] Fetching emblem from Bungie API...")
		if err := bungie.FetchEmblem(emblemHash); err != nil {
			fmt.Fprintf(os.Stderr, "Failed to fetch emblem: %v\n", err)
			os.Exit(1)
		}
		fmt.Println("✓ Emblem downloaded to data/emblem.jpg")

		// Step 4: Generate badge
		fmt.Println("[4/4] Generating badge image...")
		badgeStats := &badge.Stats{
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

		fmt.Println("\n🎉 Pipeline complete! Badge ready at badge.png")
	case "help", "--help", "-h":
		printUsage()
		os.Exit(0)
	default:
		fmt.Fprintf(os.Stderr, "Unknown command: %s\n", cmd)
		printUsage()
		os.Exit(1)
	}
}

func printUsage() {
	fmt.Fprintf(os.Stderr, "Usage: contribemblem <command>\n\n")
	fmt.Fprintf(os.Stderr, "Commands:\n")
	fmt.Fprintf(os.Stderr, "  fetch-stats      Fetch GitHub stats via GraphQL\n")
	fmt.Fprintf(os.Stderr, "  select-emblem    Select weekly emblem hash\n")
	fmt.Fprintf(os.Stderr, "  fetch-emblem     Fetch emblem image from Bungie API\n")
	fmt.Fprintf(os.Stderr, "  generate         Generate badge image\n")
	fmt.Fprintf(os.Stderr, "  run              Run full pipeline\n")
	fmt.Fprintf(os.Stderr, "  help             Show this help message\n")
}
