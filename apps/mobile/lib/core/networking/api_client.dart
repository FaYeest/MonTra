import 'dart:async';
import 'dart:convert';
import 'dart:io';
import 'package:http/http.dart' as http;

import '../storage/token_storage.dart';
import 'api_config.dart';
import 'api_exception.dart';

/// HTTP Client for communicating with the authoritative Go API backend (`/api/v1/*`).
class ApiClient {
  ApiClient({
    MontraApiConfig? config,
    TokenStorage? tokenStorage,
    http.Client? httpClient,
  })  : _config = config ?? const MontraApiConfig(),
        _tokenStorage = tokenStorage,
        _httpClient = httpClient ?? http.Client();

  final MontraApiConfig _config;
  final TokenStorage? _tokenStorage;
  final http.Client _httpClient;

  Future<dynamic> get(String path, {Map<String, String>? queryParams}) async {
    return _send('GET', path, queryParams: queryParams);
  }

  Future<dynamic> post(
    String path, {
    Map<String, dynamic>? body,
    Map<String, String>? extraHeaders,
  }) async {
    return _send('POST', path, body: body, extraHeaders: extraHeaders);
  }

  Future<dynamic> put(
    String path, {
    Map<String, dynamic>? body,
    Map<String, String>? extraHeaders,
  }) async {
    return _send('PUT', path, body: body, extraHeaders: extraHeaders);
  }

  Future<dynamic> delete(String path) async {
    return _send('DELETE', path);
  }

  Future<dynamic> _send(
    String method,
    String path, {
    Map<String, String>? queryParams,
    Map<String, dynamic>? body,
    Map<String, String>? extraHeaders,
  }) async {
    final uri = _buildUri(path, queryParams);
    final headers = <String, String>{
      'Content-Type': 'application/json',
      'Accept': 'application/json',
      if (extraHeaders != null) ...extraHeaders,
    };

    final token = _tokenStorage?.getToken();
    if (token != null && token.isNotEmpty) {
      headers['Authorization'] = 'Bearer $token';
    }

    http.Response response;
    try {
      final request = http.Request(method, uri);
      request.headers.addAll(headers);
      if (body != null) {
        request.body = jsonEncode(body);
      }

      final streamedResponse = await _httpClient
          .send(request)
          .timeout(_config.requestTimeout);
      response = await http.Response.fromStream(streamedResponse);
    } on SocketException {
      throw const ApiException(
        message: 'Tidak dapat terhubung ke server. Periksa koneksi internet Anda.',
        code: 'NETWORK_ERROR',
        statusCode: 0,
      );
    } on TimeoutException {
      throw const ApiException(
        message: 'Permintaan waktu habis. Silakan coba kembali sesaat lagi.',
        code: 'TIMEOUT_ERROR',
        statusCode: 408,
      );
    } catch (e) {
      throw ApiException(
        message: 'Terjadi kesalahan jaringan: ${e.toString()}',
        code: 'CLIENT_ERROR',
        statusCode: 0,
      );
    }

    return _handleResponse(response);
  }

  Uri _buildUri(String path, Map<String, String>? queryParams) {
    var base = _config.baseUrl;
    if (base.endsWith('/')) {
      base = base.substring(0, base.length - 1);
    }
    var cleanPath = path;
    if (!cleanPath.startsWith('/')) {
      cleanPath = '/$cleanPath';
    }

    final fullUrl = '$base$cleanPath';
    final parsed = Uri.parse(fullUrl);
    if (queryParams != null && queryParams.isNotEmpty) {
      return parsed.replace(
        queryParameters: <String, String>{
          ...parsed.queryParameters,
          ...queryParams,
        },
      );
    }
    return parsed;
  }

  dynamic _handleResponse(http.Response response) {
    Map<String, dynamic>? json;
    if (response.body.isNotEmpty) {
      try {
        json = jsonDecode(response.body) as Map<String, dynamic>;
      } catch (_) {
        // Not a JSON response
      }
    }

    if (response.statusCode >= 200 && response.statusCode < 300) {
      if (json != null && json.containsKey('data')) {
        return json['data'];
      }
      return json;
    }

    if (json != null && json.containsKey('error')) {
      final errorObj = json['error'];
      if (errorObj is Map<String, dynamic>) {
        final code = errorObj['code'] as String? ?? 'ERROR';
        final message = errorObj['message'] as String? ?? 'Terjadi kesalahan.';
        final details = errorObj['details'];
        throw ApiException(
          message: message,
          code: code,
          statusCode: response.statusCode,
          details: details,
        );
      }
    }

    throw ApiException(
      message: 'Permintaan gagal dengan kode HTTP ${response.statusCode}.',
      code: 'HTTP_ERROR_${response.statusCode}',
      statusCode: response.statusCode,
    );
  }
}
