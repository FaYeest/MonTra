/// API configuration for the MonTra mobile client.
///
/// The mobile client communicates exclusively with the authoritative Go backend (`/api/v1/*`)
/// and never directly contacts the internal Python ML service.
class MontraApiConfig {
  const MontraApiConfig({
    this.baseUrl = const String.fromEnvironment(
      'MONTRA_API_BASE_URL',
      defaultValue: 'http://10.0.2.2:8080/api/v1',
    ),
    this.requestTimeout = const Duration(seconds: 10),
  });

  final String baseUrl;
  final Duration requestTimeout;
}
