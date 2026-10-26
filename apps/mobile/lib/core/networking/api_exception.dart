/// Represents an error returned by the MonTra API or mobile network transport.
class ApiException implements Exception {
  const ApiException({
    required this.message,
    this.code = 'UNKNOWN_ERROR',
    this.statusCode = 500,
    this.details,
  });

  final String message;
  final String code;
  final int statusCode;
  final dynamic details;

  bool get isUnauthorized => statusCode == 401;
  bool get isNotFound => statusCode == 404;
  bool get isConflict => statusCode == 409;
  bool get isValidationError => statusCode == 422 || statusCode == 400;

  @override
  String toString() => 'ApiException($code, $statusCode): $message';
}
