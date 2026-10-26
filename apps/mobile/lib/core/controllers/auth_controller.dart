import 'package:flutter/foundation.dart';

import '../models/user.dart';
import '../networking/api_client.dart';
import '../networking/api_exception.dart';
import '../storage/token_storage.dart';

enum AuthStatus {
  initial,
  authenticated,
  unauthenticated,
  loading,
}

/// Controller for managing user authentication state, sign-in, registration, and logout.
class AuthController extends ChangeNotifier {
  AuthController({
    required ApiClient apiClient,
    required TokenStorage tokenStorage,
  })  : _apiClient = apiClient,
        _tokenStorage = tokenStorage;

  final ApiClient _apiClient;
  final TokenStorage _tokenStorage;

  AuthStatus _status = AuthStatus.initial;
  UserModel? _currentUser;
  String? _errorMessage;

  AuthStatus get status => _status;
  UserModel? get currentUser => _currentUser;
  String? get errorMessage => _errorMessage;
  bool get isAuthenticated => _status == AuthStatus.authenticated;
  bool get isLoading => _status == AuthStatus.loading;

  /// Check initial auth state on app launch.
  Future<void> checkAuth() async {
    final token = _tokenStorage.getToken();
    if (token == null || token.isEmpty) {
      _status = AuthStatus.unauthenticated;
      _currentUser = null;
      notifyListeners();
      return;
    }

    _currentUser = _tokenStorage.getUser();
    _status = AuthStatus.loading;
    notifyListeners();

    try {
      dynamic res;
      try {
        res = await _apiClient.get('/auth/me');
      } catch (_) {
        res = await _apiClient.get('/me');
      }
      if (res is Map<String, dynamic>) {
        final userMap = res.containsKey('user') ? res['user'] : res;
        _currentUser = UserModel.fromJson(userMap as Map<String, dynamic>);
        await _tokenStorage.saveUser(_currentUser!);
        _status = AuthStatus.authenticated;
      } else {
        _status = AuthStatus.authenticated;
      }
    } on ApiException catch (e) {
      if (e.isUnauthorized) {
        await _tokenStorage.clearAll();
        _currentUser = null;
        _status = AuthStatus.unauthenticated;
      } else {
        // Network issue or offline - if we have a saved token and user, remain authenticated
        _status = _currentUser != null
            ? AuthStatus.authenticated
            : AuthStatus.unauthenticated;
      }
    } catch (_) {
      _status = _currentUser != null
          ? AuthStatus.authenticated
          : AuthStatus.unauthenticated;
    }

    notifyListeners();
  }

  /// Authenticate with email and password.
  Future<bool> login(String email, String password) async {
    _status = AuthStatus.loading;
    _errorMessage = null;
    notifyListeners();

    try {
      final res = await _apiClient.post(
        '/auth/login',
        body: <String, dynamic>{
          'email': email.trim(),
          'password': password,
        },
      );

      if (res is Map<String, dynamic>) {
        final token = (res['token'] ?? res['access_token']) as String?;
        final userMap = res['user'] as Map<String, dynamic>?;

        if (token != null && token.isNotEmpty) {
          await _tokenStorage.saveToken(token);
        }
        if (userMap != null) {
          _currentUser = UserModel.fromJson(userMap);
          await _tokenStorage.saveUser(_currentUser!);
        }

        _status = AuthStatus.authenticated;
        _errorMessage = null;
        notifyListeners();
        return true;
      }

      throw const ApiException(
        message: 'Respon server tidak valid saat login.',
        code: 'INVALID_RESPONSE',
      );
    } on ApiException catch (e) {
      _status = AuthStatus.unauthenticated;
      _errorMessage = e.message;
      notifyListeners();
      return false;
    } catch (e) {
      _status = AuthStatus.unauthenticated;
      _errorMessage = 'Gagal masuk: ${e.toString()}';
      notifyListeners();
      return false;
    }
  }

  /// Register a new user account.
  Future<bool> register({
    required String fullName,
    required String email,
    required String password,
  }) async {
    _status = AuthStatus.loading;
    _errorMessage = null;
    notifyListeners();

    try {
      final res = await _apiClient.post(
        '/auth/register',
        body: <String, dynamic>{
          'display_name': fullName.trim(),
          'full_name': fullName.trim(),
          'email': email.trim(),
          'password': password,
        },
      );

      if (res is Map<String, dynamic>) {
        final token = (res['token'] ?? res['access_token']) as String?;
        final userMap = res['user'] as Map<String, dynamic>?;

        if (token != null && token.isNotEmpty) {
          await _tokenStorage.saveToken(token);
        }
        if (userMap != null) {
          _currentUser = UserModel.fromJson(userMap);
          await _tokenStorage.saveUser(_currentUser!);
        }

        _status = AuthStatus.authenticated;
        _errorMessage = null;
        notifyListeners();
        return true;
      }

      throw const ApiException(
        message: 'Respon server tidak valid saat pendaftaran.',
        code: 'INVALID_RESPONSE',
      );
    } on ApiException catch (e) {
      _status = AuthStatus.unauthenticated;
      _errorMessage = e.message;
      notifyListeners();
      return false;
    } catch (e) {
      _status = AuthStatus.unauthenticated;
      _errorMessage = 'Gagal mendaftar: ${e.toString()}';
      notifyListeners();
      return false;
    }
  }

  /// Log out and clear stored session.
  Future<void> logout() async {
    try {
      await _apiClient.post('/auth/logout');
    } catch (_) {
      // Best-effort logout to server
    }
    await _tokenStorage.clearAll();
    _currentUser = null;
    _errorMessage = null;
    _status = AuthStatus.unauthenticated;
    notifyListeners();
  }

  void clearError() {
    _errorMessage = null;
    notifyListeners();
  }
}
