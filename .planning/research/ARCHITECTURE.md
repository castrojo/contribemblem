# Architecture Research

**Domain:** GitHub Actions with Image Generation and Auto-Commit
**Researched:** January 27, 2026
**Confidence:** HIGH

## Standard Architecture

### System Overview

```
┌─────────────────────────────────────────────────────────────────┐
│                     GitHub Action Runner                         │
├─────────────────────────────────────────────────────────────────┤
│  ┌────────────┐  ┌────────────┐  ┌────────────┐  ┌────────────┐│
│  │   Config   │  │  GitHub    │  │   Bungie   │  │   Image    ││
│  │   Parser   │  │   API      │  │   API      │  │  Renderer  ││
│  │            │  │  Client    │  │  Client    │  │            ││
│  └─────┬──────┘  └─────┬──────┘  └─────┬──────┘  └─────┬──────┘│
│        │               │               │               │        │
│        └───────────────┴───────────────┴───────────────┘        │
│                            │                                    │
├────────────────────────────┼────────────────────────────────────┤
│                     Data Aggregation Layer                      │
│                            │                                    │
│  ┌─────────────────────────┴────────────────────────────┐      │
│  │              Image Generation Engine                  │      │
│  │    (Composite background + stats + emblem icons)      │      │
│  └─────────────────────┬────────────────────────────────┘      │
├────────────────────────┼────────────────────────────────────────┤
│                     Output Layer                                │
│  ┌──────────────┐  ┌──────────────┐  ┌──────────────┐          │
│  │   PNG File   │  │  README.md   │  │     Git      │          │
│  │   Writer     │  │   Updater    │  │  Committer   │          │
│  └──────────────┘  └──────────────┘  └──────────────┘          │
└─────────────────────────────────────────────────────────────────┘
```

### Component Responsibilities

| Component | Responsibility | Typical Implementation |
|-----------|----------------|------------------------|
| Config Parser | Reads and validates YAML config from repo | YAML parser (PyYAML/js-yaml), validates structure and rotation list |
| GitHub API Client | Fetches GitHub contribution stats using Octokit | REST client with GITHUB_TOKEN, scoped to current calendar year |
| Bungie API Client | Fetches emblem artwork from Bungie Destiny API | REST client with Bungie API key, handles rate limiting |
| Image Renderer | Generates composite PNG with background, stats, icons | Canvas library (node-canvas/Pillow), composites layers |
| PNG Writer | Writes generated image to filesystem | File I/O, typically to project root or assets/ folder |
| README Updater | Updates README.md with image reference | String manipulation with marker comments or prepend strategy |
| Git Committer | Commits and pushes changes back to repo | git CLI or libgit2, uses GITHUB_TOKEN for auth |

## Recommended Project Structure

### For JavaScript/TypeScript Action

```
src/
├── config/             # Configuration handling
│   ├── parser.js       # Reads and validates YAML config
│   └── validator.js    # Schema validation logic
├── api/                # External API clients
│   ├── github.js       # GitHub API (Octokit wrapper)
│   └── bungie.js       # Bungie API client
├── image/              # Image generation
│   ├── renderer.js     # Main image composition logic
│   ├── background.js   # Background layer handling
│   ├── stats.js        # Stats overlay rendering
│   └── emblems.js      # Emblem icon placement
├── git/                # Git operations
│   ├── committer.js    # Commit and push logic
│   └── readme.js       # README.md update logic
├── utils/              # Shared utilities
│   ├── logger.js       # Logging utility
│   └── date.js         # Date/calendar year helpers
└── index.js            # Entry point (reads inputs, orchestrates flow)

action.yml              # GitHub Action metadata
package.json            # Dependencies
README.md               # Documentation
```

### For Python Action

```
src/
├── config/             # Configuration handling
│   ├── parser.py       # Reads and validates YAML config
│   └── validator.py    # Schema validation logic
├── api/                # External API clients
│   ├── github.py       # GitHub API (PyGithub wrapper)
│   └── bungie.py       # Bungie API client
├── image/              # Image generation
│   ├── renderer.py     # Main image composition (Pillow)
│   ├── background.py   # Background layer handling
│   ├── stats.py        # Stats overlay rendering
│   └── emblems.py      # Emblem icon placement
├── git/                # Git operations
│   ├── committer.py    # Commit and push logic
│   └── readme.py       # README.md update logic
├── utils/              # Shared utilities
│   ├── logger.py       # Logging utility
│   └── date.py         # Date/calendar year helpers
└── main.py             # Entry point (reads inputs, orchestrates flow)

action.yml              # GitHub Action metadata
requirements.txt        # Python dependencies
Dockerfile              # Container definition (if using Docker)
README.md               # Documentation
```

### Structure Rationale

- **config/**: Isolates configuration parsing from business logic, makes testing easier
- **api/**: Separate clients for each external service, allows mocking during tests
- **image/**: Image generation is complex enough to warrant its own module with sub-components
- **git/**: Git operations are sensitive (commits, pushes) and should be isolated
- **utils/**: Shared utilities prevent duplication and centralize common functionality
- **Entry point at root**: Single file orchestrates the flow, making the action's execution path clear

## Architectural Patterns

### Pattern 1: Pipeline Architecture (Recommended)

**What:** Linear data flow where each stage processes and passes data to the next stage

**When to use:** When operations have clear dependencies (can't render image until data is fetched)

**Trade-offs:**
- ✅ Simple to understand and debug
- ✅ Easy to add logging between stages
- ✅ Natural error propagation
- ❌ No parallelization of independent tasks
- ❌ Failure in one stage blocks entire pipeline

**Example:**
```javascript
async function main() {
  // Stage 1: Configuration
  const config = await parseConfig();
  
  // Stage 2: Data Fetching (sequential or parallel)
  const [githubStats, emblemData] = await Promise.all([
    fetchGitHubStats(config.username),
    fetchBungieEmblems(config.rotation)
  ]);
  
  // Stage 3: Selection
  const selectedEmblem = selectRandomEmblem(emblemData);
  
  // Stage 4: Image Generation
  const imagePath = await renderImage(githubStats, selectedEmblem);
  
  // Stage 5: File Operations
  await updateReadme(imagePath);
  await commitChanges([imagePath, 'README.md']);
}
```

### Pattern 2: Action Input/Output Pattern

**What:** GitHub Actions communicate via inputs (action.yml) and outputs (set-output commands)

**When to use:** Always - this is the GitHub Actions standard

**Trade-offs:**
- ✅ Standard GitHub Actions pattern
- ✅ Clear contract between action and workflow
- ✅ Easy to test with different inputs
- ❌ Limited to string inputs/outputs

**Example:**
```yaml
# action.yml
name: 'ContribEmblem'
inputs:
  github_token:
    description: 'GitHub token for API access'
    required: true
  bungie_api_key:
    description: 'Bungie API key'
    required: true
  config_path:
    description: 'Path to config YAML'
    required: false
    default: '.contribemblem.yml'
outputs:
  image_path:
    description: 'Path to generated image'
runs:
  using: 'node20'
  main: 'dist/index.js'
```

### Pattern 3: Marker-Based README Update

**What:** Use HTML comments as markers to identify where to inject content in README

**When to use:** When you want to preserve user's README content around the image

**Trade-offs:**
- ✅ Non-destructive to existing README content
- ✅ User controls placement with markers
- ✅ Common pattern in GitHub Actions ecosystem
- ❌ Requires user to add markers manually
- ❌ Can fail silently if markers are missing

**Example:**
```markdown
# My Profile

Some custom content here...

<!-- CONTRIBEMBLEM:START -->
![Stats](./contribemblem-stats.png)
Last updated: 2026-01-27
<!-- CONTRIBEMBLEM:END -->

More custom content...
```

```javascript
function updateReadme(imagePath) {
  const readme = fs.readFileSync('README.md', 'utf8');
  const startMarker = '<!-- CONTRIBEMBLEM:START -->';
  const endMarker = '<!-- CONTRIBEMBLEM:END -->';
  
  const newContent = `![Stats](${imagePath})\nLast updated: ${new Date().toISOString()}`;
  
  const updated = readme.replace(
    new RegExp(`${startMarker}[\\s\\S]*${endMarker}`),
    `${startMarker}\n${newContent}\n${endMarker}`
  );
  
  fs.writeFileSync('README.md', updated);
}
```

### Pattern 4: Temporary File Management

**What:** Use runner's temp directory for intermediate files, only commit final outputs

**When to use:** When generating images that require multiple processing steps

**Trade-offs:**
- ✅ Keeps repo clean
- ✅ No need to .gitignore temporary files
- ✅ Automatic cleanup by runner
- ❌ Slightly more complex path management

**Example:**
```javascript
const os = require('os');
const path = require('path');

const tempDir = os.tmpdir();
const backgroundPath = path.join(tempDir, 'background.png');
const overlayPath = path.join(tempDir, 'overlay.png');
const finalPath = path.join(process.cwd(), 'contribemblem-stats.png');

// Process in temp, output to repo
await downloadBackground(backgroundPath);
await renderOverlay(overlayPath);
await composite([backgroundPath, overlayPath], finalPath);
```

## Data Flow

### Request Flow

```
[Scheduled Trigger / Manual Dispatch]
    ↓
[Action Runner Starts] → Read action.yml inputs
    ↓
[Parse Config YAML] → Validate structure, read rotation list
    ↓
[Fetch GitHub Stats] → Octokit API, current year contributions
    ↓
[Select Emblem] → Random selection from rotation list
    ↓
[Fetch Emblem Data] → Bungie API, emblem artwork
    ↓
[Render Image] → Canvas/Pillow composite
    ├─ Load background
    ├─ Render stats overlay
    └─ Place emblem icons
    ↓
[Write PNG] → Save to repo directory
    ↓
[Update README] → Insert/update image reference
    ↓
[Git Commit] → Stage files, commit, push
    ↓
[Action Completes]
```

### State Management

GitHub Actions are stateless by default. Each run starts fresh:

1. **Input State**: Passed via `action.yml` inputs and secrets
2. **Persistent State**: Stored in repo (config YAML, previous images)
3. **Ephemeral State**: In-memory during action execution, lost after run

### Key Data Flows

1. **Config → Validation → Execution**: Config YAML is parsed and validated before any API calls
2. **API Data → Aggregation → Rendering**: Stats and emblem data are fetched separately, then combined for rendering
3. **Image → File System → Git**: Generated image is written to disk, then committed via Git

## Scaling Considerations

| Scale | Architecture Adjustments |
|-------|--------------------------|
| Single user (personal profile) | Simple pipeline, direct API calls, no caching needed |
| Organization use (10-100 repos) | Consider caching emblem artwork, reusable composite action |
| High frequency runs (hourly) | Cache API responses when possible, use conditional execution (only run if stats changed) |

### Scaling Priorities

1. **First bottleneck: API Rate Limits**
   - GitHub API: 5000 req/hr authenticated, 60 req/hr unauthenticated
   - Bungie API: Varies by endpoint, typically generous for emblem lookups
   - **Solution**: Cache emblem artwork locally, only fetch when rotation changes

2. **Second bottleneck: Image Generation Performance**
   - Rendering can take 1-3 seconds depending on complexity
   - **Solution**: Use pre-rendered backgrounds, optimize overlay operations
   - Not critical for scheduled runs (weekly is typical), but matters for manual dispatch

## Anti-Patterns

### Anti-Pattern 1: Committing on Every Run Even Without Changes

**What people do:** Action commits even if stats/emblem haven't changed

**Why it's wrong:** Creates unnecessary commit noise, uses GitHub Actions minutes unnecessarily

**Do this instead:**
```javascript
// Compare new image hash to previous
const newImageHash = hashFile(newImagePath);
const oldImageHash = fs.existsSync(oldImagePath) ? hashFile(oldImagePath) : null;

if (newImageHash !== oldImageHash) {
  await commitChanges();
} else {
  console.log('No changes detected, skipping commit');
}
```

### Anti-Pattern 2: Hardcoding Paths Instead of Using GitHub Context

**What people do:** Hardcode repo name, paths, or assume directory structure

**Why it's wrong:** Breaks when repo is renamed or forked

**Do this instead:**
```javascript
// Use GitHub Action context
const repoName = process.env.GITHUB_REPOSITORY; // owner/repo
const workspace = process.env.GITHUB_WORKSPACE; // /home/runner/work/repo/repo
const configPath = path.join(workspace, '.contribemblem.yml');
```

### Anti-Pattern 3: Using Commit-to-Branch Without Protection Checks

**What people do:** Action commits directly to main branch without checking branch protection rules

**Why it's wrong:** Can fail on protected branches, or bypass required reviews

**Do this instead:**
```javascript
// Check if branch allows direct commits
const branch = process.env.GITHUB_REF_NAME;
if (branch === 'main' || branch === 'master') {
  // Verify we have write permissions
  // Consider using a different strategy (PR vs direct commit)
}

// Or use bot account with bypass permissions
// Or use workflow_dispatch with manual approval
```

### Anti-Pattern 4: Storing Secrets in Code or Config

**What people do:** Put API keys in YAML config or code

**Why it's wrong:** Security vulnerability, keys get committed to history

**Do this instead:**
```yaml
# In workflow file (.github/workflows/update-stats.yml)
- uses: username/contribemblem@v1
  with:
    github_token: ${{ secrets.GITHUB_TOKEN }}    # From secrets
    bungie_api_key: ${{ secrets.BUNGIE_API_KEY }} # From secrets
    config_path: '.contribemblem.yml'             # Non-sensitive config
```

### Anti-Pattern 5: Not Handling API Failures Gracefully

**What people do:** Action fails completely if GitHub or Bungie API is down

**Why it's wrong:** Breaks the entire workflow for transient issues

**Do this instead:**
```javascript
async function fetchWithRetry(fn, retries = 3) {
  for (let i = 0; i < retries; i++) {
    try {
      return await fn();
    } catch (error) {
      if (i === retries - 1) throw error;
      console.log(`Retry ${i + 1}/${retries} after error:`, error.message);
      await sleep(1000 * (i + 1)); // Exponential backoff
    }
  }
}

// Use fallback data if API fails
let githubStats;
try {
  githubStats = await fetchWithRetry(() => getGitHubStats());
} catch (error) {
  console.warn('Using cached stats due to API failure');
  githubStats = loadCachedStats();
}
```

## Integration Points

### External Services

| Service | Integration Pattern | Notes |
|---------|---------------------|-------|
| GitHub API | REST API via Octokit/PyGithub | Requires GITHUB_TOKEN with `repo` and `user` scope |
| Bungie API | REST API with API key | Rate limits vary by endpoint, generous for emblem lookups |
| WakaTime (optional) | REST API with API key | If tracking coding time stats |
| shields.io (for badges) | URL-based badge generation | Can generate dynamic badges for stats |

### Internal Boundaries

| Boundary | Communication | Notes |
|----------|---------------|-------|
| Config Parser ↔ Validator | Direct function call | Sync validation, throws on invalid config |
| API Clients ↔ Renderer | JSON data objects | API clients return normalized data structures |
| Renderer ↔ File Writer | File path string | Renderer produces file, writer handles I/O |
| File Writer ↔ Git Committer | File path array | Committer stages multiple files |
| All ↔ Logger | Event/method calls | Centralized logging for debugging |

## GitHub Actions Specific Considerations

### action.yml Structure

**Composite Action** (recommended for flexibility):
```yaml
name: 'ContribEmblem'
description: 'Generate and commit GitHub contribution stats with Destiny emblems'
author: 'YourUsername'
branding:
  icon: 'award'
  color: 'purple'

inputs:
  github_token:
    description: 'GitHub token with repo and user scope'
    required: true
  bungie_api_key:
    description: 'Bungie API key for emblem artwork'
    required: true
  config_path:
    description: 'Path to config YAML'
    required: false
    default: '.contribemblem.yml'

runs:
  using: 'node20'
  main: 'dist/index.js'
```

**Docker Action** (for Python or complex dependencies):
```yaml
name: 'ContribEmblem'
description: 'Generate and commit GitHub contribution stats with Destiny emblems'
runs:
  using: 'docker'
  image: 'Dockerfile'
```

### Permissions Required

```yaml
# In workflow file (.github/workflows/update-stats.yml)
permissions:
  contents: write  # For committing changes
  pull-requests: write  # If creating PRs instead of direct commits
```

### Workflow Triggers

```yaml
on:
  schedule:
    - cron: '0 0 * * 1'  # Weekly on Monday at midnight UTC
  workflow_dispatch:  # Manual trigger via GitHub UI
```

### Commit Strategy Options

1. **Direct Commit (Simplest)**
   - Use git commands directly in action
   - Requires `contents: write` permission
   - Works with GITHUB_TOKEN
   
2. **git-auto-commit Action (Recommended)**
   - Use proven action: `stefanzweifel/git-auto-commit-action@v5`
   - Handles edge cases (empty commits, branch protection)
   - More reliable than custom git commands

3. **Pull Request Strategy**
   - Create PR instead of direct commit
   - Better for protected branches
   - Use `peter-evans/create-pull-request@v6`

### File Paths in Actions

```javascript
// ✅ Correct: Use GITHUB_WORKSPACE
const workspace = process.env.GITHUB_WORKSPACE || process.cwd();
const configPath = path.join(workspace, '.contribemblem.yml');

// ❌ Wrong: Assume relative paths
const configPath = './.contribemblem.yml'; // May not work in action runner
```

### Debugging Actions

```yaml
# Enable debug logging
- uses: actions/checkout@v4
- uses: username/contribemblem@v1
  with:
    github_token: ${{ secrets.GITHUB_TOKEN }}
  env:
    ACTIONS_STEP_DEBUG: true  # Enables debug output
```

## Build Order and Dependencies

### Component Build Order

1. **Config Parser & Validator** (no dependencies)
   - Can be built and tested independently
   - Required by all other components

2. **API Clients** (depend on config)
   - GitHub client: Uses token from config
   - Bungie client: Uses API key from config
   - Can be built in parallel

3. **Image Renderer** (depends on API clients)
   - Requires data structures from API clients
   - Can develop with mock data initially

4. **File Operations** (depend on renderer)
   - PNG Writer: Waits for renderer output
   - README Updater: Can be developed independently with mock image path

5. **Git Committer** (depends on file operations)
   - Needs files to be written before committing
   - Final component in the chain

### Suggested Implementation Order

**Phase 1: Core Infrastructure**
1. Set up action.yml with inputs/outputs
2. Create config parser with validation
3. Set up logging and error handling utilities

**Phase 2: Data Fetching**
4. Implement GitHub API client (can use mock data for testing)
5. Implement Bungie API client with emblem rotation selection
6. Test data fetching with real APIs

**Phase 3: Image Generation**
7. Implement background rendering
8. Implement stats overlay
9. Implement emblem icon placement
10. Test image generation with mock data

**Phase 4: File Operations**
11. Implement PNG file writer
12. Implement README updater with marker pattern
13. Test file operations in temp directory

**Phase 5: Git Integration**
14. Implement git committer
15. Test full flow in a test repository
16. Add conditional commit logic (only commit if changed)

**Phase 6: Polish**
17. Add comprehensive error handling
18. Improve logging and debugging output
19. Write documentation and examples
20. Set up CI/CD for the action itself

### Testing Strategy by Component

- **Config Parser**: Unit tests with various YAML inputs
- **API Clients**: Integration tests with real APIs (use rate limit carefully) + mocks for unit tests
- **Image Renderer**: Visual regression tests (compare generated images)
- **File Operations**: Unit tests with temp directories
- **Git Committer**: Integration tests in test repo (careful not to spam commits)

## Sources

- GitHub Actions Documentation - Custom Actions: https://docs.github.com/en/actions/concepts/workflows-and-actions/custom-actions (HIGH confidence)
- GitHub Actions Metadata Syntax: https://docs.github.com/en/actions/reference/workflows-and-actions/metadata-syntax (HIGH confidence)
- waka-readme-stats Action (Reference Implementation): https://github.com/anmol098/waka-readme-stats (MEDIUM confidence - community project, 3.9k stars)
- GitHub Actions Marketplace - README Generator Actions: https://github.com/marketplace?type=actions&query=readme+generator (MEDIUM confidence - ecosystem survey)

---
*Architecture research for: GitHub Actions with Image Generation and Auto-Commit*
*Researched: January 27, 2026*
