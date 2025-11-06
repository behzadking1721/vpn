import 'dart:async';
import 'package:flutter/material.dart';
import 'package:vpn_client/models/server.dart';
import 'package:vpn_client/models/connection_status.dart';
import 'package:vpn_client/services/connection_service.dart';
import 'package:vpn_client/services/server_service.dart';

/// Provider برای مدیریت وضعیت کلی اپلیکیشن
class AppProvider with ChangeNotifier {
  final ConnectionService _connectionService = ConnectionService();
  final ServerService _serverService = ServerService();
  
  Server? _selectedServer;
  bool _isConnecting = false;
  bool _isScanningQR = false;
  
  // Getters
  ConnectionStatus get connectionStatus => _connectionService.currentStatus;
  Server? get currentServer => _connectionService.currentServer;
  Server? get selectedServer => _selectedServer;
  bool get isConnecting => _isConnecting;
  bool get isScanningQR => _isScanningQR;
  List<Server> get allServers => _serverService.getAllServers();
  
  // Streams
  Stream<ConnectionStatus> get connectionStatusStream => _connectionService.statusStream;
  Stream<List<Server>> get serversStream => _serverService.serversStream;
  
  AppProvider() {
    // شروع به‌روزرسانی آمار مصرف داده
    _startDataUsageUpdater();
  }
  
  /// شروع به‌روزرسانی آمار مصرف داده
  void _startDataUsageUpdater() {
    Timer.periodic(const Duration(seconds: 1), (timer) {
      if (connectionStatus.isConnected) {
        _connectionService.updateDataUsage();
      }
    });
  }
  
  /// انتخاب سرور
  void selectServer(Server server) {
    _selectedServer = server;
    notifyListeners();
  }
  
  /// اتصال به سرور انتخاب شده
  Future<bool> connectToSelectedServer() async {
    if (_selectedServer == null) return false;
    
    _isConnecting = true;
    notifyListeners();
    
    final result = await _connectionService.connectToServer(_selectedServer!);
    
    _isConnecting = false;
    notifyListeners();
    
    return result;
  }
  
  /// اتصال به سرور با ID
  Future<bool> connectToServerById(String serverId) async {
    final server = _serverService.getServerById(serverId);
    if (server == null) return false;
    
    selectServer(server);
    return await connectToSelectedServer();
  }
  
  /// قطع اتصال
  Future<bool> disconnect() async {
    final result = await _connectionService.disconnect();
    notifyListeners();
    return result;
  }
  
  /// افزودن سرور جدید
  Future<void> addServer(Server server) async {
    await _serverService.addServer(server);
    notifyListeners();
  }
  
  /// به‌روزرسانی سرور
  Future<void> updateServer(Server server) async {
    await _serverService.updateServer(server);
    notifyListeners();
  }
  
  /// حذف سرور
  Future<void> removeServer(String serverId) async {
    await _serverService.removeServer(serverId);
    if (_selectedServer?.id == serverId) {
      _selectedServer = null;
    }
    notifyListeners();
  }
  
  /// پیدا کردن سرور با کمترین پینگ
  Server? findFastestServer() {
    return _serverService.findFastestServer();
  }
  
  /// اتصال به سریع‌ترین سرور
  Future<bool> connectToFastestServer() async {
    final fastestServer = findFastestServer();
    if (fastestServer == null) return false;
    
    selectServer(fastestServer);
    return await connectToSelectedServer();
  }
  
  /// شبیه‌سازی پینگ سرورها
  Future<void> pingServers() async {
    await _serverService.pingServers();
    notifyListeners();
  }
  
  /// پارس کردن لینک اشتراک
  Future<List<Server>> parseSubscriptionLink(String link) async {
    return await _serverService.parseSubscriptionLink(link);
  }
  
  /// شروع اسکن QR
  void startQRScan() {
    _isScanningQR = true;
    notifyListeners();
  }
  
  /// پایان اسکن QR
  void endQRScan() {
    _isScanningQR = false;
    notifyListeners();
  }
  
  /// آزاد کردن منابع
  void dispose() {
    _connectionService.dispose();
    _serverService.dispose();
    super.dispose();
  }
}