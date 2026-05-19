import 'package:flutter/material.dart';

import '../../../ads/ad_banner.dart';
import 'category_details_screen.dart';

class CategoriesScreen extends StatelessWidget {
  const CategoriesScreen({super.key});

  static const categories = ['Allowed', 'Not Allowed'];

  @override
  Widget build(BuildContext context) {
    return Scaffold(
      appBar: AppBar(title: const Text('Categories')),
      body: SafeArea(
        child: ListView(
          padding: const EdgeInsets.all(12),
          children: [
            for (final category in categories)
              Card(
                child: ListTile(
                  title: Text(category),
                  trailing: const Icon(Icons.chevron_right),
                  onTap: () {
                    Navigator.of(context).push(
                      MaterialPageRoute<void>(
                        builder: (_) => CategoryDetailsScreen(
                          category: category,
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
