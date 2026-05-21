#!/usr/bin/env bash

set -euo pipefail

REPO_ROOT="$(cd "$(dirname "${BASH_SOURCE[0]}")/.." && pwd)"
IMAGE_NAME="${AIP_API_IMAGE:-aip-food-lookup-api:latest}"
TRANSFER_DIR="${TRANSFER_DIR:-/mnt/c/transfer}"
TAR_PATH="$REPO_ROOT/docker/aip-food-lookup-api-latest.tar"

echo "Building aip_food_lookup API binary"
cd "$REPO_ROOT/cmd/aip_food_lookup"
CGO_ENABLED=0 GOOS=linux GOARCH=amd64 go build -a -installsuffix cgo -o "$REPO_ROOT/docker/aip_food_lookup" .

echo "Building Docker image $IMAGE_NAME"
cd "$REPO_ROOT"
docker build -t "$IMAGE_NAME" -f docker/Dockerfile docker

docker image prune -f --filter "dangling=true"
docker save "$IMAGE_NAME" -o "$TAR_PATH"
gzip -f "$TAR_PATH"

mkdir -p "$TRANSFER_DIR"
cp "$TAR_PATH.gz" "$TRANSFER_DIR/"
echo "Wrote $TAR_PATH.gz"
