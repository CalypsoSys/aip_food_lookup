import 'package:aip_food_lookup/features/feedback/data/feedback_api.dart';
import 'package:aip_food_lookup/features/feedback/models/feedback_request.dart';
import 'package:aip_food_lookup/features/feedback/presentation/feedback_controller.dart';
import 'package:flutter/foundation.dart';
import 'package:flutter_test/flutter_test.dart';

void main() {
  test('validate requires a message', () {
    final controller = FeedbackController(feedbackApi: _FakeFeedbackApi());

    expect(
      controller.validate(name: 'Joe', email: '', message: ''),
      'Message is required.',
    );
  });

  test('validate requires name or email', () {
    final controller = FeedbackController(feedbackApi: _FakeFeedbackApi());

    expect(
      controller.validate(name: '', email: '', message: 'Hello'),
      'Enter a name or email so we know who sent this.',
    );
  });

  test('submit sends default subject and trimmed fields', () async {
    final api = _FakeFeedbackApi();
    final controller = FeedbackController(feedbackApi: api);

    final sent = await controller.submit(
      name: ' Joe ',
      email: '',
      subject: '',
      message: ' Nice app ',
    );

    expect(sent, isTrue);
    expect(api.submitCalls, 1);
    expect(api.lastRequest?.name, 'Joe');
    expect(api.lastRequest?.subject, 'App feedback');
    expect(api.lastRequest?.message, 'Nice app');
    expect(controller.value.errorMessage, isNull);
  });

  test('submit reports API failure', () async {
    final originalDebugPrint = debugPrint;
    debugPrint = (message, {wrapWidth}) {};
    addTearDown(() => debugPrint = originalDebugPrint);
    final api = _FakeFeedbackApi(shouldFail: true);
    final controller = FeedbackController(feedbackApi: api);

    final sent = await controller.submit(
      name: 'Joe',
      email: '',
      subject: 'Help',
      message: 'Broken',
    );

    expect(sent, isFalse);
    expect(controller.value.isSubmitting, isFalse);
    expect(controller.value.errorMessage, 'Feedback could not be sent.');
  });
}

class _FakeFeedbackApi extends FeedbackApi {
  _FakeFeedbackApi({this.shouldFail = false});

  final bool shouldFail;
  int submitCalls = 0;
  FeedbackRequest? lastRequest;

  @override
  Future<void> submit(FeedbackRequest request) async {
    submitCalls++;
    lastRequest = request;
    if (shouldFail) {
      throw Exception('feedback failed');
    }
  }
}
