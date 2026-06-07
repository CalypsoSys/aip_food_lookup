import 'package:aip_food_lookup/features/search/data/food_api.dart';
import 'package:aip_food_lookup/features/search/models/search_result.dart';
import 'package:aip_food_lookup/features/search/presentation/search_controller.dart'
    as feature;
import 'package:aip_food_lookup/features/search/presentation/search_screen.dart';
import 'package:flutter/material.dart';
import 'package:flutter_test/flutter_test.dart';

void main() {
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
    expect(find.text('Allowed on AIP (1)'), findsOneWidget);
    expect(find.text('Apples'), findsOneWidget);
    expect(find.text('Allowed match found'), findsNothing);
    expect(find.text('Suggest as allowed'), findsNothing);
    expect(find.text('Suggest as not allowed'), findsNothing);
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

    expect(find.text('No catalog match yet'), findsOneWidget);
    expect(find.text('Missing from the catalog?'), findsOneWidget);
    expect(find.text('Suggest as allowed'), findsOneWidget);
    expect(find.text('Suggest as not allowed'), findsOneWidget);
  });
}

class _FakeFoodApi extends FoodApi {}
