import 'dart:convert';
import 'dart:io';

import 'package:aip_food_lookup/app/config.dart';
import 'package:aip_food_lookup/features/diagnostics/presentation/diagnostics_controller.dart';
import 'package:flutter/foundation.dart';
import 'package:flutter_test/flutter_test.dart';

void main() {
  test('testConnection probes the categories API endpoint', () async {
    final server = await HttpServer.bind(InternetAddress.loopbackIPv4, 0);
    final paths = <String>[];
    addTearDown(() => server.close(force: true));
    server.listen((request) {
      paths.add(request.uri.path);
      request.response.headers.contentType = ContentType.json;
      request.response.write(jsonEncode({
        'allowed': <String>[],
        'not_allowed': <String>[],
      }));
      request.response.close();
    });

    final controller = DiagnosticsController(
      config: AppConfig(
        backendBaseUrl: 'http://${server.address.host}:${server.port}',
      ),
    );

    await controller.testConnection();

    expect(paths, ['/categories']);
    expect(controller.value.healthMessage, 'Backend responded.');
    expect(controller.value.errorMessage, isNull);
    expect(controller.value.isChecking, isFalse);
  });

  test('testConnection stores successful health response', () async {
    final controller = DiagnosticsController(
      config: const AppConfig(backendBaseUrl: 'http://example.test'),
      healthCheck: () async => 'AIP Food Lookup API',
    );

    await controller.testConnection();

    expect(controller.value.backendUrl, 'http://example.test');
    expect(controller.value.healthMessage, 'AIP Food Lookup API');
    expect(controller.value.errorMessage, isNull);
    expect(controller.value.isChecking, isFalse);
  });

  test('testConnection stores failed health response', () async {
    final originalDebugPrint = debugPrint;
    debugPrint = (message, {wrapWidth}) {};
    addTearDown(() => debugPrint = originalDebugPrint);
    final controller = DiagnosticsController(
      healthCheck: () async => throw Exception('offline'),
    );

    await controller.testConnection();

    expect(controller.value.healthMessage, isNull);
    expect(controller.value.errorMessage, 'Could not reach the backend.');
    expect(controller.value.isChecking, isFalse);
  });
}
