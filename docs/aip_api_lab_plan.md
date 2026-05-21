# AIP Food Lookup API lab modernization plan

This document records the target production shape for AIP Food Lookup.

## Goals

- Serve the public API at `api.hashimojoe.com`.
- Run the Go API in Joe's lab on host port `8084`.
- Put Cloudflare Worker in front of the mobile API so the mobile app never contains the internal gateway secret.
- Route the Worker to the lab through Cloudflare Tunnel and host-installed Caddy.
- Use Docker Compose rendered from YAML by `babalu_yaml_env`.
- Keep file-based food data and runtime suggestion/feedback files; do not add a database.
- Send feedback to Slack, with local JSONL fallback.
- Add API access/error logs, logrotate, CORS, request body limits, gateway-secret checks, and fixed-window rate limits.

## Target topology

```text
Flutter Android app
  -> https://api.hashimojoe.com
  -> Cloudflare Worker
  -> Cloudflare Tunnel origin hostname
  -> host Caddy
  -> http://127.0.0.1:8084
  -> Docker container running aip-food-lookup-api
  -> /srv/stacks/aip-food-lookup/data mounted into the container
```

Future web presence:

```text
https://hashimojoe.com
  -> static site or Pages app
  -> https://api.hashimojoe.com for lookups
```

## Mobile gateway key rule

Do not ship `AIP_GATEWAY_SECRET` in Flutter. Mobile apps are inspectable, so the internal key belongs only in
Cloudflare Worker configuration and on the API host.

The mobile app may send public diagnostic headers such as:

- `X-AIP-Client`
- `X-AIP-App-Version`

Those headers are for logs and troubleshooting, not security.

## Production config

Use `scripts/aip/config.example.yaml` as the template. Production config lives on the host at:

```text
/srv/stacks/aip-food-lookup/api/config.yaml
```

The current Go API reads flattened environment variables such as:

- `AIP__API__DataFolder`
- `AIP__API__GatewaySecret`
- `AIP__API__SlackFeedbackWebhookUrl`
- `AIP__API__RateLimit__Enabled`

## Acceptance checks

- `go test ./...` passes in `cmd/aip_food_lookup`.
- Flutter tests pass in `flutter-app`.
- Docker Compose renders from YAML.
- Caddy routes the origin hostname to `127.0.0.1:8084`.
- Protected API routes return `401` without `X-Internal-Api-Key` in production.
- Worker-injected requests reach the API.
- Feedback reaches Slack or writes one JSON object per line to `feedback.jsonl`.
