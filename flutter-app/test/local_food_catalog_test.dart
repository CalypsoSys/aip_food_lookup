import 'package:aip_food_lookup/features/search/data/local_food_catalog.dart';
import 'package:flutter_test/flutter_test.dart';

void main() {
  test('loads categories from bundled-style catalog JSON', () {
    final catalog = LocalFoodCatalog.fromJson({
      'allowed': {
        'Fruits': ['Apples'],
        'Vegetables': ['Mushrooms (All)'],
      },
      'not_allowed': {
        'Grains': ['Wheat'],
      },
    });

    final result = catalog.categories();

    expect(result.allowed, ['Fruits', 'Vegetables']);
    expect(result.notAllowed, ['Grains']);
  });

  test('searches local foods with cleaned labels', () {
    final catalog = LocalFoodCatalog.fromJson({
      'allowed': {
        'Vegetables': ['Mushrooms (All)', 'Sweet Potato'],
      },
      'not_allowed': {
        'Nightshades': ['Potatoes'],
      },
    });

    final result = catalog.search('mush', 'Search by Text and Sound');

    expect(result.allowed, ['Mushrooms']);
    expect(result.notAllowed, isEmpty);
  });

  test('matches multi-word local foods by token prefix', () {
    final catalog = LocalFoodCatalog.fromJson({
      'allowed': {
        'Vegetables': ['Sweet Potato'],
      },
      'not_allowed': {
        'Nightshades': ['Potatoes'],
      },
    });

    final result = catalog.search('sweet pot', 'Search by Text and Sound');

    expect(result.allowed, ['Sweet Potato']);
    expect(result.notAllowed, isEmpty);
  });

  test('loads foods for route-normalized subcategories', () {
    final catalog = LocalFoodCatalog.fromJson({
      'allowed': {
        'Herbs and Spices': ['Turmeric'],
      },
      'not_allowed': {
        'Grains': ['Wheat'],
      },
    });

    final result = catalog.subcategory('Allowed', 'herbs_spices');

    expect(result.allowed, ['Turmeric']);
    expect(result.notAllowed, isEmpty);
  });

  test('bundled catalog snapshot has corrected category and food data',
      () async {
    TestWidgetsFlutterBinding.ensureInitialized();

    final catalog = await LocalFoodCatalog.load();
    final categories = catalog.categories();
    final fruits = catalog.subcategory('Allowed', 'Fruits');

    expect(categories.allowed, contains('Vegetables'));
    expect(categories.allowed, contains('Herbs and Spices'));
    expect(categories.allowed, isNot(contains('Vegtables')));
    expect(categories.notAllowed, isNot(contains('Herds and Spices')));
    expect(fruits.allowed, contains('Apples'));
    expect(fruits.allowed, isNot(contains('pples')));
  });
}
