/// Deterministic Indonesian Rupiah (`IDR`) formatter defined in `docs/design.md` Section 5.3, 5.4, and 13.
///
/// Rules:
/// - Keep currency formatting consistent across the app: `Rp2.840.000`.
/// - Income and expenses must include an explicit sign (`+Rp500.000`, `−Rp25.000`) so color is never the sole indicator.
/// - Support balance privacy concealment (`Rp••••••`).
abstract final class IdrFormatter {
  static const String maskedBalance = 'Rp••••••';
  static const String minusSign = '−'; // Unicode minus (U+2212) per design.md Section 5.3

  /// Formats an integer Rupiah amount (`int`) as `Rp2.840.000`.
  /// If [concealed] is true, returns [maskedBalance].
  static String format(int units, {bool concealed = false}) {
    if (concealed) {
      return maskedBalance;
    }
    final isNegative = units < 0;
    final grouped = _groupThousands(units.abs());
    if (isNegative) {
      return '${minusSign}Rp$grouped';
    }
    return 'Rp$grouped';
  }

  /// Formats a signed cash-flow transaction amount with an explicit `+` or `−` prefix
  /// (`+Rp500.000` for income, `−Rp25.000` for expense).
  static String formatSigned(
    int positiveUnits, {
    required bool isExpense,
    bool concealed = false,
  }) {
    if (concealed) {
      return maskedBalance;
    }
    final grouped = _groupThousands(positiveUnits.abs());
    final prefix = isExpense ? minusSign : '+';
    return '${prefix}Rp$grouped';
  }

  static String _groupThousands(int nonNegativeValue) {
    final digits = nonNegativeValue.toString();
    if (digits.length <= 3) {
      return digits;
    }
    final buffer = StringBuffer();
    final firstGroupLen = digits.length % 3;
    if (firstGroupLen > 0) {
      buffer.write(digits.substring(0, firstGroupLen));
    }
    for (var i = firstGroupLen; i < digits.length; i += 3) {
      if (buffer.isNotEmpty) {
        buffer.write('.');
      }
      buffer.write(digits.substring(i, i + 3));
    }
    return buffer.toString();
  }
}
