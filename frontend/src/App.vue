<script setup lang="ts">
import { RouterLink, RouterView } from 'vue-router';
</script>

<template>
  <div class="app-shell">
    <header class="top-bar">
      <div>
        <p class="eyebrow">AIP Food Lookup</p>
        <h1>Find ingredient guidance quickly</h1>
      </div>
      <nav class="top-nav" aria-label="Primary">
        <RouterLink to="/search">Search</RouterLink>
        <RouterLink to="/categories">Categories</RouterLink>
        <RouterLink to="/about">About</RouterLink>
      </nav>
    </header>

    <main class="page-wrap">
      <RouterView />
    </main>
  </div>
</template>
