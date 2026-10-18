import 'package:flutter/material.dart';

/// Border radius and outline rules defined in `docs/design.md` Section 5.6.
///
/// Prefer a 1 dp outline or surface contrast to heavy drop shadows.
abstract final class MontraShapes {
  static const double radiusControl = 12.0;
  static const double radiusSurface = 16.0;
  static const double radiusProminent = 20.0;
  static const double radiusBottomSheet = 24.0;
  static const double outlineWidth = 1.0;

  static const BorderRadius controlBorderRadius =
      BorderRadius.all(Radius.circular(radiusControl));
  static const BorderRadius surfaceBorderRadius =
      BorderRadius.all(Radius.circular(radiusSurface));
  static const BorderRadius prominentBorderRadius =
      BorderRadius.all(Radius.circular(radiusProminent));
  static const BorderRadius bottomSheetBorderRadius = BorderRadius.vertical(
    top: Radius.circular(radiusBottomSheet),
  );
}
