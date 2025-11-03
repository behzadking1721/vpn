import 'package:flutter/material.dart';

class QuickConnect extends StatelessWidget {
  const QuickConnect({Key? key}) : super(key: key);

  @override
  Widget build(BuildContext context) {
    return Scaffold(
      body: Padding(
        padding: const EdgeInsets.all(16.0),
        child: Column(
          crossAxisAlignment: CrossAxisAlignment.start,
          children: [
            const Text(
              'Quick Connect',
              style: TextStyle(
                fontSize: 24,
                fontWeight: FontWeight.bold,
              ),
            ),
            
            const SizedBox(height: 20),
            
            // Fastest server card
            Card(
              child: Padding(
                padding: const EdgeInsets.all(16.0),
                child: Column(
                  crossAxisAlignment: CrossAxisAlignment.start,
                  children: [
                    const Text(
                      'Fastest Server',
                      style: TextStyle(
                        fontSize: 20,
                        fontWeight: FontWeight.bold,
                      ),
                    ),
                    const SizedBox(height: 10),
                    const Text('🇯🇵 Japan Server'),
                    const Text('Ping: 22 ms'),
                    const SizedBox(height: 10),
                    ElevatedButton(
                      onPressed: () {
                        // Connect to fastest server
                      },
                      child: const Text('Connect'),
                    ),
                  ],
                ),
              ),
            ),
            
            const SizedBox(height: 20),
            
            // Recent connections
            const Text(
              'Recent Connections',
              style: TextStyle(
                fontSize: 20,
                fontWeight: FontWeight.bold,
              ),
            ),
            
            const SizedBox(height: 10),
            
            Expanded(
              child: ListView(
                children: [
                  Card(
                    child: ListTile(
                      title: const Text('🇺🇸 USA Server'),
                      subtitle: const Text('Shadowsocks - 2 hours ago'),
                      trailing: const Icon(Icons.link),
                      onTap: () {
                        // Reconnect to this server
                      },
                    ),
                  ),
                  Card(
                    child: ListTile(
                      title: const Text('🇬🇧 UK Server'),
                      subtitle: const Text('Trojan - 1 day ago'),
                      trailing: const Icon(Icons.link),
                      onTap: () {
                        // Reconnect to this server
                      },
                    ),
                  ),
                  Card(
                    child: ListTile(
                      title: const Text('🇩🇪 Germany Server'),
                      subtitle: const Text('VMess - 3 days ago'),
                      trailing: const Icon(Icons.link),
                      onTap: () {
                        // Reconnect to this server
                      },
                    ),
                  ),
                ],
              ),
            ),
          ],
        ),
      ),
    );
  }
}