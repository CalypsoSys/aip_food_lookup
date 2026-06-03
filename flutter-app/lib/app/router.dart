import 'package:flutter/material.dart';

import '../ads/ad_banner.dart';
import '../features/about/presentation/about_screen.dart';
import '../features/categories/presentation/categories_screen.dart';
import '../features/search/presentation/search_screen.dart';

class AppRouter extends StatefulWidget {
  const AppRouter({super.key});

  @override
  State<AppRouter> createState() => _AppRouterState();
}

class _AppRouterState extends State<AppRouter> {
  int _selectedIndex = 0;

  @override
  Widget build(BuildContext context) {
    final screens = <Widget>[
      const SearchScreen(),
      const CategoriesScreen(),
      const AboutScreen(),
    ];

    return Scaffold(
      body: IndexedStack(index: _selectedIndex, children: screens),
      bottomNavigationBar: Column(
        mainAxisSize: MainAxisSize.min,
        children: [
          if (_selectedIndex == 0)
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
                icon: Icon(Icons.search),
                label: 'Search',
              ),
              NavigationDestination(
                icon: Icon(Icons.category_outlined),
                label: 'Categories',
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
