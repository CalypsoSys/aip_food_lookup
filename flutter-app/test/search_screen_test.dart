import 'package:aip_food_lookup/features/search/data/food_api.dart';
import 'package:aip_food_lookup/features/search/models/search_result.dart';
import 'package:aip_food_lookup/features/search/presentation/search_controller.dart'
    as feature;
import 'package:aip_food_lookup/features/search/presentation/search_screen.dart';
import 'package:flutter/material.dart';
import 'package:flutter_test/flutter_test.dart';

void main() {
  testWidgets('starts with ingredient guidance and example chips',
      (tester) async {
    final controller = feature.SearchController(foodApi: _FakeFoodApi());
    addTearDown(controller.dispose);

    await tester.pumpWidget(
      MaterialApp(home: SearchScreen(controller: controller)),
    );

    expect(find.text('Search an ingredient'), findsOneWidget);
    expect(
      find.text(
        'Best for single ingredients. For prepared foods, check the ingredient list one item at a time.',
      ),
      findsOneWidget,
    );
    expect(find.text('Apple'), findsOneWidget);
    expect(find.text('Potato'), findsOneWidget);
    expect(find.text('Coconut milk'), findsOneWidget);
  });

  testWidgets('shows results before suggestion actions when matches exist',
      (tester) async {
    final controller = feature.SearchController(foodApi: _FakeFoodApi());
    addTearDown(controller.dispose);
    controller.value = const feature.SearchState(
      query: 'apple',
      hasSearched: true,
      result: SearchResult(allowed: ['Apples'], notAllowed: []),
    );

    await tester.pumpWidget(
      MaterialApp(home: SearchScreen(controller: controller)),
    );

    expect(find.text('Allowed'), findsOneWidget);
    expect(find.text('1'), findsOneWidget);
    expect(find.text('Allowed ingredients (1)'), findsOneWidget);
    expect(find.text('Apples'), findsOneWidget);
    expect(find.text('Allowed match found'), findsNothing);
    expect(find.text('Not seeing the ingredient you meant?'), findsNothing);
    expect(find.text('Suggest as allowed'), findsNothing);
    expect(find.text('Suggest as not allowed'), findsNothing);
  });

  testWidgets('shows a compact suggestion prompt for partial matches',
      (tester) async {
    final controller = feature.SearchController(foodApi: _FakeFoodApi());
    addTearDown(controller.dispose);
    controller.value = const feature.SearchState(
      query: 'turmeric powder',
      hasSearched: true,
      result: SearchResult(allowed: ['Turmeric'], notAllowed: []),
    );

    await tester.pumpWidget(
      MaterialApp(home: SearchScreen(controller: controller)),
    );

    expect(find.text('Turmeric'), findsOneWidget);
    expect(find.text('Not seeing the ingredient you meant?'), findsOneWidget);
    expect(find.text('Suggest "turmeric powder"'), findsOneWidget);
    expect(find.text('Suggest as allowed'), findsNothing);
    expect(find.text('Suggest as not allowed'), findsNothing);

    await tester.tap(find.text('Not seeing the ingredient you meant?'));
    await tester.pump();

    expect(
      find.widgetWithText(FilledButton, 'Suggest allowed'),
      findsOneWidget,
    );
    expect(
      find.widgetWithText(OutlinedButton, 'Suggest not allowed'),
      findsOneWidget,
    );
  });

  testWidgets('shows explicit suggestion actions after no matches',
      (tester) async {
    final controller = feature.SearchController(foodApi: _FakeFoodApi());
    addTearDown(controller.dispose);
    controller.value = const feature.SearchState(
      query: 'new food',
      hasSearched: true,
      result: SearchResult(allowed: [], notAllowed: []),
    );

    await tester.pumpWidget(
      MaterialApp(home: SearchScreen(controller: controller)),
    );

    expect(find.text('Prepared foods vary by recipe'), findsOneWidget);
    expect(
      find.text(
        'Search the ingredient list one item at a time, like cherry, wheat, sugar, potato, or oil.',
      ),
      findsOneWidget,
    );
    expect(find.text('Missing from the catalog?'), findsOneWidget);
    expect(find.text('Suggest as allowed'), findsOneWidget);
    expect(find.text('Suggest as not allowed'), findsOneWidget);
  });
}

class _FakeFoodApi extends FoodApi {}
