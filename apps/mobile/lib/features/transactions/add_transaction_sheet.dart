import 'package:flutter/material.dart';

import '../../app/theme/app_colors.dart';
import '../../app/theme/app_shapes.dart';
import '../../app/theme/app_spacing.dart';
import '../../core/formatting/date_formatter.dart';
import '../../core/widgets/montra_scope.dart';

enum TransactionFormType {
  expense,
  income,
  transfer,
}

/// Modal bottom sheet for adding an expense, income, or internal transfer.
/// Implements `docs/design.md` Section 6.3 & Section 7.6.
class AddTransactionSheet extends StatefulWidget {
  const AddTransactionSheet({super.key});

  static Future<void> show(BuildContext context) {
    return showModalBottomSheet<void>(
      context: context,
      isScrollControlled: true,
      useSafeArea: true,
      backgroundColor: Colors.transparent,
      builder: (_) => const AddTransactionSheet(),
    );
  }

  @override
  State<AddTransactionSheet> createState() => _AddTransactionSheetState();
}

class _AddTransactionSheetState extends State<AddTransactionSheet> {
  final _formKey = GlobalKey<FormState>();
  final _amountController = TextEditingController();
  final _feeController = TextEditingController();
  final _descriptionController = TextEditingController();
  final _notesController = TextEditingController();

  TransactionFormType _formType = TransactionFormType.expense;
  String? _selectedAccountId;
  String? _selectedCategoryId;
  String? _selectedDestinationAccountId;
  DateTime _selectedDate = DateTime.now();
  bool _isSubmitting = false;
  String? _errorMessage;

  @override
  void initState() {
    super.initState();
    WidgetsBinding.instance.addPostFrameCallback((_) {
      _initSelections();
    });
  }

  void _initSelections() {
    final ledger = MontraScope.ledgerOf(context);
    final accounts = ledger.accounts;
    if (accounts.isNotEmpty) {
      setState(() {
        _selectedAccountId = accounts.first.id;
        if (accounts.length > 1) {
          _selectedDestinationAccountId = accounts[1].id;
        }
      });
    }

    _updateDefaultCategory();
  }

  void _updateDefaultCategory() {
    final ledger = MontraScope.ledgerOf(context);
    final targetType =
        _formType == TransactionFormType.income ? 'income' : 'expense';
    final categories = ledger.categoriesForType(targetType);
    if (categories.isNotEmpty) {
      setState(() {
        _selectedCategoryId = categories.first.id;
      });
    } else {
      setState(() {
        _selectedCategoryId = null;
      });
    }
  }

  @override
  void dispose() {
    _amountController.dispose();
    _feeController.dispose();
    _descriptionController.dispose();
    _notesController.dispose();
    super.dispose();
  }

  Future<void> _pickDate() async {
    final now = DateTime.now();
    final picked = await showDatePicker(
      context: context,
      initialDate: _selectedDate,
      firstDate: DateTime(2020),
      lastDate: DateTime(now.year + 1),
    );
    if (picked != null) {
      setState(() {
        _selectedDate = picked;
      });
    }
  }

  Future<void> _handleSubmit() async {
    if (!_formKey.currentState!.validate()) return;

    final rawAmount = int.tryParse(
      _amountController.text.replaceAll(RegExp(r'[^0-9]'), ''),
    );
    if (rawAmount == null || rawAmount <= 0) {
      setState(() {
        _errorMessage = 'Nominal harus lebih dari 0';
      });
      return;
    }

    final ledger = MontraScope.ledgerOf(context);
    setState(() {
      _isSubmitting = true;
      _errorMessage = null;
    });

    bool success = false;
    if (_formType == TransactionFormType.transfer) {
      if (_selectedAccountId == _selectedDestinationAccountId) {
        setState(() {
          _errorMessage = 'Akun sumber dan akun tujuan tidak boleh sama.';
          _isSubmitting = false;
        });
        return;
      }

      final feeAmount = int.tryParse(
            _feeController.text.replaceAll(RegExp(r'[^0-9]'), ''),
          ) ??
          0;

      success = await ledger.createTransfer(
        sourceAccountId: _selectedAccountId!,
        destinationAccountId: _selectedDestinationAccountId!,
        amount: rawAmount,
        feeAmount: feeAmount,
        transferDate: _selectedDate,
        notes: _notesController.text.trim().isNotEmpty
            ? _notesController.text.trim()
            : null,
      );
    } else {
      if (_selectedCategoryId == null) {
        setState(() {
          _errorMessage = 'Silakan pilih kategori transaksi.';
          _isSubmitting = false;
        });
        return;
      }

      final typeStr =
          _formType == TransactionFormType.income ? 'income' : 'expense';

      success = await ledger.createTransaction(
        accountId: _selectedAccountId!,
        categoryId: _selectedCategoryId!,
        type: typeStr,
        amount: rawAmount,
        description: _descriptionController.text.trim(),
        transactionDate: _selectedDate,
        notes: _notesController.text.trim().isNotEmpty
            ? _notesController.text.trim()
            : null,
      );
    }

    if (!mounted) return;

    setState(() {
      _isSubmitting = false;
      _errorMessage = ledger.errorMessage;
    });

    if (success) {
      Navigator.of(context).pop();
      ScaffoldMessenger.of(context).showSnackBar(
        SnackBar(
          content: Text(
            _formType == TransactionFormType.transfer
                ? 'Transfer berhasil dicatat'
                : 'Transaksi berhasil disimpan',
          ),
          behavior: SnackBarBehavior.floating,
        ),
      );
    }
  }

  @override
  Widget build(BuildContext context) {
    final tokens = Theme.of(context).extension<MontraColorTokens>()!;
    final textTheme = Theme.of(context).textTheme;
    final ledger = MontraScope.ledgerOf(context);
    final accounts = ledger.accounts;
    final targetCategoryType =
        _formType == TransactionFormType.income ? 'income' : 'expense';
    final categories = ledger.categoriesForType(targetCategoryType);
    final bottomInset = MediaQuery.of(context).viewInsets.bottom;

    return Container(
      decoration: BoxDecoration(
        color: tokens.surface,
        borderRadius: MontraShapes.bottomSheetBorderRadius,
      ),
      padding: EdgeInsets.fromLTRB(
        MontraSpacing.pageHorizontal,
        MontraSpacing.md,
        MontraSpacing.pageHorizontal,
        MontraSpacing.md + bottomInset,
      ),
      child: Form(
        key: _formKey,
        child: SingleChildScrollView(
          child: Column(
            mainAxisSize: MainAxisSize.min,
            crossAxisAlignment: CrossAxisAlignment.stretch,
            children: <Widget>[
              // Drag handle
              Center(
                child: Container(
                  width: 40,
                  height: 4,
                  decoration: BoxDecoration(
                    color: tokens.outline,
                    borderRadius: BorderRadius.circular(2),
                  ),
                ),
              ),
              const SizedBox(height: MontraSpacing.sm),

              // Title row with close button
              Row(
                mainAxisAlignment: MainAxisAlignment.spaceBetween,
                children: <Widget>[
                  Text(
                    'Catat Transaksi',
                    style: textTheme.titleLarge?.copyWith(
                      fontWeight: FontWeight.bold,
                    ),
                  ),
                  IconButton(
                    icon: const Icon(Icons.close),
                    tooltip: 'Tutup',
                    onPressed: () => Navigator.of(context).pop(),
                  ),
                ],
              ),
              const SizedBox(height: MontraSpacing.sm),

              // Segmented switcher (Pengeluaran, Pemasukan, Transfer)
              SegmentedButton<TransactionFormType>(
                segments: const <ButtonSegment<TransactionFormType>>[
                  ButtonSegment(
                    value: TransactionFormType.expense,
                    label: Text('Pengeluaran'),
                    icon: Icon(Icons.arrow_upward_rounded, size: 16),
                  ),
                  ButtonSegment(
                    value: TransactionFormType.income,
                    label: Text('Pemasukan'),
                    icon: Icon(Icons.arrow_downward_rounded, size: 16),
                  ),
                  ButtonSegment(
                    value: TransactionFormType.transfer,
                    label: Text('Transfer'),
                    icon: Icon(Icons.swap_horiz_rounded, size: 16),
                  ),
                ],
                selected: <TransactionFormType>{_formType},
                onSelectionChanged: (Set<TransactionFormType> selection) {
                  setState(() {
                    _formType = selection.first;
                    _errorMessage = null;
                  });
                  _updateDefaultCategory();
                },
              ),
              const SizedBox(height: MontraSpacing.md),

              // Error display
              if (_errorMessage != null) ...[
                Container(
                  padding: const EdgeInsets.all(MontraSpacing.sm),
                  decoration: BoxDecoration(
                    color: tokens.expense.withValues(alpha: 0.1),
                    borderRadius: BorderRadius.circular(MontraShapes.radiusSm),
                    border: Border.all(
                      color: tokens.expense.withValues(alpha: 0.3),
                      width: MontraShapes.outlineWidth,
                    ),
                  ),
                  child: Row(
                    children: <Widget>[
                      Icon(Icons.error_outline, color: tokens.expense, size: 18),
                      const SizedBox(width: MontraSpacing.xs),
                      Expanded(
                        child: Text(
                          _errorMessage!,
                          style: textTheme.bodySmall?.copyWith(
                            color: tokens.expense,
                            fontWeight: FontWeight.w500,
                          ),
                        ),
                      ),
                    ],
                  ),
                ),
                const SizedBox(height: MontraSpacing.sm),
              ],

              // Amount input
              TextFormField(
                controller: _amountController,
                keyboardType: TextInputType.number,
                autofocus: true,
                style: textTheme.headlineMedium?.copyWith(
                  fontWeight: FontWeight.bold,
                  color: _formType == TransactionFormType.expense
                      ? tokens.expense
                      : (_formType == TransactionFormType.income
                          ? tokens.income
                          : tokens.primary),
                ),
                decoration: const InputDecoration(
                  labelText: 'Nominal Transaksi',
                  prefixText: 'Rp ',
                  hintText: '0',
                ),
                validator: (val) {
                  if (val == null || val.trim().isEmpty) {
                    return 'Nominal wajib diisi';
                  }
                  final numVal = int.tryParse(
                    val.replaceAll(RegExp(r'[^0-9]'), ''),
                  );
                  if (numVal == null || numVal <= 0) {
                    return 'Nominal harus lebih dari 0';
                  }
                  return null;
                },
              ),
              const SizedBox(height: MontraSpacing.md),

              // Description input (for income & expense)
              if (_formType != TransactionFormType.transfer) ...[
                TextFormField(
                  controller: _descriptionController,
                  textInputAction: TextInputAction.next,
                  decoration: const InputDecoration(
                    labelText: 'Keterangan / Merchant',
                    hintText: 'Contoh: Nasi Padang, Belanja Mingguan',
                    prefixIcon: Icon(Icons.edit_note_outlined),
                  ),
                  validator: (val) {
                    if (val == null || val.trim().isEmpty) {
                      return 'Keterangan wajib diisi';
                    }
                    return null;
                  },
                ),
                const SizedBox(height: MontraSpacing.md),
              ],

              // Source Account selection
              if (accounts.isNotEmpty) ...[
                DropdownButtonFormField<String>(
                  initialValue: _selectedAccountId,
                  decoration: InputDecoration(
                    labelText: _formType == TransactionFormType.transfer
                        ? 'Akun Sumber (Dari)'
                        : 'Akun / Dompet',
                    prefixIcon: const Icon(Icons.account_balance_wallet_outlined),
                  ),
                  items: accounts.map((acc) {
                    return DropdownMenuItem<String>(
                      value: acc.id,
                      child: Text('${acc.name} (${acc.currency})'),
                    );
                  }).toList(),
                  onChanged: (val) {
                    setState(() {
                      _selectedAccountId = val;
                    });
                  },
                  validator: (val) =>
                      val == null ? 'Pilih akun pembayaran' : null,
                ),
                const SizedBox(height: MontraSpacing.md),
              ],

              // Category selection (for income & expense)
              if (_formType != TransactionFormType.transfer &&
                  categories.isNotEmpty) ...[
                DropdownButtonFormField<String>(
                  initialValue: _selectedCategoryId,
                  decoration: const InputDecoration(
                    labelText: 'Kategori',
                    prefixIcon: Icon(Icons.category_outlined),
                  ),
                  items: categories.map((cat) {
                    return DropdownMenuItem<String>(
                      value: cat.id,
                      child: Text(cat.name),
                    );
                  }).toList(),
                  onChanged: (val) {
                    setState(() {
                      _selectedCategoryId = val;
                    });
                  },
                  validator: (val) => val == null ? 'Pilih kategori' : null,
                ),
                const SizedBox(height: MontraSpacing.md),
              ],

              // Transfer Destination Account & Fee
              if (_formType == TransactionFormType.transfer) ...[
                DropdownButtonFormField<String>(
                  initialValue: _selectedDestinationAccountId,
                  decoration: const InputDecoration(
                    labelText: 'Akun Tujuan (Ke)',
                    prefixIcon: Icon(Icons.move_to_inbox_outlined),
                  ),
                  items: accounts.map((acc) {
                    return DropdownMenuItem<String>(
                      value: acc.id,
                      child: Text('${acc.name} (${acc.currency})'),
                    );
                  }).toList(),
                  onChanged: (val) {
                    setState(() {
                      _selectedDestinationAccountId = val;
                    });
                  },
                  validator: (val) {
                    if (val == null) return 'Pilih akun tujuan';
                    if (val == _selectedAccountId) {
                      return 'Akun tujuan tidak boleh sama dengan akun sumber';
                    }
                    return null;
                  },
                ),
                const SizedBox(height: MontraSpacing.md),

                // Transfer fee
                TextFormField(
                  controller: _feeController,
                  keyboardType: TextInputType.number,
                  decoration: const InputDecoration(
                    labelText: 'Biaya Admin / Transfer (Opsional)',
                    prefixText: 'Rp ',
                    hintText: '0',
                    prefixIcon: Icon(Icons.receipt_outlined),
                  ),
                ),
                const SizedBox(height: MontraSpacing.md),
              ],

              // Date Picker Field
              InkWell(
                onTap: _pickDate,
                borderRadius: BorderRadius.circular(MontraShapes.radiusSm),
                child: InputDecorator(
                  decoration: const InputDecoration(
                    labelText: 'Tanggal Transaksi',
                    prefixIcon: Icon(Icons.calendar_today_outlined),
                  ),
                  child: Text(
                    MontraDateFormatter.formatFull(_selectedDate),
                    style: textTheme.bodyMedium,
                  ),
                ),
              ),
              const SizedBox(height: MontraSpacing.md),

              // Notes field
              TextFormField(
                controller: _notesController,
                textInputAction: TextInputAction.done,
                decoration: const InputDecoration(
                  labelText: 'Catatan Tambahan (Opsional)',
                  hintText: 'Tambahkan catatan jika perlu',
                  prefixIcon: Icon(Icons.notes_outlined),
                ),
              ),
              const SizedBox(height: MontraSpacing.lg),

              // Submit button
              SizedBox(
                height: MontraSpacing.minTouchTarget,
                child: ElevatedButton(
                  onPressed: _isSubmitting ? null : _handleSubmit,
                  child: _isSubmitting
                      ? const SizedBox(
                          width: 24,
                          height: 24,
                          child: CircularProgressIndicator(
                            strokeWidth: 2.5,
                            color: Colors.white,
                          ),
                        )
                      : Text(
                          _formType == TransactionFormType.transfer
                              ? 'Simpan Transfer'
                              : 'Simpan Transaksi',
                        ),
                ),
              ),
              const SizedBox(height: MontraSpacing.xs),
            ],
          ),
        ),
      ),
    );
  }
}
