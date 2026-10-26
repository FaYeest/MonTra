import 'package:flutter/material.dart';
import 'package:flutter_test/flutter_test.dart';
import 'package:montra_mobile/app/theme/app_theme.dart';
import 'package:montra_mobile/core/controllers/auth_controller.dart';
import 'package:montra_mobile/core/controllers/ledger_controller.dart';
import 'package:montra_mobile/core/networking/api_client.dart';
import 'package:montra_mobile/core/storage/token_storage.dart';
import 'package:montra_mobile/core/widgets/montra_scope.dart';
import 'package:montra_mobile/features/transactions/add_transaction_sheet.dart';
import 'package:shared_preferences/shared_preferences.dart';

void main() {
  setUp(() {
    SharedPreferences.setMockInitialValues(<String, Object>{});
  });

  Widget buildSheetWidget() {
    return FutureBuilder<TokenStorage>(
      future: TokenStorage.create(),
      builder: (context, snapshot) {
        if (!snapshot.hasData) return const SizedBox();
        final storage = snapshot.data!;
        final client = ApiClient(tokenStorage: storage);
        return MontraScope(
          authController:
              AuthController(apiClient: client, tokenStorage: storage),
          ledgerController: LedgerController(apiClient: client),
          child: MaterialApp(
            theme: MontraTheme.light(),
            home: const Scaffold(
              body: AddTransactionSheet(),
            ),
          ),
        );
      },
    );
  }

  testWidgets(
      'AddTransactionSheet displays segment switchers and changes form fields on transfer',
      (WidgetTester tester) async {
    await tester.pumpWidget(buildSheetWidget());
    await tester.pumpAndSettle();

    expect(find.text('Catat Transaksi'), findsOneWidget);
    expect(find.text('Pengeluaran'), findsOneWidget);
    expect(find.text('Pemasukan'), findsOneWidget);
    expect(find.text('Transfer'), findsOneWidget);

    expect(find.text('Nominal Transaksi'), findsOneWidget);
    expect(find.text('Keterangan / Merchant'), findsOneWidget);

    // Switch to Transfer segment
    await tester.tap(find.text('Transfer'));
    await tester.pumpAndSettle();

    expect(find.text('Simpan Transfer'), findsOneWidget);
    expect(find.text('Biaya Admin / Transfer (Opsional)'), findsOneWidget);
  });
}
