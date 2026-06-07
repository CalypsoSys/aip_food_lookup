import 'dart:convert';

import 'package:flutter/services.dart';

import '../models/search_result.dart';

class LocalFoodCatalog {
  LocalFoodCatalog({
    required Map<String, List<String>> allowedByCategory,
    required Map<String, List<String>> notAllowedByCategory,
  })  : _allowedByCategory = _sortedCategoryMap(allowedByCategory),
        _notAllowedByCategory = _sortedCategoryMap(notAllowedByCategory);

  factory LocalFoodCatalog.fromJson(Map<String, dynamic> json) {
    return LocalFoodCatalog(
      allowedByCategory: _readCategoryMap(json['allowed']),
      notAllowedByCategory: _readCategoryMap(json['not_allowed']),
    );
  }

  static const defaultAssetPath = 'assets/catalog/catalog_snapshot.json';

  static Future<LocalFoodCatalog> load({
    AssetBundle? bundle,
    String assetPath = defaultAssetPath,
  }) async {
    final text = await (bundle ?? rootBundle).loadString(assetPath);
    final decoded = jsonDecode(text);
    if (decoded is! Map<String, dynamic>) {
      throw const FormatException('Catalog snapshot must be a JSON object.');
    }
    return LocalFoodCatalog.fromJson(decoded);
  }

  final Map<String, List<String>> _allowedByCategory;
  final Map<String, List<String>> _notAllowedByCategory;

  SearchResult categories() {
    return SearchResult(
      allowed: _allowedByCategory.keys.toList(),
      notAllowed: _notAllowedByCategory.keys.toList(),
    );
  }

  SearchResult subcategory(String category, String subcategory) {
    final categoryMap =
        category == 'Allowed' ? _allowedByCategory : _notAllowedByCategory;
    final selectedCategory = _findCategory(categoryMap, subcategory);
    if (selectedCategory == null) {
      return const SearchResult(allowed: [], notAllowed: []);
    }

    final items = List<String>.from(categoryMap[selectedCategory]!)
      ..sort(_compareFoodLabels);
    if (category == 'Allowed') {
      return SearchResult(allowed: items, notAllowed: const []);
    }
    return SearchResult(allowed: const [], notAllowed: items);
  }

  SearchResult search(String text, String searchType) {
    final query = _normalizeFoodKey(text);
    if (query.length < 3) {
      return const SearchResult(allowed: [], notAllowed: []);
    }

    return SearchResult(
      allowed: _searchCategoryMap(_allowedByCategory, query),
      notAllowed: _searchCategoryMap(_notAllowedByCategory, query),
    );
  }

  static Map<String, List<String>> _readCategoryMap(Object? value) {
    if (value is! Map<String, dynamic>) {
      return const {};
    }

    final output = <String, List<String>>{};
    for (final entry in value.entries) {
      final items = entry.value;
      if (items is! List) {
        continue;
      }
      output[entry.key] = items.whereType<String>().toList();
    }
    return output;
  }

  static Map<String, List<String>> _sortedCategoryMap(
    Map<String, List<String>> input,
  ) {
    final entries = input.entries.toList()
      ..sort((a, b) => a.key.compareTo(b.key));
    return Map.unmodifiable({
      for (final entry in entries)
        entry.key: List<String>.unmodifiable(
          List<String>.from(entry.value)..sort(_compareFoodLabels),
        ),
    });
  }

  static String? _findCategory(
    Map<String, List<String>> categoryMap,
    String subcategory,
  ) {
    final normalized = _normalizeSubcategory(subcategory);
    for (final category in categoryMap.keys) {
      if (_normalizeSubcategory(category) == normalized) {
        return category;
      }
    }
    return null;
  }

  static List<String> _searchCategoryMap(
    Map<String, List<String>> categoryMap,
    String query,
  ) {
    final matches = <String>{};
    for (final items in categoryMap.values) {
      for (final food in items) {
        if (_matchesFoodQuery(query, food)) {
          matches.add(cleanFoodLabel(food));
        }
      }
    }

    return matches.toList()..sort(_compareFoodLabels);
  }

  static bool _matchesFoodQuery(String query, String food) {
    final candidate = _normalizeFoodKey(food);
    if (candidate.isEmpty) {
      return false;
    }
    if (candidate.startsWith(query)) {
      return true;
    }

    final queryTokens = query.split(' ');
    final candidateTokens = candidate.split(' ');
    if (queryTokens.length > 1) {
      return queryTokens.every(
        (queryToken) => candidateTokens.any(
          (candidateToken) => candidateToken.startsWith(queryToken),
        ),
      );
    }

    return candidateTokens.any((token) => token.startsWith(query)) ||
        (query.length >= 4 && candidate.contains(query));
  }

  static String _normalizeFoodKey(String value) {
    return cleanFoodLabel(value)
        .toLowerCase()
        .replaceAll(RegExp(r'[^a-z0-9]+'), ' ')
        .trim()
        .replaceAll(RegExp(r'\s+'), ' ');
  }

  static String _normalizeSubcategory(String value) {
    return value.toLowerCase().replaceAll(' and ', '_').replaceAll(' ', '_');
  }

  static int _compareFoodLabels(String a, String b) {
    return a.toLowerCase().compareTo(b.toLowerCase());
  }
}
