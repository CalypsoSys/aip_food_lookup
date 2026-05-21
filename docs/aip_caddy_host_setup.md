# AIP Caddy host setup

This document describes host-installed Caddy for AIP Food Lookup behind Cloudflare Tunnel.

## Goal

Cloudflare Worker handles `api.hashimojoe.com` and forwards to the tunnel origin hostname. Caddy receives the origin
request and proxies it to the local Docker-published API port.

```text
api.hashimojoe.com Worker
  -> aip-origin.hashimojoe.com tunnel
  -> Caddy
  -> 127.0.0.1:8084
```

## Host log directory

```bash
sudo mkdir -p /srv/logs/caddy
sudo chown -R caddy:caddy /srv/logs/caddy
sudo chmod 755 /srv/logs/caddy
```

## Caddyfile

For the tunnel pattern, Caddy can listen on plain HTTP locally:

```caddy
{
    auto_https off

    log {
        output file /srv/logs/caddy/caddy.log
        format console
    }
}

http://aip-origin.hashimojoe.com {
    reverse_proxy 127.0.0.1:8084
}
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
curl -i -H "Host: aip-origin.hashimojoe.com" http://127.0.0.1:80/
curl -i -H "Host: aip-origin.hashimojoe.com" http://127.0.0.1:80/search?key=apple
curl -i -H "Host: aip-origin.hashimojoe.com" -H "X-Internal-Api-Key: $AIP_GATEWAY_SECRET" "http://127.0.0.1:80/search?key=apple"
```

The unkeyed protected-route request should return `401` in production.
