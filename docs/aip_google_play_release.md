# AIP Google Play release plan

This plan tracks the work needed to move the Flutter Android app from local development to Google Play release with
AdMob support.

## Recommended account path

Use the organizational Google Play developer account for this app. It is the better fit for a Calypso Systems package
name, production ownership, and future handoff. Keep the personal account available for personal experiments, but do not
split this app across accounts once package ownership starts.

If the app is published through a newly created personal account, Google Play may require a closed test with at least 12
testers opted in for 14 continuous days before production access. The organizational account is the cleaner release path
when available.

## Privacy policy hosting

Host the privacy policy on the same public domain used by the app and API:

```text
https://hashimojoe.com/privacy/aip-food-lookup
```

The policy URL should be stable before the first Play Console submission and before AdMob review. It should remain
available without sign-in, geoblocking, or JavaScript-only rendering.

## Current app data inventory

The Flutter app and backend currently expose these user-data surfaces:

- Search requests send the food query to the backend.
- Suggestions send the submitted food text and whether the user marked it allowed or not allowed.
- Feedback can include name, email, subject, source, and message.
- Mobile requests can include `X-AIP-Client` and `X-AIP-App-Version` headers when configured.
- Production feedback and suggestions may be posted to Slack through server-side configuration.
- Backend access/error logs may include request metadata such as path, status, and client network information.
- The app currently uses `INTERNET` and `ACCESS_NETWORK_STATE` Android permissions.

The app should not collect medical records, health measurements, precise location, contacts, photos, advertising ID, or
device identifiers beyond ordinary networking and configured client/version headers.

## Store listing draft inputs

Working values for the first Play Console draft:

- App name: `AIP Food Lookup`
- Package name: `com.calypsosystems.aipfoodlookup`
- Category: `Food & Drink` or `Health & Fitness`; choose `Food & Drink` if positioning as a food lookup tool rather than
  medical guidance.
- Short description: `Quickly check foods against an Autoimmune Protocol food list.`
- Medical disclaimer: `AIP Food Lookup is informational only and is not medical advice.`
- Support/contact path: use the app feedback form and a public support email or web contact page.

## Release build checklist

- Confirm production backend path is `https://hashimojoe.com/api`.
- Keep Android cleartext traffic restricted to debug builds.
- Confirm target SDK satisfies current Google Play requirements.
- Configure Play App Signing and upload signing keys outside the repository.
- Build a release Android App Bundle with production dart defines.
- Upload to internal testing before closed or production tracks.
- Complete Play Console app access, ads, content rating, target audience, privacy policy, and Data Safety sections.
- Keep `AIP_GATEWAY_SECRET`, Slack webhook URLs, signing keys, keystores, and AdMob production IDs out of tracked files.

## AdMob readiness checklist

- Restore or create the AdMob account.
- Create the AdMob app using package `com.calypsosystems.aipfoodlookup`.
- Publish `app-ads.txt` at the root of the developer website domain used in store listings.
- Add the Flutter Google Mobile Ads SDK with test ad unit IDs first.
- Add the Android AdMob app ID as manifest metadata through non-secret build configuration.
- Use production ad unit IDs only after AdMob app verification and Play review are ready.

## First implementation sequence

1. Finalize privacy policy and store listing copy.
2. Add a release-readiness issue/checklist for Android signing, target SDK, and cleartext traffic.
3. Add AdMob test-banner plumbing behind a config flag or test ad unit default.
4. Build and test a production-configured Android App Bundle.
5. Upload to Play internal testing.
6. Complete AdMob app verification and then switch to production ad unit IDs.

## Android release build notes

Use Play App Signing and keep the upload keystore outside git. A typical local signing setup uses an untracked
`flutter-app/android/key.properties` file that points to an untracked keystore file.

Production release builds should use HTTPS only:

```powershell
cd flutter-app
flutter build appbundle --release --dart-define=AIP_BACKEND_URL=https://hashimojoe.com/api --dart-define=AIP_CLIENT_NAME=android --dart-define=AIP_APP_VERSION=prod
```

Local debug builds can still use `http://10.0.2.2:8080` or a LAN backend URL because debug builds carry a debug-only
manifest override for cleartext traffic. The main manifest is release-safe.
