<script setup lang="ts">
defineProps<{
  title: string;
  items: string[];
  tone: 'allowed' | 'not-allowed';
  emptyText: string;
}>();
</script>

<template>
  <section class="result-column" :class="tone">
    <h3>{{ title }} · {{ items.length }}</h3>
    <ul v-if="items.length">
      <li v-for="item in items" :key="item">{{ item }}</li>
    </ul>
    <p v-else class="empty-copy">{{ emptyText }}</p>
  </section>
</template>
