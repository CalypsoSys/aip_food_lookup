#!/usr/bin/env bash

set -euo pipefail

REPO_ROOT="$(cd "$(dirname "${BASH_SOURCE[0]}")/../.." && pwd)"
OUTPUT_PATH="${SEARCH_COVERAGE_OUTPUT:-$REPO_ROOT/bin/search_coverage}"

mkdir -p "$(dirname "$OUTPUT_PATH")"

echo "Building search coverage tool"
cd "$REPO_ROOT"
CGO_ENABLED=0 GOOS="${GOOS:-linux}" GOARCH="${GOARCH:-amd64}" \
  go build -o "$OUTPUT_PATH" ./cmd/search_coverage

echo "Wrote $OUTPUT_PATH"
