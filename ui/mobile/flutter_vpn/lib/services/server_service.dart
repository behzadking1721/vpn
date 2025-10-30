import 'dart:async';
import 'dart:math';
import 'package:vpn_client/models/server.dart';

/// سرویس مدیریت سرورها
class ServerService {
  final List<Server> _servers = [];
  final StreamController<List<Server>> _serversController = 
      StreamController<List<Server>>.broadcast();

  Stream<List<Server>> get serversStream => _serversController.stream;

  ServerService() {
    // بارگذاری سرورهای نمونه
    _loadSampleServers();
  }

  /// بارگذاری سرورهای نمونه
  void _loadSampleServers() {
    _servers.addAll([
      Server(
        id: '1',
        name: '🇯🇵 Japan Server',
        host: 'jp.example.com',
        port: 443,
        protocol: 'VMess',
        encryption: 'auto',
        password: 'sample-uuid',
        tls: true,
        remark: 'Tokyo Server',
        enabled: true,
        ping: 22,
      ),
      Server(
        id: '2',
        name: '🇺🇸 USA Server',
        host: 'us.example.com',
        port: 8388,
        protocol: 'Shadowsocks',
        method: 'aes-256-gcm',
        password: 'sample-password',
        remark: 'New York Server',
        enabled: true,
        ping: 45,
      ),
      Server(
        id: '3',
        name: '🇬🇧 UK Server',
        host: 'uk.example.com',
        port: 443,
        protocol: 'Trojan',
        password: 'sample-password',
        tls: true,
        remark: 'London Server',
        enabled: true,
        ping: 68,
      ),
      Server(
        id: '4',
        name: '🇩🇪 Germany Server',
        host: 'de.example.com',
        port: 443,
        protocol: 'VLESS',
        password: 'sample-uuid',
        tls: true,
        remark: 'Frankfurt Server',
        enabled: true,
        ping: 35,
      ),
    ]);
    
    _serversController.add(List.unmodifiable(_servers));
  }

  /// دریافت تمام سرورها
  List<Server> getAllServers() {
    return List.unmodifiable(_servers);
  }

  /// افزودن سرور جدید
  Future<void> addServer(Server server) async {
    _servers.add(server);
    _serversController.add(List.unmodifiable(_servers));
  }

  /// به‌روزرسانی سرور
  Future<void> updateServer(Server server) async {
    final index = _servers.indexWhere((s) => s.id == server.id);
    if (index != -1) {
      _servers[index] = server;
      _serversController.add(List.unmodifiable(_servers));
    }
  }

  /// حذف سرور
  Future<void> removeServer(String serverId) async {
    _servers.removeWhere((server) => server.id == serverId);
    _serversController.add(List.unmodifiable(_servers));
  }

  /// دریافت سرور با ID
  Server? getServerById(String id) {
    try {
      return _servers.firstWhere((server) => server.id == id);
    } catch (e) {
      return null;
    }
  }

  /// پیدا کردن سرور با کمترین پینگ
  Server? findFastestServer() {
    final enabledServers = _servers.where((server) => server.enabled).toList();
    if (enabledServers.isEmpty) return null;
    
    enabledServers.sort((a, b) => a.ping.compareTo(b.ping));
    return enabledServers.first;
  }

  /// شبیه‌سازی پینگ سرورها
  Future<void> pingServers() async {
    for (var server in _servers) {
      if (server.enabled) {
        // شبیه‌سازی پینگ
        await Future.delayed(const Duration(milliseconds: 100));
        
        // تولید یک مقدار پینگ تصادفی بین 10 تا 200
        final newPing = Random().nextInt(190) + 10;
        
        final updatedServer = server.copyWith(ping: newPing);
        final index = _servers.indexWhere((s) => s.id == server.id);
        if (index != -1) {
          _servers[index] = updatedServer;
        }
      }
    }
    
    _serversController.add(List.unmodifiable(_servers));
  }

  /// پارس کردن لینک اشتراک
  Future<List<Server>> parseSubscriptionLink(String link) async {
    // شبیه‌سازی پردازش لینک اشتراک
    await Future.delayed(const Duration(seconds: 1));
    
    // در یک اپلیکیشن واقعی، اینجا لینک اشتراک پردازش می‌شود
    // و سرورهای جدید استخراج می‌شوند
    
    return [
      Server(
        id: DateTime.now().millisecondsSinceEpoch.toString(),
        name: ' Imported Server',
        host: 'imported.example.com',
        port: 443,
        protocol: 'VMess',
        encryption: 'auto',
        password: 'imported-uuid',
        tls: true,
        remark: 'Imported from subscription',
        enabled: true,
        ping: 0,
      ),
    ];
  }

  /// بستن سرویس و آزاد کردن منابع
  void dispose() {
    _serversController.close();
  }
}