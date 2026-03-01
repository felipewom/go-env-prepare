#!/usr/bin/env bash
# update-formula.sh — Update formula.rb with a new version tag and sha256 checksum.
#
# Usage:
#   ./scripts/update-formula.sh <tag>   <sha256>
#   ./scripts/update-formula.sh v0.1.0  abc123...
#
# The script rewrites the `url` and `sha256` fields in formula.rb in-place.

set -euo pipefail

TAG="${1:?Usage: $0 <tag> <sha256>}"
SHA256="${2:?Usage: $0 <tag> <sha256>}"

REPO_ROOT="$(cd "$(dirname "${BASH_SOURCE[0]}")/.." && pwd)"
FORMULA="${REPO_ROOT}/formula.rb"

if [[ ! -f "${FORMULA}" ]]; then
  echo "❌ formula.rb not found at ${FORMULA}" >&2
  exit 1
fi

# Strip leading 'v' for bare version number if present
VERSION="${TAG#v}"

# Update url line
sed -i.bak \
  -e "s|refs/tags/[^\"]*|refs/tags/${TAG}|" \
  -e "s|sha256 \"[^\"]*\"|sha256 \"${SHA256}\"|" \
  "${FORMULA}"

rm -f "${FORMULA}.bak"

echo "✅ formula.rb updated → tag=${TAG}  version=${VERSION}  sha256=${SHA256}"
