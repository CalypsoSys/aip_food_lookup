import 'package:flutter/foundation.dart';

import '../data/feedback_api.dart';
import '../models/feedback_request.dart';

class FeedbackState {
  const FeedbackState({
    this.isSubmitting = false,
    this.errorMessage,
  });

  final bool isSubmitting;
  final String? errorMessage;

  FeedbackState copyWith({
    bool? isSubmitting,
    String? errorMessage,
    bool clearError = false,
  }) {
    return FeedbackState(
      isSubmitting: isSubmitting ?? this.isSubmitting,
      errorMessage: clearError ? null : errorMessage ?? this.errorMessage,
    );
  }
}

class FeedbackController extends ValueNotifier<FeedbackState> {
  FeedbackController({FeedbackApi? feedbackApi})
      : _feedbackApi = feedbackApi ?? FeedbackApi(),
        super(const FeedbackState());

  final FeedbackApi _feedbackApi;

  String? validate({
    required String name,
    required String email,
    required String message,
  }) {
    final trimmedName = name.trim();
    final trimmedEmail = email.trim();
    final trimmedMessage = message.trim();
    if (trimmedMessage.isEmpty) {
      return 'Message is required.';
    }
    if (trimmedName.isEmpty && trimmedEmail.isEmpty) {
      return 'Enter a name or email so we know who sent this.';
    }
    if (trimmedEmail.isNotEmpty && !_looksLikeEmail(trimmedEmail)) {
      return 'Enter a valid email address or leave it blank.';
    }
    return null;
  }

  Future<bool> submit({
    required String name,
    required String email,
    required String subject,
    required String message,
  }) async {
    final validationMessage = validate(
      name: name,
      email: email,
      message: message,
    );
    if (validationMessage != null) {
      value = value.copyWith(errorMessage: validationMessage);
      return false;
    }
    if (value.isSubmitting) {
      return false;
    }

    value = value.copyWith(isSubmitting: true, clearError: true);
    try {
      await _feedbackApi.submit(
        FeedbackRequest(
          name: name.trim(),
          email: email.trim(),
          subject: subject.trim().isEmpty ? 'App feedback' : subject.trim(),
          message: message.trim(),
        ),
      );
      value = value.copyWith(isSubmitting: false);
      return true;
    } catch (error, stackTrace) {
      debugPrint('Feedback failed: $error\n$stackTrace');
      value = value.copyWith(
        isSubmitting: false,
        errorMessage: 'Feedback could not be sent.',
      );
      return false;
    }
  }

  bool _looksLikeEmail(String email) {
    final emailPattern = RegExp(r'^[^@\s]+@[^@\s]+\.[^@\s]+$');
    return emailPattern.hasMatch(email);
  }
}
