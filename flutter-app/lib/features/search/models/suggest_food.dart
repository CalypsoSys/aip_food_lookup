class SuggestFoodRequest {
  const SuggestFoodRequest({
    required this.inputText,
    required this.allowed,
  });

  final String inputText;
  final bool allowed;

  Map<String, dynamic> toJson() {
    return {
      'inputText': inputText,
      'allowed': allowed,
    };
  }
}
