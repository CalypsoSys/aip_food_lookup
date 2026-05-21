<script setup lang="ts">
import { useFoodStore } from '@/stores/foodStore';

const store = useFoodStore();
const apiBaseUrl = import.meta.env.VITE_AIP_API_BASE_URL || '/api';
const appVersion = import.meta.env.VITE_AIP_APP_VERSION || 'dev';
</script>

<template>
  <section class="panel narrow-panel" aria-labelledby="diagnostics-title">
    <div class="section-heading">
      <div>
        <h2 id="diagnostics-title">Diagnostics</h2>
        <p>Check local API connectivity and browser-facing configuration.</p>
      </div>
    </div>

    <dl class="diagnostics-list">
      <div>
        <dt>API base URL</dt>
        <dd>{{ apiBaseUrl }}</dd>
      </div>
      <div>
        <dt>App version</dt>
        <dd>{{ appVersion }}</dd>
      </div>
    </dl>

    <button class="primary-button" type="button" :disabled="store.checkingHealth" @click="store.checkApiHealth">
      {{ store.checkingHealth ? 'Checking...' : 'Check API health' }}
    </button>
    <p v-if="store.diagnosticsMessage" class="success-message">{{ store.diagnosticsMessage }}</p>
    <p v-if="store.errorMessage" class="error-message">{{ store.errorMessage }}</p>
  </section>
</template>
