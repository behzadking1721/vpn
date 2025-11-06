// Configuration constants for the app.
// The API base URL can be overridden at build/run time using
// `--dart-define=API_BASE_URL=http://10.0.2.2:8080/api` for Android emulator.
const String apiBaseUrl = String.fromEnvironment(
  'API_BASE_URL',
  defaultValue: 'http://localhost:8080/api',
);
