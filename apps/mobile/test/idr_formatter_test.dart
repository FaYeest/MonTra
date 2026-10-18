import 'package:flutter_test/flutter_test.dart';
import 'package:montra_mobile/core/formatting/idr_formatter.dart';

void main() {
  group('IdrFormatter', () {
    test('formats standard Rupiah integers with dot thousand separators', () {
      expect(IdrFormatter.format(0), 'Rp0');
      expect(IdrFormatter.format(500), 'Rp500');
      expect(IdrFormatter.format(25000), 'Rp25.000');
      expect(IdrFormatter.format(2840000), 'Rp2.840.000');
    });

    test('formats signed income and expense values per design.md Section 5.3', () {
      expect(IdrFormatter.formatSigned(500000, isExpense: false), '+Rp500.000');
      expect(IdrFormatter.formatSigned(25000, isExpense: true), '−Rp25.000');
    });

    test('masks balances when concealed is true', () {
      expect(IdrFormatter.format(2840000, concealed: true), 'Rp••••••');
      expect(
        IdrFormatter.formatSigned(25000, isExpense: true, concealed: true),
        'Rp••••••',
      );
    });
  });
}
