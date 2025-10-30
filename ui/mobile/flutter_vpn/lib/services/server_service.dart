import 'dart:async';
import 'dart:convert';
import 'package:http/http.dart' as http;
import 'package:vpn_client/models/server.dart';

/// سرویس سرور برای مدیریت سرورهای VPN
class ServerService {
  final String _apiBaseUrl = 'http://localhost:8080/api'; // آدرس پایه API
  List<Server> _servers = [];
  
  final StreamController<List<Server>> _serversController = 
      StreamController<List<Server>>.broadcast();
  
  Stream<List<Server>> get serversStream => _serversController.stream;

  ServerService() {
    // بارگذاری سرورها از API
    loadServers();
  }

  /// بارگذاری سرورها از API
  Future<void> loadServers() async {
    try {
      final response = await http.get(
        Uri.parse('$_apiBaseUrl/servers'),
        headers: {'Content-Type': 'application/json'},
      );
      
      if (response.statusCode == 200) {
        final List<dynamic> data = jsonDecode(response.body);
        _servers = data.map((json) => Server.fromJson(json)).toList();
        _serversController.add(_servers);
      }
    } catch (e) {
      // در صورت خطا، لیست خالی ارسال می‌شود
      _servers = [];
      _serversController.add(_servers);
    }
  }

  /// دریافت تمام سرورها
  List<Server> getAllServers() {
    return List.unmodifiable(_servers);
  }

  /// افزودن سرور جدید
  Future<void> addServer(Server server) async {
    try {
      final response = await http.post(
        Uri.parse('$_apiBaseUrl/servers'),
        headers: {'Content-Type': 'application/json'},
        body: jsonEncode(server.toJson()),
      );
      
      if (response.statusCode == 201) {
        final newServer = Server.fromJson(jsonDecode(response.body));
        _servers.add(newServer);
        _serversController.add(_servers);
      }
    } catch (e) {
      // خطا در افزودن سرور
      rethrow;
    }
  }

  /// به‌روزرسانی سرور
  Future<void> updateServer(Server server) async {
    try {
      final response = await http.put(
        Uri.parse('$_apiBaseUrl/servers/${server.id}'),
        headers: {'Content-Type': 'application/json'},
        body: jsonEncode(server.toJson()),
      );
      
      if (response.statusCode == 200) {
        final updatedServer = Server.fromJson(jsonDecode(response.body));
        final index = _servers.indexWhere((s) => s.id == server.id);
        if (index != -1) {
          _servers[index] = updatedServer;
          _serversController.add(_servers);
        }
      }
    } catch (e) {
      // خطا در به‌روزرسانی سرور
      rethrow;
    }
  }

  /// حذف سرور
  Future<void> removeServer(String serverId) async {
    try {
      final response = await http.delete(
        Uri.parse('$_apiBaseUrl/servers/$serverId'),
        headers: {'Content-Type': 'application/json'},
      );
      
      if (response.statusCode == 200) {
        _servers.removeWhere((server) => server.id == serverId);
        _serversController.add(_servers);
      }
    } catch (e) {
      // خطا در حذف سرور
      rethrow;
    }
  }

  /// دریافت سرور با ID
  Server? getServerById(String id) {
    try {
      return _servers.firstWhere((server) => server.id == id);
    } catch (e) {
      return null;
    }
  }

  /// پیدا کردن سریع‌ترین سرور
  Server? findFastestServer() {
    if (_servers.isEmpty) return null;
    
    // فیلتر کردن سرورهای فعال
    final activeServers = _servers.where((server) => server.enabled).toList();
    if (activeServers.isEmpty) return null;
    
    // مرتب‌سازی بر اساس پینگ (کمترین پینگ = سریع‌ترین)
    activeServers.sort((a, b) => a.ping.compareTo(b.ping));
    
    return activeServers.first;
  }

  /// شبیه‌سازی پینگ سرورها
  Future<void> pingServers() async {
    // در یک اپلیکیشن واقعی، اینجا پینگ واقعی سرورها انجام می‌شود
    for (var server in _servers) {
      // شبیه‌سازی پینگ
      final newPing = (10 + (DateTime.now().millisecond % 200));
      final updatedServer = server.copyWith(ping: newPing);
      
      final index = _servers.indexWhere((s) => s.id == server.id);
      if (index != -1) {
        _servers[index] = updatedServer;
      }
    }
    
    _serversController.add(_servers);
  }

  /// پارس کردن لینک اشتراک
  Future<List<Server>> parseSubscriptionLink(String link) async {
    // در یک اپلیکیشن واقعی، اینجا لینک اشتراک پارس می‌شود
    // برای حالا یک لیست خالی برمی‌گردانیم
    return [];
  }

  /// بستن سرویس و آزاد کردن منابع
  void dispose() {
    _serversController.close();
  }
}