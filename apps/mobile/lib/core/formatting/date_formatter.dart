/// Safe, zero-dependency, deterministic Indonesian date formatter for MonTra.
/// Defined in `docs/design.md` Section 5.3 & Heuristic 2.
abstract final class MontraDateFormatter {
  static const List<String> _days = <String>[
    'Senin',
    'Selasa',
    'Rabu',
    'Kamis',
    'Jumat',
    'Sabtu',
    'Minggu',
  ];

  static const List<String> _months = <String>[
    'Januari',
    'Februari',
    'Maret',
    'April',
    'Mei',
    'Juni',
    'Juli',
    'Agustus',
    'September',
    'Oktober',
    'November',
    'Desember',
  ];

  /// Formats date as `EEEE, d MMMM yyyy` (e.g. `Senin, 26 Oktober 2026`).
  static String formatFull(DateTime date) {
    final dayName = _days[date.weekday - 1];
    final monthName = _months[date.month - 1];
    return '$dayName, ${date.day} $monthName ${date.year}';
  }

  /// Formats date as `d MMMM yyyy` (e.g. `26 Oktober 2026`).
  static String formatMedium(DateTime date) {
    final monthName = _months[date.month - 1];
    return '${date.day} $monthName ${date.year}';
  }

  /// Relative date label: 'Hari Ini', 'Kemarin', or `formatFull(date)`.
  static String formatRelative(DateTime date, {DateTime? relativeTo}) {
    final ref = relativeTo ?? DateTime.now();
    final today = DateTime(ref.year, ref.month, ref.day);
    final target = DateTime(date.year, date.month, date.day);
    final diffDays = today.difference(target).inDays;

    if (diffDays == 0) {
      return 'Hari Ini';
    } else if (diffDays == 1) {
      return 'Kemarin';
    }
    return formatFull(date);
  }
}
