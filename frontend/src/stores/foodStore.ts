import { defineStore } from 'pinia';

import {
  loadCategories,
  searchFoods,
  searchTypes,
  submitFeedback,
} from '@/services/apiClient';
import type { FeedbackRequest, SearchResult, SearchType } from '@/types';

interface FoodState {
  query: string;
  searchType: SearchType;
  result: SearchResult;
  categories: SearchResult;
  loadingSearch: boolean;
  loadingCategories: boolean;
  submittingFeedback: boolean;
  errorMessage: string;
  feedbackMessage: string;
}

export const useFoodStore = defineStore('food', {
  state: (): FoodState => ({
    query: '',
    searchType: searchTypes[0],
    result: { allowed: [], notAllowed: [] },
    categories: { allowed: [], notAllowed: [] },
    loadingSearch: false,
    loadingCategories: false,
    submittingFeedback: false,
    errorMessage: '',
    feedbackMessage: '',
  }),
  actions: {
    async search() {
      const trimmed = this.query.trim();
      this.errorMessage = '';
      if (trimmed.length < 3) {
        this.result = { allowed: [], notAllowed: [] };
        return;
      }

      this.loadingSearch = true;
      try {
        this.result = await searchFoods(trimmed, this.searchType);
      } catch {
        this.errorMessage = 'Search failed. Check that the API is reachable.';
      } finally {
        this.loadingSearch = false;
      }
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
  },
});
