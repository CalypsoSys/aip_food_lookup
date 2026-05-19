class AppConfig {
  const AppConfig({required this.backendBaseUrl});

  static const backendUrlFromDefine = String.fromEnvironment(
    'AIP_BACKEND_URL',
    defaultValue: 'http://10.0.2.2:8080',
  );

  static const dev = AppConfig(backendBaseUrl: backendUrlFromDefine);

  final String backendBaseUrl;
}
