package github

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"os"
	"time"
)

// Stats represents GitHub contribution statistics
type Stats struct {
	Year          int    `json:"year"`
	UpdatedAt     string `json:"updated_at"`
	Commits       int    `json:"commits"`
	PullRequests  int    `json:"pull_requests"`
	Issues        int    `json:"issues"`
	Reviews       int    `json:"reviews"`
	StarsReceived int    `json:"stars_received"`
}

// GraphQL response structures (nested for JSON unmarshaling)
type graphQLResponse struct {
	Data struct {
		User struct {
			ContributionsCollection struct {
				TotalCommitContributions            int `json:"totalCommitContributions"`
				TotalPullRequestContributions       int `json:"totalPullRequestContributions"`
				TotalIssueContributions             int `json:"totalIssueContributions"`
				TotalPullRequestReviewContributions int `json:"totalPullRequestReviewContributions"`
			} `json:"contributionsCollection"`
			Repositories struct {
				Nodes []struct {
					StargazerCount int `json:"stargazerCount"`
				} `json:"nodes"`
			} `json:"repositories"`
		} `json:"user"`
	} `json:"data"`
}

// FetchStats queries GitHub GraphQL API for user contribution stats
// Requires GITHUB_TOKEN env var
// If username is empty, falls back to GITHUB_ACTOR env var
// If client is nil, uses default http.Client with 30s timeout
func FetchStats(username string, client *http.Client) (*Stats, error) {
	token := os.Getenv("GITHUB_TOKEN")
	if token == "" {
		return nil, fmt.Errorf("GITHUB_TOKEN environment variable not set")
	}

	// Fall back to GITHUB_ACTOR env var if username not provided
	if username == "" {
		username = os.Getenv("GITHUB_ACTOR")
		if username == "" {
			return nil, fmt.Errorf("username not provided and GITHUB_ACTOR environment variable not set")
		}
	}

	// Calculate current year boundaries in UTC (matches GitHub's contribution logic)
	now := time.Now().UTC()
	currentYear := now.Year()
	yearStart := fmt.Sprintf("%d-01-01T00:00:00Z", currentYear)
	yearEnd := fmt.Sprintf("%d-12-31T23:59:59Z", currentYear)

	// Build GraphQL query
	query := map[string]interface{}{
		"query": "query($username: String!, $from: DateTime!, $to: DateTime!) { user(login: $username) { contributionsCollection(from: $from, to: $to) { totalCommitContributions totalPullRequestContributions totalIssueContributions totalPullRequestReviewContributions } repositories(ownerAffiliations: OWNER, first: 100) { nodes { stargazerCount } } } }",
		"variables": map[string]string{
			"username": username,
			"from":     yearStart,
			"to":       yearEnd,
		},
	}

	body, err := json.Marshal(query)
	if err != nil {
		return nil, fmt.Errorf("failed to marshal query: %w", err)
	}

	// Execute GraphQL request
	req, err := http.NewRequest("POST", "https://api.github.com/graphql", bytes.NewReader(body))
	if err != nil {
		return nil, fmt.Errorf("failed to create request: %w", err)
	}

	req.Header.Set("Authorization", "bearer "+token)
	req.Header.Set("Content-Type", "application/json")

	// Use provided client or create default
	if client == nil {
		client = &http.Client{Timeout: 30 * time.Second}
	}
	resp, err := client.Do(req)
	if err != nil {
		return nil, fmt.Errorf("GraphQL request failed: %w", err)
	}
	defer resp.Body.Close()

	// Log rate limit (non-blocking)
	if remaining := resp.Header.Get("X-Ratelimit-Remaining"); remaining != "" {
		fmt.Fprintf(os.Stderr, "Rate limit remaining: %s\n", remaining)
	}

	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("GraphQL request returned status %d", resp.StatusCode)
	}

	// Parse response
	respBody, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, fmt.Errorf("failed to read response: %w", err)
	}

	var gqlResp graphQLResponse
	if err := json.Unmarshal(respBody, &gqlResp); err != nil {
		return nil, fmt.Errorf("failed to parse response: %w", err)
	}

	// Sum stars from all repositories
	totalStars := 0
	for _, repo := range gqlResp.Data.User.Repositories.Nodes {
		totalStars += repo.StargazerCount
	}

	// Transform to Stats struct (equivalent to process-stats.sh)
	stats := &Stats{
		Year:          currentYear,
		UpdatedAt:     now.Format("2006-01-02T15:04:05Z"),
		Commits:       gqlResp.Data.User.ContributionsCollection.TotalCommitContributions,
		PullRequests:  gqlResp.Data.User.ContributionsCollection.TotalPullRequestContributions,
		Issues:        gqlResp.Data.User.ContributionsCollection.TotalIssueContributions,
		Reviews:       gqlResp.Data.User.ContributionsCollection.TotalPullRequestReviewContributions,
		StarsReceived: totalStars,
	}

	return stats, nil
}
