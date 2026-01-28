#!/usr/bin/env bash
set -euo pipefail

# fetch-emblem.sh
# Downloads Destiny 2 emblem artwork from Bungie API
# Accepts emblem hash via stdin

# Configuration
BUNGIE_BASE_URL="https://www.bungie.net"
MANIFEST_API="${BUNGIE_BASE_URL}/Platform/Destiny2/Manifest/"
OUTPUT_FILE="data/emblem.jpg"
MANIFEST_CACHE="data/manifest.json"
USER_AGENT="ContribEmblem/1.0 (+https://github.com/castrojo/contribemblem)"

# Check for API key
if [ -z "${BUNGIE_API_KEY:-}" ]; then
  echo "ERROR: BUNGIE_API_KEY environment variable not set" >&2
  echo "Please add your Bungie API key to GitHub repository secrets" >&2
  echo "Get API key at: https://www.bungie.net/en/Application" >&2
  exit 1
fi

# Read emblem hash from stdin
read -r EMBLEM_HASH

echo "Fetching emblem hash: $EMBLEM_HASH"

# Fetch manifest metadata (to get current database URL)
echo "Fetching Bungie manifest metadata..."
MANIFEST_RESPONSE=$(curl -s -H "X-API-Key: $BUNGIE_API_KEY" \
                          -H "User-Agent: $USER_AGENT" \
                          "$MANIFEST_API")

# Check for API errors
ERROR_CODE=$(echo "$MANIFEST_RESPONSE" | jq -r '.ErrorCode')
if [ "$ERROR_CODE" != "1" ]; then
  ERROR_STATUS=$(echo "$MANIFEST_RESPONSE" | jq -r '.ErrorStatus')
  echo "ERROR: Bungie API returned error code $ERROR_CODE: $ERROR_STATUS" >&2
  exit 1
fi

# Log rate limit status (non-blocking monitoring)
RATE_LIMIT_REMAINING=$(curl -s -I -H "X-API-Key: $BUNGIE_API_KEY" \
                            -H "User-Agent: $USER_AGENT" \
                            "$MANIFEST_API" | grep -i "x-ratelimit-remaining" | cut -d: -f2 | tr -d ' \r' || echo "unknown")
echo "ℹ️  Rate limit remaining: $RATE_LIMIT_REMAINING"

# Extract DestinyInventoryItemDefinition URL
MANIFEST_URL=$(echo "$MANIFEST_RESPONSE" | jq -r '.Response.jsonWorldComponentContentPaths.en.DestinyInventoryItemDefinition')

if [ "$MANIFEST_URL" = "null" ] || [ -z "$MANIFEST_URL" ]; then
  echo "ERROR: Could not extract manifest URL from response" >&2
  exit 1
fi

MANIFEST_FULL_URL="${BUNGIE_BASE_URL}${MANIFEST_URL}"
echo "Manifest URL: $MANIFEST_FULL_URL"

# Download manifest if not cached
if [ ! -f "$MANIFEST_CACHE" ]; then
  echo "Downloading manifest database (~100MB, this may take a moment)..."
  curl -s -H "User-Agent: $USER_AGENT" \
       -o "$MANIFEST_CACHE" \
       "$MANIFEST_FULL_URL"
  echo "✓ Manifest cached"
else
  echo "✓ Using cached manifest"
fi

# Look up emblem in manifest
echo "Looking up emblem $EMBLEM_HASH in manifest..."
EMBLEM_DATA=$(jq -r --arg hash "$EMBLEM_HASH" '.[$hash]' "$MANIFEST_CACHE")

if [ "$EMBLEM_DATA" = "null" ] || [ -z "$EMBLEM_DATA" ]; then
  echo "ERROR: Emblem hash $EMBLEM_HASH not found in manifest" >&2
  exit 1
fi

# Extract icon path
ICON_PATH=$(echo "$EMBLEM_DATA" | jq -r '.displayProperties.icon')

if [ "$ICON_PATH" = "null" ] || [ -z "$ICON_PATH" ]; then
  echo "ERROR: Could not extract icon path from emblem data" >&2
  exit 1
fi

ICON_URL="${BUNGIE_BASE_URL}${ICON_PATH}"
echo "Downloading emblem image from: $ICON_URL"

# Download emblem image
HTTP_STATUS=$(curl -s -w "%{http_code}" -o "$OUTPUT_FILE" \
                   -H "User-Agent: $USER_AGENT" \
                   "$ICON_URL")

if [ "$HTTP_STATUS" != "200" ]; then
  echo "ERROR: Failed to download emblem image (HTTP $HTTP_STATUS)" >&2
  exit 1
fi

echo "✓ Emblem image saved to $OUTPUT_FILE"

# Verify it's a valid image
if file "$OUTPUT_FILE" | grep -qE "JPEG|PNG"; then
  echo "✓ Valid image format confirmed"
else
  echo "ERROR: Downloaded file is not a valid JPEG or PNG image" >&2
  exit 1
fi

exit 0
