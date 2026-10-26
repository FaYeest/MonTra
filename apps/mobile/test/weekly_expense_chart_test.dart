import 'package:flutter/material.dart';
import 'package:flutter_test/flutter_test.dart';
import 'package:montra_mobile/app/theme/app_theme.dart';
import 'package:montra_mobile/core/formatting/idr_formatter.dart';
import 'package:montra_mobile/core/models/transaction.dart';
import 'package:montra_mobile/features/home/models/daily_expense.dart';
import 'package:montra_mobile/features/home/widgets/weekly_expense_chart_card.dart';

void main() {
  group('DailyExpense computation', () {
    test('computes exactly 7 days ending today', () {
      final now = DateTime(2026, 10, 26, 12, 0, 0);
      final txs = <TransactionModel>[
        TransactionModel(
          id: 'tx-1',
          accountId: 'acc-1',
          categoryId: 'cat-1',
          type: 'expense',
          amount: 25000,
          currency: 'IDR',
          description: 'Makan Siang',
          transactionDate: DateTime(2026, 10, 26, 11, 0, 0),
        ),
        TransactionModel(
          id: 'tx-2',
          accountId: 'acc-1',
          categoryId: 'cat-1',
          type: 'expense',
          amount: 15000,
          currency: 'IDR',
          description: 'Kopi',
          transactionDate: DateTime(2026, 10, 25, 9, 0, 0),
        ),
        TransactionModel(
          id: 'tx-3',
          accountId: 'acc-1',
          categoryId: 'cat-inc',
          type: 'income',
          amount: 1000000,
          currency: 'IDR',
          description: 'Gaji',
          transactionDate: DateTime(2026, 10, 26, 8, 0, 0),
        ),
      ];

      final days = DailyExpense.computeLast7Days(txs, anchorDate: now);
      expect(days.length, 7);
      expect(days.last.isToday, isTrue);
      expect(days.last.amount, 25000); // Only expenses counted
      expect(days[5].amount, 15000);
      expect(days[0].amount, 0);
    });

    test('IdrFormatter formatCompact formats correctly', () {
      expect(IdrFormatter.formatCompact(500), '500');
      expect(IdrFormatter.formatCompact(15000), '15rb');
      expect(IdrFormatter.formatCompact(15500), '15.5rb');
      expect(IdrFormatter.formatCompact(1000000), '1jt');
      expect(IdrFormatter.formatCompact(1500000), '1.5jt');
      expect(IdrFormatter.formatCompact(15000, concealed: true), '••••');
    });
  });

  group('WeeklyExpenseChartCard Widget', () {
    testWidgets('renders title, bars, and triggers onTap callback',
        (tester) async {
      bool tapped = false;
      final txs = <TransactionModel>[
        TransactionModel(
          id: 'tx-1',
          accountId: 'acc-1',
          categoryId: 'cat-1',
          type: 'expense',
          amount: 35000,
          currency: 'IDR',
          description: 'Bensin',
          transactionDate: DateTime.now(),
        ),
      ];

      await tester.pumpWidget(
        MaterialApp(
          theme: MontraTheme.light(),
          home: Scaffold(
            body: WeeklyExpenseChartCard(
              transactions: txs,
              concealBalances: false,
              onTap: () {
                tapped = true;
              },
            ),
          ),
        ),
      );

      expect(find.text('Tren Pengeluaran 7 Hari'), findsOneWidget);
      expect(find.text('Rincian'), findsOneWidget);
      expect(find.byType(WeeklyExpenseChartCard), findsOneWidget);

      await tester.tap(find.byType(WeeklyExpenseChartCard));
      await tester.pumpAndSettle();

      expect(tapped, isTrue);
    });
  });
}
