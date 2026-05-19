import 'package:flutter/material.dart';

import '../../../ads/ad_banner.dart';
import 'search_controller.dart' as feature;

class SearchScreen extends StatefulWidget {
  const SearchScreen({super.key});

  @override
  State<SearchScreen> createState() => _SearchScreenState();
}

class _SearchScreenState extends State<SearchScreen> {
  late final feature.SearchController _controller;

  @override
  void initState() {
    super.initState();
    _controller = feature.SearchController();
  }

  @override
  void dispose() {
    _controller.dispose();
    super.dispose();
  }

  @override
  Widget build(BuildContext context) {
    return ValueListenableBuilder<feature.SearchState>(
      valueListenable: _controller,
      builder: (context, state, _) {
        return Scaffold(
          appBar: AppBar(
            title: const Text('AIP Food Lookup'),
            actions: [
              IconButton(
                tooltip: 'Suggest allowed',
                onPressed: () => _suggest(context, allowed: true),
                icon: const Icon(Icons.add_comment_outlined),
              ),
              IconButton(
                tooltip: 'Suggest not allowed',
                onPressed: () => _suggest(context, allowed: false),
                icon: const Icon(Icons.comments_disabled_outlined),
              ),
            ],
          ),
          body: SafeArea(
            child: ListView(
              padding: const EdgeInsets.all(12),
              children: [
                TextField(
                  decoration: const InputDecoration(
                    labelText: 'Food',
                    hintText: 'Enter a food to check',
                  ),
                  textInputAction: TextInputAction.search,
                  onChanged: _controller.updateQuery,
                ),
                const SizedBox(height: 12),
                DropdownButtonFormField<String>(
                  initialValue: state.searchType,
                  decoration: const InputDecoration(labelText: 'Search type'),
                  items: feature.searchTypes
                      .map(
                        (type) => DropdownMenuItem(
                          value: type,
                          child: Text(type),
                        ),
                      )
                      .toList(),
                  onChanged: (value) {
                    if (value != null) {
                      _controller.updateSearchType(value);
                    }
                  },
                ),
                if (state.isLoading) ...[
                  const SizedBox(height: 16),
                  const LinearProgressIndicator(),
                ],
                if (state.errorMessage != null) ...[
                  const SizedBox(height: 16),
                  _ErrorBanner(message: state.errorMessage!),
                ],
                const SizedBox(height: 16),
                _ResultSection(
                  title: 'Allowed on AIP:',
                  fallback: 'Swipe right to suggest as allowed',
                  items: state.result.allowed,
                ),
                const SizedBox(height: 12),
                _ResultSection(
                  title: 'NOT Allowed on AIP:',
                  fallback: 'Swipe left to suggest as NOT allowed',
                  items: state.result.notAllowed,
                ),
                const SizedBox(height: 16),
                const AdBannerPlaceholder(),
              ],
            ),
          ),
        );
      },
    );
  }

  Future<void> _suggest(BuildContext context, {required bool allowed}) async {
    final ok = await _controller.suggestCurrentFood(allowed: allowed);
    if (!context.mounted) {
      return;
    }

    final message = ok
        ? 'We will look at your suggestion promptly and add to our catalog.'
        : 'Enter at least 3 characters before suggesting a food.';
    await showDialog<void>(
      context: context,
      builder: (context) => AlertDialog(
        title: const Text('Attention'),
        content: Text(message),
        actions: [
          TextButton(
            onPressed: () => Navigator.of(context).pop(),
            child: const Text('OK'),
          ),
        ],
      ),
    );
  }
}

class _ResultSection extends StatelessWidget {
  const _ResultSection({
    required this.title,
    required this.fallback,
    required this.items,
  });

  final String title;
  final String fallback;
  final List<String> items;

  @override
  Widget build(BuildContext context) {
    if (items.isEmpty) {
      return Card(
        child: ListTile(
          title: Text(title),
          subtitle: Text(fallback),
        ),
      );
    }

    return Column(
      crossAxisAlignment: CrossAxisAlignment.start,
      children: [
        Text(title, style: Theme.of(context).textTheme.titleMedium),
        const SizedBox(height: 6),
        ...items.map((item) => Card(child: ListTile(title: Text(item)))),
      ],
    );
  }
}

class _ErrorBanner extends StatelessWidget {
  const _ErrorBanner({required this.message});

  final String message;

  @override
  Widget build(BuildContext context) {
    return Material(
      color: Theme.of(context).colorScheme.errorContainer,
      borderRadius: BorderRadius.circular(8),
      child: Padding(
        padding: const EdgeInsets.all(12),
        child: Text(
          message,
          style: TextStyle(
            color: Theme.of(context).colorScheme.onErrorContainer,
          ),
        ),
      ),
    );
  }
}
