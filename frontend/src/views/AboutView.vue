<template>
  <section class="panel narrow-panel" aria-labelledby="about-title">
    <div class="section-heading">
      <div>
        <h2 id="about-title">About AIP Food Lookup</h2>
        <p>Quick food lookup for Autoimmune Protocol catalog guidance.</p>
      </div>
    </div>

    <div class="copy-stack">
      <p>
        Use Search for fast lookup, Categories to browse foods, and Suggestions or Feedback to
        help improve the catalog.
      </p>
      <p>This app is informational only and is not medical advice.</p>
    </div>

    <div class="action-grid">
      <RouterLink class="primary-link" to="/feedback">Feedback</RouterLink>
      <RouterLink class="secondary-link" to="/diagnostics">Diagnostics</RouterLink>
    </div>
  </section>
</template>
