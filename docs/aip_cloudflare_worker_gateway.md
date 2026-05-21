# AIP Cloudflare Worker gateway

The public mobile API should be `api.hashimojoe.com`, backed by a Cloudflare Worker. The Worker injects the private
gateway header before forwarding to the lab origin.

## Hostname pattern

Recommended public and origin hostnames:

```text
api.hashimojoe.com          public Worker route used by Flutter and future web clients
aip-origin.hashimojoe.com   Cloudflare Tunnel hostname that reaches host Caddy
```

The origin hostname should not be used by clients. If it is reached directly, the Go API should reject protected routes
because the gateway secret header is missing.

## Worker secrets

Configure these in Cloudflare, not in Git:

```text
AIP_GATEWAY_SECRET
AIP_ORIGIN_BASE_URL=https://aip-origin.hashimojoe.com
```

## Worker sketch

```js
export default {
  async fetch(request, env) {
    const url = new URL(request.url);
    const origin = new URL(env.AIP_ORIGIN_BASE_URL);
    url.protocol = origin.protocol;
    url.hostname = origin.hostname;
    url.port = origin.port;

    const headers = new Headers(request.headers);
    headers.delete("X-Internal-Api-Key");
    headers.set("X-Internal-Api-Key", env.AIP_GATEWAY_SECRET);
    headers.set("X-Forwarded-Host", "api.hashimojoe.com");

    return fetch(new Request(url.toString(), {
      method: request.method,
      headers,
      body: request.body,
      redirect: "manual",
    }));
  },
};
```

If a future static site at `hashimojoe.com` calls this API from a browser, keep the Go API `AllowedOrigins` set to the
site origins.
