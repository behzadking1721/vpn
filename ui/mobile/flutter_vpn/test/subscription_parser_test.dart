import 'dart:convert';
import 'package:flutter_test/flutter_test.dart';
import 'package:vpn_client/services/server_service.dart';

void main() {
  test('vmess: base64 JSON with tls and string port', () async {
    final vmessJson = '{"add":"vmess.example.com","port":"443","id":"vm-uuid-1","ps":"vm-1","tls":"tls"}';
    final vmess = 'vmess://' + base64Encode(utf8.encode(vmessJson));
    final parsed = await ServerService.parseSubscriptionLinkStatic(vmess);
    expect(parsed, isNotNull);
    expect(parsed.length, 1);
    final s = parsed.first;
    expect(s.protocol, 'vmess');
    expect(s.host, 'vmess.example.com');
    expect(s.port, 443);
    expect(s.tls, isTrue);
  });

  test('vmess: JSON with numeric port and host key', () async {
    final vmessJson = '{"host":"v2.example.net","port":8443,"uuid":"vm-uuid-2","name":"v2node"}';
    final vmess = 'vmess://' + base64Encode(utf8.encode(vmessJson));
    final parsed = await ServerService.parseSubscriptionLinkStatic(vmess);
    expect(parsed.length, 1);
    final s = parsed.first;
    expect(s.host, 'v2.example.net');
    expect(s.port, 8443);
    expect(s.name, isNotEmpty);
  });

  test('vless: uri without explicit port but with security=tls should default to 443', () async {
    final vless = 'vless://11112222-3333-4444-5555-666677778888@vless.example.com?security=tls#vless-node';
    final parsed = await ServerService.parseSubscriptionLinkStatic(vless);
    expect(parsed.length, 1);
    final s = parsed.first;
    expect(s.protocol, 'vless');
    expect(s.host, 'vless.example.com');
    expect(s.port, 443);
    expect(s.tls, isTrue);
    expect(s.name, 'vless-node');
  });

  test('vless: full URI with port, path and headers present in query', () async {
    final vless = 'vless://abcd-ef01@vless2.example.org:8443?path=%2Fws&header=host%3Dvless2.example.org#node2';
    final parsed = await ServerService.parseSubscriptionLinkStatic(vless);
    expect(parsed.length, 1);
    final s = parsed.first;
    expect(s.host, 'vless2.example.org');
    expect(s.port, 8443);
    expect(s.name, 'node2');
  });

  test('ss: raw form method:pass@host:port and base64 variant', () async {
    final raw = 'ss://aes-128-gcm:pw123@ss.example.org:1194#raw';
    final parsedRaw = await ServerService.parseSubscriptionLinkStatic(raw);
    expect(parsedRaw.length, 1);
    expect(parsedRaw.first.protocol, 'ss');
    expect(parsedRaw.first.port, 1194);

    final base = base64Encode(utf8.encode('aes-128-gcm:pw123@ss.example.org:1194'));
    final ssEncoded = 'ss://$base#encoded';
    final parsedEnc = await ServerService.parseSubscriptionLinkStatic(ssEncoded);
    expect(parsedEnc.length, 1);
    expect(parsedEnc.first.port, 1194);
  });

  test('multiline/base64 subscription list mixing vmess/vless/ss', () async {
    final vmess = 'vmess://' + base64Encode(utf8.encode('{"add":"multi-vm.example","port":"443","id":"m1","ps":"multi-vm"}'));
    final vless = 'vless://zzz@multi-vless.example:8443?security=none#mvless';
    final ssRaw = 'ss://aes-256-cfb:pw@multi-ss.example:8388#ss1';
    final combined = '$vmess\n$vless\n$ssRaw';
    final listBase = base64Encode(utf8.encode(combined));
    final parsed = await ServerService.parseSubscriptionLinkStatic(listBase);
    expect(parsed.length, 3);
    expect(parsed.map((s) => s.protocol), containsAll(['vmess', 'vless', 'ss']));
  });
}

