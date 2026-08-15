# AIP Food Lookup Apple release plan

Status: Planning  
Last updated: 2026-08-15

This plan tracks the work needed to release the Flutter app through the same
no-Mac path used for Body Flow and Go:

```text
GitHub -> Codemagic -> App Store Connect -> TestFlight -> App Store
```

The Go API and Cloudflare gateway remain the production backend. Apple release
work applies to the Flutter client under `flutter-app/`.

## Confirmed release decisions

- App name: `AIP Food Lookup`
- Apple Bundle ID: `com.calypsosystems.aipfoodlookup`
- Apple App ID description: `AIP Food Lookup App`
- Supported devices: iPhone and iPad
- Distribution owner: Calypso Systems LLC Apple Developer organization
- Initial distribution: internal TestFlight
- Ads: disabled for the first TestFlight build
- Shared Android/iOS marketing version: `1.0.0`
- Planned shared build number: `10`
- Planned Flutter version value: `1.0.0+10`

The Bundle ID has already been registered in Apple Certificates, Identifiers &
Profiles. Body Flow and Go uses a separate Bundle ID and signing identity.

The repository currently uses `0.1.8+9`. Build number `10` is planned instead
of `1` because Android build numbers must increase after a build is uploaded to
Google Play. Using `1.0.0+10` is safe whether or not build 9 was uploaded and
keeps both platforms on the same version. Do not change `pubspec.yaml` until the
release implementation begins.

## Current repository state

- The Flutter app has an Android target but no `flutter-app/ios/` target.
- The Android application ID already matches the Apple Bundle ID:
  `com.calypsosystems.aipfoodlookup`.
- Production API requests default to `https://hashimojoe.com/api`.
- The app includes a bundled offline catalog fallback.
- The About screen includes an informational-only medical disclaimer.
- The privacy policy is still a draft and has unresolved contact and effective
  date fields.
- App configuration defaults the client name to `android`.
- Google Mobile Ads is initialized unconditionally and currently uses Android
  test configuration.
- Existing store screenshots can guide the Apple listing, but Apple-device
  screenshots must be captured from the iOS build.

## Phase 1: Create the iOS Flutter target

1. Confirm the working tree is clean and run the existing Flutter tests.
2. Generate the iOS platform files from `flutter-app/` using the installed
   Flutter stable toolchain.
3. Review the generated diff to ensure existing Android, Dart, and asset files
   were not overwritten unexpectedly.
4. Configure the Runner target with:
   - Bundle ID `com.calypsosystems.aipfoodlookup`
   - Display name `AIP Food Lookup`
   - iPhone and iPad support
   - An initial iOS 15 deployment target, unless current Flutter/plugin
     requirements require a higher target
   - Release version inherited from `pubspec.yaml`
5. Generate Apple app icons and launch assets from the committed identity
   artwork.
6. Keep production networking HTTPS-only. Do not add broad App Transport
   Security exceptions.

Acceptance gate:

- `flutter-app/ios/` is committed and contains no signing secrets.
- Existing Flutter tests and analysis pass.
- Codemagic can compile the unsigned iOS project using a current stable Xcode.

## Phase 2: Make application behavior platform-aware

1. Change client identification so iOS builds send `X-AIP-Client: ios` and
   Android builds continue sending `X-AIP-Client: android`.
2. Add iOS-specific Google Mobile Ads configuration.
3. Do not initialize or load Mobile Ads when `AIP_ADS_ENABLED=false`.
4. Use Google's iOS test App ID and banner unit ID during ad-enabled testing.
   Do not commit production AdMob identifiers.
5. Replace Android-only diagnostics imagery with a platform-neutral icon.
6. Add an easily accessible privacy-policy link to the About screen.
7. Add or update tests for every configuration and behavior change.

Acceptance gate:

- `flutter analyze` passes.
- `flutter test` passes.
- Android behavior is unchanged except for intentional platform-neutral UI.
- An ads-disabled iOS build launches without initializing the ads SDK.

## Phase 3: Finalize privacy and support material

1. Finalize and publish the privacy policy at:

   ```text
   https://hashimojoe.com/privacy/aip-food-lookup
   ```

2. Replace the draft effective date and contact placeholders.
3. Confirm the public policy opens without authentication or geoblocking.
4. Inventory data handling for Apple's App Privacy questionnaire, including:
   - Food search text sent to the API
   - Food suggestions and allowed/not-allowed selections
   - Optional feedback name, email, subject, and message
   - Backend access and error logging
   - Server-side Slack delivery
   - AdMob data practices when ads are enabled
5. Confirm retention and deletion-request language matches actual operations.
6. Keep the in-app and store listing language clear that the app is
   informational and is not medical advice.

Acceptance gate:

- The published policy matches the app, backend, logging, Slack, and selected
  advertising configuration.
- The policy is linked from the app and App Store Connect.
- A public support contact is available.

## Phase 4: Configure App Store Connect and Codemagic

### App Store Connect

1. Confirm all current Apple agreements are accepted.
2. Create the iOS app record using:

   ```text
   Name: AIP Food Lookup
   Bundle ID: com.calypsosystems.aipfoodlookup
   SKU: aip-food-lookup-ios
   Primary language: English (U.S.)
   ```

3. Create a dedicated App Store Connect API key for Codemagic with App Manager
   access.
4. Download the `.p8` file once and store it securely outside the repository.

### Codemagic

1. Connect the GitHub repository to Codemagic.
2. Set the Flutter project directory to `flutter-app`.
3. Select a current macOS build image with Xcode 26 or later and an iOS 26 SDK
   or later. Recheck Apple's upload requirement before each release.
4. Configure automatic signing with:
   - The dedicated App Store Connect API key
   - Bundle ID `com.calypsosystems.aipfoodlookup`
   - App Store provisioning profile type
   - Release build mode
5. Run Flutter analysis and tests before creating the release artifact.
6. Supply release configuration without committing secrets:

   ```text
   AIP_BACKEND_URL=https://hashimojoe.com/api
   AIP_CLIENT_NAME=ios
   AIP_APP_VERSION=1.0.0
   AIP_ADS_ENABLED=false
   ```

7. Publish the signed IPA to App Store Connect, but do not automatically submit
   it for App Review.

Acceptance gate:

- Codemagic produces a signed `.ipa` artifact.
- App Store Connect accepts and processes the build.
- No signing keys, API keys, AdMob IDs, or backend secrets appear in source or
  build logs.

## Phase 5: Prepare the App Store listing

Complete the following before public review:

- App description and promotional text
- Subtitle and keywords
- `Food & Drink` primary category
- Privacy policy and support URLs
- Copyright and content-rights responses
- Age rating
- Advertising declaration matching the submitted build
- App Privacy responses
- Export compliance
- Review contact and review notes
- At least one current iPhone screenshot
- Current iPad screenshots because the app supports iPad

Use screenshots captured from the actual iOS build. Avoid Android status bars,
navigation controls, or device frames that misrepresent the Apple experience.

Acceptance gate:

- App Store Connect shows no missing metadata.
- Listing statements and screenshots match the submitted build.

## Phase 6: Internal TestFlight validation

Assign the processed build to a Calypso Systems internal testing group and test
on at least one physical iPhone and one iPad.

Validate:

- Fresh installation and cold launch
- Search by text and sound
- Allowed and not-allowed results
- Categories and category details
- Offline catalog fallback
- Suggestions and feedback
- Production API connectivity
- About disclaimer and privacy-policy link
- Diagnostics content
- Portrait and landscape layouts where supported
- Larger Dynamic Type settings
- Ads remain disabled in the first build
- No Android-specific labels, icons, or behavior

Acceptance gate:

- No release-blocking crashes or layout defects.
- Core flows work on both iPhone and iPad.
- Backend logs show the iOS client identifier correctly.

## Phase 7: Enable and validate AdMob

Do this only after the ads-disabled TestFlight build is stable.

1. Register the iOS app in AdMob using the Apple Bundle ID.
2. Create an iOS banner unit.
3. Supply the production iOS AdMob App ID and banner unit ID through secure
   build configuration.
4. Update `Info.plist` with the required Google Mobile Ads and SKAdNetwork
   configuration.
5. Determine consent and App Tracking Transparency behavior from the final ad
   configuration and supported regions.
6. Update the privacy policy and App Privacy answers before distributing an
   ad-enabled build.
7. Upload a new build with an incremented build number and test ads through
   TestFlight before App Review.

Acceptance gate:

- Test ads work before production identifiers are introduced.
- Production ads do not appear in development builds.
- Consent, privacy policy, and App Store disclosures match actual SDK behavior.

## Phase 8: Submit for App Review

1. Select the approved TestFlight build.
2. Reconfirm privacy, advertising, age-rating, and export-compliance answers.
3. Add review notes explaining the informational food-reference purpose and
   how to exercise search, offline, suggestion, and feedback features.
4. Use manual release after approval for the first public version.
5. Monitor App Store correspondence, crash reports, feedback, and backend logs.

## Resume checklist

When implementation resumes, begin here:

- [ ] Confirm `1.0.0+10` is the desired shared Android/iOS version.
- [ ] Confirm the public support email or contact page.
- [ ] Confirm whether the App Store Connect app record already exists.
- [x] Run the existing Flutter analysis and test baseline.
- [x] Generate and inspect the iOS Flutter target.
- [x] Configure the Runner Bundle ID, iOS 15 baseline, iPhone/iPad support,
      display name, and identity assets.
- [ ] Implement the platform-aware client and ads-disabled startup behavior.
- [ ] Add the in-app privacy-policy link.
- [ ] Configure the first Codemagic unsigned iOS compile.
- [ ] Configure Apple automatic signing only after compilation succeeds.
- [ ] Upload the first ads-disabled internal TestFlight build.

## Secrets and generated artifacts

Never commit:

- App Store Connect `.p8` keys, Key IDs, or Issuer IDs
- Distribution certificates or provisioning profiles
- Production AdMob identifiers if repository policy keeps them external
- Backend gateway secrets or Slack webhook URLs
- Generated IPA files, Xcode archives, or signing exports

## Reference documentation

- Apple submission requirements:
  <https://developer.apple.com/app-store/submitting/>
- Apple App Review Guidelines:
  <https://developer.apple.com/app-store/review/guidelines/>
- Apple App Privacy management:
  <https://developer.apple.com/help/app-store-connect/manage-app-information/manage-app-privacy/>
- Flutter iOS release guide:
  <https://docs.flutter.dev/deployment/ios>
- Codemagic iOS signing:
  <https://docs.codemagic.io/flutter-code-signing/ios-code-signing/>
- Google Mobile Ads Flutter setup:
  <https://developers.google.com/admob/flutter/quick-start>
