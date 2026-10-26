/// Account model representing a virtual or financial account (e.g. Kas, Bank, Dompet Digital).
class AccountModel {
  const AccountModel({
    required this.id,
    required this.name,
    required this.type,
    required this.currency,
    required this.balance,
    required this.openingBalance,
    this.isArchived = false,
  });

  final String id;
  final String name;
  final String type;
  final String currency;
  final int balance;
  final int openingBalance;
  final bool isArchived;

  factory AccountModel.fromJson(Map<String, dynamic> json) {
    return AccountModel(
      id: json['id'] as String? ?? '',
      name: json['name'] as String? ?? '',
      type: json['type'] as String? ?? 'cash',
      currency: json['currency'] as String? ?? 'IDR',
      balance: (json['balance'] as num?)?.toInt() ?? 0,
      openingBalance: (json['opening_balance'] as num?)?.toInt() ?? 0,
      isArchived: json['is_archived'] as bool? ?? false,
    );
  }

  Map<String, dynamic> toJson() {
    return <String, dynamic>{
      'id': id,
      'name': name,
      'type': type,
      'currency': currency,
      'balance': balance,
      'opening_balance': openingBalance,
      'is_archived': isArchived,
    };
  }
}
