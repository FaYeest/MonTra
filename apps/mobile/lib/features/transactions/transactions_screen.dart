import 'package:flutter/material.dart';

import '../../app/theme/app_colors.dart';
import '../../app/theme/app_spacing.dart';

/// Transactions (`Transaksi`) foundation screen defined in `docs/design.md` Section 6.2 & 12.
class TransactionsScreen extends StatelessWidget {
  const TransactionsScreen({
    super.key,
    required this.onRequestAddTransaction,
  });

  final VoidCallback onRequestAddTransaction;

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
        Text('Transaksi', style: textTheme.headlineLarge),
        const SizedBox(height: MontraSpacing.xxs),
        Text(
          'Cari, periksa, dan kelola catatan transaksi Anda.',
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
                  'Belum ada transaksi untuk periode ini',
                  style: textTheme.titleMedium,
                ),
                const SizedBox(height: MontraSpacing.xs),
                Text(
                  'Transaksi yang telah disimpan ke server akan ditampilkan di sini '
                  'beserta pencarian dan filter kategori.',
                  style: textTheme.bodyMedium?.copyWith(
                    color: tokens.textSecondary,
                  ),
                ),
                const SizedBox(height: MontraSpacing.md),
                SizedBox(
                  height: MontraSpacing.minTouchTarget,
                  child: OutlinedButton.icon(
                    onPressed: onRequestAddTransaction,
                    icon: const Icon(Icons.add),
                    label: const Text('Catat transaksi'),
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
