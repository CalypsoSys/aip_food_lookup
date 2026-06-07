# AIP Food Lookup Flutter App

Flutter migration of the existing .NET MAUI mobile frontend. Android is the first supported target; iOS is planned for a later macOS/cloud build phase.

## Windows Android Development

Install Flutter stable for Windows, Android Studio, and an Android emulator or device. Then run:

```powershell
cd C:\CalypsoSystems\aip_food_lookup\flutter-app
flutter pub get
flutter analyze
flutter test
flutter run --dart-define=AIP_BACKEND_URL=http://10.0.2.2:8080 --dart-define=AIP_CLIENT_NAME=android --dart-define=AIP_APP_VERSION=dev
```

For a physical Android device, replace `10.0.2.2` with the Windows machine LAN IP address that can reach the Go backend.

## Configuration

Production builds default to `https://hashimojoe.com/api`. Local development should override the backend URL at run time:

```powershell
flutter run --dart-define=AIP_BACKEND_URL=http://10.0.2.2:8080 --dart-define=AIP_CLIENT_NAME=android --dart-define=AIP_APP_VERSION=dev
```

Do not commit private backend URLs, production AdMob IDs, signing keys, tokens, or certificates.
Do not put the internal API gateway secret in Flutter; production builds use `https://hashimojoe.com/api`,
where the Cloudflare Pages Function injects the secret server-side. Override the production default only when intentionally building against another backend.

```powershell
flutter build apk --dart-define=AIP_CLIENT_NAME=android --dart-define=AIP_APP_VERSION=prod
```

## Ads

The Search screen includes a Google Mobile Ads banner. By default it uses Google's Android test app ID and test banner ad
unit ID. Keep those defaults for local development and pre-store verification.

Production ad IDs should be supplied at build time and should not be committed:

```powershell
flutter build appbundle --release --dart-define=AIP_CLIENT_NAME=android --dart-define=AIP_APP_VERSION=prod --dart-define=AIP_ADMOB_BANNER_AD_UNIT_ID=<ad-unit-id> --android-project-arg=AIP_ADMOB_APP_ID=<app-id>
```

Set `--dart-define=AIP_ADS_ENABLED=false` to hide the banner in a build.

The app includes an About tab with links to Feedback and Diagnostics. Diagnostics shows the active backend URL and can test the `/` health endpoint. Use it when switching between emulator URLs and physical-phone LAN URLs.

Feedback posts to the Go backend `/feedback` endpoint; Slack delivery is controlled by server-side webhook
configuration and falls back to local JSONL storage when Slack is unavailable.

## App Identity Assets

The Flutter app now carries migrated MAUI/recovered assets under `assets/identity/` and `assets/images/`. Android launcher icons are generated into `android/app/src/main/res/mipmap-*` from `assets/identity/app_icon.png`, and the native splash screen uses a padded app-mark image generated from the same source.

Android uses the package/application ID `com.calypsosystems.aipfoodlookup` and the launcher label `AIP Food Lookup`. iOS remains planned for a later macOS/cloud build phase; the same committed identity assets should be used when generating the iOS runner.

## Current Milestone

Milestone 1 includes a Flutter project foundation, bottom navigation, app theme, API client, DTO models, Search flow, Categories route scaffold, test ad placeholder, and unit tests for DTO behavior.

The Go backend now exposes the MAUI app's expected `/search`, `/suggest`, `/categories`, and `/subcategory` endpoints. Start it before running backend-connected Flutter screens.
