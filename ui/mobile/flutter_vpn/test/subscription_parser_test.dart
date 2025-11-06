import 'dart:convert';
import 'package:flutter_test/flutter_test.dart';
import 'package:vpn_client/services/server_service.dart';

void main() {
  test('parseSubscriptionLink handles vmess simple payload', () async {
    final service = ServerService();
    // a minimal vmess json encoded example (may need tuning)
  final vmessJson = '{"add":"example.com","port":"443","id":"uuid1","ps":"vm1","tls":"tls"}';
  final vmessEncoded = 'vmess://' + base64Encode(utf8.encode(vmessJson));
    final parsed = await service.parseSubscriptionLink(vmessEncoded);
    expect(parsed, isNotNull);
    expect(parsed.length, greaterThan(0));
  });
}
