class Server {
  final String id;
  final String name;
  final String host;
  final int port;
  final String protocol;
  final String? encryption;
  final String? password;
  final String? method;
  final bool tls;
  final String? sni;
  final String? fingerprint;
  final String? remark;
  final bool enabled;
  final int ping;
  final bool connected;

  Server({
    required this.id,
    required this.name,
    required this.host,
    required this.port,
    required this.protocol,
    this.encryption,
    this.password,
    this.method,
    this.tls = false,
    this.sni,
    this.fingerprint,
    this.remark,
    this.enabled = true,
    this.ping = 0,
    this.connected = false,
  });

  Server copyWith({
    String? id,
    String? name,
    String? host,
    int? port,
    String? protocol,
    String? encryption,
    String? password,
    String? method,
    bool? tls,
    String? sni,
    String? fingerprint,
    String? remark,
    bool? enabled,
    int? ping,
    bool? connected,
  }) {
    return Server(
      id: id ?? this.id,
      name: name ?? this.name,
      host: host ?? this.host,
      port: port ?? this.port,
      protocol: protocol ?? this.protocol,
      encryption: encryption ?? this.encryption,
      password: password ?? this.password,
      method: method ?? this.method,
      tls: tls ?? this.tls,
      sni: sni ?? this.sni,
      fingerprint: fingerprint ?? this.fingerprint,
      remark: remark ?? this.remark,
      enabled: enabled ?? this.enabled,
      ping: ping ?? this.ping,
      connected: connected ?? this.connected,
    );
  }

  Map<String, dynamic> toJson() {
    return {
      'id': id,
      'name': name,
      'host': host,
      'port': port,
      'protocol': protocol,
      'encryption': encryption,
      'password': password,
      'method': method,
      'tls': tls,
      'sni': sni,
      'fingerprint': fingerprint,
      'remark': remark,
      'enabled': enabled,
      'ping': ping,
      'connected': connected,
    };
  }

  factory Server.fromJson(Map<String, dynamic> json) {
    return Server(
      id: json['id'] as String,
      name: json['name'] as String,
      host: json['host'] as String,
      port: json['port'] as int,
      protocol: json['protocol'] as String,
      encryption: json['encryption'] as String?,
      password: json['password'] as String?,
      method: json['method'] as String?,
      tls: json['tls'] as bool? ?? false,
      sni: json['sni'] as String?,
      fingerprint: json['fingerprint'] as String?,
      remark: json['remark'] as String?,
      enabled: json['enabled'] as bool? ?? true,
      ping: json['ping'] as int? ?? 0,
      connected: json['connected'] as bool? ?? false,
    );
  }

  @override
  bool operator ==(Object other) {
    if (identical(this, other)) return true;
    return other is Server && other.id == id;
  }

  @override
  int get hashCode => id.hashCode;
}