# AIP VS Code launch profiles

The workspace includes local launch tasks for mobile and web development.

## Local: frontend + API

Use the `Local: frontend + API` compound launch to work on the Vue web frontend.

It runs:

- hidden Go API task at `http://127.0.0.1:8080`
- hidden Vite task at `http://127.0.0.1:5173`
- visible browser debug session for the frontend

The frontend calls `/api/*`; Vite proxies those calls to the Go API and strips the `/api` prefix.

## Flutter: Galaxy S23

Use `Flutter: Galaxy S23` for Android device debugging. It starts the Go API locally, runs `flutter pub get`, then
launches the Flutter app with `.vscode/aip-local.env`.

Do not put `AIP_GATEWAY_SECRET`, Slack webhooks, or private backend URLs in tracked VS Code files.
