import 'package:flutter/material.dart';

import '../core/controllers/auth_controller.dart';
import '../core/controllers/ledger_controller.dart';
import '../core/networking/api_client.dart';
import '../core/networking/api_config.dart';
import '../core/storage/token_storage.dart';
import '../core/widgets/montra_app_shell.dart';
import '../core/widgets/montra_scope.dart';
import '../features/auth/login_screen.dart';
import 'theme/app_colors.dart';
import 'theme/app_theme.dart';

/// Root MonTra Flutter application widget.
class MontraApp extends StatefulWidget {
  const MontraApp({
    super.key,
    this.initialThemeMode = ThemeMode.light,
    this.apiClient,
    this.tokenStorage,
    this.authController,
    this.ledgerController,
  });

  final ThemeMode initialThemeMode;
  final ApiClient? apiClient;
  final TokenStorage? tokenStorage;
  final AuthController? authController;
  final LedgerController? ledgerController;

  @override
  State<MontraApp> createState() => _MontraAppState();
}

class _MontraAppState extends State<MontraApp> {
  late ThemeMode _themeMode = widget.initialThemeMode;
  late final ApiClient _apiClient;
  late final AuthController _authController;
  late final LedgerController _ledgerController;
  bool _initialized = false;

  @override
  void initState() {
    super.initState();
    _bootstrap();
  }

  Future<void> _bootstrap() async {
    final storage = widget.tokenStorage ?? await TokenStorage.create();
    _apiClient = widget.apiClient ??
        ApiClient(
          config: const MontraApiConfig(),
          tokenStorage: storage,
        );

    _authController = widget.authController ??
        AuthController(
          apiClient: _apiClient,
          tokenStorage: storage,
        );

    _ledgerController = widget.ledgerController ??
        LedgerController(
          apiClient: _apiClient,
        );

    await _authController.checkAuth();
    if (_authController.isAuthenticated) {
      await _ledgerController.refreshAll();
    }

    if (mounted) {
      setState(() {
        _initialized = true;
      });
    }
  }

  void _toggleTheme() {
    setState(() {
      _themeMode =
          _themeMode == ThemeMode.dark ? ThemeMode.light : ThemeMode.dark;
    });
  }

  @override
  Widget build(BuildContext context) {
    if (!_initialized) {
      return MaterialApp(
        title: 'MonTra',
        debugShowCheckedModeBanner: false,
        theme: MontraTheme.light(),
        darkTheme: MontraTheme.dark(),
        themeMode: _themeMode,
        home: Scaffold(
          body: Center(
            child: Column(
              mainAxisAlignment: MainAxisAlignment.center,
              children: <Widget>[
                Icon(
                  Icons.account_balance_wallet_rounded,
                  size: 48,
                  color: MontraColorTokens.light.primary,
                ),
                const SizedBox(height: 16),
                const CircularProgressIndicator(),
              ],
            ),
          ),
        ),
      );
    }

    return MontraScope(
      authController: _authController,
      ledgerController: _ledgerController,
      child: ListenableBuilder(
        listenable: _authController,
        builder: (context, _) {
          return MaterialApp(
            title: 'MonTra',
            debugShowCheckedModeBanner: false,
            theme: MontraTheme.light(),
            darkTheme: MontraTheme.dark(),
            themeMode: _themeMode,
            home: _authController.isAuthenticated
                ? MontraAppShell(
                    themeMode: _themeMode,
                    onToggleTheme: _toggleTheme,
                  )
                : const LoginScreen(),
          );
        },
      ),
    );
  }
}
