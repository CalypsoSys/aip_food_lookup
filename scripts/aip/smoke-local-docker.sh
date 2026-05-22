#!/usr/bin/env bash

set -euo pipefail

SCRIPT_DIR="$(cd "$(dirname "${BASH_SOURCE[0]}")" && pwd)"
REPO_ROOT="$(cd "$SCRIPT_DIR/../.." && pwd)"

if [ -f "$SCRIPT_DIR/../docker-compose.yml" ]; then
    DEFAULT_CONFIG_FILE="$SCRIPT_DIR/../config.yaml"
else
    DEFAULT_CONFIG_FILE="$SCRIPT_DIR/config.local.yaml"
fi

CONFIG_FILE="${CONFIG_FILE:-$DEFAULT_CONFIG_FILE}"
HOST="${AIP_SMOKE_HOST:-127.0.0.1}"
PORT="${AIP_SMOKE_PORT:-8084}"
SECRET="${AIP_GATEWAY_SECRET:-local-smoke-secret}"
BASE_URL="http://$HOST:$PORT"

export AIP_GATEWAY_SECRET="$SECRET"
export AIP_SLACK_FEEDBACK_WEBHOOK_URL="${AIP_SLACK_FEEDBACK_WEBHOOK_URL:-}"
export AIP_REPO_ROOT="${AIP_REPO_ROOT:-$REPO_ROOT}"
export CONFIG_FILE

require_command() {
    if ! command -v "$1" >/dev/null 2>&1; then
        echo "Required command not found: $1" >&2
        exit 1
    fi
}

expect_status() {
    local expected="$1"
    local url="$2"
    shift 2
    local status
    status="$(curl -sS -o /tmp/aip-smoke-response.txt -w "%{http_code}" "$@" "$url")"
    if [ "$status" != "$expected" ]; then
        echo "Expected HTTP $expected from $url, got $status" >&2
        cat /tmp/aip-smoke-response.txt >&2
        exit 1
    fi
}

require_command curl
require_command docker

if [ ! -f "$CONFIG_FILE" ]; then
    echo "Config file not found: $CONFIG_FILE" >&2
    echo "Copy scripts/aip/config.example.yaml to scripts/aip/config.local.yaml first." >&2
    exit 1
fi

"$SCRIPT_DIR/compose-aip.sh" up -d

echo "Checking health at $BASE_URL/"
expect_status 200 "$BASE_URL/"

echo "Checking protected search rejects missing gateway key"
expect_status 401 "$BASE_URL/search?key=apple"

echo "Checking protected search accepts gateway key"
expect_status 200 "$BASE_URL/search?key=apple" -H "X-Internal-Api-Key: $SECRET"

echo "Checking suggestion endpoint accepts gateway key"
expect_status 200 "$BASE_URL/suggest" \
    -H "X-Internal-Api-Key: $SECRET" \
    -H "Content-Type: application/json" \
    --data '{"inputText":"smoke-test-food-unique","allowed":true}'

echo "Checking feedback endpoint accepts gateway key"
expect_status 200 "$BASE_URL/feedback" \
    -H "X-Internal-Api-Key: $SECRET" \
    -H "Content-Type: application/json" \
    --data '{"name":"Smoke Test","email":"","subject":"Local smoke","message":"Local Docker smoke test","source":"smoke"}'

echo "AIP local Docker smoke test passed."
