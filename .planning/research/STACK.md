# Stack Research

**Domain:** GitHub Actions + Image Generation + Node.js/TypeScript
**Researched:** January 27, 2026
**Confidence:** HIGH

## Recommended Stack

### Core Technologies

| Technology | Version | Purpose | Why Recommended |
|------------|---------|---------|-----------------|
| **Node.js** | >= 20.3.0 | Runtime environment | Required for GitHub Actions JavaScript actions (node20 runner). GitHub officially recommends Node 20 for actions as of 2025. |
| **TypeScript** | ^5.x | Type-safe development | Standard for modern GitHub Actions development. Provides better IDE support and catches errors at compile-time. |
| **@actions/core** | ^1.11.1 | GitHub Actions SDK (inputs/outputs/logging) | Official GitHub Actions toolkit for handling action inputs, outputs, state, and logging. Essential for any JavaScript action. |
| **@actions/github** | ^6.0.0 | GitHub API client (Octokit wrapper) | Official GitHub Actions toolkit providing pre-authenticated Octokit instance with action context. |

### GitHub API & Stats Collection

| Library | Version | Purpose | When to Use |
|---------|---------|---------|-------------|
| **octokit** | ^5.0.5 | GitHub REST/GraphQL API client | For fetching GitHub contribution stats. Includes authentication, pagination, and error handling. This is the all-in-one SDK. |
| **@octokit/rest** | (included in octokit) | REST API typed methods | Preferred over raw `octokit.request()` for better type safety and autocomplete. |
| **@octokit/graphql** | (included in octokit) | GraphQL API queries | For complex queries requiring multiple resources. More efficient than multiple REST calls. |

### Image Generation

**RECOMMENDED: sharp**

| Library | Version | Purpose | When to Use |
|---------|---------|---------|-------------|
| **sharp** | ^0.34.4 | High-performance image processing | **Primary choice.** 4-5x faster than ImageMagick. Pure Node.js, no system dependencies in GitHub Actions. Pre-built binaries for all platforms. Ideal for compositing, resizing, overlaying text. |
| **@napi-rs/canvas** | ^0.1.x (alternative) | Canvas API for Node.js | Alternative if you need Canvas 2D API specifically. Uses Rust/NAPI (faster than node-canvas). |

**NOT RECOMMENDED: node-canvas** 

| Library | Version | Avoid Because |
|---------|---------|---------------|
| ~~node-canvas~~ | v3.x | Requires Cairo system dependencies (not available in standard GitHub Actions runners). Slower than sharp. Build complexity. |

### Configuration & Data Parsing

| Library | Version | Purpose | When to Use |
|---------|---------|---------|-------------|
| **js-yaml** | ^4.1.0 | YAML parsing | For reading action configuration from YAML files (e.g., emblem mappings, user preferences). |
| **zod** | ^3.22.0 (optional) | Schema validation | Recommended for validating user configuration input. Provides type-safe parsing with detailed error messages. |

### Bungie API Client

| Approach | Library | Version | Why |
|----------|---------|---------|-----|
| **Custom client** | axios | ^1.6.0 | Bungie API is straightforward REST. Use axios for HTTP requests with interceptors for auth token refresh. Smaller bundle than a full SDK. |
| **Typed requests** | bungie-api-ts | ^5.x (optional) | Community TypeScript definitions for Bungie API. Provides excellent type safety but adds 2MB+ to bundle. |

### File Operations & Git

| Library | Version | Purpose | When to Use |
|---------|---------|---------|-------------|
| **@actions/exec** | ^1.1.1 | Execute shell commands | For running git commands (`git add`, `git commit`, `git push`). Provides better error handling than child_process. |
| **@actions/io** | ^1.1.3 | File system utilities | For cross-platform file operations (cp, mv, rmRF). Handles Windows path issues. |

### Development Tools

| Tool | Purpose | Notes |
|------|---------|-------|
| **@vercel/ncc** | Bundle TypeScript action into single file | GitHub Actions best practice. Bundles all dependencies so users don't need `node_modules`. Alternative to Rollup. |
| **ts-node** | TypeScript execution for development | Useful for local testing. Not needed in production action. |
| **prettier** | Code formatting | Standard formatter for TypeScript projects. |
| **eslint** | Linting with TypeScript support | Use `@typescript-eslint/parser` and `@typescript-eslint/eslint-plugin`. |

## Installation

```bash
# Core dependencies
npm install @actions/core@^1.11.1 @actions/github@^6.0.0 octokit@^5.0.5

# Image generation (sharp - recommended)
npm install sharp@^0.34.4

# Configuration parsing
npm install js-yaml@^4.1.0

# Bungie API (HTTP client)
npm install axios@^1.6.0

# Optional: Bungie API TypeScript definitions
npm install bungie-api-ts@^5.x

# Optional: Schema validation
npm install zod@^3.22.0

# Dev dependencies
npm install -D typescript@^5.x @vercel/ncc@^0.38.0 @types/node@^20.x @types/js-yaml@^4.0.0
npm install -D prettier eslint @typescript-eslint/parser @typescript-eslint/eslint-plugin
```

## Stack Architecture Decision: sharp vs node-canvas

### Why sharp is the Clear Winner

| Criterion | sharp ✅ | node-canvas ❌ |
|-----------|----------|----------------|
| **Performance** | 4-5x faster (libvips-based) | Slower (Cairo-based) |
| **Dependencies** | Pre-built binaries, zero system deps | Requires Cairo, Pango, libpng, libjpeg - NOT available in GitHub Actions by default |
| **Installation** | `npm install sharp` works everywhere | Requires `apt-get install` or build from source |
| **GitHub Actions compatibility** | Works out-of-box on ubuntu-latest, macos-latest, windows-latest | Requires custom Docker container or system package installation step |
| **Bundle size** | ~10MB (pre-built binary) | ~15MB + system libraries |
| **API complexity** | Simple, chainable API | Canvas 2D API (more verbose for simple operations) |
| **Text rendering** | Supports text overlays via SVG or Pango | Full Canvas 2D text API |
| **Maintenance** | Actively maintained, v0.34.4 (Nov 2025) | Maintained, but v3.x breaking changes |

### Use node-canvas ONLY if:
- You specifically need Canvas 2D API compatibility
- You're running in a custom Docker container where you control system packages
- You need complex text rendering with multiple fonts and bidirectional text

### For ContribEmblem, use sharp because:
1. **Zero system dependencies** - works immediately in GitHub Actions
2. **Performance** - generating 100+ emblems in a workflow will be 4-5x faster
3. **Simpler API** - compositing emblem background + stats overlay is straightforward
4. **Text rendering is sufficient** - sharp supports SVG text overlays or Pango markup for stat numbers

## Alternatives Considered

| Recommended | Alternative | When to Use Alternative |
|-------------|-------------|-------------------------|
| sharp | Jimp | If you need pure JavaScript (no native binaries), but 10x slower |
| sharp | node-canvas | If you need Canvas 2D API and can manage system dependencies |
| @vercel/ncc | Rollup | If you need more control over bundling (e.g., code splitting) |
| axios | node-fetch | If you want native fetch API, but axios has better error handling and interceptors |
| octokit | @octokit/rest only | If you want to minimize bundle size (~200KB savings) |

## What NOT to Use

| Avoid | Why | Use Instead |
|-------|-----|-------------|
| **ImageMagick/GraphicsMagick** | Requires system binaries not available in GitHub Actions. Slower than sharp. | sharp |
| **node-canvas** | Requires Cairo system libraries. See decision matrix above. | sharp |
| **canvas npm package** | Deprecated, unmaintained. | sharp or @napi-rs/canvas |
| **puppeteer for image gen** | 300MB+ Chrome dependency. Massive overkill for image generation. | sharp |
| **gm (GraphicsMagick wrapper)** | Requires GraphicsMagick binary installation. | sharp |
| **bungie-net-core** | Abandoned, last update 2019. | Custom axios client + bungie-api-ts types |

## Stack Patterns by Use Case

**For maximum performance and simplicity (recommended for ContribEmblem):**
- Use `sharp` for all image operations
- Use `@actions/github` for GitHub API calls (pre-authenticated)
- Use `js-yaml` for config parsing
- Bundle with `@vercel/ncc`

**If you need Canvas 2D API specifically:**
- Use `@napi-rs/canvas` (Rust-based, faster than node-canvas)
- OR use `node-canvas` in a custom Docker container with Cairo pre-installed
- Accept the performance and complexity trade-off

**For Bungie API:**
- Use `axios` with custom client (best control, smallest bundle)
- Add `bungie-api-ts` for TypeScript types (optional, adds 2MB+)
- Implement token refresh logic with axios interceptors

## Version Compatibility

### Critical Dependencies

| Package | Version | Compatible With | Notes |
|---------|---------|-----------------|-------|
| Node.js | 18.17.0+ or 20.3.0+ | GitHub Actions node20 runner | node16 runner is deprecated as of 2025 |
| sharp | 0.34.4 | Node.js 20.9.0+ | Dropped Node 18 support in v0.35.0-rc |
| @actions/core | 1.11.1 | Node.js 18+ | Latest stable |
| @actions/github | 6.0.0 | Node.js 18+ | Uses octokit v5.x |
| octokit | 5.0.5 | Node.js 18+ | All-in-one SDK (includes rest, graphql, auth) |

### Breaking Changes to Watch

- **sharp v0.35.0** (currently in RC): Drops Node.js 18, requires 20.9.0+. Wait for stable release.
- **node-canvas v3.0.0**: Breaking changes to N-API migration. Use v3.2.1 (latest stable).
- **@actions/github v6.0.0**: Changed from octokit v3 to v5. Most code compatible.

## GitHub Actions Runtime Considerations

### What's Available in ubuntu-latest Runner (2025)

✅ **Available out-of-box:**
- Node.js 20.x
- npm, yarn, pnpm
- Git
- Python 3.x (for node-gyp if needed)
- Basic build tools (gcc, make)

❌ **NOT Available (requires manual install):**
- Cairo (required for node-canvas)
- Pango (required for node-canvas)
- libvips (bundled with sharp pre-built binary)

### Optimization Tips

1. **Use Pre-built Binaries**: sharp provides pre-built binaries for Linux x64 (glibc), macOS x64/arm64, Windows x64. No compilation needed.

2. **Bundle with @vercel/ncc**: 
   ```bash
   ncc build src/index.ts -o dist
   ```
   This creates a single `dist/index.js` with all dependencies. Users don't need to run `npm install`.

3. **Cache node_modules**: Use actions/cache for faster workflow runs during development.

4. **Minimize API Calls**: Use GraphQL for complex queries to reduce API call count and avoid rate limits.

5. **Parallel Image Generation**: If generating multiple emblems, use `Promise.all()` or worker threads.

## Sources

**HIGH confidence sources:**
- GitHub Actions Documentation: https://docs.github.com/en/actions/creating-actions/creating-a-javascript-action (Official, verified Jan 2026)
- GitHub Actions Toolkit: https://github.com/actions/toolkit (Official, active development)
- sharp Documentation: https://sharp.pixelplumbing.com/ (Official, verified v0.34.4 release Nov 2025)
- sharp GitHub Releases: https://github.com/lovell/sharp/releases (Verified latest versions)
- node-canvas Releases: https://github.com/Automattic/node-canvas/releases (Verified v3.2.1 Jan 2026)
- Octokit Documentation: https://github.com/octokit/octokit.js (Official, verified v5.0.5 Oct 2025)

**MEDIUM confidence sources:**
- bungie-api-ts: Community-maintained, not official Bungie library
- Performance comparisons: Based on sharp documentation claims and community benchmarks

**LOW confidence (NOT USED):**
- Stack Overflow posts (not relied upon for version recommendations)
- Blog posts older than 2024

---
*Stack research for: GitHub Actions + Image Generation*
*Researched: January 27, 2026*
