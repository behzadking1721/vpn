//go:build ignore

package main

import (
	"fmt"
	"vpn-client/internal/managers"
)

func main() {
	fmt.Println("Testing managers package...")

	// Create a connection manager instance
	cm := managers.NewConnectionManager()
	if cm == nil {
		fmt.Println("❌ Failed to create ConnectionManager")
		return
	}

	fmt.Println("✅ ConnectionManager created successfully")
	fmt.Printf("Initial status: %v\n", cm.GetStatus())

	fmt.Println("🎉 Example completed!")
}