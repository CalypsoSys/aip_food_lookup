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
Production feedback and suggestions post to Slack when `AIP__API__SlackFeedbackWebhookUrl` is configured. Feedback falls
back to `data/feedback.jsonl` if Slack is unavailable; suggestions succeed when either the local suggestion file write or
Slack delivery succeeds.

Run locally:

```powershell
cd cmd\aip_food_lookup
$env:AIP_DATA_FOLDER='..\..\data'
go run .
```

## Lab deployment target

The planned client-facing API endpoint is:

```text
https://hashimojoe.com/api
```

Recommended request path:

```text
Flutter Android app or future hashimojoe.com site
  -> Cloudflare Pages Function at hashimojoe.com/api/*
  -> Cloudflare Tunnel origin hostname
  -> host Caddy
  -> 127.0.0.1:8084
  -> Docker container running the Go API
```

The Flutter app must not contain the internal gateway secret. The Pages Function injects `X-Internal-Api-Key`, and the
Go API requires it for protected routes in production.

Deployment/config docs live in `docs/`, with YAML config in `scripts/aip/config.example.yaml`.

## Search coverage analyzer

The search coverage tool uses a two-step workflow: extract search terms on the production server, then compare them with the local catalog. It reads rotated and gzip-compressed API access logs without contacting the production API.

### Build and install the server tool

From the repository root in WSL, build the Linux binary:

```bash
bash scripts/aip/build-search-coverage.sh
```

Copy it to the server and install it in the utilities directory:

```bash
scp bin/search_coverage joe@YOUR_SERVER:/tmp/search_coverage
ssh joe@YOUR_SERVER
sudo cp /tmp/search_coverage /srv/utilities/bin/search_coverage
sudo chmod +x /srv/utilities/bin/search_coverage
```

### Extract searches on the server

Run this on the server:

```bash
/srv/utilities/bin/search_coverage extract \
  --logs /srv/logs/aip-food-lookup/api \
  --output /tmp/aip-searches.tsv
```

The output is a plain TSV file containing unique search keys, request counts, timestamps, and HTTP status counts. Copy it back to the repository's ignored `output` directory:

```powershell
scp joe@YOUR_SERVER:/tmp/aip-searches.tsv .\output\aip-searches.tsv
```

### Compare searches with the local catalog

From the repository root on Windows:

```powershell
go run .\cmd\search_coverage check `
  --input .\output\aip-searches.tsv `
  --catalog .\data
```

The report lists the total search keys, covered keys, uncovered keys, and details for each uncovered search. The downloaded export and local build artifacts are ignored by Git.
