import 'package:flutter/material.dart';

class AssetHeader extends StatelessWidget {
  const AssetHeader({
    super.key,
    required this.assetName,
    required this.height,
    this.semanticLabel = 'AIP Food Lookup',
  });

  final String assetName;
  final double height;
  final String semanticLabel;

  @override
  Widget build(BuildContext context) {
    return Semantics(
      label: semanticLabel,
      image: true,
      child: Image.asset(
        assetName,
        height: height,
        fit: BoxFit.contain,
      ),
    );
  }
}
