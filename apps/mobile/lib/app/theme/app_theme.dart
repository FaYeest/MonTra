import 'package:flutter/material.dart';

import 'app_colors.dart';
import 'app_shapes.dart';
import 'app_typography.dart';

/// Builds the light and dark `ThemeData` instances for MonTra (`docs/design.md` Section 5).
abstract final class MontraTheme {
  static ThemeData light() {
    return _build(
      brightness: Brightness.light,
      tokens: MontraColorTokens.light,
    );
  }

  static ThemeData dark() {
    return _build(
      brightness: Brightness.dark,
      tokens: MontraColorTokens.dark,
    );
  }

  static ThemeData _build({
    required Brightness brightness,
    required MontraColorTokens tokens,
  }) {
    final textTheme = MontraTypography.buildTextTheme(
      tokens.textPrimary,
      tokens.textSecondary,
    );

    final colorScheme = ColorScheme(
      brightness: brightness,
      primary: tokens.primary,
      onPrimary: brightness == Brightness.light
          ? tokens.surface
          : const Color(0xFF151B18),
      primaryContainer: tokens.primaryContainer,
      onPrimaryContainer: tokens.textPrimary,
      secondary: tokens.accent,
      onSecondary: const Color(0xFF202B25),
      error: tokens.expense,
      onError: tokens.surface,
      errorContainer: tokens.expenseContainer,
      onErrorContainer: tokens.textPrimary,
      surface: tokens.surface,
      onSurface: tokens.textPrimary,
      outline: tokens.outline,
    );

    return ThemeData(
      useMaterial3: true,
      brightness: brightness,
      scaffoldBackgroundColor: tokens.canvas,
      colorScheme: colorScheme,
      textTheme: textTheme,
      extensions: <ThemeExtension<dynamic>>[tokens],
      appBarTheme: AppBarTheme(
        backgroundColor: tokens.canvas,
        foregroundColor: tokens.textPrimary,
        elevation: 0,
        scrolledUnderElevation: 0,
        centerTitle: false,
        titleTextStyle: textTheme.titleLarge,
      ),
      cardTheme: CardThemeData(
        color: tokens.surface,
        elevation: 0,
        margin: EdgeInsets.zero,
        shape: RoundedRectangleBorder(
          borderRadius: MontraShapes.surfaceBorderRadius,
          side: BorderSide(
            color: tokens.outline,
            width: MontraShapes.outlineWidth,
          ),
        ),
      ),
      dividerTheme: DividerThemeData(
        color: tokens.outline,
        thickness: MontraShapes.outlineWidth,
        space: 1.0,
      ),
      bottomSheetTheme: BottomSheetThemeData(
        backgroundColor: tokens.surface,
        elevation: 2,
        shape: const RoundedRectangleBorder(
          borderRadius: MontraShapes.bottomSheetBorderRadius,
        ),
      ),
    );
  }
}
