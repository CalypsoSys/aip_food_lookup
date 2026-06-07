import 'package:flutter/material.dart';

import '../ads/ad_banner.dart';
import '../features/about/presentation/about_screen.dart';
import '../features/categories/presentation/categories_screen.dart';
import '../features/search/presentation/search_screen.dart';

class AppRouter extends StatefulWidget {
  const AppRouter({super.key, this.showAds = true});

  final bool showAds;

  @override
  State<AppRouter> createState() => _AppRouterState();
}

class _AppRouterState extends State<AppRouter> {
  static const _searchIndex = 1;

  int _selectedIndex = _searchIndex;

  @override
  Widget build(BuildContext context) {
    final screens = <Widget>[
      const CategoriesScreen(),
      const SearchScreen(),
      const AboutScreen(),
    ];

    return Scaffold(
      body: IndexedStack(index: _selectedIndex, children: screens),
      bottomNavigationBar: Column(
        mainAxisSize: MainAxisSize.min,
        children: [
          if (widget.showAds && _selectedIndex == _searchIndex)
            const SafeArea(
              top: false,
              bottom: false,
              child: SizedBox(
                height: 66,
                child: Center(child: SearchAdBanner()),
              ),
            ),
          NavigationBar(
            selectedIndex: _selectedIndex,
            onDestinationSelected: (index) {
              setState(() => _selectedIndex = index);
            },
            destinations: const [
              NavigationDestination(
                icon: Icon(Icons.category_outlined),
                label: 'Categories',
              ),
              NavigationDestination(
                icon: _SearchNavLogo(isSelected: false),
                selectedIcon: _SearchNavLogo(isSelected: true),
                label: 'Search',
              ),
              NavigationDestination(
                icon: Icon(Icons.info_outline),
                label: 'About',
              ),
            ],
          ),
        ],
      ),
    );
  }
}

class _SearchNavLogo extends StatelessWidget {
  const _SearchNavLogo({required this.isSelected});

  final bool isSelected;

  @override
  Widget build(BuildContext context) {
    final colorScheme = Theme.of(context).colorScheme;
    final size = isSelected ? 48.0 : 42.0;
    return Container(
      width: size,
      height: size,
      padding: const EdgeInsets.all(4),
      decoration: BoxDecoration(
        shape: BoxShape.circle,
        color: colorScheme.surface,
        border: Border.all(
          color: isSelected ? colorScheme.primary : colorScheme.outlineVariant,
          width: isSelected ? 2 : 1,
        ),
        boxShadow: [
          if (isSelected)
            BoxShadow(
              color: colorScheme.shadow.withValues(alpha: 0.16),
              blurRadius: 10,
              offset: const Offset(0, 4),
            ),
        ],
      ),
      child: ClipOval(
        child: Image.asset(
          'assets/identity/adaptive_icon.png',
          fit: BoxFit.cover,
        ),
      ),
    );
  }
}
