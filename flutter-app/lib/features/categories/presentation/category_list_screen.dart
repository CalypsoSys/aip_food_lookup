import 'package:flutter/material.dart';

import '../../../ads/ad_banner.dart';

class CategoryListScreen extends StatelessWidget {
  const CategoryListScreen({
    super.key,
    required this.category,
    required this.subcategory,
  });

  final String category;
  final String subcategory;

  @override
  Widget build(BuildContext context) {
    return Scaffold(
      appBar: AppBar(title: Text('Category: $category')),
      body: SafeArea(
        child: ListView(
          padding: const EdgeInsets.all(12),
          children: [
            Text(
              '$category: $subcategory',
              style: Theme.of(context).textTheme.titleMedium,
            ),
            const SizedBox(height: 12),
            const Card(
              child: ListTile(
                title: Text('Food list endpoint pending'),
                subtitle: Text('This route is scaffolded for milestone 1.'),
              ),
            ),
            const SizedBox(height: 16),
            const AdBannerPlaceholder(),
          ],
        ),
      ),
    );
  }
}
