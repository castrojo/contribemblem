#!/usr/bin/env bash
set -euo pipefail

# select-emblem.sh
# Deterministic weekly emblem selection based on ISO week number
# Uses SHA256-seeded selection for consistency within each calendar week

CONFIG_FILE="data/emblem-config.json"
FALLBACK_EMBLEM="1409726931"  # "The Seventh Column" emblem

# Check if config file exists
if [ ! -f "$CONFIG_FILE" ]; then
  echo "$FALLBACK_EMBLEM"
  exit 0
fi

# Parse rotation array from config
ROTATION=$(jq -r '.rotation[]' "$CONFIG_FILE" 2>/dev/null || echo "")

# If rotation is empty, use fallback
if [ -z "$ROTATION" ]; then
  echo "$FALLBACK_EMBLEM"
  exit 0
fi

# Convert rotation to array
EMBLEM_ARRAY=($ROTATION)
ARRAY_LENGTH=${#EMBLEM_ARRAY[@]}

# If array is empty, use fallback
if [ "$ARRAY_LENGTH" -eq 0 ]; then
  echo "$FALLBACK_EMBLEM"
  exit 0
fi

# Calculate current ISO week using UTC date
# Format: YYYY-Www (e.g., "2026-W04")
# %G = ISO year (handles year boundaries correctly)
# %V = ISO week number (01-53)
ISO_WEEK=$(date -u +%G-W%V)

# Generate deterministic seed from ISO week
# Use SHA256 hash to convert week string to numeric value
HASH=$(echo -n "$ISO_WEEK" | sha256sum | cut -c1-8)

# Convert hex hash to decimal
HASH_DECIMAL=$(printf '%d' "0x$HASH")

# Select emblem using modulo operation
INDEX=$((HASH_DECIMAL % ARRAY_LENGTH))
SELECTED_EMBLEM=${EMBLEM_ARRAY[$INDEX]}

echo "$SELECTED_EMBLEM"
