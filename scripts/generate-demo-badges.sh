#!/bin/bash
# Generate static demo badges for example users
# This script creates demonstration badges with sample stats and different emblems

set -e

# Require BUNGIE_API_KEY
if [ -z "$BUNGIE_API_KEY" ]; then
  echo "Error: BUNGIE_API_KEY environment variable not set"
  echo "Get your API key from https://www.bungie.net/en/Application"
  exit 1
fi

echo "Creating examples directory..."
mkdir -p examples data

echo "Generating demo badges..."

# Define users with their stats and unique emblem hashes
declare -A users=(
  ["castrojo"]="1538938257"  # Seventh Column Projection
  ["jeefy"]="2962058744"     # Different emblem from rotation
  ["mrbobbytables"]="2962058745"  # Another emblem from rotation
)

declare -A commits=(["castrojo"]=842 ["jeefy"]=567 ["mrbobbytables"]=423)
declare -A prs=(["castrojo"]=156 ["jeefy"]=123 ["mrbobbytables"]=98)
declare -A issues=(["castrojo"]=89 ["jeefy"]=67 ["mrbobbytables"]=156)
declare -A reviews=(["castrojo"]=234 ["jeefy"]=189 ["mrbobbytables"]=312)
declare -A stars=(["castrojo"]=1247 ["jeefy"]=892 ["mrbobbytables"]=634)

for user in castrojo jeefy mrbobbytables; do
  echo ""
  echo "=== Generating badge for @${user} ==="
  
  emblem_hash="${users[$user]}"
  
  # Create stats file
  cat > data/stats.json <<EOF
{
  "year": 2026,
  "updated_at": "2026-02-07T00:00:00Z",
  "commits": ${commits[$user]},
  "pull_requests": ${prs[$user]},
  "issues": ${issues[$user]},
  "reviews": ${reviews[$user]},
  "stars_received": ${stars[$user]}
}
EOF
  
  # Delete cached emblem to force fresh fetch
  rm -f data/emblem.jpg
  
  # Fetch emblem from Bungie API
  echo "Fetching emblem ${emblem_hash}..."
  export GITHUB_ACTOR="$user"
  ./contribemblem fetch-emblem "$emblem_hash"
  
  # Generate badge
  echo "Generating badge..."
  ./contribemblem generate
  
  # Copy to examples
  cp badge.png "examples/${user}.png"
  echo "✓ Badge saved to examples/${user}.png"
done

echo ""
echo "✨ All demo badges generated successfully!"
echo "   examples/castrojo.png"
echo "   examples/jeefy.png"
echo "   examples/mrbobbytables.png"
