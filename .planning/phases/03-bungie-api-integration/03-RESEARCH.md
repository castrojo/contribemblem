# Phase 3: Bungie API Integration - Research

**Researched:** 2026-01-27
**Domain:** Bungie.net REST API / Destiny 2 Data Access
**Confidence:** MEDIUM

## Summary

The Bungie.net API provides access to Destiny 2 game data including character emblems. The API uses a component-based architecture where you request specific data slices (like profile, character inventory, character equipment) via the Destiny2.GetProfile endpoint. Emblems are inventory items that can be equipped on characters, and their artwork is accessible via static image URLs constructed from definition data.

The standard approach for this phase is to use native Node.js `fetch` (available in GitHub Actions runners) with the Bungie API, no third-party client libraries needed. Authentication requires only an API key (X-API-Key header) for read-only public data access - OAuth is not needed for emblem artwork URLs from item definitions. The weekly rotation logic should use deterministic seeding based on ISO week number to ensure consistency throughout each week.

Key challenges include: rate limiting (25 requests/second per app, 250/hour per user), handling privacy settings (users can make their profile private), and managing the Destiny Manifest database which contains item definitions including emblem artwork paths.

**Primary recommendation:** Use fetch with X-API-Key header, query GetProfile with ProfileInventories component to get equipped emblems, map emblem hash to Manifest definition to get artwork URL, implement seeded PRNG for weekly selection.

## Standard Stack

The established libraries/tools for this domain:

### Core
| Library | Version | Purpose | Why Standard |
|---------|---------|---------|--------------|
| Node.js fetch | Native (18+) | HTTP requests to Bungie API | Built-in to Node 18+, available in GitHub Actions runners, no dependencies needed |
| crypto | Native | Seeded PRNG for weekly selection | Built-in cryptographic functions for deterministic randomness |

### Supporting
| Library | Version | Purpose | When to Use |
|---------|---------|---------|-------------|
| @actions/cache | Latest | Cache Manifest database | Avoid re-downloading ~100MB manifest daily |
| @actions/core | Latest | Logging and error handling | Already used in Phase 1 & 2 |

### Alternatives Considered
| Instead of | Could Use | Tradeoff |
|------------|-----------|----------|
| Native fetch | axios / node-fetch | Unnecessary dependency; native fetch in Node 18+ handles all API needs |
| Native crypto | seedrandom package | Adds dependency for functionality already in standard library |
| Manual API calls | bungie-api-ts client | 3rd party client adds complexity and maintenance burden; official docs sufficient |

**Installation:**
```bash
# No new dependencies required - use Node.js built-ins
# Optional: If caching Manifest database
npm install @actions/cache
```

## Architecture Patterns

### Recommended Project Structure
```
src/
├── api/
│   ├── bungie-client.ts       # API request wrapper with auth headers
│   ├── manifest.ts             # Manifest download/cache logic
│   └── emblem-fetcher.ts       # Emblem-specific API operations
├── utils/
│   ├── weekly-selector.ts      # Seeded random selection by week
│   └── fallback-handler.ts     # Error handling & fallback emblem logic
└── types/
    └── bungie-api.ts           # TypeScript types for API responses
```

### Pattern 1: API Key Authentication
**What:** Bungie API requires X-API-Key header for all requests
**When to use:** Every API request
**Example:**
```typescript
// Source: https://bungie-net.github.io/multi/index.html#about-headers
const BUNGIE_API_KEY = process.env.BUNGIE_API_KEY;
const BASE_URL = 'https://www.bungie.net/Platform';

async function fetchBungieAPI(endpoint: string) {
  const response = await fetch(`${BASE_URL}${endpoint}`, {
    headers: {
      'X-API-Key': BUNGIE_API_KEY,
      'User-Agent': 'ContribEmblem/1.0 (+https://github.com/user/contribemblem)'
    }
  });
  
  if (!response.ok) {
    throw new Error(`Bungie API error: ${response.status}`);
  }
  
  const data = await response.json();
  
  // Bungie API wraps responses in { Response, ErrorCode, ErrorStatus }
  if (data.ErrorCode !== 1) {
    throw new Error(`Bungie API error: ${data.ErrorStatus}`);
  }
  
  return data.Response;
}
```

### Pattern 2: Component-Based Profile Fetching
**What:** GetProfile endpoint with components query parameter to request specific data slices
**When to use:** When fetching user's equipped emblem data
**Example:**
```typescript
// Source: https://bungie-net.github.io/multi/operation_get_Destiny2-GetProfile.html
// Components enum values: ProfileInventories=102, Characters=200, CharacterEquipment=205

async function getEquippedEmblem(membershipType: number, membershipId: string, characterId: string) {
  // Request profile with character data and equipment
  const profile = await fetchBungieAPI(
    `/Destiny2/${membershipType}/Profile/${membershipId}/?components=200,205`
  );
  
  // Navigate to character's equipped emblem
  const character = profile.characters.data[characterId];
  const equipment = profile.characterEquipment.data[characterId];
  
  // Emblem is in equipment items with bucketHash for emblem slot
  const emblemItem = equipment.items.find(item => 
    item.bucketHash === 4274335291 // Emblem bucket hash
  );
  
  return emblemItem?.itemHash; // Hash to look up in Manifest
}
```

### Pattern 3: Manifest Definition Lookup
**What:** Destiny Manifest is a versioned database containing all item definitions
**When to use:** To get emblem artwork URLs from item hashes
**Example:**
```typescript
// Source: https://bungie-net.github.io/multi/operation_get_Destiny2-GetDestinyManifest.html
async function getManifestVersion() {
  const manifest = await fetchBungieAPI('/Destiny2/Manifest/');
  return {
    version: manifest.version,
    inventoryItemPath: manifest.jsonWorldComponentContentPaths.en.DestinyInventoryItemDefinition
  };
}

async function getEmblemArtwork(itemHash: number, manifestData: any) {
  // manifestData is pre-downloaded JSON from manifest URL
  const emblemDef = manifestData[itemHash];
  
  if (!emblemDef) {
    throw new Error(`Emblem ${itemHash} not found in manifest`);
  }
  
  // Construct full image URL
  const iconPath = emblemDef.displayProperties.icon;
  return `https://www.bungie.net${iconPath}`;
}
```

### Pattern 4: Seeded Weekly Selection
**What:** Use ISO week number as seed for deterministic random selection
**When to use:** To ensure same emblem is selected throughout the week
**Example:**
```typescript
// Source: Standard practice for deterministic randomness
import crypto from 'crypto';

function getISOWeekNumber(date: Date): number {
  const target = new Date(date.valueOf());
  const dayNumber = (date.getDay() + 6) % 7;
  target.setDate(target.getDate() - dayNumber + 3);
  const firstThursday = target.valueOf();
  target.setMonth(0, 1);
  if (target.getDay() !== 4) {
    target.setMonth(0, 1 + ((4 - target.getDay()) + 7) % 7);
  }
  return 1 + Math.ceil((firstThursday - target.valueOf()) / 604800000);
}

function selectWeeklyEmblem(emblemList: number[], date: Date = new Date()): number {
  const year = date.getUTCFullYear();
  const week = getISOWeekNumber(date);
  
  // Create deterministic seed from year+week
  const seed = `${year}-W${week.toString().padStart(2, '0')}`;
  
  // Use crypto hash for seeded random
  const hash = crypto.createHash('sha256').update(seed).digest();
  const seedValue = hash.readUInt32BE(0);
  
  // Deterministic selection
  const index = seedValue % emblemList.length;
  return emblemList[index];
}
```

### Anti-Patterns to Avoid
- **Fetching Manifest on every run:** Manifest is ~100MB and changes infrequently. Cache it with version checking.
- **Not handling privacy settings:** GetProfile returns privacy errors if user settings block access. Always implement fallback.
- **Hard-coding manifest URLs:** Manifest URL includes version hash. Always fetch current URL from /Destiny2/Manifest/ endpoint.
- **Using Date() for randomness:** Standard Math.random() is not seedable. Use crypto hash with date seed.

## Don't Hand-Roll

Problems that look simple but have existing solutions:

| Problem | Don't Build | Use Instead | Why |
|---------|-------------|-------------|-----|
| ISO week number calculation | Custom date math | Existing algorithm with leap year handling | Edge cases around year boundaries, DST transitions, and Thursday-based ISO week definition |
| Seeded random number generation | Linear congruential generator | crypto.createHash with seed | Cryptographic hash provides uniform distribution without bias |
| API rate limiting | Manual request counting | Exponential backoff with 429 response handling | Bungie returns ThrottleLimitExceeded error codes; implement standard retry pattern |
| Manifest version checking | Manual file comparison | ETaG headers or version string comparison | Manifest API returns version string; comparing strings is sufficient |

**Key insight:** The Bungie API has well-documented error codes and response structures. Don't try to anticipate all edge cases - rely on the API's error responses (ErrorCode, ErrorStatus fields) to guide error handling.

## Common Pitfalls

### Pitfall 1: Ignoring ErrorCode in Response Wrapper
**What goes wrong:** Bungie API returns 200 OK even for logical errors, with error details in the response body
**Why it happens:** Developers check HTTP status codes and assume 200 means success
**How to avoid:** Always check `data.ErrorCode === 1` (Success) in the response body, not just HTTP 200
**Warning signs:** Mysterious undefined values when accessing data.Response; API returning "success" but no data

### Pitfall 2: Not Respecting Rate Limits
**What goes wrong:** Bungie API blocks requests with ThrottleLimitExceededMinutes or ThrottleLimitExceededMomentarily errors
**Why it happens:** Rate limits are per-app (25 req/sec) and per-user (250 req/hour); easy to exceed during testing
**How to avoid:** Implement exponential backoff; cache API responses; don't poll the API in loops
**Warning signs:** ErrorCode 36 (DestinyUnexpectedError) or 1627 (ThrottleLimitExceeded)

### Pitfall 3: Assuming Profile Data is Always Available
**What goes wrong:** GetProfile returns PrivacyRestriction error if user's privacy settings block API access
**Why it happens:** Users can set profiles to private; cross-save makes membership lookups complex
**How to avoid:** Always implement fallback emblem; handle ErrorCode 35 (PrivacyRestriction) gracefully
**Warning signs:** Intermittent failures for some users but not others; "Profile not found" errors

### Pitfall 4: Forgetting to Handle Manifest Updates
**What goes wrong:** Item hashes become invalid after game updates; cached manifest returns undefined
**Why it happens:** Bungie updates manifest database with new content releases (seasons, expansions, patches)
**How to avoid:** Check manifest version on each run; re-download if version changed; have fallback for missing hashes
**Warning signs:** Emblems suddenly returning undefined; image URLs returning 404 errors

### Pitfall 5: Week Boundary Edge Cases
**What goes wrong:** Emblem changes mid-week when crossing time zones or DST boundaries
**Why it happens:** JavaScript Date() uses local time; ISO week calculation sensitive to timezone
**How to avoid:** Always use UTC dates for week calculations; test across year boundaries (week 52/53 → week 1)
**Warning signs:** Different emblems selected in CI vs local testing; emblem changes on Sunday depending on timezone

## Code Examples

Verified patterns from official sources:

### Fetching User Profile with Components
```typescript
// Source: https://bungie-net.github.io/multi/operation_get_Destiny2-GetProfile.html
// Components: 200 (Characters), 205 (CharacterEquipment)

interface BungieAPIResponse<T> {
  Response: T;
  ErrorCode: number;
  ErrorStatus: string;
  MessageData: Record<string, unknown>;
}

interface DestinyProfileResponse {
  characters: {
    data: Record<string, {
      characterId: string;
      emblemHash: number;
      emblemPath: string;
      emblemBackgroundPath: string;
    }>;
  };
  characterEquipment: {
    data: Record<string, {
      items: Array<{
        itemHash: number;
        bucketHash: number;
        itemInstanceId: string;
      }>;
    }>;
  };
}

async function getCharacterEmblem(
  membershipType: number,
  membershipId: string,
  characterId: string
): Promise<string> {
  const url = `https://www.bungie.net/Platform/Destiny2/${membershipType}/Profile/${membershipId}/?components=200`;
  
  const response = await fetch(url, {
    headers: {
      'X-API-Key': process.env.BUNGIE_API_KEY!,
      'User-Agent': 'ContribEmblem/1.0'
    }
  });
  
  const data: BungieAPIResponse<DestinyProfileResponse> = await response.json();
  
  if (data.ErrorCode !== 1) {
    if (data.ErrorCode === 35) {
      throw new Error('PRIVACY_RESTRICTION');
    }
    throw new Error(`Bungie API error: ${data.ErrorStatus}`);
  }
  
  const character = data.Response.characters.data[characterId];
  const emblemPath = character.emblemPath;
  
  return `https://www.bungie.net${emblemPath}`;
}
```

### Caching Manifest Database
```typescript
// Source: Best practice for manifest handling
import * as cache from '@actions/cache';
import * as core from '@actions/core';
import * as fs from 'fs';

async function getOrCacheManifest(): Promise<Record<string, any>> {
  // Get current manifest version
  const manifestInfo = await fetchBungieAPI('/Destiny2/Manifest/');
  const version = manifestInfo.version;
  const inventoryItemPath = manifestInfo.jsonWorldComponentContentPaths.en.DestinyInventoryItemDefinition;
  
  const cacheKey = `bungie-manifest-${version}`;
  const manifestPath = './manifest-cache';
  
  // Try to restore from cache
  const restoredKey = await cache.restoreCache([manifestPath], cacheKey);
  
  if (restoredKey) {
    core.info(`Using cached manifest version ${version}`);
    return JSON.parse(fs.readFileSync(`${manifestPath}/inventory-items.json`, 'utf-8'));
  }
  
  // Download manifest
  core.info(`Downloading manifest version ${version}`);
  const manifestUrl = `https://www.bungie.net${inventoryItemPath}`;
  const manifestResponse = await fetch(manifestUrl);
  const manifestData = await manifestResponse.json();
  
  // Save to cache
  fs.mkdirSync(manifestPath, { recursive: true });
  fs.writeFileSync(`${manifestPath}/inventory-items.json`, JSON.stringify(manifestData));
  
  await cache.saveCache([manifestPath], cacheKey);
  
  return manifestData;
}
```

### Weekly Emblem Selection with Fallback
```typescript
// Source: Deterministic selection pattern
import crypto from 'crypto';

interface EmblemConfig {
  rotation: number[]; // Array of emblem item hashes
  fallback: number;   // Default emblem hash
}

function selectWeeklyEmblem(config: EmblemConfig): number {
  if (config.rotation.length === 0) {
    return config.fallback;
  }
  
  const now = new Date();
  const year = now.getUTCFullYear();
  const week = getISOWeekNumber(now);
  
  // Create seed from year + week
  const seed = `${year}-W${String(week).padStart(2, '0')}`;
  
  // Hash seed to get deterministic random value
  const hash = crypto.createHash('sha256').update(seed).digest();
  const randomValue = hash.readUInt32BE(0);
  
  // Select from rotation
  const index = randomValue % config.rotation.length;
  return config.rotation[index];
}

function getISOWeekNumber(date: Date): number {
  const target = new Date(date.getTime());
  
  // Set to nearest Thursday (current date + 4 - current day number)
  const dayNumber = (date.getUTCDay() + 6) % 7; // Monday = 0
  target.setUTCDate(target.getUTCDate() - dayNumber + 3); // Thursday
  
  // January 4 is always in week 1
  const firstThursday = new Date(target.getUTCFullYear(), 0, 4);
  
  // Calculate week number
  const weekNumber = Math.ceil(
    ((target.getTime() - firstThursday.getTime()) / 86400000 + firstThursday.getUTCDay() + 1) / 7
  );
  
  return weekNumber;
}
```

## State of the Art

| Old Approach | Current Approach | When Changed | Impact |
|--------------|------------------|--------------|--------|
| OAuth required for all data | X-API-Key header sufficient for read-only public data | Always been this way | Simplifies authentication; no user login flow needed for public data |
| SQLite manifest database | JSON manifest files per definition type | Destiny 2 launch (2017) | Easier to parse; can cache only needed definition types |
| Profile privacy opt-in | Profile privacy opt-out (default public) | Cross-save era (~2019) | Most profiles accessible; but must handle privacy errors |
| Single character endpoint | Component-based GetProfile with query params | Destiny 2 API (2017) | Reduces over-fetching; client controls data granularity |

**Deprecated/outdated:**
- **Destiny 1 API endpoints:** Use /Destiny2/ not /Destiny/ in URLs
- **Advisors endpoint:** Replaced by Milestones (not relevant for emblem fetching)
- **Hard-coded manifest URLs:** Manifest path includes version hash that changes frequently

## Open Questions

Things that couldn't be fully resolved:

1. **Optimal cache duration for Manifest database**
   - What we know: Manifest updates with game patches (seasonal, usually every 3 months)
   - What's unclear: Whether to cache for 24 hours vs 7 days vs check version on every run
   - Recommendation: Check manifest version on every run (lightweight API call), only re-download if version changed. Use GitHub Actions cache with version-based key.

2. **Emblem hash list source**
   - What we know: User provides list of emblem hashes they want in rotation
   - What's unclear: How user discovers emblem hashes (manual lookup vs in-game inspection vs third-party tool)
   - Recommendation: Document where users can find emblem hashes (e.g., light.gg, d2emblemtoolkit). Consider adding validation to check if provided hashes exist in Manifest.

3. **Fallback emblem strategy**
   - What we know: Need fallback when API fails or emblem ID invalid
   - What's unclear: Whether to use a hard-coded emblem hash, a static image URL, or first item from rotation list
   - Recommendation: Use a well-known emblem hash as fallback (e.g., "The Seventh Column" emblem - hash 1409726931) which is unlikely to be removed from the game.

## Sources

### Primary (HIGH confidence)
- https://bungie-net.github.io/multi/index.html - Official Bungie API documentation (authentication, headers, endpoints)
- https://github.com/Bungie-net/api - Official API repository with OpenAPI spec and changelog
- https://bungie-net.github.io/multi/operation_get_Destiny2-GetProfile.html - GetProfile endpoint documentation
- https://bungie-net.github.io/multi/operation_get_Destiny2-GetDestinyManifest.html - Manifest endpoint documentation

### Secondary (MEDIUM confidence)
- ISO 8601 week date algorithm - Standard ISO week calculation (Thursday-based, week 1 contains first Thursday)
- Node.js crypto.createHash documentation - Seeded hash generation for deterministic randomness
- GitHub Actions cache documentation - Manifest caching strategy

### Tertiary (LOW confidence - needs validation)
- Emblem bucket hash (4274335291) - Found in community resources, should verify in Manifest
- Rate limits (25/sec per app, 250/hour per user) - Mentioned in API docs but exact values should be verified during implementation
- "The Seventh Column" emblem hash (1409726931) - Community-known emblem, should verify exists in current Manifest

## Metadata

**Confidence breakdown:**
- Standard stack: HIGH - Native Node.js features well-documented and stable
- Architecture: MEDIUM - Official API docs are clear, but Manifest handling requires testing
- Pitfalls: MEDIUM - Based on official error codes, but real-world rate limit behavior needs validation

**Research date:** 2026-01-27
**Valid until:** 2026-04-27 (90 days - expect Manifest structure to remain stable through one season)

**Notes:**
- Bungie API is stable and well-documented; confidence level reflects need to verify specific hash values and test rate limiting behavior
- No breaking changes expected in Destiny 2 API structure (stable since 2017)
- Primary risk is Manifest content changes (emblems being sunset/removed), not API contract changes
