# AIP Food Lookup production runbook

This runbook is the operational document for deploying, validating, rolling back, and maintaining the lab-hosted AIP
Food Lookup Go API stack.

Related docs:

- [aip_api_lab_plan.md](aip_api_lab_plan.md)
- [aip_ubuntu_host_preparation.md](aip_ubuntu_host_preparation.md)
- [aip_cloudflare_pages_gateway.md](aip_cloudflare_pages_gateway.md)
- [aip_caddy_host_setup.md](aip_caddy_host_setup.md)
- [aip_local_development.md](aip_local_development.md)

## Steady-state topology

Public client surface on Cloudflare:

- `hashimojoe.com`
- `https://hashimojoe.com/api` for Flutter production builds

Private lab origin behind Cloudflare Tunnel:

- `aip.hashimojoe.com`

Recommended request path:

- browser or Flutter client -> `https://hashimojoe.com/api`
- `/api/*` -> Cloudflare Pages Functions
- Pages Functions -> Cloudflare Tunnel hostname for the API origin
- Cloudflare Tunnel -> host-installed Caddy on the Ubuntu host
- Caddy -> `aip-food-lookup-api`

At minimum, host Caddy should include:

```caddy
{
    auto_https off
}

http://aip.hashimojoe.com {
    reverse_proxy 127.0.0.1:8084
}
```

## Server layout

Expected structure:

```text
/srv/stacks/aip-food-lookup/api
  docker-compose.yml
  config.yaml
  aip-food-lookup-api-latest.tar.gz
  scripts/
    compose-aip.sh
    smoke-local-docker.sh
    render-config-env
    aip.logrotate
    caddy.logrotate

/srv/stacks/aip-food-lookup/data
/srv/backups/aip-food-lookup
/srv/logs/aip-food-lookup/api
/srv/logs/caddy
```

Create the required directories if they do not already exist:

```bash
sudo mkdir -p /srv/stacks/aip-food-lookup/api/scripts
sudo mkdir -p /srv/stacks/aip-food-lookup/data
sudo mkdir -p /srv/backups/aip-food-lookup
sudo mkdir -p /srv/logs/aip-food-lookup/api
sudo mkdir -p /srv/logs/caddy

sudo chown -R $USER:$USER /srv/stacks/aip-food-lookup
sudo chown -R $USER:$USER /srv/backups/aip-food-lookup
sudo chown -R $USER:$USER /srv/logs/aip-food-lookup
sudo chown -R caddy:caddy /srv/logs/caddy
```

The API container mounts:

```text
/srv/stacks/aip-food-lookup/data -> /app/data
/srv/logs/aip-food-lookup/api   -> /app/logs
```

## Files from this repo

Copy or derive these from the repo:

- `docker/Dockerfile`
- `docker/docker-compose.yml`
- `scripts/aip/compose-aip.sh`
- `scripts/aip/smoke-local-docker.sh`
- `scripts/aip/config.example.yaml`
- `scripts/aip/aip.logrotate`
- `scripts/caddy/caddy.logrotate`
- the built API image tarball you create locally

Server-local files that must not come from git:

- `/srv/stacks/aip-food-lookup/api/config.yaml`
- Cloudflare Tunnel credentials
- real secrets and webhook URLs

## Required secret inputs

Keep real values in the host shell environment, a password manager, or the server-local `config.yaml`. Do not commit
them.

| Name | Purpose |
| --- | --- |
| `AIP_GATEWAY_SECRET` | Internal API key injected by Cloudflare Pages Functions |
| `AIP_SLACK_FEEDBACK_WEBHOOK_URL` | Slack webhook for feedback and suggestions |

For the Cloudflare Pages project, configure these environment bindings in Cloudflare, not in Git:

| Name | Purpose |
| --- | --- |
| `AIP_ORIGIN_BASE_URL` | Tunnel/Caddy origin URL, for example `https://aip.hashimojoe.com` |
| `AIP_GATEWAY_SECRET` | Same internal key configured on the Go API host |

## Build the API image locally

From the repo root in WSL/Linux:

```bash
mkdir -p /mnt/c/transfer
if [ -f /mnt/c/transfer/aip-food-lookup-api-latest.tar.gz ]; then mv /mnt/c/transfer/aip-food-lookup-api-latest.tar.gz /mnt/c/transfer/aip-food-lookup-api-latest.lastgood.tar.gz; fi
TRANSFER_DIR=/mnt/c/transfer scripts/build.sh
```

That leaves:

```text
C:\transfer\aip-food-lookup-api-latest.tar.gz
```

## Build the shared YAML-to-env renderer

Build the shared renderer from its repo in WSL/Linux so the server receives a Linux binary:

```bash
cd ~/work/calypsosys-workbench/repos/babalu-yaml-env
mkdir -p /mnt/c/transfer
if [ -f /mnt/c/transfer/render-config-env ]; then mv /mnt/c/transfer/render-config-env /mnt/c/transfer/render-config-env.lastgood; fi
go build -o /mnt/c/transfer/render-config-env ./cmd/babalu_yaml_env
```

That gives you:

```text
C:\transfer\render-config-env
```

## Create the server config.yaml

Before any stack command, create the server-local `config.yaml` using `scripts/aip/config.example.yaml` as the
reference, then fill in the real production values. Use `${VARIABLE_NAME}` for secrets so the YAML remains the single
source of truth while secrets still come from the host environment at runtime.

Then run on the Ubuntu host:

```bash
cd /srv/stacks/aip-food-lookup/api
vi config.yaml
chmod 600 config.yaml
```

Minimum structure:

```yaml
AIP_API_IMAGE: aip-food-lookup-api:latest
AIP_API_HOST_BIND: 127.0.0.1
AIP_API_HOST_PORT: 8084
AIP_DATA_HOST_PATH: /srv/stacks/aip-food-lookup/data
AIP_LOGS_HOST_PATH: /srv/logs/aip-food-lookup/api

AIP:
  API:
    ListenAddress: :8080
    DataFolder: /app/data
    AccessLogPath: /app/logs/access.log
    ErrorLogPath: /app/logs/errors.log
    AllowedOrigins:
      - https://hashimojoe.com
      - https://www.hashimojoe.com
    RequireGatewaySecret: true
    GatewaySecretHeaderName: X-Internal-Api-Key
    GatewaySecret: ${AIP_GATEWAY_SECRET}
    SlackFeedbackWebhookUrl: ${AIP_SLACK_FEEDBACK_WEBHOOK_URL}
    FeedbackJSONLPath: /app/data/feedback.jsonl
    RequestBodyLimitBytes: 32768
    RateLimit:
      Enabled: true
      SearchPermitLimit: 300
      WritePermitLimit: 60
      FeedbackPermitLimit: 10
      WindowSeconds: 60
```

Notes:

- food data and runtime feedback files persist in `/srv/stacks/aip-food-lookup/data`
- access/error logs are written under `/srv/logs/aip-food-lookup/api`
- leave `SlackFeedbackWebhookUrl` empty only when Slack feedback and suggestion delivery is intentionally disabled

## Stage and copy artifacts to the server

The AIP repo lives in WSL/Linux for this workflow, while Windows PowerShell owns the SSH key context. Stage the repo
files that PowerShell must copy into a Windows-visible transfer subfolder from WSL:

```bash
mkdir -p /mnt/c/transfer/aip-deploy

cp docker/docker-compose.yml /mnt/c/transfer/aip-deploy/docker-compose.yml
cp scripts/aip/compose-aip.sh /mnt/c/transfer/aip-deploy/compose-aip.sh
cp scripts/aip/smoke-local-docker.sh /mnt/c/transfer/aip-deploy/smoke-local-docker.sh
cp scripts/aip/aip.logrotate /mnt/c/transfer/aip-deploy/aip.logrotate
cp scripts/caddy/caddy.logrotate /mnt/c/transfer/aip-deploy/caddy.logrotate
```

Then copy from Windows PowerShell:

```powershell
$server = "replace_with_user@replace_with_server"
$transfer = "C:\transfer\aip-deploy"

ssh ${server} "cd /srv/stacks/aip-food-lookup/api && if [ -f aip-food-lookup-api-latest.tar.gz ]; then mv aip-food-lookup-api-latest.tar.gz aip-food-lookup-api-latest.lastgood.tar.gz; fi"
scp C:\transfer\aip-food-lookup-api-latest.tar.gz ${server}:/srv/stacks/aip-food-lookup/api/
scp C:\transfer\render-config-env ${server}:/srv/stacks/aip-food-lookup/api/scripts/render-config-env
scp "$transfer\docker-compose.yml" ${server}:/srv/stacks/aip-food-lookup/api/docker-compose.yml
scp "$transfer\compose-aip.sh" ${server}:/srv/stacks/aip-food-lookup/api/scripts/compose-aip.sh
scp "$transfer\smoke-local-docker.sh" ${server}:/srv/stacks/aip-food-lookup/api/scripts/smoke-local-docker.sh
scp "$transfer\aip.logrotate" ${server}:/srv/stacks/aip-food-lookup/api/scripts/aip.logrotate
scp "$transfer\caddy.logrotate" ${server}:/srv/stacks/aip-food-lookup/api/scripts/caddy.logrotate
```

After copying artifacts and editing `config.yaml`, on the Ubuntu host:

```bash
chmod +x /srv/stacks/aip-food-lookup/api/scripts/compose-aip.sh
chmod +x /srv/stacks/aip-food-lookup/api/scripts/smoke-local-docker.sh
chmod +x /srv/stacks/aip-food-lookup/api/scripts/render-config-env
chmod 600 /srv/stacks/aip-food-lookup/api/config.yaml
```

Install logrotate policies:

```bash
sudo cp /srv/stacks/aip-food-lookup/api/scripts/aip.logrotate /etc/logrotate.d/aip-food-lookup-api
sudo cp /srv/stacks/aip-food-lookup/api/scripts/caddy.logrotate /etc/logrotate.d/caddy
sudo chmod 644 /etc/logrotate.d/aip-food-lookup-api /etc/logrotate.d/caddy
sudo logrotate -d /etc/logrotate.d/aip-food-lookup-api
sudo logrotate -d /etc/logrotate.d/caddy
```

## Preflight checks on the server

Run on the Ubuntu host:

```bash
cd /srv/stacks/aip-food-lookup/api
docker version
docker compose version
test -f config.yaml && echo "config.yaml present"
test -f docker-compose.yml && echo "compose file present"
test -x scripts/compose-aip.sh && echo "compose wrapper present"
test -x scripts/smoke-local-docker.sh && echo "smoke test present"
test -x scripts/render-config-env && echo "render binary present"
test -d /srv/stacks/aip-food-lookup/data && echo "data directory present"
sudo caddy validate --config /etc/caddy/Caddyfile
systemctl status caddy --no-pager
systemctl is-active --quiet cloudflared && echo "cloudflared running"
```

If this reports that `cloudflared` is not running, or if a direct status check reports
`Unit cloudflared.service could not be found`, complete the Cloudflare Tunnel service setup in
[aip_ubuntu_host_preparation.md](aip_ubuntu_host_preparation.md) before continuing. Avoid pasting cloudflared status
output into tickets or docs because it can expose the tunnel token in the process arguments.

Validate the rendered compose config:

```bash
export AIP_GATEWAY_SECRET=replace_me
export AIP_SLACK_FEEDBACK_WEBHOOK_URL=https://hooks.slack.com/services/replace/me
./scripts/compose-aip.sh config >/tmp/aip-compose.out
tail -n 30 /tmp/aip-compose.out
```

Use real values in the active shell session before running the wrapper for deployment. This `config` command is only a
sanity check that `config.yaml` and the exported shell variables render into a valid Docker Compose config; it does not
start the stack.

## Load the API image

On the Ubuntu host:

```bash
cd /srv/stacks/aip-food-lookup/api
gzip -dc aip-food-lookup-api-latest.tar.gz | docker load
```

## Bring up the stack

On the Ubuntu host:

```bash
cd /srv/stacks/aip-food-lookup/api
./scripts/compose-aip.sh up -d
./scripts/compose-aip.sh ps
./scripts/compose-aip.sh logs aip-food-lookup-api --tail=100
```

Check the API directly on the host:

```bash
gatewaySecret="replace_with_real_gateway_secret"
curl -i http://127.0.0.1:8084/
curl -i "http://127.0.0.1:8084/search?key=apple"
curl -i -H "X-Internal-Api-Key: ${gatewaySecret}" "http://127.0.0.1:8084/search?key=apple"
```

If gateway-secret enforcement is enabled, direct protected requests without `X-Internal-Api-Key` should return `401`.
The keyed request should return JSON results.

Check the Caddy path:

```bash
curl -i -H "Host: aip.hashimojoe.com" http://127.0.0.1:80/
curl -i -H "Host: aip.hashimojoe.com" -H "X-Internal-Api-Key: ${gatewaySecret}" "http://127.0.0.1:80/search?key=apple"
```

If either request returns `500` and `/srv/logs/aip-food-lookup/api/errors.log` does not exist, check the container logs:

```bash
./scripts/compose-aip.sh logs aip-food-lookup-api --tail=200
```

## Run the smoke test

On the Ubuntu host:

```bash
cd /srv/stacks/aip-food-lookup/api
AIP_GATEWAY_SECRET="${gatewaySecret}" ./scripts/smoke-local-docker.sh
```

The smoke test starts Docker Compose and verifies:

- `GET /` returns `200`
- unkeyed `GET /search?key=apple` returns `401`
- keyed `GET /search?key=apple` returns `200`
- keyed `POST /feedback` returns `200`

## Cloudflare Tunnel

Recommended ingress shape:

```yaml
tunnel: <your-tunnel-id>
credentials-file: /etc/cloudflared/<your-tunnel-id>.json

ingress:
  - hostname: aip.hashimojoe.com
    service: http://127.0.0.1:80
  - service: http_status:404
```

Run or restart the tunnel:

```bash
sudo systemctl restart cloudflared
systemctl is-active --quiet cloudflared && echo "cloudflared running"
```

## Cloudflare Pages project

### `hashimojoe.com`

- Root directory: `frontend`
- Build command: `pnpm install --frozen-lockfile && pnpm run build`
- Output: `dist`
- Custom domains:
  - `hashimojoe.com`
  - `www.hashimojoe.com`
- Environment bindings:
  - `AIP_ORIGIN_BASE_URL=https://aip.hashimojoe.com`
  - `AIP_GATEWAY_SECRET=<same as AIP__API__GatewaySecret>`
  - `VITE_AIP_API_BASE_URL=/api`

## Validate public paths

From a workstation:

```bash
curl -i https://hashimojoe.com/
curl -i https://hashimojoe.com/api/
curl -i "https://hashimojoe.com/api/search?key=apple"
```

Protected public API calls should pass through the Pages Function, which injects `X-Internal-Api-Key` server-side.

Flutter production builds should use:

```text
https://hashimojoe.com/api
```

## Logs and maintenance

On the Ubuntu host:

```bash
cd /srv/stacks/aip-food-lookup/api
./scripts/compose-aip.sh ps
./scripts/compose-aip.sh logs aip-food-lookup-api --tail=100
docker ps --filter name=aip-food-lookup-api

ls -l /srv/logs/aip-food-lookup/api
tail -n 50 /srv/logs/aip-food-lookup/api/access.log
tail -n 50 /srv/logs/aip-food-lookup/api/errors.log
tail -n 50 /srv/logs/caddy/caddy.log
```

Restart the stack:

```bash
cd /srv/stacks/aip-food-lookup/api
./scripts/compose-aip.sh restart
```

Roll back to the last good image tarball if you staged one:

```bash
cd /srv/stacks/aip-food-lookup/api
gzip -dc aip-food-lookup-api-latest.lastgood.tar.gz | docker load
./scripts/compose-aip.sh up -d
```

## Backup

Back up file-backed data from the Ubuntu host:

```bash
sudo tar -czf /srv/backups/aip-food-lookup/aip-data-$(date +%Y%m%d).tar.gz -C /srv/stacks/aip-food-lookup data
```

The backup includes runtime feedback JSONL, suggestion files, and any file-backed data updates under
`/srv/stacks/aip-food-lookup/data`.
