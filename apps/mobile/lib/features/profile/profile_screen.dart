import 'package:flutter/material.dart';

import '../../app/theme/app_colors.dart';
import '../../app/theme/app_shapes.dart';
import '../../app/theme/app_spacing.dart';
import '../../core/widgets/montra_scope.dart';

/// Profile and settings (`Profil`) screen defined in `docs/design.md` Section 6.6.
class ProfileScreen extends StatelessWidget {
  const ProfileScreen({super.key});

  Future<void> _confirmLogout(BuildContext context) async {
    final confirmed = await showDialog<bool>(
      context: context,
      builder: (dialogCtx) {
        return AlertDialog(
          title: const Text('Keluar dari Akun?'),
          content: const Text(
            'Sesi masuk Anda pada perangkat ini akan diakhiri. '
            'Catatan transaksi Anda tetap tersimpan aman di server.',
          ),
          actions: <Widget>[
            TextButton(
              onPressed: () => Navigator.of(dialogCtx).pop(false),
              child: const Text('Batal'),
            ),
            FilledButton(
              onPressed: () => Navigator.of(dialogCtx).pop(true),
              child: const Text('Keluar'),
            ),
          ],
        );
      },
    );

    if (confirmed == true && context.mounted) {
      final auth = MontraScope.authOf(context);
      await auth.logout();
      if (context.mounted && Navigator.of(context).canPop()) {
        Navigator.of(context).pop();
      }
    }
  }

  @override
  Widget build(BuildContext context) {
    final tokens = Theme.of(context).extension<MontraColorTokens>()!;
    final textTheme = Theme.of(context).textTheme;
    final auth = MontraScope.authOf(context);
    final user = auth.currentUser;

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
            // User Card
            Card(
              child: Padding(
                padding: const EdgeInsets.all(MontraSpacing.lg),
                child: Row(
                  children: <Widget>[
                    CircleAvatar(
                      radius: 28,
                      backgroundColor: tokens.primary.withValues(alpha: 0.15),
                      child: Text(
                        (user?.fullName.isNotEmpty ?? false)
                            ? user!.fullName[0].toUpperCase()
                            : 'M',
                        style: textTheme.titleLarge?.copyWith(
                          color: tokens.primary,
                          fontWeight: FontWeight.bold,
                        ),
                      ),
                    ),
                    const SizedBox(width: MontraSpacing.md),
                    Expanded(
                      child: Column(
                        crossAxisAlignment: CrossAxisAlignment.start,
                        children: <Widget>[
                          Text(
                            user?.fullName ?? 'Pengguna MonTra',
                            style: textTheme.titleMedium?.copyWith(
                              fontWeight: FontWeight.bold,
                            ),
                          ),
                          const SizedBox(height: MontraSpacing.xxs),
                          Text(
                            user?.email ?? 'Tidak terhubung',
                            style: textTheme.bodySmall?.copyWith(
                              color: tokens.textSecondary,
                            ),
                          ),
                        ],
                      ),
                    ),
                  ],
                ),
              ),
            ),
            const SizedBox(height: MontraSpacing.md),

            // Regional & Defaults Card
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
                      'Format: Rp2.840.000\n'
                      'Zona waktu default: Asia/Jakarta (WIB)\n'
                      'Versi aplikasi: 0.1.0 (Fase 1 Core Ledger)',
                      style: textTheme.bodyMedium?.copyWith(
                        color: tokens.textSecondary,
                        height: 1.5,
                      ),
                    ),
                  ],
                ),
              ),
            ),
            const SizedBox(height: MontraSpacing.lg),

            // Logout Button
            SizedBox(
              height: MontraSpacing.minTouchTarget,
              child: OutlinedButton.icon(
                style: OutlinedButton.styleFrom(
                  foregroundColor: tokens.expense,
                  side: BorderSide(
                    color: tokens.expense.withValues(alpha: 0.5),
                    width: MontraShapes.outlineWidth,
                  ),
                ),
                onPressed: () => _confirmLogout(context),
                icon: const Icon(Icons.logout_rounded),
                label: const Text('Keluar dari Akun'),
              ),
            ),
          ],
        ),
      ),
    );
  }
}
