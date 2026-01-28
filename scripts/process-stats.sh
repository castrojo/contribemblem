#!/bin/bash
set -euo pipefail

# Process GraphQL response into standardized stats.json format
# Input: GraphQL response JSON from stdin
# Output: Simplified stats.json with 5 metrics + metadata

jq '{
  year: (now | strftime("%Y") | tonumber),
  updated_at: (now | strftime("%Y-%m-%dT%H:%M:%SZ")),
  commits: .data.user.contributionsCollection.totalCommitContributions,
  pull_requests: .data.user.contributionsCollection.totalPullRequestContributions,
  issues: .data.user.contributionsCollection.totalIssueContributions,
  reviews: .data.user.contributionsCollection.totalPullRequestReviewContributions,
  stars_received: ([.data.user.repositories.nodes[].stargazerCount] | add // 0)
}'
