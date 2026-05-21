# AIP local development

## Go API

Run from the repo root:

```bash
cd cmd/aip_food_lookup
AIP_DATA_FOLDER=../../data go run .
```

The API listens on `:8080` by default. The local health check is:

```bash
curl -i http://127.0.0.1:8080/
```

For local gateway-secret testing:

```bash
cd cmd/aip_food_lookup
AIP__API__RequireGatewaySecret=true AIP__API__GatewaySecret=local-secret AIP_DATA_FOLDER=../../data go run .
```

Then:

```bash
curl -i "http://127.0.0.1:8080/search?key=apple"
curl -i -H "X-Internal-Api-Key: local-secret" "http://127.0.0.1:8080/search?key=apple"
```

## Flutter Android

The Android emulator can reach the host on `10.0.2.2`:

```bash
cd flutter-app
flutter pub get
flutter test
flutter run --dart-define=AIP_BACKEND_URL=http://10.0.2.2:8080 --dart-define=AIP_CLIENT_NAME=android --dart-define=AIP_APP_VERSION=dev
```

For production builds, point the app at:

```text
https://api.hashimojoe.com
```

Do not pass or compile `AIP_GATEWAY_SECRET` into Flutter. The Cloudflare Worker owns that internal key.

## Docker config render

Build `render-config-env` from `repos/babalu-yaml-env`, then place it at:

```text
scripts/aip/render-config-env
```

Copy the sample config:

```bash
cp scripts/aip/config.example.yaml scripts/aip/config.local.yaml
vi scripts/aip/config.local.yaml
```

Then test compose rendering:

```bash
scripts/aip/compose-aip.sh config
```
