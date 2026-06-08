<script setup lang="ts">
import { computed, ref, watch } from 'vue';

import ResultColumn from '@/components/ResultColumn.vue';
import { searchTypes } from '@/services/apiClient';
import { useFoodStore } from '@/stores/foodStore';

const store = useFoodStore();
let searchTimer: number | undefined;
const suggestionExpanded = ref(false);
const suggestionQuery = ref('');

const query = computed(() => store.query.trim());
const hasMinimumQuery = computed(() => query.value.length > 2);
const allowedCount = computed(() => store.result.allowed.length);
const notAllowedCount = computed(() => store.result.notAllowed.length);
const hasMatches = computed(() => allowedCount.value > 0 || notAllowedCount.value > 0);
const hasSearchedCurrentQuery = computed(
  () => store.lastSearchedQuery.toLowerCase() === query.value.toLowerCase(),
);
const isLikelyPreparedFood = computed(() => /\s+/.test(query.value));
const suggestionChoicesVisible = computed(
  () =>
    !hasMatches.value ||
    store.submittingSuggestion ||
    (suggestionExpanded.value && suggestionQuery.value === query.value),
);
const showSuggestions = computed(
  () =>
    hasMinimumQuery.value &&
    hasSearchedCurrentQuery.value &&
    !store.loadingSearch,
);
const lookupMessage = computed(() => {
  if (store.errorMessage) {
    return {
      tone: 'error',
      title: store.errorMessage,
      subtitle: 'Check that your device can reach the configured backend URL.',
    };
  }

  if (!hasMinimumQuery.value) {
    return {
      tone: 'info',
      title: 'Search the ingredient catalog',
      subtitle: 'Type at least 3 characters. Results appear automatically.',
    };
  }

  if (store.loadingSearch || !hasSearchedCurrentQuery.value) {
    return {
      tone: 'info',
      title: 'Checking the catalog',
      subtitle: 'Looking for allowed and not allowed ingredients.',
    };
  }

  if (allowedCount.value > 0 && notAllowedCount.value > 0) {
    return {
      tone: 'mixed',
      title: 'Mixed catalog results',
      subtitle: `${allowedCount.value} allowed and ${notAllowedCount.value} not allowed matches found.`,
    };
  }

  if (hasMatches.value) {
    return null;
  }

  return {
    tone: 'info',
    title: isLikelyPreparedFood.value ? 'Prepared foods vary by recipe' : 'No ingredient match yet',
    subtitle: isLikelyPreparedFood.value
      ? 'Search the ingredient list one item at a time, like cherry, wheat, sugar, potato, or oil.'
      : 'Suggest this ingredient for review below.',
  };
});

watch(
  () => [store.query, store.searchType],
  () => {
    window.clearTimeout(searchTimer);
    suggestionExpanded.value = false;
    searchTimer = window.setTimeout(() => {
      void store.search();
    }, 350);
  },
);

function showSuggestionChoices() {
  suggestionExpanded.value = true;
  suggestionQuery.value = query.value;
}
</script>

<template>
  <section class="panel lookup-panel" aria-labelledby="lookup-title">
    <div class="section-heading">
      <div>
        <h2 id="lookup-title">Search</h2>
        <p>Search the ingredient catalog by text, sound, or both.</p>
      </div>
      <span v-if="store.loadingSearch" class="status-pill">Searching</span>
    </div>

    <div class="field-row">
      <label class="field">
        <span>Ingredient</span>
        <input v-model="store.query" type="search" placeholder="Search an ingredient" />
      </label>
      <button
        class="icon-button compact-icon"
        :class="{ invisible: !store.query }"
        type="button"
        :disabled="!store.query"
        aria-label="Clear search"
        title="Clear search"
        @click="store.clearSearch"
      >
        x
      </button>
    </div>

    <label class="field">
      <span>Search type</span>
      <select v-model="store.searchType">
        <option v-for="searchType in searchTypes" :key="searchType" :value="searchType">
          {{ searchType }}
        </option>
      </select>
    </label>

    <p class="helper">
      Best for single ingredients. For prepared foods, check the ingredient list one item at a
      time.
    </p>
    <p v-if="store.suggestionMessage" class="success-message">{{ store.suggestionMessage }}</p>

    <div class="lookup-counts" aria-label="Current result counts">
      <div class="lookup-count allowed" :class="{ active: allowedCount > 0 }">
        <span>Allowed</span>
        <strong>{{ allowedCount }}</strong>
      </div>
      <div class="lookup-count not-allowed" :class="{ active: notAllowedCount > 0 }">
        <span>Not allowed</span>
        <strong>{{ notAllowedCount }}</strong>
      </div>
    </div>

    <div
      v-if="lookupMessage"
      class="lookup-message"
      :class="lookupMessage.tone"
      role="status"
    >
      <strong>{{ lookupMessage.title }}</strong>
      <p>{{ lookupMessage.subtitle }}</p>
    </div>

    <div v-if="store.recentSearches.length" class="recent-searches-bar">
      <div class="recent-row" aria-label="Recent searches">
        <button
          v-for="recentSearch in store.recentSearches"
          :key="recentSearch"
          class="chip-button"
          type="button"
          @click="store.selectRecentSearch(recentSearch)"
        >
          {{ recentSearch }}
        </button>
      </div>
      <button
        class="icon-button compact-icon"
        type="button"
        aria-label="Clear recent searches"
        title="Clear recent searches"
        @click="store.clearRecentSearches"
      >
        x
      </button>
    </div>

    <div v-if="showSuggestions" class="suggestion-card">
      <button
        v-if="hasMatches && !suggestionChoicesVisible"
        class="suggestion-prompt"
        type="button"
        :disabled="store.submittingSuggestion"
        @click="showSuggestionChoices"
      >
        <span>
          <strong>Not seeing the ingredient you meant?</strong>
          <small>Suggest "{{ query }}"</small>
        </span>
        <span aria-hidden="true">v</span>
      </button>
      <div v-else>
        <h3>{{ hasMatches ? `Suggest "${query}"` : 'Missing from the catalog?' }}</h3>
        <p>
          {{
            hasMatches
              ? 'Send this ingredient for review if the exact item is missing.'
              : 'Send a suggestion only when the lookup does not find a clear match.'
          }}
        </p>
        <div class="suggestion-actions">
          <button
            class="secondary-button"
            type="button"
            :disabled="store.submittingSuggestion || store.query.trim().length < 3"
            @click="store.suggestCurrentQuery(true)"
          >
            {{ hasMatches ? 'Suggest allowed' : 'Suggest as allowed' }}
          </button>
          <button
            class="secondary-button danger"
            type="button"
            :disabled="store.submittingSuggestion || store.query.trim().length < 3"
            @click="store.suggestCurrentQuery(false)"
          >
            {{ hasMatches ? 'Suggest not allowed' : 'Suggest as not allowed' }}
          </button>
        </div>
        <div v-if="store.submittingSuggestion" class="progress-line" aria-label="Submitting" />
      </div>
    </div>

    <div class="results-grid">
      <ResultColumn
        title="Allowed ingredients"
        :items="store.result.allowed"
        tone="allowed"
        empty-text="Allowed ingredient matches will appear here."
        :show-count="false"
      />
      <ResultColumn
        title="Not allowed ingredients"
        :items="store.result.notAllowed"
        tone="not-allowed"
        empty-text="Not allowed ingredient matches will appear here."
        :show-count="false"
      />
    </div>
  </section>
</template>
