class ConnectionStatus {
  final bool isConnected;
  final String? serverId;
  final DateTime? startTime;
  final int dataSent;
  final int dataReceived;
  final String? protocol;
  final String? serverName;

  ConnectionStatus({
    required this.isConnected,
    this.serverId,
    this.startTime,
    this.dataSent = 0,
    this.dataReceived = 0,
    this.protocol,
    this.serverName,
  });

  ConnectionStatus copyWith({
    bool? isConnected,
    String? serverId,
    DateTime? startTime,
    int? dataSent,
    int? dataReceived,
    String? protocol,
    String? serverName,
  }) {
    return ConnectionStatus(
      isConnected: isConnected ?? this.isConnected,
      serverId: serverId ?? this.serverId,
      startTime: startTime ?? this.startTime,
      dataSent: dataSent ?? this.dataSent,
      dataReceived: dataReceived ?? this.dataReceived,
      protocol: protocol ?? this.protocol,
      serverName: serverName ?? this.serverName,
    );
  }

  Map<String, dynamic> toJson() {
    return {
      'isConnected': isConnected,
      'serverId': serverId,
      'startTime': startTime?.toIso8601String(),
      'dataSent': dataSent,
      'dataReceived': dataReceived,
      'protocol': protocol,
      'serverName': serverName,
    };
  }

  factory ConnectionStatus.fromJson(Map<String, dynamic> json) {
    return ConnectionStatus(
      isConnected: json['isConnected'] as bool,
      serverId: json['serverId'] as String?,
      startTime: json['startTime'] != null
          ? DateTime.parse(json['startTime'] as String)
          : null,
      dataSent: json['dataSent'] as int? ?? 0,
      dataReceived: json['dataReceived'] as int? ?? 0,
      protocol: json['protocol'] as String?,
      serverName: json['serverName'] as String?,
    );
  }
}