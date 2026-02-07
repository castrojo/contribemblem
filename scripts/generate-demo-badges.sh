#!/bin/bash
# Generate static demo badges for example users
# This script creates demonstration badges with sample stats

set -e

echo "Creating examples directory..."
mkdir -p examples

echo "Generating demo badges..."

# Sample stats for demonstration (not real data)
# Format: commits, pull_requests, issues, reviews, stars

# @castrojo - High activity across all areas
cat > /tmp/demo-stats-castrojo.json <<EOF
{
  "commits": 842,
  "pull_requests": 156,
  "issues": 89,
  "reviews": 234,
  "stars_received": 1247
}
EOF

# @jeefy - Balanced contributions
cat > /tmp/demo-stats-jeefy.json <<EOF
{
  "commits": 567,
  "pull_requests": 123,
  "issues": 67,
  "reviews": 189,
  "stars_received": 892
}
EOF

# @mrbobbytables - Heavy reviewer and issue creator
cat > /tmp/demo-stats-mrbobbytables.json <<EOF
{
  "commits": 423,
  "pull_requests": 98,
  "issues": 156,
  "reviews": 312,
  "stars_received": 634
}
EOF

# Note: This script requires actual emblem images to be present
# You can fetch different emblems using:
# contribemblem fetch-emblem <hash>

# Example emblem hashes:
# 1538938257 - Seventh Column Projection
# 1409726931 - Warlock Bond
# 2216440108 - Hunter Cloak

echo ""
echo "To complete demo badge generation:"
echo "1. Fetch emblem images for each user:"
echo "   contribemblem fetch-emblem 1538938257"
echo "   cp data/emblem.jpg examples/emblem-castrojo.jpg"
echo ""
echo "2. Generate badges:"
echo "   # For castrojo"
echo "   cp /tmp/demo-stats-castrojo.json data/stats.json"
echo "   contribemblem generate"
echo "   cp badge.png examples/castrojo.png"
echo ""
echo "3. Repeat for other users with different emblems and stats"
echo ""
echo "Demo stats files created in /tmp/demo-stats-*.json"
