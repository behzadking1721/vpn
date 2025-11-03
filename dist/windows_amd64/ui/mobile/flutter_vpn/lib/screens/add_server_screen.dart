import 'package:flutter/material.dart';
import 'package:provider/provider.dart';
import 'package:vpn_client/models/server.dart';
import 'package:vpn_client/providers/app_provider.dart';
import 'package:vpn_client/utils/format_utils.dart';

class AddServerScreen extends StatefulWidget {
  const AddServerScreen({Key? key}) : super(key: key);

  @override
  State<AddServerScreen> createState() => _AddServerScreenState();
}

class _AddServerScreenState extends State<AddServerScreen> {
  final _formKey = GlobalKey<FormState>();
  
  // Controllers
  final _nameController = TextEditingController();
  final _hostController = TextEditingController();
  final _portController = TextEditingController(text: '443');
  final _passwordController = TextEditingController();
  final _encryptionController = TextEditingController(text: 'auto');
  final _methodController = TextEditingController();
  final _sniController = TextEditingController();
  final _fingerprintController = TextEditingController();
  final _remarkController = TextEditingController();
  
  // Values
  String _selectedProtocol = 'VMess';
  bool _tlsEnabled = false;
  bool _serverEnabled = true;
  
  final List<String> _protocols = [
    'VMess',
    'VLESS',
    'Trojan',
    'Shadowsocks',
    'Reality',
    'Hysteria2',
    'TUIC',
    'SSH',
  ];

  @override
  void dispose() {
    _nameController.dispose();
    _hostController.dispose();
    _portController.dispose();
    _passwordController.dispose();
    _encryptionController.dispose();
    _methodController.dispose();
    _sniController.dispose();
    _fingerprintController.dispose();
    _remarkController.dispose();
    super.dispose();
  }

  @override
  Widget build(BuildContext context) {
    return Scaffold(
      appBar: AppBar(
        title: const Text('Add Server'),
        actions: [
          IconButton(
            icon: const Icon(Icons.save),
            onPressed: _saveServer,
          ),
        ],
      ),
      body: Form(
        key: _formKey,
        child: SingleChildScrollView(
          padding: const EdgeInsets.all(16.0),
          child: Column(
            crossAxisAlignment: CrossAxisAlignment.start,
            children: [
              // Name field
              TextFormField(
                controller: _nameController,
                decoration: const InputDecoration(
                  labelText: 'Server Name',
                  hintText: 'Enter server name',
                ),
                validator: (value) {
                  if (value == null || value.isEmpty) {
                    return 'Please enter a server name';
                  }
                  return null;
                },
              ),
              
              const SizedBox(height: 16),
              
              // Host field
              TextFormField(
                controller: _hostController,
                decoration: const InputDecoration(
                  labelText: 'Host',
                  hintText: 'Enter server host or IP',
                ),
                validator: (value) {
                  if (value == null || value.isEmpty) {
                    return 'Please enter a host';
                  }
                  return null;
                },
              ),
              
              const SizedBox(height: 16),
              
              // Port field
              TextFormField(
                controller: _portController,
                decoration: const InputDecoration(
                  labelText: 'Port',
                  hintText: 'Enter port number',
                ),
                keyboardType: TextInputType.number,
                validator: (value) {
                  if (value == null || value.isEmpty) {
                    return 'Please enter a port';
                  }
                  final port = int.tryParse(value);
                  if (port == null || port < 1 || port > 65535) {
                    return 'Please enter a valid port (1-65535)';
                  }
                  return null;
                },
              ),
              
              const SizedBox(height: 16),
              
              // Protocol dropdown
              InputDecorator(
                decoration: const InputDecoration(
                  labelText: 'Protocol',
                ),
                child: DropdownButtonHideUnderline(
                  child: DropdownButton<String>(
                    value: _selectedProtocol,
                    isDense: true,
                    onChanged: (String? newValue) {
                      if (newValue != null) {
                        setState(() {
                          _selectedProtocol = newValue;
                        });
                      }
                    },
                    items: _protocols.map<DropdownMenuItem<String>>((String value) {
                      return DropdownMenuItem<String>(
                        value: value,
                        child: Text(value),
                      );
                    }).toList(),
                  ),
                ),
              ),
              
              const SizedBox(height: 16),
              
              // Password/UUID field
              TextFormField(
                controller: _passwordController,
                decoration: InputDecoration(
                  labelText: _selectedProtocol == 'Shadowsocks' ? 'Password' : 'UUID/Password',
                  hintText: _selectedProtocol == 'Shadowsocks' 
                    ? 'Enter password' 
                    : 'Enter UUID or password',
                ),
                validator: (value) {
                  if (value == null || value.isEmpty) {
                    return 'Please enter a password or UUID';
                  }
                  return null;
                },
              ),
              
              const SizedBox(height: 16),
              
              // Method field (for Shadowsocks)
              if (_selectedProtocol == 'Shadowsocks') ...[
                TextFormField(
                  controller: _methodController,
                  decoration: const InputDecoration(
                    labelText: 'Encryption Method',
                    hintText: 'Enter encryption method (e.g., aes-256-gcm)',
                  ),
                  validator: (value) {
                    if (value == null || value.isEmpty) {
                      return 'Please enter an encryption method';
                    }
                    return null;
                  },
                ),
                const SizedBox(height: 16),
              ] else ...[
                // Encryption field (for other protocols)
                TextFormField(
                  controller: _encryptionController,
                  decoration: const InputDecoration(
                    labelText: 'Encryption',
                    hintText: 'Enter encryption method',
                  ),
                ),
                const SizedBox(height: 16),
              ],
              
              // TLS switch
              SwitchListTile(
                title: const Text('TLS Enabled'),
                value: _tlsEnabled,
                onChanged: (bool value) {
                  setState(() {
                    _tlsEnabled = value;
                  });
                },
              ),
              
              const SizedBox(height: 16),
              
              // SNI field
              if (_tlsEnabled) ...[
                TextFormField(
                  controller: _sniController,
                  decoration: const InputDecoration(
                    labelText: 'SNI (optional)',
                    hintText: 'Enter Server Name Indication',
                  ),
                ),
                const SizedBox(height: 16),
                
                // Fingerprint field
                TextFormField(
                  controller: _fingerprintController,
                  decoration: const InputDecoration(
                    labelText: 'Fingerprint (optional)',
                    hintText: 'Enter TLS fingerprint',
                  ),
                ),
                const SizedBox(height: 16),
              ],
              
              // Remark field
              TextFormField(
                controller: _remarkController,
                decoration: const InputDecoration(
                  labelText: 'Remark (optional)',
                  hintText: 'Enter any remarks',
                ),
                maxLines: 3,
              ),
              
              const SizedBox(height: 16),
              
              // Enabled switch
              SwitchListTile(
                title: const Text('Server Enabled'),
                value: _serverEnabled,
                onChanged: (bool value) {
                  setState(() {
                    _serverEnabled = value;
                  });
                },
              ),
              
              const SizedBox(height: 32),
              
              // Save button
              Center(
                child: ElevatedButton(
                  onPressed: _saveServer,
                  child: const Text('Save Server'),
                ),
              ),
            ],
          ),
        ),
      ),
    );
  }
  
  void _saveServer() {
    if (_formKey.currentState!.validate()) {
      // Create new server
      final newServer = Server(
        id: DateTime.now().millisecondsSinceEpoch.toString(),
        name: _nameController.text,
        host: _hostController.text,
        port: int.parse(_portController.text),
        protocol: _selectedProtocol,
        password: _passwordController.text,
        method: _selectedProtocol == 'Shadowsocks' ? _methodController.text : null,
        encryption: _selectedProtocol != 'Shadowsocks' ? _encryptionController.text : null,
        tls: _tlsEnabled,
        sni: _sniController.text.isNotEmpty ? _sniController.text : null,
        fingerprint: _fingerprintController.text.isNotEmpty ? _fingerprintController.text : null,
        remark: _remarkController.text.isNotEmpty ? _remarkController.text : null,
        enabled: _serverEnabled,
        ping: 0,
      );
      
      // Add server using provider
      final appProvider = Provider.of<AppProvider>(context, listen: false);
      appProvider.addServer(newServer);
      
      // Show success message
      ScaffoldMessenger.of(context).showSnackBar(
        const SnackBar(
          content: Text('Server added successfully'),
          duration: Duration(seconds: 2),
        ),
      );
      
      // Navigate back
      Navigator.pop(context);
    }
  }
}