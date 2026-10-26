/// User model representing an authenticated MonTra account.
class UserModel {
  const UserModel({
    required this.id,
    required this.email,
    required this.fullName,
    this.createdAt,
  });

  final String id;
  final String email;
  final String fullName;
  final DateTime? createdAt;

  factory UserModel.fromJson(Map<String, dynamic> json) {
    return UserModel(
      id: json['id'] as String? ?? '',
      email: json['email'] as String? ?? '',
      fullName: (json['display_name'] ?? json['full_name']) as String? ?? '',
      createdAt: json['created_at'] != null
          ? DateTime.tryParse(json['created_at'] as String)
          : null,
    );
  }

  Map<String, dynamic> toJson() {
    return <String, dynamic>{
      'id': id,
      'email': email,
      'full_name': fullName,
      if (createdAt != null) 'created_at': createdAt!.toIso8601String(),
    };
  }
}
