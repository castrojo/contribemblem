#!/bin/bash
set -euo pipefail

# Fetch GitHub stats via GraphQL API with UTC date boundaries and rate limit monitoring
# Output: GraphQL response JSON (piped to process-stats.sh)

# Calculate current year boundaries in UTC (matches GitHub's contribution counting logic)
CURRENT_YEAR=$(date -u +"%Y")
YEAR_START="${CURRENT_YEAR}-01-01T00:00:00Z"
YEAR_END="${CURRENT_YEAR}-12-31T23:59:59Z"

# Build GraphQL query with variables
query=$(cat <<EOF
{
  "query": "query(\$username: String!, \$from: DateTime!, \$to: DateTime!) { user(login: \$username) { contributionsCollection(from: \$from, to: \$to) { totalCommitContributions totalPullRequestContributions totalIssueContributions totalPullRequestReviewContributions } repositories(ownerAffiliations: OWNER, first: 100) { nodes { stargazerCount } } } }",
  "variables": {
    "username": "$GITHUB_ACTOR",
    "from": "$YEAR_START",
    "to": "$YEAR_END"
  }
}
EOF
)

# Execute GraphQL query with rate limit monitoring
response=$(curl -s -H "Authorization: bearer $GITHUB_TOKEN" \
  -H "Content-Type: application/json" \
  -X POST \
  -d "$query" \
  https://api.github.com/graphql \
  -i)

# Extract and log rate limit remaining (case-insensitive, strip carriage returns)
remaining=$(echo "$response" | grep -i "x-ratelimit-remaining:" | awk '{print $2}' | tr -d '\r')
if [ -n "$remaining" ]; then
  echo "Rate limit remaining: $remaining" >&2
  if [ "$remaining" -lt 100 ]; then
    echo "Warning: Rate limit running low (< 100 remaining)" >&2
  fi
fi

# Extract response body (JSON after blank line separating headers)
body=$(echo "$response" | sed -n '/^{/,$ p')

# Pipe to process-stats.sh for transformation
echo "$body" | $(dirname "$0")/process-stats.sh
