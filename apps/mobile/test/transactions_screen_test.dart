import 'package:flutter/material.dart';
import 'package:flutter_test/flutter_test.dart';
import 'package:montra_mobile/app/theme/app_theme.dart';
import 'package:montra_mobile/core/controllers/auth_controller.dart';
import 'package:montra_mobile/core/controllers/ledger_controller.dart';
import 'package:montra_mobile/core/networking/api_client.dart';
import 'package:montra_mobile/core/storage/token_storage.dart';
import 'package:montra_mobile/core/widgets/montra_scope.dart';
import 'package:montra_mobile/features/transactions/transactions_screen.dart';
import 'package:shared_preferences/shared_preferences.dart';

void main() {
  setUp(() {
    SharedPreferences.setMockInitialValues(<String, Object>{});
  });

  Widget buildTransactionsWidget() {
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
            home: Scaffold(
              body: TransactionsScreen(
                onRequestAddTransaction: () {},
              ),
            ),
          ),
        );
      },
    );
  }

  testWidgets(
      'TransactionsScreen displays filter chips, search field, and empty state',
      (WidgetTester tester) async {
    await tester.pumpWidget(buildTransactionsWidget());
    await tester.pumpAndSettle();

    expect(find.text('Transaksi'), findsOneWidget);
    expect(find.text('Semua'), findsOneWidget);
    expect(find.text('Pengeluaran'), findsNWidgets(2)); // summary card + chip
    expect(find.text('Pemasukan'), findsNWidgets(2)); // summary card + chip
    expect(find.text('Transfer'), findsOneWidget);

    expect(find.text('Belum ada transaksi tercatat'), findsOneWidget);
  });
}
