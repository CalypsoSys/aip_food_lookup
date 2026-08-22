import 'package:aip_food_lookup/app/config.dart';
import 'package:flutter/foundation.dart';
import 'package:flutter_test/flutter_test.dart';

void main() {
  test('default backend points at the production API gateway', () {
    expect(AppConfig.backendUrlFromDefine, 'https://hashimojoe.com/api');
    expect(AppConfig.appVersionFromDefine, 'prod');
  });

  test('publicHeaders includes mobile diagnostic headers', () {
    const config = AppConfig(
      backendBaseUrl: 'https://hashimojoe.com/api',
      clientName: 'android',
      appVersion: 'dev',
    );

    expect(config.publicHeaders, {
      'X-AIP-Client': 'android',
      'X-AIP-App-Version': 'dev',
    });
  });

  test('platform defaults identify iOS separately from Android', () {
    expect(
      AppConfig.clientNameForPlatform(TargetPlatform.iOS),
      'ios',
    );
    expect(
      AppConfig.clientNameForPlatform(TargetPlatform.android),
      'android',
    );
  });

  test('default config headers use the current platform client name', () {
    const config = AppConfig(backendBaseUrl: 'https://hashimojoe.com/api');

    expect(
      config.publicHeaders['X-AIP-Client'],
      AppConfig.clientNameForPlatform(defaultTargetPlatform),
    );
  });

  test('publicHeaders omits blank values', () {
    const config = AppConfig(
      backendBaseUrl: 'https://hashimojoe.com/api',
      clientName: ' ',
      appVersion: '',
    );

    expect(config.publicHeaders, isEmpty);
  });

  test('ad config defaults to Google test banner values', () {
    expect(AppConfig.adsEnabled, isFalse);
    expect(
      AppConfig.adMobBannerAdUnitId,
      anyOf(
        AppConfig.androidAdMobBannerAdUnitId,
        AppConfig.iosAdMobBannerAdUnitId,
      ),
    );
  });

  test('privacy policy URL is available to the app', () {
    expect(AppConfig.privacyPolicyUrl, startsWith('https://'));
  });
}
