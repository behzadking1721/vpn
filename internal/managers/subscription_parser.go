package managers

import (
	"context"
	"crypto/rand"
	"encoding/base64"
	"encoding/json"
	"fmt"
	"net/url"
	"strings"
	"sync"
	"time"

	"vpnclient/src/core"
)

// SubscriptionParser parses subscription links and extracts server information
type SubscriptionParser struct {
	cache    map[string][]*core.Server
	cacheMux sync.RWMutex
	cacheTTL time.Duration
}

// NewSubscriptionParser creates a new subscription parser
func NewSubscriptionParser() *SubscriptionParser {
	return &SubscriptionParser{
		cache:    make(map[string][]*core.Server),
		cacheTTL: 10 * time.Minute, // Default cache TTL of 10 minutes
	}
}

// Parse parses a subscription link and returns a list of servers
func (sp *SubscriptionParser) Parse(subLink string) ([]*core.Server, error) {
	link := strings.TrimSpace(subLink)
	
	// Check cache first
	sp.cacheMux.RLock()
	if servers, exists := sp.cache[link]; exists {
		sp.cacheMux.RUnlock()
		return servers, nil
	}
	sp.cacheMux.RUnlock()
	
	// Handle different protocol formats
	var servers []*core.Server
	var err error
	
	if strings.HasPrefix(link, "vmess://") {
		server, parseErr := sp.parseVMessLink(link)
		if parseErr != nil {
			return nil, parseErr
		}
		servers = []*core.Server{server}
	} else if strings.HasPrefix(link, "ss://") {
		server, parseErr := sp.parseShadowsocksLink(link)
		if parseErr != nil {
			return nil, parseErr
		}
		servers = []*core.Server{server}
	} else if strings.HasPrefix(link, "trojan://") {
		server, parseErr := sp.parseTrojanLink(link)
		if parseErr != nil {
			return nil, parseErr
		}
		servers = []*core.Server{server}
	} else {
		// Assume it's a base64-encoded subscription link
		servers, err = sp.parseBase64Subscription(link)
		if err != nil {
			return nil, err
		}
	}
	
	// Update cache
	sp.cacheMux.Lock()
	sp.cache[link] = servers
	sp.cacheMux.Unlock()
	
	return servers, nil
}

// parseVMessLink parses a VMess link
func (sp *SubscriptionParser) parseVMessLink(link string) (*core.Server, error) {
	// Remove protocol prefix
	base64Str := strings.TrimPrefix(link, "vmess://")
	
	// Decode base64
	jsonStr, err := base64.StdEncoding.DecodeString(base64Str)
	if err != nil {
		return nil, fmt.Errorf("failed to decode VMess link: %v", err)
	}
	
	// Parse JSON
	var vmessConfig map[string]interface{}
	err = json.Unmarshal(jsonStr, &vmessConfig)
	if err != nil {
		return nil, fmt.Errorf("failed to parse VMess JSON: %v", err)
	}
	
	// Create server from configuration
	server := &core.Server{
		ID:        generateID(),
		Name:      getStringValue(vmessConfig, "ps", "VMess Server"),
		Host:      getStringValue(vmessConfig, "add", ""),
		Port:      getIntValue(vmessConfig, "port", 0),
		Protocol:  "vmess",
		Config: map[string]interface{}{
			"encryption": getStringValue(vmessConfig, "scy", "auto"),
			"user_id":    getStringValue(vmessConfig, "id", ""),
			"alter_id":   getIntValue(vmessConfig, "aid", 0),
			"tls":        getStringValue(vmessConfig, "tls", "") == "tls",
			"sni":        getStringValue(vmessConfig, "sni", ""),
		},
	}
	
	// Set server name if empty
	if server.Name == "" {
		server.Name = fmt.Sprintf("VMess %s:%d", server.Host, server.Port)
	}
	
	// Enable server by default
	server.Enabled = true
	
	return server, nil
}

// parseShadowsocksLink parses a Shadowsocks link
func (sp *SubscriptionParser) parseShadowsocksLink(link string) (*core.Server, error) {
	// Remove protocol prefix
	base64Part := strings.TrimPrefix(link, "ss://")
	
	// Handle different formats
	var decodedStr string
	var err error
	
	// Try to decode as base64
	decodedBytes, decodeErr := base64.URLEncoding.DecodeString(base64Part)
	if decodeErr != nil {
		decodedBytes, decodeErr = base64.StdEncoding.DecodeString(base64Part)
	}
	
	if decodeErr == nil {
		decodedStr = string(decodedBytes)
	} else {
		decodedStr = base64Part
	}
	
	// Parse the link
	var host, port, password, method, name string
	
	// Handle SIP002 format (ss://method:password@host:port#name)
	if strings.Contains(decodedStr, "@") {
		// Split method:password and host:port
		parts := strings.Split(decodedStr, "@")
		if len(parts) != 2 {
			return nil, fmt.Errorf("invalid Shadowsocks link format")
		}
		
		// Parse method:password
		authPart := parts[0]
		authParts := strings.Split(authPart, ":")
		if len(authParts) != 2 {
			return nil, fmt.Errorf("invalid Shadowsocks authentication format")
		}
		method = authParts[0]
		password = authParts[1]
		
		// Parse host:port and name
		hostPortPart := parts[1]
		if strings.Contains(hostPortPart, "#") {
			namePart := strings.Split(hostPortPart, "#")
			hostPortPart = namePart[0]
			name = namePart[1]
		}
		
		hostPortParts := strings.Split(hostPortPart, ":")
		if len(hostPortParts) != 2 {
			return nil, fmt.Errorf("invalid Shadowsocks host:port format")
		}
		host = hostPortParts[0]
		port = hostPortParts[1]
	} else {
		// Handle legacy format
		// Split by #
		parts := strings.Split(decodedStr, "#")
		if len(parts) > 1 {
			name = parts[1]
		}
		
		// Split by : for host:port
		mainPart := parts[0]
		mainParts := strings.Split(mainPart, ":")
		if len(mainParts) < 3 {
			return nil, fmt.Errorf("invalid Shadowsocks link format")
		}
		
		// Last two parts are host and port
		port = mainParts[len(mainParts)-1]
		host = mainParts[len(mainParts)-2]
		
		// Everything else is method:password
		methodPassword := strings.Join(mainParts[:len(mainParts)-2], ":")
		methodPasswordParts := strings.Split(methodPassword, ":")
		if len(methodPasswordParts) != 2 {
			return nil, fmt.Errorf("invalid Shadowsocks method:password format")
		}
		method = methodPasswordParts[0]
		password = methodPasswordParts[1]
	}
	
	// Convert port to int
	var portInt int
	fmt.Sscanf(port, "%d", &portInt)
	
	// Create server
	server := &core.Server{
		ID:       generateID(),
		Name:     name,
		Host:     host,
		Port:     portInt,
		Protocol: "shadowsocks",
		Config: map[string]interface{}{
			"method":   method,
			"password": password,
		},
		Enabled: true,
	}
	
	// Set default name if empty
	if server.Name == "" {
		server.Name = fmt.Sprintf("Shadowsocks %s:%d", host, portInt)
	}
	
	return server, nil
}

// parseTrojanLink parses a Trojan link
func (sp *SubscriptionParser) parseTrojanLink(link string) (*core.Server, error) {
	// Remove protocol prefix
	link = strings.TrimPrefix(link, "trojan://")
	
	// Parse URL
	u, err := url.Parse("trojan://" + link)
	if err != nil {
		return nil, fmt.Errorf("failed to parse Trojan link: %v", err)
	}
	
	// Extract password (username)
	password := u.User.Username()
	
	// Extract host and port
	host := u.Hostname()
	port := 443 // Default port
	if u.Port() != "" {
		fmt.Sscanf(u.Port(), "%d", &port)
	}
	
	// Extract name from fragment
	name := "Trojan Server"
	if u.Fragment != "" {
		decodedName, _ := url.QueryUnescape(u.Fragment)
		if decodedName != "" {
			name = decodedName
		}
	}
	
	server := &core.Server{
		ID:       generateID(),
		Name:     name,
		Host:     host,
		Port:     port,
		Protocol: "trojan",
		Config: map[string]interface{}{
			"password": password,
		},
		Enabled: true,
	}
	
	return server, nil
}

// parseBase64Subscription parses a base64-encoded subscription
func (sp *SubscriptionParser) parseBase64Subscription(link string) ([]*core.Server, error) {
	// Decode base64
	decodedBytes, err := base64.StdEncoding.DecodeString(link)
	if err != nil {
		// Try URL encoding
		decodedBytes, err = base64.URLEncoding.DecodeString(link)
		if err != nil {
			return nil, fmt.Errorf("failed to decode subscription link: %v", err)
		}
	}
	
	// Split lines
	links := strings.Split(string(decodedBytes), "\n")
	
	// Parse each link
	var servers []*core.Server
	for _, subLink := range links {
		subLink = strings.TrimSpace(subLink)
		if subLink == "" {
			continue
		}
		
		// Parse the link
		parser := NewSubscriptionParser()
		serverList, err := parser.Parse(subLink)
		if err != nil {
			// Skip invalid links
			continue
		}
		
		servers = append(servers, serverList...)
	}
	
	return servers, nil
}

// Helper functions
func getStringValue(m map[string]interface{}, key, defaultValue string) string {
	if val, ok := m[key]; ok {
		if str, ok := val.(string); ok {
			return str
		}
	}
	return defaultValue
}

// getIntValue gets an int value from a map with a default
func getIntValue(m map[string]interface{}, key string, defaultValue int) int {
	if val, ok := m[key]; ok {
		if num, ok := val.(float64); ok {
			return int(num)
		}
		if num, ok := val.(int); ok {
			return num
		}
	}
	return defaultValue
}

// generateID generates a random ID
func generateID() string {
	bytes := make([]byte, 16)
	if _, err := rand.Read(bytes); err != nil {
		// Fallback to timestamp if random generation fails
		return fmt.Sprintf("server_%d", time.Now().UnixNano())
	}
	
	// Convert to hex string
	return fmt.Sprintf("%x", bytes)
}