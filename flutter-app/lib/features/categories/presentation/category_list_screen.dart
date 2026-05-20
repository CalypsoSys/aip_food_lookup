import 'package:flutter/material.dart';

import '../../../ads/ad_banner.dart';
import '../../../widgets/status_card.dart';
import 'category_controller.dart';

class CategoryListScreen extends StatefulWidget {
  const CategoryListScreen({
    super.key,
    required this.category,
    required this.subcategory,
  });

  final String category;
  final String subcategory;

  @override
  State<CategoryListScreen> createState() => _CategoryListScreenState();
}

class _CategoryListScreenState extends State<CategoryListScreen> {
  late final CategoryController _controller;

  @override
  void initState() {
    super.initState();
    _controller = CategoryController();
    _controller.loadFoods(widget.category, widget.subcategory);
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
              onRefresh: () => _controller.loadFoods(
                widget.category,
                widget.subcategory,
              ),
              child: ListView(
                padding: const EdgeInsets.all(12),
                children: [
                  Text(
                    state.isLoading
                        ? '${widget.category}: ${widget.subcategory}'
                        : '${widget.category}: ${widget.subcategory} (${state.items.length})',
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
                      onAction: () => _controller.loadFoods(
                        widget.category,
                        widget.subcategory,
                      ),
                    ),
                  if (!state.isLoading &&
                      state.errorMessage == null &&
                      state.items.isEmpty)
                    const StatusCard(
                      title: 'No foods found.',
                      icon: Icons.search_off_outlined,
                    ),
                  for (final food in state.items)
                    Card(child: ListTile(title: Text(food))),
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
