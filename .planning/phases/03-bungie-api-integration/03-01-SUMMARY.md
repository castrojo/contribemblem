---
phase: 03-bungie-api-integration
plan: 01
status: complete
duration: 5 minutes
completed: 2026-01-28

tech_stack:
  added: []
  modified: []

patterns:
  - ISO week-based deterministic selection (UTC timezone)
  - SHA256 seeded randomness for consistency within week
  - Fallback emblem on API failure (resilience pattern)
  - Manifest caching with date-based keys (~100MB file)
  - Non-blocking rate limit monitoring

files_created:
  - scripts/select-emblem.sh
  - scripts/fetch-emblem.sh
  - data/emblem-config.json

files_modified:
  - .github/workflows/update-emblem.yml

key_types_exported: []

key_functions_exported:
  - select-emblem.sh (weekly emblem selection)
  - fetch-emblem.sh (Bungie API fetcher)

decisions:
  - Bash scripts maintain consistency with Phase 2 patterns
  - SHA256(ISO week) provides deterministic weekly rotation
  - Fallback emblem: "The Seventh Column" (hash 1409726931)
  - Manifest caching critical due to 100MB size
  - Error handling: exit 1 on failure, workflow catches with fallback

blockers: []

affects:
  - phase: 04-image-generation
    reason: "data/emblem.jpg now available as background source"
  - phase: 05-readme-update
    reason: "Emblem rotation visible in weekly badge updates"

subsystem: data-pipeline
---

# Phase 3, Plan 1 Summary: Bungie API Integration

## What Was Built

**Weekly Emblem Rotation System**
- `scripts/select-emblem.sh` — Deterministic emblem selection based on ISO week number
- `scripts/fetch-emblem.sh` — Bungie API client with authentication and caching
- `data/emblem-config.json` — Emblem rotation configuration with fallback
- Updated `.github/workflows/update-emblem.yml` with emblem fetching pipeline

## Key Technical Decisions

### 1. ISO Week-Based Deterministic Selection
**Decision:** Use `date -u +%G-W%V` for ISO week calculation, hash with SHA256, select via modulo

**Why:** Ensures same emblem throughout each calendar week (Sunday-Saturday), changes predictably on Sunday midnight UTC. UTC timezone matches Phase 1 schedule and prevents mid-week changes due to timezone differences.

**Pattern:**
```bash
ISO_WEEK=$(date -u +%G-W%V)  # e.g., "2026-W04"
HASH=$(echo -n "$ISO_WEEK" | sha256sum | cut -c1-8)
INDEX=$((HASH_DECIMAL % ARRAY_LENGTH))
```

### 2. Bash Scripts (Not Node.js/TypeScript)
**Decision:** Implement in bash using curl + jq

**Why:** Maintains pattern consistency with Phase 2 (fetch-stats.sh, process-stats.sh). Bungie API is RESTful JSON — no need for Node.js runtime. Reduces dependencies and GitHub Actions startup time.

### 3. Manifest Caching Strategy
**Decision:** Cache manifest.json with date-based key (`bungie-manifest-YYYYMMDD`)

**Why:** Manifest is ~100MB and changes infrequently (~every 3 months). Daily cache check is sufficient. Follows Phase 2 cache key pattern for consistency.

### 4. Fallback Emblem on API Failure
**Decision:** Workflow catches fetch-emblem.sh exit code 1, downloads fallback emblem directly

**Why:** Resilience pattern from Phase 2 (non-blocking errors). Action should never fail due to external API issues. "The Seventh Column" emblem (1409726931) is iconic and always available.

**Workflow pattern:**
```yaml
if echo "$HASH" | ./scripts/fetch-emblem.sh; then
  echo "✓ Success"
else
  curl -o data/emblem.jpg [fallback-url]
fi
```

## Architecture Integration

### Data Pipeline Flow
```
Phase 2 (Stats) ──────┐
                       ├──> Phase 4 (Image Generation)
Phase 3 (Emblem) ─────┘

Phase 3 Output: data/emblem.jpg (artwork background)
Phase 2 Output: data/stats.json (metrics overlay)
Phase 4 Input: Composite both into final badge
```

### Workflow Position
```yaml
1. Fetch GitHub Stats (Phase 2)
2. Select Weekly Emblem (Phase 3) ← NEW
3. Fetch Emblem Image (Phase 3) ← NEW
4. Generate Badge (Phase 4) ← NEXT
5. Update README (Phase 5)
6. Commit Changes (Phase 6)
```

## Verification Results

### Task 1: Emblem Selector
✅ Deterministic selection (same output within week): `2448092419`
✅ Valid numeric hash format
✅ Fallback working when config missing: `1409726931`
✅ Executable permissions set

### Task 2: Bungie API Fetcher
✅ Error handling for missing BUNGIE_API_KEY
✅ Executable permissions set
✅ Proper error messages with setup instructions
⚠️ API testing requires user setup (BUNGIE_API_KEY secret)

### Task 3: Workflow Integration
✅ BUNGIE_API_KEY configured in workflow
✅ Selection step present (select-emblem.sh)
✅ Fetch step present (fetch-emblem.sh)
✅ Manifest cache configured with date-based key
✅ Default emblem config creation step
✅ Fallback handling on fetch failure

## User Setup Required

**Before this workflow runs successfully, user must:**

1. **Create Bungie.net Developer Account**
   - Visit: https://www.bungie.net/en/Application
   - Sign in with Bungie.net account (or create one)

2. **Register Application for API Key**
   - Application Name: `ContribEmblem`
   - OAuth Client Type: Not Applicable (read-only public data)
   - Redirect URL: Not required
   - Scope: Not required (public manifest data only)
   - Click "Create New App" to get API key

3. **Add API Key to GitHub Repository Secrets**
   - Navigate to: Repository Settings → Secrets and variables → Actions
   - Click "New repository secret"
   - Name: `BUNGIE_API_KEY`
   - Value: [paste API key from Bungie.net]
   - Click "Add secret"

**Why required:** Bungie API requires X-API-Key header for all requests. This is free and does not require OAuth since we only access public manifest data.

## Requirements Satisfied

✅ **EMBLM-01:** Action fetches emblem artwork from Bungie API with proper authentication
- fetch-emblem.sh includes X-API-Key header
- Workflow provides BUNGIE_API_KEY from secrets

✅ **EMBLM-02:** Action randomly selects one emblem from user's rotation list weekly
- select-emblem.sh reads data/emblem-config.json rotation array
- Workflow executes selection in workflow step

✅ **EMBLM-03:** Random selection seeded by week number for consistency within week
- ISO week calculation: `date -u +%G-W%V`
- SHA256 hash of week string provides deterministic seed
- UTC timezone ensures consistency with Sunday midnight schedule

✅ **EMBLM-04:** Fallback emblem used if fetch fails or emblem ID invalid
- select-emblem.sh returns fallback when config missing/empty
- fetch-emblem.sh exits with code 1 on errors
- Workflow catches failure and downloads fallback emblem directly

## Files for Phase 4

Phase 4 (Image Generation) can now access:
- `data/stats.json` — GitHub contribution metrics with Power Level
- `data/emblem.jpg` — Destiny 2 emblem artwork (800x400px typically)
- Both files regenerated/updated weekly by workflow

## Known Limitations

1. **Manifest Size:** ~100MB download on first run. Cached for 24 hours, but initial workflow run will be slower.
2. **API Dependency:** If Bungie.net is down, fallback emblem always used. No retry logic currently.
3. **Fixed Rotation List:** User must manually edit `data/emblem-config.json` to change rotation. Could be enhanced with UI in future.
4. **No Emblem Metadata:** Not fetching emblem name/description for badge. Only using artwork background.

## Next Phase Preview

**Phase 4: Image Generation with Power Level**
- Use `sharp` library to composite emblem + stats
- Render Power Level prominently on emblem background
- Output: 800x400px PNG badge
- Text contrast/stroke for readability over varying emblem colors

**Input files ready:**
- `data/emblem.jpg` ← Phase 3
- `data/stats.json` ← Phase 2

**Technical considerations:**
- Node.js + TypeScript for sharp library (unlike bash scripts for data fetching)
- Font selection for Power Level display
- Text positioning and contrast handling
- Image size optimization for GitHub embedding
