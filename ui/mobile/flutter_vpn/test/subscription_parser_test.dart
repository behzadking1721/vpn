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

  test('vmess: websocket full fields (path, headers, sni, aid)', () async {
    final vmessObj = {
      'v': '2',
      'ps': 'ws-full',
      'add': 'ws.example.com',
      'port': '443',
      'id': 'uuid-ws-1',
      'aid': 0,
      'net': 'ws',
      'path': '/websocket',
      'host': 'host.example.com',
      'tls': 'tls',
      'sni': 'sni.example.com',
      'headers': {
        'Host': 'host.example.com',
        'X-Custom': 'value'
      }
    };
    final vmess = 'vmess://' + base64Encode(utf8.encode(jsonEncode(vmessObj)));
    final parsed = await ServerService.parseSubscriptionLinkStatic(vmess);
    expect(parsed.length, 1);
    final s = parsed.first;
    expect(s.protocol, 'vmess');
    expect(s.network, 'ws');
    expect(s.wsPath, '/websocket');
    expect(s.wsHeaders, isNotNull);
    expect(s.wsHeaders!['Host'], 'host.example.com');
    expect(s.alterId, 0);
    expect(s.sni, 'sni.example.com');
    expect(s.tls, isTrue);
  });

  test('vmess: tcp network and kcp', () async {
    final tcp = {'add': 'tcp.example', 'port': 1194, 'id': 'tcp-id', 'net': 'tcp', 'ps': 'tcp-node'};
    final kcp = {'add': 'kcp.example', 'port': 29900, 'id': 'kcp-id', 'net': 'kcp', 'ps': 'kcp-node'};
    final p = 'vmess://' + base64Encode(utf8.encode(jsonEncode(tcp))) + '\n' + 'vmess://' + base64Encode(utf8.encode(jsonEncode(kcp)));
    final parsed = await ServerService.parseSubscriptionLinkStatic(base64Encode(utf8.encode(p)));
    expect(parsed.length, 2);
    expect(parsed[0].network, 'tcp');
    expect(parsed[1].network, 'kcp');
  });

  test('vless: ws path and sni in query params', () async {
    final vless = 'vless://abcd-ef01@vless-ws.example.org:443?security=tls&path=%2Fws&sni=vless-sni.example#vless-ws';
    final parsed = await ServerService.parseSubscriptionLinkStatic(vless);
    expect(parsed.length, 1);
    final s = parsed.first;
    expect(s.protocol, 'vless');
    expect(s.wsPath, '/ws');
    expect(s.sni, 'vless-sni.example');
    expect(s.port, 443);
  });

  test('vless: header query params mapped to wsHeaders', () async {
    final vless = 'vless://abcd-ef01@vless-headers.example.org:443?security=tls&path=%2Fws&header=Host%3Dvless-headers.example&header=X-Custom%3Dvalue#vless-headers';
    final parsed = await ServerService.parseSubscriptionLinkStatic(vless);
    expect(parsed.length, 1);
    final s = parsed.first;
    expect(s.protocol, 'vless');
    expect(s.wsPath, '/ws');
    expect(s.port, 443);
    expect(s.wsHeaders, isNotNull);
    expect(s.wsHeaders!['Host'], 'vless-headers.example');
    expect(s.wsHeaders!['X-Custom'], 'value');
  });
}

