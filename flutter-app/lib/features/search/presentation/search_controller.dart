import 'package:flutter/foundation.dart';

import '../data/food_api.dart';
import '../models/search_result.dart';
import '../models/suggest_food.dart';

const searchTypes = [
  'Search by Text and Sound',
  'Search by Text',
  'Search by Sound',
];
const defaultSearchType = 'Search by Text and Sound';

class SearchState {
  const SearchState({
    this.query = '',
    this.searchType = defaultSearchType,
    this.result = const SearchResult(allowed: [], notAllowed: []),
    this.isLoading = false,
    this.errorMessage,
  });

  final String query;
  final String searchType;
  final SearchResult result;
  final bool isLoading;
  final String? errorMessage;

  SearchState copyWith({
    String? query,
    String? searchType,
    SearchResult? result,
    bool? isLoading,
    String? errorMessage,
    bool clearError = false,
  }) {
    return SearchState(
      query: query ?? this.query,
      searchType: searchType ?? this.searchType,
      result: result ?? this.result,
      isLoading: isLoading ?? this.isLoading,
      errorMessage: clearError ? null : errorMessage ?? this.errorMessage,
    );
  }
}

class SearchController extends ValueNotifier<SearchState> {
  SearchController({FoodApi? foodApi})
      : _foodApi = foodApi ?? FoodApi(),
        super(const SearchState());

  final FoodApi _foodApi;

  Future<void> updateQuery(String query) async {
    value = value.copyWith(query: query, clearError: true);
    await _runSearchIfValid();
  }

  Future<void> updateSearchType(String searchType) async {
    value = value.copyWith(searchType: searchType, clearError: true);
    await _runSearchIfValid();
  }

  Future<bool> suggestCurrentFood({required bool allowed}) async {
    final query = value.query.trim();
    if (!_hasSearchMinimum(query)) {
      return false;
    }

    try {
      await _foodApi.suggest(
        SuggestFoodRequest(inputText: query, allowed: allowed),
      );
      return true;
    } catch (error, stackTrace) {
      debugPrint('Suggestion failed: $error\n$stackTrace');
      value = value.copyWith(
        errorMessage: 'Suggestion could not be made.',
      );
      return false;
    }
  }

  Future<void> _runSearchIfValid() async {
    final query = value.query.trim();
    if (!_hasSearchMinimum(query)) {
      value = value.copyWith(
        result: const SearchResult(allowed: [], notAllowed: []),
        isLoading: false,
        clearError: true,
      );
      return;
    }

    value = value.copyWith(isLoading: true, clearError: true);
    try {
      final result = await _foodApi.search(query, value.searchType);
      value = value.copyWith(result: result, isLoading: false);
    } catch (error, stackTrace) {
      debugPrint('Search failed: $error\n$stackTrace');
      value = value.copyWith(
        isLoading: false,
        errorMessage: 'An unexpected error occurred.',
      );
    }
  }

  bool _hasSearchMinimum(String query) => query.length > 2;
}
