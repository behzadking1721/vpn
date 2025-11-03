package main

import (
	"fmt"
	"os"
	
	// Import managers package using the correct path
	"vpnclient/internal/managers"
)

var version = "dev"

func main() {
	fmt.Printf("VPN Client Application - Version %s\n", version)
	
	// Create a connection manager instance
	cm := managers.NewConnectionManager()
	if cm == nil {
		fmt.Println("Failed to create connection manager")
		os.Exit(1)
	}
	
	fmt.Printf("Connection manager created with status: %v\n", cm.GetStatus())
	
	// Try to connect (this will likely fail in test environment)
	err := cm.Connect(nil)
	if err != nil {
		fmt.Printf("Connection attempt failed (expected in test environment): %v\n", err)
	}
	
	fmt.Println("Application finished successfully")
}