# AIP Caddy host setup

This document describes host-installed Caddy for AIP Food Lookup behind Cloudflare Tunnel.

## Goal

Cloudflare Pages Functions handle same-origin `/api/*` requests from `hashimojoe.com` and forward to the tunnel origin
hostname. Caddy receives the origin request and proxies it to the local Docker-published API port.

```text
hashimojoe.com Pages Function /api/*
  -> api.hashimojoe.com tunnel
  -> Caddy
  -> 127.0.0.1:8084
```

## Host log directory

```bash
sudo mkdir -p /srv/logs/caddy
sudo chown -R caddy:caddy /srv/logs/caddy
sudo chmod 755 /srv/logs/caddy
```

## Shared Host Caddyfile

The authoritative Caddyfile is host-owned, not repo-owned. Keep the deployable
`/etc/caddy/Caddyfile` on the server and use the shared workbench reference as the
starting point:

```text
CalypsoSys operations workbench:
  docs/caddy.md
  templates/caddy/calypsosys-host.Caddyfile.example
```

AIP Food Lookup needs this route in the shared host Caddyfile:

```text
api.hashimojoe.com -> 127.0.0.1:8084
```

## Validate

```bash
sudo vi /etc/caddy/Caddyfile
sudo caddy validate --config /etc/caddy/Caddyfile
sudo systemctl reload caddy
sudo systemctl status caddy --no-pager
```

Local checks:

```bash
curl -i -H "Host: api.hashimojoe.com" http://127.0.0.1:80/
curl -i -H "Host: api.hashimojoe.com" http://127.0.0.1:80/search?key=apple
curl -i -H "Host: api.hashimojoe.com" -H "X-Internal-Api-Key: $AIP_GATEWAY_SECRET" "http://127.0.0.1:80/search?key=apple"
```

The unkeyed protected-route request should return `401` in production.
