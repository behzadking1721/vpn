import 'dart:convert';
import 'package:flutter_test/flutter_test.dart';
import 'package:vpn_client/services/server_service.dart';

void main() {
  test('parseSubscriptionLink handles vmess simple payload', () async {
      // a minimal vmess json encoded example (may need tuning)
      final vmessJson = '{"add":"example.com","port":"443","id":"uuid1","ps":"vm1","tls":"tls"}';
      final vmessEncoded = 'vmess://' + base64Encode(utf8.encode(vmessJson));
      final parsed = await ServerService.parseSubscriptionLinkStatic(vmessEncoded);
    expect(parsed, isNotNull);
    expect(parsed.length, greaterThan(0));
  });

  test('parseSubscriptionLink handles base64 subscription list', () async {
    // simulate a subscription that is a base64 list of vmess lines
    final vmess1 = 'vmess://' + base64Encode(utf8.encode('{"add":"a.example.com","port":"443","id":"u1","ps":"node-a","tls":"tls"}'));
    final vmess2 = 'vmess://' + base64Encode(utf8.encode('{"add":"b.example.com","port":"8443","id":"u2","ps":"node-b","tls":"tls"}'));
    final listPayload = base64Encode(utf8.encode('$vmess1\n$vmess2'));
    final parsed = await ServerService.parseSubscriptionLinkStatic(listPayload);
    expect(parsed.length, equals(2));
    expect(parsed[0].host, contains('a.example.com'));
    expect(parsed[1].port, 8443);
  });

    test('parseSubscriptionLink handles vless simple URI', () async {
      final vless = 'vless://uuid123@example.com:443?encrypt=none#MyNode';
      final parsed = await ServerService.parseSubscriptionLinkStatic(vless);
      expect(parsed, isNotNull);
      expect(parsed.length, greaterThan(0));
      expect(parsed.first.protocol, 'vless');
    });

    test('parseSubscriptionLink handles ss base64', () async {
      // base64 of "aes-256-cfb:password@shadowsocks.example.com:8388"
      final payload = base64Encode(utf8.encode('aes-256-cfb:password@shadowsocks.example.com:8388'));
      final ss = 'ss://$payload#test';
      final parsed = await ServerService.parseSubscriptionLinkStatic(ss);
      expect(parsed, isNotNull);
      expect(parsed.length, greaterThan(0));
      expect(parsed.first.protocol, 'ss');
      expect(parsed.first.port, 8388);
    });
}
