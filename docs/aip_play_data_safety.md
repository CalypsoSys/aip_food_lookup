# AIP Google Play Data Safety worksheet

This worksheet prepares the Google Play Data Safety form for AIP Food Lookup. It is a working draft for Play Console
entry; verify final answers against the published app build, privacy policy, backend configuration, and any SDKs enabled
at release time.

Google Play asks developers to disclose app data collection, sharing, security practices, and purposes across every
version, region, and user age distributed under the package name.

## Baseline release assumption

Use this worksheet for the first Android release before production AdMob IDs are enabled.

- App package: `com.calypsosystems.aipfoodlookup`
- Public API: `https://hashimojoe.com/api`
- Privacy policy: `https://hashimojoe.com/privacy/aip-food-lookup`
- Ads: not enabled in production until AdMob setup is complete
- User accounts: none
- Login: none
- Payments: none
- Location permissions: none
- Contacts/calendar/photo permissions: none

If production AdMob is enabled before Play submission, update the advertising section before submitting Data Safety.

## App-level answers

| Question | Draft answer | Notes |
| --- | --- | --- |
| Does the app collect or share any required user data types? | Yes | Search, suggestions, feedback, and operational request metadata are sent to the backend. |
| Is all collected user data encrypted in transit? | Yes | Production API uses HTTPS through `https://hashimojoe.com/api`. Local dev builds are not Play release builds. |
| Can users request deletion of data? | Yes | Privacy policy should include a support contact for deletion requests. |
| Is the app independently validated against a global security standard? | No | Do not claim independent validation unless completed. |
| Does the app allow users to create accounts? | No | No account creation or login flow exists. |

## Data types to disclose

### Personal info

| Data type | Collected? | Shared? | Required? | Purpose | Notes |
| --- | --- | --- | --- | --- | --- |
| Name | Yes, if submitted | Yes, service provider | Optional | App functionality, developer communications | Feedback form can include name. Slack/server tooling may receive it. |
| Email address | Yes, if submitted | Yes, service provider | Optional | App functionality, developer communications | Feedback form can include email for replies. |
| Other info | Yes, if submitted | Yes, service provider | Optional | App functionality, developer communications | Feedback subject/source/message can include user-provided personal info. |

### App activity

| Data type | Collected? | Shared? | Required? | Purpose | Notes |
| --- | --- | --- | --- | --- | --- |
| App interactions | Yes | Yes, service provider | Required for API use | App functionality, analytics, fraud prevention/security | Search terms, suggestion actions, endpoint paths, status codes, and app requests are processed by the backend. |
| Other user-generated content | Yes, if submitted | Yes, service provider | Optional for suggestions/feedback | App functionality, developer communications | Food suggestions and feedback messages are user-submitted content. |

### App info and performance

| Data type | Collected? | Shared? | Required? | Purpose | Notes |
| --- | --- | --- | --- | --- | --- |
| Diagnostics | Yes | Yes, service provider | Required for service operation | Analytics, fraud prevention/security | Backend access/error logs may capture failures, status, request metadata, and operational diagnostics. |

## Data types not expected in baseline release

Do not disclose these as collected unless code, SDKs, backend behavior, or Play services change before release:

- Approximate or precise location
- User IDs
- Phone number
- Payment information
- Purchase history
- Contacts
- Photos or videos
- Audio files
- Files and docs
- Calendar events
- Health and fitness data
- SMS or emails
- Installed apps
- Web browsing history
- Device or other IDs, except if introduced by AdMob or another SDK

## Sharing notes

Treat Slack, hosting, logs, and backend infrastructure as service providers when they process data on behalf of the app.
If an SDK or service uses app data for its own purposes, advertising profiles, cross-app measurement, or non-service
provider activity, disclose that as sharing according to Google Play guidance.

## Advertising update path

When AdMob is enabled in production:

- Review Google Mobile Ads SDK data disclosures for the exact SDK version.
- Add any required data types such as device or other IDs, app interactions, diagnostics, or advertising data if the SDK
  collects them.
- Mark advertising-related purposes where applicable.
- Confirm whether data is shared with Google for ads, measurement, fraud prevention, or analytics.
- Update the privacy policy before switching from test ad units to production ad units.

## Open decisions before Play submission

- Final support/privacy contact email.
- Final privacy policy effective date.
- Whether initial category is `Food & Drink` or `Health & Fitness`.
- Whether production ads are enabled before first Play submission or deferred until after initial approval.
- Whether backend logs include IP addresses in a way that should be documented more explicitly in the privacy policy.

## Verification checklist

- Build release app with production dart defines.
- Confirm Android manifest permissions match this worksheet.
- Confirm no production AdMob SDK or ad unit IDs are active unless the advertising update path is completed.
- Confirm privacy policy matches Play Data Safety answers.
- Confirm Play Console Data Safety answers match the final uploaded artifact and backend behavior.
