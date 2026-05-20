import 'package:flutter/material.dart';

import '../../../ads/ad_banner.dart';
import '../../../widgets/status_card.dart';
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
                  Text(
                    state.isLoading
                        ? widget.category
                        : '${widget.category} (${state.items.length})',
                    style: Theme.of(context).textTheme.titleMedium,
                  ),
                  const SizedBox(height: 12),
                  if (state.isLoading) const LinearProgressIndicator(),
                  if (state.errorMessage != null)
                    StatusCard(
                      title: state.errorMessage!,
                      subtitle:
                          'Pull to refresh or check the backend connection.',
                      icon: Icons.wifi_off_outlined,
                      actionLabel: 'Retry',
                      onAction: () => _controller.loadTopLevelCategories(
                        widget.category,
                      ),
                    ),
                  if (!state.isLoading &&
                      state.errorMessage == null &&
                      state.items.isEmpty)
                    const StatusCard(
                      title: 'No categories found.',
                      icon: Icons.search_off_outlined,
                    ),
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
