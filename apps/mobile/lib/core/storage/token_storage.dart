import 'dart:convert';
import 'package:shared_preferences/shared_preferences.dart';

import '../models/user.dart';

/// Manages persistence of the JWT authentication token and cached user profile.
class TokenStorage {
  TokenStorage(this._prefs);

  static const String _keyToken = 'montra_auth_token';
  static const String _keyUser = 'montra_auth_user';

  final SharedPreferences _prefs;

  static Future<TokenStorage> create() async {
    final prefs = await SharedPreferences.getInstance();
    return TokenStorage(prefs);
  }

  String? getToken() {
    return _prefs.getString(_keyToken);
  }

  Future<void> saveToken(String token) async {
    await _prefs.setString(_keyToken, token);
  }

  Future<void> removeToken() async {
    await _prefs.remove(_keyToken);
  }

  UserModel? getUser() {
    final raw = _prefs.getString(_keyUser);
    if (raw == null || raw.isEmpty) return null;
    try {
      final json = jsonDecode(raw) as Map<String, dynamic>;
      return UserModel.fromJson(json);
    } catch (_) {
      return null;
    }
  }

  Future<void> saveUser(UserModel user) async {
    await _prefs.setString(_keyUser, jsonEncode(user.toJson()));
  }

  Future<void> removeUser() async {
    await _prefs.remove(_keyUser);
  }

  Future<void> clearAll() async {
    await _prefs.remove(_keyToken);
    await _prefs.remove(_keyUser);
  }
}
