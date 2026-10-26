/// Transfer model representing an atomic transfer between two user accounts.
class TransferModel {
  const TransferModel({
    required this.id,
    required this.sourceAccountId,
    required this.destinationAccountId,
    required this.amount,
    this.feeAmount = 0,
    required this.currency,
    required this.transferDate,
    this.notes,
    this.sourceAccountName,
    this.destinationAccountName,
  });

  final String id;
  final String sourceAccountId;
  final String destinationAccountId;
  final int amount;
  final int feeAmount;
  final String currency;
  final DateTime transferDate;
  final String? notes;
  final String? sourceAccountName;
  final String? destinationAccountName;

  factory TransferModel.fromJson(Map<String, dynamic> json) {
    return TransferModel(
      id: json['id'] as String? ?? '',
      sourceAccountId: json['source_account_id'] as String? ?? '',
      destinationAccountId: json['destination_account_id'] as String? ?? '',
      amount: (json['amount'] as num?)?.toInt() ?? 0,
      feeAmount: (json['fee_amount'] as num?)?.toInt() ?? 0,
      currency: json['currency'] as String? ?? 'IDR',
      transferDate: json['transfer_date'] != null
          ? DateTime.tryParse(json['transfer_date'] as String) ?? DateTime.now()
          : DateTime.now(),
      notes: json['notes'] as String?,
      sourceAccountName: json['source_account_name'] as String?,
      destinationAccountName: json['destination_account_name'] as String?,
    );
  }

  Map<String, dynamic> toJson() {
    return <String, dynamic>{
      'id': id,
      'source_account_id': sourceAccountId,
      'destination_account_id': destinationAccountId,
      'amount': amount,
      'fee_amount': feeAmount,
      'currency': currency,
      'transfer_date': transferDate.toIso8601String(),
      if (notes != null) 'notes': notes,
    };
  }
}
