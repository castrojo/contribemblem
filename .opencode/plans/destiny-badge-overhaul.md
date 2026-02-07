# Plan: Destiny 2-Style Emblem Badge Visual Overhaul

## Goal
Make the badge generator output look like actual in-game Destiny 2 emblem banners.

## Current State
The current badge is a raw emblem background with floating text (username, power level, stats). It lacks the layered, polished look of Destiny 2's in-game emblem banners.

## Phases

### Phase 1: Dark Gradient Overlay
- Draw a horizontal gradient overlay on the emblem background
- Left side: fully transparent (emblem art shows through)
- Right side: semi-opaque black (~100 alpha at far right), starting at 40% from left
- Creates readable area for the power level number

### Phase 2: Bottom Stat Bar
- Draw a semi-transparent dark bar (rgba(0,0,0,150)) across bottom 36px
- Render stats centered in cells within the bar
- Add thin vertical divider lines (rgba(255,255,255,50)) between each stat
- Center each stat label+value within its cell

### Phase 3: Top Accent Line
- Draw a 3px gold accent line (#F4D03F) along the top edge
- Matches Destiny 2's colored accent on emblem banners

### Phase 4: Power Level Improvements
- Right-align the power level so it sits flush with right edge (minus margin)
- Add diamond icon (◆) to the left of the power level number
- Use measureText() to dynamically calculate positions

### Phase 5: Typography Improvements
- 3 font sizes instead of 2:
  - Large (48pt): Power level
  - Medium (26pt): Username (was 20pt)
  - Small (14pt): Stat labels and values
- Stat labels ("COMMITS", "PRS", etc.) in dim white (180, 180, 190)
- Stat values in bright white
- Labels and values rendered separately for different colors

### Phase 6: Overall Polish
- 1px border around entire badge in dark grey (60, 60, 65)
- Slight overall darken overlay (rgba(0,0,0,25)) for Destiny dark UI feel
- Consistent 20px horizontal margins

## Files to Modify

### `internal/badge/generator.go`
- Add new color constants: DimWhiteColor, AccentColor, StatBarColor, DividerColor, BorderColor, OverlayDark
- Add layout constants: marginX, marginTop, statBarHeight, accentHeight, borderWidth, statDividerW
- Change loadFonts() to return FontFaces struct with 3 sizes (Large/Medium/Small)
- Add helper functions: drawRect, drawHorizontalGradient, blendOver, drawBorder, measureText
- Rewrite Generate() to use all 6 phases in order:
  1. Scale emblem background
  2. Overall darken overlay
  3. Horizontal gradient overlay
  4. Stat bar rectangle
  5. Accent line rectangle
  6. Border
  7. Render username with Medium face at (marginX+4, accentHeight+marginTop+28)
  8. Render power level right-aligned with diamond icon using Large face
  9. Render stats centered in stat bar cells with Small face, labels dim, values bright

### `internal/badge/generator_test.go`
- Update test to work with new loadFonts() signature (returns *FontFaces instead of two faces)
- Verify tests still pass

### `internal/badge/text_test.go`
- May need updates if test calls loadFonts directly
- Verify tests pass

### Example regeneration
- Create temporary cmd/generate-examples/main.go
- Generate 3 example badges with same stats as before
- Delete the temp file
- Verify visual quality

## Execution Order
1. Modify generator.go with all 6 phases
2. Update generator_test.go for new font API
3. Run tests to verify
4. Regenerate example badges
5. Verify visuals
6. Commit and push
