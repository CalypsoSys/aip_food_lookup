import 'package:flutter/material.dart';

import '../../../features/diagnostics/presentation/diagnostics_screen.dart';
import '../../../features/feedback/presentation/feedback_screen.dart';
import '../../../widgets/asset_header.dart';
import '../../../widgets/status_card.dart';

class AboutScreen extends StatelessWidget {
  const AboutScreen({super.key});

  @override
  Widget build(BuildContext context) {
    return Scaffold(
      appBar: AppBar(title: const Text('About')),
      body: SafeArea(
        child: ListView(
          padding: const EdgeInsets.all(12),
          children: [
            const AssetHeader(
              assetName: 'assets/identity/adaptive_icon.png',
              height: 72,
            ),
            const SizedBox(height: 12),
            Text(
              'AIP Food Lookup',
              textAlign: TextAlign.center,
              style: Theme.of(context).textTheme.headlineSmall,
            ),
            const SizedBox(height: 12),
            const StatusCard(
              title: 'Quick food lookup',
              subtitle:
                  'AIP Food Lookup helps you quickly check whether foods are commonly listed as allowed or not allowed for the Autoimmune Protocol diet.',
              icon: Icons.search,
            ),
            const StatusCard(
              title: 'Help improve the catalog',
              subtitle:
                  'Use Search for quick lookup, Categories to browse foods, and Suggestions or Feedback to help improve the catalog.',
              icon: Icons.volunteer_activism_outlined,
            ),
            const StatusCard(
              title: 'Informational only',
              subtitle:
                  'This app is informational only and is not medical advice.',
              icon: Icons.health_and_safety_outlined,
            ),
            const SizedBox(height: 8),
            FilledButton.icon(
              onPressed: () {
                Navigator.of(context).push(
                  MaterialPageRoute<void>(
                    builder: (_) => const FeedbackScreen(),
                  ),
                );
              },
              icon: const Icon(Icons.feedback_outlined),
              label: const Text('Send feedback'),
            ),
            const SizedBox(height: 8),
            OutlinedButton.icon(
              onPressed: () {
                Navigator.of(context).push(
                  MaterialPageRoute<void>(
                    builder: (_) => const DiagnosticsScreen(),
                  ),
                );
              },
              icon: const Icon(Icons.settings_outlined),
              label: const Text('Diagnostics'),
            ),
          ],
        ),
      ),
    );
  }
}
