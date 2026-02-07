# ContribEmblem

Generate Destiny 2-style emblem badges that showcase your GitHub contribution statistics.

## Overview

ContribEmblem is a GitHub Action that automatically generates and updates a personalized emblem badge featuring:

- **Weekly rotating emblems** from Destiny 2's extensive collection
- **Power Level calculation** based on your GitHub activity metrics
- **Auto-updating images** that refresh weekly via GitHub Actions
- **Stable embed URLs** for use in profiles, READMEs, and websites

## Examples

Here are some example badges showcasing different contribution profiles:

### @castrojo (Power Level: 2568)
![castrojo's ContribEmblem](examples/castrojo.png)

*High activity across all contribution areas - commits, PRs, issues, reviews, and popular repositories.*

### @jeefy (Power Level: 1838)
![jeefy's ContribEmblem](examples/jeefy.png)

*Balanced contributions with strong review activity and consistent commits.*

### @mrbobbytables (Power Level: 1623)
![mrbobbytables' ContribEmblem](examples/mrbobbytables.png)

*Community-focused contributor with emphasis on code reviews and issue engagement.*

> **Note:** These are demonstration badges with sample stats to showcase the visual style. Your actual badge will reflect your real GitHub contribution data.

## Features

### 🎮 Destiny 2 Aesthetic
Your emblem rotates weekly from Destiny 2's emblem collection, giving your profile a fresh look while maintaining the iconic Destiny visual style.

### 📊 GitHub Stats Integration
Your "Power Level" is calculated from 5 key GitHub metrics:
- Commits this year
- Pull requests
- Issues created/commented
- Code reviews
- Stars received on your repositories

### 🔄 Automatic Updates
The emblem updates every Sunday at midnight UTC via GitHub Actions. No manual intervention required.

### 🖼️ Easy Embedding
Use stable GitHub URLs to embed your emblem anywhere:

```markdown
![My ContribEmblem](https://raw.githubusercontent.com/castrojo/contribemblem/main/badge.png)
```

## Setup

### Prerequisites
- A GitHub account
- A Bungie.net API key (free) - [Get one here](https://www.bungie.net/en/Application)

### Installation

1. **Fork or clone this repository**

2. **Add secrets to your repository:**
   - Go to Settings → Secrets and variables → Actions
   - Add `BUNGIE_API_KEY` with your Bungie API key

3. **Configure emblem rotation (optional):**
   - Edit `data/emblem-config.json` to customize your emblem rotation
   - Find emblem hashes at [Destiny 2 API documentation](https://bungie-net.github.io/)

4. **Enable GitHub Actions:**
   - The workflow will run automatically every Sunday at midnight UTC
   - Or manually trigger it from the Actions tab

## Local Development

### Building

```bash
# Build the binary
make build

# Run tests
make test

# Run the full pipeline locally
export GITHUB_TOKEN=your_token
export GITHUB_ACTOR=your_username
export BUNGIE_API_KEY=your_key
./contribemblem run
```

### CLI Commands

The `contribemblem` binary supports these subcommands:

```bash
contribemblem fetch-stats      # Fetch GitHub stats via GraphQL
contribemblem select-emblem    # Select weekly emblem hash
contribemblem fetch-emblem     # Fetch emblem image from Bungie API
contribemblem generate         # Generate badge image
contribemblem run              # Run full pipeline
contribemblem help             # Show help message
```

## Configuration

Edit `data/emblem-config.json` to customize your emblem rotation:

```json
{
  "rotation": [
    "1538938257",
    "2962058744",
    "2962058745"
  ],
  "fallback": "1538938257"
}
```

- `rotation`: Array of Bungie emblem hashes to rotate through weekly
- `fallback`: Emblem to use if rotation is empty or config is missing

## Technical Details

- **Language:** Pure Go (no CGo dependencies)
- **Image Generation:** stdlib + `golang.org/x/image`
- **Badge Size:** 800×400px PNG
- **Font:** Rajdhani Bold (embedded via `go:embed`)
- **Caching:** GitHub Actions cache for manifest and stats (24 hours)

## License

See [LICENSE](LICENSE) for details.

## Inspiration

Inspired by Destiny 2's iconic emblem system and the desire to bring that visual flair to GitHub profiles.

<!-- CONTRIBEMBLEM:START -->
![ContribEmblem](badge.png)

*Last updated: February 7, 2026*
<!-- CONTRIBEMBLEM:END -->
