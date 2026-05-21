import 'dart:convert';
import 'dart:io';

import 'api_exception.dart';

class ApiClient {
  ApiClient({
    required String baseUrl,
    Map<String, String> defaultHeaders = const {},
  })  : _baseUri = Uri.parse(baseUrl),
        _defaultHeaders = Map.unmodifiable(defaultHeaders);

  final Uri _baseUri;
  final Map<String, String> _defaultHeaders;
  final HttpClient _client = HttpClient();

  Future<Map<String, dynamic>> getJson(
    String path, {
    Map<String, String>? query,
  }) async {
    final uri = _buildUri(path, query);
    final request = await _client.getUrl(uri);
    _applyDefaultHeaders(request);
    request.headers.set(HttpHeaders.acceptHeader, 'application/json');
    final response = await request.close();
    return _decodeObjectResponse(response);
  }

  Future<void> postJson(String path, Map<String, dynamic> body) async {
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
  }

  Future<String> getText(String path) async {
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
  }

  Uri _buildUri(String path, [Map<String, String>? query]) {
    final normalizedPath = path.startsWith('/') ? path : '/$path';
    return _baseUri.replace(
      path: normalizedPath,
      queryParameters: query,
    );
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
}
