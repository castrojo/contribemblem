# Feature Research: GitHub Contribution Badge Generators

**Domain:** GitHub contribution visualization / profile readme badges
**Researched:** January 27, 2026
**Confidence:** HIGH

## Feature Landscape

### Table Stakes (Users Expect These)

Features users assume exist. Missing these = product feels incomplete.

| Feature | Why Expected | Complexity | API Source | Notes |
|---------|--------------|------------|------------|-------|
| Total commits | Core contribution metric shown everywhere | LOW | GraphQL User.contributionsCollection | Standard across all badge generators |
| Pull requests created | Primary collaboration metric | LOW | GraphQL User.pullRequests | github-readme-stats, profile-summary-cards show this |
| Issues opened | Community engagement signal | LOW | GraphQL User.issues | Standard metric |
| Public repositories | Shows project creation activity | LOW | REST API User.public_repos | Available in basic user endpoint |
| Stars received | Recognition metric | MEDIUM | GraphQL: must aggregate across repos | Requires repo iteration |
| Followers count | Social proof metric | LOW | REST API User.followers | Basic user field |
| Account age | Tenure signal | LOW | REST API User.created_at | Simple calculation from created date |
| Top languages | Shows technical diversity | MEDIUM | GitHub Linguist data via repos API | Requires aggregation across repos |

### Differentiators (Competitive Advantage)

Features that set the product apart. Not required, but valuable.

| Feature | Value Proposition | Complexity | Implementation Notes |
|---------|-------------------|------------|---------------------|
| Calendar year windowing | Shows current year activity only | LOW | Filter contributionsCollection by date range |
| Emblem artwork backgrounds | Gaming aesthetic (Destiny 2 themed) | MEDIUM | Custom PNG composition with imagemagick/canvas |
| Configurable metric selection | User picks 5 most impressive stats | LOW | YAML config, simple templating |
| Weekly emblem rotation | Fresh visual every week, from user's collection | LOW | Random selection seeded by week number |
| Emblem catalog browser | Discover available artworks | MEDIUM | Static asset gallery with thumbnails |
| Private repo stats | Include private contributions (self-hosted only) | MEDIUM | Requires authenticated API, user provides token |
| Commit streak tracking | Longest/current contribution streak | MEDIUM | Calculate from contributionsCollection |
| Review count | Shows code review activity | MEDIUM | GraphQL PullRequestReview count |
| Merged PR count | Success rate indicator | MEDIUM | Filter pullRequests by merged state |
| Org/repo filtering | Focus on specific projects (e.g., CNCF only) | MEDIUM | Filter contributions by org or repo list |

### Anti-Features (Commonly Requested, Often Problematic)

Features that seem good but create problems.

| Feature | Why Requested | Why Problematic | Alternative |
|---------|---------------|-----------------|-------------|
| Real-time updates | "I want live stats" | GitHub rate limits (5000 req/hour), API latency, unnecessary load | Cache for 24 hours, document refresh schedule |
| Historical data (all-time) | "Show my entire GitHub history" | Massive API calls, slow generation, data inflation (old activity less relevant) | Current calendar year only, keeps cards relevant |
| Social media integration | "Show Twitter followers too" | Scope creep, maintenance burden, API dependency hell | GitHub-only focus maintains simplicity |
| Skill endorsements | "Let others vote on my skills" | Becomes popularity contest, requires database/state, moderation burden | Top languages from actual code |
| Animated badges | "Make it move/flash" | GitHub markdown doesn't support video, accessibility issues, distraction | Static high-quality PNG with visual polish |
| Custom theming | "Let me change all colors" | Configuration complexity explosion, design inconsistency | Curated emblem catalog with coherent aesthetics |
| Leaderboards | "Rank me against others" | Encourages gaming metrics, not collaboration quality | Personal showcase only |

## Metric Inventory: GitHub API

### Fetchable via GraphQL (Preferred - More Efficient)

| Metric | GraphQL Field | Filtering Options | Complexity |
|--------|--------------|-------------------|------------|
| Total commits | `contributionsCollection.totalCommitContributions` | Date range (from/to) | EASY |
| Total PRs | `pullRequests.totalCount` | States, date filters | EASY |
| Total issues | `issues.totalCount` | States, date filters | EASY |
| PR reviews | `contributionsCollection.totalPullRequestReviewContributions` | Date range | EASY |
| Repositories contributed to | `contributionsCollection.totalRepositoriesWithContributedCommits` | Date range | EASY |
| Public repositories | `repositories.totalCount` | Privacy filter | EASY |
| Stars received | Sum of `repositories.stargazerCount` | Requires iteration | MEDIUM |
| Followers | `followers.totalCount` | N/A | EASY |
| Top languages | `repositories.languages` edges | Requires aggregation | MEDIUM |
| Commit streak | Manual calc from `contributionsCollection.contributionCalendar` | Date math | MEDIUM |
| Gists | `gists.totalCount` | Public/private filter | EASY |

### Fetchable via REST API (Simpler but Less Efficient)

| Metric | REST Endpoint | Notes | Complexity |
|--------|--------------|-------|------------|
| User profile data | `GET /users/{username}` | public_repos, followers, created_at, etc. | EASY |
| Events (activity stream) | `GET /users/{username}/events` | Last 90 days only, paginated | MEDIUM |
| Repository list | `GET /users/{username}/repos` | Paginated, can filter by type | MEDIUM |
| Starred repos | `GET /users/{username}/starred` | Count only, no metadata | EASY |

### NOT Easily Fetchable

| Metric | Why Hard | Workaround |
|--------|----------|-----------|
| Code review quality | No score/rating in API | Count only, not quality |
| Contribution to specific orgs | Must filter manually | Slow, requires many API calls |
| Lines of code | Not exposed by API | Language stats are proxy |
| Collaboration breadth | No direct metric | "Repos contributed to" is closest |
| Merged vs closed PRs ratio | Must calculate manually | Fetch all PRs and filter by state |

## Metric Display Conventions

### Number Formatting (from github-readme-stats patterns)

| Display Type | Example | When to Use |
|--------------|---------|-------------|
| Short format | `6.6k`, `1.2M` | Default for large numbers (saves space) |
| Long format | `6626`, `1,234,567` | When precision matters or <1000 |
| Percentages | `42%` | Ratios (merged PR rate, language distribution) |
| Rates | `12/week` | Velocity metrics (commits per week) |

### Typical Thresholds (Community Standards)

| Metric | Good | Great | Exceptional | Notes |
|--------|------|-------|-------------|-------|
| Annual commits | 500+ | 1000+ | 2000+ | Varies by role (maintainer vs contributor) |
| PRs opened | 50+ | 100+ | 200+ | Quality > quantity |
| Stars received | 100+ | 1000+ | 10k+ | Project-dependent |
| Repos contributed to | 10+ | 25+ | 50+ | Breadth signal |
| Followers | 50+ | 200+ | 1000+ | Community influence |

## Feature Dependencies

```
[Emblem artwork system]
    ├──requires──> [PNG generation (imagemagick/canvas)]
    ├──requires──> [Emblem asset library (static files)]
    └──enables──> [Weekly rotation feature]

[Configurable metric selection]
    ├──requires──> [YAML config parsing]
    ├──requires──> [Metric fetching abstraction]
    └──enables──> [User customization]

[Calendar year windowing]
    └──enhances──> [All time-based metrics]

[Private repo stats]
    ├──requires──> [Self-hosted deployment]
    ├──requires──> [User PAT storage]
    └──conflicts──> [Public Vercel hosting]

[GitHub Action automation]
    ├──requires──> [Scheduled workflow]
    ├──enables──> [Weekly emblem rotation]
    └──enables──> [Automatic README updates]
```

### Critical Path (Must Build in Order)

1. **Metric fetching layer** → GraphQL API client with caching
2. **PNG generation** → Template system with metric overlay
3. **Emblem catalog** → Static asset management
4. **Config parsing** → YAML schema for user preferences
5. **GitHub Action** → Automation orchestration
6. **Rotation logic** → Weekly emblem selection

## MVP Definition

### Launch With (v1.0)

Minimum viable product — what's needed to validate "Destiny-themed GitHub badge generator" concept.

- [x] **Fetch 5 core metrics** — Commits, PRs, Issues, Stars, Repos (GraphQL)
- [x] **Emblem PNG generation** — Overlay metrics on Destiny artwork background
- [x] **5 emblem templates** — Small curated set of Destiny 2 emblems
- [x] **Calendar year filtering** — Show current year stats only
- [x] **GitHub Action** — Weekly scheduled run, commit PNG to repo
- [x] **Basic YAML config** — Username, emblem IDs, metric selection
- [x] **README embedding** — Image tag with relative path

**Why these?** Proves core concept (gaming aesthetic + GitHub stats) with minimal complexity. Users can immediately see value.

### Add After Validation (v1.x)

Features to add once core is working and users are engaged.

- [ ] **Expanded emblem catalog** — 20+ Destiny emblems (HIGH priority)
- [ ] **Review count metric** — Shows code review activity (MEDIUM effort)
- [ ] **Commit streak** — Longest/current contribution streak (MEDIUM effort)
- [ ] **Merged PR count** — Success rate indicator (LOW effort)
- [ ] **Language breakdown** — Top 3 languages with percentages (MEDIUM effort)
- [ ] **Emblem preview tool** — Web UI to browse available emblems (MEDIUM effort)

**Triggers:** 
- 50+ users → Expand emblem catalog
- User requests → Add review/streak metrics
- Configuration friction → Build preview tool

### Future Consideration (v2.0+)

Features to defer until product-market fit is established.

- [ ] **Custom emblem upload** — User-provided artwork (requires moderation, storage)
- [ ] **Multi-user badges** — Team contribution cards (complex aggregation)
- [ ] **Org-specific filtering** — "Show only CNCF contributions" (API-intensive)
- [ ] **Historical year comparison** — "2025 vs 2024" view (storage required)
- [ ] **Achievement system** — Gamified milestones (scope creep risk)

**Why defer:** These add significant complexity without validating core value prop first.

## Feature Prioritization Matrix

| Feature | User Value | Implementation Cost | Priority | Phase |
|---------|------------|---------------------|----------|-------|
| Emblem PNG generation | HIGH | MEDIUM | P1 | MVP |
| 5 core metrics (commits, PRs, issues, stars, repos) | HIGH | LOW | P1 | MVP |
| Calendar year windowing | HIGH | LOW | P1 | MVP |
| GitHub Action automation | HIGH | LOW | P1 | MVP |
| Weekly emblem rotation | MEDIUM | LOW | P1 | MVP |
| Expanded emblem catalog (20+) | HIGH | MEDIUM | P2 | v1.x |
| Review count metric | MEDIUM | MEDIUM | P2 | v1.x |
| Commit streak tracking | MEDIUM | MEDIUM | P2 | v1.x |
| Merged PR count | MEDIUM | LOW | P2 | v1.x |
| Top 3 languages | MEDIUM | MEDIUM | P2 | v1.x |
| Emblem preview web UI | MEDIUM | HIGH | P2 | v1.x |
| Private repo stats | LOW | HIGH | P3 | v2.0+ |
| Custom emblem upload | LOW | HIGH | P3 | v2.0+ |
| Multi-user badges | LOW | HIGH | P3 | v2.0+ |
| Org filtering | LOW | HIGH | P3 | v2.0+ |

**Priority key:**
- P1: Must have for launch (MVP)
- P2: Should have, add when possible (v1.x)
- P3: Nice to have, future consideration (v2.0+)

## Competitor Feature Analysis

| Feature | github-readme-stats | profile-summary-cards | contribcard | ContribEmblem (Ours) |
|---------|---------------------|------------------------|-------------|----------------------|
| Total commits | ✅ (all-time) | ✅ (weekly graph) | ✅ (year filter) | ✅ (calendar year) |
| Pull requests | ✅ (count) | ✅ (with chart) | ✅ | ✅ (count + merged) |
| Issues | ✅ (count) | ✅ (with chart) | ✅ | ✅ (count) |
| Stars received | ✅ (total) | ❌ | ✅ | ✅ (total) |
| Top languages | ✅ (separate card) | ✅ (pie chart) | ❌ | ✅ (top 3-5) |
| Commit streak | ❌ | ❌ | ❌ | ✅ (v1.x) |
| Review count | ❌ | ❌ | ❌ | ✅ (v1.x) |
| **Visual theming** | Text themes (40+) | Predefined styles (15) | CNCF-branded | **Destiny emblem art** |
| **Time windowing** | All-time only | Weekly/monthly | Project-specific | **Calendar year** |
| **Customization** | URL params | GitHub Action config | Limited | **YAML config + emblem catalog** |
| **Update frequency** | On-demand (Vercel) | GitHub Action (daily) | Manual | **Weekly (Action)** |
| **Unique value** | Most popular, live | Detailed charts | CNCF focus | **Gaming aesthetic** |

### Key Insights

**What we're NOT competing on:**
- Comprehensiveness (github-readme-stats has 50+ metrics)
- Real-time updates (Vercel hosting enables this for others)
- Chart variety (profile-summary-cards has graphs/pie charts)

**What we ARE competing on:**
- **Aesthetic differentiation**: Destiny 2 emblem artwork (no one else does this)
- **Curation**: 5 carefully chosen metrics vs overwhelming choice
- **Relevance**: Calendar year focus (not all-time stats that get stale)
- **Rotation**: Weekly emblem change keeps profile fresh
- **Gaming community appeal**: Attracts a specific, passionate user base

## Metric Showcase Strategy

### Most Impressive Stats (High Perceived Value)

Users want to show off these metrics prominently:

1. **Total commits (calendar year)** — Shows sustained activity
2. **Stars received (total)** — Social proof of impact
3. **Pull requests merged** — Collaboration success rate
4. **Repositories contributed to** — Breadth signal
5. **Commit streak (current)** — Consistency indicator

### Supporting Metrics (Context)

These provide context but aren't primary showcase items:

- Issues opened (shows engagement, less impressive than PRs)
- Review count (good for maintainers, less relevant for contributors)
- Followers (social metric, not contribution metric)
- Top languages (interesting but not quantitative achievement)

### Default 5-Metric Template

**Recommended starter configuration:**

```yaml
metrics:
  - commits_year        # Top-left: Current year total
  - pull_requests       # Top-right: All-time PRs
  - stars               # Center: Most impressive number
  - repos_contributed   # Bottom-left: Breadth
  - streak_current      # Bottom-right: Consistency
```

**Rationale:** Balances recency (year commits), impact (stars), collaboration (PRs), breadth (repos), and consistency (streak).

## Sources

**HIGH confidence (official documentation):**
- GitHub GraphQL API Explorer: https://docs.github.com/en/graphql
- GitHub REST API Users endpoint: https://docs.github.com/en/rest/users/users
- GitHub contributionsCollection schema: https://docs.github.com/en/graphql/reference/objects#contributioncollection

**HIGH confidence (competitor analysis):**
- github-readme-stats (78k stars): https://github.com/anuraghazra/github-readme-stats
  - Shows: commits, PRs, issues, stars, languages, rank
  - Format: Short (6.6k) or long (6626)
  - Update: On-demand (Vercel serverless)
  - Themes: 40+ predefined color schemes
  
- github-profile-summary-cards (3.3k stars): https://github.com/vn7n24fzkq/github-profile-summary-cards
  - Shows: profile details graph, language pie charts, productive time, stats summary
  - Format: Multiple cards with charts/graphs
  - Update: GitHub Action (daily/weekly)
  - Unique: Time-based activity heatmap
  
- CNCF contribcard (28 stars): https://github.com/cncf/contribcard
  - Shows: Kubernetes project contributions
  - Format: Contributor recognition cards
  - Focus: CNCF/K8s ecosystem only
  - Unique: Project-specific branding (Kubernetes theme)

**MEDIUM confidence (API limitations observed):**
- GitHub rate limits: 5000 requests/hour (authenticated), forces caching strategy
- Event API limitation: Only 90 days of history, not suitable for all-time metrics
- Language stats: Requires iteration across all repos, computationally expensive

---
*Feature research for: ContribEmblem - GitHub contribution badge generator with Destiny 2 emblem aesthetics*  
*Researched: January 27, 2026*
