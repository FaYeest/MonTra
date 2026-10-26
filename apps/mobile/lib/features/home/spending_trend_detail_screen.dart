import 'dart:math' as math;
import 'package:flutter/material.dart';

import '../../app/theme/app_colors.dart';
import '../../app/theme/app_shapes.dart';
import '../../app/theme/app_spacing.dart';
import '../../app/theme/app_typography.dart';
import '../../core/formatting/idr_formatter.dart';
import '../../core/widgets/montra_scope.dart';
import '../transactions/add_transaction_sheet.dart';
import 'models/daily_expense.dart';

/// Full-page interactive detail screen for 7-day spending trends.
///
/// Opened when the user taps on [WeeklyExpenseChartCard] in Beranda.
/// Allows inspecting daily expenditures, peak spending, and day-by-day transaction breakdown.
class SpendingTrendDetailScreen extends StatefulWidget {
  const SpendingTrendDetailScreen({
    super.key,
    required this.concealBalances,
  });

  final bool concealBalances;

  @override
  State<SpendingTrendDetailScreen> createState() =>
      _SpendingTrendDetailScreenState();
}

class _SpendingTrendDetailScreenState extends State<SpendingTrendDetailScreen> {
  int _selectedDayIndex = 6; // Default to today (last index in 7 days)

  @override
  Widget build(BuildContext context) {
    final tokens = Theme.of(context).extension<MontraColorTokens>()!;
    final textTheme = Theme.of(context).textTheme;
    final ledger = MontraScope.ledgerOf(context);

    final days = DailyExpense.computeLast7Days(ledger.transactions);
    final totalExpense = days.fold<int>(0, (sum, d) => sum + d.amount);
    final maxExpense = days.fold<int>(0, (prev, d) => math.max(prev, d.amount));
    final averageExpense = (totalExpense / 7).round();

    // Identify peak day
    DailyExpense? peakDay;
    for (final d in days) {
      if (d.amount > 0 && (peakDay == null || d.amount > peakDay.amount)) {
        peakDay = d;
      }
    }

    final safeIndex = _selectedDayIndex.clamp(0, days.length - 1);
    final selectedDay = days[safeIndex];

    return Scaffold(
      appBar: AppBar(
        title: const Text('Rincian Tren Pengeluaran'),
        leading: IconButton(
          icon: const Icon(Icons.arrow_back_rounded),
          tooltip: 'Kembali ke Beranda',
          onPressed: () => Navigator.of(context).pop(),
        ),
      ),
      body: ListView(
        padding: const EdgeInsets.symmetric(
          horizontal: MontraSpacing.pageHorizontal,
          vertical: MontraSpacing.md,
        ),
        children: <Widget>[
          // Summary Banner Card
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
                Text(
                  'Total Pengeluaran 7 Hari Terakhir',
                  style: textTheme.bodySmall?.copyWith(
                    color: tokens.textSecondary,
                  ),
                ),
                const SizedBox(height: MontraSpacing.xxs),
                Text(
                  IdrFormatter.format(
                    totalExpense,
                    concealed: widget.concealBalances,
                  ),
                  style: MontraTypography.primaryFinancialValue(
                    tokens.textPrimary,
                  ),
                ),
                const SizedBox(height: MontraSpacing.md),
                Row(
                  children: <Widget>[
                    Expanded(
                      child: _SummaryPill(
                        label: 'Rata-rata Harian',
                        value: IdrFormatter.format(
                          averageExpense,
                          concealed: widget.concealBalances,
                        ),
                        icon: Icons.calendar_today_rounded,
                        tokens: tokens,
                        textTheme: textTheme,
                      ),
                    ),
                    const SizedBox(width: MontraSpacing.sm),
                    Expanded(
                      child: _SummaryPill(
                        label: 'Hari Tertinggi',
                        value: peakDay != null
                            ? '${peakDay.dayLabel} (${IdrFormatter.format(peakDay.amount, concealed: widget.concealBalances)})'
                            : 'Rp0',
                        icon: Icons.trending_up_rounded,
                        tokens: tokens,
                        textTheme: textTheme,
                      ),
                    ),
                  ],
                ),
              ],
            ),
          ),
          const SizedBox(height: MontraSpacing.lg),

          // Interactive Chart Card
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
                      'Pilih Hari untuk Rincian',
                      style: textTheme.titleMedium?.copyWith(
                        fontWeight: FontWeight.bold,
                      ),
                    ),
                    Container(
                      padding: const EdgeInsets.symmetric(
                        horizontal: MontraSpacing.sm,
                        vertical: MontraSpacing.xxs,
                      ),
                      decoration: BoxDecoration(
                        color: tokens.primaryContainer,
                        borderRadius:
                            BorderRadius.circular(MontraShapes.radiusSm),
                      ),
                      child: Text(
                        '${selectedDay.dayLabel}, ${selectedDay.formattedDate}',
                        style: textTheme.labelSmall?.copyWith(
                          fontWeight: FontWeight.bold,
                          color: tokens.primary,
                        ),
                      ),
                    ),
                  ],
                ),
                const SizedBox(height: MontraSpacing.xs),
                Text(
                  'Sentuh batang grafik untuk melihat daftar transaksi pada tanggal tersebut.',
                  style: textTheme.bodySmall?.copyWith(
                    color: tokens.textSecondary,
                  ),
                ),
                const SizedBox(height: MontraSpacing.lg),

                // Interactive Bar Chart
                SizedBox(
                  height: 150,
                  child: Row(
                    crossAxisAlignment: CrossAxisAlignment.end,
                    children: List.generate(days.length, (idx) {
                      final day = days[idx];
                      final isSelected = idx == safeIndex;
                      final ratio = maxExpense > 0
                          ? (day.amount / maxExpense).clamp(0.0, 1.0)
                          : 0.0;

                      return Expanded(
                        child: GestureDetector(
                          behavior: HitTestBehavior.opaque,
                          onTap: () {
                            setState(() {
                              _selectedDayIndex = idx;
                            });
                          },
                          child: _InteractiveBarColumn(
                            day: day,
                            fillRatio: ratio,
                            isSelected: isSelected,
                            concealBalances: widget.concealBalances,
                            tokens: tokens,
                            textTheme: textTheme,
                          ),
                        ),
                      );
                    }),
                  ),
                ),
              ],
            ),
          ),
          const SizedBox(height: MontraSpacing.lg),

          // Selected Day Transaction Breakdown Section
          Row(
            mainAxisAlignment: MainAxisAlignment.spaceBetween,
            children: <Widget>[
              Text(
                'Transaksi ${selectedDay.dayLabel} (${selectedDay.formattedDate})',
                style: textTheme.titleMedium?.copyWith(
                  fontWeight: FontWeight.bold,
                ),
              ),
              Text(
                IdrFormatter.format(
                  selectedDay.amount,
                  concealed: widget.concealBalances,
                ),
                style: textTheme.titleMedium?.copyWith(
                  fontWeight: FontWeight.bold,
                  color: selectedDay.amount > 0
                      ? tokens.expense
                      : tokens.textSecondary,
                ),
              ),
            ],
          ),
          const SizedBox(height: MontraSpacing.xs),

          if (selectedDay.transactions.isEmpty)
            Card(
              child: Padding(
                padding: const EdgeInsets.symmetric(
                  horizontal: MontraSpacing.md,
                  vertical: MontraSpacing.xl,
                ),
                child: Center(
                  child: Column(
                    children: <Widget>[
                      Icon(
                        Icons.receipt_long_outlined,
                        size: 40,
                        color: tokens.textSecondary.withValues(alpha: 0.4),
                      ),
                      const SizedBox(height: MontraSpacing.sm),
                      Text(
                        'Tidak ada pengeluaran pada ${selectedDay.dayLabel}, ${selectedDay.formattedDate}',
                        style: textTheme.bodyMedium?.copyWith(
                          color: tokens.textSecondary,
                          fontWeight: FontWeight.w500,
                        ),
                        textAlign: TextAlign.center,
                      ),
                    ],
                  ),
                ),
              ),
            )
          else
            Card(
              child: ListView.separated(
                shrinkWrap: true,
                physics: const NeverScrollableScrollPhysics(),
                itemCount: selectedDay.transactions.length,
                separatorBuilder: (_, __) => Divider(
                  height: 1,
                  color: tokens.outline.withValues(alpha: 0.5),
                ),
                itemBuilder: (context, idx) {
                  final tx = selectedDay.transactions[idx];
                  return ListTile(
                    leading: Container(
                      padding: const EdgeInsets.all(MontraSpacing.xs),
                      decoration: BoxDecoration(
                        color: tokens.expense.withValues(alpha: 0.12),
                        borderRadius:
                            BorderRadius.circular(MontraShapes.radiusSm),
                      ),
                      child: Icon(
                        Icons.arrow_upward_rounded,
                        color: tokens.expense,
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
                      tx.categoryName ?? 'Pengeluaran',
                      style: textTheme.bodySmall?.copyWith(
                        color: tokens.textSecondary,
                      ),
                    ),
                    trailing: Text(
                      IdrFormatter.formatSigned(
                        tx.amount,
                        isExpense: true,
                        concealed: widget.concealBalances,
                      ),
                      style: textTheme.titleMedium?.copyWith(
                        fontWeight: FontWeight.bold,
                        color: tokens.expense,
                      ),
                    ),
                  );
                },
              ),
            ),

          const SizedBox(height: MontraSpacing.xl),

          // Add Transaction Shortcut Button
          SizedBox(
            height: MontraSpacing.minTouchTarget,
            child: OutlinedButton.icon(
              onPressed: () => AddTransactionSheet.show(context),
              icon: const Icon(Icons.add),
              label: const Text('Catat Transaksi Baru'),
            ),
          ),
          const SizedBox(height: MontraSpacing.xl),
        ],
      ),
    );
  }
}

class _SummaryPill extends StatelessWidget {
  const _SummaryPill({
    required this.label,
    required this.value,
    required this.icon,
    required this.tokens,
    required this.textTheme,
  });

  final String label;
  final String value;
  final IconData icon;
  final MontraColorTokens tokens;
  final TextTheme textTheme;

  @override
  Widget build(BuildContext context) {
    return Container(
      padding: const EdgeInsets.all(MontraSpacing.sm),
      decoration: BoxDecoration(
        color: tokens.surfaceSubtle,
        borderRadius: BorderRadius.circular(MontraShapes.radiusSm),
      ),
      child: Column(
        crossAxisAlignment: CrossAxisAlignment.start,
        children: <Widget>[
          Row(
            children: <Widget>[
              Icon(icon, size: 14, color: tokens.textSecondary),
              const SizedBox(width: 4),
              Expanded(
                child: Text(
                  label,
                  style: textTheme.labelSmall?.copyWith(
                    color: tokens.textSecondary,
                  ),
                  overflow: TextOverflow.ellipsis,
                ),
              ),
            ],
          ),
          const SizedBox(height: MontraSpacing.xxs),
          Text(
            value,
            style: textTheme.bodyMedium?.copyWith(
              fontWeight: FontWeight.bold,
              color: tokens.textPrimary,
            ),
            overflow: TextOverflow.ellipsis,
          ),
        ],
      ),
    );
  }
}

class _InteractiveBarColumn extends StatelessWidget {
  const _InteractiveBarColumn({
    required this.day,
    required this.fillRatio,
    required this.isSelected,
    required this.concealBalances,
    required this.tokens,
    required this.textTheme,
  });

  final DailyExpense day;
  final double fillRatio;
  final bool isSelected;
  final bool concealBalances;
  final MontraColorTokens tokens;
  final TextTheme textTheme;

  static const double _barMaxHeight = 90.0;
  static const double _barWidth = 18.0;

  @override
  Widget build(BuildContext context) {
    final barFillHeight = math.max(
      day.amount > 0 ? 8.0 : 4.0,
      fillRatio * _barMaxHeight,
    );

    final barColor = isSelected
        ? tokens.primary
        : (day.amount > 0
            ? tokens.primary.withValues(alpha: 0.5)
            : tokens.outline.withValues(alpha: 0.5));

    return Column(
      mainAxisAlignment: MainAxisAlignment.end,
      children: <Widget>[
        // Selected Indicator badge on top
        if (isSelected)
          Container(
            margin: const EdgeInsets.only(bottom: 4),
            padding: const EdgeInsets.symmetric(horizontal: 4, vertical: 2),
            decoration: BoxDecoration(
              color: tokens.primary,
              borderRadius: BorderRadius.circular(4),
            ),
            child: Text(
              IdrFormatter.formatCompact(day.amount, concealed: concealBalances),
              style: textTheme.labelSmall?.copyWith(
                color: tokens.canvas,
                fontSize: 9,
                fontWeight: FontWeight.bold,
              ),
            ),
          )
        else
          const SizedBox(height: 18),

        // Bar Track & Fill
        Container(
          width: _barWidth,
          height: _barMaxHeight,
          alignment: Alignment.bottomCenter,
          decoration: BoxDecoration(
            color: isSelected
                ? tokens.primary.withValues(alpha: 0.12)
                : tokens.surfaceSubtle,
            borderRadius: BorderRadius.circular(_barWidth / 2),
            border: isSelected
                ? Border.all(color: tokens.primary, width: 1.5)
                : null,
          ),
          child: AnimatedContainer(
            duration: const Duration(milliseconds: 300),
            curve: Curves.easeOutCubic,
            width: _barWidth,
            height: barFillHeight,
            decoration: BoxDecoration(
              color: barColor,
              borderRadius: BorderRadius.circular(_barWidth / 2),
            ),
          ),
        ),
        const SizedBox(height: MontraSpacing.xs),

        // Day label
        Text(
          day.dayLabel,
          style: textTheme.labelSmall?.copyWith(
            fontWeight: isSelected ? FontWeight.bold : FontWeight.w500,
            color: isSelected ? tokens.primary : tokens.textSecondary,
          ),
        ),

        // Date number
        Text(
          '${day.date.day}',
          style: textTheme.labelSmall?.copyWith(
            fontSize: 10,
            fontWeight: isSelected ? FontWeight.bold : FontWeight.normal,
            color: isSelected
                ? tokens.primary
                : tokens.textSecondary.withValues(alpha: 0.7),
          ),
        ),
      ],
    );
  }
}
