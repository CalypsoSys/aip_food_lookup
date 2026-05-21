import type { CategoriesResult, FeedbackRequest, SearchResult, SearchType } from '@/types';

const apiBaseUrl = (import.meta.env.VITE_AIP_API_BASE_URL || '/api').replace(/\/$/, '');

export const searchTypes: SearchType[] = [
  'Search by Text and Sound',
  'Search by Text',
  'Search by Sound',
];

export function normalizeSearchType(value: SearchType): string {
  return value.replace(/\s/g, '').toLowerCase();
}

export function readSearchResult(json: unknown): SearchResult {
  const payload = json as { allowed?: unknown; not_allowed?: unknown };
  return {
    allowed: Array.isArray(payload.allowed) ? payload.allowed.filter(isString) : [],
    notAllowed: Array.isArray(payload.not_allowed) ? payload.not_allowed.filter(isString) : [],
  };
}

export async function searchFoods(query: string, searchType: SearchType): Promise<SearchResult> {
  const params = new URLSearchParams({
    key: query,
    type: normalizeSearchType(searchType),
  });
  const json = await getJson(`${apiBaseUrl}/search?${params.toString()}`);
  return readSearchResult(json);
}

export async function loadCategories(): Promise<CategoriesResult> {
  const json = await getJson(`${apiBaseUrl}/categories`);
  return readSearchResult(json);
}

export async function submitFeedback(request: FeedbackRequest): Promise<void> {
  await postJson(`${apiBaseUrl}/feedback`, request);
}

async function getJson(url: string): Promise<unknown> {
  const response = await fetch(url, {
    headers: {
      Accept: 'application/json',
      'X-AIP-Client': 'web',
      'X-AIP-App-Version': import.meta.env.VITE_AIP_APP_VERSION || 'dev',
    },
  });
  if (!response.ok) {
    throw new Error(`API request failed with status ${response.status}`);
  }
  return response.json();
}

async function postJson(url: string, body: unknown): Promise<void> {
  const response = await fetch(url, {
    method: 'POST',
    headers: {
      Accept: 'application/json',
      'Content-Type': 'application/json',
      'X-AIP-Client': 'web',
      'X-AIP-App-Version': import.meta.env.VITE_AIP_APP_VERSION || 'dev',
    },
    body: JSON.stringify(body),
  });
  if (!response.ok) {
    throw new Error(`API request failed with status ${response.status}`);
  }
}

function isString(value: unknown): value is string {
  return typeof value === 'string';
}
