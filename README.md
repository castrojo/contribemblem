# ContribEmblem

Generate Destiny 2-style emblem badges that showcase your GitHub contribution statistics.

## Overview

ContribEmblem is a GitHub Action that automatically generates and updates a personalized emblem badge featuring:

- **Weekly rotating emblems** from Destiny 2's extensive collection
- **Power Level calculation** based on your GitHub activity metrics
- **Auto-updating images** that refresh weekly via GitHub Actions
- **Stable embed URLs** for use in profiles, READMEs, and websites

## Features

### 🎮 Destiny 2 Aesthetic
Your emblem rotates weekly from Destiny 2's emblem collection, giving your profile a fresh look while maintaining the iconic Destiny visual style.

### 📊 GitHub Stats Integration
Your "Power Level" is calculated from 5 key GitHub metrics:
- Contributions this year
- Commit count
- Pull requests
- Issues created/commented
- Repositories contributed to

### 🔄 Automatic Updates
The emblem updates every Sunday at midnight UTC via GitHub Actions. No manual intervention required.

### 🖼️ Easy Embedding
Use stable GitHub URLs to embed your emblem anywhere:

```markdown
![My ContribEmblem](https://raw.githubusercontent.com/castrojo/contribemblem/main/emblem.png)
```

## Status

**🚧 Under Development**

This project is currently being built using the GSD (Get Shit Done) workflow:

- ✅ Phase 0: Project Planning & Research
- 🔄 Phase 1: GitHub Actions Foundation (In Progress)
- ⏳ Phase 2: GitHub Stats Collection
- ⏳ Phase 3: Bungie API Integration
- ⏳ Phase 4: Image Generation with Power Level
- ⏳ Phase 5: README Update & Commit
- ⏳ Phase 6: Configuration & Validation

## Roadmap

### Phase 1: GitHub Actions Foundation
Set up the scheduled workflow infrastructure with proper permissions and loop prevention.

### Phase 2: GitHub Stats Collection
Implement GitHub GraphQL API integration to fetch current year statistics.

### Phase 3: Bungie API Integration
Connect to Bungie's API to fetch emblem images and metadata for weekly rotation.

### Phase 4: Image Generation
Use the `sharp` library to composite emblem images with Power Level overlays.

### Phase 5: README Update & Commit
Automatically inject emblem images into README with marker-based updates.

### Phase 6: Configuration & Validation
Add YAML configuration for customization and comprehensive error handling.

## License

See [LICENSE](LICENSE) for details.

## Inspiration

Inspired by Destiny 2's iconic emblem system and the desire to bring that visual flair to GitHub profiles.
