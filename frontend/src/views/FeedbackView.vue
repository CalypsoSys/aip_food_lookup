<script setup lang="ts">
import FeedbackForm from '@/components/FeedbackForm.vue';
import { useFoodStore } from '@/stores/foodStore';

const store = useFoodStore();
</script>

<template>
  <section class="panel narrow-panel" aria-labelledby="feedback-title">
    <div class="section-heading">
      <div>
        <h2 id="feedback-title">Feedback</h2>
        <p>Send a note for catalog improvements or app issues.</p>
      </div>
    </div>

    <p v-if="store.errorMessage" class="error-message">{{ store.errorMessage }}</p>
    <FeedbackForm />
  </section>
</template>
