import 'package:flutter/material.dart';

import '../../app/theme/app_colors.dart';
import '../../app/theme/app_shapes.dart';
import '../../app/theme/app_spacing.dart';
import '../../core/formatting/date_formatter.dart';
import '../../core/formatting/idr_formatter.dart';
import '../../core/models/transaction.dart';
import '../../core/models/transfer.dart';
import '../../core/widgets/montra_scope.dart';

enum TransactionFilter {
  all,
  expense,
  income,
  transfer,
}

/// Transactions (`Transaksi`) screen defined in `docs/design.md` Section 6.2.
class TransactionsScreen extends StatefulWidget {
  const TransactionsScreen({
    super.key,
    required this.onRequestAddTransaction,
  });

  final VoidCallback onRequestAddTransaction;

  @override
  State<TransactionsScreen> createState() => _TransactionsScreenState();
}

class _TransactionsScreenState extends State<TransactionsScreen> {
  final _searchController = TextEditingController();
  TransactionFilter _currentFilter = TransactionFilter.all;
  String _searchQuery = '';

  @override
  void dispose() {
    _searchController.dispose();
    super.dispose();
  }

  @override
  Widget build(BuildContext context) {
    final tokens = Theme.of(context).extension<MontraColorTokens>()!;
    final textTheme = Theme.of(context).textTheme;
    final ledger = MontraScope.ledgerOf(context);

    // Filter items
    final transactions = ledger.transactions;
    final transfers = ledger.transfers;

    // Combine transactions and transfers into unified display list
    final allEntries = <_LedgerItem>[];

    if (_currentFilter != TransactionFilter.transfer) {
      for (final tx in transactions) {
        if (_currentFilter == TransactionFilter.income && !tx.isIncome) continue;
        if (_currentFilter == TransactionFilter.expense && !tx.isExpense) continue;

        if (_searchQuery.isNotEmpty) {
          final q = _searchQuery.toLowerCase();
          final descMatches = tx.description.toLowerCase().contains(q);
          final catMatches = (tx.categoryName ?? '').toLowerCase().contains(q);
          final accMatches = (tx.accountName ?? '').toLowerCase().contains(q);
          if (!descMatches && !catMatches && !accMatches) continue;
        }

        allEntries.add(_LedgerItem.fromTransaction(tx));
      }
    }

    if (_currentFilter == TransactionFilter.all ||
        _currentFilter == TransactionFilter.transfer) {
      for (final tr in transfers) {
        if (_searchQuery.isNotEmpty) {
          final q = _searchQuery.toLowerCase();
          final srcMatches =
              (tr.sourceAccountName ?? '').toLowerCase().contains(q);
          final dstMatches =
              (tr.destinationAccountName ?? '').toLowerCase().contains(q);
          final notesMatches = (tr.notes ?? '').toLowerCase().contains(q);
          if (!srcMatches && !dstMatches && !notesMatches && !'transfer'.contains(q)) {
            continue;
          }
        }

        allEntries.add(_LedgerItem.fromTransfer(tr));
      }
    }

    // Sort by date descending
    allEntries.sort((a, b) => b.date.compareTo(a.date));

    // Group by Date formatted string
    final grouped = <String, List<_LedgerItem>>{};

    for (final entry in allEntries) {
      final groupKey = MontraDateFormatter.formatRelative(entry.date);
      grouped.putIfAbsent(groupKey, () => <_LedgerItem>[]).add(entry);
    }

    return RefreshIndicator(
      onRefresh: () => ledger.refreshAll(),
      child: ListView(
        padding: const EdgeInsets.symmetric(
          horizontal: MontraSpacing.pageHorizontal,
          vertical: MontraSpacing.md,
        ),
        children: <Widget>[
          // Header
          Text('Transaksi', style: textTheme.headlineLarge),
          const SizedBox(height: MontraSpacing.xxs),
          Text(
            'Cari, periksa, dan kelola catatan transaksi Anda.',
            style: textTheme.bodySmall,
          ),
          const SizedBox(height: MontraSpacing.md),

          // Period Summary Row (Income vs Expense)
          Row(
            children: <Widget>[
              Expanded(
                child: Container(
                  padding: const EdgeInsets.all(MontraSpacing.md),
                  decoration: BoxDecoration(
                    color: tokens.income.withValues(alpha: 0.08),
                    borderRadius: BorderRadius.circular(MontraShapes.radiusMd),
                    border: Border.all(
                      color: tokens.income.withValues(alpha: 0.25),
                      width: MontraShapes.outlineWidth,
                    ),
                  ),
                  child: Column(
                    crossAxisAlignment: CrossAxisAlignment.start,
                    children: <Widget>[
                      Row(
                        children: <Widget>[
                          Icon(Icons.arrow_downward_rounded,
                              size: 16, color: tokens.income),
                          const SizedBox(width: MontraSpacing.xxs),
                          Text(
                            'Pemasukan',
                            style: textTheme.bodySmall?.copyWith(
                              color: tokens.income,
                              fontWeight: FontWeight.w600,
                            ),
                          ),
                        ],
                      ),
                      const SizedBox(height: MontraSpacing.xxs),
                      Text(
                        IdrFormatter.formatSigned(ledger.totalIncome, isExpense: false),
                        style: textTheme.titleMedium?.copyWith(
                          color: tokens.income,
                          fontWeight: FontWeight.bold,
                        ),
                      ),
                    ],
                  ),
                ),
              ),
              const SizedBox(width: MontraSpacing.md),
              Expanded(
                child: Container(
                  padding: const EdgeInsets.all(MontraSpacing.md),
                  decoration: BoxDecoration(
                    color: tokens.expense.withValues(alpha: 0.08),
                    borderRadius: BorderRadius.circular(MontraShapes.radiusMd),
                    border: Border.all(
                      color: tokens.expense.withValues(alpha: 0.25),
                      width: MontraShapes.outlineWidth,
                    ),
                  ),
                  child: Column(
                    crossAxisAlignment: CrossAxisAlignment.start,
                    children: <Widget>[
                      Row(
                        children: <Widget>[
                          Icon(Icons.arrow_upward_rounded,
                              size: 16, color: tokens.expense),
                          const SizedBox(width: MontraSpacing.xxs),
                          Text(
                            'Pengeluaran',
                            style: textTheme.bodySmall?.copyWith(
                              color: tokens.expense,
                              fontWeight: FontWeight.w600,
                            ),
                          ),
                        ],
                      ),
                      const SizedBox(height: MontraSpacing.xxs),
                      Text(
                        IdrFormatter.formatSigned(ledger.totalExpense, isExpense: true),
                        style: textTheme.titleMedium?.copyWith(
                          color: tokens.expense,
                          fontWeight: FontWeight.bold,
                        ),
                      ),
                    ],
                  ),
                ),
              ),
            ],
          ),
          const SizedBox(height: MontraSpacing.md),

          // Search Field
          TextField(
            controller: _searchController,
            onChanged: (val) {
              setState(() {
                _searchQuery = val.trim();
              });
            },
            decoration: InputDecoration(
              hintText: 'Cari keterangan, kategori, atau akun...',
              prefixIcon: const Icon(Icons.search),
              suffixIcon: _searchQuery.isNotEmpty
                  ? IconButton(
                      icon: const Icon(Icons.clear),
                      onPressed: () {
                        _searchController.clear();
                        setState(() {
                          _searchQuery = '';
                        });
                      },
                    )
                  : null,
            ),
          ),
          const SizedBox(height: MontraSpacing.sm),

          // Filter Chips
          SingleChildScrollView(
            scrollDirection: Axis.horizontal,
            child: Row(
              children: <Widget>[
                _buildFilterChip('Semua', TransactionFilter.all),
                const SizedBox(width: MontraSpacing.xs),
                _buildFilterChip('Pengeluaran', TransactionFilter.expense),
                const SizedBox(width: MontraSpacing.xs),
                _buildFilterChip('Pemasukan', TransactionFilter.income),
                const SizedBox(width: MontraSpacing.xs),
                _buildFilterChip('Transfer', TransactionFilter.transfer),
              ],
            ),
          ),
          const SizedBox(height: MontraSpacing.md),

          // Grouped List or Empty State
          if (allEntries.isEmpty) ...[
            Card(
              child: Padding(
                padding: const EdgeInsets.all(MontraSpacing.xl),
                child: Column(
                  crossAxisAlignment: CrossAxisAlignment.center,
                  children: <Widget>[
                    Icon(
                      _searchQuery.isNotEmpty
                          ? Icons.search_off_rounded
                          : Icons.receipt_long_outlined,
                      size: 48,
                      color: tokens.textSecondary,
                    ),
                    const SizedBox(height: MontraSpacing.sm),
                    Text(
                      _searchQuery.isNotEmpty
                          ? 'Tidak ada transaksi yang cocok'
                          : 'Belum ada transaksi tercatat',
                      style: textTheme.titleMedium,
                      textAlign: TextAlign.center,
                    ),
                    const SizedBox(height: MontraSpacing.xs),
                    Text(
                      _searchQuery.isNotEmpty
                          ? 'Coba gunakan kata kunci pencarian yang berbeda atau reset filter.'
                          : 'Mulai catat pengeluaran, pemasukan, atau transfer akun Anda.',
                      style: textTheme.bodySmall?.copyWith(
                        color: tokens.textSecondary,
                      ),
                      textAlign: TextAlign.center,
                    ),
                    if (_searchQuery.isEmpty) ...[
                      const SizedBox(height: MontraSpacing.md),
                      ElevatedButton.icon(
                        onPressed: widget.onRequestAddTransaction,
                        icon: const Icon(Icons.add),
                        label: const Text('Catat Transaksi'),
                      ),
                    ],
                  ],
                ),
              ),
            ),
          ] else ...[
            for (final entry in grouped.entries) ...[
              Padding(
                padding: const EdgeInsets.only(
                  top: MontraSpacing.sm,
                  bottom: MontraSpacing.xs,
                ),
                child: Text(
                  entry.key,
                  style: textTheme.labelLarge?.copyWith(
                    color: tokens.textSecondary,
                    fontWeight: FontWeight.bold,
                  ),
                ),
              ),
              Card(
                margin: const EdgeInsets.only(bottom: MontraSpacing.xs),
                child: ListView.separated(
                  shrinkWrap: true,
                  physics: const NeverScrollableScrollPhysics(),
                  itemCount: entry.value.length,
                  separatorBuilder: (_, __) => Divider(
                    height: 1,
                    color: tokens.outline.withValues(alpha: 0.5),
                  ),
                  itemBuilder: (context, idx) {
                    final item = entry.value[idx];
                    return _buildTransactionTile(context, item, tokens, textTheme);
                  },
                ),
              ),
            ],
          ],
        ],
      ),
    );
  }

  Widget _buildFilterChip(String label, TransactionFilter filter) {
    final isSelected = _currentFilter == filter;
    return FilterChip(
      label: Text(label),
      selected: isSelected,
      onSelected: (_) {
        setState(() {
          _currentFilter = filter;
        });
      },
    );
  }

  Widget _buildTransactionTile(
    BuildContext context,
    _LedgerItem item,
    MontraColorTokens tokens,
    TextTheme textTheme,
  ) {
    final IconData iconData;
    final Color iconColor;
    final Color iconBg;
    final String amountText;
    final Color amountColor;

    if (item.isTransfer) {
      iconData = Icons.swap_horiz_rounded;
      iconColor = tokens.primary;
      iconBg = tokens.primary.withValues(alpha: 0.12);
      amountText = IdrFormatter.format(item.amount);
      amountColor = tokens.textPrimary;
    } else if (item.isIncome) {
      iconData = Icons.arrow_downward_rounded;
      iconColor = tokens.income;
      iconBg = tokens.income.withValues(alpha: 0.12);
      amountText = IdrFormatter.formatSigned(item.amount, isExpense: false);
      amountColor = tokens.income;
    } else {
      iconData = Icons.arrow_upward_rounded;
      iconColor = tokens.expense;
      iconBg = tokens.expense.withValues(alpha: 0.12);
      amountText = IdrFormatter.formatSigned(item.amount, isExpense: true);
      amountColor = tokens.expense;
    }

    return ListTile(
      leading: Container(
        padding: const EdgeInsets.all(MontraSpacing.xs),
        decoration: BoxDecoration(
          color: iconBg,
          borderRadius: BorderRadius.circular(MontraShapes.radiusSm),
        ),
        child: Icon(iconData, color: iconColor, size: 20),
      ),
      title: Text(
        item.title,
        style: textTheme.bodyMedium?.copyWith(
          fontWeight: FontWeight.w600,
        ),
      ),
      subtitle: Text(
        item.subtitle,
        style: textTheme.bodySmall?.copyWith(
          color: tokens.textSecondary,
        ),
      ),
      trailing: Text(
        amountText,
        style: textTheme.titleMedium?.copyWith(
          fontWeight: FontWeight.bold,
          color: amountColor,
        ),
      ),
    );
  }
}

class _LedgerItem {
  _LedgerItem({
    required this.id,
    required this.title,
    required this.subtitle,
    required this.amount,
    required this.date,
    required this.isTransfer,
    required this.isIncome,
  });

  final String id;
  final String title;
  final String subtitle;
  final int amount;
  final DateTime date;
  final bool isTransfer;
  final bool isIncome;

  factory _LedgerItem.fromTransaction(TransactionModel tx) {
    final cat = tx.categoryName ?? (tx.isIncome ? 'Pemasukan' : 'Pengeluaran');
    final acc = tx.accountName != null ? ' • ${tx.accountName}' : '';
    return _LedgerItem(
      id: tx.id,
      title: tx.description,
      subtitle: '$cat$acc',
      amount: tx.amount,
      date: tx.transactionDate,
      isTransfer: false,
      isIncome: tx.isIncome,
    );
  }

  factory _LedgerItem.fromTransfer(TransferModel tr) {
    final src = tr.sourceAccountName ?? 'Sumber';
    final dst = tr.destinationAccountName ?? 'Tujuan';
    final fee = tr.feeAmount > 0
        ? ' (Biaya: ${IdrFormatter.format(tr.feeAmount)})'
        : '';
    return _LedgerItem(
      id: tr.id,
      title: 'Transfer: $src → $dst',
      subtitle: (tr.notes != null && tr.notes!.isNotEmpty)
          ? '${tr.notes!}$fee'
          : 'Transfer Antar-Akun$fee',
      amount: tr.amount,
      date: tr.transferDate,
      isTransfer: true,
      isIncome: false,
    );
  }
}
