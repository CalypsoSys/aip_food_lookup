import '../../../app/config.dart';
import '../../../core/networking/api_client.dart';
import '../models/search_result.dart';
import '../models/suggest_food.dart';
import 'local_food_catalog.dart';

typedef LocalFoodCatalogLoader = Future<LocalFoodCatalog> Function();

class FoodApi {
  FoodApi({
    ApiClient? client,
    AppConfig config = AppConfig.dev,
    LocalFoodCatalogLoader? fallbackCatalogLoader,
  })  : _client = client ??
            ApiClient(
              baseUrl: config.backendBaseUrl,
              defaultHeaders: config.publicHeaders,
            ),
        _fallbackCatalogLoader =
            fallbackCatalogLoader ?? (() => LocalFoodCatalog.load());

  final ApiClient _client;
  final LocalFoodCatalogLoader _fallbackCatalogLoader;
  Future<LocalFoodCatalog>? _fallbackCatalog;

  Future<SearchResult> search(String text, String searchType) async {
    try {
      final json = await _client.getJson(
        '/search',
        query: {
          'key': text,
          'type': normalizeSearchType(searchType),
        },
      );
      return SearchResult.fromJson(json);
    } catch (_) {
      return (await _loadFallbackCatalog()).search(text, searchType);
    }
  }

  Future<void> suggest(SuggestFoodRequest request) {
    return _client.postJson('/suggest', request.toJson());
  }

  Future<SearchResult> categories() async {
    try {
      final json = await _client.getJson('/categories');
      return SearchResult.fromJson(json);
    } catch (_) {
      return (await _loadFallbackCatalog()).categories();
    }
  }

  Future<SearchResult> subcategory(String category, String subcategory) async {
    try {
      final json = await _client.getJson(
        '/subcategory',
        query: {
          'cat': category,
          'sub': normalizeSubcategory(subcategory),
        },
      );
      return SearchResult.fromJson(json);
    } catch (_) {
      return (await _loadFallbackCatalog()).subcategory(category, subcategory);
    }
  }

  Future<LocalFoodCatalog> _loadFallbackCatalog() {
    return _fallbackCatalog ??= _fallbackCatalogLoader();
  }
}

String normalizeSearchType(String value) {
  return value.replaceAll(' ', '').toLowerCase();
}

String normalizeSubcategory(String value) {
  return value.replaceAll(' and ', '_').toLowerCase();
}
