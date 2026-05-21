# aip_food_lookup

Go API backend and Flutter migration workspace for AIP Food Lookup.

## Backend

The Go backend lives in `cmd/aip_food_lookup` and serves:

- `GET /search?key=<text>&type=<searchbytextandsound|searchbytext|searchbysound>`
- `POST /suggest`
- `POST /feedback`
- `GET /categories`
- `GET /subcategory?cat=<Allowed|Not Allowed>&sub=<subcategory>`

Food data is stored in `data/allowed` and `data/not_allowed`. Runtime suggestion and feedback files are ignored by git.
Production feedback posts to Slack when `AIP__API__SlackFeedbackWebhookUrl` is configured and falls back to
`data/feedback.jsonl` if Slack is unavailable.

Run locally:

```powershell
cd cmd\aip_food_lookup
$env:AIP_DATA_FOLDER='..\..\data'
go run .
```

## Lab deployment target

The planned lab API endpoint is:

```text
https://api.hashimojoe.com
```

Recommended request path:

```text
Flutter Android app or future hashimojoe.com site
  -> Cloudflare Worker at api.hashimojoe.com
  -> Cloudflare Tunnel origin hostname
  -> host Caddy
  -> 127.0.0.1:8084
  -> Docker container running the Go API
```

The Flutter app must not contain the internal gateway secret. The Worker injects `X-Internal-Api-Key`, and the Go API
requires it for protected routes in production.

Deployment/config docs live in `docs/`, with YAML config in `scripts/aip/config.example.yaml`.
