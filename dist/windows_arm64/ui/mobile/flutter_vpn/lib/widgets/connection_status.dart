import 'package:flutter/material.dart';
import 'package:provider/provider.dart';
import 'package:vpn_client/providers/app_provider.dart';
import 'package:vpn_client/utils/format_utils.dart';

class ConnectionStatusWidget extends StatelessWidget {
  const ConnectionStatusWidget({Key? key}) : super(key: key);

  @override
  Widget build(BuildContext context) {
    return Consumer<AppProvider>(
      builder: (context, appProvider, child) {
        final status = appProvider.connectionStatus;
        final server = appProvider.currentServer;
        
        return Padding(
          padding: const EdgeInsets.all(16.0),
          child: Column(
            crossAxisAlignment: CrossAxisAlignment.start,
            children: [
              // Connection status card
              Card(
                child: Padding(
                  padding: const EdgeInsets.all(16.0),
                  child: Column(
                    crossAxisAlignment: CrossAxisAlignment.start,
                    children: [
                      const Text(
                        'Connection Status',
                        style: TextStyle(
                          fontSize: 20,
                          fontWeight: FontWeight.bold,
                        ),
                      ),
                      const SizedBox(height: 10),
                      Row(
                        children: [
                          Container(
                            width: 12,
                            height: 12,
                            decoration: BoxDecoration(
                              color: status.isConnected ? Colors.green : Colors.red,
                              shape: BoxShape.circle,
                            ),
                          ),
                          const SizedBox(width: 10),
                          Text(
                            status.isConnected ? 'Connected' : 'Disconnected',
                            style: TextStyle(
                              fontSize: 16,
                              color: status.isConnected ? Colors.green : Colors.red,
                            ),
                          ),
                        ],
                      ),
                      const SizedBox(height: 10),
                      if (server != null) ...[
                        Text('Server: ${server.name}'),
                        Text('Protocol: ${server.protocol}'),
                        Text('Host: ${server.host}:${server.port}'),
                      ] else if (status.isConnected) ...[
                        Text('Server: ${status.serverName ?? 'Unknown'}'),
                        Text('Protocol: ${status.protocol ?? 'Unknown'}'),
                      ] else ...[
                        const Text('No server connected'),
                      ],
                    ],
                  ),
                ),
              ),
              
              const SizedBox(height: 20),
              
              // Data usage
              Card(
                child: Padding(
                  padding: const EdgeInsets.all(16.0),
                  child: Column(
                    crossAxisAlignment: CrossAxisAlignment.start,
                    children: [
                      const Text(
                        'Data Usage',
                        style: TextStyle(
                          fontSize: 20,
                          fontWeight: FontWeight.bold,
                        ),
                      ),
                      const SizedBox(height: 10),
                      _buildDataUsageRow(
                        'Upload', 
                        FormatUtils.formatBytes(status.dataSent), 
                        Colors.blue
                      ),
                      const SizedBox(height: 10),
                      _buildDataUsageRow(
                        'Download', 
                        FormatUtils.formatBytes(status.dataReceived), 
                        Colors.green
                      ),
                      const SizedBox(height: 10),
                      _buildDataUsageRow(
                        'Total', 
                        FormatUtils.formatBytes(status.dataSent + status.dataReceived), 
                        Colors.purple
                      ),
                    ],
                  ),
                ),
              ),
              
              const SizedBox(height: 20),
              
              // Connection time
              if (status.isConnected && status.startTime != null) ...[
                Card(
                  child: Padding(
                    padding: const EdgeInsets.all(16.0),
                    child: Column(
                      crossAxisAlignment: CrossAxisAlignment.start,
                      children: [
                        const Text(
                          'Connection Time',
                          style: TextStyle(
                            fontSize: 20,
                            fontWeight: FontWeight.bold,
                          ),
                        ),
                        const SizedBox(height: 10),
                        Text(
                          FormatUtils.formatDuration(
                            DateTime.now().difference(status.startTime!)
                          )
                        ),
                      ],
                    ),
                  ),
                ),
                
                const SizedBox(height: 20),
              ],
              
              // Disconnect/connect button
              Center(
                child: ElevatedButton(
                  style: ElevatedButton.styleFrom(
                    backgroundColor: status.isConnected ? Colors.red : Colors.blue,
                    padding: const EdgeInsets.symmetric(horizontal: 50, vertical: 15),
                  ),
                  onPressed: status.isConnected 
                    ? () => _disconnect(context, appProvider)
                    : () => _connect(context, appProvider),
                  child: Text(
                    status.isConnected ? 'Disconnect' : 'Connect',
                    style: const TextStyle(fontSize: 18),
                  ),
                ),
              ),
            ],
          ),
        );
      },
    );
  }
  
  Widget _buildDataUsageRow(String label, String value, Color color) {
    return Row(
      mainAxisAlignment: MainAxisAlignment.spaceBetween,
      children: [
        Text(label),
        Text(
          value,
          style: TextStyle(
            color: color,
            fontWeight: FontWeight.bold,
          ),
        ),
      ],
    );
  }
  
  void _connect(BuildContext context, AppProvider appProvider) {
    if (appProvider.selectedServer != null) {
      appProvider.connectToSelectedServer().then((success) {
        if (!success) {
          ScaffoldMessenger.of(context).showSnackBar(
            const SnackBar(
              content: Text('Failed to connect'),
              duration: Duration(seconds: 2),
            ),
          );
        }
      });
    } else {
      ScaffoldMessenger.of(context).showSnackBar(
        const SnackBar(
          content: Text('Please select a server first'),
          duration: Duration(seconds: 2),
        ),
      );
    }
  }
  
  void _disconnect(BuildContext context, AppProvider appProvider) {
    appProvider.disconnect().then((success) {
      if (success) {
        ScaffoldMessenger.of(context).showSnackBar(
          const SnackBar(
            content: Text('Disconnected successfully'),
            duration: Duration(seconds: 2),
          ),
        );
      }
    });
  }
}