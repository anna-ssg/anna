#!/usr/bin/env bash
set -euo pipefail
REPO="anna-ssg/anna"
ARCH="Linux_x86_64"

echo "Fetching latest release info for $REPO..."
LATEST_TAG=$(
  curl -s https://api.github.com/repos/$REPO/releases/latest \
  | sed -n 's/.*"tag_name":[[:space:]]*"\([^"]*\)".*/\1/p'
)

if [[ -z "$LATEST_TAG" ]]; then
  echo "Failed to determine latest release tag"
  exit 1
fi

echo "Latest release: $LATEST_TAG"
TARBALL="anna_${ARCH}.tar.gz"
URL="https://github.com/$REPO/releases/download/$LATEST_TAG/$TARBALL"

echo "Downloading $URL..."
curl -L "$URL" | tar -xz

if [[ ! -f anna ]]; then
  echo "anna binary not found after extraction"
  exit 1
fi

chmod +x anna
echo "Running anna..."
./anna
