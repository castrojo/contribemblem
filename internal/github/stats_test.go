package github

import (
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"os"
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

// TestFetchStats_Success tests successful GraphQL response parsing
func TestFetchStats_Success(t *testing.T) {
	// Mock successful GitHub GraphQL response
	mockResponse := `{
		"data": {
			"user": {
				"contributionsCollection": {
					"totalCommitContributions": 42,
					"totalPullRequestContributions": 15,
					"totalIssueContributions": 8,
					"totalPullRequestReviewContributions": 23
				},
				"repositories": {
					"nodes": [
						{"stargazerCount": 100},
						{"stargazerCount": 50},
						{"stargazerCount": 25}
					]
				}
			}
		}
	}`

	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		// Verify request headers
		if auth := r.Header.Get("Authorization"); auth != "bearer test-token" {
			t.Errorf("Expected Authorization header 'bearer test-token', got '%s'", auth)
		}
		if ct := r.Header.Get("Content-Type"); ct != "application/json" {
			t.Errorf("Expected Content-Type 'application/json', got '%s'", ct)
		}

		// Set rate limit header
		w.Header().Set("X-Ratelimit-Remaining", "4999")
		w.WriteHeader(http.StatusOK)
		w.Write([]byte(mockResponse))
	}))
	defer server.Close()

	// Set required env vars
	os.Setenv("GITHUB_TOKEN", "test-token")
	os.Setenv("GITHUB_ACTOR", "testuser")
	defer func() {
		os.Unsetenv("GITHUB_TOKEN")
		os.Unsetenv("GITHUB_ACTOR")
	}()

	// Create client that uses the test server
	client := &http.Client{
		Transport: &mockTransport{
			handler: func(req *http.Request) (*http.Response, error) {
				// Redirect to test server
				req.URL.Scheme = "http"
				req.URL.Host = server.URL[7:] // Remove "http://"
				return http.DefaultTransport.RoundTrip(req)
			},
		},
	}

	stats, err := FetchStats("testuser", client)
	if err != nil {
		t.Fatalf("FetchStats() failed: %v", err)
	}

	// Verify parsed stats
	if stats.Commits != 42 {
		t.Errorf("Expected Commits=42, got %d", stats.Commits)
	}
	if stats.PullRequests != 15 {
		t.Errorf("Expected PullRequests=15, got %d", stats.PullRequests)
	}
	if stats.Issues != 8 {
		t.Errorf("Expected Issues=8, got %d", stats.Issues)
	}
	if stats.Reviews != 23 {
		t.Errorf("Expected Reviews=23, got %d", stats.Reviews)
	}
	if stats.StarsReceived != 175 { // 100 + 50 + 25
		t.Errorf("Expected StarsReceived=175, got %d", stats.StarsReceived)
	}
}

// TestFetchStats_RateLimitHeader tests rate limit header logging
func TestFetchStats_RateLimitHeader(t *testing.T) {
	mockResponse := `{
		"data": {
			"user": {
				"contributionsCollection": {
					"totalCommitContributions": 10,
					"totalPullRequestContributions": 5,
					"totalIssueContributions": 2,
					"totalPullRequestReviewContributions": 3
				},
				"repositories": {
					"nodes": []
				}
			}
		}
	}`

	rateLimitCalled := false
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("X-Ratelimit-Remaining", "100")
		rateLimitCalled = true
		w.WriteHeader(http.StatusOK)
		w.Write([]byte(mockResponse))
	}))
	defer server.Close()

	os.Setenv("GITHUB_TOKEN", "test-token")
	os.Setenv("GITHUB_ACTOR", "testuser")
	defer func() {
		os.Unsetenv("GITHUB_TOKEN")
		os.Unsetenv("GITHUB_ACTOR")
	}()

	client := &http.Client{
		Transport: &mockTransport{
			handler: func(req *http.Request) (*http.Response, error) {
				req.URL.Scheme = "http"
				req.URL.Host = server.URL[7:]
				return http.DefaultTransport.RoundTrip(req)
			},
		},
	}

	_, err := FetchStats("testuser", client)
	if err != nil {
		t.Fatalf("FetchStats() failed: %v", err)
	}

	if !rateLimitCalled {
		t.Error("Rate limit header was not checked")
	}
}

// TestFetchStats_Non200StatusCode tests handling of non-200 HTTP status codes
func TestFetchStats_Non200StatusCode(t *testing.T) {
	testCases := []struct {
		name       string
		statusCode int
	}{
		{"Unauthorized", http.StatusUnauthorized},
		{"Forbidden", http.StatusForbidden},
		{"Internal Server Error", http.StatusInternalServerError},
		{"Bad Gateway", http.StatusBadGateway},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
				w.WriteHeader(tc.statusCode)
				w.Write([]byte(`{"error": "test error"}`))
			}))
			defer server.Close()

			os.Setenv("GITHUB_TOKEN", "test-token")
			os.Setenv("GITHUB_ACTOR", "testuser")
			defer func() {
				os.Unsetenv("GITHUB_TOKEN")
				os.Unsetenv("GITHUB_ACTOR")
			}()

			client := &http.Client{
				Transport: &mockTransport{
					handler: func(req *http.Request) (*http.Response, error) {
						req.URL.Scheme = "http"
						req.URL.Host = server.URL[7:]
						return http.DefaultTransport.RoundTrip(req)
					},
				},
			}

			_, err := FetchStats("testuser", client)
			if err == nil {
				t.Fatalf("Expected error for status code %d, got nil", tc.statusCode)
			}

			expectedErr := "GraphQL request returned status"
			if err.Error()[:len(expectedErr)] != expectedErr {
				t.Errorf("Expected error containing '%s', got '%v'", expectedErr, err)
			}
		})
	}
}

// TestFetchStats_MalformedJSON tests handling of malformed JSON responses
func TestFetchStats_MalformedJSON(t *testing.T) {
	testCases := []struct {
		name     string
		response string
	}{
		{"Invalid JSON", `{"data": invalid json}`},
		{"Truncated JSON", `{"data": {"user": {`},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
				w.WriteHeader(http.StatusOK)
				w.Write([]byte(tc.response))
			}))
			defer server.Close()

			os.Setenv("GITHUB_TOKEN", "test-token")
			os.Setenv("GITHUB_ACTOR", "testuser")
			defer func() {
				os.Unsetenv("GITHUB_TOKEN")
				os.Unsetenv("GITHUB_ACTOR")
			}()

			client := &http.Client{
				Transport: &mockTransport{
					handler: func(req *http.Request) (*http.Response, error) {
						req.URL.Scheme = "http"
						req.URL.Host = server.URL[7:]
						return http.DefaultTransport.RoundTrip(req)
					},
				},
			}

			_, err := FetchStats("", client)
			if err == nil {
				t.Fatal("Expected error for malformed JSON, got nil")
			}

			expectedErr := "failed to parse response"
			if err.Error()[:len(expectedErr)] != expectedErr {
				t.Errorf("Expected error containing '%s', got '%v'", expectedErr, err)
			}
		})
	}
}

// TestFetchStats_MissingEnvVars tests handling of missing environment variables
func TestFetchStats_MissingEnvVars(t *testing.T) {
	testCases := []struct {
		name        string
		token       string
		actor       string
		expectedErr string
	}{
		{
			name:        "Missing GITHUB_TOKEN",
			token:       "",
			actor:       "testuser",
			expectedErr: "GITHUB_TOKEN environment variable not set",
		},
		{
			name:        "Missing GITHUB_ACTOR and no username provided",
			token:       "test-token",
			actor:       "",
			expectedErr: "username not provided and GITHUB_ACTOR environment variable not set",
		},
		{
			name:        "Missing both",
			token:       "",
			actor:       "",
			expectedErr: "GITHUB_TOKEN environment variable not set",
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			// Clear env vars first
			os.Unsetenv("GITHUB_TOKEN")
			os.Unsetenv("GITHUB_ACTOR")

			// Set only the non-empty ones
			if tc.token != "" {
				os.Setenv("GITHUB_TOKEN", tc.token)
				defer os.Unsetenv("GITHUB_TOKEN")
			}
			if tc.actor != "" {
				os.Setenv("GITHUB_ACTOR", tc.actor)
				defer os.Unsetenv("GITHUB_ACTOR")
			}

			_, err := FetchStats("", nil)
			if err == nil {
				t.Fatal("Expected error for missing env var, got nil")
			}

			if err.Error() != tc.expectedErr {
				t.Errorf("Expected error '%s', got '%v'", tc.expectedErr, err)
			}
		})
	}
}

// mockTransport implements http.RoundTripper for testing
type mockTransport struct {
	handler func(*http.Request) (*http.Response, error)
}

func (m *mockTransport) RoundTrip(req *http.Request) (*http.Response, error) {
	return m.handler(req)
}
