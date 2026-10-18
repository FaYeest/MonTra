import 'package:flutter/material.dart';

import '../../app/theme/app_colors.dart';
import '../../app/theme/app_spacing.dart';

/// Insights (`Insight`) foundation screen defined in `docs/design.md` Section 6.5 & 12.
///
/// Explicitly avoids inventing sample insights or unsupported AI forecasts when history is empty.
class InsightsScreen extends StatelessWidget {
  const InsightsScreen({super.key});

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
        Text('Insight', style: textTheme.headlineLarge),
        const SizedBox(height: MontraSpacing.xxs),
        Text(
          'Pola pengeluaran, perbandingan periode, dan perkiraan berbasis data.',
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
                  'Insight akan muncul setelah tersedia cukup riwayat transaksi',
                  style: textTheme.titleMedium,
                ),
                const SizedBox(height: MontraSpacing.xs),
                Text(
                  'MonTra hanya menampilkan temuan statistik dan perkiraan akhir bulan '
                  'ketika data transaksi Anda mencukupi dan dapat diverifikasi.',
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
