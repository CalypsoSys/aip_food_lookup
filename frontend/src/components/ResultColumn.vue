<script setup lang="ts">
withDefaults(
  defineProps<{
    title: string;
    items: string[];
    tone: 'allowed' | 'not-allowed';
    emptyText: string;
    showCount?: boolean;
  }>(),
  {
    showCount: true,
  },
);
</script>

<template>
  <section class="result-column" :class="tone">
    <h3>
      {{ title }}
      <span v-if="showCount" class="result-count">- {{ items.length }}</span>
    </h3>
    <ul v-if="items.length">
      <li v-for="item in items" :key="item">{{ item }}</li>
    </ul>
    <p v-else class="empty-copy">{{ emptyText }}</p>
  </section>
</template>
