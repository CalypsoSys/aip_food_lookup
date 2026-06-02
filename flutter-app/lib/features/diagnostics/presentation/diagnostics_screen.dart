import 'package:flutter/foundation.dart';
import 'package:flutter/material.dart';

import '../../../widgets/status_card.dart';
import 'diagnostics_controller.dart';

class DiagnosticsScreen extends StatefulWidget {
  const DiagnosticsScreen({super.key});

  @override
  State<DiagnosticsScreen> createState() => _DiagnosticsScreenState();
}

class _DiagnosticsScreenState extends State<DiagnosticsScreen> {
  late final DiagnosticsController _controller;

  @override
  void initState() {
    super.initState();
    _controller = DiagnosticsController();
  }

  @override
  void dispose() {
    _controller.dispose();
    super.dispose();
  }

  @override
  Widget build(BuildContext context) {
    return ValueListenableBuilder<DiagnosticsState>(
      valueListenable: _controller,
      builder: (context, state, _) {
        return Scaffold(
          appBar: AppBar(title: const Text('Diagnostics')),
          body: SafeArea(
            child: ListView(
              padding: const EdgeInsets.all(12),
              children: [
                StatusCard(
                  title: 'Backend URL',
                  subtitle: state.backendUrl,
                  icon: Icons.dns_outlined,
                ),
                StatusCard(
                  title: 'Platform',
                  subtitle: defaultTargetPlatform.name,
                  icon: Icons.phone_android,
                ),
                const StatusCard(
                  title: 'Android URL tip',
                  subtitle:
                      'Emulators use 10.0.2.2 for the Windows host. Physical phones need this PC\'s LAN IP address.',
                  icon: Icons.info_outline,
                ),
                const SizedBox(height: 8),
                FilledButton.icon(
                  onPressed:
                      state.isChecking ? null : _controller.testConnection,
                  icon: state.isChecking
                      ? const SizedBox.square(
                          dimension: 18,
                          child: CircularProgressIndicator(strokeWidth: 2),
                        )
                      : const Icon(Icons.wifi_find_outlined),
                  label: Text(
                    state.isChecking ? 'Testing...' : 'Test connection',
                  ),
                ),
                if (state.healthMessage != null) ...[
                  const SizedBox(height: 12),
                  StatusCard(
                    title: 'Backend reachable',
                    subtitle: state.healthMessage,
                    icon: Icons.check_circle_outline,
                  ),
                ],
                if (state.errorMessage != null) ...[
                  const SizedBox(height: 12),
                  StatusCard(
                    title: state.errorMessage!,
                    subtitle:
                        'Confirm your phone can reach the configured backend URL.',
                    icon: Icons.wifi_off_outlined,
                    actionLabel: 'Retry',
                    onAction: _controller.testConnection,
                  ),
                ],
              ],
            ),
          ),
        );
      },
    );
  }
}
