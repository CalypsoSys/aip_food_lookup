# AIP Food Lookup Frontend

Vue 3 + TypeScript + Vite starter for the future `hashimojoe.com` web presence.

The web app calls the API through `VITE_AIP_API_BASE_URL`, defaulting to `/api`.
Local Vite development proxies `/api/*` to the Go API at `http://127.0.0.1:8080`
and strips the `/api` prefix.

Production is intended for Cloudflare Pages. Keep the default same-origin `/api`
base path and configure the Pages Function proxy with server-side environment
bindings for the origin URL and gateway secret.

```text
VITE_AIP_API_BASE_URL=/api
```

Do not put internal gateway headers, server-side secrets, Slack URLs, or private
backend URLs in frontend code or browser-exposed environment files.

## Commands

```powershell
pnpm install --frozen-lockfile
pnpm run dev
pnpm run test
pnpm run build
```

Standard frontend verification:

```powershell
pnpm install --frozen-lockfile && pnpm run build
```

The app uses Vue Router with history URLs:

```text
/search
/categories
/about
/feedback
/diagnostics
```

Cloudflare Pages should publish the Vite build output from `dist/`. The
`public/_redirects` file keeps routed browser refreshes on `index.html`.

Run the Go API locally before starting Vite:

```powershell
cd ..\cmd\aip_food_lookup
$env:AIP_DATA_FOLDER='..\..\data'
go run .
```
