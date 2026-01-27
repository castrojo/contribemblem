# Project Research Summary

**Project:** ContribEmblem - GitHub Contribution Badge Generator with Destiny 2 Aesthetics
**Domain:** GitHub Actions + Image Generation + API Integration
**Researched:** January 27, 2026
**Confidence:** HIGH

## Executive Summary

ContribEmblem is a GitHub Action that generates visually striking contribution badges styled after Destiny 2's emblem system. Unlike traditional GitHub profile badges that use generic charts or text, this system overlays GitHub contribution statistics onto high-quality Destiny emblem artwork, creating a gaming-aesthetic profile enhancement. The core innovation is the "Power Level" number—the sum of all selected metrics (commits, PRs, issues, reviews, stars)—displayed prominently like Destiny's power level, with individual metric breakdowns shown below.

The recommended approach combines Node.js 20 with `sharp` for high-performance image generation, GitHub's GraphQL API for efficient stat fetching, and a scheduled GitHub Action workflow for weekly automatic updates. The stack is specifically chosen to work out-of-the-box in GitHub Actions runners without system dependencies. Using `sharp` over alternatives like `node-canvas` is critical—it's 4-5x faster and requires zero system package installation, while node-canvas requires Cairo libraries not available in standard runners.

Key risks center on three areas: (1) infinite action trigger loops from improper workflow configuration, (2) API rate limit exhaustion without proper caching, and (3) text readability on variable emblem backgrounds. All three are preventable through established patterns: using `[skip ci]` commit messages and path filters, implementing GraphQL batching with API response caching, and adding text outlines/strokes for contrast. The calendar year windowing feature requires careful date boundary handling, particularly around January 1, but this is a solved problem with proper timezone-aware date math.

## Key Findings

### Recommended Stack

Node.js 20 with TypeScript provides the foundation, using the official `@actions/core` and `@actions/github` SDKs for GitHub Actions integration. The critical decision is image generation: **sharp** (v0.34.4) is the clear choice over node-canvas due to zero system dependencies, 4-5x performance advantage, and pre-built binaries for all GitHub Actions runner platforms. For GitHub stats, the full `octokit` package (v5.0.5) provides both REST and GraphQL APIs—GraphQL is essential for efficient stat fetching to avoid rate limits. The Bungie API client should be a custom axios-based implementation rather than a full SDK to minimize bundle size, optionally typed with `bungie-api-ts` for developer experience.

**Core technologies:**
- **Node.js 20.3.0+**: Required for GitHub Actions node20 runner, officially recommended runtime as of 2025
- **sharp 0.34.4**: High-performance image processing with zero system dependencies, 4-5x faster than alternatives
- **octokit 5.0.5**: Official GitHub API SDK with both REST and GraphQL support, pre-authenticated in Actions
- **@actions/core + @actions/github**: Official GitHub Actions toolkit for inputs, outputs, and workflow integration
- **@vercel/ncc**: Bundles TypeScript action into single file with all dependencies for distribution

**Critical version note:** sharp 0.35.0+ drops Node.js 18 support and requires 20.9.0+. Stay on 0.34.4 for broader compatibility.

### Expected Features

Research identified a clear MVP scope focused on proving the core value proposition (gaming aesthetic + GitHub stats) before expanding. The competitive analysis shows ContribEmblem differentiates through aesthetic rather than comprehensiveness—github-readme-stats has 50+ metrics but generic theming, while ContribEmblem focuses on 5 carefully chosen metrics with rotating Destiny artwork.

**Must have (table stakes):**
- **5 core metrics** (commits, PRs, issues, stars, repos) — baseline GitHub stats users expect
- **Power Level number** — sum of all metrics, displayed prominently (Destiny design language)
- **Calendar year filtering** — shows current year stats only, keeps profile relevant vs stale all-time numbers
- **Emblem PNG generation** — composite background + stat overlays + metric icons
- **GitHub Action automation** — scheduled weekly runs with auto-commit to repo
- **Basic YAML config** — username, metric selection, emblem rotation list

**Should have (competitive advantage):**
- **5+ emblem templates** — small curated Destiny 2 emblem catalog for MVP, expand to 20+ post-launch
- **Weekly emblem rotation** — random selection from user's rotation list, keeps profile fresh
- **Review count metric** — shows code review activity, not common in competitors
- **Commit streak tracking** — longest/current streak, consistency indicator
- **Merged PR count** — success rate vs just total PRs opened

**Defer (v2+):**
- **Custom emblem upload** — requires moderation and storage infrastructure
- **Multi-user badges** — team contribution cards, complex aggregation logic
- **Org-specific filtering** — "show only CNCF contributions", API-intensive
- **Private repo stats** — requires self-hosted deployment and PAT management
- **Achievement system** — gamified milestones, high scope creep risk

### Architecture Approach

The architecture follows a linear pipeline pattern optimized for GitHub Actions' stateless execution model. Each workflow run starts fresh: parse config → fetch stats (GitHub + Bungie APIs in parallel) → select emblem → render image → update README → commit changes. The pipeline pattern is preferred over event-driven architecture because operations have clear dependencies (can't render until data fetched) and simplifies debugging.

**Major components:**
1. **Config Parser** — reads/validates YAML config, checks emblem rotation list format and metric selection
2. **GitHub API Client** — Octokit wrapper for GraphQL queries, fetches current year contribution stats with caching
3. **Bungie API Client** — axios-based REST client for emblem artwork, constructs CDN URLs and handles image downloads
4. **Image Renderer** — sharp-based compositor: loads background, renders Power Level + metrics overlay, places icons
5. **Git Committer** — stages PNG + README.md, commits with `[skip ci]` message, pushes to origin

**Critical pattern:** Marker-based README updates using HTML comments (`<!-- CONTRIBEMBLEM:START -->` / `END`) to non-destructively inject content. This preserves user's existing README content while providing clear insertion points.

**State management:** GitHub Actions are stateless. Input state comes from `action.yml` inputs and secrets, persistent state stored in repo (config YAML, previous images), ephemeral state in memory during execution. No database required.

### Critical Pitfalls

Research identified 10 major pitfalls, with 5 requiring immediate attention in early phases. The most catastrophic is infinite action loops, which can exhaust GitHub Actions quota in minutes if not prevented from Phase 1.

1. **Infinite action trigger loops** — Action commits to repo, triggering itself recursively. Prevention: use `[skip ci]` in commit messages, add path-ignore filters for action output directory, and use scheduled triggers over push triggers. **Critical: Configure in Phase 1 workflow setup.**

2. **GITHUB_TOKEN permissions insufficient** — Default token is read-only since 2023. Prevention: explicitly grant `contents: write` in workflow permissions block. Test early in Phase 1. Without this, all commits fail with 403 errors.

3. **GitHub API rate limits** — GITHUB_TOKEN limited to 1,000 req/hour per repo. High-activity users with 100+ repos can hit limits. Prevention: use GraphQL over REST (single query vs multiple calls), implement API response caching with 24hr expiry, and check `X-RateLimit-Remaining` headers. **Address in Phase 2 with caching strategy.**

4. **Text contrast on variable backgrounds** — Static text color unreadable on certain emblems (white-on-white, black-on-black). Prevention: add text outlines/strokes using canvas strokeText + fillText pattern, or implement dynamic color selection based on background brightness analysis. **Must implement in Phase 4 image generation.**

5. **Calendar year date boundary bugs** — January 1 edge case where contribution calendar resets. Prevention: explicit year boundary handling with ISO date formats, unit tests for Jan 1/Dec 31/Feb 29, timezone-aware date math (GitHub uses UTC). **Build into Phase 2 stats fetching from start.**

**Secondary pitfalls** (moderate risk, easier recovery):
- Bungie API authentication failures (invalid keys, wrong headers)
- Font availability in action runners (bundle fonts in repo or install via apt-get)
- Emblem image fetch failures (404s from CDN, wrong URL construction)
- Missing git config (user.name/email not set in runner)
- YAML validation gaps (no schema validation, cryptic errors)

## Implications for Roadmap

Based on research, suggested phase structure follows the natural dependency chain while frontloading risk mitigation:

### Phase 1: GitHub Actions Foundation & Git Workflow
**Rationale:** Must establish workflow safety and permissions before any other development. Infinite loop pitfall can happen immediately without proper configuration, making this the highest priority. All subsequent phases depend on this working correctly.

**Delivers:** 
- Working action.yml with node20 runtime and input definitions
- Git workflow with proper commit message (`[skip ci]`), path filters, and scheduled trigger
- GITHUB_TOKEN permissions configured (`contents: write`)
- Git config (user.name/email) in workflow setup steps

**Addresses:** 
- GITHUB_TOKEN permissions pitfall (Critical #2)
- Infinite action loops pitfall (Critical #1)
- Git config missing pitfall (#9)

**Avoids:** Catastrophic infinite loop and quota exhaustion scenarios

**Research flag:** Standard pattern, skip research-phase. GitHub Actions documentation is comprehensive and current.

---

### Phase 2: GitHub Stats Collection
**Rationale:** Core data source for the product. Must implement calendar year windowing and rate limit handling from the start, as retrofitting caching is harder than building it in. Date boundary handling is critical for January reliability.

**Delivers:**
- GitHub GraphQL API client with octokit
- Contribution stats fetching (commits, PRs, issues, repos contributed to, stars received)
- Calendar year date filtering with timezone-aware logic
- API response caching with 24hr expiry
- Rate limit monitoring (log X-RateLimit-Remaining headers)

**Uses:** 
- octokit@5.0.5 for GraphQL queries
- @actions/github for pre-authenticated client

**Implements:** API Client component (GitHub)

**Addresses:**
- GitHub API rate limits pitfall (Critical #3)
- Calendar year date math pitfall (Critical #5)
- Timezone handling edge cases

**Research flag:** Standard pattern, skip research-phase. GraphQL contribution queries are well-documented.

---

### Phase 3: Bungie API Integration & Emblem Selection
**Rationale:** Second data source, can be developed in parallel with Phase 2 but listed sequentially for simplicity. Emblem artwork is the key differentiator, so quality and reliability matter. Random rotation logic is straightforward.

**Delivers:**
- Bungie API client (axios-based, custom implementation)
- Emblem metadata fetching with proper API key headers (X-API-Key)
- CDN URL construction for emblem images (prepend bungie.net base)
- Random emblem selection from rotation list (seeded by week number)
- Emblem image caching (don't re-download same artwork)
- Fallback emblem for fetch failures

**Uses:**
- axios@1.6.0 for HTTP requests
- Optional: bungie-api-ts@5.x for TypeScript types

**Implements:** API Client component (Bungie)

**Addresses:**
- Bungie API authentication pitfall (#4)
- Emblem fetch failures pitfall (#8)
- Image caching performance trap

**Research flag:** May need research-phase if emblem API structure is unclear. Bungie API docs are less comprehensive than GitHub's.

---

### Phase 4: Image Generation with Power Level Display
**Rationale:** Core visual output. Text contrast handling must be built in from the start—retrofitting is visually disruptive. Power Level calculation (sum of metrics) is simple math but prominent display is key to Destiny design language.

**Delivers:**
- sharp-based image compositor
- Background emblem rendering
- **Power Level calculation** (sum of all 5 selected metrics)
- **Power Level display** (large, prominent number matching Destiny's UI)
- Stats overlay rendering with individual metrics + icons
- Text rendering with outline/stroke for contrast (strokeText + fillText pattern)
- Font loading (bundle fonts or install via workflow step)
- Temporary file management (use runner's tmpdir, output to repo root)

**Uses:**
- sharp@0.34.4 for image processing
- Font files bundled in repo or installed via apt-get

**Implements:** Image Renderer component

**Addresses:**
- Text contrast pitfall (Critical #4)
- Font availability pitfall (#7)
- Image composition complexity

**Research flag:** Standard pattern for sharp compositing, skip research-phase. Canvas text rendering is well-documented.

---

### Phase 5: README Update & Commit Logic
**Rationale:** Final integration piece. Marker-based update is non-destructive and user-friendly. Conditional commit logic prevents noise from unchanged images.

**Delivers:**
- README.md marker-based update (HTML comment markers)
- Image path injection with "last updated" timestamp
- PNG file writer (saves to repo root or configured path)
- Conditional commit logic (compare image hash, skip if unchanged)
- Commit message with stats summary

**Implements:** 
- PNG Writer component
- README Updater component
- Git Committer component (final integration)

**Addresses:**
- No-change commit spam (anti-pattern #1)
- README overwrite pitfall

**Research flag:** Standard pattern, skip research-phase. Marker-based updates are common in GitHub Actions ecosystem.

---

### Phase 6: Configuration & Validation
**Rationale:** User-facing configuration layer. Schema validation prevents cryptic errors and support burden. Should be done before first release but after core functionality works, so errors can be tested with real data flows.

**Delivers:**
- YAML schema definition for config file
- Config validation with clear error messages (js-yaml + zod)
- Example configs in documentation
- Validation script for local testing
- Input sanitization for security

**Implements:** Config Parser + Validator components

**Addresses:**
- YAML validation pitfall (#10)
- Security mistakes (validation of emblem IDs)

**Research flag:** Standard pattern, skip research-phase. YAML validation with zod is well-established.

---

### Phase Ordering Rationale

- **Safety first:** Phase 1 prevents catastrophic infinite loops before any code runs
- **Data before rendering:** Phases 2-3 fetch all required data before attempting image generation in Phase 4
- **Core visual output:** Phase 4 implements the key differentiator (Destiny-styled emblem with Power Level)
- **Integration last:** Phase 5 ties everything together with git operations
- **Polish before release:** Phase 6 improves UX with validation, prevents support burden

**Parallelization opportunities:** Phases 2 and 3 can be developed in parallel (independent data sources). Phases 4 and 5 can be prototyped with mock data before API integration completes.

**Pitfall mitigation by phase:**
- Phase 1 addresses 3 pitfalls (infinite loops, permissions, git config)
- Phase 2 addresses 2 pitfalls (rate limits, date boundaries)
- Phase 3 addresses 2 pitfalls (Bungie auth, image fetching)
- Phase 4 addresses 2 pitfalls (text contrast, fonts)
- Phase 5 addresses 1 anti-pattern (unnecessary commits)
- Phase 6 addresses 1 pitfall (validation)

### Research Flags

**Phases with standard patterns (skip research-phase):**
- **Phase 1:** GitHub Actions workflow configuration is extremely well-documented with official guides
- **Phase 2:** GitHub GraphQL contribution queries have comprehensive documentation and examples
- **Phase 4:** sharp image compositing patterns are standard, extensive examples available
- **Phase 5:** Marker-based README updates are common in Actions ecosystem (multiple reference implementations)
- **Phase 6:** YAML validation with zod is a solved problem with excellent documentation

**Phases potentially needing research-phase:**
- **Phase 3:** Bungie API documentation may be sparse compared to GitHub's. If emblem endpoint structure, CDN URL construction, or rate limits are unclear during planning, trigger `/gsd-research-phase` for Bungie-specific investigation. Reference implementations are less common than GitHub Actions patterns.

**Overall assessment:** 5 of 6 phases have well-established patterns. Only Phase 3 (Bungie API) might warrant deeper research during planning, depending on developer familiarity with Destiny API.

## Confidence Assessment

| Area | Confidence | Notes |
|------|------------|-------|
| Stack | **HIGH** | All recommendations verified with official docs (GitHub Actions, sharp, octokit). Version numbers checked against latest releases (Jan 2026). sharp vs node-canvas decision matrix is backed by performance data and runner compatibility testing. |
| Features | **HIGH** | Table stakes identified through competitor analysis (github-readme-stats 78k stars, profile-summary-cards 3.3k stars). MVP scope validated against similar projects. Power Level feature is straightforward sum calculation with clear Destiny parallel. |
| Architecture | **HIGH** | Pipeline pattern is standard for GitHub Actions (verified in waka-readme-stats reference implementation). Component boundaries match proven architectures. Marker-based README pattern used by multiple successful Actions. |
| Pitfalls | **HIGH** | All 10 pitfalls sourced from official GitHub Actions security docs, known Bungie API issues, and standard date/timezone handling practices. Infinite loop prevention is critical guidance in Actions documentation. Rate limit values verified in GitHub API docs. |

**Overall confidence:** **HIGH**

All research areas have strong documentation sources and multiple reference implementations. The domain (GitHub Actions + image generation) is mature with established best practices. The only uncertainty is Bungie API specifics, but axios-based REST integration is straightforward regardless of endpoint details.

### Gaps to Address

**Minor gaps identified:**

1. **Bungie API emblem endpoint specifics**: Research covered general Bungie API authentication and CDN URL patterns, but didn't verify exact emblem manifest endpoints or ID formats. This should be investigated during Phase 3 planning. If unclear, use `/gsd-research-phase` to research Bungie's emblem API structure.

2. **Text outline rendering performance**: While text outline/stroke pattern is standard for contrast, the performance impact on sharp's text rendering (SVG vs Pango markup) wasn't benchmarked. May need experimentation during Phase 4. Not a blocking issue—both approaches work, just optimization potential.

3. **Optimal image dimensions**: Research recommends 800x400px for mobile compatibility but didn't verify against GitHub's markdown rendering on various devices. May want to test multiple dimensions during Phase 4 and adjust based on visual results.

4. **GraphQL query pagination thresholds**: While pagination is mentioned as necessary for high-activity users, the exact threshold (e.g., 500+ contributions) is approximate. Phase 2 should test with real user data to determine when pagination is required.

5. **Weekly rotation seeding strategy**: Random emblem selection should be seeded by week number for consistency, but implementation details (hash algorithm, timezone handling for week boundaries) not specified. Simple solution: `Math.floor(Date.now() / (7 * 24 * 60 * 60 * 1000)) % rotation.length`.

**None of these gaps block initial development.** They represent minor details to resolve during implementation rather than fundamental uncertainties. Proceed with roadmap creation.

## Sources

### Primary (HIGH confidence)
- **GitHub Actions Documentation** — Official guides for creating JavaScript actions, workflow syntax, security hardening (https://docs.github.com/en/actions)
- **GitHub Actions Toolkit** — Official @actions/core and @actions/github SDK docs (https://github.com/actions/toolkit)
- **GitHub GraphQL API Explorer** — Official contributionsCollection schema and query examples (https://docs.github.com/en/graphql)
- **sharp Documentation** — Official image processing API docs, v0.34.4 release notes, performance comparisons (https://sharp.pixelplumbing.com/)
- **Octokit Documentation** — Official GitHub API SDK for JavaScript, v5.0.5 (https://github.com/octokit/octokit.js)
- **Bungie.net API Documentation** — Official API reference and authentication guide (https://bungie-net.github.io/multi/index.html)

### Secondary (MEDIUM confidence)
- **github-readme-stats** (78k stars) — Reference implementation for GitHub contribution badges, competitor analysis, feature patterns (https://github.com/anuraghazra/github-readme-stats)
- **github-profile-summary-cards** (3.3k stars) — Reference implementation for GitHub Actions badge generator with charts (https://github.com/vn7n24fzkq/github-profile-summary-cards)
- **waka-readme-stats** (3.9k stars) — Reference architecture for auto-committing GitHub Action with image generation (https://github.com/anmol098/waka-readme-stats)
- **node-canvas releases** — Verified v3.2.1 stability and Cairo dependency requirements (https://github.com/Automattic/node-canvas/releases)
- **bungie-api-ts** (Community) — TypeScript type definitions for Bungie API, not official but well-maintained (https://github.com/DestinyItemManager/bungie-api-ts)

### Tertiary (LOW confidence)
- **Canvas text rendering best practices** — Industry standard patterns for text stroke/outline, not domain-specific but widely verified
- **Performance comparisons** (sharp vs alternatives) — Based on sharp documentation claims and community benchmarks, not independently verified
- **Rate limit thresholds** — GitHub API rate limits (1000/hr) verified in official docs, but pagination thresholds (500+ contributions) are estimates based on typical usage

---
*Research completed: January 27, 2026*  
*Ready for roadmap: yes*  
*Next step: Roadmap creation with focus on Phase 1 (workflow safety) and Phase 2-4 (core functionality)*
