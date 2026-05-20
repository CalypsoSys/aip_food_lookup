import 'package:flutter/material.dart';

import '../../../ads/ad_banner.dart';
import 'category_controller.dart';
import 'category_list_screen.dart';

class CategoryDetailsScreen extends StatefulWidget {
  const CategoryDetailsScreen({super.key, required this.category});

  final String category;

  @override
  State<CategoryDetailsScreen> createState() => _CategoryDetailsScreenState();
}

class _CategoryDetailsScreenState extends State<CategoryDetailsScreen> {
  late final CategoryController _controller;

  @override
  void initState() {
    super.initState();
    _controller = CategoryController();
    _controller.loadTopLevelCategories(widget.category);
  }

  @override
  void dispose() {
    _controller.dispose();
    super.dispose();
  }

  @override
  Widget build(BuildContext context) {
    return Scaffold(
      appBar: AppBar(title: Text('Category: ${widget.category}')),
      body: SafeArea(
        child: ValueListenableBuilder<CategoryListState>(
          valueListenable: _controller,
          builder: (context, state, _) {
            return RefreshIndicator(
              onRefresh: () => _controller.loadTopLevelCategories(
                widget.category,
              ),
              child: ListView(
                padding: const EdgeInsets.all(12),
                children: [
                  if (state.isLoading) const LinearProgressIndicator(),
                  if (state.errorMessage != null)
                    _MessageCard(
                      title: state.errorMessage!,
                      actionLabel: 'Retry',
                      onAction: () => _controller.loadTopLevelCategories(
                        widget.category,
                      ),
                    ),
                  if (!state.isLoading &&
                      state.errorMessage == null &&
                      state.items.isEmpty)
                    const _MessageCard(title: 'No categories found.'),
                  for (final subcategory in state.items)
                    Card(
                      child: ListTile(
                        title: Text(subcategory),
                        trailing: const Icon(Icons.chevron_right),
                        onTap: () {
                          Navigator.of(context).push(
                            MaterialPageRoute<void>(
                              builder: (_) => CategoryListScreen(
                                category: widget.category,
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
            );
          },
        ),
      ),
    );
  }
}

class _MessageCard extends StatelessWidget {
  const _MessageCard({
    required this.title,
    this.actionLabel,
    this.onAction,
  });

  final String title;
  final String? actionLabel;
  final VoidCallback? onAction;

  @override
  Widget build(BuildContext context) {
    return Card(
      child: ListTile(
        title: Text(title),
        trailing: actionLabel == null
            ? null
            : TextButton(
                onPressed: onAction,
                child: Text(actionLabel!),
              ),
      ),
    );
  }
}
