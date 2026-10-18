import 'package:flutter/material.dart';

/// Typography scale defined in `docs/design.md` Section 5.4.
///
/// Uses scalable `sp` units (mapped to Flutter logical pixels with TextScaler)
/// and tabular numerals (`FontFeature.tabularFigures()`) for monetary values.
abstract final class MontraTypography {
  static const List<FontFeature> tabularNumerals = <FontFeature>[
    FontFeature.tabularFigures(),
  ];

  static TextTheme buildTextTheme(Color primaryColor, Color secondaryColor) {
    return TextTheme(
      headlineLarge: TextStyle(
        fontSize: 28.0,
        fontWeight: FontWeight.w700,
        height: 1.25,
        color: primaryColor,
      ),
      titleLarge: TextStyle(
        fontSize: 19.0,
        fontWeight: FontWeight.w700,
        height: 1.3,
        color: primaryColor,
      ),
      titleMedium: TextStyle(
        fontSize: 16.0,
        fontWeight: FontWeight.w600,
        height: 1.35,
        color: primaryColor,
      ),
      bodyLarge: TextStyle(
        fontSize: 15.0,
        fontWeight: FontWeight.w400,
        height: 1.45,
        color: primaryColor,
      ),
      bodyMedium: TextStyle(
        fontSize: 14.0,
        fontWeight: FontWeight.w400,
        height: 1.45,
        color: primaryColor,
      ),
      bodySmall: TextStyle(
        fontSize: 13.0,
        fontWeight: FontWeight.w400,
        height: 1.4,
        color: secondaryColor,
      ),
      labelLarge: TextStyle(
        fontSize: 14.0,
        fontWeight: FontWeight.w600,
        height: 1.2,
        color: primaryColor,
      ),
      labelMedium: TextStyle(
        fontSize: 12.0,
        fontWeight: FontWeight.w600,
        height: 1.2,
        color: secondaryColor,
      ),
    );
  }

  /// Large financial value style (28–34 sp, weight 700–800, tabular figures).
  static TextStyle primaryFinancialValue(Color color) {
    return TextStyle(
      fontSize: 30.0,
      fontWeight: FontWeight.w700,
      height: 1.2,
      color: color,
      fontFeatures: tabularNumerals,
    );
  }

  /// Secondary financial value style (18–22 sp, weight 650–750, tabular figures).
  static TextStyle secondaryFinancialValue(Color color) {
    return TextStyle(
      fontSize: 20.0,
      fontWeight: FontWeight.w700,
      height: 1.25,
      color: color,
      fontFeatures: tabularNumerals,
    );
  }
}
