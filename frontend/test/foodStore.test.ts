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
    submitFeedback: vi.fn().mockResolvedValue(undefined),
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
});
