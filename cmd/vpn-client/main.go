package main

import (
	"fmt"
	"log"
	"os"
	"os/signal"
	"path/filepath"
	"syscall"
	"time"

	"vpnclient/internal/api"
	"vpnclient/internal/database"
	"vpnclient/internal/logging"
	"vpnclient/internal/managers"
	"vpnclient/internal/notifications"
	"vpnclient/internal/stats"
	"vpnclient/internal/updater"
)

// Version information set at build time
var (
	version   = "dev"
	buildTime = "unknown"
	gitCommit = "unknown"
)

func main() {
	// Handle version flag
	for _, arg := range os.Args[1:] {
		if arg == "--version" || arg == "-v" {
			fmt.Printf("VPN Client %s\nBuild: %s\nCommit: %s\n", version, buildTime, gitCommit)
			return
		}
	}

	// Create data directory if it doesn't exist
	dataDir := "data"
	if err := os.MkdirAll(dataDir, 0755); err != nil {
		log.Fatalf("Failed to create data directory: %v", err)
	}

	// Initialize logger
	logFilePath := filepath.Join(dataDir, "vpn-client.log")
	logger, err := logging.NewLogger(logging.Config{
		Level:     logging.INFO,
		Output:    logFilePath,
		Timestamp: true,
	})
	if err != nil {
		log.Fatalf("Failed to initialize logger: %v", err)
	}
	defer logger.Close()

	logger.Info("Starting VPN Client version %s (build: %s, commit: %s)", version, buildTime, gitCommit)

	// Initialize database
	dbPath := filepath.Join(dataDir, "vpn.db")
	db, err := database.NewDB(dbPath)
	if err != nil {
		logger.Fatal("Failed to initialize database: %v", err)
	}
	logger.Info("Database initialized successfully")

	// Initialize managers
	serverManager := managers.NewServerManager(db)
	connectionManager := managers.NewConnectionManager()
	subscriptionManager := managers.NewSubscriptionManager(serverManager, db)

	// Initialize notification manager
	notificationManager := notifications.NewNotificationManager(100) // Keep max 100 notifications

	// Initialize stats manager
	statsManager := stats.NewStatsManager()

	// Set notification manager for managers that need it
	connectionManager.SetNotificationManager(notificationManager)
	serverManager.SetNotificationManager(notificationManager)
	subscriptionManager.SetNotificationManager(notificationManager)

	// Set logger for managers that need it
	connectionManager.SetLogger(logger)
	serverManager.SetLogger(logger)
	subscriptionManager.SetLogger(logger)

	// Set stats manager for connection manager
	connectionManager.SetStatsManager(statsManager)

	// Initialize subscription parser (not used directly here)
	_ = managers.NewSubscriptionParser()

	// Initialize updater
	updaterConfig := updater.Config{
		Interval: 24 * time.Hour, // Update once per day by default
		Enabled:  true,
	}

	updater := updater.NewUpdater(serverManager, subscriptionManager, updaterConfig, logger)

	// Start the updater
	updater.Start()
	defer updater.Stop()

	// Create API server
	apiServer := api.NewServer(
		":8080",
		serverManager,
		connectionManager,
		notificationManager,
		statsManager,
		updater,
		logger,
		logFilePath,
	)

	// Start the API server
	go func() {
		if err := apiServer.Start(); err != nil {
			logger.Fatal("Failed to start API server: %v", err)
		}
	}()

	// Setup signal handling for graceful shutdown
	sigChan := make(chan os.Signal, 1)
	signal.Notify(sigChan, syscall.SIGINT, syscall.SIGTERM)

	// Add a small delay to ensure server starts
	time.Sleep(100 * time.Millisecond)
	logger.Info("VPN Client started successfully!")
	fmt.Println("VPN Client started successfully!")
	fmt.Printf("VPN Client %s started\n", version)
	fmt.Println("API Server running on http://localhost:8080")

	// Log server startup
	logger.Info("API Server running on http://localhost:8080")

	// Wait for shutdown signal
	<-sigChan
	logger.Info("Shutting down...")
	fmt.Println("\nShutting down...")

	// Gracefully shutdown the server
	if err := apiServer.Shutdown(); err != nil {
		logger.Error("Error shutting down server: %v", err)
		log.Printf("Error shutting down server: %v", err)
	}

	logger.Info("VPN Client stopped.")
	fmt.Println("VPN Client stopped.")
}
