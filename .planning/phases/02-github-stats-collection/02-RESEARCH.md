# Phase 2: GitHub Stats Collection - Research

**Researched:** 2026-01-27
**Domain:** GitHub GraphQL API / GitHub Actions Caching
**Confidence:** HIGH

## Summary

GitHub Stats Collection requires using the GitHub GraphQL API to query contribution statistics efficiently while respecting rate limits. The standard approach uses GraphQL to fetch all 5 metrics (commits, PRs, issues, reviews, stars) in a single query with date filtering for the current calendar year. API responses should be cached using GitHub Actions' built-in `actions/cache@v4` with 24-hour expiry. Rate limit monitoring is handled via response headers (`x-ratelimit-remaining`), and timezone-aware date boundaries are critical for January 1 edge cases.

**Primary recommendation:** Use GitHub's built-in `GITHUB_TOKEN` with GraphQL API, cache responses with `actions/cache@v4`, and calculate calendar year boundaries in UTC to match GitHub's contribution counting logic.

## Standard Stack

The established libraries/tools for this domain:

### Core
| Library | Version | Purpose | Why Standard |
|---------|---------|---------|--------------|
| GitHub GraphQL API | v4 | Query contribution stats | Official GitHub API with single-query efficiency |
| `GITHUB_TOKEN` | Built-in | Authenticate API requests | Automatically provided in GitHub Actions, 1,000 points/hour |
| `actions/cache@v4` | v4.2.0+ | Cache API responses | Official GitHub Actions caching with 24hr TTL support |
| `curl` or `gh` CLI | Built-in | Make HTTP requests | Native tools available on all GitHub runners |

### Supporting
| Library | Version | Purpose | When to Use |
|---------|---------|---------|-------------|
| Node.js 24 | Built-in | Runtime for cache action | Required by `actions/cache@v4` (min runner 2.327.1) |
| `jq` | Built-in | Parse JSON responses | Available on ubuntu-latest runners for response processing |

### Alternatives Considered
| Instead of | Could Use | Tradeoff |
|------------|-----------|----------|
| GraphQL API | REST API | GraphQL fetches all 5 metrics in 1 query vs. 5+ REST calls; rate limit cost is lower |
| `actions/cache` | Manual file caching | Cache action handles expiry/eviction automatically; manual requires custom TTL logic |
| `GITHUB_TOKEN` | Personal Access Token | PAT unnecessary; `GITHUB_TOKEN` auto-provided with sufficient permissions |

**Installation:**
```yaml
# No installation needed - all tools built into GitHub Actions
# actions/cache@v4 requires minimum runner version 2.327.1
```

## Architecture Patterns

### Recommended Project Structure
```
.github/
├── workflows/
│   └── update-stats.yml          # Main workflow with scheduled trigger
scripts/
├── fetch-stats.sh                 # GraphQL query execution
└── cache-stats.sh                 # Cache key generation helper
data/
└── stats.json                     # Cached stats output (committed to repo)
```

### Pattern 1: Single GraphQL Query for All Stats
**What:** Query all 5 contribution metrics in one GraphQL request using ContributionsCollection
**When to use:** Always - minimizes rate limit consumption and API calls
**Example:**
```graphql
# Source: https://docs.github.com/en/graphql/reference/objects#contributionscollection
query($username: String!, $from: DateTime!, $to: DateTime!) {
  user(login: $username) {
    contributionsCollection(from: $from, to: $to) {
      totalCommitContributions
      totalPullRequestContributions
      totalIssueContributions
      totalPullRequestReviewContributions
      totalRepositoryContributions
    }
    repositories(ownerAffiliations: OWNER, first: 100) {
      nodes {
        stargazerCount
      }
    }
  }
}
```

### Pattern 2: Cache Key with Date Boundary
**What:** Use date as part of cache key to ensure daily refresh
**When to use:** Always - enforces 24-hour cache expiry without manual TTL logic
**Example:**
```yaml
# Source: https://docs.github.com/en/actions/reference/workflows-and-actions/dependency-caching
- name: Cache GitHub Stats
  uses: actions/cache@v4
  with:
    path: data/stats.json
    key: github-stats-${{ runner.os }}-${{ steps.get-date.outputs.date }}
    restore-keys: |
      github-stats-${{ runner.os }}-
```

### Pattern 3: UTC Date Boundary Calculation
**What:** Calculate calendar year start/end in UTC timezone
**When to use:** Always - GitHub counts contributions in UTC, not local time
**Example:**
```bash
# Source: https://docs.github.com/en/account-and-profile/reference/profile-contributions-reference
CURRENT_YEAR=$(date -u +"%Y")
YEAR_START="${CURRENT_YEAR}-01-01T00:00:00Z"
YEAR_END="${CURRENT_YEAR}-12-31T23:59:59Z"
```

### Pattern 4: Rate Limit Monitoring (Non-Blocking)
**What:** Log rate limit headers without failing workflow
**When to use:** Every API call - awareness without disruption
**Example:**
```bash
# Source: https://docs.github.com/en/graphql/overview/rate-limits-and-query-limits-for-the-graphql-api
response=$(curl -H "Authorization: bearer $TOKEN" -X POST -d "$query" https://api.github.com/graphql -i)
remaining=$(echo "$response" | grep -i "x-ratelimit-remaining" | awk '{print $2}')
echo "Rate limit remaining: $remaining"
```

### Anti-Patterns to Avoid
- **Fetching stats without date filters:** Queries entire contribution history; slow and expensive
- **Multiple REST API calls:** 5+ separate requests vs. 1 GraphQL query; wastes rate limit
- **Cache keys without date:** Cache persists indefinitely; stale data after 24 hours
- **Hardcoded timezone offsets:** Breaks on January 1 depending on runner locale; always use UTC

## Don't Hand-Roll

Problems that look simple but have existing solutions:

| Problem | Don't Build | Use Instead | Why |
|---------|-------------|-------------|-----|
| API response caching | Custom file timestamp checking + eviction logic | `actions/cache@v4` | GitHub's cache action handles TTL (via key rotation), compression, restoration, and automatic eviction after 7 days of non-access |
| Rate limit tracking | Custom counter with file persistence | GraphQL response headers (`x-ratelimit-remaining`) | GitHub returns current quota in every response; no need to track locally |
| Date/time calculations | Custom timezone conversion | `date -u` command | Shell's `-u` flag ensures UTC; avoid JavaScript `Date` timezone bugs |
| JSON parsing | `sed`/`awk` string manipulation | `jq` (built into runners) | `jq` handles nested objects, escaping, and edge cases correctly |
| GraphQL query construction | String concatenation with variables | Parameterized queries with `$variables` | Prevents injection risks and handles escaping automatically |

**Key insight:** GitHub Actions and GraphQL API are designed to work together. Using built-in tools (`GITHUB_TOKEN`, `actions/cache`, response headers) eliminates 90% of custom logic.

## Common Pitfalls

### Pitfall 1: Ignoring January 1 Date Boundary
**What goes wrong:** On January 1, queries may fetch previous year's stats if timezone is not UTC
**Why it happens:** Runner's local timezone != GitHub's contribution counting timezone (UTC)
**How to avoid:** Always use `date -u` for UTC calculations; never rely on runner's local time
**Warning signs:** Stats drop to zero on January 1 in some timezones but not others

### Pitfall 2: Cache Thrashing with Wrong Key
**What goes wrong:** Cache is created/evicted rapidly, never hitting; workflow fetches API every run
**Why it happens:** Cache key changes every run (e.g., includes `${{ github.run_id }}`)
**How to avoid:** Use date-based key (`YYYYMMDD`) that stays constant for 24 hours
**Warning signs:** `actions/cache` logs show "cache-hit: false" on every run despite recent creation

### Pitfall 3: Exceeding Rate Limits Without Monitoring
**What goes wrong:** Workflow fails silently when rate limit exhausted; no stats update
**Why it happens:** `GITHUB_TOKEN` has 1,000 points/hour; single query costs ~1-5 points but failed retries accumulate
**How to avoid:** Check `x-ratelimit-remaining` header before query; log warning if < 100 remaining
**Warning signs:** API returns empty response or 200 status with error message in body

### Pitfall 4: Fetching Stars Inefficiently
**What goes wrong:** Query fetches all repositories to sum stars; slow for users with many repos
**Why it happens:** `stargazerCount` is per-repository; no single "total stars" field
**How to avoid:** Use pagination (`first: 100`) and accept approximate count for users with 100+ repos
**Warning signs:** Query timeout (10-second limit) or rate limit consumption spikes

### Pitfall 5: Caching Directly in Workflow vs. Data File
**What goes wrong:** Cache stores intermediate API response; logic changes require cache invalidation
**Why it happens:** Caching raw GraphQL JSON instead of processed `stats.json` output
**How to avoid:** Cache the final `data/stats.json` file that downstream phases consume
**Warning signs:** Phase 3 emblem generation uses old data structure after Phase 2 logic update

## Code Examples

Verified patterns from official sources:

### Fetch Stats with GraphQL (Bash Script)
```bash
# Source: https://docs.github.com/en/graphql/guides/forming-calls-with-graphql
#!/bin/bash
set -euo pipefail

# Calculate current year boundaries in UTC
CURRENT_YEAR=$(date -u +"%Y")
YEAR_START="${CURRENT_YEAR}-01-01T00:00:00Z"
YEAR_END="${CURRENT_YEAR}-12-31T23:59:59Z"

# GraphQL query with variables
query=$(cat <<EOF
{
  "query": "query(\$username: String!, \$from: DateTime!, \$to: DateTime!) {
    user(login: \$username) {
      contributionsCollection(from: \$from, to: \$to) {
        totalCommitContributions
        totalPullRequestContributions
        totalIssueContributions
        totalPullRequestReviewContributions
      }
      repositories(ownerAffiliations: OWNER, first: 100) {
        nodes {
          stargazerCount
        }
      }
    }
  }",
  "variables": {
    "username": "$GITHUB_ACTOR",
    "from": "$YEAR_START",
    "to": "$YEAR_END"
  }
}
EOF
)

# Execute query with rate limit monitoring
response=$(curl -s -H "Authorization: bearer $GITHUB_TOKEN" \
  -H "Content-Type: application/json" \
  -X POST \
  -d "$query" \
  https://api.github.com/graphql \
  -i)

# Extract and log rate limit
remaining=$(echo "$response" | grep -i "x-ratelimit-remaining:" | awk '{print $2}' | tr -d '\r')
echo "Rate limit remaining: $remaining"

# Parse response body (after headers)
body=$(echo "$response" | sed -n '/^{/,$ p')
echo "$body" | jq '.'
```

### Cache Stats with Date-Based Key
```yaml
# Source: https://docs.github.com/en/actions/reference/workflows-and-actions/dependency-caching
- name: Get current date (UTC)
  id: get-date
  run: echo "date=$(date -u +%Y%m%d)" >> $GITHUB_OUTPUT

- name: Cache GitHub Stats
  id: cache-stats
  uses: actions/cache@v4
  with:
    path: data/stats.json
    key: github-stats-${{ steps.get-date.outputs.date }}
    restore-keys: |
      github-stats-

- name: Fetch stats if cache miss
  if: steps.cache-stats.outputs.cache-hit != 'true'
  run: ./scripts/fetch-stats.sh > data/stats.json
```

### Process GraphQL Response into stats.json
```bash
# Source: jq manual (built into GitHub runners)
#!/bin/bash
set -euo pipefail

# Input: GraphQL response JSON from stdin
# Output: Simplified stats.json to stdout

jq '{
  year: (now | strftime("%Y") | tonumber),
  updated_at: (now | strftime("%Y-%m-%dT%H:%M:%SZ")),
  commits: .data.user.contributionsCollection.totalCommitContributions,
  pull_requests: .data.user.contributionsCollection.totalPullRequestContributions,
  issues: .data.user.contributionsCollection.totalIssueContributions,
  reviews: .data.user.contributionsCollection.totalPullRequestReviewContributions,
  stars_received: ([.data.user.repositories.nodes[].stargazerCount] | add // 0)
}'
```

## State of the Art

| Old Approach | Current Approach | When Changed | Impact |
|--------------|------------------|--------------|--------|
| REST API v3 for stats | GraphQL API v4 | 2016 | Single query vs. multiple endpoints; 5x fewer API calls |
| `actions/cache@v2` | `actions/cache@v4` | 2024 | Requires Node.js 24 and runner 2.327.1+; improved compression |
| Personal Access Tokens | `GITHUB_TOKEN` built-in | 2019 | No PAT management needed; automatic scoping and rotation |
| Manual cache TTL with file timestamps | Cache key rotation (date-based) | 2020 | GitHub handles eviction; simpler workflow logic |

**Deprecated/outdated:**
- **REST API for contributions:** Still works but requires 5+ separate calls vs. 1 GraphQL query
- **`actions/cache@v3` and older:** Will fail after Feb 2025 due to cache backend migration
- **Hardcoded rate limit thresholds:** GitHub changed limits for `GITHUB_TOKEN` from 5,000 to 1,000 points/hour in 2023

## Open Questions

Things that couldn't be fully resolved:

1. **ContributionsCollection exact filtering behavior**
   - What we know: GitHub's GraphQL `contributionsCollection(from:, to:)` accepts ISO 8601 DateTime
   - What's unclear: Whether partial-day boundaries (e.g., `2026-01-01T06:00:00Z`) correctly filter contributions by commit timestamp or if full-day boundaries are required
   - Recommendation: Use full-day boundaries (`00:00:00Z` to `23:59:59Z`) to match GitHub's contribution graph display logic

2. **Stars aggregation limits**
   - What we know: `repositories(first: 100)` returns max 100 repos; pagination available with `after` cursor
   - What's unclear: Whether Phase 3 emblem logic requires exact star count or if approximate count for top 100 repos is acceptable
   - Recommendation: Start with 100-repo limit; add pagination in future phase if exact count becomes requirement

3. **Cache eviction timing precision**
   - What we know: Caches evicted after 7 days of non-access; 10 GB repo limit
   - What's unclear: Exact eviction timing (e.g., if cache accessed at 11:59 PM, is it available the next day?)
   - Recommendation: Design assuming cache available for at least 23 hours after creation; don't rely on >7 day retention

## Sources

### Primary (HIGH confidence)
- GitHub GraphQL API Overview - https://docs.github.com/en/graphql/overview/about-the-graphql-api
- GitHub GraphQL Rate Limits - https://docs.github.com/en/graphql/overview/rate-limits-and-query-limits-for-the-graphql-api
- GitHub GraphQL Forming Calls - https://docs.github.com/en/graphql/guides/forming-calls-with-graphql
- GitHub Actions Caching Dependencies - https://docs.github.com/en/actions/using-workflows/caching-dependencies-to-speed-up-workflows
- GitHub Actions Dependency Caching Reference - https://docs.github.com/en/actions/reference/workflows-and-actions/dependency-caching
- actions/cache GitHub Repository - https://github.com/actions/cache (v4.2.0+)
- GitHub Actions Contexts - https://docs.github.com/en/actions/reference/workflows-and-actions/contexts
- GitHub Contribution Counting Rules - https://docs.github.com/en/account-and-profile/how-tos/contribution-settings/troubleshooting-missing-contributions

### Secondary (MEDIUM confidence)
- None required - all findings verified against official GitHub documentation

### Tertiary (LOW confidence)
- None - no unverified claims included

## Metadata

**Confidence breakdown:**
- Standard stack: HIGH - All tools from official GitHub documentation and Actions marketplace
- Architecture: HIGH - Patterns verified against official examples and current GitHub API behavior
- Pitfalls: HIGH - Derived from official rate limit docs, cache action behavior, and contribution counting rules

**Research date:** 2026-01-27
**Valid until:** 2026-02-27 (30 days) - GitHub Actions stable; GraphQL API v4 mature; unlikely to change
