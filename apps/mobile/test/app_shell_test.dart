import 'package:flutter/material.dart';
import 'package:flutter_test/flutter_test.dart';
import 'package:montra_mobile/app/app.dart';

void main() {
  testWidgets(
    'MontraAppShell renders 5 primary navigation items and toggles balance visibility',
    (WidgetTester tester) async {
      await tester.pumpWidget(const MontraApp());

      expect(find.text('Beranda'), findsOneWidget);
      expect(find.text('Transaksi'), findsOneWidget);
      expect(find.text('Tambah'), findsOneWidget);
      expect(find.text('Anggaran'), findsOneWidget);
      expect(find.text('Insight'), findsOneWidget);

      expect(find.text('Rp0'), findsOneWidget);

      await tester.tap(find.byIcon(Icons.visibility_outlined));
      await tester.pump();

      expect(find.text('Rp••••••'), findsWidgets);
    },
  );
}
