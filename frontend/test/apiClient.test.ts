import { describe, expect, it, vi } from 'vitest';

import {
  checkHealth,
  loadSubcategory,
  normalizeSearchType,
  normalizeSubcategory,
  readSearchResult,
  searchFoods,
  submitFeedback,
  submitSuggestion,
} from '../src/services/apiClient';

describe('apiClient', () => {
  it('normalizes search types for the Go API', () => {
    expect(normalizeSearchType('Search by Text and Sound')).toBe('searchbytextandsound');
    expect(normalizeSearchType('Search by Text')).toBe('searchbytext');
    expect(normalizeSearchType('Search by Sound')).toBe('searchbysound');
  });

  it('normalizes category labels for the Go API', () => {
    expect(normalizeSubcategory('Herbs and Spices')).toBe('herbs_spices');
    expect(normalizeSubcategory('Fruits')).toBe('fruits');
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

  it('loadSubcategory calls the category endpoint with route-compatible values', async () => {
    const fetchMock = vi.fn().mockResolvedValue({
      ok: true,
      json: async () => ({ allowed: ['Apples'], not_allowed: [] }),
    });
    vi.stubGlobal('fetch', fetchMock);

    await loadSubcategory('Allowed', 'Herbs and Spices');

    expect(fetchMock).toHaveBeenCalledWith(
      '/api/subcategory?cat=Allowed&sub=herbs_spices',
      expect.objectContaining({
        headers: expect.objectContaining({
          'X-AIP-Client': 'web',
        }),
      }),
    );
  });

  it('submitSuggestion posts the MAUI-compatible suggestion shape', async () => {
    const fetchMock = vi.fn().mockResolvedValue({ ok: true });
    vi.stubGlobal('fetch', fetchMock);

    await submitSuggestion({ inputText: 'cassava chips', allowed: true });

    expect(fetchMock).toHaveBeenCalledWith(
      '/api/suggest',
      expect.objectContaining({
        method: 'POST',
        body: JSON.stringify({ inputText: 'cassava chips', allowed: true }),
      }),
    );
  });

  it('checkHealth reads the API health text', async () => {
    const fetchMock = vi.fn().mockResolvedValue({
      ok: true,
      text: async () => 'AIP Food Lookup API',
    });
    vi.stubGlobal('fetch', fetchMock);

    await expect(checkHealth()).resolves.toBe('AIP Food Lookup API');
    expect(fetchMock).toHaveBeenCalledWith(
      '/api/',
      expect.objectContaining({
        headers: expect.objectContaining({
          Accept: 'text/plain',
        }),
      }),
    );
  });
});
