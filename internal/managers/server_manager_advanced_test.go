package managers

import (
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"

	"vpnclient/internal/database"
	"vpnclient/internal/logging"
	"vpnclient/internal/notifications"
	"vpnclient/mocks"
	"vpnclient/src/core"
)

// TestServerManagerWithNotifications tests server manager with notifications integration
func TestServerManagerWithNotifications(t *testing.T) {
	// Setup
	mockStore := &mocks.MockStore{}
	sm := NewServerManager(mockStore)
	notificationManager := notifications.NewNotificationManager(10)
	sm.SetNotificationManager(notificationManager)

	// Mock expectations
	server := &core.Server{
		ID:       "test-server",
		Name:     "Test Server",
		Host:     "example.com",
		Port:     8080,
		Protocol: "vmess",
		Enabled:  true,
	}
	
	mockStore.On("Add", mock.Anything, mock.AnythingOfType("*core.Server")).Return(nil)
	mockStore.On("GetAll", mock.Anything).Return([]interface{}{server}, nil)

	// Test adding server
	err := sm.AddServer(server)
	assert.NoError(t, err)

	// Check notifications
	notifications := notificationManager.GetNotifications()
	assert.NotEmpty(t, notifications)
	
	// Look for server added notification
	found := false
	for _, notif := range notifications {
		if notif.Title == "Server Added" {
			found = true
			assert.Equal(t, "New server Test Server added successfully", notif.Message)
			assert.Equal(t, notifications.Success, notif.Type)
			break
		}
	}
	assert.True(t, found)
}

// TestServerManagerWithLogger tests server manager with logger integration
func TestServerManagerWithLogger(t *testing.T) {
	// Setup
	mockStore := &mocks.MockStore{}
	sm := NewServerManager(mockStore)
	logger := &logging.Logger{}
	sm.SetLogger(logger)

	// Mock expectations
	server := &core.Server{
		ID:       "test-server",
		Name:     "Test Server",
		Host:     "example.com",
		Port:     8080,
		Protocol: "vmess",
		Enabled:  true,
	}
	
	mockStore.On("Add", mock.Anything, mock.AnythingOfType("*core.Server")).Return(nil)
	mockStore.On("Get", mock.Anything, "test-server").Return(server, nil)
	mockStore.On("Update", mock.Anything, mock.AnythingOfType("*core.Server")).Return(nil)
	mockStore.On("Delete", mock.Anything, "test-server").Return(nil)

	// Test adding server
	err := sm.AddServer(server)
	assert.NoError(t, err)

	// Test getting server
	retrievedServer, err := sm.GetServer("test-server")
	assert.NoError(t, err)
	assert.Equal(t, server, retrievedServer)

	// Test updating server
	server.Name = "Updated Server"
	err = sm.UpdateServer(server)
	assert.NoError(t, err)

	// Test deleting server
	err = sm.DeleteServer("test-server")
	assert.NoError(t, err)
}

// TestServerManagerCache tests server manager caching functionality
func TestServerManagerCache(t *testing.T) {
	// Setup
	mockStore := &mocks.MockStore{}
	sm := NewServerManager(mockStore)

	// Mock expectations
	server := &core.Server{
		ID:       "test-server",
		Name:     "Test Server",
		Host:     "example.com",
		Port:     8080,
		Protocol: "vmess",
		Enabled:  true,
	}
	
	// First call to GetServer should hit the store
	mockStore.On("Get", mock.Anything, "test-server").Return(server, nil).Once()

	// Test getting server (should populate cache)
	retrievedServer, err := sm.GetServer("test-server")
	assert.NoError(t, err)
	assert.Equal(t, server, retrievedServer)

	// Second call to GetServer should use cache (no additional mock calls)
	retrievedServer, err = sm.GetServer("test-server")
	assert.NoError(t, err)
	assert.Equal(t, server, retrievedServer)

	// Updating server should update cache
	server.Name = "Updated Server"
	mockStore.On("Update", mock.Anything, mock.AnythingOfType("*core.Server")).Return(nil)
	err = sm.UpdateServer(server)
	assert.NoError(t, err)

	// Getting server again should return updated version from cache
	retrievedServer, err = sm.GetServer("test-server")
	assert.NoError(t, err)
	assert.Equal(t, "Updated Server", retrievedServer.Name)

	// Deleting server should remove it from cache
	mockStore.On("Delete", mock.Anything, "test-server").Return(nil)
	err = sm.DeleteServer("test-server")
	assert.NoError(t, err)
}

// TestServerManagerGetBestServer tests getting the best server
func TestServerManagerGetBestServer(t *testing.T) {
	// Setup
	mockStore := &mocks.MockStore{}
	sm := NewServerManager(mockStore)

	// Mock expectations
	servers := []interface{}{
		&core.Server{
			ID:       "server-1",
			Name:     "Server 1",
			Host:     "example1.com",
			Port:     8080,
			Protocol: "vmess",
			Enabled:  true,
			Ping:     100,
		},
		&core.Server{
			ID:       "server-2",
			Name:     "Server 2",
			Host:     "example2.com",
			Port:     8080,
			Protocol: "vmess",
			Enabled:  true,
			Ping:     50, // Best ping
		},
		&core.Server{
			ID:       "server-3",
			Name:     "Server 3",
			Host:     "example3.com",
			Port:     8080,
			Protocol: "vmess",
			Enabled:  true,
			Ping:     75,
		},
		&core.Server{
			ID:       "server-4",
			Name:     "Server 4",
			Host:     "example4.com",
			Port:     8080,
			Protocol: "vmess",
			Enabled:  false, // Disabled server
			Ping:     25,
		},
	}
	
	mockStore.On("GetAll", mock.Anything).Return(servers, nil)

	// Test getting best server
	bestServer, err := sm.GetBestServer()
	assert.NoError(t, err)
	assert.Equal(t, "server-2", bestServer.ID) // Server 2 has the best (lowest) ping
	assert.Equal(t, 50, bestServer.Ping)
}

// TestServerManagerGetBestServerNoEnabledServers tests getting best server when no servers are enabled
func TestServerManagerGetBestServerNoEnabledServers(t *testing.T) {
	// Setup
	mockStore := &mocks.MockStore{}
	sm := NewServerManager(mockStore)

	// Mock expectations - all servers disabled
	servers := []interface{}{
		&core.Server{
			ID:       "server-1",
			Name:     "Server 1",
			Host:     "example1.com",
			Port:     8080,
			Protocol: "vmess",
			Enabled:  false,
		},
		&core.Server{
			ID:       "server-2",
			Name:     "Server 2",
			Host:     "example2.com",
			Port:     8080,
			Protocol: "vmess",
			Enabled:  false,
		},
	}
	
	mockStore.On("GetAll", mock.Anything).Return(servers, nil)

	// Test getting best server
	bestServer, err := sm.GetBestServer()
	assert.Error(t, err)
	assert.Nil(t, bestServer)
	assert.Contains(t, err.Error(), "no enabled servers available")
}

// TestServerManagerTestAllServersPing tests testing ping for all servers
func TestServerManagerTestAllServersPing(t *testing.T) {
	// Setup
	mockStore := &mocks.MockStore{}
	sm := NewServerManager(mockStore)

	// Mock expectations
	servers := []interface{}{
		&core.Server{
			ID:       "server-1",
			Name:     "Server 1",
			Host:     "127.0.0.1", // Localhost should respond
			Port:     8080,
			Protocol: "vmess",
			Enabled:  true,
		},
		&core.Server{
			ID:       "server-2",
			Name:     "Server 2",
			Host:     "invalid.host.that.does.not.exist", // This should fail
			Port:     8080,
			Protocol: "vmess",
			Enabled:  true,
		},
	}
	
	mockStore.On("GetAll", mock.Anything).Return(servers, nil)
	mockStore.On("Update", mock.Anything, mock.AnythingOfType("*core.Server")).Return(nil)

	// Test testing all servers ping
	err := sm.TestAllServersPing()
	assert.NoError(t, err) // Should not return error even if some pings fail
}

// TestServerManagerValidateServer tests server validation
func TestServerManagerValidateServer(t *testing.T) {
	// Setup
	mockStore := &mocks.MockStore{}
	sm := NewServerManager(mockStore)

	// Test server with missing host
	invalidServer1 := &core.Server{
		ID:       "invalid-server-1",
		Name:     "Invalid Server 1",
		Host:     "", // Missing host
		Port:     8080,
		Protocol: "vmess",
		Enabled:  true,
	}

	err := sm.AddServer(invalidServer1)
	assert.Error(t, err)
	assert.Contains(t, err.Error(), "server host is required")

	// Test server with invalid port
	invalidServer2 := &core.Server{
		ID:       "invalid-server-2",
		Name:     "Invalid Server 2",
		Host:     "example.com",
		Port:     99999, // Invalid port
		Protocol: "vmess",
		Enabled:  true,
	}

	err = sm.AddServer(invalidServer2)
	assert.Error(t, err)
	assert.Contains(t, err.Error(), "server port must be between 1 and 65535")

	// Test server with missing protocol
	invalidServer3 := &core.Server{
		ID:       "invalid-server-3",
		Name:     "Invalid Server 3",
		Host:     "example.com",
		Port:     8080,
		Protocol: "", // Missing protocol
		Enabled:  true,
	}

	err = sm.AddServer(invalidServer3)
	assert.Error(t, err)
	assert.Contains(t, err.Error(), "server protocol is required")

	// Test server with invalid protocol
	invalidServer4 := &core.Server{
		ID:       "invalid-server-4",
		Name:     "Invalid Server 4",
		Host:     "example.com",
		Port:     8080,
		Protocol: "invalid-protocol", // Invalid protocol
		Enabled:  true,
	}

	err = sm.AddServer(invalidServer4)
	assert.Error(t, err)
	assert.Contains(t, err.Error(), "invalid protocol: invalid-protocol")
}

// TestServerManagerEnableDisable tests enabling and disabling servers
func TestServerManagerEnableDisable(t *testing.T) {
	// Setup
	mockStore := &mocks.MockStore{}
	sm := NewServerManager(mockStore)

	// Mock expectations
	server := &core.Server{
		ID:       "test-server",
		Name:     "Test Server",
		Host:     "example.com",
		Port:     8080,
		Protocol: "vmess",
		Enabled:  true,
	}
	
	mockStore.On("Get", mock.Anything, "test-server").Return(server, nil)
	mockStore.On("Update", mock.Anything, mock.AnythingOfType("*core.Server")).Return(nil)

	// Test that server is initially enabled
	assert.True(t, server.Enabled)

	// Test disabling server
	err := sm.DisableServer("test-server")
	assert.NoError(t, err)
	assert.False(t, server.Enabled)

	// Test enabling server
	err = sm.EnableServer("test-server")
	assert.NoError(t, err)
	assert.True(t, server.Enabled)
}