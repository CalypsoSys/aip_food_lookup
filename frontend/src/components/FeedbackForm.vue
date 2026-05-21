<script setup lang="ts">
import { ref } from 'vue';

import { useFoodStore } from '@/stores/foodStore';

const store = useFoodStore();
const feedbackName = ref('');
const feedbackEmail = ref('');
const feedbackSubject = ref('');
const feedbackMessage = ref('');

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
  <form class="form-stack" @submit.prevent="sendFeedback">
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
    <button class="primary-button" type="submit" :disabled="store.submittingFeedback">
      {{ store.submittingFeedback ? 'Sending...' : 'Send feedback' }}
    </button>
    <p v-if="store.feedbackMessage" class="success-message">{{ store.feedbackMessage }}</p>
  </form>
</template>
