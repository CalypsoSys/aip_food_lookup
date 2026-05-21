import { defineStore } from 'pinia';

import {
  checkHealth,
  loadCategories,
  loadSubcategory,
  searchFoods,
  searchTypes,
  submitFeedback,
  submitSuggestion,
} from '@/services/apiClient';
import type { CategoryKind, FeedbackRequest, SearchResult, SearchType } from '@/types';

interface FoodState {
  query: string;
  searchType: SearchType;
  result: SearchResult;
  categories: SearchResult;
  categoryFoods: SearchResult;
  selectedCategoryKind: CategoryKind;
  selectedSubcategory: string;
  recentSearches: string[];
  loadingSearch: boolean;
  loadingCategories: boolean;
  loadingCategoryFoods: boolean;
  submittingFeedback: boolean;
  submittingSuggestion: boolean;
  checkingHealth: boolean;
  errorMessage: string;
  feedbackMessage: string;
  suggestionMessage: string;
  diagnosticsMessage: string;
}

export const useFoodStore = defineStore('food', {
  state: (): FoodState => ({
    query: '',
    searchType: searchTypes[0],
    result: { allowed: [], notAllowed: [] },
    categories: { allowed: [], notAllowed: [] },
    categoryFoods: { allowed: [], notAllowed: [] },
    selectedCategoryKind: 'Allowed',
    selectedSubcategory: '',
    recentSearches: [],
    loadingSearch: false,
    loadingCategories: false,
    loadingCategoryFoods: false,
    submittingFeedback: false,
    submittingSuggestion: false,
    checkingHealth: false,
    errorMessage: '',
    feedbackMessage: '',
    suggestionMessage: '',
    diagnosticsMessage: '',
  }),
  actions: {
    clearSearch() {
      this.query = '';
      this.result = { allowed: [], notAllowed: [] };
      this.errorMessage = '';
      this.suggestionMessage = '';
    },
    async search() {
      const trimmed = this.query.trim();
      this.errorMessage = '';
      this.suggestionMessage = '';
      if (trimmed.length < 3) {
        this.result = { allowed: [], notAllowed: [] };
        return;
      }

      this.loadingSearch = true;
      try {
        this.result = await searchFoods(trimmed, this.searchType);
        this.recentSearches = updatedRecentSearches(this.recentSearches, trimmed);
      } catch {
        this.errorMessage = 'Search failed. Check that the API is reachable.';
      } finally {
        this.loadingSearch = false;
      }
    },
    async selectRecentSearch(query: string) {
      this.query = query;
      await this.search();
    },
    async loadCategories() {
      this.loadingCategories = true;
      this.errorMessage = '';
      try {
        this.categories = await loadCategories();
      } catch {
        this.errorMessage = 'Categories failed to load.';
      } finally {
        this.loadingCategories = false;
      }
    },
    async loadCategoryFoods(kind: CategoryKind, subcategory: string) {
      this.selectedCategoryKind = kind;
      this.selectedSubcategory = subcategory;
      this.loadingCategoryFoods = true;
      this.errorMessage = '';
      try {
        this.categoryFoods = await loadSubcategory(kind, subcategory);
      } catch {
        this.errorMessage = 'Could not load foods for this category.';
      } finally {
        this.loadingCategoryFoods = false;
      }
    },
    async suggestCurrentQuery(allowed: boolean) {
      const trimmed = this.query.trim();
      this.errorMessage = '';
      this.suggestionMessage = '';
      if (trimmed.length < 3) {
        this.errorMessage = 'Enter at least 3 characters before suggesting a food.';
        return;
      }

      this.submittingSuggestion = true;
      try {
        await submitSuggestion({ inputText: trimmed, allowed });
        this.suggestionMessage = `Thanks. We will review "${trimmed}" for the catalog.`;
      } catch {
        this.errorMessage = 'Suggestion could not be made.';
      } finally {
        this.submittingSuggestion = false;
      }
    },
    async sendFeedback(request: FeedbackRequest) {
      this.submittingFeedback = true;
      this.feedbackMessage = '';
      this.errorMessage = '';
      try {
        await submitFeedback(request);
        this.feedbackMessage = 'Thanks. Your feedback was sent.';
      } catch {
        this.errorMessage = 'Feedback could not be sent.';
      } finally {
        this.submittingFeedback = false;
      }
    },
    async checkApiHealth() {
      this.checkingHealth = true;
      this.diagnosticsMessage = '';
      this.errorMessage = '';
      try {
        this.diagnosticsMessage = await checkHealth();
      } catch {
        this.errorMessage = 'Could not reach the food lookup API.';
      } finally {
        this.checkingHealth = false;
      }
    },
  },
});

function updatedRecentSearches(current: string[], query: string): string[] {
  const normalized = query.trim();
  const withoutDuplicate = current.filter(
    (item) => item.toLowerCase() !== normalized.toLowerCase(),
  );
  return [normalized, ...withoutDuplicate].slice(0, 5);
}
