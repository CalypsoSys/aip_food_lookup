import 'package:flutter/foundation.dart';

import '../../search/data/food_api.dart';

class CategoryListState {
  const CategoryListState({
    this.items = const [],
    this.isLoading = false,
    this.errorMessage,
  });

  final List<String> items;
  final bool isLoading;
  final String? errorMessage;

  CategoryListState copyWith({
    List<String>? items,
    bool? isLoading,
    String? errorMessage,
    bool clearError = false,
  }) {
    return CategoryListState(
      items: items ?? this.items,
      isLoading: isLoading ?? this.isLoading,
      errorMessage: clearError ? null : errorMessage ?? this.errorMessage,
    );
  }
}

class CategoryController extends ValueNotifier<CategoryListState> {
  CategoryController({FoodApi? foodApi})
      : _foodApi = foodApi ?? FoodApi(),
        super(const CategoryListState());

  final FoodApi _foodApi;

  Future<void> loadTopLevelCategories(String category) async {
    value = value.copyWith(isLoading: true, clearError: true);
    try {
      final result = await _foodApi.categories();
      final items = category == 'Allowed' ? result.allowed : result.notAllowed;
      value = CategoryListState(items: List<String>.from(items)..sort());
    } catch (error, stackTrace) {
      debugPrint('Category load failed: $error\n$stackTrace');
      value = value.copyWith(
        isLoading: false,
        errorMessage: 'Could not load categories.',
      );
    }
  }

  Future<void> loadFoods(String category, String subcategory) async {
    value = value.copyWith(isLoading: true, clearError: true);
    try {
      final result = await _foodApi.subcategory(category, subcategory);
      final items = category == 'Allowed' ? result.allowed : result.notAllowed;
      value = CategoryListState(items: List<String>.from(items)..sort());
    } catch (error, stackTrace) {
      debugPrint('Food list load failed: $error\n$stackTrace');
      value = value.copyWith(
        isLoading: false,
        errorMessage: 'Could not load foods.',
      );
    }
  }
}
