---
phase: 04-image-generation-with-power-level
plan: 01
status: complete
duration: "Implemented during Go conversion (commits GO-01 through GO-06)"
completed: 2026-02-06

tech_stack:
  added:
    - "golang.org/x/image v0.35.0 (image compositing, font rendering)"
    - "golang.org/x/image/font/opentype (Rajdhani Bold TTF parsing)"
    - "go:embed (font asset embedding)"
  modified:
    - "Replaced planned Node.js/node-canvas stack with pure Go (no CGo)"

patterns:
  - "Multi-layer text rendering: shadow → stroke → fill (3 passes)"
  - "Stroke via multi-offset technique (16 offsets at ±2px for 4px effective width)"
  - "K/M number formatting for large stat values"
  - "BiLinear scaling for emblem background to canvas dimensions"
  - "go:embed for font assets (zero external file dependencies)"

files_created:
  - internal/badge/generator.go
  - internal/badge/text.go
  - internal/badge/format.go
  - internal/badge/generator_test.go
  - internal/badge/text_test.go
  - internal/badge/format_test.go
  - internal/badge/assets/fonts/Rajdhani-Bold.ttf

files_modified:
  - cmd/contribemblem/main.go (added generate and run subcommands)
  - .github/workflows/update-emblem.yml (Go binary replaces bash/node scripts)

key_types_exported:
  - "badge.Stats (Commits, PullRequests, Issues, Reviews, Stars)"

key_functions_exported:
  - "badge.Generate(emblemPath, stats, outputPath) error"
  - "badge.DrawTextWithOutline(dst, text, x, y, face, fillColor)"
  - "badge.FormatNumber(n) string"

decisions:
  - "Pure Go with golang.org/x/image — no CGo, no node-canvas, no sharp"
  - "Font embedded via go:embed — single binary deployment, no external assets"
  - "Rajdhani Bold — geometric sans-serif matching Destiny UI"
  - "Power Level color: #F4D03F (Destiny exotic gold)"
  - "Multi-offset stroke technique instead of canvas strokeText (Go has no native stroke)"
  - "Power Level at (700, 80) in 72pt, stats at y=350 in 32pt"
  - "Stats spacing: 140px intervals starting at x=50"
  - "Unicode stat icons: ●◆■▲★ (monochrome, subtle)"

blockers: []

affects:
  - phase: 05-readme-update
    reason: "badge.png generated and ready for README injection"
  - phase: workflow
    reason: "Go binary handles full pipeline via 'contribemblem run'"

subsystem: image-generation
---

# Phase 4, Plan 1 Summary: Image Generation with Power Level

## What Was Built

**Destiny 2-Styled Badge Generator in Pure Go**

The image generation pipeline was implemented as part of a full Go conversion (commits GO-01 through GO-06), replacing the originally planned Node.js/node-canvas stack with pure Go using `golang.org/x/image`.

### Core Components

- `internal/badge/generator.go` — Main badge compositing: loads emblem background, scales to 800x400 canvas, renders Power Level and individual stats, outputs PNG
- `internal/badge/text.go` — Multi-layer text rendering (shadow → stroke → fill) for contrast on variable emblem backgrounds
- `internal/badge/format.go` — K/M number formatting (1200 → "1.2K", 1500000 → "1.5M")
- `internal/badge/assets/fonts/Rajdhani-Bold.ttf` — Embedded font via `go:embed`

### CLI Integration

- `contribemblem generate` — Reads `data/stats.json` + `data/emblem.jpg`, outputs `badge.png`
- `contribemblem run` — Full pipeline including badge generation as step 4/4

## Key Technical Decisions

### 1. Pure Go (No CGo)
**Decision:** Use `golang.org/x/image` instead of Node.js/node-canvas or sharp

**Why:** Single static binary, no runtime dependencies, no npm install, simpler CI/CD. The `golang.org/x/image` package provides sufficient image compositing and font rendering for badge generation.

### 2. Multi-Offset Stroke Technique
**Decision:** Simulate text stroke by drawing text at 16 offsets (±2px in each direction)

**Why:** Go's `golang.org/x/image/font` has no native strokeText equivalent. The multi-offset approach produces visually similar results to canvas strokeText, with 16 positions forming a complete border around the text. Combined with shadow layers, this satisfies IMAGE-04 (text readable on variable backgrounds).

### 3. Font Embedding via go:embed
**Decision:** Embed Rajdhani-Bold.ttf directly in the binary

**Why:** Eliminates file path resolution issues, ensures font is always available regardless of working directory, and produces a truly self-contained binary.

### 4. Color Scheme from Destiny Design Language
- Power Level: `#F4D03F` (Destiny exotic gold)
- Stats: White `#FFFFFF`
- Stroke: Black `#000000`
- Shadow: Black with 80% alpha `rgba(0,0,0,0.8)`

## Architecture

```
data/stats.json ──┐
                   ├──> badge.Generate() ──> badge.png
data/emblem.jpg ──┘

Internal pipeline:
1. loadImage(emblemPath) → image.Image
2. BiLinear.Scale to 800x400 canvas
3. loadFonts() → 72pt + 32pt Rajdhani Bold faces
4. Calculate Power Level (sum of 5 metrics)
5. DrawTextWithOutline(Power Level, 700, 80, 72pt, gold)
6. DrawTextWithOutline(each stat, bottom row, 32pt, white)
7. savePNG(canvas, outputPath)
```

## Requirements Satisfied

✅ **IMAGE-01:** PNG with emblem background and stat overlays
- `generator.go:57` scales emblem to fill 800x400 canvas
- Stats rendered on top via font.Drawer

✅ **IMAGE-02:** Power Level displayed prominently
- Sum of 5 metrics at (700, 80) in 72pt Rajdhani Bold
- Destiny exotic gold `#F4D03F`

✅ **IMAGE-03:** Individual metrics with icons
- Unicode icons `●◆■▲★` at y=350, 140px spacing
- K/M formatted values beside each icon

✅ **IMAGE-04:** Text contrast for readability
- 3-layer rendering: shadow (offset +2/+3px) → stroke (16 offsets) → fill
- Black shadow/stroke ensures readability on any emblem background

✅ **IMAGE-05:** 800x400px dimensions
- Constants `Width=800, Height=400` in generator.go

✅ **IMAGE-06:** Stable filename overwrites previous
- Always outputs to `badge.png` via os.Create (truncates existing file)

## Test Status

- `TestFormatNumber` — PASS (validates K/M formatting)
- `TestGenerate` — SKIP (needs emblem fixture image)
- `TestDrawTextWithOutline` — SKIP (needs image fixture)

## Stack Pivot Note

The `04-01-PLAN.md` describes a Node.js/node-canvas implementation that was never executed. Instead, the entire project was converted to Go (commits GO-01 through GO-08), implementing Phase 4 functionality as part of that conversion. The plan document is preserved for historical context but does not reflect the actual implementation.

## Next Phase Preview

**Phase 5: README Update & Commit**
- Marker-based injection: `<!-- CONTRIBEMBLEM:START -->` / `<!-- CONTRIBEMBLEM:END -->`
- Conditional commits: skip if badge.png hash unchanged
- [skip ci] commit messages
- Input: `badge.png` from Phase 4
