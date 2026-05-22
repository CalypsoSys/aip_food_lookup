# AIP Cloudflare Pages API gateway

The web app should be deployed as a Cloudflare Pages app. Browser code calls same-origin `/api/*`, and Flutter
production builds use `https://hashimojoe.com/api`. A Pages Function injects the private gateway header before
forwarding to the lab origin.

## Hostname pattern

Recommended web and origin hostnames:

```text
hashimojoe.com              Cloudflare Pages site and /api gateway used by browser and Flutter clients
api.hashimojoe.com          Cloudflare Tunnel hostname that reaches host Caddy
```

The origin hostname should not be used by clients. If it is reached directly, the Go API should reject protected routes
because the gateway secret header is missing.

The web frontend should use the same-origin API path in production:

```text
VITE_AIP_API_BASE_URL=/api
```

Flutter production builds should use the absolute API base URL:

```text
https://hashimojoe.com/api
```

Local Vite development also uses `/api/*` and proxies to `http://127.0.0.1:8080`, so the internal gateway secret is
never compiled into browser code.

Flutter/mobile can still move to a separate public API hostname later if usage or client needs justify it, but phase 1
uses the Pages Function gateway on `hashimojoe.com`.

## Pages environment bindings

Configure these in Cloudflare Pages, not in Git:

```text
AIP_GATEWAY_SECRET
AIP_ORIGIN_BASE_URL=https://api.hashimojoe.com
```

## Pages Function implementation

The repo includes the server-side proxy at:

```text
frontend/functions/api/[[path]].ts
```

The Pages Function forwards `/api/search`, `/api/categories`, `/api/feedback`, and future `/api/*` paths to
`AIP_ORIGIN_BASE_URL`, preserving the method, query string, and request body. It deletes any client-supplied
`X-Internal-Api-Key` header, then injects the configured `AIP_GATEWAY_SECRET`.

Keep the Go API `AllowedOrigins` set to the deployed Pages origins.
