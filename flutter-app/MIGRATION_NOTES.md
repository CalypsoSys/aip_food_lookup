# Migration Notes

## Assumptions

- Flutter work lives in `C:\CalypsoSystems\aip_food_lookup\flutter-app`.
- The MAUI reference app lives in `C:\CalypsoSystems\AIPFoodLookup` and remains read-only.
- The visible Go backend lives in `C:\CalypsoSystems\aip_food_lookup\cmd\aip_food_lookup`.
- Android emulator development uses `http://10.0.2.2:8080` by default. A physical Android device should use the Windows host LAN IP instead of `localhost`.
- User-facing spelling is corrected from MAUI typos, for example `Catagories` to `Categories`.

## Backend Status

The Go backend has been restored with the MAUI-expected endpoints:

- `GET /search?key=<text>&type=<searchbytextandsound|searchbytext|searchbysound>`
- `POST /suggest` with `{ "inputText": "food", "allowed": true }`
- `GET /categories`
- `GET /subcategory?cat=<Allowed|Not Allowed>&sub=<subcategory>`

For this milestone, the recovered `Nonsense-I-Know` private-key header check was intentionally removed so the Flutter client can call the API during local development. A future milestone should add a maintainable auth or abuse-prevention mechanism before production deployment.

## AdMob

Milestone 1 includes a UI placeholder only. Production AdMob IDs must stay out of source control. A later milestone can add `google_mobile_ads` with test IDs first.

## iOS

iOS is planned structurally, but this milestone does not require a Mac. Once the Flutter SDK is available, iOS project files can be generated from macOS or a cloud build environment.
