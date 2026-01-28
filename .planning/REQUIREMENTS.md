# Requirements: ContribEmblem

**Defined:** January 27, 2026
**Core Value:** Developers get beautiful, game-inspired visual badges of their GitHub contributions that update and rotate automatically every week, embeddable anywhere on GitHub.

## v1 Requirements

Requirements for initial release. Each maps to roadmap phases.

### GitHub Actions Infrastructure

- [ ] **ACTNS-01**: Action runs on scheduled trigger (weekly) without manual intervention
- [ ] **ACTNS-02**: Action prevents infinite trigger loops via [skip ci] commit message and path filters
- [ ] **ACTNS-03**: Action has contents:write permission for committing generated images
- [ ] **ACTNS-04**: Action configures git user.name and user.email in workflow setup

### Configuration Management

- [ ] **CONFG-01**: User can define username in YAML config file
- [ ] **CONFG-02**: User can select which metrics to display (commits, PRs, issues, reviews, stars)
- [ ] **CONFG-03**: User can define emblem rotation list in YAML config
- [ ] **CONFG-04**: Config validation provides clear error messages for invalid formats
- [ ] **CONFG-05**: Config parser validates emblem IDs exist in catalog
- [ ] **CONFG-06**: Example config file provided in repository

### GitHub Stats Collection

- [ ] **STATS-01**: Action fetches GitHub contribution stats via GraphQL API
- [ ] **STATS-02**: Stats filtered to current calendar year only (Jan 1 - Dec 31)
- [ ] **STATS-03**: API responses cached with 24hr expiry to avoid rate limits
- [ ] **STATS-04**: Rate limit headers monitored and logged (X-RateLimit-Remaining)
- [ ] **STATS-05**: Stats include: commits, PRs, issues, code reviews, stars received

### Emblem Integration

- [ ] **EMBLM-01**: Action fetches emblem artwork from Bungie API with proper authentication
- [ ] **EMBLM-02**: Action randomly selects one emblem from user's rotation list weekly
- [ ] **EMBLM-03**: Random selection seeded by week number for consistency within week
- [ ] **EMBLM-04**: Fallback emblem used if fetch fails or emblem ID invalid

### Image Generation

- [ ] **IMAGE-01**: Action generates PNG image with emblem background and stat overlays
- [ ] **IMAGE-02**: Power Level (sum of all 5 metrics) displayed prominently on image
- [ ] **IMAGE-03**: Individual metric values and icons displayed below Power Level
- [ ] **IMAGE-04**: Text rendered with outline/stroke for contrast on variable backgrounds
- [ ] **IMAGE-05**: Image dimensions optimized for GitHub markdown display (800x400px)
- [ ] **IMAGE-06**: Generated image saved with stable filename that overwrites previous version

### File Operations

- [ ] **FILES-01**: Action commits generated PNG to repository
- [ ] **FILES-02**: Action updates README.md with marker-based injection (HTML comments)
- [ ] **FILES-03**: README update includes emblem image and last updated timestamp
- [ ] **FILES-04**: Commit skipped if image hash unchanged from previous run

## v2 Requirements

Deferred to future release. Tracked but not in current roadmap.

### Custom Emblems

- **CUSTM-01**: User can upload custom emblem artwork (requires moderation)
- **CUSTM-02**: Custom emblems stored in repository assets directory

### Achievement System

- **ACHVT-01**: User earns achievement badges for milestones (100 commits, 10 PRs merged, etc.)
- **ACHVT-02**: Achievement icons displayed alongside stats

### Multi-User Support

- **MULTI-01**: Single config can generate badges for multiple team members
- **MULTI-02**: Team aggregate stats displayed on separate team badge

### Private Repository Stats

- **PRIVT-01**: Action supports private repository contributions via PAT
- **PRIVT-02**: Self-hosted deployment guide for private repo stats

### Organization Filtering

- **ORGFL-01**: User can filter stats to specific organizations (e.g., "only CNCF contributions")
- **ORGFL-02**: Multiple organization filters supported in config

## Out of Scope

Explicitly excluded. Documented to prevent scope creep.

| Feature | Reason |
|---------|--------|
| All-time stats | Calendar year filtering keeps profile relevant, annual reset prevents stale numbers |
| Real-time updates | Weekly schedule sufficient, hourly/daily updates waste API quota |
| Custom metric formulas | 5 core metrics proven in research, custom formulas add complexity without value |
| Social sharing to other platforms | GitHub-focused tool, embedding via raw.githubusercontent.com is sufficient |
| Historical trend graphs | Adds complexity, focus is on current year snapshot with Destiny aesthetic |

## Traceability

Which phases cover which requirements. Updated during roadmap creation.

| Requirement | Phase | Status |
|-------------|-------|--------|
| ACTNS-01 | Phase 1 | Pending |
| ACTNS-02 | Phase 1 | Pending |
| ACTNS-03 | Phase 1 | Pending |
| ACTNS-04 | Phase 1 | Pending |
| STATS-01 | Phase 2 | Pending |
| STATS-02 | Phase 2 | Pending |
| STATS-03 | Phase 2 | Pending |
| STATS-04 | Phase 2 | Pending |
| STATS-05 | Phase 2 | Pending |
| EMBLM-01 | Phase 3 | Pending |
| EMBLM-02 | Phase 3 | Pending |
| EMBLM-03 | Phase 3 | Pending |
| EMBLM-04 | Phase 3 | Pending |
| IMAGE-01 | Phase 4 | Pending |
| IMAGE-02 | Phase 4 | Pending |
| IMAGE-03 | Phase 4 | Pending |
| IMAGE-04 | Phase 4 | Pending |
| IMAGE-05 | Phase 4 | Pending |
| IMAGE-06 | Phase 4 | Pending |
| FILES-01 | Phase 5 | Pending |
| FILES-02 | Phase 5 | Pending |
| FILES-03 | Phase 5 | Pending |
| FILES-04 | Phase 5 | Pending |
| CONFG-01 | Phase 6 | Pending |
| CONFG-02 | Phase 6 | Pending |
| CONFG-03 | Phase 6 | Pending |
| CONFG-04 | Phase 6 | Pending |
| CONFG-05 | Phase 6 | Pending |
| CONFG-06 | Phase 6 | Pending |

**Coverage:**
- v1 requirements: 30 total
- Mapped to phases: 30
- Unmapped: 0 ✓

---
*Requirements defined: January 27, 2026*
*Last updated: January 27, 2026 after initial definition*
