/// Category model representing income/expense categories (e.g. Makanan, Gaji, Transportasi).
class CategoryModel {
  const CategoryModel({
    required this.id,
    required this.name,
    required this.type,
    this.icon = 'category',
    this.color = '#1D70B8',
    this.isSystem = false,
  });

  final String id;
  final String name;
  final String type; // 'income' or 'expense'
  final String icon;
  final String color;
  final bool isSystem;

  factory CategoryModel.fromJson(Map<String, dynamic> json) {
    return CategoryModel(
      id: json['id'] as String? ?? '',
      name: json['name'] as String? ?? '',
      type: json['type'] as String? ?? 'expense',
      icon: json['icon'] as String? ?? 'category',
      color: json['color'] as String? ?? '#1D70B8',
      isSystem: json['is_system'] as bool? ?? false,
    );
  }

  Map<String, dynamic> toJson() {
    return <String, dynamic>{
      'id': id,
      'name': name,
      'type': type,
      'icon': icon,
      'color': color,
      'is_system': isSystem,
    };
  }
}
