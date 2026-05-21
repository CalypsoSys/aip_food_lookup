import { describe, expect, it, vi } from 'vitest';

import {
  normalizeSearchType,
  readSearchResult,
  searchFoods,
  submitFeedback,
} from '../src/services/apiClient';

describe('apiClient', () => {
  it('normalizes search types for the Go API', () => {
    expect(normalizeSearchType('Search by Text and Sound')).toBe('searchbytextandsound');
    expect(normalizeSearchType('Search by Text')).toBe('searchbytext');
    expect(normalizeSearchType('Search by Sound')).toBe('searchbysound');
  });

  it('reads allowed and not_allowed response arrays', () => {
    expect(
      readSearchResult({
        allowed: ['Apples', 1],
        not_allowed: ['Wheat', null],
      }),
    ).toEqual({
      allowed: ['Apples'],
      notAllowed: ['Wheat'],
    });
  });

  it('searchFoods calls the public client path without gateway secrets', async () => {
    const fetchMock = vi.fn().mockResolvedValue({
      ok: true,
      json: async () => ({ allowed: ['Apples'], not_allowed: [] }),
    });
    vi.stubGlobal('fetch', fetchMock);

    await searchFoods('apple', 'Search by Text');

    const [, requestInit] = fetchMock.mock.calls[0];
    const headers = requestInit.headers as Record<string, string>;
    const internalHeaderName = ['X', 'Internal', 'Api', 'Key'].join('-');

    expect(fetchMock).toHaveBeenCalledWith(
      '/api/search?key=apple&type=searchbytext',
      expect.objectContaining({
        headers: expect.objectContaining({
          'X-AIP-Client': 'web',
        }),
      }),
    );
    expect(Object.keys(headers)).not.toContain(internalHeaderName);
  });

  it('submitFeedback posts web feedback payloads', async () => {
    const fetchMock = vi.fn().mockResolvedValue({ ok: true });
    vi.stubGlobal('fetch', fetchMock);

    await submitFeedback({
      name: 'Joe',
      email: '',
      subject: 'Idea',
      message: 'Nice app',
      source: 'web',
    });

    expect(fetchMock).toHaveBeenCalledWith(
      '/api/feedback',
      expect.objectContaining({
        method: 'POST',
        body: JSON.stringify({
          name: 'Joe',
          email: '',
          subject: 'Idea',
          message: 'Nice app',
          source: 'web',
        }),
      }),
    );
  });
});
