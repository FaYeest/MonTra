import 'package:flutter/widgets.dart';

import '../controllers/auth_controller.dart';
import '../controllers/ledger_controller.dart';

/// Provides [AuthController] and [LedgerController] down the widget tree.
class MontraScope extends InheritedWidget {
  const MontraScope({
    super.key,
    required this.authController,
    required this.ledgerController,
    required super.child,
  });

  final AuthController authController;
  final LedgerController ledgerController;

  static MontraScope of(BuildContext context) {
    final scope = context.dependOnInheritedWidgetOfExactType<MontraScope>();
    assert(scope != null, 'No MontraScope found in context');
    return scope!;
  }

  static AuthController authOf(BuildContext context) {
    return of(context).authController;
  }

  static LedgerController ledgerOf(BuildContext context) {
    return of(context).ledgerController;
  }

  @override
  bool updateShouldNotify(MontraScope oldWidget) {
    return authController != oldWidget.authController ||
        ledgerController != oldWidget.ledgerController;
  }
}
