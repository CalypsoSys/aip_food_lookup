import 'package:flutter/material.dart';

import '../../../ads/ad_banner.dart';
import '../../../widgets/status_card.dart';
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
                const SizedBox(height: 8),
                Text(
                  'Type at least 3 characters. Results update after you pause typing.',
                  style: Theme.of(context).textTheme.bodySmall,
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
                  StatusCard(
                    title: state.errorMessage!,
                    subtitle:
                        'Check that the Go backend is running and your phone can reach this PC on port 8080.',
                    icon: Icons.wifi_off_outlined,
                  ),
                ],
                const SizedBox(height: 16),
                _ResultSection(
                  title: 'Allowed on AIP:',
                  fallback: _fallbackText(
                    state,
                    suggestion: 'Use the plus button to suggest as allowed.',
                  ),
                  items: state.result.allowed,
                  hasSearched: state.hasSearched,
                ),
                const SizedBox(height: 12),
                _ResultSection(
                  title: 'NOT Allowed on AIP:',
                  fallback: _fallbackText(
                    state,
                    suggestion:
                        'Use the minus button to suggest as not allowed.',
                  ),
                  items: state.result.notAllowed,
                  hasSearched: state.hasSearched,
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
    final query = _controller.value.query.trim();
    final ok = await _controller.suggestCurrentFood(allowed: allowed);
    if (!context.mounted) {
      return;
    }

    final message = ok
        ? 'Thanks. We will review "$query" and add it to the catalog when appropriate.'
        : 'Enter at least 3 characters and make sure the backend is reachable.';
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

  String _fallbackText(
    feature.SearchState state, {
    required String suggestion,
  }) {
    if (state.query.trim().length < 3) {
      return 'Enter a food name to search.';
    }
    if (state.isLoading) {
      return 'Searching...';
    }
    if (state.hasSearched) {
      return 'No matches found. $suggestion';
    }
    return suggestion;
  }
}

class _ResultSection extends StatelessWidget {
  const _ResultSection({
    required this.title,
    required this.fallback,
    required this.items,
    required this.hasSearched,
  });

  final String title;
  final String fallback;
  final List<String> items;
  final bool hasSearched;

  @override
  Widget build(BuildContext context) {
    if (items.isEmpty) {
      return Card(
        child: ListTile(
          title: Text(hasSearched ? '$title 0' : title),
          subtitle: Text(fallback),
        ),
      );
    }

    return Column(
      crossAxisAlignment: CrossAxisAlignment.start,
      children: [
        Text(
          '$title ${items.length}',
          style: Theme.of(context).textTheme.titleMedium,
        ),
        const SizedBox(height: 6),
        ...items.map((item) => Card(child: ListTile(title: Text(item)))),
      ],
    );
  }
}
