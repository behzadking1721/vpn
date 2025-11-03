import 'dart:async';
import 'dart:convert';
import 'package:http/http.dart' as http;
import 'package:vpn_client/models/server.dart';
import 'package:vpn_client/models/connection_status.dart';

/// سرویس اتصال برای مدیریت اتصال به سرورهای VPN
class ConnectionService {
  ConnectionStatus _currentStatus = ConnectionStatus(isConnected: false);
  Server? _currentServer;
  final String _apiBaseUrl = 'http://localhost:8080/api'; // آدرس پایه API
  
  final StreamController<ConnectionStatus> _statusController = 
      StreamController<ConnectionStatus>.broadcast();
  
  Stream<ConnectionStatus> get statusStream => _statusController.stream;

  ConnectionStatus get currentStatus => _currentStatus;
  Server? get currentServer => _currentServer;

  /// اتصال به یک سرور
  Future<bool> connectToServer(Server server) async {
    try {
      // ارسال درخواست اتصال به API
      final response = await http.post(
        Uri.parse('$_apiBaseUrl/connect'),
        headers: {'Content-Type': 'application/json'},
        body: jsonEncode({'server_id': server.id}),
      );
      
      if (response.statusCode == 200) {
        _currentServer = server;
        _currentStatus = ConnectionStatus(
          isConnected: true,
          serverId: server.id,
          serverName: server.name,
          protocol: server.protocol,
          startTime: DateTime.now(),
          dataSent: 0,
          dataReceived: 0,
        );
        
        // ارسال وضعیت به stream
        _statusController.add(_currentStatus);
        return true;
      } else {
        // خطا در اتصال
        _currentServer = null;
        _currentStatus = ConnectionStatus(isConnected: false);
        _statusController.add(_currentStatus);
        return false;
      }
    } catch (e) {
      // خطا در ارتباط با API
      _currentServer = null;
      _currentStatus = ConnectionStatus(isConnected: false);
      _statusController.add(_currentStatus);
      return false;
    }
  }

  /// قطع اتصال از سرور فعلی
  Future<bool> disconnect() async {
    try {
      // ارسال درخواست قطع اتصال به API
      final response = await http.post(
        Uri.parse('$_apiBaseUrl/disconnect'),
        headers: {'Content-Type': 'application/json'},
      );
      
      if (response.statusCode == 200) {
        _currentServer = null;
        _currentStatus = ConnectionStatus(isConnected: false);
        
        // ارسال وضعیت به stream
        _statusController.add(_currentStatus);
        return true;
      } else {
        return false;
      }
    } catch (e) {
      return false;
    }
  }

  /// دریافت آمار مصرف داده
  Future<Map<String, int>> getDataUsage() async {
    try {
      // دریافت آمار از API
      final response = await http.get(
        Uri.parse('$_apiBaseUrl/stats'),
        headers: {'Content-Type': 'application/json'},
      );
      
      if (response.statusCode == 200) {
        final data = jsonDecode(response.body);
        return {
          'sent': data['data_sent'] as int? ?? 0,
          'received': data['data_received'] as int? ?? 0,
        };
      } else {
        return {'sent': 0, 'received': 0};
      }
    } catch (e) {
      return {'sent': 0, 'received': 0};
    }
  }

  /// بروزرسانی آمار مصرف داده
  Future<void> updateDataUsage() async {
    if (_currentStatus.isConnected) {
      final usage = await getDataUsage();
      _currentStatus = _currentStatus.copyWith(
        dataSent: usage['sent']!,
        dataReceived: usage['received']!,
      );
      
      _statusController.add(_currentStatus);
    }
  }

  /// دریافت وضعیت اتصال از API
  Future<ConnectionStatus> getConnectionStatus() async {
    try {
      final response = await http.get(
        Uri.parse('$_apiBaseUrl/status'),
        headers: {'Content-Type': 'application/json'},
      );
      
      if (response.statusCode == 200) {
        final data = jsonDecode(response.body);
        final status = data['status'] as String;
        
        return ConnectionStatus(
          isConnected: status == 'connected',
          serverId: _currentServer?.id,
          serverName: _currentServer?.name,
          protocol: _currentServer?.protocol,
        );
      } else {
        return ConnectionStatus(isConnected: false);
      }
    } catch (e) {
      return ConnectionStatus(isConnected: false);
    }
  }

  /// بستن سرویس و آزاد کردن منابع
  void dispose() {
    _statusController.close();
  }
}