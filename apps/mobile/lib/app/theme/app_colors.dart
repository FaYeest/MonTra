import 'package:flutter/material.dart';

/// Centralized color tokens for MonTra defined in `docs/design.md` Section 5.1 & 5.2.
///
/// Do not define ad-hoc color hex literals inside individual widgets.
/// Never use color as the sole indicator of income, expense, warning, or success.
@immutable
class MontraColorTokens extends ThemeExtension<MontraColorTokens> {
  const MontraColorTokens({
    required this.canvas,
    required this.surface,
    required this.surfaceSubtle,
    required this.textPrimary,
    required this.textSecondary,
    required this.outline,
    required this.primary,
    required this.primaryAlt,
    required this.primaryContainer,
    required this.accent,
    required this.expense,
    required this.expenseContainer,
    required this.warning,
  });

  final Color canvas;
  final Color surface;
  final Color surfaceSubtle;
  final Color textPrimary;
  final Color textSecondary;
  final Color outline;
  final Color primary;
  final Color primaryAlt;
  final Color primaryContainer;
  final Color accent;
  final Color expense;
  final Color expenseContainer;
  final Color warning;

  /// Light palette (`docs/design.md` Section 5.1).
  static const MontraColorTokens light = MontraColorTokens(
    canvas: Color(0xFFF4F3ED),
    surface: Color(0xFFFFFEFA),
    surfaceSubtle: Color(0xFFEEEDE5),
    textPrimary: Color(0xFF202B25),
    textSecondary: Color(0xFF747B73),
    outline: Color(0xFFDEDFD5),
    primary: Color(0xFF244D3D),
    primaryAlt: Color(0xFF315F49),
    primaryContainer: Color(0xFFDCE9D7),
    accent: Color(0xFFD7F17E),
    expense: Color(0xFFB95745),
    expenseContainer: Color(0xFFF6E4DF),
    warning: Color(0xFF956D2D),
  );

  /// Dark palette (`docs/design.md` Section 5.2).
  static const MontraColorTokens dark = MontraColorTokens(
    canvas: Color(0xFF151B18),
    surface: Color(0xFF1D2621),
    surfaceSubtle: Color(0xFF26312A),
    textPrimary: Color(0xFFF2F1E9),
    textSecondary: Color(0xFFA3AEA4),
    outline: Color(0xFF354139),
    primary: Color(0xFFD7F17E),
    primaryAlt: Color(0xFFB8D462),
    primaryContainer: Color(0xFF32483A),
    accent: Color(0xFFD7F17E),
    expense: Color(0xFFE58A76),
    expenseContainer: Color(0xFF3D2722),
    warning: Color(0xFFE2B875),
  );

  @override
  MontraColorTokens copyWith({
    Color? canvas,
    Color? surface,
    Color? surfaceSubtle,
    Color? textPrimary,
    Color? textSecondary,
    Color? outline,
    Color? primary,
    Color? primaryAlt,
    Color? primaryContainer,
    Color? accent,
    Color? expense,
    Color? expenseContainer,
    Color? warning,
  }) {
    return MontraColorTokens(
      canvas: canvas ?? this.canvas,
      surface: surface ?? this.surface,
      surfaceSubtle: surfaceSubtle ?? this.surfaceSubtle,
      textPrimary: textPrimary ?? this.textPrimary,
      textSecondary: textSecondary ?? this.textSecondary,
      outline: outline ?? this.outline,
      primary: primary ?? this.primary,
      primaryAlt: primaryAlt ?? this.primaryAlt,
      primaryContainer: primaryContainer ?? this.primaryContainer,
      accent: accent ?? this.accent,
      expense: expense ?? this.expense,
      expenseContainer: expenseContainer ?? this.expenseContainer,
      warning: warning ?? this.warning,
    );
  }

  @override
  MontraColorTokens lerp(ThemeExtension<MontraColorTokens>? other, double t) {
    if (other is! MontraColorTokens) {
      return this;
    }
    return MontraColorTokens(
      canvas: Color.lerp(canvas, other.canvas, t)!,
      surface: Color.lerp(surface, other.surface, t)!,
      surfaceSubtle: Color.lerp(surfaceSubtle, other.surfaceSubtle, t)!,
      textPrimary: Color.lerp(textPrimary, other.textPrimary, t)!,
      textSecondary: Color.lerp(textSecondary, other.textSecondary, t)!,
      outline: Color.lerp(outline, other.outline, t)!,
      primary: Color.lerp(primary, other.primary, t)!,
      primaryAlt: Color.lerp(primaryAlt, other.primaryAlt, t)!,
      primaryContainer: Color.lerp(primaryContainer, other.primaryContainer, t)!,
      accent: Color.lerp(accent, other.accent, t)!,
      expense: Color.lerp(expense, other.expense, t)!,
      expenseContainer: Color.lerp(expenseContainer, other.expenseContainer, t)!,
      warning: Color.lerp(warning, other.warning, t)!,
    );
  }
}
