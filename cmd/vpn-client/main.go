package main

import (
	"fmt"
	"os"
	"path/filepath"

	"vpnclient/internal/api"
	"vpnclient/internal/database"
	"vpnclient/internal/managers"
)

func main() {
	fmt.Println("🚀 VPN Client Starting...")

	// Determine data directory
	dataDir := filepath.Join(".", "data")
	if envDataDir := os.Getenv("VPN_DATA_DIR"); envDataDir != "" {
		dataDir = envDataDir
	}

	// Initialize database/store
	fmt.Println("📦 Initializing database...")
	store, err := database.NewDB(dataDir)
	if err != nil {
		fmt.Printf("❌ Failed to initialize database: %v\n", err)
		os.Exit(1)
	}
	defer store.Close()

	// Initialize managers
	fmt.Println("⚙️  Initializing managers...")
	serverManager := managers.NewServerManager(store)
	connectionManager := managers.NewConnectionManager()

	// Create API server
	fmt.Println("🌐 Creating API server...")
	apiServer := api.NewServer(
		":8080",
		serverManager,
		connectionManager,
	)

	// Start server
	fmt.Println("✅ VPN Client started successfully")
	fmt.Println("📡 API Server listening on http://localhost:8080")
	fmt.Println("🔗 API Endpoints:")
	fmt.Println("   GET    /api/servers           - List all servers")
	fmt.Println("   POST   /api/servers           - Add a server")
	fmt.Println("   GET    /api/servers/{id}      - Get server details")
	fmt.Println("   PUT    /api/servers/{id}      - Update server")
	fmt.Println("   DELETE /api/servers/{id}      - Delete server")
	fmt.Println("   POST   /api/connect           - Connect to server")
	fmt.Println("   POST   /api/disconnect        - Disconnect")
	fmt.Println("   GET    /api/status            - Get connection status")
	fmt.Println("   GET    /api/stats             - Get connection statistics")
	fmt.Println("   GET    /health                - Health check")
	fmt.Println()
	fmt.Println("Press Ctrl+C to stop...")

	if err := apiServer.Start(); err != nil {
		fmt.Printf("❌ Server error: %v\n", err)
		os.Exit(1)
	}
}
