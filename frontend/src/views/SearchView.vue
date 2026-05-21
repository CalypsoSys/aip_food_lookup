<script setup lang="ts">
import { watch } from 'vue';

import ResultColumn from '@/components/ResultColumn.vue';
import { searchTypes } from '@/services/apiClient';
import { useFoodStore } from '@/stores/foodStore';

const store = useFoodStore();
let searchTimer: number | undefined;

watch(
  () => [store.query, store.searchType],
  () => {
    window.clearTimeout(searchTimer);
    searchTimer = window.setTimeout(() => {
      void store.search();
    }, 350);
  },
);
</script>

<template>
  <section class="panel lookup-panel" aria-labelledby="lookup-title">
    <div class="section-heading">
      <div>
        <h2 id="lookup-title">Search</h2>
        <p>Search foods by text, sound, or both.</p>
      </div>
      <span v-if="store.loadingSearch" class="status-pill">Searching</span>
    </div>

    <div class="field-row">
      <label class="field">
        <span>Food</span>
        <input v-model="store.query" type="search" placeholder="Type at least 3 characters" />
      </label>
      <button class="icon-button" type="button" :disabled="!store.query" @click="store.clearSearch">
        Clear
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

    <p class="helper">Type at least 3 characters. Results update after you pause typing.</p>
    <p v-if="store.errorMessage" class="error-message">{{ store.errorMessage }}</p>
    <p v-if="store.suggestionMessage" class="success-message">{{ store.suggestionMessage }}</p>

    <div v-if="store.recentSearches.length" class="recent-row" aria-label="Recent searches">
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

    <div class="suggestion-actions">
      <button
        class="secondary-button"
        type="button"
        :disabled="store.submittingSuggestion || store.query.trim().length < 3"
        @click="store.suggestCurrentQuery(true)"
      >
        Suggest allowed
      </button>
      <button
        class="secondary-button danger"
        type="button"
        :disabled="store.submittingSuggestion || store.query.trim().length < 3"
        @click="store.suggestCurrentQuery(false)"
      >
        Suggest not allowed
      </button>
    </div>

    <div class="results-grid">
      <ResultColumn
        title="Allowed on AIP"
        :items="store.result.allowed"
        tone="allowed"
        empty-text="Allowed matches will appear here."
      />
      <ResultColumn
        title="Not allowed on AIP"
        :items="store.result.notAllowed"
        tone="not-allowed"
        empty-text="Not allowed matches will appear here."
      />
    </div>
  </section>
</template>
