import React, { useState, useEffect } from 'react';
import {
  StyleSheet,
  View,
  Text,
  TouchableOpacity,
  ScrollView,
  Alert,
  RefreshControl,
  TextInput,
  Switch,
  StatusBar,
  Platform,
} from 'react-native';
import AsyncStorage from '@react-native-async-storage/async-storage';
import PushNotification from 'react-native-push-notification';

const API_BASE_URL = 'http://localhost:8080/api';

// Configure push notifications
PushNotification.configure({
  onRegister: function (token) {
    console.log("TOKEN:", token);
  },
  onNotification: function (notification) {
    console.log("NOTIFICATION:", notification);
  },
  permissions: {
    alert: true,
    badge: true,
    sound: true,
  },
  popInitialNotification: true,
  requestPermissions: Platform.OS === 'ios',
});

const App = () => {
  const [servers, setServers] = useState([]);
  const [subscriptions, setSubscriptions] = useState([]);
  const [connectionStatus, setConnectionStatus] = useState('disconnected');
  const [currentServer, setCurrentServer] = useState(null);
  const [stats, setStats] = useState({ data_sent: 0, data_recv: 0, time: 0 });
  const [loading, setLoading] = useState(false);
  const [theme, setTheme] = useState('dark');
  const [language, setLanguage] = useState('en');
  const [activeTab, setActiveTab] = useState('dashboard');
  const [subscriptionUrl, setSubscriptionUrl] = useState('');
  const [showAddSubscription, setShowAddSubscription] = useState(false);

  // Translations
  const translations = {
    en: {
      vpnClient: 'VPN Client',
      dashboard: 'Dashboard',
      servers: 'Servers',
      subscriptions: 'Subscriptions',
      settings: 'Settings',
      connect: 'CONNECT',
      disconnect: 'DISCONNECT',
      connecting: 'CONNECTING...',
      connected: 'CONNECTED',
      disconnected: 'DISCONNECTED',
      download: 'Download',
      upload: 'Upload',
      connectionTime: 'Connection Time',
      currentSession: 'Current Session',
      refresh: 'Refresh',
      testAll: 'Test All',
      addSubscription: 'Add Subscription',
      subscriptionUrl: 'Subscription URL',
      enterSubscriptionLink: 'Enter subscription link',
      cancel: 'Cancel',
      submit: 'Submit',
      name: 'Name',
      url: 'URL',
      serversCount: 'Servers',
      lastUpdate: 'Last Update',
      actions: 'Actions',
      update: 'Update',
      delete: 'Delete',
      test: 'Test',
      ping: 'Ping',
      protocol: 'Protocol',
      status: 'Status',
      enabled: 'Enabled',
      disabled: 'Disabled',
      bestServer: 'Best Server',
      error: 'Error',
      success: 'Success',
      confirmDelete: 'Are you sure you want to delete this subscription?',
      subscriptionAdded: 'Subscription added successfully',
      subscriptionAddFailed: 'Failed to add subscription',
      subscriptionDeleted: 'Subscription deleted successfully',
      subscriptionDeleteFailed: 'Failed to delete subscription',
      serverTestComplete: 'Server testing complete',
      connectionSuccess: 'Successfully connected to the server',
      connectionFailed: 'Failed to connect to the server',
      disconnectionSuccess: 'Successfully disconnected from the server',
      disconnectionFailed: 'Failed to disconnect from the server',
      notifications: 'Notifications',
      enableNotifications: 'Enable Notifications',
      notificationTitle: 'VPN Client Notification',
    },
    fa: {
      vpnClient: 'کلاینت VPN',
      dashboard: 'داشبورد',
      servers: 'سرورها',
      subscriptions: 'اشتراک‌ها',
      settings: 'تنظیمات',
      connect: 'اتصال',
      disconnect: 'قطع اتصال',
      connecting: 'در حال اتصال...',
      connected: 'متصل شد',
      disconnected: 'قطع شده',
      download: 'دانلود',
      upload: 'آپلود',
      connectionTime: 'زمان اتصال',
      currentSession: 'جلسه فعلی',
      refresh: 'تازه‌سازی',
      testAll: 'تست همه',
      addSubscription: 'افزودن اشتراک',
      subscriptionUrl: 'آدرس اشتراک',
      enterSubscriptionLink: 'لینک اشتراک را وارد کنید',
      cancel: 'لغو',
      submit: 'ثبت',
      name: 'نام',
      url: 'آدرس',
      serversCount: 'سرورها',
      lastUpdate: 'آخرین به‌روزرسانی',
      actions: 'عملیات',
      update: 'بروزرسانی',
      delete: 'حذف',
      test: 'تست',
      ping: 'پینگ',
      protocol: 'پروتکل',
      status: 'وضعیت',
      enabled: 'فعال',
      disabled: 'غیرفعال',
      bestServer: 'بهترین سرور',
      error: 'خطا',
      success: 'موفقیت',
      confirmDelete: 'آیا مطمئن هستید که می‌خواهید این اشتراک را حذف کنید؟',
      subscriptionAdded: 'اشتراک با موفقیت اضافه شد',
      subscriptionAddFailed: 'افزودن اشتراک ناموفق بود',
      subscriptionDeleted: 'اشتراک با موفقیت حذف شد',
      subscriptionDeleteFailed: 'حذف اشتراک ناموفق بود',
      serverTestComplete: 'تست سرورها تکمیل شد',
      connectionSuccess: 'با موفقیت به سرور متصل شد',
      connectionFailed: 'اتصال به سرور ناموفق بود',
      disconnectionSuccess: 'اتصال از سرور با موفقیت قطع شد',
      disconnectionFailed: 'قطع اتصال از سرور ناموفق بود',
      notifications: 'اعلان‌ها',
      enableNotifications: 'فعال کردن اعلان‌ها',
      notificationTitle: 'اعلان کلاینت VPN',
    },
    zh: {
      vpnClient: 'VPN客户端',
      dashboard: '仪表板',
      servers: '服务器',
      subscriptions: '订阅',
      settings: '设置',
      connect: '连接',
      disconnect: '断开连接',
      connecting: '连接中...',
      connected: '已连接',
      disconnected: '已断开连接',
      download: '下载',
      upload: '上传',
      connectionTime: '连接时间',
      currentSession: '当前会话',
      refresh: '刷新',
      testAll: '全部测试',
      addSubscription: '添加订阅',
      subscriptionUrl: '订阅URL',
      enterSubscriptionLink: '输入订阅链接',
      cancel: '取消',
      submit: '提交',
      name: '名称',
      url: '网址',
      serversCount: '服务器',
      lastUpdate: '最后更新',
      actions: '操作',
      update: '更新',
      delete: '删除',
      test: '测试',
      ping: '延迟',
      protocol: '协议',
      status: '状态',
      enabled: '启用',
      disabled: '禁用',
      bestServer: '最佳服务器',
      error: '错误',
      success: '成功',
      confirmDelete: '您确定要删除此订阅吗？',
      subscriptionAdded: '订阅添加成功',
      subscriptionAddFailed: '订阅添加失败',
      subscriptionDeleted: '订阅删除成功',
      subscriptionDeleteFailed: '订阅删除失败',
      serverTestComplete: '服务器测试完成',
      connectionSuccess: '成功连接到服务器',
      connectionFailed: '连接服务器失败',
      disconnectionSuccess: '成功断开服务器连接',
      disconnectionFailed: '断开服务器连接失败',
      notifications: '通知',
      enableNotifications: '启用通知',
      notificationTitle: 'VPN客户端通知',
    }
  };

  const t = (key) => {
    return translations[language][key] || key;
  };

  // Load settings on app start
  useEffect(() => {
    loadSettings();
    configurePushNotifications();
  }, []);

  // Configure push notifications
  const configurePushNotifications = () => {
    PushNotification.createChannel(
      {
        channelId: "vpn-notifications",
        channelName: "VPN Notifications",
        channelDescription: "Notifications for VPN connection status",
        playSound: true,
        soundName: "default",
        importance: 4,
        vibrate: true,
      },
      (created) => console.log(`createChannel returned '${created}'`)
    );
  };

  // Show local notification
  const showLocalNotification = (title, message) => {
    PushNotification.localNotification({
      channelId: "vpn-notifications",
      title: title,
      message: message,
      playSound: true,
      soundName: "default",
      actions: '["OK"]',
    });
  };

  // Load settings from AsyncStorage
  const loadSettings = async () => {
    try {
      const savedTheme = await AsyncStorage.getItem('theme');
      const savedLanguage = await AsyncStorage.getItem('language');
      const savedNotifications = await AsyncStorage.getItem('notifications');
      
      if (savedTheme) setTheme(savedTheme);
      if (savedLanguage) setLanguage(savedLanguage);
      if (savedNotifications !== null) {
        setNotificationsEnabled(savedNotifications === 'true');
      }
    } catch (error) {
      console.error('Error loading settings:', error);
    }
  };

  // Save settings to AsyncStorage
  const saveSettings = async () => {
    try {
      await AsyncStorage.setItem('theme', theme);
      await AsyncStorage.setItem('language', language);
      await AsyncStorage.setItem('notifications', notificationsEnabled.toString());
    } catch (error) {
      console.error('Error saving settings:', error);
    }
  };

  // Save settings when they change
  useEffect(() => {
    saveSettings();
  }, [theme, language, notificationsEnabled]);

  // Load data based on active tab
  useEffect(() => {
    switch (activeTab) {
      case 'dashboard':
        loadDashboardData();
        break;
      case 'servers':
        loadAllServers();
        break;
      case 'subscriptions':
        loadSubscriptions();
        break;
    }
  }, [activeTab]);

  // Load dashboard data
  const loadDashboardData = async () => {
    setLoading(true);
    try {
      const response = await fetch(`${API_BASE_URL}/servers/enabled`);
      const data = await response.json();
      setServers(data);
    } catch (error) {
      console.error('Error loading dashboard data:', error);
      Alert.alert(t('error'), 'Failed to load server data');
    } finally {
      setLoading(false);
    }
  };

  // Load all servers
  const loadAllServers = async () => {
    setLoading(true);
    try {
      const response = await fetch(`${API_BASE_URL}/servers`);
      const data = await response.json();
      setServers(data);
    } catch (error) {
      console.error('Error loading servers:', error);
      Alert.alert(t('error'), 'Failed to load servers');
    } finally {
      setLoading(false);
    }
  };

  // Load subscriptions
  const loadSubscriptions = async () => {
    setLoading(true);
    try {
      const response = await fetch(`${API_BASE_URL}/subscriptions`);
      const data = await response.json();
      setSubscriptions(data);
    } catch (error) {
      console.error('Error loading subscriptions:', error);
      Alert.alert(t('error'), 'Failed to load subscriptions');
    } finally {
      setLoading(false);
    }
  };

  // Connect to a server
  const connectToServer = async (serverId) => {
    try {
      setConnectionStatus('connecting');
      const response = await fetch(`${API_BASE_URL}/connect`, {
        method: 'POST',
        headers: {
          'Content-Type': 'application/json',
        },
        body: JSON.stringify({ server_id: serverId }),
      });

      if (response.ok) {
        setConnectionStatus('connected');
        const server = servers.find(s => s.id === serverId);
        setCurrentServer(server);
        if (notificationsEnabled) {
          showLocalNotification(t('notificationTitle'), t('connectionSuccess'));
        }
        Alert.alert(t('success'), t('connectionSuccess'));
        // Start stats updates
        updateStats();
      } else {
        setConnectionStatus('disconnected');
        if (notificationsEnabled) {
          showLocalNotification(t('notificationTitle'), t('connectionFailed'));
        }
        Alert.alert(t('error'), t('connectionFailed'));
      }
    } catch (error) {
      setConnectionStatus('disconnected');
      console.error('Error connecting to server:', error);
      if (notificationsEnabled) {
        showLocalNotification(t('notificationTitle'), t('connectionFailed'));
      }
      Alert.alert(t('error'), t('connectionFailed'));
    }
  };

  // Connect to best server
  const connectToBestServer = async () => {
    try {
      setConnectionStatus('connecting');
      const response = await fetch(`${API_BASE_URL}/connect/best`, {
        method: 'POST',
      });

      if (response.ok) {
        setConnectionStatus('connected');
        if (notificationsEnabled) {
          showLocalNotification(t('notificationTitle'), t('connectionSuccess'));
        }
        Alert.alert(t('success'), t('connectionSuccess'));
        // Start stats updates
        updateStats();
      } else {
        setConnectionStatus('disconnected');
        if (notificationsEnabled) {
          showLocalNotification(t('notificationTitle'), t('connectionFailed'));
        }
        Alert.alert(t('error'), t('connectionFailed'));
      }
    } catch (error) {
      setConnectionStatus('disconnected');
      console.error('Error connecting to best server:', error);
      if (notificationsEnabled) {
        showLocalNotification(t('notificationTitle'), t('connectionFailed'));
      }
      Alert.alert(t('error'), t('connectionFailed'));
    }
  };

  // Disconnect from current server
  const disconnect = async () => {
    try {
      const response = await fetch(`${API_BASE_URL}/disconnect`, {
        method: 'POST',
      });

      if (response.ok) {
        setConnectionStatus('disconnected');
        setCurrentServer(null);
        setStats({ data_sent: 0, data_recv: 0, time: 0 });
        if (notificationsEnabled) {
          showLocalNotification(t('notificationTitle'), t('disconnectionSuccess'));
        }
        Alert.alert(t('success'), t('disconnectionSuccess'));
      } else {
        if (notificationsEnabled) {
          showLocalNotification(t('notificationTitle'), t('disconnectionFailed'));
        }
        Alert.alert(t('error'), t('disconnectionFailed'));
      }
    } catch (error) {
      console.error('Error disconnecting:', error);
      if (notificationsEnabled) {
        showLocalNotification(t('notificationTitle'), t('disconnectionFailed'));
      }
      Alert.alert(t('error'), t('disconnectionFailed'));
    }
  };

  // Update connection stats
  const updateStats = async () => {
    if (connectionStatus !== 'connected') return;

    try {
      const response = await fetch(`${API_BASE_URL}/stats`);
      if (response.ok) {
        const data = await response.json();
        setStats({
          data_sent: data.data_sent || 0,
          data_recv: data.data_recv || 0,
          time: data.started_at ? Math.floor((Date.now() - new Date(data.started_at).getTime()) / 1000) : 0
        });
      }
    } catch (error) {
      console.error('Error updating stats:', error);
    }

    // Schedule next update
    setTimeout(updateStats, 1000);
  };

  // Format bytes to human readable format
  const formatBytes = (bytes) => {
    if (bytes === 0) return '0 Bytes';
    const k = 1024;
    const sizes = ['Bytes', 'KB', 'MB', 'GB'];
    const i = Math.floor(Math.log(bytes) / Math.log(k));
    return parseFloat((bytes / Math.pow(k, i)).toFixed(2)) + ' ' + sizes[i];
  };

  // Format time to HH:MM:SS
  const formatTime = (seconds) => {
    const h = Math.floor(seconds / 3600);
    const m = Math.floor((seconds % 3600) / 60);
    const s = seconds % 60;
    return `${h.toString().padStart(2, '0')}:${m.toString().padStart(2, '0')}:${s.toString().padStart(2, '0')}`;
  };

  // Add subscription
  const addSubscription = async () => {
    if (!subscriptionUrl.trim()) {
      Alert.alert(t('error'), t('enterSubscriptionLink'));
      return;
    }

    try {
      const response = await fetch(`${API_BASE_URL}/subscriptions`, {
        method: 'POST',
        headers: {
          'Content-Type': 'application/json',
        },
        body: JSON.stringify({ url: subscriptionUrl }),
      });

      if (response.ok) {
        setSubscriptionUrl('');
        setShowAddSubscription(false);
        loadSubscriptions();
        if (notificationsEnabled) {
          showLocalNotification(t('notificationTitle'), t('subscriptionAdded'));
        }
        Alert.alert(t('success'), t('subscriptionAdded'));
      } else {
        const error = await response.json();
        if (notificationsEnabled) {
          showLocalNotification(t('notificationTitle'), t('subscriptionAddFailed'));
        }
        Alert.alert(t('error'), `${t('subscriptionAddFailed')}: ${error.error}`);
      }
    } catch (error) {
      console.error('Error adding subscription:', error);
      if (notificationsEnabled) {
        showLocalNotification(t('notificationTitle'), t('subscriptionAddFailed'));
      }
      Alert.alert(t('error'), t('subscriptionAddFailed'));
    }
  };

  // Delete subscription
  const deleteSubscription = async (id) => {
    Alert.alert(
      t('confirmDelete'),
      '',
      [
        {
          text: t('cancel'),
          style: 'cancel',
        },
        {
          text: t('delete'),
          style: 'destructive',
          onPress: async () => {
            try {
              const response = await fetch(`${API_BASE_URL}/subscriptions/${id}`, {
                method: 'DELETE',
              });

              if (response.ok) {
                loadSubscriptions();
                if (notificationsEnabled) {
                  showLocalNotification(t('notificationTitle'), t('subscriptionDeleted'));
                }
                Alert.alert(t('success'), t('subscriptionDeleted'));
              } else {
                const error = await response.json();
                if (notificationsEnabled) {
                  showLocalNotification(t('notificationTitle'), t('subscriptionDeleteFailed'));
                }
                Alert.alert(t('error'), `${t('subscriptionDeleteFailed')}: ${error.error}`);
              }
            } catch (error) {
              console.error('Error deleting subscription:', error);
              if (notificationsEnabled) {
                showLocalNotification(t('notificationTitle'), t('subscriptionDeleteFailed'));
              }
              Alert.alert(t('error'), t('subscriptionDeleteFailed'));
            }
          },
        },
      ]
    );
  };

  // Test all servers
  const testAllServers = async () => {
    try {
      const response = await fetch(`${API_BASE_URL}/servers/test-all-ping`, {
        method: 'POST',
      });

      if (response.ok) {
        loadAllServers();
        if (notificationsEnabled) {
          showLocalNotification(t('notificationTitle'), t('serverTestComplete'));
        }
        Alert.alert(t('success'), t('serverTestComplete'));
      } else {
        const error = await response.json();
        Alert.alert(t('error'), `Error: ${error.error}`);
      }
    } catch (error) {
      console.error('Error testing servers:', error);
      Alert.alert(t('error'), 'Failed to test servers');
    }
  };

  // Get ping class for styling
  const getPingClass = (ping) => {
    if (ping <= 0) return styles.pingBad;
    if (ping < 50) return styles.pingGood;
    if (ping < 150) return styles.pingMedium;
    return styles.pingBad;
  };

  // State for notification settings
  const [notificationsEnabled, setNotificationsEnabled] = useState(true);

  // Toggle notifications
  const toggleNotifications = async (value) => {
    setNotificationsEnabled(value);
    await AsyncStorage.setItem('notifications', value.toString());
    
    if (value) {
      // Request notification permission
      PushNotification.requestPermissions();
    }
  };

  // Render dashboard tab
  const renderDashboard = () => (
    <ScrollView 
      style={styles.container}
      refreshControl={
        <RefreshControl refreshing={loading} onRefresh={loadDashboardData} />
      }
    >
      {/* Stats Cards */}
      <View style={styles.statsContainer}>
        <View style={[styles.card, theme === 'dark' ? styles.cardDark : styles.cardLight]}>
          <Text style={styles.statTitle}>{t('download')}</Text>
          <Text style={styles.statValue}>{formatBytes(stats.data_recv)}</Text>
          <Text style={styles.statLabel}>{t('currentSession')}</Text>
        </View>

        <View style={[styles.card, theme === 'dark' ? styles.cardDark : styles.cardLight]}>
          <Text style={styles.statTitle}>{t('upload')}</Text>
          <Text style={styles.statValue}>{formatBytes(stats.data_sent)}</Text>
          <Text style={styles.statLabel}>{t('currentSession')}</Text>
        </View>

        <View style={[styles.card, theme === 'dark' ? styles.cardDark : styles.cardLight]}>
          <Text style={styles.statTitle}>{t('connectionTime')}</Text>
          <Text style={styles.statValue}>{formatTime(stats.time)}</Text>
          <Text style={styles.statLabel}>{t('currentSession')}</Text>
        </View>
      </View>

      {/* Connection Panel */}
      <View style={[styles.card, theme === 'dark' ? styles.cardDark : styles.cardLight, styles.connectionPanel]}>
        <View style={[
          styles.connectionStatus,
          connectionStatus === 'connected' ? styles.statusConnected : 
          connectionStatus === 'connecting' ? styles.statusConnecting : styles.statusDisconnected
        ]}>
          <Text style={styles.connectionStatusText}>
            {connectionStatus === 'connected' ? t('connected') :
             connectionStatus === 'connecting' ? t('connecting') : t('disconnected')}
          </Text>
        </View>

        <Text style={styles.connectionTitle}>{t('secureConnection')}</Text>
        <Text style={styles.connectionSubtitle}>{t('protectPrivacy')}</Text>

        <View style={styles.connectionButtons}>
          {connectionStatus === 'disconnected' ? (
            <TouchableOpacity 
              style={[styles.connectButton, styles.primaryButton]} 
              onPress={connectToBestServer}
            >
              <Text style={styles.buttonText}>{t('connect')}</Text>
            </TouchableOpacity>
          ) : (
            <TouchableOpacity 
              style={[styles.connectButton, styles.dangerButton]} 
              onPress={disconnect}
            >
              <Text style={styles.buttonText}>{t('disconnect')}</Text>
            </TouchableOpacity>
          )}
        </View>
      </View>

      {/* Server List */}
      <View style={styles.sectionHeader}>
        <Text style={styles.sectionTitle}>{t('servers')}</Text>
        <TouchableOpacity style={styles.refreshButton} onPress={loadDashboardData}>
          <Text style={styles.refreshButtonText}>{t('refresh')}</Text>
        </TouchableOpacity>
      </View>

      <View style={[styles.card, theme === 'dark' ? styles.cardDark : styles.cardLight]}>
        {servers.filter(s => s.enabled).map((server) => (
          <View 
            key={server.id} 
            style={[styles.serverItem, theme === 'dark' ? styles.serverItemDark : styles.serverItemLight]}
          >
            <View style={styles.serverInfo}>
              <Text style={styles.serverName}>{server.name || `${server.host}:${server.port}`}</Text>
              <Text style={styles.serverDetails}>
                {server.protocol} • {server.enabled ? t('enabled') : t('disabled')}
              </Text>
            </View>
            <View style={styles.serverActions}>
              <View style={[styles.pingIndicator, getPingClass(server.ping)]}>
                <Text style={styles.pingText}>
                  {server.ping > 0 ? `${server.ping} ms` : 'N/A'}
                </Text>
              </View>
              <TouchableOpacity 
                style={[styles.actionButton, styles.primaryButtonSmall]}
                onPress={() => connectToServer(server.id)}
                disabled={connectionStatus === 'connecting'}
              >
                <Text style={styles.buttonTextSmall}>{t('connect')}</Text>
              </TouchableOpacity>
            </View>
          </View>
        ))}
      </View>
    </ScrollView>
  );

  // Render servers tab
  const renderServers = () => (
    <ScrollView 
      style={styles.container}
      refreshControl={
        <RefreshControl refreshing={loading} onRefresh={loadAllServers} />
      }
    >
      <View style={styles.sectionHeader}>
        <Text style={styles.sectionTitle}>{t('allServers')}</Text>
        <View style={styles.sectionActions}>
          <TouchableOpacity style={styles.refreshButton} onPress={loadAllServers}>
            <Text style={styles.refreshButtonText}>{t('refresh')}</Text>
          </TouchableOpacity>
          <TouchableOpacity style={styles.testButton} onPress={testAllServers}>
            <Text style={styles.testButtonText}>{t('testAll')}</Text>
          </TouchableOpacity>
        </View>
      </View>

      <View style={[styles.card, theme === 'dark' ? styles.cardDark : styles.cardLight]}>
        {servers.map((server) => (
          <View 
            key={server.id} 
            style={[styles.serverItem, theme === 'dark' ? styles.serverItemDark : styles.serverItemLight]}
          >
            <View style={styles.serverInfo}>
              <Text style={styles.serverName}>{server.name || `${server.host}:${server.port}`}</Text>
              <Text style={styles.serverDetails}>
                {server.protocol} • {server.enabled ? t('enabled') : t('disabled')}
              </Text>
            </View>
            <View style={styles.serverActions}>
              <View style={[styles.pingIndicator, getPingClass(server.ping)]}>
                <Text style={styles.pingText}>
                  {server.ping > 0 ? `${server.ping} ms` : 'N/A'}
                </Text>
              </View>
              <TouchableOpacity 
                style={[styles.actionButton, styles.primaryButtonSmall]}
                onPress={() => connectToServer(server.id)}
                disabled={connectionStatus === 'connecting'}
              >
                <Text style={styles.buttonTextSmall}>{t('connect')}</Text>
              </TouchableOpacity>
            </View>
          </View>
        ))}
      </View>
    </ScrollView>
  );

  // Render subscriptions tab
  const renderSubscriptions = () => (
    <ScrollView 
      style={styles.container}
      refreshControl={
        <RefreshControl refreshing={loading} onRefresh={loadSubscriptions} />
      }
    >
      <View style={styles.sectionHeader}>
        <Text style={styles.sectionTitle}>{t('manageSubscriptions')}</Text>
        <View style={styles.sectionActions}>
          <TouchableOpacity style={styles.refreshButton} onPress={loadSubscriptions}>
            <Text style={styles.refreshButtonText}>{t('refresh')}</Text>
          </TouchableOpacity>
          <TouchableOpacity 
            style={[styles.addButton, styles.primaryButton]} 
            onPress={() => setShowAddSubscription(true)}
          >
            <Text style={styles.buttonText}>{t('addSubscription')}</Text>
          </TouchableOpacity>
        </View>
      </View>

      {showAddSubscription && (
        <View style={[styles.card, theme === 'dark' ? styles.cardDark : styles.cardLight, styles.subscriptionForm]}>
          <Text style={styles.formTitle}>{t('addSubscription')}</Text>
          <TextInput
            style={[styles.input, theme === 'dark' ? styles.inputDark : styles.inputLight]}
            placeholder={t('enterSubscriptionLink')}
            placeholderTextColor={theme === 'dark' ? '#aaa' : '#666'}
            value={subscriptionUrl}
            onChangeText={setSubscriptionUrl}
          />
          <View style={styles.formActions}>
            <TouchableOpacity 
              style={[styles.formButton, styles.primaryButton]} 
              onPress={addSubscription}
            >
              <Text style={styles.buttonText}>{t('submit')}</Text>
            </TouchableOpacity>
            <TouchableOpacity 
              style={[styles.formButton, styles.secondaryButton]} 
              onPress={() => setShowAddSubscription(false)}
            >
              <Text style={styles.buttonText}>{t('cancel')}</Text>
            </TouchableOpacity>
          </View>
        </View>
      )}

      <View style={[styles.card, theme === 'dark' ? styles.cardDark : styles.cardLight]}>
        {subscriptions.map((sub) => (
          <View 
            key={sub.id} 
            style={[styles.subscriptionItem, theme === 'dark' ? styles.subscriptionItemDark : styles.subscriptionItemLight]}
          >
            <View style={styles.subscriptionInfo}>
              <Text style={styles.subscriptionName}>{sub.name || 'Unnamed Subscription'}</Text>
              <Text style={styles.subscriptionUrl} numberOfLines={1}>{sub.url}</Text>
              <Text style={styles.subscriptionDetails}>
                {sub.serverCount || 0} {t('servers')} • {sub.lastUpdate ? new Date(sub.lastUpdate).toLocaleString() : t('never')}
              </Text>
            </View>
            <View style={styles.subscriptionActions}>
              <TouchableOpacity 
                style={[styles.actionButton, styles.secondaryButtonSmall]}
                onPress={() => {}}
              >
                <Text style={styles.buttonTextSmall}>{t('update')}</Text>
              </TouchableOpacity>
              <TouchableOpacity 
                style={[styles.actionButton, styles.dangerButtonSmall]}
                onPress={() => deleteSubscription(sub.id)}
              >
                <Text style={styles.buttonTextSmall}>{t('delete')}</Text>
              </TouchableOpacity>
            </View>
          </View>
        ))}
      </View>
    </ScrollView>
  );

  // Render settings tab
  const renderSettings = () => (
    <ScrollView style={styles.container}>
      <View style={[styles.card, theme === 'dark' ? styles.cardDark : styles.cardLight]}>
        <View style={styles.settingItem}>
          <Text style={styles.settingLabel}>{t('language')}</Text>
          <View style={styles.settingControl}>
            <TouchableOpacity 
              style={[styles.languageButton, language === 'en' && styles.activeLanguageButton]}
              onPress={() => setLanguage('en')}
            >
              <Text style={[styles.languageText, language === 'en' && styles.activeLanguageText]}>EN</Text>
            </TouchableOpacity>
            <TouchableOpacity 
              style={[styles.languageButton, language === 'fa' && styles.activeLanguageButton]}
              onPress={() => setLanguage('fa')}
            >
              <Text style={[styles.languageText, language === 'fa' && styles.activeLanguageText]}>FA</Text>
            </TouchableOpacity>
            <TouchableOpacity 
              style={[styles.languageButton, language === 'zh' && styles.activeLanguageButton]}
              onPress={() => setLanguage('zh')}
            >
              <Text style={[styles.languageText, language === 'zh' && styles.activeLanguageText]}>ZH</Text>
            </TouchableOpacity>
          </View>
        </View>

        <View style={styles.settingItem}>
          <Text style={styles.settingLabel}>{t('theme')}</Text>
          <View style={styles.settingControl}>
            <TouchableOpacity 
              style={[styles.themeButton, theme === 'light' && styles.activeThemeButton]}
              onPress={() => setTheme('light')}
            >
              <Text style={[styles.themeText, theme === 'light' && styles.activeThemeText]}>{t('light')}</Text>
            </TouchableOpacity>
            <TouchableOpacity 
              style={[styles.themeButton, theme === 'dark' && styles.activeThemeButton]}
              onPress={() => setTheme('dark')}
            >
              <Text style={[styles.themeText, theme === 'dark' && styles.activeThemeText]}>{t('dark')}</Text>
            </TouchableOpacity>
          </View>
        </View>

        <View style={[styles.settingItem, styles.settingItemLast]}>
          <Text style={styles.settingLabel}>{t('notifications')}</Text>
          <Switch
            value={notificationsEnabled}
            onValueChange={toggleNotifications}
            trackColor={{ false: "#767577", true: "#81b0ff" }}
            thumbColor={notificationsEnabled ? "#f5dd4b" : "#f4f3f4"}
          />
        </View>
      </View>
    </ScrollView>
  );

  return (
    <View style={[styles.appContainer, theme === 'dark' ? styles.appContainerDark : styles.appContainerLight]}>
      <StatusBar 
        barStyle={theme === 'dark' ? 'light-content' : 'dark-content'} 
        backgroundColor={theme === 'dark' ? '#121826' : '#f8fafc'} 
      />
      
      {/* Header */}
      <View style={[styles.header, theme === 'dark' ? styles.headerDark : styles.headerLight]}>
        <Text style={styles.headerTitle}>{t('vpnClient')}</Text>
      </View>

      {/* Content */}
      {activeTab === 'dashboard' && renderDashboard()}
      {activeTab === 'servers' && renderServers()}
      {activeTab === 'subscriptions' && renderSubscriptions()}
      {activeTab === 'settings' && renderSettings()}

      {/* Navigation */}
      <View style={[styles.navigation, theme === 'dark' ? styles.navigationDark : styles.navigationLight]}>
        <TouchableOpacity 
          style={[styles.navItem, activeTab === 'dashboard' && styles.activeNavItem]}
          onPress={() => setActiveTab('dashboard')}
        >
          <Text style={[styles.navText, activeTab === 'dashboard' && styles.activeNavText]}>
            {t('dashboard')}
          </Text>
        </TouchableOpacity>
        
        <TouchableOpacity 
          style={[styles.navItem, activeTab === 'servers' && styles.activeNavItem]}
          onPress={() => setActiveTab('servers')}
        >
          <Text style={[styles.navText, activeTab === 'servers' && styles.activeNavText]}>
            {t('servers')}
          </Text>
        </TouchableOpacity>
        
        <TouchableOpacity 
          style={[styles.navItem, activeTab === 'subscriptions' && styles.activeNavItem]}
          onPress={() => setActiveTab('subscriptions')}
        >
          <Text style={[styles.navText, activeTab === 'subscriptions' && styles.activeNavText]}>
            {t('subscriptions')}
          </Text>
        </TouchableOpacity>
        
        <TouchableOpacity 
          style={[styles.navItem, activeTab === 'settings' && styles.activeNavItem]}
          onPress={() => setActiveTab('settings')}
        >
          <Text style={[styles.navText, activeTab === 'settings' && styles.activeNavText]}>
            {t('settings')}
          </Text>
        </TouchableOpacity>
      </View>
    </View>
  );
};

const styles = StyleSheet.create({
  appContainer: {
    flex: 1,
  },
  appContainerLight: {
    backgroundColor: '#f8fafc',
  },
  appContainerDark: {
    backgroundColor: '#121826',
  },
  header: {
    paddingTop: 80,
    paddingBottom: 20,
    paddingHorizontal: 20,
    borderBottomWidth: 1,
  },
  headerLight: {
    backgroundColor: '#ffffff',
    borderBottomColor: '#e2e8f0',
  },
  headerDark: {
    backgroundColor: '#1e293b',
    borderBottomColor: '#334155',
  },
  headerTitle: {
    fontSize: 24,
    fontWeight: 'bold',
  },
  container: {
    flex: 1,
    padding: 20,
  },
  statsContainer: {
    flexDirection: 'row',
    flexWrap: 'wrap',
    justifyContent: 'space-between',
    marginBottom: 20,
  },
  card: {
    borderRadius: 12,
    padding: 20,
    marginBottom: 20,
    shadowColor: '#000',
    shadowOffset: {
      width: 0,
      height: 2,
    },
    shadowOpacity: 0.1,
    shadowRadius: 3.84,
    elevation: 5,
  },
  cardLight: {
    backgroundColor: '#ffffff',
  },
  cardDark: {
    backgroundColor: '#1e293b',
  },
  statTitle: {
    fontSize: 16,
    marginBottom: 10,
  },
  statValue: {
    fontSize: 24,
    fontWeight: 'bold',
    marginVertical: 10,
  },
  statLabel: {
    fontSize: 14,
    opacity: 0.7,
  },
  connectionPanel: {
    alignItems: 'center',
  },
  connectionStatus: {
    width: 150,
    height: 150,
    borderRadius: 75,
    justifyContent: 'center',
    alignItems: 'center',
    marginBottom: 20,
  },
  statusConnected: {
    backgroundColor: 'rgba(6, 214, 160, 0.15)',
  },
  statusConnecting: {
    backgroundColor: 'rgba(255, 209, 102, 0.15)',
  },
  statusDisconnected: {
    backgroundColor: 'rgba(239, 71, 111, 0.15)',
  },
  connectionStatusText: {
    fontSize: 18,
    fontWeight: 'bold',
  },
  connectionTitle: {
    fontSize: 20,
    fontWeight: 'bold',
    marginBottom: 10,
  },
  connectionSubtitle: {
    fontSize: 16,
    marginBottom: 20,
    textAlign: 'center',
  },
  connectionButtons: {
    flexDirection: 'row',
    justifyContent: 'center',
  },
  connectButton: {
    paddingHorizontal: 30,
    paddingVertical: 15,
    borderRadius: 30,
  },
  primaryButton: {
    backgroundColor: '#4361ee',
  },
  dangerButton: {
    backgroundColor: '#ef476f',
  },
  buttonText: {
    color: 'white',
    fontWeight: 'bold',
    fontSize: 16,
  },
  sectionHeader: {
    flexDirection: 'row',
    justifyContent: 'space-between',
    alignItems: 'center',
    marginBottom: 20,
  },
  sectionTitle: {
    fontSize: 20,
    fontWeight: 'bold',
  },
  sectionActions: {
    flexDirection: 'row',
  },
  refreshButton: {
    paddingHorizontal: 15,
    paddingVertical: 8,
    borderRadius: 20,
    backgroundColor: '#4361ee',
    marginLeft: 10,
  },
  testButton: {
    paddingHorizontal: 15,
    paddingVertical: 8,
    borderRadius: 20,
    backgroundColor: '#7209b7',
    marginLeft: 10,
  },
  addButton: {
    paddingHorizontal: 15,
    paddingVertical: 8,
    borderRadius: 20,
  },
  refreshButtonText: {
    color: 'white',
    fontWeight: 'bold',
  },
  testButtonText: {
    color: 'white',
    fontWeight: 'bold',
  },
  buttonTextSmall: {
    color: 'white',
    fontWeight: 'bold',
    fontSize: 12,
  },
  serverItem: {
    flexDirection: 'row',
    justifyContent: 'space-between',
    alignItems: 'center',
    padding: 15,
    borderRadius: 8,
    marginBottom: 10,
  },
  serverItemLight: {
    backgroundColor: 'rgba(0, 0, 0, 0.05)',
  },
  serverItemDark: {
    backgroundColor: 'rgba(255, 255, 255, 0.05)',
  },
  serverInfo: {
    flex: 1,
  },
  serverName: {
    fontSize: 16,
    fontWeight: 'bold',
    marginBottom: 5,
  },
  serverDetails: {
    fontSize: 14,
    opacity: 0.7,
  },
  serverActions: {
    flexDirection: 'row',
    alignItems: 'center',
  },
  pingIndicator: {
    paddingHorizontal: 10,
    paddingVertical: 5,
    borderRadius: 20,
    marginRight: 10,
  },
  pingGood: {
    backgroundColor: 'rgba(6, 214, 160, 0.15)',
  },
  pingMedium: {
    backgroundColor: 'rgba(255, 209, 102, 0.15)',
  },
  pingBad: {
    backgroundColor: 'rgba(239, 71, 111, 0.15)',
  },
  pingText: {
    fontSize: 12,
    fontWeight: 'bold',
  },
  actionButton: {
    paddingHorizontal: 15,
    paddingVertical: 8,
    borderRadius: 20,
  },
  primaryButtonSmall: {
    backgroundColor: '#4361ee',
  },
  secondaryButtonSmall: {
    backgroundColor: '#7209b7',
  },
  dangerButtonSmall: {
    backgroundColor: '#ef476f',
  },
  subscriptionForm: {
    marginBottom: 20,
  },
  formTitle: {
    fontSize: 18,
    fontWeight: 'bold',
    marginBottom: 15,
  },
  input: {
    borderWidth: 1,
    borderRadius: 8,
    paddingHorizontal: 15,
    paddingVertical: 12,
    marginBottom: 15,
    fontSize: 16,
  },
  inputLight: {
    borderColor: '#cbd5e1',
    backgroundColor: '#f8fafc',
    color: '#1e293b',
  },
  inputDark: {
    borderColor: '#334155',
    backgroundColor: '#1e293b',
    color: '#f1f5f9',
  },
  formActions: {
    flexDirection: 'row',
    justifyContent: 'space-between',
  },
  formButton: {
    flex: 1,
    paddingHorizontal: 20,
    paddingVertical: 12,
    borderRadius: 8,
    marginHorizontal: 5,
  },
  secondaryButton: {
    backgroundColor: '#64748b',
  },
  subscriptionItem: {
    flexDirection: 'row',
    justifyContent: 'space-between',
    alignItems: 'center',
    padding: 15,
    borderRadius: 8,
    marginBottom: 10,
  },
  subscriptionItemLight: {
    backgroundColor: 'rgba(0, 0, 0, 0.05)',
  },
  subscriptionItemDark: {
    backgroundColor: 'rgba(255, 255, 255, 0.05)',
  },
  subscriptionInfo: {
    flex: 1,
  },
  subscriptionName: {
    fontSize: 16,
    fontWeight: 'bold',
    marginBottom: 5,
  },
  subscriptionUrl: {
    fontSize: 14,
    marginBottom: 5,
    opacity: 0.7,
  },
  subscriptionDetails: {
    fontSize: 12,
    opacity: 0.7,
  },
  subscriptionActions: {
    flexDirection: 'row',
  },
  settingItem: {
    flexDirection: 'row',
    justifyContent: 'space-between',
    alignItems: 'center',
    paddingVertical: 15,
    borderBottomWidth: 1,
  },
  settingItemLast: {
    borderBottomWidth: 0,
  },
  settingLabel: {
    fontSize: 16,
    fontWeight: 'bold',
  },
  settingControl: {
    flexDirection: 'row',
  },
  languageButton: {
    paddingHorizontal: 15,
    paddingVertical: 8,
    borderRadius: 20,
    marginHorizontal: 5,
  },
  activeLanguageButton: {
    backgroundColor: '#4361ee',
  },
  languageText: {
    fontWeight: 'bold',
  },
  activeLanguageText: {
    color: 'white',
  },
  themeButton: {
    paddingHorizontal: 15,
    paddingVertical: 8,
    borderRadius: 20,
    marginHorizontal: 5,
  },
  activeThemeButton: {
    backgroundColor: '#4361ee',
  },
  themeText: {
    fontWeight: 'bold',
  },
  activeThemeText: {
    color: 'white',
  },
  navigation: {
    flexDirection: 'row',
    borderTopWidth: 1,
    paddingVertical: 10,
  },
  navigationLight: {
    backgroundColor: '#ffffff',
    borderTopColor: '#e2e8f0',
  },
  navigationDark: {
    backgroundColor: '#1e293b',
    borderTopColor: '#334155',
  },
  navItem: {
    flex: 1,
    alignItems: 'center',
    paddingVertical: 10,
  },
  activeNavItem: {
    // Add active styling if needed
  },
  navText: {
    fontSize: 12,
    fontWeight: 'bold',
  },
  activeNavText: {
    color: '#4361ee',
  },
});

export default App;

