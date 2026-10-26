import 'dart:math' as math;
import 'package:flutter/material.dart';

import '../../../app/theme/app_colors.dart';
import '../../../app/theme/app_shapes.dart';
import '../../../app/theme/app_spacing.dart';
import '../../../core/formatting/idr_formatter.dart';
import '../../../core/models/transaction.dart';
import '../models/daily_expense.dart';

/// Card widget displaying the 7-day expense bar chart on the Home screen.
///
/// Follows `docs/design.md` Section 6.1 (Spending trend for selected period)
/// and Nielsen usability heuristics with clear titles, units, and tap-for-detail interaction.
class WeeklyExpenseChartCard extends StatelessWidget {
  const WeeklyExpenseChartCard({
    super.key,
    required this.transactions,
    required this.concealBalances,
    required this.onTap,
  });

  final List<TransactionModel> transactions;
  final bool concealBalances;
  final VoidCallback onTap;

  @override
  Widget build(BuildContext context) {
    final tokens = Theme.of(context).extension<MontraColorTokens>()!;
    final textTheme = Theme.of(context).textTheme;

    final days = DailyExpense.computeLast7Days(transactions);
    final totalExpense = days.fold<int>(0, (sum, d) => sum + d.amount);
    final maxExpense = days.fold<int>(0, (prev, d) => math.max(prev, d.amount));
    final averageExpense = (totalExpense / 7).round();

    return Material(
      color: tokens.surface,
      shape: RoundedRectangleBorder(
        borderRadius: MontraShapes.prominentBorderRadius,
        side: BorderSide(
          color: tokens.outline,
          width: MontraShapes.outlineWidth,
        ),
      ),
      child: InkWell(
        onTap: onTap,
        borderRadius: MontraShapes.prominentBorderRadius,
        child: Padding(
          padding: const EdgeInsets.all(MontraSpacing.lg),
          child: Column(
            crossAxisAlignment: CrossAxisAlignment.start,
            children: <Widget>[
              // Header Row
              Row(
                mainAxisAlignment: MainAxisAlignment.spaceBetween,
                children: <Widget>[
                  Expanded(
                    child: Column(
                      crossAxisAlignment: CrossAxisAlignment.start,
                      children: <Widget>[
                        Row(
                          children: <Widget>[
                            Icon(
                              Icons.bar_chart_rounded,
                              size: 20,
                              color: tokens.primary,
                            ),
                            const SizedBox(width: MontraSpacing.xs),
                            Text(
                              'Tren Pengeluaran 7 Hari',
                              style: textTheme.titleMedium?.copyWith(
                                fontWeight: FontWeight.bold,
                              ),
                            ),
                          ],
                        ),
                        const SizedBox(height: MontraSpacing.xxs),
                        Text(
                          'Total 7 hari: ${IdrFormatter.format(totalExpense, concealed: concealBalances)}',
                          style: textTheme.bodySmall?.copyWith(
                            color: tokens.textSecondary,
                          ),
                        ),
                      ],
                    ),
                  ),
                  Row(
                    mainAxisSize: MainAxisSize.min,
                    children: <Widget>[
                      Text(
                        'Rincian',
                        style: textTheme.labelMedium?.copyWith(
                          color: tokens.primary,
                          fontWeight: FontWeight.w600,
                        ),
                      ),
                      const SizedBox(width: 2),
                      Icon(
                        Icons.chevron_right_rounded,
                        size: 20,
                        color: tokens.primary,
                      ),
                    ],
                  ),
                ],
              ),
              const SizedBox(height: MontraSpacing.sm),

              // Rata-rata mini badge
              Container(
                padding: const EdgeInsets.symmetric(
                  horizontal: MontraSpacing.sm,
                  vertical: MontraSpacing.xxs,
                ),
                decoration: BoxDecoration(
                  color: tokens.surfaceSubtle,
                  borderRadius: BorderRadius.circular(MontraShapes.radiusSm),
                ),
                child: Row(
                  mainAxisSize: MainAxisSize.min,
                  children: <Widget>[
                    Icon(
                      Icons.insights_rounded,
                      size: 14,
                      color: tokens.textSecondary,
                    ),
                    const SizedBox(width: MontraSpacing.xxs),
                    Text(
                      'Rata-rata: ${IdrFormatter.format(averageExpense, concealed: concealBalances)} / hari',
                      style: textTheme.labelSmall?.copyWith(
                        color: tokens.textSecondary,
                        fontWeight: FontWeight.w500,
                      ),
                    ),
                  ],
                ),
              ),
              const SizedBox(height: MontraSpacing.md),

              // 7-day Bar Chart
              SizedBox(
                height: 110,
                child: Row(
                  crossAxisAlignment: CrossAxisAlignment.end,
                  children: days.map((day) {
                    final ratio = maxExpense > 0
                        ? (day.amount / maxExpense).clamp(0.0, 1.0)
                        : 0.0;
                    return Expanded(
                      child: _DailyBarItem(
                        day: day,
                        fillRatio: ratio,
                        tokens: tokens,
                        textTheme: textTheme,
                      ),
                    );
                  }).toList(),
                ),
              ),
              const SizedBox(height: MontraSpacing.xs),

              // Footer Tap Hint
              Center(
                child: Text(
                  'Ketuk untuk melihat riwayat pengeluaran per hari',
                  style: textTheme.labelSmall?.copyWith(
                    color: tokens.textSecondary.withValues(alpha: 0.8),
                    fontSize: 11,
                  ),
                ),
              ),
            ],
          ),
        ),
      ),
    );
  }
}

class _DailyBarItem extends StatelessWidget {
  const _DailyBarItem({
    required this.day,
    required this.fillRatio,
    required this.tokens,
    required this.textTheme,
  });

  final DailyExpense day;
  final double fillRatio;
  final MontraColorTokens tokens;
  final TextTheme textTheme;

  static const double _barMaxHeight = 72.0;
  static const double _barWidth = 14.0;

  @override
  Widget build(BuildContext context) {
    final barFillHeight = math.max(
      day.amount > 0 ? 6.0 : 3.0,
      fillRatio * _barMaxHeight,
    );

    final barColor = day.amount > 0
        ? (day.isToday ? tokens.primary : tokens.primary.withValues(alpha: 0.6))
        : tokens.outline.withValues(alpha: 0.5);

    return Column(
      mainAxisAlignment: MainAxisAlignment.end,
      children: <Widget>[
        // Bar Column
        Container(
          width: _barWidth,
          height: _barMaxHeight,
          alignment: Alignment.bottomCenter,
          decoration: BoxDecoration(
            color: tokens.surfaceSubtle,
            borderRadius: BorderRadius.circular(_barWidth / 2),
          ),
          child: AnimatedContainer(
            duration: const Duration(milliseconds: 350),
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
            fontWeight: day.isToday ? FontWeight.bold : FontWeight.w500,
            color: day.isToday ? tokens.primary : tokens.textSecondary,
          ),
        ),

        // Small indicator dot for Today
        if (day.isToday)
          Container(
            margin: const EdgeInsets.only(top: 2),
            width: 4,
            height: 4,
            decoration: BoxDecoration(
              shape: BoxShape.circle,
              color: tokens.primary,
            ),
          )
        else
          const SizedBox(height: 6),
      ],
    );
  }
}
