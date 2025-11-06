import 'package:flutter_test/flutter_test.dart';
import 'package:vpn_client/models/server.dart';

void main() {
  test('Server.fromJson handles address/uuid fields', () {
    final json1 = {
      'uuid': 'u1',
      'name': 'Test1',
      'address': 'example.com',
      'port': 443,
      'protocol': 'vless',
      'security': 'tls',
    };

    final s1 = Server.fromJson(json1);
    expect(s1.id, contains('u1'));
    expect(s1.host, 'example.com');
    expect(s1.port, 443);
    expect(s1.protocol, 'vless');
  });

  test('Server.fromJson parses numeric port strings', () {
    final json2 = {
      'id': 'x',
      'name': 'Test2',
      'host': '1.2.3.4',
      'port': '1194',
      'protocol': 'openvpn',
    };
    final s2 = Server.fromJson(json2);
    expect(s2.port, 1194);
  });
}
