package managers

import (
	"encoding/base64"
	"encoding/json"
	"testing"
)

// TestSubscriptionParserVMess tests parsing of VMess subscription links
func TestSubscriptionParserVMess(t *testing.T) {
	parser := NewSubscriptionParser()
	
	// Create a sample VMess configuration
	vmessConfig := map[string]interface{}{
		"v":    "2",
		"ps":   "Test VMess Server",
		"add":  "example.com",
		"port": "443",
		"id":   "12345678-1234-1234-1234-123456789012",
		"aid":  "0",
		"net":  "ws",
		"type": "none",
		"host": "example.com",
		"path": "/v2ray",
		"tls":  "tls",
	}
	
	// Convert to JSON
	jsonData, err := json.Marshal(vmessConfig)
	if err != nil {
		t.Fatalf("Failed to marshal VMess config: %v", err)
	}
	
	// Encode as base64
	encoded := base64.StdEncoding.EncodeToString(jsonData)
	vmessLink := "vmess://" + encoded
	
	// Parse the link
	servers, err := parser.Parse(vmessLink)
	if err != nil {
		t.Fatalf("Failed to parse VMess link: %v", err)
	}
	
	// Verify results
	if len(servers) != 1 {
		t.Fatalf("Expected 1 server, got %d", len(servers))
	}
	
	server := servers[0]
	if server.Name != "Test VMess Server" {
		t.Errorf("Expected server name 'Test VMess Server', got '%s'", server.Name)
	}
	
	if server.Host != "example.com" {
		t.Errorf("Expected host 'example.com', got '%s'", server.Host)
	}
	
	if server.Port != 443 {
		t.Errorf("Expected port 443, got %d", server.Port)
	}
	
	if server.Protocol != "vmess" {
		t.Errorf("Expected protocol 'vmess', got '%s'", server.Protocol)
	}
}

// TestSubscriptionParserShadowsocks tests parsing of Shadowsocks subscription links
func TestSubscriptionParserShadowsocks(t *testing.T) {
	parser := NewSubscriptionParser()
	
	// Create a sample Shadowsocks link with plain text format
	// Format: ss://method:password@server:port#name
	ssLink := "ss://aes-128-gcm:testpassword@ss.example.com:8388#Test%20SS%20Server"
	
	// Parse the link
	servers, err := parser.Parse(ssLink)
	if err != nil {
		t.Fatalf("Failed to parse Shadowsocks link: %v", err)
	}
	
	// Verify results
	if len(servers) != 1 {
		t.Fatalf("Expected 1 server, got %d", len(servers))
	}
	
	server := servers[0]
	if server.Name != "Test SS Server" {
		t.Errorf("Expected server name 'Test SS Server', got '%s'", server.Name)
	}
	
	if server.Host != "ss.example.com" {
		t.Errorf("Expected host 'ss.example.com', got '%s'", server.Host)
	}
	
	if server.Port != 8388 {
		t.Errorf("Expected port 8388, got %d", server.Port)
	}
	
	if server.Protocol != "shadowsocks" {
		t.Errorf("Expected protocol 'shadowsocks', got '%s'", server.Protocol)
	}
	
	// Test base64 encoded format
	ssLink2 := "ss://YWVzLTEyOC1nY206dGVzdHBhc3N3b3JkQHNzMi5leGFtcGxlLmNvbTo4Mzg4#Test%20SS%20Server%202"
	
	// Parse the link
	servers2, err := parser.Parse(ssLink2)
	if err != nil {
		t.Fatalf("Failed to parse Shadowsocks link: %v", err)
	}
	
	// Verify results
	if len(servers2) != 1 {
		t.Fatalf("Expected 1 server, got %d", len(servers2))
	}
	
	server2 := servers2[0]
	if server2.Name != "Test SS Server 2" {
		t.Errorf("Expected server name 'Test SS Server 2', got '%s'", server2.Name)
	}
	
	if server2.Host != "ss2.example.com" {
		t.Errorf("Expected host 'ss2.example.com', got '%s'", server2.Host)
	}
	
	if server2.Port != 8388 {
		t.Errorf("Expected port 8388, got %d", server2.Port)
	}
	
	if server2.Protocol != "shadowsocks" {
		t.Errorf("Expected protocol 'shadowsocks', got '%s'", server2.Protocol)
	}
}

// TestSubscriptionParserTrojan tests parsing of Trojan subscription links
func TestSubscriptionParserTrojan(t *testing.T) {
	parser := NewSubscriptionParser()
	
	// Create a sample Trojan link
	trojanLink := "trojan://password@trojan.example.com:443#Test%20Trojan%20Server"
	
	// Parse the link
	servers, err := parser.Parse(trojanLink)
	if err != nil {
		t.Fatalf("Failed to parse Trojan link: %v", err)
	}
	
	// Verify results
	if len(servers) != 1 {
		t.Fatalf("Expected 1 server, got %d", len(servers))
	}
	
	server := servers[0]
	if server.Name != "Test Trojan Server" {
		t.Errorf("Expected server name 'Test Trojan Server', got '%s'", server.Name)
	}
	
	if server.Host != "trojan.example.com" {
		t.Errorf("Expected host 'trojan.example.com', got '%s'", server.Host)
	}
	
	if server.Port != 443 {
		t.Errorf("Expected port 443, got %d", server.Port)
	}
	
	if server.Protocol != "trojan" {
		t.Errorf("Expected protocol 'trojan', got '%s'", server.Protocol)
	}
}

// TestSubscriptionParserBase64Subscription tests parsing of base64-encoded subscription
func TestSubscriptionParserBase64Subscription(t *testing.T) {
	parser := NewSubscriptionParser()
	
	// Create sample server configurations
	vmessConfig := map[string]interface{}{
		"v":    "2",
		"ps":   "Test VMess Server",
		"add":  "vmess.example.com",
		"port": "443",
		"id":   "12345678-1234-1234-1234-123456789012",
		"aid":  "0",
		"net":  "ws",
		"type": "none",
		"host": "vmess.example.com",
		"path": "/v2ray",
		"tls":  "tls",
	}
	
	// Convert to JSON
	jsonData, err := json.Marshal(vmessConfig)
	if err != nil {
		t.Fatalf("Failed to marshal VMess config: %v", err)
	}
	
	// Encode as base64
	vmessEncoded := base64.StdEncoding.EncodeToString(jsonData)
	
	// Create a multi-line subscription with different protocols
	subscriptionContent := "vmess://" + vmessEncoded + "\n" +
		"ss://YWVzLTEyOC1nY206dGVzdHBhc3N3b3JkQHNzMi5leGFtcGxlLmNvbTo4Mzg4\n" +
		"trojan://password@trojan2.example.com:443\n"
	
	// Encode the entire subscription as base64
	subscriptionEncoded := base64.StdEncoding.EncodeToString([]byte(subscriptionContent))
	
	// Parse the subscription
	servers, err := parser.Parse(subscriptionEncoded)
	if err != nil {
		t.Fatalf("Failed to parse base64 subscription: %v", err)
	}
	
	// Verify results
	if len(servers) != 3 {
		t.Fatalf("Expected 3 servers, got %d", len(servers))
	}
	
	// Check VMess server
	vmessServer := servers[0]
	if vmessServer.Protocol != "vmess" {
		t.Errorf("Expected first server to be VMess, got '%s'", vmessServer.Protocol)
	}
	
	// Check Shadowsocks server
	ssServer := servers[1]
	if ssServer.Protocol != "shadowsocks" {
		t.Errorf("Expected second server to be Shadowsocks, got '%s'", ssServer.Protocol)
	}
	
	// Check Trojan server
	trojanServer := servers[2]
	if trojanServer.Protocol != "trojan" {
		t.Errorf("Expected third server to be Trojan, got '%s'", trojanServer.Protocol)
	}
}

// TestSubscriptionParserInvalidLinks tests parsing of invalid subscription links
func TestSubscriptionParserInvalidLinks(t *testing.T) {
	parser := NewSubscriptionParser()
	
	// Test invalid base64
	_, err := parser.Parse("invalid_base64")
	if err == nil {
		t.Error("Expected error when parsing invalid base64, got nil")
	}
	
	// Test unsupported protocol
	_, err = parser.Parse("unsupported://data")
	if err == nil {
		t.Error("Expected error when parsing unsupported protocol, got nil")
	}
	
	// Test empty link - this should not cause an error as it might be a valid empty subscription
	_, err = parser.Parse("")
	if err != nil {
		t.Errorf("Did not expect error when parsing empty link, got %v", err)
	}
}