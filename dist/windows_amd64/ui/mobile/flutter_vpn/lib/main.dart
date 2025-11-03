import 'package:flutter/material.dart';
import 'package:provider/provider.dart';
import 'package:vpn_client/providers/app_provider.dart';
import 'package:vpn_client/screens/home_screen.dart';
import 'package:vpn_client/screens/settings_screen.dart';

void main() {
  runApp(
    ChangeNotifierProvider(
      create: (context) => AppProvider(),
      child: const VPNClientApp(),
    ),
  );
}

class VPNClientApp extends StatelessWidget {
  const VPNClientApp({Key? key}) : super(key: key);

  @override
  Widget build(BuildContext context) {
    return MaterialApp(
      title: 'VPN Client',
      theme: ThemeData(
        primarySwatch: Colors.blue,
        visualDensity: VisualDensity.adaptivePlatformDensity,
      ),
      home: const HomeScreen(),
      routes: {
        '/settings': (context) => const SettingsScreen(),
      },
      debugShowCheckedModeBanner: false,
    );
  }
}