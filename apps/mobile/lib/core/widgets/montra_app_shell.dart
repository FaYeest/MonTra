import 'package:flutter/material.dart';

import '../../app/router.dart';
import '../../app/theme/app_colors.dart';
import '../../app/theme/app_shapes.dart';
import '../../app/theme/app_spacing.dart';
import '../../features/budgets/budgets_screen.dart';
import '../../features/home/home_screen.dart';
import '../../features/insights/insights_screen.dart';
import '../../features/profile/profile_screen.dart';
import '../../features/transactions/add_transaction_sheet.dart';
import '../../features/transactions/transactions_screen.dart';

/// Persistent 5-destination shell defined in `docs/design.md` Section 4.1 & 7.1:
/// - Beranda
/// - Transaksi
/// - Tambah (central emphasized action)
/// - Anggaran
/// - Insight
/// Plus top app bar avatar navigating to `Profil`.
class MontraAppShell extends StatefulWidget {
  const MontraAppShell({
    super.key,
    required this.themeMode,
    required this.onToggleTheme,
  });

  final ThemeMode themeMode;
  final VoidCallback onToggleTheme;

  @override
  State<MontraAppShell> createState() => _MontraAppShellState();
}

class _MontraAppShellState extends State<MontraAppShell> {
  MontraPrimaryTab _selectedTab = MontraPrimaryTab.beranda;
  bool _concealBalances = false;

  void _openAddTransactionSheet(BuildContext context) {
    AddTransactionSheet.show(context);
  }

  void _openProfileScreen(BuildContext context) {
    Navigator.of(context).push(
      MaterialPageRoute<void>(
        builder: (_) => const ProfileScreen(),
      ),
    );
  }

  @override
  Widget build(BuildContext context) {
    final tokens = Theme.of(context).extension<MontraColorTokens>()!;

    return Scaffold(
      appBar: AppBar(
        title: const Text('MonTra'),
        actions: <Widget>[
          Semantics(
            button: true,
            label: 'Ubah tema tampilan terang atau gelap',
            child: IconButton(
              constraints: const BoxConstraints(
                minWidth: MontraSpacing.minTouchTarget,
                minHeight: MontraSpacing.minTouchTarget,
              ),
              tooltip: 'Ubah tema',
              onPressed: widget.onToggleTheme,
              icon: Icon(
                widget.themeMode == ThemeMode.dark
                    ? Icons.light_mode_outlined
                    : Icons.dark_mode_outlined,
              ),
            ),
          ),
          Semantics(
            button: true,
            label: 'Buka Profil dan Pengaturan',
            child: IconButton(
              constraints: const BoxConstraints(
                minWidth: MontraSpacing.minTouchTarget,
                minHeight: MontraSpacing.minTouchTarget,
              ),
              tooltip: 'Profil',
              onPressed: () => _openProfileScreen(context),
              icon: const Icon(Icons.person_outline),
            ),
          ),
        ],
      ),
      body: SafeArea(
        child: IndexedStack(
          index: _selectedTab.index,
          children: <Widget>[
            HomeScreen(
              concealBalances: _concealBalances,
              onToggleConcealBalances: () {
                setState(() {
                  _concealBalances = !_concealBalances;
                });
              },
              onRequestAddTransaction: () => _openAddTransactionSheet(context),
              onViewAllTransactions: () {
                setState(() {
                  _selectedTab = MontraPrimaryTab.transaksi;
                });
              },
            ),
            TransactionsScreen(
              onRequestAddTransaction: () => _openAddTransactionSheet(context),
            ),
            const BudgetsScreen(),
            const InsightsScreen(),
          ],
        ),
      ),
      bottomNavigationBar: Container(
        decoration: BoxDecoration(
          color: tokens.surface,
          border: Border(
            top: BorderSide(
              color: tokens.outline,
              width: MontraShapes.outlineWidth,
            ),
          ),
        ),
        child: SafeArea(
          top: false,
          child: Padding(
            padding: const EdgeInsets.symmetric(
              horizontal: MontraSpacing.xs,
              vertical: MontraSpacing.xxs,
            ),
            child: Row(
              mainAxisAlignment: MainAxisAlignment.spaceAround,
              children: <Widget>[
                _NavButton(
                  tab: MontraPrimaryTab.beranda,
                  icon: Icons.home_outlined,
                  activeIcon: Icons.home,
                  isSelected: _selectedTab == MontraPrimaryTab.beranda,
                  onTap: () => setState(
                    () => _selectedTab = MontraPrimaryTab.beranda,
                  ),
                ),
                _NavButton(
                  tab: MontraPrimaryTab.transaksi,
                  icon: Icons.receipt_long_outlined,
                  activeIcon: Icons.receipt_long,
                  isSelected: _selectedTab == MontraPrimaryTab.transaksi,
                  onTap: () => setState(
                    () => _selectedTab = MontraPrimaryTab.transaksi,
                  ),
                ),
                _CenterAddButton(
                  onTap: () => _openAddTransactionSheet(context),
                ),
                _NavButton(
                  tab: MontraPrimaryTab.anggaran,
                  icon: Icons.pie_chart_outline,
                  activeIcon: Icons.pie_chart,
                  isSelected: _selectedTab == MontraPrimaryTab.anggaran,
                  onTap: () => setState(
                    () => _selectedTab = MontraPrimaryTab.anggaran,
                  ),
                ),
                _NavButton(
                  tab: MontraPrimaryTab.insight,
                  icon: Icons.insights_outlined,
                  activeIcon: Icons.insights,
                  isSelected: _selectedTab == MontraPrimaryTab.insight,
                  onTap: () => setState(
                    () => _selectedTab = MontraPrimaryTab.insight,
                  ),
                ),
              ],
            ),
          ),
        ),
      ),
    );
  }
}

class _NavButton extends StatelessWidget {
  const _NavButton({
    required this.tab,
    required this.icon,
    required this.activeIcon,
    required this.isSelected,
    required this.onTap,
  });

  final MontraPrimaryTab tab;
  final IconData icon;
  final IconData activeIcon;
  final bool isSelected;
  final VoidCallback onTap;

  @override
  Widget build(BuildContext context) {
    final tokens = Theme.of(context).extension<MontraColorTokens>()!;
    final color = isSelected ? tokens.primary : tokens.textSecondary;

    return Semantics(
      button: true,
      selected: isSelected,
      label: tab.semanticLabel,
      child: InkWell(
        borderRadius: MontraShapes.controlBorderRadius,
        onTap: onTap,
        child: ConstrainedBox(
          constraints: const BoxConstraints(
            minWidth: MontraSpacing.minTouchTarget,
            minHeight: MontraSpacing.minTouchTarget,
          ),
          child: Padding(
            padding: const EdgeInsets.symmetric(
              horizontal: MontraSpacing.xs,
              vertical: MontraSpacing.xxs,
            ),
            child: Column(
              mainAxisSize: MainAxisSize.min,
              children: <Widget>[
                Container(
                  padding: const EdgeInsets.symmetric(
                    horizontal: MontraSpacing.sm,
                    vertical: 2.0,
                  ),
                  decoration: BoxDecoration(
                    color: isSelected
                        ? tokens.primaryContainer
                        : Colors.transparent,
                    borderRadius: MontraShapes.controlBorderRadius,
                  ),
                  child: Icon(
                    isSelected ? activeIcon : icon,
                    size: 22.0,
                    color: color,
                  ),
                ),
                const SizedBox(height: 2.0),
                Text(
                  tab.label,
                  style: Theme.of(context).textTheme.labelMedium?.copyWith(
                        color: color,
                        fontWeight:
                            isSelected ? FontWeight.w700 : FontWeight.w500,
                      ),
                ),
              ],
            ),
          ),
        ),
      ),
    );
  }
}

class _CenterAddButton extends StatelessWidget {
  const _CenterAddButton({required this.onTap});

  final VoidCallback onTap;

  @override
  Widget build(BuildContext context) {
    final tokens = Theme.of(context).extension<MontraColorTokens>()!;

    return Semantics(
      button: true,
      label: 'Tambah transaksi baru',
      child: InkWell(
        borderRadius: MontraShapes.controlBorderRadius,
        onTap: onTap,
        child: ConstrainedBox(
          constraints: const BoxConstraints(
            minWidth: MontraSpacing.minTouchTarget,
            minHeight: MontraSpacing.minTouchTarget,
          ),
          child: Column(
            mainAxisSize: MainAxisSize.min,
            children: <Widget>[
              Container(
                width: 44.0,
                height: 32.0,
                decoration: BoxDecoration(
                  color: tokens.accent,
                  borderRadius: MontraShapes.controlBorderRadius,
                  border: Border.all(
                    color: tokens.outline,
                    width: MontraShapes.outlineWidth,
                  ),
                ),
                child: const Icon(
                  Icons.add,
                  color: Color(0xFF202B25),
                  size: 22.0,
                ),
              ),
              const SizedBox(height: 2.0),
              Text(
                'Tambah',
                style: Theme.of(context).textTheme.labelMedium?.copyWith(
                      color: tokens.textPrimary,
                      fontWeight: FontWeight.w700,
                    ),
              ),
            ],
          ),
        ),
      ),
    );
  }
}
