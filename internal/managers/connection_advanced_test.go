package managers

import (
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"

	"vpnclient/internal/logging"
	"vpnclient/internal/notifications"
	"vpnclient/internal/stats"
	"vpnclient/mocks"
	"vpnclient/src/core"
)

// TestConnectionManagerWithStats tests connection manager with stats integration
func TestConnectionManagerWithStats(t *testing.T) {
	// Setup
	cm := NewConnectionManager()
	statsManager := stats.NewStatsManager()
	cm.SetStatsManager(statsManager)

	// Test connecting
	server := &core.Server{
		ID:       "test-server",
		Name:     "Test Server",
		Host:     "example.com",
		Port:     8080,
		Protocol: "vmess",
	}

	err := cm.Connect(server)
	assert.NoError(t, err)
	assert.Equal(t, Connected, cm.GetStatus())

	// Test updating stats
	cm.UpdateStats(1024, 2048)
	
	// Check current connection stats
	currentStats := statsManager.GetCurrentConnection()
	assert.NotNil(t, currentStats)
	assert.Equal(t, int64(1024), currentStats.DataSent)
	assert.Equal(t, int64(2048), currentStats.DataRecv)

	// Test disconnecting
	err = cm.Disconnect()
	assert.NoError(t, err)
	assert.Equal(t, Disconnected, cm.GetStatus())

	// Check that connection was ended in stats
	currentStats = statsManager.GetCurrentConnection()
	assert.Nil(t, currentStats)
}

// TestConnectionManagerWithLogger tests connection manager with logger integration
func TestConnectionManagerWithLogger(t *testing.T) {
	// Setup
	cm := NewConnectionManager()
	
	// Create a mock logger or use a real one for this test
	logger := &logging.Logger{}
	cm.SetLogger(logger)

	// Test connecting
	server := &core.Server{
		ID:       "test-server",
		Name:     "Test Server",
		Host:     "example.com",
		Port:     8080,
		Protocol: "vmess",
	}

	err := cm.Connect(server)
	assert.NoError(t, err)
	assert.Equal(t, Connected, cm.GetStatus())

	// Test disconnecting
	err = cm.Disconnect()
	assert.NoError(t, err)
	assert.Equal(t, Disconnected, cm.GetStatus())
}

// TestConnectionManagerWithNotifications tests connection manager with notifications integration
func TestConnectionManagerWithNotifications(t *testing.T) {
	// Setup
	cm := NewConnectionManager()
	notificationManager := notifications.NewNotificationManager(10)
	cm.SetNotificationManager(notificationManager)

	// Test connecting
	server := &core.Server{
		ID:       "test-server",
		Name:     "Test Server",
		Host:     "example.com",
		Port:     8080,
		Protocol: "vmess",
	}

	err := cm.Connect(server)
	assert.NoError(t, err)
	assert.Equal(t, Connected, cm.GetStatus())

	// Check notifications
	notifications := notificationManager.GetNotifications()
	assert.NotEmpty(t, notifications)
	
	// Look for connection notification
	found := false
	for _, notif := range notifications {
		if notif.Title == "Connected" {
			found = true
			assert.Equal(t, "Successfully connected to Test Server", notif.Message)
			assert.Equal(t, notifications.Success, notif.Type)
			break
		}
	}
	assert.True(t, found)

	// Test disconnecting
	err = cm.Disconnect()
	assert.NoError(t, err)
	assert.Equal(t, Disconnected, cm.GetStatus())

	// Check for disconnection notification
	notifications = notificationManager.GetNotifications()
	found = false
	for _, notif := range notifications {
		if notif.Title == "Disconnected" {
			found = true
			assert.Equal(t, "Successfully disconnected from Test Server", notif.Message)
			assert.Equal(t, notifications.Success, notif.Type)
			break
		}
	}
	assert.True(t, found)
}

// TestConnectionManagerConcurrentConnections tests that concurrent connections are handled properly
func TestConnectionManagerConcurrentConnections(t *testing.T) {
	// Setup
	cm := NewConnectionManager()

	server1 := &core.Server{
		ID:       "server-1",
		Name:     "Server 1",
		Host:     "example1.com",
		Port:     8080,
		Protocol: "vmess",
	}

	server2 := &core.Server{
		ID:       "server-2",
		Name:     "Server 2",
		Host:     "example2.com",
		Port:     8080,
		Protocol: "shadowsocks",
	}

	// Connect to first server
	err := cm.Connect(server1)
	assert.NoError(t, err)
	assert.Equal(t, Connected, cm.GetStatus())
	assert.Equal(t, server1, cm.GetCurrentServer())

	// Try to connect to second server while already connected
	err = cm.Connect(server2)
	assert.Error(t, err)
	assert.Contains(t, err.Error(), "already connected or connecting to a server")
	assert.Equal(t, server1, cm.GetCurrentServer()) // Should still be connected to first server
}

// TestConnectionManagerDisconnectWhenNotConnected tests disconnecting when not connected
func TestConnectionManagerDisconnectWhenNotConnected(t *testing.T) {
	// Setup
	cm := NewConnectionManager()
	assert.Equal(t, Disconnected, cm.GetStatus())

	// Try to disconnect when not connected
	err := cm.Disconnect()
	assert.Error(t, err)
	assert.Contains(t, err.Error(), "not connected to any server")
	assert.Equal(t, Disconnected, cm.GetStatus())
}

// TestConnectionManagerConnectionInfo tests getting connection info
func TestConnectionManagerConnectionInfo(t *testing.T) {
	// Setup
	cm := NewConnectionManager()

	// Test getting connection info when disconnected
	info := cm.GetConnectionInfo()
	assert.Nil(t, info)

	// Connect to a server
	server := &core.Server{
		ID:       "test-server",
		Name:     "Test Server",
		Host:     "example.com",
		Port:     8080,
		Protocol: "vmess",
	}

	err := cm.Connect(server)
	assert.NoError(t, err)
	assert.Equal(t, Connected, cm.GetStatus())

	// Test getting connection info when connected
	info = cm.GetConnectionInfo()
	assert.NotNil(t, info)
	assert.Equal(t, server.ID, info.ID)
	assert.WithinDuration(t, time.Now(), info.StartedAt, time.Second)
	assert.Equal(t, int64(0), info.DataSent)
	assert.Equal(t, int64(0), info.DataRecv)

	// Update stats and check again
	cm.UpdateStats(1000, 2000)
	info = cm.GetConnectionInfo()
	assert.Equal(t, int64(1000), info.DataSent)
	assert.Equal(t, int64(2000), info.DataRecv)

	// Disconnect and check again
	err = cm.Disconnect()
	assert.NoError(t, err)
	info = cm.GetConnectionInfo()
	assert.Nil(t, info)
}

// TestConnectionManagerUnsupportedProtocol tests connecting with unsupported protocol
func TestConnectionManagerUnsupportedProtocol(t *testing.T) {
	// Setup
	cm := NewConnectionManager()

	server := &core.Server{
		ID:       "test-server",
		Name:     "Test Server",
		Host:     "example.com",
		Port:     8080,
		Protocol: "unsupported-protocol",
	}

	// Try to connect with unsupported protocol
	err := cm.Connect(server)
	assert.Error(t, err)
	assert.Contains(t, err.Error(), "unsupported protocol: unsupported-protocol")
	assert.Equal(t, Error, cm.GetStatus())
}