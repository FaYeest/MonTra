import 'package:flutter/material.dart';

import '../../app/theme/app_colors.dart';
import '../../app/theme/app_shapes.dart';
import '../../app/theme/app_spacing.dart';
import '../../app/theme/app_typography.dart';
import '../../core/formatting/idr_formatter.dart';
import '../../core/widgets/montra_scope.dart';
import 'spending_trend_detail_screen.dart';
import 'widgets/weekly_expense_chart_card.dart';

/// Home (`Beranda`) screen defined in `docs/design.md` Section 6.1 & 12.
///
/// Displays the recorded balance surface with privacy toggle, cash-flow summary,
/// recent transactions, and an honest empty state when no transactions exist.
class HomeScreen extends StatelessWidget {
  const HomeScreen({
    super.key,
    required this.concealBalances,
    required this.onToggleConcealBalances,
    required this.onRequestAddTransaction,
    this.onViewAllTransactions,
  });

  final bool concealBalances;
  final VoidCallback onToggleConcealBalances;
  final VoidCallback onRequestAddTransaction;
  final VoidCallback? onViewAllTransactions;

  @override
  Widget build(BuildContext context) {
    final tokens = Theme.of(context).extension<MontraColorTokens>()!;
    final textTheme = Theme.of(context).textTheme;
    final ledger = MontraScope.ledgerOf(context);

    final transactions = ledger.transactions;
    final transfers = ledger.transfers;
    final hasRecords = transactions.isNotEmpty || transfers.isNotEmpty;

    return RefreshIndicator(
      onRefresh: () => ledger.refreshAll(),
      child: ListView(
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

          // Total Balance Card
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
                  IdrFormatter.format(
                    ledger.totalBalance,
                    concealed: concealBalances,
                  ),
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
                          ledger.totalIncome,
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
                          -ledger.totalExpense,
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
          const SizedBox(height: MontraSpacing.lg),

          // 7-Day Expense Trend Chart Card (Interactive)
          WeeklyExpenseChartCard(
            transactions: transactions,
            concealBalances: concealBalances,
            onTap: () {
              Navigator.of(context).push(
                MaterialPageRoute<void>(
                  builder: (_) => SpendingTrendDetailScreen(
                    concealBalances: concealBalances,
                  ),
                ),
              );
            },
          ),
          const SizedBox(height: MontraSpacing.lg),

          // Recent Transactions Section or Empty State
          if (!hasRecords) ...[
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
          ] else ...[
            Row(
              mainAxisAlignment: MainAxisAlignment.spaceBetween,
              children: <Widget>[
                Text(
                  'Transaksi Terbaru',
                  style: textTheme.titleMedium?.copyWith(
                    fontWeight: FontWeight.bold,
                  ),
                ),
                if (onViewAllTransactions != null)
                  TextButton(
                    onPressed: onViewAllTransactions,
                    child: const Text('Lihat semua'),
                  ),
              ],
            ),
            const SizedBox(height: MontraSpacing.xs),
            Card(
              child: ListView.separated(
                shrinkWrap: true,
                physics: const NeverScrollableScrollPhysics(),
                itemCount: transactions.length > 5 ? 5 : transactions.length,
                separatorBuilder: (_, __) => Divider(
                  height: 1,
                  color: tokens.outline.withValues(alpha: 0.5),
                ),
                itemBuilder: (context, idx) {
                  final tx = transactions[idx];
                  return ListTile(
                    leading: Container(
                      padding: const EdgeInsets.all(MontraSpacing.xs),
                      decoration: BoxDecoration(
                        color: tx.isIncome
                            ? tokens.income.withValues(alpha: 0.12)
                            : tokens.expense.withValues(alpha: 0.12),
                        borderRadius:
                            BorderRadius.circular(MontraShapes.radiusSm),
                      ),
                      child: Icon(
                        tx.isIncome
                            ? Icons.arrow_downward_rounded
                            : Icons.arrow_upward_rounded,
                        color: tx.isIncome ? tokens.income : tokens.expense,
                        size: 20,
                      ),
                    ),
                    title: Text(
                      tx.description,
                      style: textTheme.bodyMedium?.copyWith(
                        fontWeight: FontWeight.w600,
                      ),
                    ),
                    subtitle: Text(
                      tx.categoryName ??
                          (tx.isIncome ? 'Pemasukan' : 'Pengeluaran'),
                      style: textTheme.bodySmall?.copyWith(
                        color: tokens.textSecondary,
                      ),
                    ),
                    trailing: Text(
                      IdrFormatter.formatSigned(
                        tx.amount,
                        isExpense: !tx.isIncome,
                        concealed: concealBalances,
                      ),
                      style: textTheme.titleMedium?.copyWith(
                        fontWeight: FontWeight.bold,
                        color: tx.isIncome ? tokens.income : tokens.expense,
                      ),
                    ),
                  );
                },
              ),
            ),
          ],
        ],
      ),
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
