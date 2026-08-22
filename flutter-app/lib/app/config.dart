import 'package:flutter/foundation.dart';

class AppConfig {
  const AppConfig({
    required this.backendBaseUrl,
    this.clientName,
    this.appVersion = appVersionFromDefine,
  });

  static const productionBackendBaseUrl = 'https://hashimojoe.com/api';

  static const backendUrlFromDefine = String.fromEnvironment(
    'AIP_BACKEND_URL',
    defaultValue: productionBackendBaseUrl,
  );

  static const clientNameFromDefine = String.fromEnvironment(
    'AIP_CLIENT_NAME',
    defaultValue: '',
  );

  static const appVersionFromDefine = String.fromEnvironment(
    'AIP_APP_VERSION',
    defaultValue: 'prod',
  );

  static const adsEnabled = bool.fromEnvironment(
    'AIP_ADS_ENABLED',
    defaultValue: false,
  );

  static const adMobBannerAdUnitIdFromDefine = String.fromEnvironment(
    'AIP_ADMOB_BANNER_AD_UNIT_ID',
    defaultValue: '',
  );

  static const androidAdMobBannerAdUnitId =
      'ca-app-pub-3940256099942544/6300978111';
  static const iosAdMobBannerAdUnitId =
      'ca-app-pub-3940256099942544/2934735716';

  static const privacyPolicyUrl =
      'https://hashimojoe.com/privacy/aip-food-lookup';

  static const dev = AppConfig(
    backendBaseUrl: backendUrlFromDefine,
    clientName: clientNameFromDefine == '' ? null : clientNameFromDefine,
  );

  final String backendBaseUrl;
  final String? clientName;
  final String appVersion;

  static String clientNameForPlatform(TargetPlatform platform) {
    return platform == TargetPlatform.iOS ? 'ios' : 'android';
  }

  static String get adMobBannerAdUnitId {
    if (adMobBannerAdUnitIdFromDefine.trim().isNotEmpty) {
      return adMobBannerAdUnitIdFromDefine;
    }

    return defaultTargetPlatform == TargetPlatform.iOS
        ? iosAdMobBannerAdUnitId
        : androidAdMobBannerAdUnitId;
  }

  Map<String, String> get publicHeaders {
    final headers = <String, String>{};
    final resolvedClientName =
        clientName?.trim() ?? clientNameForPlatform(defaultTargetPlatform);
    if (resolvedClientName.isNotEmpty) {
      headers['X-AIP-Client'] = resolvedClientName;
    }
    if (appVersion.trim().isNotEmpty) {
      headers['X-AIP-App-Version'] = appVersion.trim();
    }
    return headers;
  }
}
