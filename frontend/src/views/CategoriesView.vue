<script setup lang="ts">
import { computed, onMounted, ref, watch } from 'vue';
import { RouterLink, useRoute } from 'vue-router';

import ResultColumn from '@/components/ResultColumn.vue';
import { useFoodStore } from '@/stores/foodStore';
import type { CategoryKind } from '@/types';

const route = useRoute();
const store = useFoodStore();
const filterText = ref('');

const routeKind = computed<CategoryKind | null>(() => {
  if (route.params.kind === 'allowed') {
    return 'Allowed';
  }
  if (route.params.kind === 'not-allowed') {
    return 'Not Allowed';
  }
  return null;
});
const routeSubcategory = computed(() =>
  typeof route.params.subcategory === 'string' ? route.params.subcategory : '',
);
const filteredFoods = computed(() => {
  const query = filterText.value.trim().toLowerCase();
  const items =
    store.selectedCategoryKind === 'Allowed'
      ? store.categoryFoods.allowed
      : store.categoryFoods.notAllowed;
  if (!query) {
    return items;
  }
  return items.filter((item) => item.toLowerCase().includes(query));
});

onMounted(() => {
  void store.loadCategories();
});

watch(
  (): [CategoryKind | null, string] => [routeKind.value, routeSubcategory.value],
  ([kind, subcategory]) => {
    if (kind && subcategory) {
      filterText.value = '';
      void store.loadCategoryFoods(kind, subcategory);
    }
  },
  { immediate: true },
);

function kindPath(kind: CategoryKind) {
  return kind === 'Allowed' ? 'allowed' : 'not-allowed';
}
</script>

<template>
  <div class="content-grid">
    <section class="panel" aria-labelledby="categories-title">
      <div class="section-heading">
        <div>
          <h2 id="categories-title">Categories</h2>
          <p>Browse allowed and not allowed foods by catalog group.</p>
        </div>
        <span v-if="store.loadingCategories" class="status-pill">Loading</span>
      </div>

      <p v-if="store.errorMessage" class="error-message">{{ store.errorMessage }}</p>

      <div class="category-columns">
        <div>
          <h3>Allowed</h3>
          <ul>
            <li v-for="category in store.categories.allowed" :key="category">
              <RouterLink :to="`/categories/${kindPath('Allowed')}/${encodeURIComponent(category)}`">
                {{ category }}
              </RouterLink>
            </li>
          </ul>
        </div>
        <div>
          <h3>Not allowed</h3>
          <ul>
            <li v-for="category in store.categories.notAllowed" :key="category">
              <RouterLink
                :to="`/categories/${kindPath('Not Allowed')}/${encodeURIComponent(category)}`"
              >
                {{ category }}
              </RouterLink>
            </li>
          </ul>
        </div>
      </div>
    </section>

    <section class="panel" aria-labelledby="category-detail-title">
      <div class="section-heading">
        <div>
          <h2 id="category-detail-title">
            {{ routeSubcategory || 'Select a category' }}
          </h2>
          <p>
            {{
              routeKind
                ? `${routeKind} foods in this category.`
                : 'Choose a category to view matching foods.'
            }}
          </p>
        </div>
        <span v-if="store.loadingCategoryFoods" class="status-pill">Loading</span>
      </div>

      <label v-if="routeKind" class="field">
        <span>Filter foods</span>
        <input v-model="filterText" type="search" placeholder="Filter this category" />
      </label>

      <ResultColumn
        v-if="routeKind"
        :title="store.selectedCategoryKind"
        :items="filteredFoods"
        :tone="store.selectedCategoryKind === 'Allowed' ? 'allowed' : 'not-allowed'"
        empty-text="No foods matched this filter."
      />
      <p v-else class="empty-copy">Category foods will appear here.</p>
    </section>
  </div>
</template>
