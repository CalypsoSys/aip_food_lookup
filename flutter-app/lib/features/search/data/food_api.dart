import '../../../app/config.dart';
import '../../../core/networking/api_client.dart';
import '../models/search_result.dart';
import '../models/suggest_food.dart';

class FoodApi {
  FoodApi({
    ApiClient? client,
    AppConfig config = AppConfig.dev,
  }) : _client = client ?? ApiClient(baseUrl: config.backendBaseUrl);

  final ApiClient _client;

  Future<SearchResult> search(String text, String searchType) async {
    final json = await _client.getJson(
      '/search',
      query: {
        'key': text,
        'type': normalizeSearchType(searchType),
      },
    );
    return SearchResult.fromJson(json);
  }

  Future<void> suggest(SuggestFoodRequest request) {
    return _client.postJson('/suggest', request.toJson());
  }

  Future<SearchResult> categories() async {
    final json = await _client.getJson('/categories');
    return SearchResult.fromJson(json);
  }

  Future<SearchResult> subcategory(String category, String subcategory) async {
    final json = await _client.getJson(
      '/subcategory',
      query: {
        'cat': category,
        'sub': normalizeSubcategory(subcategory),
      },
    );
    return SearchResult.fromJson(json);
  }
}

String normalizeSearchType(String value) {
  return value.replaceAll(' ', '').toLowerCase();
}

String normalizeSubcategory(String value) {
  return value.replaceAll(' and ', '_').toLowerCase();
}
