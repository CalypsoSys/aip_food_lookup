import 'package:aip_food_lookup/features/categories/presentation/category_controller.dart';
import 'package:aip_food_lookup/features/search/data/food_api.dart';
import 'package:aip_food_lookup/features/search/models/search_result.dart';
import 'package:flutter_test/flutter_test.dart';

void main() {
  test('loadTopLevelCategories selects allowed categories', () async {
    final controller = CategoryController(
      foodApi: _FakeFoodApi(
        categoriesResult: const SearchResult(
          allowed: ['Fruits'],
          notAllowed: ['Grains'],
        ),
      ),
    );

    await controller.loadTopLevelCategories('Allowed');

    expect(controller.value.items, ['Fruits']);
    expect(controller.value.isLoading, isFalse);
    expect(controller.value.errorMessage, isNull);
  });

  test('loadFoods selects not allowed foods', () async {
    final controller = CategoryController(
      foodApi: _FakeFoodApi(
        subcategoryResult: const SearchResult(
          allowed: ['Apples'],
          notAllowed: ['Wheat'],
        ),
      ),
    );

    await controller.loadFoods('Not Allowed', 'Grains');

    expect(controller.value.items, ['Wheat']);
    expect(controller.value.isLoading, isFalse);
    expect(controller.value.errorMessage, isNull);
  });
}

class _FakeFoodApi extends FoodApi {
  _FakeFoodApi({
    this.categoriesResult = const SearchResult(allowed: [], notAllowed: []),
    this.subcategoryResult = const SearchResult(allowed: [], notAllowed: []),
  });

  final SearchResult categoriesResult;
  final SearchResult subcategoryResult;

  @override
  Future<SearchResult> categories() async {
    return categoriesResult;
  }

  @override
  Future<SearchResult> subcategory(String category, String subcategory) async {
    return subcategoryResult;
  }
}
