import 'package:flutter/foundation.dart';

import '../../../app/config.dart';
import '../../../core/networking/api_client.dart';

typedef HealthCheck = Future<String> Function();

const _diagnosticsTimeout = Duration(seconds: 2);

class DiagnosticsState {
  const DiagnosticsState({
    required this.backendUrl,
    this.isChecking = false,
    this.healthMessage,
    this.errorMessage,
  });

  final String backendUrl;
  final bool isChecking;
  final String? healthMessage;
  final String? errorMessage;

  DiagnosticsState copyWith({
    bool? isChecking,
    String? healthMessage,
    String? errorMessage,
    bool clearResult = false,
  }) {
    return DiagnosticsState(
      backendUrl: backendUrl,
      isChecking: isChecking ?? this.isChecking,
      healthMessage: clearResult ? null : healthMessage ?? this.healthMessage,
      errorMessage: clearResult ? null : errorMessage ?? this.errorMessage,
    );
  }
}

class DiagnosticsController extends ValueNotifier<DiagnosticsState> {
  DiagnosticsController({
    AppConfig config = AppConfig.dev,
    HealthCheck? healthCheck,
  })  : _healthCheck = healthCheck ??
            (() async {
              await ApiClient(
                baseUrl: config.backendBaseUrl,
                defaultHeaders: config.publicHeaders,
              ).getJson('/categories', timeout: _diagnosticsTimeout);
              return 'Backend responded.';
            }),
        super(DiagnosticsState(backendUrl: config.backendBaseUrl));

  final HealthCheck _healthCheck;

  Future<void> testConnection() async {
    value = value.copyWith(isChecking: true, clearResult: true);
    try {
      final response = await _healthCheck();
      value = value.copyWith(
        isChecking: false,
        healthMessage:
            response.trim().isEmpty ? 'Backend responded.' : response,
      );
    } catch (error, stackTrace) {
      debugPrint('Diagnostics health check failed: $error\n$stackTrace');
      value = value.copyWith(
        isChecking: false,
        errorMessage: 'Could not reach the backend.',
      );
    }
  }
}
