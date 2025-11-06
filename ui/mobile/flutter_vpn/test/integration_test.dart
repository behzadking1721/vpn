import 'package:flutter_test/flutter_test.dart';
import 'package:http/http.dart' as http;
import 'dart:convert';

void main() {
  const base = String.fromEnvironment('API_BASE_URL', defaultValue: 'http://localhost:8080');

  test('mock server status and stats', () async {
    final statusResp = await http.get(Uri.parse('$base/api/status'));
    expect(statusResp.statusCode, anyOf(200, 404));

    final statsResp = await http.get(Uri.parse('$base/api/stats'));
    expect(statsResp.statusCode, anyOf(200, 404));
    if (statsResp.statusCode == 200) {
      final obj = jsonDecode(statsResp.body);
      expect(obj, contains('data_sent'));
      expect(obj, contains('data_received'));
    }
  });
}
