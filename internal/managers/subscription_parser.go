package managers

import (
	"encoding/base64"
	"encoding/json"
	"fmt"
	"math/rand"
	"net/url"
	"strconv"
	"strings"
	"time"

	"vpnclient/src/core"
)

// SubscriptionParser handles parsing of subscription links
type SubscriptionParser struct{}

// NewSubscriptionParser creates a new subscription parser
func NewSubscriptionParser() *SubscriptionParser {
	return &SubscriptionParser{}
}

// ParseSubscription parses a subscription URL and returns a list of servers
func (sp *SubscriptionParser) ParseSubscription(subscriptionURL string) ([]*core.Server, error) {
	// Parse the URL
	parsedURL, err := url.Parse(subscriptionURL)
	if err != nil {
		return nil, fmt.Errorf("failed to parse subscription URL: %w", err)
	}

	// Handle different subscription types
	switch parsedURL.Scheme {
	case "http", "https":
		return sp.parseHTTPSubscription(subscriptionURL)
	case "vmess":
		server, err := sp.parseVMessLink(subscriptionURL)
		if err != nil {
			return nil, err
		}
		return []*core.Server{server}, nil
	case "ss":
		server, err := sp.parseShadowsocksLink(subscriptionURL)
		if err != nil {
			return nil, err
		}
		return []*core.Server{server}, nil
	default:
		return nil, fmt.Errorf("unsupported subscription scheme: %s", parsedURL.Scheme)
	}
}

// parseHTTPSubscription parses an HTTP(S) subscription
func (sp *SubscriptionParser) parseHTTPSubscription(subscriptionURL string) ([]*core.Server, error) {
	// In a real implementation, we would fetch the content from the URL
	// For now, we'll return an empty list as a placeholder
	return []*core.Server{}, nil
}

// parseVMessLink parses a vmess:// link
func (sp *SubscriptionParser) parseVMessLink(link string) (*core.Server, error) {
	// Remove the scheme
	data := strings.TrimPrefix(link, "vmess://")
	
	// Decode base64
	decoded, err := base64.StdEncoding.DecodeString(data)
	if err != nil {
		return nil, fmt.Errorf("failed to decode vmess link: %w", err)
	}
	
	// Parse JSON
	var config map[string]interface{}
	if err := json.Unmarshal(decoded, &config); err != nil {
		return nil, fmt.Errorf("failed to parse vmess config: %w", err)
	}
	
	// Create server
	server := &core.Server{
		ID:        generateID(),
		Name:      getString(config, "ps", "VMess Server"),
		Host:      getString(config, "add", ""),
		Port:      getInt(config, "port", 0),
		Protocol:  "vmess",
		Enabled:   true,
		CreatedAt: time.Now(),
		UpdatedAt: time.Now(),
	}
	
	return server, nil
}

// parseShadowsocksLink parses an ss:// link
func (sp *SubscriptionParser) parseShadowsocksLink(link string) (*core.Server, error) {
	// Remove the scheme
	data := strings.TrimPrefix(link, "ss://")
	
	// Split by # to get the name
	parts := strings.Split(data, "#")
	configPart := parts[0]
	name := "Shadowsocks Server"
	if len(parts) > 1 {
		name = parts[1]
	}
	
	// Decode base64
	decoded, err := base64.StdEncoding.DecodeString(configPart)
	if err != nil {
		// Try without base64 decoding
		decoded = []byte(configPart)
	}
	
	// Parse the config part (method:password@host:port)
	configStr := string(decoded)
	atIndex := strings.LastIndex(configStr, "@")
	if atIndex == -1 {
		return nil, fmt.Errorf("invalid shadowsocks link format")
	}
	
	hostPort := configStr[atIndex+1:]
	hostPortParts := strings.Split(hostPort, ":")
	if len(hostPortParts) != 2 {
		return nil, fmt.Errorf("invalid host:port format")
	}
	
	port, err := strconv.Atoi(hostPortParts[1])
	if err != nil {
		return nil, fmt.Errorf("invalid port: %w", err)
	}
	
	server := &core.Server{
		ID:        generateID(),
		Name:      name,
		Host:      hostPortParts[0],
		Port:      port,
		Protocol:  "shadowsocks",
		Enabled:   true,
		CreatedAt: time.Now(),
		UpdatedAt: time.Now(),
	}
	
	return server, nil
}

// getString gets a string value from a map with a default
func getString(m map[string]interface{}, key, defaultValue string) string {
	if val, ok := m[key]; ok {
		if str, ok := val.(string); ok {
			return str
		}
	}
	return defaultValue
}

// getInt gets an int value from a map with a default
func getInt(m map[string]interface{}, key string, defaultValue int) int {
	if val, ok := m[key]; ok {
		if num, ok := val.(float64); ok {
			return int(num)
		}
		if str, ok := val.(string); ok {
			if num, err := strconv.Atoi(str); err == nil {
				return num
			}
		}
	}
	return defaultValue
}

// generateID generates a simple ID (in a real implementation, you might want to use UUID)
func generateID() string {
	// This is a placeholder implementation
	// In a real implementation, you would generate a proper unique ID
	rand.Seed(time.Now().UnixNano())
	return fmt.Sprintf("server-%d", rand.Intn(1000000))
}