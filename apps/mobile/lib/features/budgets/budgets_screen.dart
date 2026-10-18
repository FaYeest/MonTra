import 'package:flutter/material.dart';

import '../../app/theme/app_colors.dart';
import '../../app/theme/app_spacing.dart';

/// Budgets (`Anggaran`) foundation screen defined in `docs/design.md` Section 6.4 & 12.
class BudgetsScreen extends StatelessWidget {
  const BudgetsScreen({super.key});

  @override
  Widget build(BuildContext context) {
    final tokens = Theme.of(context).extension<MontraColorTokens>()!;
    final textTheme = Theme.of(context).textTheme;

    return ListView(
      padding: const EdgeInsets.symmetric(
        horizontal: MontraSpacing.pageHorizontal,
        vertical: MontraSpacing.md,
      ),
      children: <Widget>[
        Text('Anggaran', style: textTheme.headlineLarge),
        const SizedBox(height: MontraSpacing.xxs),
        Text(
          'Pantau batas pengeluaran bulanan dan sisa anggaran per kategori.',
          style: textTheme.bodySmall,
        ),
        const SizedBox(height: MontraSpacing.lg),
        Card(
          child: Padding(
            padding: const EdgeInsets.all(MontraSpacing.lg),
            child: Column(
              crossAxisAlignment: CrossAxisAlignment.start,
              children: <Widget>[
                Text(
                  'Belum ada anggaran aktif',
                  style: textTheme.titleMedium,
                ),
                const SizedBox(height: MontraSpacing.xs),
                Text(
                  'Pengaturan anggaran bulanan dan per kategori akan diaktifkan '
                  'pada Fase 2 setelah modul transaksi inti tersedia.',
                  style: textTheme.bodyMedium?.copyWith(
                    color: tokens.textSecondary,
                  ),
                ),
              ],
            ),
          ),
        ),
      ],
    );
  }
}
