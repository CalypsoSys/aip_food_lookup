import 'package:aip_food_lookup/features/search/data/food_api.dart';
import 'package:aip_food_lookup/features/search/models/search_result.dart';
import 'package:aip_food_lookup/features/search/models/suggest_food.dart';
import 'package:aip_food_lookup/features/search/presentation/search_controller.dart';
import 'package:fake_async/fake_async.dart';
import 'package:flutter/foundation.dart';
import 'package:flutter_test/flutter_test.dart';

void main() {
  test('updateQuery debounces backend search calls', () {
    fakeAsync((async) {
      final api = _FakeFoodApi();
      final controller = SearchController(foodApi: api);

      controller.updateQuery('app');
      controller.updateQuery('appl');
      async.elapse(const Duration(milliseconds: 349));

      expect(api.searchCalls, 0);

      async.elapse(const Duration(milliseconds: 1));

      expect(api.searchCalls, 1);
      expect(api.lastSearchText, 'appl');
    });
  });

  test('updateQuery clears results until MAUI minimum length is met', () {
    fakeAsync((async) {
      final api = _FakeFoodApi();
      final controller = SearchController(foodApi: api);

      controller.updateQuery('ap');
      async.elapse(const Duration(milliseconds: 350));

      expect(api.searchCalls, 0);
      expect(controller.value.result.allowed, isEmpty);
      expect(controller.value.result.notAllowed, isEmpty);
      expect(controller.value.hasSearched, isFalse);
    });
  });

  test('suggestCurrentFood requires at least 3 characters', () async {
    final api = _FakeFoodApi();
    final controller = SearchController(foodApi: api);

    controller.updateQuery('ap');
    final didSuggest = await controller.suggestCurrentFood(allowed: true);

    expect(didSuggest, isFalse);
    expect(api.suggestCalls, 0);
  });

  test('suggestCurrentFood submits valid suggestions', () async {
    final api = _FakeFoodApi();
    final controller = SearchController(foodApi: api);

    controller.updateQuery('apple');
    final didSuggest = await controller.suggestCurrentFood(allowed: true);

    expect(didSuggest, isTrue);
    expect(api.suggestCalls, 1);
    expect(api.lastSuggestion?.inputText, 'apple');
    expect(api.lastSuggestion?.allowed, isTrue);
    expect(controller.value.isSuggesting, isFalse);
    expect(controller.value.errorMessage, isNull);
  });

  test('suggestCurrentFood reports failed suggestions', () async {
    final originalDebugPrint = debugPrint;
    debugPrint = (message, {wrapWidth}) {};
    addTearDown(() => debugPrint = originalDebugPrint);
    final api = _FakeFoodApi(failSuggest: true);
    final controller = SearchController(foodApi: api);

    controller.updateQuery('apple');
    final didSuggest = await controller.suggestCurrentFood(allowed: false);

    expect(didSuggest, isFalse);
    expect(api.suggestCalls, 1);
    expect(controller.value.isSuggesting, isFalse);
    expect(controller.value.errorMessage, 'Suggestion could not be made.');
  });
}

class _FakeFoodApi extends FoodApi {
  _FakeFoodApi({this.failSuggest = false});

  final bool failSuggest;
  int searchCalls = 0;
  int suggestCalls = 0;
  String? lastSearchText;
  SuggestFoodRequest? lastSuggestion;

  @override
  Future<SearchResult> search(String text, String searchType) async {
    searchCalls++;
    lastSearchText = text;
    return const SearchResult(allowed: ['Apples'], notAllowed: []);
  }

  @override
  Future<void> suggest(SuggestFoodRequest request) async {
    suggestCalls++;
    lastSuggestion = request;
    if (failSuggest) {
      throw Exception('suggest failed');
    }
  }
}
