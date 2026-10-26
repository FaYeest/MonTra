import 'package:flutter/material.dart';
import 'package:flutter_test/flutter_test.dart';
import 'package:montra_mobile/app/theme/app_theme.dart';
import 'package:montra_mobile/core/controllers/auth_controller.dart';
import 'package:montra_mobile/core/controllers/ledger_controller.dart';
import 'package:montra_mobile/core/networking/api_client.dart';
import 'package:montra_mobile/core/storage/token_storage.dart';
import 'package:montra_mobile/core/widgets/montra_app_shell.dart';
import 'package:montra_mobile/core/widgets/montra_scope.dart';
import 'package:shared_preferences/shared_preferences.dart';

void main() {
  setUp(() {
    SharedPreferences.setMockInitialValues(<String, Object>{
      'montra_auth_token': 'mock_jwt_token',
      'montra_auth_user':
          '{"id":"usr-123","email":"test@montra.id","full_name":"Test User"}',
    });
  });

  Widget buildTestWidget({
    required AuthController authController,
    required LedgerController ledgerController,
  }) {
    return MontraScope(
      authController: authController,
      ledgerController: ledgerController,
      child: MaterialApp(
        theme: MontraTheme.light(),
        home: MontraAppShell(
          themeMode: ThemeMode.light,
          onToggleTheme: () {},
        ),
      ),
    );
  }

  testWidgets(
    'MontraAppShell renders 5 primary navigation items and toggles balance visibility',
    (WidgetTester tester) async {
      final storage = await TokenStorage.create();
      final client = ApiClient(tokenStorage: storage);
      final auth = AuthController(apiClient: client, tokenStorage: storage);
      final ledger = LedgerController(apiClient: client);

      await tester.pumpWidget(
        buildTestWidget(
          authController: auth,
          ledgerController: ledger,
        ),
      );

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
