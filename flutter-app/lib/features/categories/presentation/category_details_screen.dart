import 'package:flutter/material.dart';

import '../../../ads/ad_banner.dart';
import 'category_list_screen.dart';

class CategoryDetailsScreen extends StatelessWidget {
  const CategoryDetailsScreen({super.key, required this.category});

  final String category;

  @override
  Widget build(BuildContext context) {
    const placeholderSubcategories = <String>[
      'Backend categories pending',
    ];

    return Scaffold(
      appBar: AppBar(title: Text('Category: $category')),
      body: SafeArea(
        child: ListView(
          padding: const EdgeInsets.all(12),
          children: [
            for (final subcategory in placeholderSubcategories)
              Card(
                child: ListTile(
                  title: Text(subcategory),
                  trailing: const Icon(Icons.chevron_right),
                  onTap: () {
                    Navigator.of(context).push(
                      MaterialPageRoute<void>(
                        builder: (_) => CategoryListScreen(
                          category: category,
                          subcategory: subcategory,
                        ),
                      ),
                    );
                  },
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
