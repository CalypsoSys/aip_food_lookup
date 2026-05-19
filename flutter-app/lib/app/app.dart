import 'package:flutter/material.dart';

import 'router.dart';
import 'theme.dart';

class AipFoodLookupApp extends StatelessWidget {
  const AipFoodLookupApp({super.key});

  @override
  Widget build(BuildContext context) {
    return MaterialApp(
      title: 'AIP Food Lookup',
      debugShowCheckedModeBanner: false,
      theme: buildAppTheme(),
      home: const AppRouter(),
    );
  }
}
