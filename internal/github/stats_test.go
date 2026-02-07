package github

import (
	"encoding/json"
	"testing"
)

func TestStatsJSONMarshaling(t *testing.T) {
	stats := &Stats{
		Year:          2026,
		UpdatedAt:     "2026-02-06T22:00:00Z",
		Commits:       150,
		PullRequests:  25,
		Issues:        10,
		Reviews:       30,
		StarsReceived: 150,
	}

	data, err := json.Marshal(stats)
	if err != nil {
		t.Fatalf("Failed to marshal stats: %v", err)
	}

	var unmarshaled Stats
	if err := json.Unmarshal(data, &unmarshaled); err != nil {
		t.Fatalf("Failed to unmarshal stats: %v", err)
	}

	if unmarshaled.Commits != stats.Commits {
		t.Errorf("Expected Commits=%d, got %d", stats.Commits, unmarshaled.Commits)
	}
	if unmarshaled.PullRequests != stats.PullRequests {
		t.Errorf("Expected PullRequests=%d, got %d", stats.PullRequests, unmarshaled.PullRequests)
	}
	if unmarshaled.Issues != stats.Issues {
		t.Errorf("Expected Issues=%d, got %d", stats.Issues, unmarshaled.Issues)
	}
	if unmarshaled.Reviews != stats.Reviews {
		t.Errorf("Expected Reviews=%d, got %d", stats.Reviews, unmarshaled.Reviews)
	}
	if unmarshaled.StarsReceived != stats.StarsReceived {
		t.Errorf("Expected StarsReceived=%d, got %d", stats.StarsReceived, unmarshaled.StarsReceived)
	}
}
