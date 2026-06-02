# AIP local development

## Go API

Run from the repo root:

```bash
cd cmd/aip_food_lookup
AIP_DATA_FOLDER=../../data go run .
```

The API listens on `:8080` by default. The local health check is:

```bash
curl -i http://127.0.0.1:8080/
```

For local gateway-secret testing:

```bash
cd cmd/aip_food_lookup
AIP_GATEWAY_SECRET="$(openssl rand -hex 24)"
AIP__API__RequireGatewaySecret=true AIP__API__GatewaySecret="$AIP_GATEWAY_SECRET" AIP_DATA_FOLDER=../../data go run .
```

Then:

```bash
curl -i "http://127.0.0.1:8080/search?key=apple"
curl -i -H "X-Internal-Api-Key: $AIP_GATEWAY_SECRET" "http://127.0.0.1:8080/search?key=apple"
```

## Flutter Android

The Android emulator can reach the host on `10.0.2.2`:

```bash
cd flutter-app
flutter pub get
flutter test
flutter run --dart-define=AIP_BACKEND_URL=http://10.0.2.2:8080 --dart-define=AIP_CLIENT_NAME=android --dart-define=AIP_APP_VERSION=dev
```

For production builds, point the app at:

```text
https://hashimojoe.com/api
```

Do not pass or compile `AIP_GATEWAY_SECRET` into Flutter. The Cloudflare Pages Function owns that internal key.

## Web frontend

The web starter lives in `frontend/` and uses Vue 3, TypeScript, Vite, Pinia, and Vitest.

Run the Go API on `127.0.0.1:8080`, then:

```bash
cd frontend
pnpm install --frozen-lockfile
pnpm run dev
```

Open:

```text
http://127.0.0.1:5173
```

Local web API calls use `/api/*`. Vite proxies those requests to `http://127.0.0.1:8080` and removes the `/api`
prefix. Production web builds should keep the same-origin default:

```text
VITE_AIP_API_BASE_URL=/api
```

Cloudflare Pages Functions handle the production `/api/*` proxy and inject the internal gateway key from server-side
Pages environment bindings. Do not put `X-Internal-Api-Key`, `AIP_GATEWAY_SECRET`, or Slack webhook values in frontend
code or browser-exposed frontend env files.

Standard frontend verification is:

```bash
pnpm install --frozen-lockfile && pnpm run build
```

VS Code includes a `Local: frontend + API` compound launch. It starts the Go API and Vite as hidden background tasks,
then opens the web frontend in a visible browser debug session.

## Docker config render

Install the shared renderer locally at:

```text
/srv/utilities/bin/render-config-env
```

The compose wrapper also supports a legacy fallback at `scripts/aip/render-config-env` during migration, or a custom
path through `RENDER_BIN`.

Copy the sample config:

```bash
cp scripts/aip/config.example.yaml scripts/aip/config.local.yaml
vi scripts/aip/config.local.yaml
```

Then test compose rendering:

```bash
scripts/aip/compose-aip.sh config
```

## Local Docker smoke test

After building/loading `aip-food-lookup-api:latest`, copy the sample config:

```bash
cp scripts/aip/config.example.yaml scripts/aip/config.local.yaml
```

For local smoke testing, keep these values in `scripts/aip/config.local.yaml` or provide them through your shell:

```yaml
AIP_API_HOST_BIND: 127.0.0.1
AIP_API_HOST_PORT: 8084
AIP_DATA_HOST_PATH: ${AIP_REPO_ROOT}/data
AIP_LOGS_HOST_PATH: ${AIP_REPO_ROOT}/logs/aip-api
```

Then run:

```bash
AIP_GATEWAY_SECRET=local-smoke-secret scripts/aip/smoke-local-docker.sh
```

The smoke test starts Docker Compose and verifies:

- `GET /` returns `200`
- unkeyed `GET /search?key=apple` returns `401`
- keyed `GET /search?key=apple` returns `200`
- keyed `POST /feedback` returns `200`
