package managers

import (
	"crypto/rand"
	"encoding/base64"
	"encoding/json"
	"fmt"
	"net/url"
	"strings"

	"vpnclient/src/core"
)

// SubscriptionParser parses subscription links and extracts server information
type SubscriptionParser struct{}

// NewSubscriptionParser creates a new subscription parser
func NewSubscriptionParser() *SubscriptionParser {
	return &SubscriptionParser{}
}

// Parse parses a subscription link and returns a list of servers
func (sp *SubscriptionParser) Parse(subLink string) ([]*core.Server, error) {
	link := strings.TrimSpace(subLink)
	
	// Handle different protocol formats
	if strings.HasPrefix(link, "vmess://") {
		server, err := sp.parseVMessLink(link)
		if err != nil {
			return nil, err
		}
		return []*core.Server{server}, nil
	} else if strings.HasPrefix(link, "ss://") {
		server, err := sp.parseShadowsocksLink(link)
		if err != nil {
			return nil, err
		}
		return []*core.Server{server}, nil
	} else if strings.HasPrefix(link, "trojan://") {
		server, err := sp.parseTrojanLink(link)
		if err != nil {
			return nil, err
		}
		return []*core.Server{server}, nil
	} else {
		// Assume it's a base64-encoded subscription link
		return sp.parseBase64Subscription(link)
	}
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
	if strings.Contains(base64Part, "#") {
		// Format: ss://base64(method:password)@server:port#name
		parts := strings.SplitN(base64Part, "#", 2)
		base64Part = parts[0]
		// The name part is URL encoded
		name, _ := url.QueryUnescape(parts[1])
		
		// Try to decode base64 part
		decoded, err := base64.StdEncoding.DecodeString(base64Part)
		if err != nil {
			// If base64 decode fails, assume the part before @ is not base64 encoded
			// Try to parse as plain text
			atIndex := strings.LastIndex(base64Part, "@")
			if atIndex == -1 {
				return nil, fmt.Errorf("invalid Shadowsocks link format")
			}
			
			creds := base64Part[:atIndex]
			serverInfo := base64Part[atIndex+1:]
			
			// Try to decode credentials part
			decodedCreds, err := base64.StdEncoding.DecodeString(creds)
			if err != nil {
				// If decode fails, use as plain text
				decodedCreds = []byte(creds)
			}
			
			// Split credentials
			colonIndex := strings.Index(string(decodedCreds), ":")
			if colonIndex == -1 {
				return nil, fmt.Errorf("invalid Shadowsocks credentials format")
			}
			
			method := string(decodedCreds)[:colonIndex]
			password := string(decodedCreds)[colonIndex+1:]
			
			// Split server info
			colonIndex = strings.LastIndex(serverInfo, ":")
			if colonIndex == -1 {
				return nil, fmt.Errorf("invalid Shadowsocks server info format")
			}
			
			host := serverInfo[:colonIndex]
			var port int
			fmt.Sscanf(serverInfo[colonIndex+1:], "%d", &port)
			
			server := &core.Server{
				ID:       generateID(),
				Name:     name,
				Host:     host,
				Port:     port,
				Protocol: "shadowsocks",
				Config: map[string]interface{}{
					"method":   method,
					"password": password,
				},
				Enabled: true,
			}
			
			// Set server name if empty
			if server.Name == "" {
				server.Name = fmt.Sprintf("Shadowsocks %s:%d", server.Host, server.Port)
			}
			
			return server, nil
		} else {
			// Successfully decoded base64
			decodedStr := string(decoded)
			
			// Parse the decoded string
			atIndex := strings.LastIndex(decodedStr, "@")
			if atIndex == -1 {
				return nil, fmt.Errorf("invalid Shadowsocks link format")
			}
			
			creds := decodedStr[:atIndex]
			serverInfo := decodedStr[atIndex+1:]
			
			// Split credentials
			colonIndex := strings.Index(creds, ":")
			if colonIndex == -1 {
				return nil, fmt.Errorf("invalid Shadowsocks credentials format")
			}
			
			method := creds[:colonIndex]
			password := creds[colonIndex+1:]
			
			// Split server info
			colonIndex = strings.LastIndex(serverInfo, ":")
			if colonIndex == -1 {
				return nil, fmt.Errorf("invalid Shadowsocks server info format")
			}
			
			host := serverInfo[:colonIndex]
			var port int
			fmt.Sscanf(serverInfo[colonIndex+1:], "%d", &port)
			
			server := &core.Server{
				ID:       generateID(),
				Name:     name,
				Host:     host,
				Port:     port,
				Protocol: "shadowsocks",
				Config: map[string]interface{}{
					"method":   method,
					"password": password,
				},
				Enabled: true,
			}
			
			// Set server name if empty
			if server.Name == "" {
				server.Name = fmt.Sprintf("Shadowsocks %s:%d", server.Host, server.Port)
			}
			
			return server, nil
		}
	} else {
		// No name part, try to decode base64 part
		decoded, err := base64.StdEncoding.DecodeString(base64Part)
		if err != nil {
			// Try to parse as plain text
			atIndex := strings.LastIndex(base64Part, "@")
			if atIndex == -1 {
				return nil, fmt.Errorf("invalid Shadowsocks link format")
			}
			
			creds := base64Part[:atIndex]
			serverInfo := base64Part[atIndex+1:]
			
			// Try to decode credentials part
			decodedCreds, err := base64.StdEncoding.DecodeString(creds)
			if err != nil {
				// If decode fails, use as plain text
				decodedCreds = []byte(creds)
			}
			
			// Split credentials
			colonIndex := strings.Index(string(decodedCreds), ":")
			if colonIndex == -1 {
				return nil, fmt.Errorf("invalid Shadowsocks credentials format")
			}
			
			method := string(decodedCreds)[:colonIndex]
			password := string(decodedCreds)[colonIndex+1:]
			
			// Split server info
			colonIndex = strings.LastIndex(serverInfo, ":")
			if colonIndex == -1 {
				return nil, fmt.Errorf("invalid Shadowsocks server info format")
			}
			
			host := serverInfo[:colonIndex]
			var port int
			fmt.Sscanf(serverInfo[colonIndex+1:], "%d", &port)
			
			server := &core.Server{
				ID:       generateID(),
				Name:     fmt.Sprintf("Shadowsocks %s:%d", host, port),
				Host:     host,
				Port:     port,
				Protocol: "shadowsocks",
				Config: map[string]interface{}{
					"method":   method,
					"password": password,
				},
				Enabled: true,
			}
			
			return server, nil
		}
		
		decodedStr := string(decoded)
		
		// Parse the decoded string
		atIndex := strings.LastIndex(decodedStr, "@")
		if atIndex == -1 {
			return nil, fmt.Errorf("invalid Shadowsocks link format")
		}
		
		creds := decodedStr[:atIndex]
		serverInfo := decodedStr[atIndex+1:]
		
		// Split credentials
		colonIndex := strings.Index(creds, ":")
		if colonIndex == -1 {
			return nil, fmt.Errorf("invalid Shadowsocks credentials format")
		}
		
		method := creds[:colonIndex]
		password := creds[colonIndex+1:]
		
		// Split server info
		colonIndex = strings.LastIndex(serverInfo, ":")
		if colonIndex == -1 {
			return nil, fmt.Errorf("invalid Shadowsocks server info format")
		}
		
		host := serverInfo[:colonIndex]
		var port int
		fmt.Sscanf(serverInfo[colonIndex+1:], "%d", &port)
		
		server := &core.Server{
			ID:       generateID(),
			Name:     fmt.Sprintf("Shadowsocks %s:%d", host, port),
			Host:     host,
			Port:     port,
			Protocol: "shadowsocks",
			Config: map[string]interface{}{
				"method":   method,
				"password": password,
			},
			Enabled: true,
		}
		
		return server, nil
	}
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
	// Decode base64 subscription data
	decodedData, err := base64.StdEncoding.DecodeString(link)
	if err != nil {
		return nil, fmt.Errorf("failed to decode subscription data: %v", err)
	}

	// Split lines and parse each server
	lines := strings.Split(string(decodedData), "\n")
	var servers []*core.Server

	for _, line := range lines {
		line = strings.TrimSpace(line)
		if line == "" {
			continue
		}

		// Parse different protocol formats
		if strings.HasPrefix(line, "vmess://") {
			server, err := sp.parseVMessLink(line)
			if err != nil {
				// Log error but continue parsing other lines
				continue
			}
			servers = append(servers, server)
		} else if strings.HasPrefix(line, "ss://") {
			server, err := sp.parseShadowsocksLink(line)
			if err != nil {
				// Log error but continue parsing other lines
				continue
			}
			servers = append(servers, server)
		} else if strings.HasPrefix(line, "trojan://") {
			server, err := sp.parseTrojanLink(line)
			if err != nil {
				// Log error but continue parsing other lines
				continue
			}
			servers = append(servers, server)
		}
		// Add support for other protocols as needed
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

func getIntValue(m map[string]interface{}, key string, defaultValue int) int {
	if val, ok := m[key]; ok {
		switch v := val.(type) {
		case float64:
			return int(v)
		case int:
			return v
		case string:
			var i int
			fmt.Sscanf(v, "%d", &i)
			return i
		}
	}
	return defaultValue
}

func generateID() string {
	// Generate a random ID
	bytes := make([]byte, 16)
	_, err := rand.Read(bytes)
	if err != nil {
		// Fallback to timestamp-based ID if random generation fails
		return fmt.Sprintf("srv_%d", len(fmt.Sprintf("%v", &bytes)))
	}
	return base64.URLEncoding.EncodeToString(bytes)
}