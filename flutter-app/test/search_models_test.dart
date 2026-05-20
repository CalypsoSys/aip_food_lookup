import 'package:aip_food_lookup/features/search/data/food_api.dart';
import 'package:aip_food_lookup/features/search/models/search_result.dart';
import 'package:aip_food_lookup/features/search/models/suggest_food.dart';
import 'package:flutter_test/flutter_test.dart';

void main() {
  test('SearchResult reads allowed and not_allowed lists', () {
    final result = SearchResult.fromJson({
      'allowed': ['apple'],
      'not_allowed': ['wheat'],
    });

    expect(result.allowed, ['apple']);
    expect(result.notAllowed, ['wheat']);
  });

  test('SuggestFoodRequest matches backend DTO shape', () {
    const request = SuggestFoodRequest(inputText: 'rice', allowed: true);

    expect(request.toJson(), {
      'inputText': 'rice',
      'allowed': true,
    });
  });

  test('Search type normalization matches MAUI behavior', () {
    expect(
      normalizeSearchType('Search by Text and Sound'),
      'searchbytextandsound',
    );
  });

  test('Subcategory normalization matches MAUI behavior', () {
    expect(normalizeSubcategory('Herbs and Spices'), 'herbs_spices');
  });
}
