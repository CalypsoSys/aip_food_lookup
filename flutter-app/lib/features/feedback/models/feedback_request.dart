class FeedbackRequest {
  const FeedbackRequest({
    required this.name,
    required this.email,
    required this.subject,
    required this.message,
    this.source = 'mobile',
  });

  final String name;
  final String email;
  final String subject;
  final String message;
  final String source;

  Map<String, dynamic> toJson() {
    return {
      'name': name,
      'email': email,
      'subject': subject,
      'message': message,
      'source': source,
    };
  }
}
