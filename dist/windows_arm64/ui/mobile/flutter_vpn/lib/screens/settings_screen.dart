import 'package:flutter/material.dart';
import 'package:vpn_client/screens/subscription_screen.dart';

class SettingsScreen extends StatefulWidget {
  const SettingsScreen({Key? key}) : super(key: key);

  @override
  State<SettingsScreen> createState() => _SettingsScreenState();
}

class _SettingsScreenState extends State<SettingsScreen> {
  bool _ipv6Enabled = true;
  bool _autoConnect = false;
  bool _darkMode = false;
  String _selectedTheme = 'system';
  String _tunnelMode = 'tcp_and_udp';
  String _selectedCore = 'sing-box';

  @override
  Widget build(BuildContext context) {
    return Scaffold(
      appBar: AppBar(
        title: const Text('Settings'),
      ),
      body: ListView(
        children: [
          // Connection settings
          const ListTile(
            title: Text(
              'Connection',
              style: TextStyle(fontWeight: FontWeight.bold),
            ),
          ),
          
          SwitchListTile(
            title: const Text('IPv6 Support'),
            subtitle: const Text('Enable IPv6 connectivity'),
            value: _ipv6Enabled,
            onChanged: (bool value) {
              setState(() {
                _ipv6Enabled = value;
              });
            },
          ),
          
          SwitchListTile(
            title: const Text('Auto Connect'),
            subtitle: const Text('Automatically connect to the fastest server'),
            value: _autoConnect,
            onChanged: (bool value) {
              setState(() {
                _autoConnect = value;
              });
            },
          ),
          
          ListTile(
            title: const Text('Tunnel Mode'),
            subtitle: Text(_tunnelMode == 'tcp_and_udp' 
                ? 'TCP and UDP' 
                : _tunnelMode == 'tcp_only' 
                    ? 'TCP Only' 
                    : 'UDP Only'),
            trailing: const Icon(Icons.arrow_forward_ios),
            onTap: () {
              _showTunnelModeDialog();
            },
          ),
          
          ListTile(
            title: const Text('VPN Core'),
            subtitle: Text(_selectedCore == 'sing-box' 
                ? 'Sing-box (Recommended)' 
                : 'XRay Core'),
            trailing: const Icon(Icons.arrow_forward_ios),
            onTap: () {
              _showVPNCoredialog();
            },
          ),
          
          const Divider(),
          
          // Appearance settings
          const ListTile(
            title: Text(
              'Appearance',
              style: TextStyle(fontWeight: FontWeight.bold),
            ),
          ),
          
          SwitchListTile(
            title: const Text('Dark Mode'),
            subtitle: const Text('Enable dark theme'),
            value: _darkMode,
            onChanged: (bool value) {
              setState(() {
                _darkMode = value;
              });
            },
          ),
          
          ListTile(
            title: const Text('Theme'),
            subtitle: Text(_selectedTheme == 'system' 
                ? 'System Default' 
                : _selectedTheme == 'light' 
                    ? 'Light' 
                    : 'Dark'),
            trailing: const Icon(Icons.arrow_forward_ios),
            onTap: () {
              _showThemeDialog();
            },
          ),
          
          const Divider(),
          
          // Subscription settings
          const ListTile(
            title: Text(
              'Subscription',
              style: TextStyle(fontWeight: FontWeight.bold),
            ),
          ),
          
          ListTile(
            title: const Text('Import Subscription'),
            subtitle: const Text('Add servers from subscription link'),
            trailing: const Icon(Icons.arrow_forward_ios),
            onTap: () {
              Navigator.push(
                context,
                MaterialPageRoute(builder: (context) => const SubscriptionScreen()),
              );
            },
          ),
          
          const Divider(),
          
          // About section
          const ListTile(
            title: Text(
              'About',
              style: TextStyle(fontWeight: FontWeight.bold),
            ),
          ),
          
          ListTile(
            title: const Text('Version'),
            subtitle: const Text('1.0.0'),
          ),
          
          ListTile(
            title: const Text('License'),
            subtitle: const Text('MIT License'),
            onTap: () {
              // Show license information
              _showLicenseInfo();
            },
          ),
          
          ListTile(
            title: const Text('Source Code'),
            subtitle: const Text('github.com/vpn-client'),
            onTap: () {
              // Open source code link
              _openSourceCode();
            },
          ),
        ],
      ),
    );
  }
  
  void _showTunnelModeDialog() {
    showDialog(
      context: context,
      builder: (BuildContext context) {
        return AlertDialog(
          title: const Text('Tunnel Mode'),
          content: Column(
            mainAxisSize: MainAxisSize.min,
            children: [
              RadioListTile<String>(
                title: const Text('TCP and UDP'),
                value: 'tcp_and_udp',
                groupValue: _tunnelMode,
                onChanged: (String? value) {
                  setState(() {
                    _tunnelMode = value!;
                  });
                  Navigator.of(context).pop();
                },
              ),
              RadioListTile<String>(
                title: const Text('TCP Only'),
                value: 'tcp_only',
                groupValue: _tunnelMode,
                onChanged: (String? value) {
                  setState(() {
                    _tunnelMode = value!;
                  });
                  Navigator.of(context).pop();
                },
              ),
              RadioListTile<String>(
                title: const Text('UDP Only'),
                value: 'udp_only',
                groupValue: _tunnelMode,
                onChanged: (String? value) {
                  setState(() {
                    _tunnelMode = value!;
                  });
                  Navigator.of(context).pop();
                },
              ),
            ],
          ),
        );
      },
    );
  }
  
  void _showThemeDialog() {
    showDialog(
      context: context,
      builder: (BuildContext context) {
        return AlertDialog(
          title: const Text('Theme'),
          content: Column(
            mainAxisSize: MainAxisSize.min,
            children: [
              RadioListTile<String>(
                title: const Text('System Default'),
                value: 'system',
                groupValue: _selectedTheme,
                onChanged: (String? value) {
                  setState(() {
                    _selectedTheme = value!;
                  });
                  Navigator.of(context).pop();
                },
              ),
              RadioListTile<String>(
                title: const Text('Light'),
                value: 'light',
                groupValue: _selectedTheme,
                onChanged: (String? value) {
                  setState(() {
                    _selectedTheme = value!;
                  });
                  Navigator.of(context).pop();
                },
              ),
              RadioListTile<String>(
                title: const Text('Dark'),
                value: 'dark',
                groupValue: _selectedTheme,
                onChanged: (String? value) {
                  setState(() {
                    _selectedTheme = value!;
                  });
                  Navigator.of(context).pop();
                },
              ),
            ],
          ),
        );
      },
    );
  }
  
  void _showVPNCoredialog() {
    showDialog(
      context: context,
      builder: (BuildContext context) {
        return AlertDialog(
          title: const Text('VPN Core'),
          content: Column(
            mainAxisSize: MainAxisSize.min,
            children: [
              RadioListTile<String>(
                title: const Text('Sing-box (Recommended)'),
                value: 'sing-box',
                groupValue: _selectedCore,
                onChanged: (String? value) {
                  setState(() {
                    _selectedCore = value!;
                  });
                  Navigator.of(context).pop();
                },
              ),
              RadioListTile<String>(
                title: const Text('XRay Core'),
                value: 'xray',
                groupValue: _selectedCore,
                onChanged: (String? value) {
                  setState(() {
                    _selectedCore = value!;
                  });
                  Navigator.of(context).pop();
                },
              ),
            ],
          ),
        );
      },
    );
  }
  
  void _showLicenseInfo() {
    showDialog(
      context: context,
      builder: (BuildContext context) {
        return AlertDialog(
          title: const Text('License'),
          content: const SingleChildScrollView(
            child: Text(
              'MIT License\n\n'
              'Copyright (c) 2025 VPN Client\n\n'
              'Permission is hereby granted, free of charge, to any person obtaining a copy '
              'of this software and associated documentation files (the "Software"), to deal '
              'in the Software without restriction, including without limitation the rights '
              'to use, copy, modify, merge, publish, distribute, sublicense, and/or sell '
              'copies of the Software, and to permit persons to whom the Software is '
              'furnished to do so, subject to the following conditions:\n\n'
              'The above copyright notice and this permission notice shall be included in all '
              'copies or substantial portions of the Software.\n\n'
              'THE SOFTWARE IS PROVIDED "AS IS", WITHOUT WARRANTY OF ANY KIND, EXPRESS OR '
              'IMPLIED, INCLUDING BUT NOT LIMITED TO THE WARRANTIES OF MERCHANTABILITY, '
              'FITNESS FOR A PARTICULAR PURPOSE AND NONINFRINGEMENT. IN NO EVENT SHALL THE '
              'AUTHORS OR COPYRIGHT HOLDERS BE LIABLE FOR ANY CLAIM, DAMAGES OR OTHER '
              'LIABILITY, WHETHER IN AN ACTION OF CONTRACT, TORT OR OTHERWISE, ARISING FROM, '
              'OUT OF OR IN CONNECTION WITH THE SOFTWARE OR THE USE OR OTHER DEALINGS IN THE '
              'SOFTWARE.',
            ),
          ),
          actions: [
            TextButton(
              onPressed: () {
                Navigator.of(context).pop();
              },
              child: const Text('Close'),
            ),
          ],
        );
      },
    );
  }
  
  void _openSourceCode() {
    ScaffoldMessenger.of(context).showSnackBar(
      const SnackBar(
        content: Text('In a real app, this would open the source code repository'),
        duration: Duration(seconds: 2),
      ),
    );
  }
}