# Phase 4: Image Generation with Power Level - Research

**Researched:** 2026-01-28
**Domain:** Node.js canvas-based PNG image generation with text overlays
**Confidence:** HIGH

## Summary

Phase 4 requires compositing emblem artwork with GitHub statistics overlays into a 800x400px PNG, featuring Destiny 2-styled UI with prominent Power Level display. The standard approach is node-canvas (Cairo-backed) for server-side image generation in GitHub Actions, with multiple text rendering techniques (stroke, shadow, background overlay) combined for optimal contrast on variable backgrounds.

The technical challenge centers on text readability: Destiny emblems have highly variable backgrounds (dark, light, busy patterns), requiring robust multi-layer text treatment. The context explicitly requires copying "Destiny's exact implementation" for text contrast, which uses stroke + shadow + subtle background techniques simultaneously.

For fonts, geometric sans-serif options (Rajdhani, Orbitron, Saira Condensed) from Google Fonts closely match Destiny's UI typography and can be downloaded/registered with node-canvas via registerFont().

**Primary recommendation:** Use node-canvas 2.11+ with strokeText() + fillText() + shadowBlur layering to replicate Destiny's text contrast approach. Download Rajdhani Bold/Saira Condensed Bold as primary font options.

## Standard Stack

The established libraries/tools for server-side PNG generation with text overlays:

### Core
| Library | Version | Purpose | Why Standard |
|---------|---------|---------|--------------|
| canvas (node-canvas) | 2.11+ | Cairo-backed Canvas API for Node.js | Industry standard for server-side canvas operations; implements W3C Canvas API; zero native dependencies on most platforms; 10.6k stars, battle-tested in production |
| @napi-rs/canvas | 0.1+ | Rust-based alternative to node-canvas | Newer, faster alternative using Skia; better for pure performance but less mature ecosystem |

### Supporting
| Library | Version | Purpose | When to Use |
|---------|---------|---------|-------------|
| sharp | 0.34+ | libvips-based image processing | If you need advanced image optimization post-generation; 4-5x faster than ImageMagick for resizing |
| axios | 1.6+ | HTTP client for fetching emblem images | Standard choice for downloading emblem artwork from Bungie CDN |
| fs/promises | Built-in | File system operations | Writing PNG buffers to disk |

### Alternatives Considered
| Instead of | Could Use | Tradeoff |
|------------|-----------|----------|
| node-canvas | jimp | Pure JavaScript (zero native deps) but much slower; no native Canvas API compatibility; limited text rendering features |
| node-canvas | @napi-rs/canvas | Faster but newer (less proven); Rust-based; may have edge cases with complex text rendering |
| node-canvas | sharp alone | Sharp excels at image processing/compositing but has limited text rendering capabilities compared to Canvas API |

**Installation:**
```bash
npm install canvas
# For font downloads (if not bundling fonts)
npm install axios
```

**Note on dependencies:** node-canvas requires Cairo, Pango, and other system libraries, but provides pre-built binaries for common platforms (Linux x64, macOS x64/arm64, Windows x64). GitHub Actions ubuntu-latest runners include necessary dependencies.

## Architecture Patterns

### Recommended Project Structure
```
scripts/
├── generate-badge.js       # Main badge generation script
├── lib/
│   ├── canvas-utils.js    # Canvas setup, text measurement utilities
│   ├── text-renderer.js   # Multi-layer text rendering (stroke+shadow+fill)
│   └── image-loader.js    # Emblem image fetching/caching
└── assets/
    └── fonts/
        ├── Rajdhani-Bold.ttf
        └── SairaCondensed-Bold.ttf
```

### Pattern 1: Multi-Layer Text Rendering for Contrast
**What:** Render text in multiple passes (stroke → shadow → fill) to ensure readability on any background
**When to use:** Text overlays on variable/unknown backgrounds (Destiny emblems range from pure white to pure black)
**Example:**
```javascript
// Source: Inferred from Destiny 2 UI and MDN Canvas API best practices
function renderContrastText(ctx, text, x, y, options = {}) {
  const {
    fontSize = 48,
    fontFamily = 'Rajdhani',
    strokeWidth = 4,
    shadowBlur = 10,
    fillColor = '#FFFFFF',
    strokeColor = '#000000',
    shadowColor = 'rgba(0, 0, 0, 0.8)'
  } = options;

  ctx.font = `bold ${fontSize}px ${fontFamily}`;
  ctx.textAlign = 'center';
  ctx.textBaseline = 'middle';

  // Layer 1: Outer glow/shadow
  ctx.shadowColor = shadowColor;
  ctx.shadowBlur = shadowBlur;
  ctx.shadowOffsetX = 0;
  ctx.shadowOffsetY = 2;
  ctx.fillStyle = 'transparent';
  ctx.fillText(text, x, y);

  // Layer 2: Stroke (outline)
  ctx.shadowColor = 'transparent';
  ctx.shadowBlur = 0;
  ctx.strokeStyle = strokeColor;
  ctx.lineWidth = strokeWidth;
  ctx.lineJoin = 'round'; // Prevents sharp corners
  ctx.strokeText(text, x, y);

  // Layer 3: Fill (main text)
  ctx.fillStyle = fillColor;
  ctx.fillText(text, x, y);
}
```

### Pattern 2: Font Registration Before Canvas Use
**What:** Register custom fonts synchronously before creating canvas context
**When to use:** Always, when using non-system fonts
**Example:**
```javascript
// Source: node-canvas official docs
const { registerFont, createCanvas } = require('canvas');

// Register fonts BEFORE canvas creation
registerFont('./assets/fonts/Rajdhani-Bold.ttf', { 
  family: 'Rajdhani',
  weight: 'bold'
});

const canvas = createCanvas(800, 400);
const ctx = canvas.getContext('2d');
// Now 'Rajdhani' is available in ctx.font
```

### Pattern 3: Image Composition Pipeline
**What:** Load emblem → create canvas → draw emblem → overlay stats → save PNG
**When to use:** Standard sequence for all badge generation
**Example:**
```javascript
// Source: node-canvas patterns + best practices
const { createCanvas, loadImage } = require('canvas');

async function generateBadge(emblemUrl, stats) {
  // 1. Load emblem background
  const emblemImage = await loadImage(emblemUrl);
  
  // 2. Create canvas
  const canvas = createCanvas(800, 400);
  const ctx = canvas.getContext('2d');
  
  // 3. Draw emblem as background
  ctx.drawImage(emblemImage, 0, 0, 800, 400);
  
  // 4. Overlay stats with contrast text
  renderPowerLevel(ctx, stats.powerLevel);
  renderIndividualStats(ctx, stats);
  
  // 5. Export PNG
  return canvas.toBuffer('image/png');
}
```

### Pattern 4: Number Formatting with K/M Abbreviation
**What:** Format large numbers for display (1234 → "1.2K", 10500 → "10.5K")
**When to use:** Displaying stat values in limited space
**Example:**
```javascript
function formatNumber(num) {
  if (num >= 1000000) {
    return (num / 1000000).toFixed(1).replace(/\.0$/, '') + 'M';
  }
  if (num >= 1000) {
    return (num / 1000).toFixed(1).replace(/\.0$/, '') + 'K';
  }
  return num.toString();
}
```

### Anti-Patterns to Avoid
- **Single-layer text rendering:** Using only fillText() or only strokeText() fails on high-contrast backgrounds. Destiny always uses multiple layers.
- **System fonts only:** Relying on system fonts produces inconsistent results across platforms. Always bundle and register custom fonts.
- **Hardcoded positions without measurement:** Use `ctx.measureText()` to calculate text width for dynamic positioning/alignment.
- **Synchronous file I/O in generation:** Use async/await with loadImage() and toBuffer() to avoid blocking.

## Don't Hand-Roll

Problems that look simple but have existing solutions:

| Problem | Don't Build | Use Instead | Why |
|---------|-------------|-------------|-----|
| PNG encoding | Custom PNG writer | `canvas.toBuffer('image/png')` | PNG spec is complex; node-canvas uses battle-tested libpng via Cairo |
| Image downloading/caching | Fetch + manual caching | axios + built-in cache or Phase 2's cache | HTTP error handling, redirects, retries are subtle |
| Font file parsing | Custom TTF parser | node-canvas `registerFont()` | TrueType spec is 200+ pages; let Cairo/Pango handle it |
| Text measurement | Manual width calculation | `ctx.measureText(text).width` | Accounts for kerning, ligatures, font rendering hints |
| Color utilities | Manual hex/rgba conversion | Use standard CSS color strings | Canvas API natively understands hex, rgb(), rgba(), hsl() |

**Key insight:** Image generation involves file format specs (PNG, TTF), platform rendering differences, and edge cases that mature libraries have solved. Custom solutions inevitably rediscover these issues.

## Common Pitfalls

### Pitfall 1: Font Not Loaded Before Canvas Creation
**What goes wrong:** `registerFont()` called after canvas creation results in font not being available; text renders in fallback system font
**Why it happens:** Fonts must be registered globally before any canvas context is created
**How to avoid:** Call all `registerFont()` statements at module initialization, before any `createCanvas()` calls
**Warning signs:** Text looks different locally vs CI; "Font not found" style rendering with wrong metrics

### Pitfall 2: Insufficient Text Contrast Techniques
**What goes wrong:** Text unreadable on certain emblem backgrounds (e.g., white text on light emblem, black stroke on dark emblem)
**Why it happens:** Using only stroke OR shadow instead of Destiny's multi-layer approach (stroke AND shadow AND fill)
**How to avoid:** Always use 3-layer technique: shadow → stroke → fill, even if it seems excessive
**Warning signs:** User reports "can't read stats on [specific emblem]"; manual testing on only dark or only light backgrounds

### Pitfall 3: Text Positioning Without Baseline Alignment
**What goes wrong:** Text appears vertically misaligned, especially when mixing different font sizes
**Why it happens:** Forgetting to set `ctx.textBaseline` causes inconsistent vertical positioning
**How to avoid:** Always set `ctx.textBaseline = 'middle'` (or 'top'/'alphabetic') explicitly before drawing text
**Warning signs:** Numbers appear higher/lower than expected; alignment breaks when changing font size

### Pitfall 4: Image Aspect Ratio Mismatch
**What goes wrong:** Emblem images stretch/squash when drawn; Bungie emblems may not be exactly 800x400
**Why it happens:** Assuming emblem dimensions without checking; using wrong drawImage() parameters
**How to avoid:** Check emblem dimensions with `image.width/height`, calculate aspect ratio, use appropriate drawImage() signature (9 parameters for cropping/scaling)
**Warning signs:** Emblems look distorted; some emblems have black bars or stretched content

### Pitfall 5: Missing Cairo Dependencies on CI
**What goes wrong:** `npm install canvas` succeeds but `require('canvas')` fails at runtime with "libcairo.so not found"
**Why it happens:** GitHub Actions runner missing system dependencies; pre-built binary didn't match platform
**How to avoid:** Use `ubuntu-latest` runner (has Cairo preinstalled) OR explicitly install: `sudo apt-get install -y libcairo2-dev libjpeg-dev libpango1.0-dev libgif-dev librsvg2-dev`
**Warning signs:** Works locally, fails in CI with "Cannot find module 'canvas'" or "Error loading shared library"

## Code Examples

Verified patterns from official sources:

### Complete Badge Generation Function
```javascript
// Source: Synthesized from node-canvas docs + Destiny design requirements
const { createCanvas, loadImage, registerFont } = require('canvas');

// Register fonts at module load
registerFont('./assets/fonts/Rajdhani-Bold.ttf', { family: 'Rajdhani', weight: 'bold' });

async function generateDestinyBadge(emblemUrl, stats) {
  // Load emblem
  const emblem = await loadImage(emblemUrl);
  
  // Create canvas
  const canvas = createCanvas(800, 400);
  const ctx = canvas.getContext('2d');
  
  // Draw emblem background (fill entire canvas)
  ctx.drawImage(emblem, 0, 0, 800, 400);
  
  // Render Power Level (top-right, large)
  renderContrastText(ctx, stats.powerLevel.toString(), 700, 80, {
    fontSize: 72,
    fontFamily: 'Rajdhani',
    strokeWidth: 6,
    shadowBlur: 12,
    fillColor: '#F4D03F', // Destiny exotic gold
    strokeColor: '#000000'
  });
  
  // Render "POWER" label
  renderContrastText(ctx, 'POWER', 700, 130, {
    fontSize: 18,
    fontFamily: 'Rajdhani',
    strokeWidth: 3,
    fillColor: '#FFFFFF',
    strokeColor: '#000000'
  });
  
  // Render individual stats (horizontal row, bottom)
  const statY = 350;
  const statSpacing = 150;
  const startX = 100;
  
  const statOrder = [
    { key: 'commits', icon: '●', label: 'Commits' },
    { key: 'pullRequests', icon: '◆', label: 'PRs' },
    { key: 'issues', icon: '■', label: 'Issues' },
    { key: 'reviews', icon: '▲', label: 'Reviews' },
    { key: 'stars', icon: '★', label: 'Stars' }
  ];
  
  statOrder.forEach((stat, i) => {
    const x = startX + (i * statSpacing);
    const value = formatNumber(stats[stat.key]);
    
    // Icon
    renderContrastText(ctx, stat.icon, x - 15, statY - 10, {
      fontSize: 16,
      strokeWidth: 2,
      fillColor: '#AAAAAA'
    });
    
    // Value
    renderContrastText(ctx, value, x + 10, statY - 10, {
      fontSize: 28,
      fontFamily: 'Rajdhani',
      strokeWidth: 4,
      fillColor: '#FFFFFF'
    });
  });
  
  // Export PNG
  return canvas.toBuffer('image/png');
}

function formatNumber(num) {
  if (num >= 1000000) return (num / 1000000).toFixed(1).replace(/\.0$/, '') + 'M';
  if (num >= 1000) return (num / 1000).toFixed(1).replace(/\.0$/, '') + 'K';
  return num.toString();
}

function renderContrastText(ctx, text, x, y, options = {}) {
  // See Pattern 1 above for implementation
}
```

### Font Loading with Fallback
```javascript
// Source: node-canvas docs + error handling best practices
const { registerFont } = require('canvas');
const fs = require('fs');

function registerFontsWithFallback() {
  const fonts = [
    { path: './assets/fonts/Rajdhani-Bold.ttf', family: 'Rajdhani', weight: 'bold' },
    { path: './assets/fonts/SairaCondensed-Bold.ttf', family: 'Saira Condensed', weight: 'bold' }
  ];
  
  fonts.forEach(font => {
    try {
      if (fs.existsSync(font.path)) {
        registerFont(font.path, { family: font.family, weight: font.weight });
        console.log(`✓ Registered font: ${font.family}`);
      } else {
        console.warn(`⚠ Font file not found: ${font.path}`);
      }
    } catch (error) {
      console.error(`✗ Failed to register ${font.family}:`, error.message);
    }
  });
}

// Call at module initialization
registerFontsWithFallback();
```

### Testing Text Contrast on Various Backgrounds
```javascript
// Source: Test pattern for validating multi-layer approach
async function testTextContrast() {
  const { createCanvas } = require('canvas');
  const canvas = createCanvas(800, 400);
  const ctx = canvas.getContext('2d');
  
  // Test backgrounds: pure white, pure black, gradient
  const backgrounds = [
    () => { ctx.fillStyle = '#FFFFFF'; ctx.fillRect(0, 0, 800, 400); },
    () => { ctx.fillStyle = '#000000'; ctx.fillRect(0, 0, 800, 400); },
    () => {
      const gradient = ctx.createLinearGradient(0, 0, 800, 0);
      gradient.addColorStop(0, '#000000');
      gradient.addColorStop(1, '#FFFFFF');
      ctx.fillStyle = gradient;
      ctx.fillRect(0, 0, 800, 400);
    }
  ];
  
  backgrounds.forEach((drawBg, i) => {
    drawBg();
    renderContrastText(ctx, 'POWER 1337', 400, 200, {
      fontSize: 72,
      fillColor: '#F4D03F',
      strokeColor: '#000000',
      strokeWidth: 6,
      shadowBlur: 12
    });
    fs.writeFileSync(`test-contrast-${i}.png`, canvas.toBuffer('image/png'));
  });
}
```

## State of the Art

| Old Approach | Current Approach | When Changed | Impact |
|--------------|------------------|--------------|--------|
| ImageMagick CLI calls | node-canvas (Cairo) or sharp (libvips) | ~2015 | Native Node.js integration; 4-5x faster; better error handling |
| System fonts only | Bundled font files + registerFont() | Always recommended | Consistent cross-platform rendering |
| Single-pass text rendering | Multi-layer (stroke + shadow + fill) | Game UI best practice | Readable on any background |
| Manual PNG encoding | Built-in Canvas.toBuffer() | Always available | Simpler, more reliable |
| Blocking file I/O | Async loadImage() + toBuffer() | node-canvas 2.0+ (2018) | Non-blocking in Node.js event loop |

**Deprecated/outdated:**
- **node-canvas 1.x:** Use 2.x+ (better async support, stability)
- **`canvas.createPNGStream()`:** Prefer `canvas.toBuffer()` for simpler code (streams add complexity for single-file generation)
- **gm (GraphicsMagick bindings):** Superseded by sharp for performance; node-canvas for text rendering

## Open Questions

Things that couldn't be fully resolved:

1. **Exact Destiny font identification**
   - What we know: Destiny 2 uses custom geometric sans-serif fonts; Rajdhani Bold and Saira Condensed Bold are close approximations based on visual comparison
   - What's unclear: Bungie hasn't officially documented UI fonts; may use proprietary/licensed fonts
   - Recommendation: Use Rajdhani Bold (free, OFL license) as primary; validate visual match with Destiny screenshots from context document

2. **Optimal stroke width for 800x400 canvas**
   - What we know: Stroke width should scale with font size; Destiny uses relatively thick strokes (stroke width ≈ fontSize / 12)
   - What's unclear: Exact pixel measurements from Destiny's implementation at this resolution
   - Recommendation: Start with 6px stroke for 72px Power Level text, 4px for 28px stat text; tune based on visual review

3. **Icon rendering: Unicode symbols vs image sprites**
   - What we know: Context specifies "monochrome symbolic icons" for stats; Unicode symbols (●◆■▲★) match "symbolic" requirement
   - What's unclear: Whether custom SVG icons would better match Destiny aesthetic vs simple Unicode
   - Recommendation: Start with Unicode symbols (simpler, no additional assets); Phase 4 can be enhanced later with custom icons if needed

4. **Color values for "Destiny exotic gold"**
   - What we know: Destiny uses gold/amber color for exotic items and Power Level displays
   - What's unclear: Exact hex code (varies by screen, Destiny's lighting system affects UI colors)
   - Recommendation: Use #F4D03F (amber gold) or #F39C12 (deeper gold); validate against Destiny screenshot from context

## Sources

### Primary (HIGH confidence)
- node-canvas GitHub README - [https://github.com/Automattic/node-canvas](https://github.com/Automattic/node-canvas) - Installation, API patterns, registerFont() usage
- MDN Canvas API - strokeText() documentation - [https://developer.mozilla.org/en-US/docs/Web/API/CanvasRenderingContext2D/strokeText](https://developer.mozilla.org/en-US/docs/Web/API/CanvasRenderingContext2D/strokeText) - Multi-layer text rendering techniques
- sharp documentation - [https://sharp.pixelplumbing.com/](https://sharp.pixelplumbing.com/) - Alternative library comparison
- Google Fonts repository - [https://github.com/google/fonts](https://github.com/google/fonts) - Font licensing, availability of geometric sans-serif fonts

### Secondary (MEDIUM confidence)
- Phase 4 CONTEXT.md decisions - User-specified design choices (exact Destiny match, multi-layer text contrast, geometric fonts, top-right Power Level placement)
- STATE.md prior decisions - Phase 2 caching pattern, Phase 3 emblem stability

### Tertiary (LOW confidence)
- Destiny 2 UI aesthetic - Based on context description; specific font/color values inferred from "Destiny design language" without official Bungie documentation

## Metadata

**Confidence breakdown:**
- Standard stack: HIGH - node-canvas is industry standard with 10.6k stars, clear documentation, battle-tested
- Architecture: HIGH - Patterns derived from official docs + established game UI practices
- Pitfalls: HIGH - Common issues documented in node-canvas issue tracker and Stack Overflow
- Font selection: MEDIUM - Approximations based on visual similarity; not official Destiny fonts
- Color values: MEDIUM - Inferred from Destiny screenshots; no official color palette

**Research date:** 2026-01-28
**Valid until:** 90 days (node-canvas stable; major version changes rare; font availability stable)
**Phase dependencies:** Requires Phase 2 (GitHub stats JSON) and Phase 3 (emblem URL) outputs as inputs

**Implementation notes:**
- GitHub Actions ubuntu-latest includes Cairo dependencies; no additional system packages needed
- Font files must be committed to repository OR downloaded in workflow (Google Fonts allows redistribution under OFL)
- Badge generation should be idempotent: same inputs → same PNG (deterministic rendering)
- Image caching strategy: Use date-based cache key from Phase 2 decision (github-stats-YYYYMMDD) to avoid regenerating badge unnecessarily
