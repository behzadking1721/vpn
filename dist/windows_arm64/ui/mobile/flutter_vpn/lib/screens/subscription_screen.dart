import 'package:flutter/material.dart';
import 'package:provider/provider.dart';
import 'package:vpn_client/providers/app_provider.dart';

class SubscriptionScreen extends StatefulWidget {
  const SubscriptionScreen({Key? key}) : super(key: key);

  @override
  State<SubscriptionScreen> createState() => _SubscriptionScreenState();
}

class _SubscriptionScreenState extends State<SubscriptionScreen> {
  final _formKey = GlobalKey<FormState>();
  final _linkController = TextEditingController();
  bool _isProcessing = false;

  @override
  void dispose() {
    _linkController.dispose();
    super.dispose();
  }

  @override
  Widget build(BuildContext context) {
    return Scaffold(
      appBar: AppBar(
        title: const Text('Subscription'),
      ),
      body: Padding(
        padding: const EdgeInsets.all(16.0),
        child: Form(
          key: _formKey,
          child: Column(
            crossAxisAlignment: CrossAxisAlignment.start,
            children: [
              const Text(
                'Import Subscription',
                style: TextStyle(
                  fontSize: 24,
                  fontWeight: FontWeight.bold,
                ),
              ),
              
              const SizedBox(height: 20),
              
              const Text(
                'Enter your subscription link below to import servers:',
                style: TextStyle(
                  fontSize: 16,
                  color: Colors.grey,
                ),
              ),
              
              const SizedBox(height: 20),
              
              TextFormField(
                controller: _linkController,
                decoration: const InputDecoration(
                  labelText: 'Subscription Link',
                  hintText: 'vmess:// or ss:// or https://',
                  border: OutlineInputBorder(),
                ),
                maxLines: 3,
                validator: (value) {
                  if (value == null || value.isEmpty) {
                    return 'Please enter a subscription link';
                  }
                  // Basic validation for common protocols
                  if (!value.startsWith(RegExp(r'^(vmess|vless|ss|trojan|https?)://'))) {
                    return 'Please enter a valid subscription link';
                  }
                  return null;
                },
              ),
              
              const SizedBox(height: 20),
              
              const Text(
                'Supported formats:',
                style: TextStyle(
                  fontSize: 16,
                  fontWeight: FontWeight.bold,
                ),
              ),
              
              const SizedBox(height: 10),
              
              const Text('• VMess (vmess://)'),
              const Text('• VLESS (vless://)'),
              const Text('• Shadowsocks (ss://)'),
              const Text('• Trojan (trojan://)'),
              const Text('• HTTPS Subscription Links'),
              
              const SizedBox(height: 30),
              
              Center(
                child: ElevatedButton(
                  onPressed: _isProcessing ? null : _importSubscription,
                  style: ElevatedButton.styleFrom(
                    padding: const EdgeInsets.symmetric(horizontal: 50, vertical: 15),
                  ),
                  child: _isProcessing
                    ? const SizedBox(
                        width: 20,
                        height: 20,
                        child: CircularProgressIndicator(
                          strokeWidth: 2,
                          valueColor: AlwaysStoppedAnimation<Color>(Colors.white),
                        ),
                      )
                    : const Text(
                        'Import Servers',
                        style: TextStyle(fontSize: 16),
                      ),
                ),
              ),
              
              const SizedBox(height: 20),
              
              const Text(
                'Note:',
                style: TextStyle(
                  fontSize: 16,
                  fontWeight: FontWeight.bold,
                ),
              ),
              
              const SizedBox(height: 10),
              
              const Text(
                '• Imported servers will be added to your server list\n'
                '• Existing servers with the same configuration will be updated\n'
                '• You can manage imported servers like any other server',
                style: TextStyle(
                  color: Colors.grey,
                ),
              ),
            ],
          ),
        ),
      ),
    );
  }
  
  void _importSubscription() async {
    if (_formKey.currentState!.validate()) {
      setState(() {
        _isProcessing = true;
      });
      
      try {
        final appProvider = Provider.of<AppProvider>(context, listen: false);
        final link = _linkController.text.trim();
        
        // In a real app, you would parse the subscription link
        // For now, we'll simulate the process
        await Future.delayed(const Duration(seconds: 2));
        
        // Parse the subscription link
        final servers = await appProvider.parseSubscriptionLink(link);
        
        // Add servers
        for (final server in servers) {
          await appProvider.addServer(server);
        }
        
        if (mounted) {
          setState(() {
            _isProcessing = false;
          });
          
          ScaffoldMessenger.of(context).showSnackBar(
            SnackBar(
              content: Text('Successfully imported ${servers.length} servers'),
              duration: const Duration(seconds: 3),
            ),
          );
          
          // Clear the form
          _linkController.clear();
        }
      } catch (e) {
        if (mounted) {
          setState(() {
            _isProcessing = false;
          });
          
          ScaffoldMessenger.of(context).showSnackBar(
            SnackBar(
              content: Text('Error importing subscription: $e'),
              duration: const Duration(seconds: 3),
            ),
          );
        }
      }
    }
  }
}