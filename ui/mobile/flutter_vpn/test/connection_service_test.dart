import 'package:flutter_test/flutter_test.dart';
import 'package:vpn_client/services/connection_service.dart';
import 'package:vpn_client/models/server.dart';

void main() {
  test('ConnectionService connect/status/stats/disconnect flow', () async {
    // Use the mock server running on localhost:8080
    final cs = ConnectionService();

    // Construct a server object matching mock config
    final server = Server(
      id: 'xxxxxxxx-xxxx-xxxx-xxxx-xxxxxxxxxxxx',
      name: 'سرور تست',
      host: 'example.com',
      port: 443,
      protocol: 'vless',
    );

    // Attempt to connect (mock server should respond OK)
    final ok = await cs.connectToServer(server);
    expect(ok, isTrue, reason: 'connectToServer should return true for mock server');

    // Status should report connected
    final status = await cs.getConnectionStatus();
    expect(status.isConnected, isTrue);

    // Update and read data usage
    await cs.updateDataUsage();
    final usage = await cs.getDataUsage();
    expect(usage['sent'] is int, isTrue);
    expect(usage['received'] is int, isTrue);

    // Disconnect
    final d = await cs.disconnect();
    expect(d, isTrue);
  }, timeout: Timeout(Duration(seconds: 10)));
}
