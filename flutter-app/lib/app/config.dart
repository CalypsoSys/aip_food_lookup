class AppConfig {
  const AppConfig({
    required this.backendBaseUrl,
    this.clientName = clientNameFromDefine,
    this.appVersion = appVersionFromDefine,
  });

  static const backendUrlFromDefine = String.fromEnvironment(
    'AIP_BACKEND_URL',
    defaultValue: 'http://10.0.2.2:8080',
  );

  static const clientNameFromDefine = String.fromEnvironment(
    'AIP_CLIENT_NAME',
    defaultValue: 'android',
  );

  static const appVersionFromDefine = String.fromEnvironment(
    'AIP_APP_VERSION',
    defaultValue: 'dev',
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
