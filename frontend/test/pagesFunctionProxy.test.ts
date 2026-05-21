import { describe, expect, it } from 'vitest';

import {
  buildOriginUrl,
  createProxyRequest,
  gatewayHeaderName,
  readProxyEnv,
} from '../functions/api/proxy';

const originEnvName = ['AIP', 'ORIGIN', 'BASE', 'URL'].join('_');
const secretEnvName = ['AIP', 'GATEWAY', 'SECRET'].join('_');

describe('Pages Function API proxy', () => {
  it('builds origin URLs from same-origin API paths', () => {
    const url = buildOriginUrl(
      'https://hashimojoe.com/api/search?key=apple&type=searchbytext',
      'https://aip-origin.example.com/base',
      'search',
    );

    expect(url).toBe('https://aip-origin.example.com/base/search?key=apple&type=searchbytext');
  });

  it('requires the origin URL and gateway secret bindings', () => {
    expect(() => readProxyEnv({ [originEnvName]: 'https://aip-origin.example.com' })).toThrow(
      'Missing Cloudflare Pages API proxy configuration.',
    );
  });

  it('replaces any client-supplied gateway header with the Pages binding value', async () => {
    const request = new Request('https://hashimojoe.com/api/feedback', {
      method: 'POST',
      headers: {
        [gatewayHeaderName()]: 'client-value',
        'Content-Type': 'application/json',
      },
      body: JSON.stringify({ message: 'Nice app' }),
    });

    const proxyRequest = createProxyRequest(
      request,
      {
        [originEnvName]: 'https://aip-origin.example.com',
        [secretEnvName]: 'server-value',
      },
      'feedback',
    );

    expect(proxyRequest.url).toBe('https://aip-origin.example.com/feedback');
    expect(proxyRequest.headers.get(gatewayHeaderName())).toBe('server-value');
    expect(proxyRequest.headers.get('X-Forwarded-Host')).toBe('hashimojoe.com');
    expect(proxyRequest.method).toBe('POST');
    expect(await proxyRequest.text()).toBe(JSON.stringify({ message: 'Nice app' }));
  });
});
