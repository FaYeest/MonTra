import '../../../../core/models/transaction.dart';

/// Aggregated daily expense data for 7-day trend analysis.
class DailyExpense {
  const DailyExpense({
    required this.date,
    required this.dayLabel,
    required this.formattedDate,
    required this.amount,
    required this.isToday,
    required this.transactions,
  });

  final DateTime date;
  final String dayLabel;
  final String formattedDate;
  final int amount;
  final bool isToday;
  final List<TransactionModel> transactions;

  /// Helper to aggregate 7 days ending on [anchorDate] (defaults to DateTime.now()).
  static List<DailyExpense> computeLast7Days(
    List<TransactionModel> transactions, {
    DateTime? anchorDate,
  }) {
    final now = anchorDate ?? DateTime.now();
    final today = DateTime(now.year, now.month, now.day);

    const indonesianDays = <int, String>{
      DateTime.monday: 'Sen',
      DateTime.tuesday: 'Sel',
      DateTime.wednesday: 'Rab',
      DateTime.thursday: 'Kam',
      DateTime.friday: 'Jum',
      DateTime.saturday: 'Sab',
      DateTime.sunday: 'Min',
    };

    const monthNames = <int, String>{
      1: 'Jan',
      2: 'Feb',
      3: 'Mar',
      4: 'Apr',
      5: 'Mei',
      6: 'Jun',
      7: 'Jul',
      8: 'Agu',
      9: 'Sep',
      10: 'Okt',
      11: 'Nov',
      12: 'Des',
    };

    final result = <DailyExpense>[];

    // 7 days from (today - 6) up to today
    for (int i = 6; i >= 0; i--) {
      final dayDate = today.subtract(Duration(days: i));
      final dayLabel = indonesianDays[dayDate.weekday] ?? '';
      final formattedDate = '${dayDate.day} ${monthNames[dayDate.month] ?? ''}';

      final dayTxs = transactions.where((tx) {
        if (!tx.isExpense) return false;
        final d = tx.transactionDate.toLocal();
        return d.year == dayDate.year &&
            d.month == dayDate.month &&
            d.day == dayDate.day;
      }).toList();

      final total = dayTxs.fold<int>(0, (sum, tx) => sum + tx.amount);

      result.add(
        DailyExpense(
          date: dayDate,
          dayLabel: dayLabel,
          formattedDate: formattedDate,
          amount: total,
          isToday: i == 0,
          transactions: dayTxs,
        ),
      );
    }

    return result;
  }
}
