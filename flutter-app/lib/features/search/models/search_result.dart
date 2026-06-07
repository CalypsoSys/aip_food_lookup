class SearchResult {
  const SearchResult({
    required this.allowed,
    required this.notAllowed,
  });

  factory SearchResult.fromJson(Map<String, dynamic> json) {
    return SearchResult(
      allowed: _stringList(json['allowed']),
      notAllowed: _stringList(json['not_allowed']),
    );
  }

  final List<String> allowed;
  final List<String> notAllowed;

  static List<String> _stringList(Object? value) {
    if (value is! List) {
      return const [];
    }
    return value.whereType<String>().map(cleanFoodLabel).toList();
  }
}

String cleanFoodLabel(String value) {
  return value
      .replaceFirst(RegExp(r'\s*\(all\)\s*$', caseSensitive: false), '')
      .trim();
}
