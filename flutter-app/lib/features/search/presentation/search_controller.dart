import 'dart:async';

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
    this.isSuggesting = false,
    this.errorMessage,
    this.hasSearched = false,
  });

  final String query;
  final String searchType;
  final SearchResult result;
  final bool isLoading;
  final bool isSuggesting;
  final String? errorMessage;
  final bool hasSearched;

  SearchState copyWith({
    String? query,
    String? searchType,
    SearchResult? result,
    bool? isLoading,
    bool? isSuggesting,
    String? errorMessage,
    bool? hasSearched,
    bool clearError = false,
  }) {
    return SearchState(
      query: query ?? this.query,
      searchType: searchType ?? this.searchType,
      result: result ?? this.result,
      isLoading: isLoading ?? this.isLoading,
      isSuggesting: isSuggesting ?? this.isSuggesting,
      errorMessage: clearError ? null : errorMessage ?? this.errorMessage,
      hasSearched: hasSearched ?? this.hasSearched,
    );
  }
}

class SearchController extends ValueNotifier<SearchState> {
  SearchController({FoodApi? foodApi})
      : _foodApi = foodApi ?? FoodApi(),
        super(const SearchState());

  final FoodApi _foodApi;
  Timer? _debounce;
  int _searchGeneration = 0;

  void updateQuery(String query) {
    value = value.copyWith(query: query, clearError: true);
    _debounce?.cancel();
    _debounce = Timer(const Duration(milliseconds: 350), _runSearchIfValid);
  }

  Future<void> updateSearchType(String searchType) async {
    value = value.copyWith(searchType: searchType, clearError: true);
    _debounce?.cancel();
    await _runSearchIfValid();
  }

  Future<bool> suggestCurrentFood({required bool allowed}) async {
    final query = value.query.trim();
    if (!_hasSearchMinimum(query)) {
      return false;
    }
    if (value.isSuggesting) {
      return false;
    }

    value = value.copyWith(isSuggesting: true, clearError: true);
    try {
      await _foodApi.suggest(
        SuggestFoodRequest(inputText: query, allowed: allowed),
      );
      value = value.copyWith(isSuggesting: false);
      return true;
    } catch (error, stackTrace) {
      debugPrint('Suggestion failed: $error\n$stackTrace');
      value = value.copyWith(
        isSuggesting: false,
        errorMessage: 'Suggestion could not be made.',
      );
      return false;
    }
  }

  Future<void> _runSearchIfValid() async {
    final query = value.query.trim();
    if (!_hasSearchMinimum(query)) {
      _searchGeneration++;
      value = value.copyWith(
        result: const SearchResult(allowed: [], notAllowed: []),
        isLoading: false,
        hasSearched: false,
        clearError: true,
      );
      return;
    }

    final generation = ++_searchGeneration;
    value = value.copyWith(
      isLoading: true,
      hasSearched: true,
      clearError: true,
    );
    try {
      final result = await _foodApi.search(query, value.searchType);
      if (generation != _searchGeneration) {
        return;
      }
      value = value.copyWith(result: result, isLoading: false);
    } catch (error, stackTrace) {
      if (generation != _searchGeneration) {
        return;
      }
      debugPrint('Search failed: $error\n$stackTrace');
      value = value.copyWith(
        isLoading: false,
        errorMessage: 'Could not reach the food lookup API.',
      );
    }
  }

  bool _hasSearchMinimum(String query) => query.length > 2;

  @override
  void dispose() {
    _debounce?.cancel();
    super.dispose();
  }
}
