import 'package:flutter/material.dart';
import 'package:provider/provider.dart';
import 'package:vpn_client/models/server.dart';
import 'package:vpn_client/providers/app_provider.dart';
import 'package:vpn_client/screens/qr_scanner_screen.dart';
import 'package:vpn_client/screens/server_details_screen.dart';

class ServerList extends StatefulWidget {
  const ServerList({Key? key}) : super(key: key);

  @override
  State<ServerList> createState() => _ServerListState();
}

class _ServerListState extends State<ServerList> {
  @override
  Widget build(BuildContext context) {
    return Consumer<AppProvider>(
      builder: (context, appProvider, child) {
        final servers = appProvider.allServers;
        
        return Scaffold(
          body: RefreshIndicator(
            onRefresh: () async {
              await appProvider.pingServers();
            },
            child: ListView.builder(
              itemCount: servers.length,
              itemBuilder: (context, index) {
                final server = servers[index];
                final isSelected = appProvider.selectedServer?.id == server.id;
                final isConnected = server.id == appProvider.currentServer?.id;
                
                return Card(
                  margin: const EdgeInsets.symmetric(horizontal: 10, vertical: 5),
                  elevation: isSelected ? 4 : 2,
                  color: isSelected ? Colors.blue.withOpacity(0.1) : Colors.white,
                  child: ListTile(
                    leading: CircleAvatar(
                      backgroundColor: isConnected 
                          ? Colors.green 
                          : server.enabled 
                              ? Colors.grey 
                              : Colors.red,
                      child: Icon(
                        isConnected 
                            ? Icons.link 
                            : server.enabled 
                                ? Icons.link_off 
                                : Icons.block,
                        color: Colors.white,
                      ),
                    ),
                    title: Text(
                      server.name,
                      style: TextStyle(
                        fontWeight: isSelected ? FontWeight.bold : FontWeight.normal,
                      ),
                    ),
                    subtitle: Column(
                      crossAxisAlignment: CrossAxisAlignment.start,
                      children: [
                        Text('${server.host}:${server.port} (${server.protocol})'),
                        Text('Ping: ${server.ping} ms'),
                      ],
                    ),
                    trailing: IconButton(
                      icon: const Icon(Icons.more_vert),
                      onPressed: () {
                        // Show server options
                        _showServerOptions(context, server);
                      },
                    ),
                    onTap: () {
                      // Navigate to server details
                      Navigator.push(
                        context,
                        MaterialPageRoute(
                          builder: (context) => ServerDetailsScreen(server: server),
                        ),
                      );
                    },
                  ),
                );
              },
            ),
          ),
          floatingActionButton: Column(
            mainAxisAlignment: MainAxisAlignment.end,
            children: [
              FloatingActionButton(
                heroTag: 'qr',
                child: const Icon(Icons.qr_code_scanner),
                onPressed: () {
                  // Scan QR code
                  _scanQRCode(context, appProvider);
                },
              ),
              const SizedBox(height: 10),
              FloatingActionButton(
                heroTag: 'add',
                child: const Icon(Icons.add),
                onPressed: () {
                  // Add new server
                  _addNewServer(context);
                },
              ),
            ],
          ),
        );
      },
    );
  }

  void _scanQRCode(BuildContext context, AppProvider appProvider) async {
    appProvider.startQRScan();
    
    final result = await Navigator.push(
      context,
      MaterialPageRoute(builder: (context) => const QRScannerScreen()),
    );
    
    appProvider.endQRScan();
    
    if (result != null) {
      ScaffoldMessenger.of(context).showSnackBar(
        SnackBar(
          content: Text('Scanned: ${result.substring(0, 30)}...'),
          duration: const Duration(seconds: 3),
        ),
      );
      
      // In a real app, you would parse the result and add the server
      _parseAndAddServer(context, appProvider, result);
    }
  }

  void _addNewServer(BuildContext context) {
    ScaffoldMessenger.of(context).showSnackBar(
      const SnackBar(
        content: Text('Add server functionality would be implemented here'),
        duration: Duration(seconds: 2),
      ),
    );
    
    // In a real app, this would open the add server form
  }

  void _showServerOptions(BuildContext context, Server server) {
    showModalBottomSheet(
      context: context,
      builder: (BuildContext context) {
        return SafeArea(
          child: Column(
            mainAxisSize: MainAxisSize.min,
            children: [
              ListTile(
                leading: const Icon(Icons.info),
                title: const Text('Details'),
                onTap: () {
                  Navigator.pop(context);
                  Navigator.push(
                    context,
                    MaterialPageRoute(
                      builder: (context) => ServerDetailsScreen(server: server),
                    ),
                  );
                },
              ),
              ListTile(
                leading: const Icon(Icons.edit),
                title: const Text('Edit'),
                onTap: () {
                  Navigator.pop(context);
                  // Edit server
                  ScaffoldMessenger.of(context).showSnackBar(
                    const SnackBar(
                      content: Text('Edit functionality would be implemented here'),
                      duration: Duration(seconds: 2),
                    ),
                  );
                },
              ),
              ListTile(
                leading: const Icon(Icons.delete),
                title: const Text('Delete'),
                onTap: () {
                  Navigator.pop(context);
                  // Delete server
                  _deleteServer(context, server);
                },
              ),
            ],
          ),
        );
      },
    );
  }
  
  void _parseAndAddServer(BuildContext context, AppProvider appProvider, String qrData) {
    // In a real app, you would parse the QR data and add the server
    // This is a simplified example
    ScaffoldMessenger.of(context).showSnackBar(
      const SnackBar(
        content: Text('Server would be parsed and added in a real app'),
        duration: Duration(seconds: 2),
      ),
    );
  }
  
  void _deleteServer(BuildContext context, Server server) {
    final appProvider = Provider.of<AppProvider>(context, listen: false);
    
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
                appProvider.removeServer(server.id);
                Navigator.pop(context); // Close dialog
                
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
}