import 'dart:convert';

import 'package:flutter/services.dart';

import '../models/search_result.dart';

class LocalFoodCatalog {
  LocalFoodCatalog({
    required Map<String, List<LocalFoodEntry>> allowedByCategory,
    required Map<String, List<LocalFoodEntry>> notAllowedByCategory,
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

  final Map<String, List<LocalFoodEntry>> _allowedByCategory;
  final Map<String, List<LocalFoodEntry>> _notAllowedByCategory;

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

    final items = categoryMap[selectedCategory]!
        .map((entry) => entry.name)
        .toList()
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

  static Map<String, List<LocalFoodEntry>> _readCategoryMap(Object? value) {
    if (value is! Map<String, dynamic>) {
      return const {};
    }

    final output = <String, List<LocalFoodEntry>>{};
    for (final entry in value.entries) {
      final items = entry.value;
      if (items is! List) {
        continue;
      }
      output[entry.key] =
          items.map(_readFoodEntry).whereType<LocalFoodEntry>().toList();
    }
    return output;
  }

  static LocalFoodEntry? _readFoodEntry(Object? value) {
    if (value is String) {
      return LocalFoodEntry(name: value, aliases: const []);
    }
    if (value is Map<String, dynamic> && value['name'] is String) {
      final aliases = value['aliases'];
      return LocalFoodEntry(
        name: value['name'] as String,
        aliases: aliases is List
            ? aliases.whereType<String>().toList()
            : const [],
      );
    }
    return null;
  }

  static Map<String, List<LocalFoodEntry>> _sortedCategoryMap(
    Map<String, List<LocalFoodEntry>> input,
  ) {
    final entries = input.entries.toList()
      ..sort((a, b) => a.key.compareTo(b.key));
    return Map.unmodifiable({
      for (final entry in entries)
        entry.key: List<LocalFoodEntry>.unmodifiable(
          List<LocalFoodEntry>.from(entry.value)
            ..sort((a, b) => _compareFoodLabels(a.name, b.name)),
        ),
    });
  }

  static String? _findCategory(
    Map<String, List<LocalFoodEntry>> categoryMap,
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
    Map<String, List<LocalFoodEntry>> categoryMap,
    String query,
  ) {
    final matches = <String>{};
    for (final items in categoryMap.values) {
      for (final food in items) {
        if (_matchesFoodQuery(query, food)) {
          matches.add(food.name);
        }
      }
    }

    return matches.toList()..sort(_compareFoodLabels);
  }

  static bool _matchesFoodQuery(String query, LocalFoodEntry food) {
    for (final candidateText in [food.name, ...food.aliases]) {
      final candidate = _normalizeFoodKey(candidateText);
      if (candidate.isEmpty) {
        continue;
      }
      if (candidate.startsWith(query)) {
        return true;
      }

      final queryTokens = query.split(' ');
      final candidateTokens = candidate.split(' ');
      if (queryTokens.length > 1 &&
          queryTokens.every(
            (queryToken) => candidateTokens.any(
              (candidateToken) => candidateToken.startsWith(queryToken),
            ),
          )) {
        return true;
      }

      if (queryTokens.length == 1 &&
          (candidateTokens.any((token) => token.startsWith(query)) ||
              (query.length >= 4 && candidate.contains(query)))) {
        return true;
      }
    }
    return false;
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

class LocalFoodEntry {
  const LocalFoodEntry({required this.name, required this.aliases});

  final String name;
  final List<String> aliases;
}
