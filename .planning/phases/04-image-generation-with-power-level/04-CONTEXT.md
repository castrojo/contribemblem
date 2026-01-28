# Phase 4: Image Generation with Power Level - Context

**Gathered:** 2026-01-28
**Status:** Ready for planning

<domain>
## Phase Boundary

Generate a Destiny 2-styled badge PNG (800x400px) that composites emblem artwork from Bungie API (Phase 3 output) with GitHub contribution stats overlay (Phase 2 output), featuring prominent Power Level display in Destiny's design language.

</domain>

<decisions>
## Implementation Decisions

### Visual Hierarchy & Layout
- **Power Level placement:** Top-right corner (matching Destiny character screen locations exactly)
- **Individual stats arrangement:** Horizontal row with symbolic icons (following Destiny design language)
- **Information density:** Medium - stat values with symbolic icons from emoji palette (no text labels)
- **Canvas usage:** Full emblem background fills 800x400px, stats overlay on top

### Typography & Styling
- **Font choice:** Destiny-style geometric (bold, geometric sans-serif matching Destiny's UI fonts)
- **Text contrast:** Multiple techniques combined (stroke, background overlay, drop shadow) - copy Destiny's exact implementation for numbers and letters
- **Text sizing:** Based on Destiny design language - copy their design exactly
- **Color scheme:** Copy the exact way Destiny does it (gold for exotic/Power Level, white for stats, etc.)

### Destiny Design Language
- **Visual reference:** Character screen emblem style (top-right corner display in Destiny)
- **UI chrome:** No chrome - clean emblem + overlay only, no geometric borders or accents
- **Aesthetic fidelity:** Exact match - should look indistinguishable from Destiny emblems
- **Visual effects:** Claude's discretion to match the reference

### Stat Display Format
- **Icon style:** Monochrome symbolic icons, subtle appearance (not colorful emoji)
- **Icon placement:** Horizontal pair - icon left of number
- **Stat order:** Keep Phase 2 JSON order (commits, PRs, issues, reviews, stars)
- **Number formatting:** Abbreviate large numbers (K/M format: 1.2K, 10.5K, etc.)

### Claude's Discretion
- Exact visual effects implementation (glows, shadows, gradients) to match Destiny reference
- Specific font file selection within geometric sans-serif category
- Precise spacing and padding for visual balance
- Loading skeleton or progress indicator during generation

</decisions>

<specifics>
## Specific Ideas

**Primary Reference:** Destiny 2 character screen emblems (top-right corner display where Power Level appears)

**Design Philosophy:** "Should look exactly design-wise to the ones that come in Destiny" - exact match, not inspired by or loosely based on. The goal is for someone familiar with Destiny to recognize it immediately as matching Destiny's aesthetic.

**Text Treatment:** Study how Destiny renders numbers and letters over varying emblem backgrounds - they use multiple contrast techniques simultaneously. This exact approach should be replicated.

**Symbolic Icons:** Subtle, monochrome icons that don't overpower the numbers. Think of Destiny's stat icons in menus - understated but clear.

</specifics>

<deferred>
## Deferred Ideas

None - discussion stayed within phase scope

</deferred>

---

*Phase: 04-image-generation-with-power-level*
*Context gathered: 2026-01-28*
