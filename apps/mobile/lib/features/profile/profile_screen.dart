import 'package:flutter/material.dart';

import '../../app/theme/app_colors.dart';
import '../../app/theme/app_spacing.dart';

/// Profile and settings (`Profil`) screen defined in `docs/design.md` Section 6.6.
/// Opened via the avatar button in the top app bar.
class ProfileScreen extends StatelessWidget {
  const ProfileScreen({super.key});

  @override
  Widget build(BuildContext context) {
    final tokens = Theme.of(context).extension<MontraColorTokens>()!;
    final textTheme = Theme.of(context).textTheme;

    return Scaffold(
      appBar: AppBar(
        title: const Text('Profil & Pengaturan'),
      ),
      body: SafeArea(
        child: ListView(
          padding: const EdgeInsets.symmetric(
            horizontal: MontraSpacing.pageHorizontal,
            vertical: MontraSpacing.md,
          ),
          children: <Widget>[
            Card(
              child: Padding(
                padding: const EdgeInsets.all(MontraSpacing.lg),
                child: Column(
                  crossAxisAlignment: CrossAxisAlignment.start,
                  children: <Widget>[
                    Text(
                      'Preferensi & Privasi Data',
                      style: textTheme.titleMedium,
                    ),
                    const SizedBox(height: MontraSpacing.xs),
                    Text(
                      'Mata uang default: Rupiah Indonesia (IDR)\n'
                      'Zona waktu default: Asia/Jakarta\n'
                      'Versi aplikasi: 0.1.0 (Fase Fondasi)',
                      style: textTheme.bodyMedium?.copyWith(
                        color: tokens.textSecondary,
                      ),
                    ),
                  ],
                ),
              ),
            ),
          ],
        ),
      ),
    );
  }
}
