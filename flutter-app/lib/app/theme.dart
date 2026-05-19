import 'package:flutter/material.dart';

ThemeData buildAppTheme() {
  const primary = Color(0xFF512BD4);
  const magenta = Color(0xFFD600AA);

  return ThemeData(
    useMaterial3: true,
    colorScheme: ColorScheme.fromSeed(
      seedColor: primary,
      primary: primary,
      secondary: magenta,
    ),
    inputDecorationTheme: const InputDecorationTheme(
      border: OutlineInputBorder(),
    ),
  );
}
