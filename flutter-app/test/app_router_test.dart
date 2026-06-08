import 'package:aip_food_lookup/app/router.dart';
import 'package:flutter/material.dart';
import 'package:flutter_test/flutter_test.dart';

void main() {
  testWidgets('starts on the centered Search destination', (tester) async {
    await tester.pumpWidget(
      const MaterialApp(home: AppRouter(showAds: false)),
    );

    final logo = find.byWidgetPredicate((widget) {
      return widget is Image &&
          widget.image is AssetImage &&
          (widget.image as AssetImage).assetName ==
              'assets/identity/adaptive_icon.png';
    });

    expect(find.text('Search an ingredient'), findsOneWidget);
    expect(logo, findsAtLeastNWidgets(1));
  });
}
