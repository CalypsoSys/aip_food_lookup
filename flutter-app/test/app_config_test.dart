import 'package:aip_food_lookup/app/config.dart';
import 'package:flutter_test/flutter_test.dart';

void main() {
  test('publicHeaders includes mobile diagnostic headers', () {
    const config = AppConfig(
      backendBaseUrl: 'https://hashimojoe.com/api',
      clientName: 'android',
      appVersion: 'dev',
    );

    expect(config.publicHeaders, {
      'X-AIP-Client': 'android',
      'X-AIP-App-Version': 'dev',
    });
  });

  test('publicHeaders omits blank values', () {
    const config = AppConfig(
      backendBaseUrl: 'https://hashimojoe.com/api',
      clientName: ' ',
      appVersion: '',
    );

    expect(config.publicHeaders, isEmpty);
  });
}
