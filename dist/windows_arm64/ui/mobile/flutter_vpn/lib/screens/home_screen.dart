import 'package:flutter/material.dart';
import 'package:provider/provider.dart';
import 'package:vpn_client/providers/app_provider.dart';
import 'package:vpn_client/screens/settings_screen.dart';
import 'package:vpn_client/widgets/server_list.dart';
import 'package:vpn_client/widgets/connection_status.dart';
import 'package:vpn_client/widgets/quick_connect.dart';

class HomeScreen extends StatefulWidget {
  const HomeScreen({Key? key}) : super(key: key);

  @override
  State<HomeScreen> createState() => _HomeScreenState();
}

class _HomeScreenState extends State<HomeScreen> {
  int _selectedIndex = 0;

  static const List<Widget> _widgetOptions = <Widget>[
    ServerList(),
    ConnectionStatusWidget(),
    QuickConnect(),
  ];

  void _onItemTapped(int index) {
    setState(() {
      _selectedIndex = index;
    });
  }

  @override
  Widget build(BuildContext context) {
    return Consumer<AppProvider>(
      builder: (context, appProvider, child) {
        return Scaffold(
          appBar: AppBar(
            title: const Text('VPN Client'),
            backgroundColor: Colors.blue,
            foregroundColor: Colors.white,
            actions: [
              IconButton(
                icon: const Icon(Icons.settings),
                onPressed: () {
                  // Navigate to settings screen
                  Navigator.push(
                    context,
                    MaterialPageRoute(builder: (context) => const SettingsScreen()),
                  );
                },
              ),
            ],
          ),
          body: _widgetOptions.elementAt(_selectedIndex),
          bottomNavigationBar: BottomNavigationBar(
            items: const <BottomNavigationBarItem>[
              BottomNavigationBarItem(
                icon: Icon(Icons.list),
                label: 'Servers',
              ),
              BottomNavigationBarItem(
                icon: Icon(Icons.network_check),
                label: 'Status',
              ),
              BottomNavigationBarItem(
                icon: Icon(Icons.flash_on),
                label: 'Quick Connect',
              ),
            ],
            currentIndex: _selectedIndex,
            selectedItemColor: Colors.blue,
            onTap: _onItemTapped,
          ),
        );
      },
    );
  }
}