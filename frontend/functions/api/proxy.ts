type PagesEnv = Record<string, string | undefined>;

export type PathParam = string | string[] | undefined;

export function gatewayHeaderName(): string {
  return ['X', 'Internal', 'Api', 'Key'].join('-');
}

export function readProxyEnv(env: PagesEnv): { originBaseUrl: string; gatewaySecret: string } {
  const originBaseUrl = env[['AIP', 'ORIGIN', 'BASE', 'URL'].join('_')];
  const gatewaySecret = env[['AIP', 'GATEWAY', 'SECRET'].join('_')];

  if (!originBaseUrl || !gatewaySecret) {
    throw new Error('Missing Cloudflare Pages API proxy configuration.');
  }

  return { originBaseUrl, gatewaySecret };
}

export function buildOriginUrl(requestUrl: string, originBaseUrl: string, pathParam: PathParam): string {
  const request = new URL(requestUrl);
  const origin = new URL(originBaseUrl);
  const basePath = origin.pathname.replace(/\/$/, '');
  const forwardedPath = normalizeForwardedPath(pathParam);

  origin.pathname = `${basePath}/${forwardedPath}`.replace(/\/+$/, '') || '/';
  origin.search = request.search;

  return origin.toString();
}

export function createProxyRequest(request: Request, env: PagesEnv, pathParam: PathParam): Request {
  const { originBaseUrl, gatewaySecret } = readProxyEnv(env);
  const headers = new Headers(request.headers);
  headers.delete(gatewayHeaderName());
  headers.set(gatewayHeaderName(), gatewaySecret);
  headers.set('X-Forwarded-Host', new URL(request.url).host);
  headers.set('X-Forwarded-Proto', new URL(request.url).protocol.replace(':', ''));

  const init: RequestInit & { duplex?: 'half' } = {
    method: request.method,
    headers,
    redirect: 'manual',
  };

  if (request.method !== 'GET' && request.method !== 'HEAD') {
    init.body = request.body;
    init.duplex = 'half';
  }

  return new Request(buildOriginUrl(request.url, originBaseUrl, pathParam), init);
}

function normalizeForwardedPath(pathParam: PathParam): string {
  if (Array.isArray(pathParam)) {
    return pathParam.join('/');
  }

  return pathParam ?? '';
}
