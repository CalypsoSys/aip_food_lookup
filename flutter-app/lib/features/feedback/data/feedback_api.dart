import '../../../app/config.dart';
import '../../../core/networking/api_client.dart';
import '../models/feedback_request.dart';

class FeedbackApi {
  FeedbackApi({
    ApiClient? client,
    AppConfig config = AppConfig.dev,
  }) : _client =
            client ??
            ApiClient(
              baseUrl: config.backendBaseUrl,
              defaultHeaders: config.publicHeaders,
            );

  final ApiClient _client;

  Future<void> submit(FeedbackRequest request) {
    return _client.postJson('/feedback', request.toJson());
  }
}
