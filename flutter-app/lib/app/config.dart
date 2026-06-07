class AppConfig {
  const AppConfig({
    required this.backendBaseUrl,
    this.clientName = clientNameFromDefine,
    this.appVersion = appVersionFromDefine,
  });

  static const productionBackendBaseUrl = 'https://hashimojoe.com/api';

  static const backendUrlFromDefine = String.fromEnvironment(
    'AIP_BACKEND_URL',
    defaultValue: productionBackendBaseUrl,
  );

  static const clientNameFromDefine = String.fromEnvironment(
    'AIP_CLIENT_NAME',
    defaultValue: 'android',
  );

  static const appVersionFromDefine = String.fromEnvironment(
    'AIP_APP_VERSION',
    defaultValue: 'prod',
  );

  static const adsEnabled = bool.fromEnvironment(
    'AIP_ADS_ENABLED',
    defaultValue: true,
  );

  static const adMobBannerAdUnitId = String.fromEnvironment(
    'AIP_ADMOB_BANNER_AD_UNIT_ID',
    defaultValue: 'ca-app-pub-3940256099942544/6300978111',
  );

  static const dev = AppConfig(backendBaseUrl: backendUrlFromDefine);

  final String backendBaseUrl;
  final String clientName;
  final String appVersion;

  Map<String, String> get publicHeaders {
    final headers = <String, String>{};
    if (clientName.trim().isNotEmpty) {
      headers['X-AIP-Client'] = clientName.trim();
    }
    if (appVersion.trim().isNotEmpty) {
      headers['X-AIP-App-Version'] = appVersion.trim();
    }
    return headers;
  }
}
