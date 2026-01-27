# Pitfalls Research

**Domain:** GitHub Actions + Bungie API + Image Generation
**Researched:** Jan 27, 2026
**Confidence:** HIGH

## Critical Pitfalls

### Pitfall 1: Infinite Action Trigger Loops

**What goes wrong:**
Action commits to repo, triggering itself infinitely, exhausting quota or creating spam commits.

**Why it happens:**
Default `push` trigger runs on ALL commits, including those made by the action itself. Developers forget to filter or use special commit strategies.

**How to avoid:**
1. **Use `[skip ci]` in commit messages**: Add `[skip ci]` or `[skip actions]` to action-generated commit messages
2. **Filter trigger paths**: Only trigger on specific file changes that aren't action outputs
   ```yaml
   on:
     push:
       paths-ignore:
         - 'profile/**'  # Don't trigger on action's own output
   ```
3. **Use `if` conditions**: Check if commit author is the action bot
   ```yaml
   jobs:
     update:
       if: github.actor != 'github-actions[bot]'
   ```
4. **Branch-specific triggers**: Only run on specific branches
   ```yaml
   on:
     push:
       branches:
         - main
     schedule:
       - cron: '0 0 * * *'  # Prefer scheduled runs
   ```

**Warning signs:**
- Action runs appearing in tight succession (minutes apart)
- Commit history showing repeated identical messages
- GitHub Actions quota warnings
- Multiple workflow runs queued for same commit

**Phase to address:**
Phase 1 (Git workflow setup) - Must be configured correctly from the start

---

### Pitfall 2: GITHUB_TOKEN Permission Insufficient for Operations

**What goes wrong:**
Action fails with 403/permission errors when trying to commit, push, or access APIs, despite using `GITHUB_TOKEN`.

**Why it happens:**
Default `GITHUB_TOKEN` has restricted permissions. Since 2023, GitHub defaults to read-only for many operations.

**How to avoid:**
1. **Explicitly grant write permissions** in workflow:
   ```yaml
   permissions:
     contents: write  # Required for commits/pushes
     pull-requests: write  # If creating/commenting on PRs
   ```
2. **Use elevated PAT for sensitive operations**: For org-wide or cross-repo operations, use a PAT stored as secret
3. **Minimum permissions principle**: Only grant what's needed
4. **Test permissions early**: Verify in Phase 1 before building complex logic

**Warning signs:**
- `Resource not accessible by integration` errors
- 403 Forbidden responses
- Silent failures with no commit pushed
- "refusing to allow a GitHub App to create or update workflow" errors

**Phase to address:**
Phase 1 (initial workflow setup) - Configure before any git operations

---

### Pitfall 3: Rate Limit Exhaustion (GitHub API)

**What goes wrong:**
Action hits GitHub API rate limits when fetching contribution stats, causing failures or incomplete data.

**Why it happens:**
- `GITHUB_TOKEN` has 1,000 requests/hour/repo limit
- Fetching per-repo stats requires multiple API calls
- Calendar year queries can require pagination for high-activity users

**How to avoid:**
1. **Cache API responses**: Use Actions cache to store results
   ```yaml
   - uses: actions/cache@v3
     with:
       path: .cache/gh-api
       key: gh-api-${{ github.sha }}-${{ hashFiles('**/config.yml') }}
   ```
2. **Batch requests**: Minimize API calls by combining queries
3. **Use GraphQL over REST**: Single GraphQL query can replace multiple REST calls
4. **Monitor rate limits**: Check `X-RateLimit-Remaining` header
5. **Implement exponential backoff**: Retry with delays if approaching limits
6. **Schedule runs strategically**: Avoid running at top of hour when other actions run

**Warning signs:**
- `API rate limit exceeded` errors
- Incomplete data in generated images
- Action taking progressively longer to complete
- 403 responses with rate limit messages

**Phase to address:**
Phase 2 (GitHub API integration) - Implement caching and rate limit handling

---

### Pitfall 4: Bungie API Authentication Failures

**What goes wrong:**
Unable to fetch emblem images due to authentication errors or invalid API keys. Emblems fail to load or return 401/403.

**Why it happens:**
- Bungie API requires API key registration and proper header formatting
- API keys can expire or be revoked
- CORS/origin restrictions if testing locally
- Emblem IDs may change between game seasons

**How to avoid:**
1. **Store API key as GitHub Secret**: Never commit API keys
   ```yaml
   env:
     BUNGIE_API_KEY: ${{ secrets.BUNGIE_API_KEY }}
   ```
2. **Test API key validity**: Verify key works before using in production
3. **Include proper headers**: `X-API-Key` header required for all requests
4. **Handle 401/403 gracefully**: Fail with clear error messages
5. **Document API key setup**: README must explain how to get Bungie API key
6. **Cache emblem images**: Don't re-fetch on every run
7. **Validate emblem ID format**: Check format matches Bungie's schema before API call

**Warning signs:**
- 401 Unauthorized responses
- 403 Forbidden (wrong API key)
- Empty/broken images in output
- `ThisEndpointRequiresApiKey` error responses

**Phase to address:**
Phase 3 (Bungie API integration) - Test authentication before implementing emblem resolution

---

### Pitfall 5: Calendar Year Date Math Errors (January 1 Edge Case)

**What goes wrong:**
On January 1, stats reset to zero or display previous year's data incorrectly. Action breaks due to date boundary logic.

**Why it happens:**
- GitHub contribution calendar resets at year boundary
- Date queries may span year boundaries incorrectly
- Off-by-one errors in date range calculations
- Timezone handling across UTC vs user timezone

**How to avoid:**
1. **Explicitly handle year boundaries**: Check if current date is early January
2. **Use ISO date format**: Avoid locale-dependent date parsing
3. **Test with mocked dates**: Unit tests for Jan 1, Dec 31, Feb 29
4. **Account for timezones**: GitHub uses UTC, user may be in different TZ
5. **Grace period for new year**: Show previous year stats for first week of January
6. **Date range validation**: Ensure start date < end date

**Warning signs:**
- Zero contributions on January 1-7
- Wrong year shown in generated image
- Action crashes on January 1 specifically
- Date range errors in logs

**Phase to address:**
Phase 2 (GitHub stats fetching) - Build in date handling logic from start

---

### Pitfall 6: Text Contrast Issues on Variable Emblem Backgrounds

**What goes wrong:**
Generated text is unreadable on certain emblem backgrounds (white text on white emblem, etc.). No dynamic contrast adjustment.

**Why it happens:**
- Emblems have wildly varying colors and patterns
- Static text color doesn't work for all backgrounds
- No contrast checking between text and background
- Overlay shadows/outlines not implemented

**How to avoid:**
1. **Analyze background brightness**: Calculate average brightness of emblem region
2. **Dynamic text color**: Use white text on dark emblems, black on light
3. **Add text outlines/strokes**: Makes text readable on any background
   ```javascript
   ctx.strokeStyle = 'black';
   ctx.lineWidth = 4;
   ctx.strokeText(text, x, y);
   ctx.fillStyle = 'white';
   ctx.fillText(text, x, y);
   ```
4. **Semi-transparent overlay**: Add dark/light box behind text
5. **Test with diverse emblems**: Get emblems from different games/seasons
6. **Fallback to standard layout**: If emblem fetch fails, use safe background

**Warning signs:**
- User complaints about readability
- Text invisible in screenshots
- High contrast emblems showing white-on-white
- Need to squint to read stats

**Phase to address:**
Phase 4 (Image generation) - Implement contrast logic before text rendering

---

### Pitfall 7: Font Availability in Action Runners

**What goes wrong:**
Generated images use fallback system fonts instead of intended custom fonts, looking unprofessional.

**Why it happens:**
- GitHub Actions runners have minimal font installations
- Custom fonts must be bundled or installed
- Font loading fails silently, falling back to Arial/sans-serif
- Font paths hardcoded to local development paths

**How to avoid:**
1. **Bundle fonts in repository**: Include `.ttf` or `.woff` files
2. **Install fonts in action**: Use apt-get or runner setup
   ```yaml
   - name: Install fonts
     run: |
       sudo apt-get update
       sudo apt-get install -y fonts-noto
   ```
3. **Test font loading**: Verify font registered before rendering
4. **Use web-safe fallbacks**: Specify font stack with fallbacks
5. **Canvas font registration**: Explicitly register fonts with canvas library
6. **Relative font paths**: Use paths relative to repo root

**Warning signs:**
- Fonts look different in action output vs local testing
- Generic sans-serif fonts in generated images
- Font-related warnings in action logs
- Text spacing/sizing incorrect

**Phase to address:**
Phase 4 (Image generation) - Configure font installation in workflow before rendering

---

### Pitfall 8: Emblem Image Fetching Failures

**What goes wrong:**
Bungie API returns emblem data, but image URLs are 404/broken. Action continues with missing emblem.

**Why it happens:**
- Emblem icon paths are relative, need full CDN URL construction
- CDN URLs change between Bungie content updates
- Some emblems removed from game but still in API
- Network timeouts during image download
- Image format changes (PNG vs JPG)

**How to avoid:**
1. **Validate image URLs**: Check HTTP status before downloading
2. **Construct CDN URLs correctly**: Use Bungie's base URL + icon path
   ```javascript
   const iconUrl = `https://www.bungie.net${iconPath}`;
   ```
3. **Retry logic**: Implement retries with exponential backoff
4. **Fallback emblem**: Have default emblem if fetch fails
5. **Cache downloaded emblems**: Don't re-download same emblem
6. **Test with multiple emblem IDs**: Verify different game/season emblems work
7. **Validate image data**: Ensure downloaded file is valid image format

**Warning signs:**
- 404 errors in action logs
- Missing emblems in generated image
- Partial image corruption
- Generic placeholder showing instead of emblem

**Phase to address:**
Phase 3 (Bungie API integration) - Implement image URL construction and validation

---

### Pitfall 9: Git Configuration Missing in Action Environment

**What goes wrong:**
Git commit fails with "user.name and user.email not set" error. Action cannot commit changes.

**Why it happens:**
- GitHub Actions runners don't have git user configured by default
- Git requires author information for commits
- Using raw git commands instead of actions/* helpers

**How to avoid:**
1. **Configure git user before commits**:
   ```yaml
   - name: Configure Git
     run: |
       git config user.name "github-actions[bot]"
       git config user.email "github-actions[bot]@users.noreply.github.com"
   ```
2. **Use git actions**: `actions/checkout@v3` with token for automatic config
3. **Add to workflow template**: Include in starter code/documentation
4. **Verify configuration**: Check git config before commit step

**Warning signs:**
- `fatal: unable to auto-detect email address` errors
- Commit step fails in action
- Works locally but fails in CI
- "Please tell me who you are" git errors

**Phase to address:**
Phase 1 (Git workflow setup) - Configure in initial workflow file

---

### Pitfall 10: YAML Configuration Validation Missing

**What goes wrong:**
Users provide invalid YAML config (wrong emblem ID format, missing fields). Action fails with cryptic errors or silently uses defaults.

**Why it happens:**
- No schema validation on config file
- Emblem IDs are numeric but users enter strings
- Required fields not checked before processing
- Type mismatches not caught early

**How to avoid:**
1. **Define JSON schema**: Create schema for config validation
2. **Validate config in action**: Check schema before processing
   ```javascript
   const Ajv = require('ajv');
   const ajv = new Ajv();
   const valid = ajv.validate(schema, config);
   if (!valid) {
     throw new Error(ajv.errorsText());
   }
   ```
3. **Provide validation script**: Let users test config locally
4. **Clear error messages**: Explain what's wrong and how to fix
5. **Example configs**: Show valid config examples in README
6. **Required field defaults**: Use sensible defaults where possible
7. **Type coercion**: Convert strings to numbers if unambiguous

**Warning signs:**
- Users reporting "doesn't work" without details
- Random failures with certain configs
- Wrong emblem showing (ID mismatch)
- Action succeeds but output is incorrect

**Phase to address:**
Phase 5 (Config validation) - Implement before first release

---

## Technical Debt Patterns

Shortcuts that seem reasonable but create long-term problems.

| Shortcut | Immediate Benefit | Long-term Cost | When Acceptable |
|----------|-------------------|----------------|-----------------|
| Hardcoded emblem ID in code | Skip config complexity | Users can't customize | Never - defeats purpose |
| No API rate limiting | Simpler code, faster dev | Production failures | Never - critical for reliability |
| Single font with no fallback | Easier font management | Broken images if font missing | Never - runners vary |
| Skip commit message filtering | Simpler trigger config | Infinite loop risk | Never - catastrophic failure mode |
| No image caching | Simpler file management | Repeated API calls, slower runs | Never - waste quota |
| Text without outline/shadow | Simpler canvas code | Unreadable on many emblems | Only if emblem backgrounds are controlled |
| REST API instead of GraphQL | More familiar | Higher rate limit usage | MVP only, refactor in Phase 2 |
| No date edge case testing | Faster initial dev | January 1 failures | Never - predictable annual bug |
| Inline secrets in workflow | Faster testing | Security vulnerability | Never - even in private repos |
| Skip error messages for users | Less code to maintain | Users can't debug issues | Never - poor UX |

## Integration Gotchas

Common mistakes when connecting to external services.

| Integration | Common Mistake | Correct Approach |
|-------------|----------------|------------------|
| GitHub API | Using REST for everything | Use GraphQL for complex queries, REST for simple ops |
| GitHub API | Not checking rate limit headers | Check `X-RateLimit-Remaining` before batch operations |
| GitHub API | Assuming pagination not needed | Always paginate for user activity (100+ contributions) |
| Bungie API | Forgetting `X-API-Key` header | Include in every request, even GET |
| Bungie API | Using full icon URL from API | API returns relative path, prepend `https://www.bungie.net` |
| Bungie API | No retry on 503 (maintenance) | Bungie has weekly maintenance, implement retries |
| Git in Actions | Not setting user.name/email | Configure before any git commit operation |
| Git in Actions | Wrong GITHUB_TOKEN permissions | Set `contents: write` in workflow permissions |
| Image Generation | Using file:// paths for fonts | Fonts must be in repo or installed system-wide |
| Image Generation | Not validating image dimensions | Canvas operations fail on 0-width/height |
| Workflow Triggers | `push` without path filters | Triggers on action's own commits (infinite loop) |
| Cron Schedule | Expecting exact timing | Delay up to 10 minutes, don't rely on precise timing |

## Performance Traps

Patterns that work at small scale but fail as usage grows.

| Trap | Symptoms | Prevention | When It Breaks |
|------|----------|------------|----------------|
| Fetching all repos individually | Long action runtime, rate limits | Use GraphQL to batch repo queries | 10+ repositories |
| Re-downloading same emblem every run | Slow, unnecessary API calls | Cache emblems by ID | Every run after first |
| No API result caching | Rate limit exhaustion | Cache with daily expiry | 50+ runs per day |
| Synchronous image processing | Action timeout on large images | Stream processing or resize | Images >5MB |
| Loading entire contribution graph | Memory issues, slow queries | Limit to current year | Users with 5+ years activity |
| No pagination for contributions | Missing data for active users | Always paginate, check hasNextPage | Users with 500+ contributions |
| Re-fetching unchanged data | Wasted quota, slow runs | Use conditional requests (ETags) | Every scheduled run |
| Sequential API calls | Long wait times | Parallelize independent requests | 5+ API endpoints |

## Security Mistakes

Domain-specific security issues beyond general web security.

| Mistake | Risk | Prevention |
|---------|------|------------|
| Committing Bungie API key to repo | Key revocation, quota theft | Store in GitHub Secrets, never commit |
| Using PAT instead of GITHUB_TOKEN | Excessive permissions, security risk | Use GITHUB_TOKEN unless cross-repo access needed |
| Not validating user-provided emblem IDs | API abuse, injection attacks | Validate format (numeric, length) before API call |
| Exposing full Bungie API responses | Leaking user PII (names, IDs) | Filter response to only needed fields |
| Not sanitizing commit messages | Script injection in logs/UI | Escape special characters in generated messages |
| World-readable cache with API data | Data leakage | Use repository-scoped cache only |
| Using `pull_request` trigger from forks | Untrusted code execution | Use `pull_request_target` with caution or disable |
| No rate limiting in public workflows | Quota exhaustion by abuse | Add workflow approval for first-time contributors |

## UX Pitfalls

Common user experience mistakes in this domain.

| Pitfall | User Impact | Better Approach |
|---------|-------------|-----------------|
| No clear setup instructions | Users can't get started | Step-by-step README with screenshots |
| Cryptic error messages | Users stuck, file issues | Actionable errors: "Add BUNGIE_API_KEY to secrets" |
| Action fails silently | Users don't know it broke | Create GitHub issue on failure with logs |
| No preview before commit | Surprises in generated image | Option to generate artifact before commit |
| Stats not updating | Confusion about when action runs | Document schedule, show "last updated" in image |
| Generic emblem when fetch fails | Looks broken | Clear "emblem unavailable" message in image |
| Mobile-unfriendly image size | Hard to view on GitHub mobile | Optimize image dimensions for mobile (800x400px max) |
| No customization options | One-size-fits-all | Config for colors, layout, font size |
| Overwrites existing profile content | Destroys user's work | Append to README or use dedicated file |
| No way to disable temporarily | Must delete workflow file | Add `enabled: false` config option |

## "Looks Done But Isn't" Checklist

Things that appear complete but are missing critical pieces.

- [ ] **Git workflow:** Verified action doesn't trigger itself infinitely (test with 3 consecutive runs)
- [ ] **Permissions:** Confirmed `contents: write` permission in workflow (check action logs for 403s)
- [ ] **Rate limits:** Implemented caching and rate limit checking (verify `X-RateLimit-Remaining` logged)
- [ ] **Bungie auth:** Tested with valid and invalid API keys (verify error handling)
- [ ] **Date boundaries:** Tested on Jan 1, Dec 31, leap year (mock system date)
- [ ] **Text contrast:** Generated images with 5+ different emblems (verify readability)
- [ ] **Font loading:** Verified font renders correctly in action runner (compare local vs CI)
- [ ] **Image URLs:** Tested emblem fetch with multiple IDs including old/removed ones (check 404 handling)
- [ ] **Config validation:** Attempted invalid YAML, missing fields, wrong types (verify error messages)
- [ ] **Error messages:** Simulated all failure modes, checked user sees actionable message
- [ ] **Calendar year:** Tested contribution query with users having 0, 1, and 1000+ contributions
- [ ] **Timezone handling:** Tested with users in different timezones (UTC vs local time)

## Recovery Strategies

When pitfalls occur despite prevention, how to recover.

| Pitfall | Recovery Cost | Recovery Steps |
|---------|---------------|----------------|
| Infinite loop triggered | LOW | Cancel workflow runs, add `[skip ci]` to commit message, push fix |
| Rate limit exhausted | MEDIUM | Wait for reset (1 hour), implement caching in next version |
| Invalid API key | LOW | Rotate key in secrets, re-run workflow |
| Wrong emblem showing | LOW | Validate emblem ID format, update config, re-run |
| January 1 date bug | MEDIUM | Hotfix date logic, backfill missing data from cache |
| Unreadable text on emblem | MEDIUM | Add text outline/shadow, regenerate images |
| Font missing in CI | LOW | Install font in workflow, re-run |
| Emblem image 404 | LOW | Implement fallback emblem, document how to find valid IDs |
| Git config missing | LOW | Add git config step to workflow, re-run |
| YAML validation failure | LOW | Fix config per error message, commit, re-run |
| Quota exhausted (Bungie) | HIGH | Wait 24 hours, implement request caching for future |
| Action timeout | MEDIUM | Optimize API calls, reduce image processing, increase timeout |

## Pitfall-to-Phase Mapping

How roadmap phases should address these pitfalls.

| Pitfall | Prevention Phase | Verification |
|---------|------------------|--------------|
| Infinite action loops | Phase 1: Git workflow | Trigger 3 consecutive runs, verify stops |
| GITHUB_TOKEN permissions | Phase 1: Git workflow | Check action logs for successful commits |
| GitHub API rate limits | Phase 2: GitHub API | Monitor rate limit headers in logs |
| Bungie API auth failures | Phase 3: Bungie API | Test with invalid key, verify error |
| Calendar year date math | Phase 2: GitHub API | Unit tests for Jan 1, Dec 31, leap year |
| Text contrast issues | Phase 4: Image gen | Test with 5+ diverse emblems |
| Font availability | Phase 4: Image gen | Compare local vs CI output |
| Emblem fetch failures | Phase 3: Bungie API | Test with invalid/old emblem IDs |
| Git config missing | Phase 1: Git workflow | Verify commits have correct author |
| YAML validation | Phase 5: Config | Submit invalid config, check error |
| Timezone edge cases | Phase 2: GitHub API | Test with UTC vs local midnight |
| Image caching | Phase 3: Bungie API | Verify cache hit on second run |

## Sources

- GitHub Actions Security Documentation: https://docs.github.com/en/actions/security-guides/security-hardening-for-github-actions (HIGH confidence)
- GitHub Actions Events Documentation: https://docs.github.com/en/actions/reference/workflows-and-actions/events-that-trigger-workflows (HIGH confidence)
- Bungie.net API Documentation: https://bungie-net.github.io/multi/index.html (HIGH confidence)
- GitHub Actions rate limiting: Official docs (1000 req/hour for GITHUB_TOKEN) (HIGH confidence)
- Canvas text rendering best practices: Industry standard practices (MEDIUM confidence)
- Git commit loop prevention: Common GitHub Actions pattern (HIGH confidence)
- Date boundary edge cases: Standard software testing practices (HIGH confidence)
- Font installation in CI: GitHub Actions runner documentation (HIGH confidence)

---
*Pitfalls research for: GitHub Actions + Bungie API + Image Generation*
*Researched: Jan 27, 2026*
