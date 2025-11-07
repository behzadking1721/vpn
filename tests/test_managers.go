package main

import (
	"fmt"
	"vpnclient/internal/managers"
)

func TestNewConnectionManager() {
	cm := managers.NewConnectionManager()
	if cm == nil {
		fmt.Println("❌ Failed to create ConnectionManager")
		return
	}

	fmt.Println("✅ ConnectionManager created successfully")
	fmt.Printf("Initial status: %v\n", cm.GetStatus())
}

func main() {
	fmt.Println("Testing managers package...")

	// Run the test
	TestNewConnectionManager()

	fmt.Println("🎉 Test completed!")
}
