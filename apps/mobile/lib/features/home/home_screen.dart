import 'package:flutter/material.dart';

import '../../app/theme/app_colors.dart';
import '../../app/theme/app_shapes.dart';
import '../../app/theme/app_spacing.dart';
import '../../app/theme/app_typography.dart';
import '../../core/formatting/idr_formatter.dart';

/// Home (`Beranda`) foundation screen defined in `docs/design.md` Section 6.1 & 12.
///
/// Displays the recorded balance surface with privacy toggle, cash-flow summary,
/// and an honest empty state until backend transactions are loaded in Phase 1/2.
class HomeScreen extends StatelessWidget {
  const HomeScreen({
    super.key,
    required this.concealBalances,
    required this.onToggleConcealBalances,
    required this.onRequestAddTransaction,
  });

  final bool concealBalances;
  final VoidCallback onToggleConcealBalances;
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
        Text(
          'Ringkasan keuangan',
          style: textTheme.headlineLarge,
        ),
        const SizedBox(height: MontraSpacing.xxs),
        Text(
          'Periode berjalan • Saldo dihitung dari catatan transaksi Anda',
          style: textTheme.bodySmall,
        ),
        const SizedBox(height: MontraSpacing.md),
        Container(
          padding: const EdgeInsets.all(MontraSpacing.lg),
          decoration: BoxDecoration(
            color: tokens.surface,
            borderRadius: MontraShapes.prominentBorderRadius,
            border: Border.all(
              color: tokens.outline,
              width: MontraShapes.outlineWidth,
            ),
          ),
          child: Column(
            crossAxisAlignment: CrossAxisAlignment.start,
            children: <Widget>[
              Row(
                mainAxisAlignment: MainAxisAlignment.spaceBetween,
                children: <Widget>[
                  Text(
                    'Total saldo tercatat',
                    style: textTheme.bodyMedium?.copyWith(
                      color: tokens.textSecondary,
                    ),
                  ),
                  Semantics(
                    button: true,
                    label: concealBalances
                        ? 'Tampilkan nominal saldo'
                        : 'Sembunyikan nominal saldo',
                    child: IconButton(
                      constraints: const BoxConstraints(
                        minWidth: MontraSpacing.minTouchTarget,
                        minHeight: MontraSpacing.minTouchTarget,
                      ),
                      onPressed: onToggleConcealBalances,
                      icon: Icon(
                        concealBalances
                            ? Icons.visibility_off_outlined
                            : Icons.visibility_outlined,
                        color: tokens.textSecondary,
                      ),
                    ),
                  ),
                ],
              ),
              Text(
                IdrFormatter.format(0, concealed: concealBalances),
                style: MontraTypography.primaryFinancialValue(
                  tokens.textPrimary,
                ),
              ),
              const SizedBox(height: MontraSpacing.md),
              Row(
                children: <Widget>[
                  Expanded(
                    child: _CashFlowMetric(
                      label: 'Pemasukan',
                      value: IdrFormatter.formatSigned(
                        0,
                        isExpense: false,
                        concealed: concealBalances,
                      ),
                    ),
                  ),
                  const SizedBox(width: MontraSpacing.sm),
                  Expanded(
                    child: _CashFlowMetric(
                      label: 'Pengeluaran',
                      value: IdrFormatter.formatSigned(
                        0,
                        isExpense: true,
                        concealed: concealBalances,
                      ),
                    ),
                  ),
                ],
              ),
            ],
          ),
        ),
        const SizedBox(height: MontraSpacing.xl),
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
                  'Catat pemasukan atau pengeluaran pertama Anda untuk melihat '
                  'arus kas dan riwayat terbaru.',
                  style: textTheme.bodyMedium?.copyWith(
                    color: tokens.textSecondary,
                  ),
                ),
                const SizedBox(height: MontraSpacing.md),
                SizedBox(
                  height: MontraSpacing.minTouchTarget,
                  child: FilledButton.icon(
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

class _CashFlowMetric extends StatelessWidget {
  const _CashFlowMetric({
    required this.label,
    required this.value,
  });

  final String label;
  final String value;

  @override
  Widget build(BuildContext context) {
    final tokens = Theme.of(context).extension<MontraColorTokens>()!;

    return Container(
      padding: const EdgeInsets.all(MontraSpacing.sm),
      decoration: BoxDecoration(
        color: tokens.surfaceSubtle,
        borderRadius: MontraShapes.controlBorderRadius,
      ),
      child: Column(
        crossAxisAlignment: CrossAxisAlignment.start,
        children: <Widget>[
          Text(
            label,
            style: Theme.of(context).textTheme.bodySmall,
          ),
          const SizedBox(height: MontraSpacing.xxs),
          Text(
            value,
            style: MontraTypography.secondaryFinancialValue(
              tokens.textPrimary,
            ),
          ),
        ],
      ),
    );
  }
}
