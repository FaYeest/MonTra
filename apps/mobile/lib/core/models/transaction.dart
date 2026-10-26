/// Transaction model representing an income or expense ledger record.
class TransactionModel {
  const TransactionModel({
    required this.id,
    required this.accountId,
    required this.categoryId,
    required this.type,
    required this.amount,
    required this.currency,
    required this.description,
    required this.transactionDate,
    this.notes,
    this.categoryName,
    this.categoryIcon,
    this.accountName,
  });

  final String id;
  final String accountId;
  final String categoryId;
  final String type; // 'income' or 'expense'
  final int amount; // Exact integer in minor/base units (IDR)
  final String currency;
  final String description;
  final DateTime transactionDate;
  final String? notes;
  final String? categoryName;
  final String? categoryIcon;
  final String? accountName;

  bool get isIncome => type == 'income';
  bool get isExpense => type == 'expense';

  factory TransactionModel.fromJson(Map<String, dynamic> json) {
    return TransactionModel(
      id: json['id'] as String? ?? '',
      accountId: json['account_id'] as String? ?? '',
      categoryId: json['category_id'] as String? ?? '',
      type: json['type'] as String? ?? 'expense',
      amount: (json['amount'] as num?)?.toInt() ?? 0,
      currency: json['currency'] as String? ?? 'IDR',
      description: json['description'] as String? ?? '',
      transactionDate: json['transaction_date'] != null
          ? DateTime.tryParse(json['transaction_date'] as String) ??
              DateTime.now()
          : DateTime.now(),
      notes: json['notes'] as String?,
      categoryName: json['category_name'] as String?,
      categoryIcon: json['category_icon'] as String?,
      accountName: json['account_name'] as String?,
    );
  }

  Map<String, dynamic> toJson() {
    return <String, dynamic>{
      'id': id,
      'account_id': accountId,
      'category_id': categoryId,
      'type': type,
      'amount': amount,
      'currency': currency,
      'description': description,
      'transaction_date': transactionDate.toIso8601String(),
      if (notes != null) 'notes': notes,
    };
  }
}
