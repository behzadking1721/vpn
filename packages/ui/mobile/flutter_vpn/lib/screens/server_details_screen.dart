import 'package:flutter/material.dart';
import 'package:provider/provider.dart';
import 'package:vpn_client/models/server.dart';
import 'package:vpn_client/providers/app_provider.dart';
import 'package:vpn_client/utils/format_utils.dart';

class ServerDetailsScreen extends StatelessWidget {
  final Server server;

  const ServerDetailsScreen({Key? key, required this.server}) : super(key: key);

  @override
  Widget build(BuildContext context) {
    return Scaffold(
      appBar: AppBar(
        title: Text(server.name),
        actions: [
          IconButton(
            icon: const Icon(Icons.edit),
            onPressed: () {
              // Navigate to edit server screen
              // For now, we'll just show a message
              ScaffoldMessenger.of(context).showSnackBar(
                const SnackBar(
                  content: Text('Edit server functionality would be implemented here'),
                  duration: Duration(seconds: 2),
                ),
              );
            },
          ),
          PopupMenuButton<String>(
            onSelected: (value) {
              switch (value) {
                case 'delete':
                  _deleteServer(context);
                  break;
                case 'share':
                  _shareServer(context);
                  break;
              }
            },
            itemBuilder: (BuildContext context) {
              return [
                const PopupMenuItem(
                  value: 'delete',
                  child: Text('Delete'),
                ),
                const PopupMenuItem(
                  value: 'share',
                  child: Text('Share'),
                ),
              ];
            },
          ),
        ],
      ),
      body: SingleChildScrollView(
        padding: const EdgeInsets.all(16.0),
        child: Column(
          crossAxisAlignment: CrossAxisAlignment.start,
          children: [
            // Server info card
            Card(
              child: Padding(
                padding: const EdgeInsets.all(16.0),
                child: Column(
                  crossAxisAlignment: CrossAxisAlignment.start,
                  children: [
                    const Text(
                      'Server Information',
                      style: TextStyle(
                        fontSize: 20,
                        fontWeight: FontWeight.bold,
                      ),
                    ),
                    const SizedBox(height: 10),
                    _buildInfoRow('Name', server.name),
                    _buildInfoRow('Host', '${server.host}:${server.port}'),
                    _buildInfoRow('Protocol', server.protocol),
                    if (server.method != null)
                      _buildInfoRow('Method', server.method!),
                    if (server.encryption != null)
                      _buildInfoRow('Encryption', server.encryption!),
                    if (server.password != null)
                      _buildInfoRow('Password/UUID', '••••••••'),
                    _buildInfoRow('Ping', FormatUtils.formatPing(server.ping)),
                    _buildInfoRow('Status', server.enabled ? 'Enabled' : 'Disabled'),
                    if (server.remark != null)
                      _buildInfoRow('Remark', server.remark!),
                  ],
                ),
              ),
            ),
            
            const SizedBox(height: 20),
            
            // TLS info card
            if (server.tls) ...[
              Card(
                child: Padding(
                  padding: const EdgeInsets.all(16.0),
                  child: Column(
                    crossAxisAlignment: CrossAxisAlignment.start,
                    children: [
                      const Text(
                        'TLS Configuration',
                        style: TextStyle(
                          fontSize: 20,
                          fontWeight: FontWeight.bold,
                        ),
                      ),
                      const SizedBox(height: 10),
                      _buildInfoRow('TLS', server.tls ? 'Enabled' : 'Disabled'),
                      if (server.sni != null)
                        _buildInfoRow('SNI', server.sni!),
                      if (server.fingerprint != null)
                        _buildInfoRow('Fingerprint', server.fingerprint!),
                    ],
                  ),
                ),
              ),
              
              const SizedBox(height: 20),
            ],
            
            // Connection button
            Center(
              child: ElevatedButton(
                onPressed: () {
                  final appProvider = Provider.of<AppProvider>(context, listen: false);
                  appProvider.selectServer(server);
                  Navigator.pop(context);
                  
                  // Show connection message
                  ScaffoldMessenger.of(context).showSnackBar(
                    SnackBar(
                      content: Text('Selected ${server.name}. Tap connect in status tab.'),
                      duration: const Duration(seconds: 2),
                    ),
                  );
                },
                child: const Text('Connect to this Server'),
              ),
            ),
          ],
        ),
      ),
    );
  }
  
  Widget _buildInfoRow(String label, String value) {
    return Padding(
      padding: const EdgeInsets.symmetric(vertical: 4.0),
      child: Row(
        crossAxisAlignment: CrossAxisAlignment.start,
        children: [
          SizedBox(
            width: 100,
            child: Text(
              '$label:',
              style: const TextStyle(fontWeight: FontWeight.bold),
            ),
          ),
          Expanded(
            child: Text(value),
          ),
        ],
      ),
    );
  }
  
  void _deleteServer(BuildContext context) {
    showDialog(
      context: context,
      builder: (BuildContext context) {
        return AlertDialog(
          title: const Text('Delete Server'),
          content: Text('Are you sure you want to delete ${server.name}?'),
          actions: [
            TextButton(
              onPressed: () {
                Navigator.pop(context);
              },
              child: const Text('Cancel'),
            ),
            TextButton(
              onPressed: () {
                final appProvider = Provider.of<AppProvider>(context, listen: false);
                appProvider.removeServer(server.id);
                Navigator.pop(context); // Close dialog
                Navigator.pop(context); // Close details screen
                
                ScaffoldMessenger.of(context).showSnackBar(
                  SnackBar(
                    content: Text('Deleted ${server.name}'),
                    duration: const Duration(seconds: 2),
                  ),
                );
              },
              child: const Text('Delete', style: TextStyle(color: Colors.red)),
            ),
          ],
        );
      },
    );
  }
  
  void _shareServer(BuildContext context) {
    // In a real app, you would generate a shareable link or QR code
    ScaffoldMessenger.of(context).showSnackBar(
      const SnackBar(
        content: Text('Share functionality would be implemented here'),
        duration: Duration(seconds: 2),
      ),
    );
  }
}