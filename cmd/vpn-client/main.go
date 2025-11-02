package main

import (
	"fmt"
	"log"
	
	"vpn-client/internal/managers"
)

func main() {
	fmt.Println("🚀 VPN Client Starting...")
	
	// Initialize managers
	cm := managers.NewConnectionManager()
	
	// Start a test connection
	err := cm.Connect(nil)
	if err != nil {
		log.Fatal("Failed to connect:", err)
	}
	
	fmt.Println("✅ VPN Client started successfully")
	fmt.Println("Current status:", cm.GetStatusString())
	fmt.Println("Press Ctrl+C to stop...")
	
	// Keep the application running
	select {}
}