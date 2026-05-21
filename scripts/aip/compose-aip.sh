#!/usr/bin/env bash

set -euo pipefail

SCRIPT_DIR="$(cd "$(dirname "${BASH_SOURCE[0]}")" && pwd)"
DEPLOY_ROOT="$(cd "$SCRIPT_DIR/.." && pwd)"
REPO_ROOT="$(cd "$SCRIPT_DIR/../.." && pwd)"

if [ -f "$DEPLOY_ROOT/docker-compose.yml" ]; then
    DEFAULT_STACK_DIR="$DEPLOY_ROOT"
    DEFAULT_CONFIG_FILE="$DEPLOY_ROOT/config.yaml"
else
    DEFAULT_STACK_DIR="$REPO_ROOT/docker"
    DEFAULT_CONFIG_FILE="$SCRIPT_DIR/config.yaml"
fi

STACK_DIR="${STACK_DIR:-$DEFAULT_STACK_DIR}"
COMPOSE_FILE="${COMPOSE_FILE:-$STACK_DIR/docker-compose.yml}"
CONFIG_FILE="${CONFIG_FILE:-$DEFAULT_CONFIG_FILE}"
RENDER_BIN="$SCRIPT_DIR/render-config-env"

if [ ! -f "$COMPOSE_FILE" ]; then
    echo "Compose file not found: $COMPOSE_FILE" >&2
    exit 1
fi

if [ ! -f "$CONFIG_FILE" ]; then
    echo "Config file not found: $CONFIG_FILE" >&2
    exit 1
fi

if [ ! -x "$RENDER_BIN" ]; then
    echo "Render binary is not executable: $RENDER_BIN" >&2
    echo "Build and deploy render-config-env from ~/work/calypsosys-workbench/repos/babalu-yaml-env." >&2
    exit 1
fi

TEMP_ENV_FILE="$(mktemp)"
cleanup() {
    rm -f "$TEMP_ENV_FILE"
}
trap cleanup EXIT

"$RENDER_BIN" --format env "$CONFIG_FILE" > "$TEMP_ENV_FILE"

COMPOSE_UNSET_ARGS=()
while IFS='=' read -r key _; do
    if [ -n "$key" ]; then
        COMPOSE_UNSET_ARGS+=("-u" "$key")
    fi
done < "$TEMP_ENV_FILE"

env "${COMPOSE_UNSET_ARGS[@]}" docker compose -f "$COMPOSE_FILE" --env-file "$TEMP_ENV_FILE" "$@"
