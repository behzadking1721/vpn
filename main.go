package main

import (
	"fmt"
	"os"

	// Import managers package using the correct path
	"vpn-client/internal/managers"
)

func main() {
	fmt.Println("VPN Client Application")

	// Create a connection manager instance
	cm := managers.NewConnectionManager()
	if cm == nil {
		fmt.Println("Failed to create connection manager")
		os.Exit(1)
	}

	fmt.Printf("Connection manager created with status: %v\n", cm.GetStatus())

	// Attempt to connect
	if err := cm.Connect(nil); err != nil {
		fmt.Printf("Connection attempt failed: %v\n", err)
	} else {
		fmt.Println("Successfully connected to the VPN")
	}

	fmt.Println("Application finished")
}