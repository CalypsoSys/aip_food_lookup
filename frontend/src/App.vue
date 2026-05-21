<script setup lang="ts">
import { onMounted, ref, watch } from 'vue';

import { searchTypes } from '@/services/apiClient';
import { useFoodStore } from '@/stores/foodStore';

const store = useFoodStore();
const feedbackName = ref('');
const feedbackEmail = ref('');
const feedbackSubject = ref('');
const feedbackMessage = ref('');
let searchTimer: number | undefined;

onMounted(() => {
  void store.loadCategories();
});

watch(
  () => [store.query, store.searchType],
  () => {
    window.clearTimeout(searchTimer);
    searchTimer = window.setTimeout(() => {
      void store.search();
    }, 300);
  },
);

function clearSearch() {
  store.query = '';
  store.result = { allowed: [], notAllowed: [] };
}

function canSendFeedback() {
  return (
    feedbackMessage.value.trim().length > 0 &&
    (feedbackName.value.trim().length > 0 || feedbackEmail.value.trim().length > 0)
  );
}

async function sendFeedback() {
  if (!canSendFeedback()) {
    store.errorMessage = 'Enter a message and either a name or email.';
    return;
  }
  await store.sendFeedback({
    name: feedbackName.value,
    email: feedbackEmail.value,
    subject: feedbackSubject.value,
    message: feedbackMessage.value,
    source: 'web',
  });
  if (store.feedbackMessage) {
    feedbackSubject.value = '';
    feedbackMessage.value = '';
  }
}
</script>

<template>
  <div class="app-shell">
    <header class="top-bar">
      <div>
        <p class="eyebrow">AIP Food Lookup</p>
        <h1>Find food guidance quickly</h1>
      </div>
      <a href="#feedback" class="header-link">Feedback</a>
    </header>

    <main class="main-grid">
      <section class="panel lookup-panel" aria-labelledby="lookup-title">
        <div class="section-heading">
          <div>
            <h2 id="lookup-title">Lookup</h2>
            <p>Search foods by text, sound, or both.</p>
          </div>
          <span v-if="store.loadingSearch" class="status-pill">Searching</span>
        </div>

        <div class="field-row">
          <label class="field">
            <span>Food</span>
            <input v-model="store.query" type="search" placeholder="Type at least 3 characters" />
          </label>
          <button class="icon-button" type="button" :disabled="!store.query" @click="clearSearch">
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

        <div class="results-grid">
          <section class="result-column allowed">
            <h3>Allowed on AIP · {{ store.result.allowed.length }}</h3>
            <ul v-if="store.result.allowed.length">
              <li v-for="food in store.result.allowed" :key="food">{{ food }}</li>
            </ul>
            <p v-else class="empty-copy">Allowed matches will appear here.</p>
          </section>

          <section class="result-column not-allowed">
            <h3>Not allowed on AIP · {{ store.result.notAllowed.length }}</h3>
            <ul v-if="store.result.notAllowed.length">
              <li v-for="food in store.result.notAllowed" :key="food">{{ food }}</li>
            </ul>
            <p v-else class="empty-copy">Not allowed matches will appear here.</p>
          </section>
        </div>
      </section>

      <aside class="side-stack">
        <section class="panel" aria-labelledby="categories-title">
          <div class="section-heading">
            <div>
              <h2 id="categories-title">Categories</h2>
              <p>Browse the catalog while search evolves.</p>
            </div>
          </div>

          <div class="category-columns">
            <div>
              <h3>Allowed</h3>
              <ul>
                <li v-for="category in store.categories.allowed" :key="category">{{ category }}</li>
              </ul>
            </div>
            <div>
              <h3>Not allowed</h3>
              <ul>
                <li v-for="category in store.categories.notAllowed" :key="category">
                  {{ category }}
                </li>
              </ul>
            </div>
          </div>
        </section>

        <section id="feedback" class="panel" aria-labelledby="feedback-title">
          <div class="section-heading">
            <div>
              <h2 id="feedback-title">Feedback</h2>
              <p>Send a note for catalog improvements or app issues.</p>
            </div>
          </div>

          <div class="compact-grid">
            <label class="field">
              <span>Name</span>
              <input v-model="feedbackName" type="text" placeholder="Name or email required" />
            </label>
            <label class="field">
              <span>Email</span>
              <input v-model="feedbackEmail" type="email" placeholder="Optional if name is entered" />
            </label>
          </div>
          <label class="field">
            <span>Subject</span>
            <input v-model="feedbackSubject" type="text" placeholder="App feedback" />
          </label>
          <label class="field">
            <span>Message</span>
            <textarea v-model="feedbackMessage" rows="4" placeholder="What should we know?" />
          </label>
          <button
            class="primary-button"
            type="button"
            :disabled="store.submittingFeedback"
            @click="sendFeedback"
          >
            {{ store.submittingFeedback ? 'Sending...' : 'Send feedback' }}
          </button>
          <p v-if="store.feedbackMessage" class="success-message">{{ store.feedbackMessage }}</p>
        </section>
      </aside>
    </main>
  </div>
</template>
