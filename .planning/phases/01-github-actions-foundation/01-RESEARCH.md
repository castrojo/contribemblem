# Phase 1: GitHub Actions Foundation - Research

**Researched:** 2026-01-27
**Domain:** GitHub Actions scheduled workflows, permissions, and infinite loop prevention
**Confidence:** HIGH

## Summary

Phase 1 establishes a safe, reliable scheduled GitHub Actions workflow foundation. The research focused on four critical areas: scheduled triggers (cron), preventing infinite workflow loops, configuring write permissions, and setting up git commit metadata.

GitHub Actions uses POSIX cron syntax for scheduled workflows and runs them on the latest commit of the default branch. The minimum interval is 5 minutes, and scheduled workflows always run on the default branch. For preventing infinite loops, GitHub provides multiple mechanisms: GITHUB_TOKEN automatically prevents recursion (commits made with it don't trigger push events), commit message skip instructions ([skip ci], [ci skip], etc.), and path filters. Permission configuration is critical because since 2023, the default GITHUB_TOKEN permissions are read-only. To commit files, workflows must explicitly declare `permissions: contents: write` at the workflow or job level. Git configuration requires setting user.name and user.email before making commits, typically done via `git config` commands in the workflow.

**Primary recommendation:** Use `on.schedule` with cron syntax, add `[skip ci]` to commit messages, explicitly set `permissions.contents: write`, configure git user before commits, and combine with path filters as defense-in-depth.

## Standard Stack

The established approach for this domain:

### Core
| Library/Tool | Version | Purpose | Why Standard |
|--------------|---------|---------|--------------|
| GitHub Actions | Native | CI/CD workflow engine | Built into GitHub, no external dependencies |
| POSIX cron syntax | Standard | Schedule specification | Universal standard for scheduling |
| GITHUB_TOKEN | Native | Authentication token | Automatically provided, scoped per workflow run |
| git CLI | Pre-installed | Commit operations | Standard on GitHub-hosted runners |

### Supporting
| Tool | Version | Purpose | When to Use |
|------|---------|---------|-------------|
| actions/checkout | v4 | Repository checkout | When workflow needs repo access |
| GitHub CLI (gh) | Pre-installed | GitHub API operations | For advanced GitHub operations |
| git config | Built-in | Configure git metadata | Always before committing |

**Installation:**
No installation needed - all tools are pre-installed on GitHub-hosted runners.

## Architecture Patterns

### Recommended Workflow Structure
```yaml
.github/
├── workflows/
│   └── scheduled-workflow.yml    # Main scheduled workflow
```

### Pattern 1: Scheduled Workflow with Safe Commits
**What:** A workflow that runs on schedule, makes changes, and commits without triggering itself
**When to use:** Any scheduled maintenance or automated update workflow
**Example:**
```yaml
# Source: Official GitHub Docs - Workflow Syntax
name: Weekly Update
on:
  schedule:
    - cron: '0 0 * * 0'  # Every Sunday at midnight UTC
  push:
    paths-ignore:
      - 'generated/**'    # Don't run on changes to generated files

permissions:
  contents: write          # Required for pushing commits

jobs:
  update:
    runs-on: ubuntu-latest
    steps:
      - uses: actions/checkout@v4
      
      - name: Configure git
        run: |
          git config user.name "github-actions[bot]"
          git config user.email "github-actions[bot]@users.noreply.github.com"
      
      - name: Generate files
        run: |
          # Your generation logic
          ./scripts/generate.sh
      
      - name: Commit and push
        run: |
          git add .
          git commit -m "chore: automated update [skip ci]" || exit 0
          git push
```

### Pattern 2: Cron Schedule Syntax
**What:** POSIX cron syntax for scheduling workflows
**When to use:** All scheduled workflows
**Example:**
```yaml
# Source: GitHub Docs - Workflow Syntax
on:
  schedule:
    # ┌───────────── minute (0 - 59)
    # │ ┌───────────── hour (0 - 23)
    # │ │ ┌───────────── day of the month (1 - 31)
    # │ │ │ ┌───────────── month (1 - 12 or JAN-DEC)
    # │ │ │ │ ┌───────────── day of the week (0 - 6 or SUN-SAT)
    - cron: '0 0 * * 0'      # Weekly on Sunday at midnight
    - cron: '0 0 * * 1'      # Weekly on Monday at midnight
    - cron: '*/15 * * * *'   # Every 15 minutes (not recommended)
```

**Key points:**
- Minimum interval: 5 minutes
- Runs on latest commit of default branch
- Uses UTC timezone
- High-traffic workflows may be delayed during peak times

### Pattern 3: Permissions Configuration
**What:** Explicit permission grants for GITHUB_TOKEN
**When to use:** Always specify minimum required permissions
**Example:**
```yaml
# Source: GitHub Docs - GITHUB_TOKEN Authentication
permissions:
  contents: write    # For committing files
  issues: write      # For creating/editing issues
  pull-requests: read  # For reading PR data

# Or at job level:
jobs:
  deploy:
    permissions:
      contents: write
    runs-on: ubuntu-latest
    steps: [...]
```

**Important:** Since 2023, default permissions are read-only. Must explicitly grant write access.

### Pattern 4: Path Filters for Loop Prevention
**What:** Prevent workflow triggers based on file paths
**When to use:** Defense-in-depth alongside [skip ci]
**Example:**
```yaml
# Source: GitHub Docs - Events that Trigger Workflows
on:
  push:
    paths-ignore:
      - 'generated/**'
      - '*.md'
      - 'docs/**'
```

### Anti-Patterns to Avoid
- **Using push events without path filters:** Can trigger on own commits if GITHUB_TOKEN workaround used
- **Omitting [skip ci] in commit messages:** No explicit signal prevents recursion
- **Not setting git user.name/user.email:** Commits fail or have incorrect attribution
- **Using default read-only permissions:** Commits fail with 403 errors
- **Setting only workflow-level permissions:** Job-level overrides may revert to read-only

## Don't Hand-Roll

Problems that look simple but have existing solutions:

| Problem | Don't Build | Use Instead | Why |
|---------|-------------|-------------|-----|
| Commit attribution | Custom author strings | `git config user.name/user.email` with GitHub Actions bot | Official bot account, consistent attribution |
| Loop prevention | Custom tracking files | `[skip ci]` + path filters + GITHUB_TOKEN behavior | Built-in, well-tested, standard |
| Token management | Store PAT as secret | Use GITHUB_TOKEN | Auto-managed, scoped, secure |
| Schedule parsing | Custom cron parser | GitHub's on.schedule | Validated, tested, standard |

**Key insight:** GitHub Actions has solved these problems with robust, well-tested solutions. Custom approaches introduce edge cases and maintenance burden.

## Common Pitfalls

### Pitfall 1: Infinite Action Loop
**What goes wrong:** Workflow triggers itself recursively, exhausting quota
**Why it happens:** 
- Commit message doesn't include [skip ci]
- Using PAT instead of GITHUB_TOKEN for commits
- Path filters not configured
- Workflow triggered on `push` without guards

**How to avoid:**
- Always include `[skip ci]` or similar in automated commit messages
- Use GITHUB_TOKEN for commits (doesn't trigger push events)
- Add `paths-ignore` filters for generated files
- Test with small quotas first

**Warning signs:**
- Multiple workflow runs for same commit
- Workflow runs immediately after previous run
- Quota exhausted alerts
- GitHub Actions usage spikes

### Pitfall 2: Permission Denied (403) on Commit
**What goes wrong:** `git push` fails with "Permission denied" or 403 error
**Why it happens:** Default GITHUB_TOKEN permissions are read-only since 2023
**How to avoid:**
```yaml
permissions:
  contents: write  # Must be explicit
```
**Warning signs:**
- Workflow fails at push step
- Error message mentions permissions or 403
- Checkout succeeds but push fails

### Pitfall 3: Missing Git Configuration
**What goes wrong:** Commits fail or have wrong author information
**Why it happens:** Git requires user.name and user.email to be configured
**How to avoid:**
```yaml
- name: Configure git
  run: |
    git config user.name "github-actions[bot]"
    git config user.email "github-actions[bot]@users.noreply.github.com"
```
**Warning signs:**
- "Please tell me who you are" error
- Commits have wrong author
- Git commands fail with configuration errors

### Pitfall 4: Schedule Syntax Errors
**What goes wrong:** Workflow never runs or runs at wrong times
**Why it happens:** Invalid cron syntax or timezone confusion
**How to avoid:**
- Use https://crontab.guru/ to validate cron expressions
- Remember all times are UTC
- Test with shorter intervals first (e.g., hourly) before weekly
- Check workflow runs page to verify schedule

**Warning signs:**
- Workflow doesn't appear in scheduled runs
- Runs at unexpected times
- Syntax validation errors in workflow file

### Pitfall 5: GITHUB_TOKEN vs PAT Confusion
**What goes wrong:** Using PAT when GITHUB_TOKEN would work, or vice versa
**Why it happens:** Misunderstanding when each is needed
**How to avoid:**
- Use GITHUB_TOKEN for same-repo operations (commits, issues, etc.)
- Use PAT/GitHub App only when triggering other workflows or accessing other repos
- Never store GITHUB_TOKEN as a secret (it's automatic)

**Warning signs:**
- Workflow triggers itself despite using token for commits
- Cannot access other repositories
- Permission errors despite correct permissions key

## Code Examples

Verified patterns from official sources:

### Complete Scheduled Workflow with All Requirements
```yaml
# Source: GitHub Docs - Workflow Syntax, Skip Workflow Runs
name: Weekly Contributors Update

on:
  schedule:
    - cron: '0 2 * * 0'  # Every Sunday at 2 AM UTC
  push:
    paths-ignore:
      - 'CONTRIBUTORS.md'
      - '.github/workflows/**'

permissions:
  contents: write

jobs:
  update-contributors:
    runs-on: ubuntu-latest
    
    steps:
      - name: Checkout repository
        uses: actions/checkout@v4
        
      - name: Configure git
        run: |
          git config user.name "github-actions[bot]"
          git config user.email "github-actions[bot]@users.noreply.github.com"
      
      - name: Generate contributors list
        run: |
          # Your generation logic here
          echo "Generating contributors..."
          # Example: node scripts/generate-contributors.js
      
      - name: Check for changes
        id: changes
        run: |
          if git diff --quiet; then
            echo "has_changes=false" >> $GITHUB_OUTPUT
          else
            echo "has_changes=true" >> $GITHUB_OUTPUT
          fi
      
      - name: Commit and push if changed
        if: steps.changes.outputs.has_changes == 'true'
        run: |
          git add .
          git commit -m "chore: update contributors list [skip ci]"
          git push
```

### Multiple Schedules with Context Access
```yaml
# Source: GitHub Docs - Workflow Syntax
on:
  schedule:
    - cron: '30 5 * * 1,3'    # Mon/Wed at 5:30 AM UTC
    - cron: '30 5,17 * * 2,4' # Tue/Thu at 5:30 AM and 5:30 PM UTC

jobs:
  scheduled_task:
    runs-on: ubuntu-latest
    steps:
      - name: Check which schedule triggered
        run: |
          echo "Triggered by: ${{ github.event.schedule }}"
          # Different logic based on schedule if needed
```

### Skip Patterns
```yaml
# Source: GitHub Docs - Skip Workflow Runs
# In commit message (any of these work):
git commit -m "Update files [skip ci]"
git commit -m "Update files [ci skip]"
git commit -m "Update files [no ci]"
git commit -m "Update files [skip actions]"
git commit -m "Update files [actions skip]"

# Or with trailer (at end of commit body):
git commit -m "Update files

Generated by automated process.

skip-checks: true"
```

## State of the Art

| Old Approach | Current Approach | When Changed | Impact |
|--------------|------------------|--------------|--------|
| Permissive defaults | Read-only by default | 2023 | Must explicitly grant write permissions |
| PAT for same-repo | GITHUB_TOKEN sufficient | Always | Simpler, more secure, prevents loops |
| Manual loop prevention | [skip ci] standard | 2020+ | Built-in, reliable mechanism |
| Implicit permissions | Explicit permissions key | 2023 | Better security, clearer intent |

**Deprecated/outdated:**
- **Relying on default write permissions:** GitHub changed defaults to read-only in 2023
- **Using workflow_run event for scheduling:** Complex, use schedule directly
- **Custom commit message parsing:** Use official [skip ci] syntax
- **Setting GITHUB_TOKEN as repository secret:** Token is automatically provided

## Open Questions

1. **Cron schedule reliability during high traffic**
   - What we know: GitHub may delay scheduled workflows during peak usage
   - What's unclear: Specific delay thresholds and behavior
   - Recommendation: Don't rely on exact timing for scheduled workflows; add tolerance

2. **Path filter ordering with [skip ci]**
   - What we know: Both mechanisms can prevent workflow runs
   - What's unclear: Interaction when both are present
   - Recommendation: Use both for defense-in-depth; either will prevent execution

3. **Git configuration persistence**
   - What we know: Git config must be set in each job
   - What's unclear: Whether defaults.run can set this globally
   - Recommendation: Always set git config at job start, don't rely on defaults

## Sources

### Primary (HIGH confidence)
- GitHub Docs - Workflow Syntax for GitHub Actions: https://docs.github.com/en/actions/using-workflows/workflow-syntax-for-github-actions
- GitHub Docs - Events that Trigger Workflows: https://docs.github.com/en/actions/using-workflows/events-that-trigger-workflows
- GitHub Docs - Skipping Workflow Runs: https://docs.github.com/en/actions/managing-workflow-runs/skipping-workflow-runs
- GitHub Docs - Use GITHUB_TOKEN for Authentication: https://docs.github.com/en/actions/security-guides/automatic-token-authentication
- GitHub Docs - Triggering a Workflow: https://docs.github.com/en/actions/using-workflows/triggering-a-workflow

### Secondary (MEDIUM confidence)
- POSIX cron specification: https://pubs.opengroup.org/onlinepubs/9699919799/utilities/crontab.html
- Git config documentation: Standard git CLI documentation

### Tertiary (LOW confidence)
None - all findings verified with official documentation

## Metadata

**Confidence breakdown:**
- Standard stack: HIGH - All tools are official GitHub/git features
- Architecture: HIGH - Patterns from official GitHub documentation
- Pitfalls: HIGH - Well-documented issues with official solutions

**Research date:** 2026-01-27
**Valid until:** 2026-04-27 (90 days - stable platform features)

**Notes:**
- Phase 1 is foundational - must be correct from day 1
- All patterns are from official GitHub documentation
- No third-party libraries or custom solutions needed
- Critical security consideration: permissions must be minimal and explicit
