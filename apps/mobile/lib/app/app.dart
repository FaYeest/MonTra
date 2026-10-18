import 'package:flutter/material.dart';

import '../core/widgets/montra_app_shell.dart';
import 'theme/app_theme.dart';

/// Root MonTra Flutter application widget.
class MontraApp extends StatefulWidget {
  const MontraApp({super.key, this.initialThemeMode = ThemeMode.light});

  final ThemeMode initialThemeMode;

  @override
  State<MontraApp> createState() => _MontraAppState();
}

class _MontraAppState extends State<MontraApp> {
  late ThemeMode _themeMode = widget.initialThemeMode;

  void _toggleTheme() {
    setState(() {
      _themeMode =
          _themeMode == ThemeMode.dark ? ThemeMode.light : ThemeMode.dark;
    });
  }

  @override
  Widget build(BuildContext context) {
    return MaterialApp(
      title: 'MonTra',
      debugShowCheckedModeBanner: false,
      theme: MontraTheme.light(),
      darkTheme: MontraTheme.dark(),
      themeMode: _themeMode,
      home: MontraAppShell(
        themeMode: _themeMode,
        onToggleTheme: _toggleTheme,
      ),
    );
  }
}
