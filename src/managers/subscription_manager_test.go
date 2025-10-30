package managers

import (
	"c:/Users/behza/OneDrive/Documents/vpn/src/core"
	"net/http"
	"net/http/httptest"
	"testing"
)

func TestSubscriptionManager(t *testing.T) {
	// Create a mock HTTP server for testing
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		// Return sample subscription content
		content := "vmess://eyJ2IjoiMiIsInBzIjoiVGVzdCBTZXJ2ZXIiLCJhZGQiOiJleGFtcGxlLmNvbSIsInBvcnQiOiI0NDMiLCJpZCI6InRlc3QtdXVpZCIsImFpZCI6IjAiLCJzY3kiOiJhdXRvIiwibmV0IjoidGNwIiwidHlwZSI6Im5vbmUiLCJob3N0IjoiIiwicGF0aCI6IiIsInRscyI6InRscyJ9"
		w.Write([]byte(content))
	}))
	defer server.Close()

	// Create server manager and subscription manager
	serverManager := NewServerManager()
	subManager := NewSubscriptionManager(serverManager)

	// Test parsing subscription from URL
	servers, err := subManager.parseSubscription(server.URL)
	if err != nil {
		t.Errorf("Failed to parse subscription: %v", err)
		return
	}

	// Check that we got the expected number of servers
	if len(servers) != 1 {
		t.Errorf("Expected 1 server, got %d", len(servers))
		return
	}

	// Check server details
	server := servers[0]
	if server.Protocol != core.ProtocolVMess {
		t.Errorf("Expected protocol VMess, got %s", server.Protocol)
	}

	if server.Host != "example.com" {
		t.Errorf("Expected host example.com, got %s", server.Host)
	}

	if server.Name != "Test Server" {
		t.Errorf("Expected name 'Test Server', got %s", server.Name)
	}
}