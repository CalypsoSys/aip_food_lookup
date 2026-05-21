# AIP Food Lookup production runbook

This runbook covers deploying and validating the lab-hosted Go API.

Related docs:

- `docs/aip_api_lab_plan.md`
- `docs/aip_cloudflare_worker_gateway.md`
- `docs/aip_caddy_host_setup.md`
- `docs/aip_local_development.md`

## Server layout

```text
/srv/stacks/aip-food-lookup/api
  docker-compose.yml
  config.yaml
  aip-food-lookup-api-latest.tar.gz
  scripts/
    compose-aip.sh
    render-config-env
    aip.logrotate

/srv/stacks/aip-food-lookup/data
/srv/backups/aip-food-lookup
/srv/logs/aip-food-lookup/api
/srv/logs/caddy
```

## Host setup

```bash
sudo apt update
sudo apt install -y curl ca-certificates docker.io docker-compose-plugin caddy logrotate

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

## Config

Create host config from `scripts/aip/config.example.yaml`:

```bash
cd /srv/stacks/aip-food-lookup/api
vi config.yaml
chmod 600 config.yaml
```

Required secret inputs:

| Name | Purpose |
| --- | --- |
| `AIP_GATEWAY_SECRET` | Internal key injected by Cloudflare Worker |
| `AIP_SLACK_FEEDBACK_WEBHOOK_URL` | Slack webhook for feedback |

## Build image

From the repo root:

```bash
scripts/build.sh
```

This writes:

```text
docker/aip-food-lookup-api-latest.tar.gz
```

## Deploy

Copy these files to `/srv/stacks/aip-food-lookup/api`:

- `docker/docker-compose.yml`
- `scripts/aip/compose-aip.sh`
- `scripts/aip/aip.logrotate`
- `docker/aip-food-lookup-api-latest.tar.gz`
- the built `render-config-env` binary from `repos/babalu-yaml-env`

Then run:

```bash
cd /srv/stacks/aip-food-lookup/api
gzip -dc aip-food-lookup-api-latest.tar.gz | docker load
chmod +x scripts/compose-aip.sh scripts/render-config-env
scripts/compose-aip.sh config
scripts/compose-aip.sh up -d
scripts/compose-aip.sh ps
```

## Validate

```bash
docker ps --filter name=aip-food-lookup-api
docker logs --tail 100 aip-food-lookup-api
curl -i http://127.0.0.1:8084/
curl -i "http://127.0.0.1:8084/search?key=apple"
curl -i -H "X-Internal-Api-Key: $AIP_GATEWAY_SECRET" "http://127.0.0.1:8084/search?key=apple"
```

The unkeyed protected-route request should return `401` when `RequireGatewaySecret` is enabled. The keyed request should
return JSON results.

Check logs:

```bash
ls -l /srv/logs/aip-food-lookup/api
tail -n 50 /srv/logs/aip-food-lookup/api/access.log
tail -n 50 /srv/logs/aip-food-lookup/api/errors.log
```

## Logrotate

```bash
export AIP_REPO_ROOT=/absolute/path/to/aip_food_lookup
sudo cp "$AIP_REPO_ROOT/scripts/aip/aip.logrotate" /etc/logrotate.d/aip-food-lookup-api
sudo cp "$AIP_REPO_ROOT/scripts/caddy/caddy.logrotate" /etc/logrotate.d/aip-caddy
sudo chmod 644 /etc/logrotate.d/aip-food-lookup-api /etc/logrotate.d/aip-caddy
sudo logrotate -d /etc/logrotate.d/aip-food-lookup-api
sudo logrotate -d /etc/logrotate.d/aip-caddy
```

## Backup

```bash
sudo tar -czf /srv/backups/aip-food-lookup/aip-data-$(date +%Y%m%d).tar.gz -C /srv/stacks/aip-food-lookup data
```
