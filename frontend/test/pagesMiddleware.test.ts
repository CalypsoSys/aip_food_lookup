import { describe, expect, it } from 'vitest';

import { isCredentialProbePath } from '../functions/_middleware';

describe('Pages middleware credential probe filter', () => {
  it('blocks common secret file probes before SPA fallback or API proxying', () => {
    expect(isCredentialProbePath('/.env')).toBe(true);
    expect(isCredentialProbePath('/.env.production')).toBe(true);
    expect(isCredentialProbePath('/api/.env')).toBe(true);
    expect(isCredentialProbePath('/config/service-account.json')).toBe(true);
    expect(isCredentialProbePath('/.aws/credentials')).toBe(true);
  });

  it('allows normal app and API routes', () => {
    expect(isCredentialProbePath('/')).toBe(false);
    expect(isCredentialProbePath('/search')).toBe(false);
    expect(isCredentialProbePath('/api/search')).toBe(false);
    expect(isCredentialProbePath('/api/categories')).toBe(false);
  });
});
