import 'package:flutter/material.dart';
import 'package:flutter_test/flutter_test.dart';
import 'package:montra_mobile/app/theme/app_theme.dart';
import 'package:montra_mobile/core/controllers/auth_controller.dart';
import 'package:montra_mobile/core/controllers/ledger_controller.dart';
import 'package:montra_mobile/core/networking/api_client.dart';
import 'package:montra_mobile/core/storage/token_storage.dart';
import 'package:montra_mobile/core/widgets/montra_scope.dart';
import 'package:montra_mobile/features/auth/login_screen.dart';
import 'package:montra_mobile/features/auth/register_screen.dart';
import 'package:shared_preferences/shared_preferences.dart';

void main() {
  setUp(() {
    SharedPreferences.setMockInitialValues(<String, Object>{});
  });

  Widget buildScopeWidget(Widget child) {
    return FutureBuilder<TokenStorage>(
      future: TokenStorage.create(),
      builder: (context, snapshot) {
        if (!snapshot.hasData) return const SizedBox();
        final storage = snapshot.data!;
        final client = ApiClient(tokenStorage: storage);
        return MontraScope(
          authController: AuthController(apiClient: client, tokenStorage: storage),
          ledgerController: LedgerController(apiClient: client),
          child: MaterialApp(
            theme: MontraTheme.light(),
            home: child,
          ),
        );
      },
    );
  }

  testWidgets('LoginScreen displays required fields and validates empty inputs',
      (WidgetTester tester) async {
    await tester.pumpWidget(buildScopeWidget(const LoginScreen()));
    await tester.pumpAndSettle();

    expect(find.text('Masuk ke Akun Anda'), findsOneWidget);
    expect(find.widgetWithText(TextFormField, 'Email'), findsOneWidget);
    expect(find.widgetWithText(TextFormField, 'Kata Sandi'), findsOneWidget);
    expect(find.widgetWithText(ElevatedButton, 'Masuk'), findsOneWidget);

    // Tap submit without typing
    await tester.tap(find.widgetWithText(ElevatedButton, 'Masuk'));
    await tester.pumpAndSettle();

    expect(find.text('Email wajib diisi'), findsOneWidget);
  });

  testWidgets(
      'RegisterScreen displays required fields and validates password mismatch',
      (WidgetTester tester) async {
    await tester.pumpWidget(buildScopeWidget(const RegisterScreen()));
    await tester.pumpAndSettle();

    expect(find.text('Mulai Perjalanan Finansial Anda'), findsOneWidget);
    expect(find.widgetWithText(TextFormField, 'Nama Lengkap'), findsOneWidget);
    expect(find.widgetWithText(TextFormField, 'Email'), findsOneWidget);
    expect(find.widgetWithText(TextFormField, 'Kata Sandi'), findsOneWidget);
    expect(find.widgetWithText(TextFormField, 'Konfirmasi Kata Sandi'),
        findsOneWidget);

    // Fill form with mismatched passwords
    await tester.enterText(
        find.widgetWithText(TextFormField, 'Nama Lengkap'), 'Budi Santoso');
    await tester.enterText(
        find.widgetWithText(TextFormField, 'Email'), 'budi@example.com');
    await tester.enterText(
        find.widgetWithText(TextFormField, 'Kata Sandi'), 'Password123!');
    await tester.enterText(
        find.widgetWithText(TextFormField, 'Konfirmasi Kata Sandi'),
        'PasswordBeda!');

    await tester.tap(find.widgetWithText(ElevatedButton, 'Daftar Sekarang'));
    await tester.pumpAndSettle();

    expect(find.text('Konfirmasi kata sandi tidak cocok'), findsOneWidget);
  });
}
