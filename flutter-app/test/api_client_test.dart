import 'dart:async';
import 'dart:io';

import 'package:aip_food_lookup/core/networking/api_client.dart';
import 'package:flutter_test/flutter_test.dart';

void main() {
  test('buildApiUri preserves API base path for endpoint requests', () {
    final uri = buildApiUri(
      Uri.parse('https://hashimojoe.com/api'),
      '/search',
      {
        'key': 'apple',
        'type': 'searchbytext',
      },
    );

    expect(
      uri.toString(),
      'https://hashimojoe.com/api/search?key=apple&type=searchbytext',
    );
  });

  test('buildApiUri preserves API base path for health checks', () {
    final uri = buildApiUri(Uri.parse('https://hashimojoe.com/api'), '/');

    expect(uri.toString(), 'https://hashimojoe.com/api/');
  });

  test('buildApiUri still supports root local backend URLs', () {
    final uri = buildApiUri(Uri.parse('http://10.0.2.2:8080'), '/categories');

    expect(uri.toString(), 'http://10.0.2.2:8080/categories');
  });

  test('getJson times out when the server does not respond', () async {
    final server = await HttpServer.bind(InternetAddress.loopbackIPv4, 0);
    addTearDown(() => server.close(force: true));
    server.listen((request) {});

    final client = ApiClient(
      baseUrl: 'http://${server.address.host}:${server.port}',
      requestTimeout: const Duration(milliseconds: 20),
    );

    await expectLater(
      client.getJson('/categories'),
      throwsA(isA<TimeoutException>()),
    );
  });
}
