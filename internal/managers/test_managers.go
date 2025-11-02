package managers

import (
	"fmt"
)

func TestNewConnectionManagerExample() {
	cm := NewConnectionManager()
	if cm == nil {
		fmt.Println("❌ Failed to create ConnectionManager")
		return
	}

	fmt.Println("✅ ConnectionManager created successfully")
	fmt.Printf("Initial status: %v\n", cm.GetStatus())
}