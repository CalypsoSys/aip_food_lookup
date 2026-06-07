import 'package:flutter/material.dart';
import 'package:flutter/services.dart';

import 'search_controller.dart' as feature;

class SearchScreen extends StatefulWidget {
  const SearchScreen({super.key, this.controller});

  final feature.SearchController? controller;

  @override
  State<SearchScreen> createState() => _SearchScreenState();
}

class _SearchScreenState extends State<SearchScreen> {
  late final feature.SearchController _controller;
  late final TextEditingController _textController;
  late final bool _ownsController;

  @override
  void initState() {
    super.initState();
    _controller = widget.controller ?? feature.SearchController();
    _ownsController = widget.controller == null;
    _textController = TextEditingController();
  }

  @override
  void dispose() {
    _textController.dispose();
    if (_ownsController) {
      _controller.dispose();
    }
    super.dispose();
  }

  @override
  Widget build(BuildContext context) {
    return ValueListenableBuilder<feature.SearchState>(
      valueListenable: _controller,
      builder: (context, state, _) {
        final canSuggest = state.query.trim().length > 2 && !state.isSuggesting;
        final hasMatches = state.result.allowed.isNotEmpty ||
            state.result.notAllowed.isNotEmpty;
        final showSuggestions = state.query.trim().length > 2 &&
            state.errorMessage == null &&
            (state.isSuggesting ||
                (state.hasSearched && !state.isLoading && !hasMatches));
        return Scaffold(
          body: SafeArea(
            child: ListView(
              padding: const EdgeInsets.all(12),
              children: [
                TextField(
                  controller: _textController,
                  decoration: InputDecoration(
                    labelText: 'Search food',
                    hintText: 'Type a food to check the AIP catalog',
                    suffixIcon: state.query.isEmpty
                        ? null
                        : IconButton(
                            tooltip: 'Clear search',
                            onPressed: () {
                              _textController.clear();
                              _controller.clearQuery();
                            },
                            icon: const Icon(Icons.clear),
                          ),
                  ),
                  textInputAction: TextInputAction.search,
                  onChanged: _controller.updateQuery,
                ),
                const SizedBox(height: 12),
                _LookupSummary(state: state),
                if (state.isLoading) ...[
                  const SizedBox(height: 8),
                  const LinearProgressIndicator(),
                ],
                const SizedBox(height: 12),
                _ResultSection(
                  title: 'Allowed on AIP',
                  items: state.result.allowed,
                  status: _ResultStatus.allowed,
                  onItemTap: _copyResult,
                ),
                _ResultSection(
                  title: 'Not allowed on AIP',
                  items: state.result.notAllowed,
                  status: _ResultStatus.notAllowed,
                  onItemTap: _copyResult,
                ),
                if (showSuggestions) ...[
                  const SizedBox(height: 8),
                  _SuggestionActions(
                    canSuggest: canSuggest,
                    isSuggesting: state.isSuggesting,
                    onSuggestAllowed: () => _suggest(context, allowed: true),
                    onSuggestNotAllowed: () =>
                        _suggest(context, allowed: false),
                  ),
                ],
                if (state.recentSearches.isNotEmpty) ...[
                  const SizedBox(height: 16),
                  Text(
                    'Recent searches',
                    style: Theme.of(context).textTheme.titleSmall,
                  ),
                  const SizedBox(height: 6),
                  SizedBox(
                    height: 42,
                    child: Row(
                      children: [
                        Expanded(
                          child: ListView.separated(
                            scrollDirection: Axis.horizontal,
                            itemCount: state.recentSearches.length,
                            separatorBuilder: (_, __) =>
                                const SizedBox(width: 8),
                            itemBuilder: (context, index) {
                              final recentSearch = state.recentSearches[index];
                              return ActionChip(
                                label: Text(recentSearch),
                                avatar: const Icon(Icons.history, size: 18),
                                onPressed: () {
                                  _textController.text = recentSearch;
                                  _controller.selectRecentSearch(recentSearch);
                                },
                              );
                            },
                          ),
                        ),
                        const SizedBox(width: 4),
                        IconButton(
                          tooltip: 'Clear recent searches',
                          onPressed: _controller.clearRecentSearches,
                          icon: const Icon(Icons.delete_outline),
                        ),
                      ],
                    ),
                  ),
                ],
                const SizedBox(height: 12),
                Theme(
                  data: Theme.of(context).copyWith(
                    dividerColor: Colors.transparent,
                  ),
                  child: ExpansionTile(
                    tilePadding: EdgeInsets.zero,
                    childrenPadding: const EdgeInsets.only(bottom: 8),
                    title: const Text('Advanced search'),
                    subtitle: Text(state.searchType),
                    children: [
                      DropdownButtonFormField<String>(
                        initialValue: state.searchType,
                        decoration:
                            const InputDecoration(labelText: 'Search type'),
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
                    ],
                  ),
                ),
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
        ? 'We will review "$query" promptly and add it to our catalog when appropriate.'
        : 'Enter at least 3 characters and make sure the backend URL is reachable.';
    ScaffoldMessenger.of(context).showSnackBar(
      SnackBar(content: Text(message)),
    );
  }

  Future<void> _copyResult(String result) async {
    await Clipboard.setData(ClipboardData(text: result));
    if (!mounted) {
      return;
    }
    ScaffoldMessenger.of(context).showSnackBar(
      SnackBar(content: Text('Copied "$result"')),
    );
  }
}

enum _ResultStatus { allowed, notAllowed }

class _LookupSummary extends StatelessWidget {
  const _LookupSummary({required this.state});

  final feature.SearchState state;

  @override
  Widget build(BuildContext context) {
    final colorScheme = Theme.of(context).colorScheme;
    final allowedCount = state.result.allowed.length;
    final notAllowedCount = state.result.notAllowed.length;
    final hasMinimum = state.query.trim().length > 2;
    final hasMatches = allowedCount > 0 || notAllowedCount > 0;

    late final String title;
    late final String subtitle;
    late final IconData messageIcon;
    late final Color messageColor;
    var showMessage = true;

    if (state.errorMessage != null) {
      messageIcon = Icons.wifi_off_outlined;
      messageColor = colorScheme.error;
      title = state.errorMessage!;
      subtitle = 'Check that your phone can reach the configured backend URL.';
    } else if (!hasMinimum) {
      messageIcon = Icons.search;
      messageColor = colorScheme.primary;
      title = 'Search the AIP catalog';
      subtitle = 'Type at least 3 characters. Results appear automatically.';
    } else if (state.isLoading || !state.hasSearched) {
      messageIcon = Icons.manage_search;
      messageColor = colorScheme.primary;
      title = 'Checking the catalog';
      subtitle = 'Looking for allowed and not allowed matches.';
    } else if (allowedCount > 0 && notAllowedCount == 0) {
      messageIcon = Icons.check_circle_outline;
      messageColor = Colors.green;
      showMessage = false;
      title =
          allowedCount == 1 ? 'Allowed match found' : 'Allowed matches found';
      subtitle = allowedCount == 1
          ? '1 allowed catalog match.'
          : '$allowedCount allowed catalog matches.';
    } else if (notAllowedCount > 0 && allowedCount == 0) {
      messageIcon = Icons.block;
      messageColor = colorScheme.error;
      showMessage = false;
      title = notAllowedCount == 1
          ? 'Not allowed match found'
          : 'Not allowed matches found';
      subtitle = notAllowedCount == 1
          ? '1 not allowed catalog match.'
          : '$notAllowedCount not allowed catalog matches.';
    } else if (hasMatches) {
      messageIcon = Icons.compare_arrows;
      messageColor = colorScheme.tertiary;
      title = 'Mixed catalog results';
      subtitle =
          '$allowedCount allowed and $notAllowedCount not allowed matches found.';
    } else {
      messageIcon = Icons.help_outline;
      messageColor = colorScheme.primary;
      title = 'No catalog match yet';
      subtitle = 'Suggest this food for review below.';
    }

    return Column(
      crossAxisAlignment: CrossAxisAlignment.stretch,
      children: [
        Row(
          children: [
            Expanded(
              child: _LookupCountTile(
                label: 'Allowed',
                count: allowedCount,
                icon: Icons.check_circle_outline,
                color: Colors.green,
                isActive: allowedCount > 0,
              ),
            ),
            const SizedBox(width: 8),
            Expanded(
              child: _LookupCountTile(
                label: 'Not allowed',
                count: notAllowedCount,
                icon: Icons.block,
                color: colorScheme.error,
                isActive: notAllowedCount > 0,
              ),
            ),
          ],
        ),
        if (showMessage) ...[
          const SizedBox(height: 8),
          _LookupMessage(
            icon: messageIcon,
            color: messageColor,
            title: title,
            subtitle: subtitle,
          ),
        ],
      ],
    );
  }
}

class _LookupCountTile extends StatelessWidget {
  const _LookupCountTile({
    required this.label,
    required this.count,
    required this.icon,
    required this.color,
    required this.isActive,
  });

  final String label;
  final int count;
  final IconData icon;
  final Color color;
  final bool isActive;

  @override
  Widget build(BuildContext context) {
    final colorScheme = Theme.of(context).colorScheme;
    final textTheme = Theme.of(context).textTheme;
    final tileColor = isActive
        ? color.withValues(alpha: 0.10)
        : colorScheme.surfaceContainerHighest.withValues(alpha: 0.55);
    final borderColor = isActive ? color : colorScheme.outlineVariant;
    final foregroundColor = isActive ? color : colorScheme.onSurfaceVariant;

    return Card(
      color: tileColor,
      shape: RoundedRectangleBorder(
        borderRadius: BorderRadius.circular(8),
        side: BorderSide(color: borderColor),
      ),
      child: ConstrainedBox(
        constraints: const BoxConstraints(minHeight: 58),
        child: Padding(
          padding: const EdgeInsets.symmetric(horizontal: 10, vertical: 8),
          child: Row(
            crossAxisAlignment: CrossAxisAlignment.center,
            children: [
              Icon(icon, color: foregroundColor, size: 20),
              const SizedBox(width: 6),
              Expanded(
                child: Text(
                  label,
                  maxLines: 2,
                  overflow: TextOverflow.ellipsis,
                  style: textTheme.labelLarge?.copyWith(
                    color: colorScheme.onSurface,
                  ),
                ),
              ),
              const SizedBox(width: 8),
              Text(
                '$count',
                style: textTheme.headlineSmall?.copyWith(
                  color: foregroundColor,
                  fontWeight: FontWeight.w700,
                ),
              ),
            ],
          ),
        ),
      ),
    );
  }
}

class _LookupMessage extends StatelessWidget {
  const _LookupMessage({
    required this.icon,
    required this.color,
    required this.title,
    required this.subtitle,
  });

  final IconData icon;
  final Color color;
  final String title;
  final String subtitle;

  @override
  Widget build(BuildContext context) {
    return Card(
      color: color.withValues(alpha: 0.07),
      child: ListTile(
        dense: true,
        leading: Icon(icon, color: color),
        title: Text(title),
        subtitle: Text(subtitle),
      ),
    );
  }
}

class _SuggestionActions extends StatelessWidget {
  const _SuggestionActions({
    required this.canSuggest,
    required this.isSuggesting,
    required this.onSuggestAllowed,
    required this.onSuggestNotAllowed,
  });

  final bool canSuggest;
  final bool isSuggesting;
  final VoidCallback onSuggestAllowed;
  final VoidCallback onSuggestNotAllowed;

  @override
  Widget build(BuildContext context) {
    return Column(
      crossAxisAlignment: CrossAxisAlignment.start,
      children: [
        Text(
          'Missing from the catalog?',
          style: Theme.of(context).textTheme.titleSmall,
        ),
        const SizedBox(height: 6),
        Text(
          'Send a suggestion only when the lookup does not find a clear match.',
          style: Theme.of(context).textTheme.bodySmall,
        ),
        const SizedBox(height: 8),
        SizedBox(
          width: double.infinity,
          child: OutlinedButton.icon(
            onPressed: canSuggest ? onSuggestAllowed : null,
            icon: const Icon(Icons.check_circle_outline),
            label: const Text('Suggest as allowed'),
          ),
        ),
        const SizedBox(height: 6),
        SizedBox(
          width: double.infinity,
          child: OutlinedButton.icon(
            onPressed: canSuggest ? onSuggestNotAllowed : null,
            icon: const Icon(Icons.block),
            label: const Text('Suggest as not allowed'),
          ),
        ),
        if (isSuggesting) ...[
          const SizedBox(height: 8),
          const LinearProgressIndicator(),
        ],
      ],
    );
  }
}

class _ResultSection extends StatelessWidget {
  const _ResultSection({
    required this.title,
    required this.items,
    required this.status,
    required this.onItemTap,
  });

  final String title;
  final List<String> items;
  final _ResultStatus status;
  final ValueChanged<String> onItemTap;

  @override
  Widget build(BuildContext context) {
    final colorScheme = Theme.of(context).colorScheme;
    final statusColor =
        status == _ResultStatus.allowed ? Colors.green : colorScheme.error;
    final containerColor = statusColor.withValues(alpha: 0.08);
    final icon = status == _ResultStatus.allowed
        ? Icons.check_circle_outline
        : Icons.block;

    if (items.isEmpty) {
      return const SizedBox.shrink();
    }

    return Column(
      crossAxisAlignment: CrossAxisAlignment.start,
      children: [
        Text(
          '$title (${items.length})',
          style: Theme.of(context).textTheme.titleMedium,
        ),
        const SizedBox(height: 6),
        ...items.map(
          (item) => Card(
            margin: const EdgeInsets.symmetric(vertical: 4),
            color: containerColor,
            child: ListTile(
              dense: true,
              visualDensity: VisualDensity.compact,
              leading: Icon(icon, color: statusColor),
              title: Text(item),
              trailing: const Icon(Icons.copy, size: 18),
              onTap: () => onItemTap(item),
            ),
          ),
        ),
        const SizedBox(height: 8),
      ],
    );
  }
}
