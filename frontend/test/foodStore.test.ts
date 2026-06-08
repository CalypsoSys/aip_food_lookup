import { createPinia, setActivePinia } from 'pinia';
import { beforeEach, describe, expect, it, vi } from 'vitest';

import { useFoodStore } from '../src/stores/foodStore';

vi.mock('../src/services/apiClient', async () => {
  const actual = await vi.importActual<typeof import('../src/services/apiClient')>(
    '../src/services/apiClient',
  );
  return {
    ...actual,
    searchFoods: vi.fn().mockResolvedValue({ allowed: ['Apples'], notAllowed: [] }),
    loadCategories: vi.fn().mockResolvedValue({ allowed: ['Fruits'], notAllowed: ['Grains'] }),
    loadSubcategory: vi.fn().mockResolvedValue({ allowed: ['Apples'], notAllowed: [] }),
    submitFeedback: vi.fn().mockResolvedValue(undefined),
    submitSuggestion: vi.fn().mockResolvedValue(undefined),
    checkHealth: vi.fn().mockResolvedValue('AIP Food Lookup API'),
  };
});

describe('foodStore', () => {
  beforeEach(() => {
    setActivePinia(createPinia());
  });

  it('does not search until query has at least 3 characters', async () => {
    const store = useFoodStore();
    store.query = 'ap';

    await store.search();

    expect(store.result.allowed).toEqual([]);
    expect(store.loadingSearch).toBe(false);
  });

  it('loads categories from the API', async () => {
    const store = useFoodStore();

    await store.loadCategories();

    expect(store.categories.allowed).toEqual(['Fruits']);
    expect(store.categories.notAllowed).toEqual(['Grains']);
  });

  it('stores recent searches after a successful search', async () => {
    const store = useFoodStore();
    store.query = 'apple';

    await store.search();

    expect(store.result.allowed).toEqual(['Apples']);
    expect(store.lastSearchedQuery).toBe('apple');
    expect(store.recentSearches).toEqual(['apple']);
  });

  it('clears recent searches', async () => {
    const store = useFoodStore();
    store.query = 'apple';
    await store.search();

    store.clearRecentSearches();

    expect(store.recentSearches).toEqual([]);
  });

  it('loads ingredients for a selected category', async () => {
    const store = useFoodStore();

    await store.loadCategoryFoods('Allowed', 'Fruits');

    expect(store.selectedCategoryKind).toBe('Allowed');
    expect(store.selectedSubcategory).toBe('Fruits');
    expect(store.categoryFoods.allowed).toEqual(['Apples']);
  });

  it('submits the current query as a suggestion', async () => {
    const store = useFoodStore();
    store.query = 'cassava chips';

    await store.suggestCurrentQuery(true);

    expect(store.suggestionMessage).toContain('cassava chips');
    expect(store.submittingSuggestion).toBe(false);
  });

  it('checks API health for diagnostics', async () => {
    const store = useFoodStore();

    await store.checkApiHealth();

    expect(store.diagnosticsMessage).toBe('AIP Food Lookup API');
    expect(store.checkingHealth).toBe(false);
  });
});
