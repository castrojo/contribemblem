# Roadmap: ContribEmblem

## Overview

ContribEmblem transforms GitHub contribution stats into beautiful Destiny 2-styled emblem badges through a six-phase implementation. Starting with workflow safety (Phase 1), we build the data pipeline (Phases 2-3), create the visual rendering engine (Phase 4), integrate file operations (Phase 5), and polish the user experience with validation (Phase 6). Each phase delivers a complete, verifiable capability that builds toward the core value: developers get game-inspired badges that auto-update weekly and embed anywhere on GitHub.

## Phases

**Phase Numbering:**
- Integer phases (1, 2, 3): Planned milestone work
- Decimal phases (2.1, 2.2): Urgent insertions (marked with INSERTED)

Decimal phases appear between their surrounding integers in numeric order.

- [x] **Phase 1: GitHub Actions Foundation** - Safe workflow with scheduled trigger and commit loop prevention
- [ ] **Phase 2: GitHub Stats Collection** - Current year contribution data with rate limit protection
- [ ] **Phase 3: Bungie API Integration** - Emblem artwork fetching with weekly rotation logic
- [ ] **Phase 4: Image Generation with Power Level** - Composite rendering with prominent power display
- [ ] **Phase 5: README Update & Commit** - Marker-based injection with conditional commits
- [ ] **Phase 6: Configuration & Validation** - Schema validation with clear error messages

## Phase Details

### Phase 1: GitHub Actions Foundation
**Goal**: Workflow runs safely on schedule without triggering infinite loops or permission failures
**Depends on**: Nothing (first phase)
**Requirements**: ACTNS-01, ACTNS-02, ACTNS-03, ACTNS-04
**Success Criteria** (what must be TRUE):
  1. Action runs on weekly schedule (cron trigger) without manual intervention
  2. Action commits to repository without triggering itself recursively
  3. Action has write permissions to commit files (no 403 errors)
  4. Git commits include proper author name and email metadata
**Plans**: 1 plan

Plans:
- [x] 01-01-PLAN.md — Create scheduled workflow with loop prevention, permissions, and git config (COMPLETE: 2026-01-27)

### Phase 2: GitHub Stats Collection
**Goal**: Action fetches current year GitHub contribution stats efficiently without hitting rate limits
**Depends on**: Phase 1
**Requirements**: STATS-01, STATS-02, STATS-03, STATS-04, STATS-05
**Success Criteria** (what must be TRUE):
  1. Action fetches 5 metrics (commits, PRs, issues, reviews, stars) via GraphQL API
  2. Stats filtered to current calendar year (Jan 1 - Dec 31 of current year)
  3. API responses cached for 24 hours to prevent redundant fetches
  4. Rate limit monitoring logs remaining quota without exhausting limits
  5. Date boundary handling works correctly on January 1 (year rollover)
**Plans**: TBD

Plans:
- [ ] 02-01: [TBD during phase planning]

### Phase 3: Bungie API Integration
**Goal**: Action retrieves Destiny emblem artwork and selects one randomly per week
**Depends on**: Phase 1 (can develop parallel to Phase 2)
**Requirements**: EMBLM-01, EMBLM-02, EMBLM-03, EMBLM-04
**Success Criteria** (what must be TRUE):
  1. Action fetches emblem images from Bungie API with valid authentication headers
  2. Action selects one emblem randomly from user's rotation list each week
  3. Same emblem selected consistently throughout the week (seeded by week number)
  4. Fallback emblem displayed if API fetch fails or emblem ID invalid
**Plans**: TBD

Plans:
- [ ] 03-01: [TBD during phase planning]

### Phase 4: Image Generation with Power Level
**Goal**: Action generates PNG with Destiny-styled stat overlay featuring prominent Power Level
**Depends on**: Phases 2 and 3
**Requirements**: IMAGE-01, IMAGE-02, IMAGE-03, IMAGE-04, IMAGE-05, IMAGE-06
**Success Criteria** (what must be TRUE):
  1. Action composites emblem background with stat overlays into 800x400px PNG
  2. Power Level (sum of all 5 metrics) displayed prominently in Destiny design language
  3. Individual metric values and icons displayed clearly below Power Level
  4. Text readable on all emblem backgrounds (outline/stroke for contrast)
  5. Generated image saved with stable filename that overwrites previous version
**Plans**: TBD

Plans:
- [ ] 04-01: [TBD during phase planning]

### Phase 5: README Update & Commit
**Goal**: Action non-destructively updates README and commits generated PNG to repository
**Depends on**: Phase 4
**Requirements**: FILES-01, FILES-02, FILES-03, FILES-04
**Success Criteria** (what must be TRUE):
  1. Action commits generated PNG to repository with [skip ci] message
  2. README updated between HTML comment markers without overwriting user content
  3. README injection includes emblem image path and last updated timestamp
  4. No commit created if image hash unchanged from previous run
**Plans**: TBD

Plans:
- [ ] 05-01: [TBD during phase planning]

### Phase 6: Configuration & Validation
**Goal**: Users can configure emblem rotation with schema validation preventing cryptic errors
**Depends on**: Phase 5
**Requirements**: CONFG-01, CONFG-02, CONFG-03, CONFG-04, CONFG-05, CONFG-06
**Success Criteria** (what must be TRUE):
  1. User configures username, metrics, and emblem rotation list in YAML file
  2. Config validation provides clear error messages for invalid YAML structure
  3. Validation checks emblem IDs exist before attempting API fetch
  4. Example config file demonstrates all supported options
**Plans**: TBD

Plans:
- [ ] 06-01: [TBD during phase planning]

## Progress

**Execution Order:**
Phases execute in numeric order: 1 → 2 → 3 → 4 → 5 → 6

| Phase | Plans Complete | Status | Completed |
|-------|----------------|--------|-----------|
| 1. GitHub Actions Foundation | 0/TBD | Not started | - |
| 2. GitHub Stats Collection | 0/TBD | Not started | - |
| 3. Bungie API Integration | 0/TBD | Not started | - |
| 4. Image Generation with Power Level | 0/TBD | Not started | - |
| 5. README Update & Commit | 0/TBD | Not started | - |
| 6. Configuration & Validation | 0/TBD | Not started | - |
