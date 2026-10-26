import 'package:flutter/foundation.dart';

import '../models/account.dart';
import '../models/category.dart';
import '../models/transaction.dart';
import '../models/transfer.dart';
import '../networking/api_client.dart';
import '../networking/api_exception.dart';

/// Controller for managing ledger records: accounts, categories, transactions, and transfers.
class LedgerController extends ChangeNotifier {
  LedgerController({required ApiClient apiClient}) : _apiClient = apiClient;

  final ApiClient _apiClient;

  List<AccountModel> _accounts = <AccountModel>[];
  List<CategoryModel> _categories = <CategoryModel>[];
  List<TransactionModel> _transactions = <TransactionModel>[];
  List<TransferModel> _transfers = <TransferModel>[];

  bool _isLoading = false;
  String? _errorMessage;

  List<AccountModel> get accounts => List.unmodifiable(_accounts);
  List<CategoryModel> get categories => List.unmodifiable(_categories);
  List<TransactionModel> get transactions => List.unmodifiable(_transactions);
  List<TransferModel> get transfers => List.unmodifiable(_transfers);
  bool get isLoading => _isLoading;
  String? get errorMessage => _errorMessage;

  /// Deterministic balance sum across all non-archived user accounts.
  int get totalBalance {
    int total = 0;
    for (final acc in _accounts) {
      if (!acc.isArchived) {
        total += acc.balance;
      }
    }
    return total;
  }

  /// Deterministic income sum from recorded transactions.
  int get totalIncome {
    int sum = 0;
    for (final tx in _transactions) {
      if (tx.isIncome) {
        sum += tx.amount;
      }
    }
    return sum;
  }

  /// Deterministic expense sum from recorded transactions + transfer fees.
  int get totalExpense {
    int sum = 0;
    for (final tx in _transactions) {
      if (tx.isExpense) {
        sum += tx.amount;
      }
    }
    for (final tr in _transfers) {
      sum += tr.feeAmount;
    }
    return sum;
  }

  /// Filter categories by transaction type ('expense' or 'income').
  List<CategoryModel> categoriesForType(String type) {
    return _categories.where((c) => c.type == type).toList();
  }

  /// Refresh all ledger state concurrently.
  Future<void> refreshAll() async {
    _isLoading = true;
    _errorMessage = null;
    notifyListeners();

    try {
      final results = await Future.wait([
        _fetchAccounts(),
        _fetchCategories(),
        _fetchTransactions(),
        _fetchTransfers(),
      ]);

      _accounts = results[0] as List<AccountModel>;
      _categories = results[1] as List<CategoryModel>;
      _transactions = results[2] as List<TransactionModel>;
      _transfers = results[3] as List<TransferModel>;

      if (_accounts.isEmpty) {
        await ensureDefaultAccount();
      }
      _errorMessage = null;
    } on ApiException catch (e) {
      _errorMessage = e.message;
    } catch (e) {
      _errorMessage = 'Gagal memuat data: ${e.toString()}';
    } finally {
      _isLoading = false;
      notifyListeners();
    }
  }

  /// Ensure at least one default account exists. If user has no accounts,
  /// creates a default "Kas Utama" account via the API.
  Future<AccountModel?> ensureDefaultAccount() async {
    if (_accounts.isNotEmpty) {
      return _accounts.first;
    }
    try {
      final res = await _apiClient.post(
        '/accounts',
        body: <String, dynamic>{
          'name': 'Kas Utama',
          'account_type': 'cash',
          'opening_balance': 0,
          'currency': 'IDR',
        },
      );
      if (res is Map<String, dynamic>) {
        final newAcc = AccountModel.fromJson(res);
        _accounts = <AccountModel>[newAcc];
        notifyListeners();
        return newAcc;
      }
    } on ApiException catch (e) {
      _errorMessage = e.message;
      notifyListeners();
    } catch (e) {
      _errorMessage = 'Gagal membuat akun kas default: ${e.toString()}';
      notifyListeners();
    }
    return null;
  }

  Future<List<AccountModel>> _fetchAccounts() async {
    final res = await _apiClient.get('/accounts');
    if (res is List) {
      return res
          .map((item) => AccountModel.fromJson(item as Map<String, dynamic>))
          .toList();
    }
    return <AccountModel>[];
  }

  Future<List<CategoryModel>> _fetchCategories() async {
    final res = await _apiClient.get('/categories');
    if (res is List) {
      return res
          .map((item) => CategoryModel.fromJson(item as Map<String, dynamic>))
          .toList();
    }
    return <CategoryModel>[];
  }

  Future<List<TransactionModel>> _fetchTransactions() async {
    final res = await _apiClient.get('/transactions', queryParams: {'limit': '100'});
    if (res is List) {
      return res
          .map((item) => TransactionModel.fromJson(item as Map<String, dynamic>))
          .toList();
    }
    return <TransactionModel>[];
  }

  Future<List<TransferModel>> _fetchTransfers() async {
    final res = await _apiClient.get('/transfers', queryParams: {'limit': '100'});
    if (res is List) {
      return res
          .map((item) => TransferModel.fromJson(item as Map<String, dynamic>))
          .toList();
    }
    return <TransferModel>[];
  }

  /// Create a new transaction (pemasukan atau pengeluaran).
  Future<bool> createTransaction({
    required String accountId,
    required String categoryId,
    required String type,
    required int amount,
    required String description,
    DateTime? transactionDate,
    String? notes,
  }) async {
    final date = transactionDate ?? DateTime.now();
    final idempotencyKey = 'tx_${DateTime.now().microsecondsSinceEpoch}';

    try {
      final res = await _apiClient.post(
        '/transactions',
        extraHeaders: {'Idempotency-Key': idempotencyKey},
        body: <String, dynamic>{
          'account_id': accountId,
          'category_id': categoryId,
          'type': type,
          'amount': amount,
          'currency': 'IDR',
          'description': description.trim(),
          'transaction_date': date.toIso8601String(),
          if (notes != null && notes.trim().isNotEmpty) 'notes': notes.trim(),
        },
      );

      if (res is Map<String, dynamic>) {
        final newTx = TransactionModel.fromJson(res);
        _transactions.insert(0, newTx);
        // Refresh accounts to reflect updated balances from backend
        final updatedAccounts = await _fetchAccounts();
        _accounts = updatedAccounts;
        notifyListeners();
        return true;
      }
      return false;
    } on ApiException catch (e) {
      _errorMessage = e.message;
      notifyListeners();
      return false;
    } catch (e) {
      _errorMessage = 'Gagal menyimpan transaksi: ${e.toString()}';
      notifyListeners();
      return false;
    }
  }

  /// Create an atomic internal transfer between two accounts.
  Future<bool> createTransfer({
    required String sourceAccountId,
    required String destinationAccountId,
    required int amount,
    int feeAmount = 0,
    DateTime? transferDate,
    String? notes,
  }) async {
    final date = transferDate ?? DateTime.now();
    final idempotencyKey = 'tr_${DateTime.now().microsecondsSinceEpoch}';

    try {
      final res = await _apiClient.post(
        '/transfers',
        extraHeaders: {'Idempotency-Key': idempotencyKey},
        body: <String, dynamic>{
          'source_account_id': sourceAccountId,
          'destination_account_id': destinationAccountId,
          'amount': amount,
          'fee_amount': feeAmount,
          'currency': 'IDR',
          'transfer_date': date.toIso8601String(),
          if (notes != null && notes.trim().isNotEmpty) 'notes': notes.trim(),
        },
      );

      if (res is Map<String, dynamic>) {
        final newTr = TransferModel.fromJson(res);
        _transfers.insert(0, newTr);
        // Refresh accounts to reflect updated balances from backend
        final updatedAccounts = await _fetchAccounts();
        _accounts = updatedAccounts;
        notifyListeners();
        return true;
      }
      return false;
    } on ApiException catch (e) {
      _errorMessage = e.message;
      notifyListeners();
      return false;
    } catch (e) {
      _errorMessage = 'Gagal menyimpan transfer: ${e.toString()}';
      notifyListeners();
      return false;
    }
  }

  void clearError() {
    _errorMessage = null;
    notifyListeners();
  }
}
