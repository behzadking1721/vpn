package utils

import (
	"c:/Users/behza/OneDrive/Documents/vpn/src/core"
	"math/rand"
	"time"
)

// PingResult represents the result of a ping operation
type PingResult struct {
	ServerID string
	Ping     int
	Error    error
}

// PingServers pings a list of servers and updates their ping times
func PingServers(servers []core.Server) []PingResult {
	results := make([]PingResult, len(servers))
	
	// In a real implementation, you would send ICMP echo requests or TCP connect attempts
	// For this demo, we'll simulate ping results
	
	for i, server := range servers {
		// Simulate network delay
		time.Sleep(time.Duration(rand.Intn(100)) * time.Millisecond)
		
		// Generate a realistic ping time
		pingTime := rand.Intn(200) + 10 // Between 10-210 ms
		
		results[i] = PingResult{
			ServerID: server.ID,
			Ping:     pingTime,
			Error:    nil,
		}
	}
	
	return results
}

// PingServer pings a single server
func PingServer(server core.Server) PingResult {
	// In a real implementation, you would send ICMP echo requests or TCP connect attempts
	// For this demo, we'll simulate a ping result
	
	// Simulate network delay
	time.Sleep(time.Duration(rand.Intn(100)) * time.Millisecond)
	
	// Generate a realistic ping time
	pingTime := rand.Intn(200) + 10 // Between 10-210 ms
	
	return PingResult{
		ServerID: server.ID,
		Ping:     pingTime,
		Error:    nil,
	}
}