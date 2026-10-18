/// 4 dp base spacing scale and touch target baselines defined in `docs/design.md` Section 5.5 & 9.1.
abstract final class MontraSpacing {
  static const double xxs = 4.0;
  static const double xs = 8.0;
  static const double sm = 12.0;
  static const double md = 16.0;
  static const double lg = 20.0;
  static const double xl = 24.0;
  static const double xxl = 32.0;
  static const double xxxl = 40.0;
  static const double jumbo = 48.0;

  /// Default horizontal padding for core pages (16–20 dp).
  static const double pageHorizontal = 16.0;

  /// Minimum interactive touch target dimension (48 × 48 dp per Android accessibility baseline).
  static const double minTouchTarget = 48.0;
}
