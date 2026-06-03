<template>
  <section class="panel policy-panel" aria-labelledby="privacy-title">
    <div class="section-heading">
      <div>
        <h2 id="privacy-title">Privacy Policy</h2>
        <p>Effective date: June 3, 2026</p>
      </div>
    </div>

    <div class="copy-stack policy-copy">
      <p>
        AIP Food Lookup helps users check foods against an Autoimmune Protocol food list.
        The app is informational only and is not medical advice.
      </p>

      <h3>Information We Collect</h3>
      <p>The app may send the following information to our servers when you use it:</p>
      <ul>
        <li>Food search text that you enter.</li>
        <li>Food suggestions you submit, including the food text and whether you marked it allowed or not allowed.</li>
        <li>Feedback you choose to send, which may include your name, email address, subject, message, and app source.</li>
        <li>Basic app request details such as app client name, app version, request path, response status, and network metadata.</li>
      </ul>
      <p>
        Do not include medical records, health measurements, passwords, payment details,
        or other sensitive personal information in food searches, suggestions, or feedback messages.
      </p>

      <h3>How We Use Information</h3>
      <ul>
        <li>Return food lookup results.</li>
        <li>Review suggested food additions or corrections.</li>
        <li>Respond to support and feedback messages when contact information is provided.</li>
        <li>Monitor reliability, diagnose errors, prevent abuse, and improve the app.</li>
        <li>Maintain and improve the food catalog.</li>
      </ul>

      <h3>Feedback and Third-Party Services</h3>
      <p>
        Feedback and food suggestions may be delivered to server-side tools such as Slack so the maintainer can review
        them. If delivery to Slack is unavailable, the backend may store feedback or suggestions in server-side files for
        later review.
      </p>
      <p>
        The app does not put private backend secrets, Slack webhook URLs, or internal API keys on the user's device.
      </p>

      <h3>Advertising</h3>
      <p>
        The app may show ads. Advertising partners such as Google AdMob may process device, app, and ad interaction
        information to serve, measure, and improve ads according to their own policies and user controls.
      </p>

      <h3>Data Sharing</h3>
      <p>
        We do not sell personal information. We may share limited information with service providers that help operate
        the app, including hosting, logging, feedback routing, analytics, security, and advertising providers. These
        providers process information on our behalf or according to their own published terms.
      </p>
      <p>
        We may also disclose information if required by law, to protect users, or to protect the security and integrity
        of the app and related services.
      </p>

      <h3>Data Retention</h3>
      <p>
        We keep submitted feedback, suggestions, and operational logs only as long as needed to operate, support, secure,
        and improve the app. Retention periods may vary depending on backup, logging, and support workflows.
      </p>

      <h3>Children's Privacy</h3>
      <p>
        AIP Food Lookup is not directed to children. Users should not submit personal information about children through
        feedback or suggestion forms.
      </p>

      <h3>Security</h3>
      <p>
        We use reasonable technical and organizational measures to protect app data. No internet service can guarantee
        absolute security.
      </p>

      <h3>Your Choices</h3>
      <p>
        You can use the food lookup without submitting optional feedback contact information. If you want to request
        deletion of feedback or other information you submitted, contact us with enough detail to identify the request.
      </p>

      <h3>Contact</h3>
      <p>
        Use the <RouterLink to="/feedback">Feedback</RouterLink> page to contact us about this policy or request deletion
        of information you submitted.
      </p>

      <h3>Change History</h3>
      <p>We may update this policy as the app changes. The effective date will be updated when a new version is published.</p>
    </div>
  </section>
</template>
