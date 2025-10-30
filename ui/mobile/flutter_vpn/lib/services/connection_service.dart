import 'dart:async';
import 'package:vpn_client/models/server.dart';
import 'package:vpn_client/models/connection_status.dart';

/// سرویس اتصال برای مدیریت اتصال به سرورهای VPN
class ConnectionService {
  ConnectionStatus _currentStatus = ConnectionStatus(isConnected: false);
  Server? _currentServer;
  
  final StreamController<ConnectionStatus> _statusController = 
      StreamController<ConnectionStatus>.broadcast();
  
  Stream<ConnectionStatus> get statusStream => _statusController.stream;

  ConnectionStatus get currentStatus => _currentStatus;
  Server? get currentServer => _currentServer;

  /// اتصال به یک سرور
  Future<bool> connectToServer(Server server) async {
    // شبیه‌سازی فرآیند اتصال
    // در یک اپلیکیشن واقعی، اینجا با استفاده از FFI یا Platform Channels
    // با هسته اتصال (مانند sing-box یا xray-core) ارتباط برقرار می‌شود
    
    // شبیه‌سازی زمان اتصال
    await Future.delayed(const Duration(seconds: 2));
    
    // به صورت تصادفی موفقیت یا شکست اتصال را تعیین می‌کنیم
    final bool isConnected = DateTime.now().millisecond % 2 == 0;
    
    if (isConnected) {
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
    } else {
      _currentServer = null;
      _currentStatus = ConnectionStatus(isConnected: false);
    }
    
    // ارسال وضعیت به stream
    _statusController.add(_currentStatus);
    
    return isConnected;
  }

  /// قطع اتصال از سرور فعلی
  Future<bool> disconnect() async {
    // شبیه‌سازی فرآیند قطع اتصال
    await Future.delayed(const Duration(milliseconds: 500));
    
    _currentServer = null;
    _currentStatus = ConnectionStatus(isConnected: false);
    
    // ارسال وضعیت به stream
    _statusController.add(_currentStatus);
    
    return true;
  }

  /// دریافت آمار مصرف داده
  Future<Map<String, int>> getDataUsage() async {
    // شبیه‌سازی آمار مصرف داده
    if (_currentStatus.isConnected) {
      // شبیه‌سازی افزایش داده
      await Future.delayed(const Duration(milliseconds: 100));
      
      final int additionalSent = 1024 * (DateTime.now().millisecond % 100);
      final int additionalReceived = 2048 * (DateTime.now().millisecond % 100);
      
      return {
        'sent': _currentStatus.dataSent + additionalSent,
        'received': _currentStatus.dataReceived + additionalReceived,
      };
    }
    
    return {'sent': 0, 'received': 0};
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

  /// بستن سرویس و آزاد کردن منابع
  void dispose() {
    _statusController.close();
  }
}