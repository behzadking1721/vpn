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
    // Basic parser for common subscription formats.
    // - If link is an http(s) URL, fetch its contents and parse each line.
    // - If link is base64 or contains vmess:// or vless:// entries, parse them.
    final List<Server> parsed = [];
    try {
      if (link.startsWith('http://') || link.startsWith('https://')) {
        final resp = await http.get(Uri.parse(link));
        if (resp.statusCode == 200) {
          final body = resp.body.trim();
          // many subscription links are base64-encoded lists
          try {
            final decoded = utf8.decode(base64.decode(body));
            for (final line in decoded.split('\n')) {
              final l = line.trim();
              if (l.isEmpty) continue;
              final s = _parseSingleSubscriptionLine(l);
              if (s != null) parsed.add(s);
            }
          } catch (_) {
            // not base64, treat as newline-separated links
            for (final line in body.split('\n')) {
              final l = line.trim();
              if (l.isEmpty) continue;
              final s = _parseSingleSubscriptionLine(l);
              if (s != null) parsed.add(s);
            }
          }
        }
      } else {
        // direct link string (vmess://, vless://, ss://) or base64
        if (link.startsWith('vmess://') || link.startsWith('vless://') || link.startsWith('ss://')) {
          final s = _parseSingleSubscriptionLine(link);
          if (s != null) parsed.add(s);
        } else {
          // try base64 decode
          try {
            final decoded = utf8.decode(base64.decode(link.trim()));
            for (final line in decoded.split('\n')) {
              final l = line.trim();
              if (l.isEmpty) continue;
              final s = _parseSingleSubscriptionLine(l);
              if (s != null) parsed.add(s);
            }
          } catch (_) {
            // couldn't parse
          }
        }
      }
    } catch (e) {
      rethrow;
    }

    return parsed;
  }

  Server? _parseSingleSubscriptionLine(String line) {
    try {
      if (line.startsWith('vmess://')) {
        final payload = line.substring('vmess://'.length);
        final decoded = utf8.decode(base64.decode(payload));
        final Map<String, dynamic> obj = jsonDecode(decoded);
        return Server.fromJson({
          'id': obj['id'] ?? obj['ps'],
          'name': obj['ps'] ?? obj['name'],
          'host': obj['add'],
          'port': obj['port'] is String ? int.tryParse(obj['port']) ?? 0 : obj['port'],
          'protocol': 'vmess',
          'tls': (obj['tls'] == 'tls'),
        });
      }
      if (line.startsWith('vless://')) {
        // vless://uuid@host:port?params#name
        final withoutScheme = line.substring('vless://'.length);
        final at = withoutScheme.split('@');
        if (at.length == 2) {
          final id = at[0];
          final rest = at[1];
          final hostPort = rest.split('/')[0].split('?')[0];
          final hp = hostPort.split(':');
          final host = hp[0];
          final port = hp.length > 1 ? int.tryParse(hp[1]) ?? 0 : 0;
          return Server.fromJson({'id': id, 'host': host, 'port': port, 'protocol': 'vless'});
        }
      }
      if (line.startsWith('ss://')) {
        // ss://base64#name or ss://method:password@host:port
        final after = line.substring('ss://'.length);
        if (after.contains('@')) {
          final beforeAt = after.split('@')[0];
          final remainder = after.split('@')[1];
          final methodPass = beforeAt;
          final hp = remainder.split('#')[0];
          final host = hp.split(':')[0];
          final port = int.tryParse(hp.split(':')[1]) ?? 0;
          return Server.fromJson({'host': host, 'port': port, 'protocol': 'ss'});
        } else {
          try {
            final decoded = utf8.decode(base64.decode(after.split('#')[0]));
            // format method:password@host:port
            final parts = decoded.split('@');
            if (parts.length == 2) {
              final hp = parts[1];
              final host = hp.split(':')[0];
              final port = int.tryParse(hp.split(':')[1]) ?? 0;
              return Server.fromJson({'host': host, 'port': port, 'protocol': 'ss'});
            }
          } catch (_) {}
        }
      }
    } catch (_) {}
    return null;
  }

  /// بستن سرویس و آزاد کردن منابع
  void dispose() {
    _serversController.close();
  }
}