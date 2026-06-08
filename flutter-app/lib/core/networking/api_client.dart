import 'dart:async';
import 'dart:convert';
import 'dart:io';

import 'api_exception.dart';

class ApiClient {
  ApiClient({
    required String baseUrl,
    Map<String, String> defaultHeaders = const {},
    Duration requestTimeout = const Duration(seconds: 5),
  })  : _baseUri = Uri.parse(baseUrl),
        _defaultHeaders = Map.unmodifiable(defaultHeaders),
        _requestTimeout = requestTimeout;

  final Uri _baseUri;
  final Map<String, String> _defaultHeaders;
  final Duration _requestTimeout;
  final HttpClient _client = HttpClient();

  Future<Map<String, dynamic>> getJson(
    String path, {
    Map<String, String>? query,
    Duration? timeout,
  }) async {
    return _withTimeout(() async {
      final uri = _buildUri(path, query);
      final request = await _client.getUrl(uri);
      _applyDefaultHeaders(request);
      request.headers.set(HttpHeaders.acceptHeader, 'application/json');
      final response = await request.close();
      return _decodeObjectResponse(response);
    }, timeout);
  }

  Future<void> postJson(
    String path,
    Map<String, dynamic> body, {
    Duration? timeout,
  }) async {
    return _withTimeout(() async {
      final uri = _buildUri(path);
      final request = await _client.postUrl(uri);
      _applyDefaultHeaders(request);
      request.headers.set(HttpHeaders.acceptHeader, 'application/json');
      request.headers.set(HttpHeaders.contentTypeHeader, 'application/json');
      request.write(jsonEncode(body));
      final response = await request.close();
      if (response.statusCode < 200 || response.statusCode >= 300) {
        final text = await response.transform(utf8.decoder).join();
        throw ApiException(text, statusCode: response.statusCode);
      }
    }, timeout);
  }

  Future<String> getText(String path, {Duration? timeout}) async {
    return _withTimeout(() async {
      final uri = _buildUri(path);
      final request = await _client.getUrl(uri);
      _applyDefaultHeaders(request);
      request.headers.set(HttpHeaders.acceptHeader, 'text/plain');
      final response = await request.close();
      final text = await response.transform(utf8.decoder).join();
      if (response.statusCode < 200 || response.statusCode >= 300) {
        throw ApiException(text, statusCode: response.statusCode);
      }
      return text;
    }, timeout);
  }

  Uri _buildUri(String path, [Map<String, String>? query]) {
    return buildApiUri(_baseUri, path, query);
  }

  void _applyDefaultHeaders(HttpClientRequest request) {
    for (final entry in _defaultHeaders.entries) {
      request.headers.set(entry.key, entry.value);
    }
  }

  Future<Map<String, dynamic>> _decodeObjectResponse(
    HttpClientResponse response,
  ) async {
    final text = await response.transform(utf8.decoder).join();
    if (response.statusCode < 200 || response.statusCode >= 300) {
      throw ApiException(text, statusCode: response.statusCode);
    }

    final decoded = jsonDecode(text);
    if (decoded is Map<String, dynamic>) {
      return decoded;
    }
    throw ApiException('Expected a JSON object response.');
  }

  Future<T> _withTimeout<T>(
    Future<T> Function() operation,
    Duration? timeout,
  ) {
    return operation().timeout(timeout ?? _requestTimeout);
  }
}

Uri buildApiUri(Uri baseUri, String path, [Map<String, String>? query]) {
  final basePath = baseUri.path.endsWith('/')
      ? baseUri.path.substring(0, baseUri.path.length - 1)
      : baseUri.path;
  final relativePath = path.startsWith('/') ? path.substring(1) : path;
  final combinedPath = relativePath.isEmpty
      ? (basePath.isEmpty ? '/' : '$basePath/')
      : '${basePath.isEmpty ? '' : basePath}/$relativePath';

  return baseUri.replace(
    path: combinedPath,
    queryParameters: query,
  );
}
