import 'dart:async';

import 'package:aip_food_lookup/core/networking/api_client.dart';
import 'package:aip_food_lookup/features/search/data/food_api.dart';
import 'package:aip_food_lookup/features/search/data/local_food_catalog.dart';
import 'package:aip_food_lookup/features/search/models/suggest_food.dart';
import 'package:flutter_test/flutter_test.dart';

void main() {
  test('search uses remote API when it responds', () async {
    final api = FoodApi(
      client: _RemoteApiClient({
        'allowed': ['Apples'],
        'not_allowed': <String>[],
      }),
      fallbackCatalogLoader: () async => _fallbackCatalog(),
    );

    final result = await api.search('apple', 'Search by Text');

    expect(result.allowed, ['Apples']);
    expect(result.notAllowed, isEmpty);
  });

  test('search falls back to local catalog when remote API fails', () async {
    final api = FoodApi(
      client: _FailingApiClient(),
      fallbackCatalogLoader: () async => _fallbackCatalog(),
    );

    final result = await api.search('apple', 'Search by Text');

    expect(result.allowed, ['Apples']);
    expect(result.notAllowed, isEmpty);
  });

  test('search uses the short fallback timeout for remote reads', () async {
    final client = _RecordingApiClient();
    final api = FoodApi(
      client: client,
      fallbackCatalogLoader: () async => _fallbackCatalog(),
    );

    final result = await api.search('apple', 'Search by Text');

    expect(client.lastGetTimeout, const Duration(milliseconds: 1200));
    expect(result.allowed, ['Apples']);
    expect(result.notAllowed, isEmpty);
  });

  test('categories fall back to local catalog when remote API fails', () async {
    final api = FoodApi(
      client: _FailingApiClient(),
      fallbackCatalogLoader: () async => _fallbackCatalog(),
    );

    final result = await api.categories();

    expect(result.allowed, ['Fruits']);
    expect(result.notAllowed, ['Grains']);
  });

  test('subcategories fall back to local catalog when remote API fails',
      () async {
    final api = FoodApi(
      client: _FailingApiClient(),
      fallbackCatalogLoader: () async => _fallbackCatalog(),
    );

    final result = await api.subcategory('Not Allowed', 'Grains');

    expect(result.allowed, isEmpty);
    expect(result.notAllowed, ['Wheat']);
  });

  test('suggestions do not use the short read timeout', () async {
    final client = _RecordingApiClient();
    final api = FoodApi(
      client: client,
      fallbackCatalogLoader: () async => _fallbackCatalog(),
    );

    await expectLater(
      api.suggest(
        const SuggestFoodRequest(inputText: 'cassava chips', allowed: true),
      ),
      throwsA(isA<TimeoutException>()),
    );

    expect(client.lastPostTimeout, isNull);
  });

  test('suggestions remain server-only', () async {
    final api = FoodApi(
      client: _FailingApiClient(),
      fallbackCatalogLoader: () async => _fallbackCatalog(),
    );

    await expectLater(
      api.suggest(
        const SuggestFoodRequest(inputText: 'cassava chips', allowed: true),
      ),
      throwsException,
    );
  });
}

LocalFoodCatalog _fallbackCatalog() {
  return LocalFoodCatalog.fromJson({
    'allowed': {
      'Fruits': ['Apples'],
    },
    'not_allowed': {
      'Grains': ['Wheat'],
    },
  });
}

class _RemoteApiClient extends ApiClient {
  _RemoteApiClient(this.response) : super(baseUrl: 'https://example.test/api');

  final Map<String, dynamic> response;

  @override
  Future<Map<String, dynamic>> getJson(
    String path, {
    Map<String, String>? query,
    Duration? timeout,
  }) async {
    return response;
  }
}

class _FailingApiClient extends ApiClient {
  _FailingApiClient() : super(baseUrl: 'https://example.test/api');

  @override
  Future<Map<String, dynamic>> getJson(
    String path, {
    Map<String, String>? query,
    Duration? timeout,
  }) async {
    throw Exception('offline');
  }

  @override
  Future<void> postJson(
    String path,
    Map<String, dynamic> body, {
    Duration? timeout,
  }) async {
    throw Exception('offline');
  }
}

class _RecordingApiClient extends ApiClient {
  _RecordingApiClient() : super(baseUrl: 'https://example.test/api');

  Duration? lastGetTimeout;
  Duration? lastPostTimeout;

  @override
  Future<Map<String, dynamic>> getJson(
    String path, {
    Map<String, String>? query,
    Duration? timeout,
  }) async {
    lastGetTimeout = timeout;
    throw TimeoutException('read timed out');
  }

  @override
  Future<void> postJson(
    String path,
    Map<String, dynamic> body, {
    Duration? timeout,
  }) async {
    lastPostTimeout = timeout;
    throw TimeoutException('write timed out');
  }
}
