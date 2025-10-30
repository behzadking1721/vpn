package main

import (
	"c:/Users/behza/OneDrive/Documents/vpn/src/api"
	"c:/Users/behza/OneDrive/Documents/vpn/src/cli"
	"c:/Users/behza/OneDrive/Documents/vpn/src/core"
	"c:/Users/behza/OneDrive/Documents/vpn/src/managers"
	"c:/Users/behza/OneDrive/Documents/vpn/src/protocols"
	"c:/Users/behza/OneDrive/Documents/vpn/src/updater"
	"c:/Users/behza/OneDrive/Documents/vpn/src/utils"
	"fmt"
	"log"
	"os"
	"path/filepath"
)

const (
	// Application version
	AppVersion = "1.0.0"
	
	// Update API URL (this should point to your update server)
	UpdateAPIURL = "https://your-update-server.com/api/latest-release"
)

func main() {
	fmt.Println("VPN Client Application")
	fmt.Println("======================")

	// Initialize logger
	logger := utils.NewLogger(&utils.LoggerConfig{
		Level:     utils.LogLevelInfo,
		File:      "./logs/app.log",
		Timestamp: true,
	})
	defer logger.Close()

	logger.Info("Starting VPN Client v%s", AppVersion)

	// Check for updates
	checkForUpdates(logger)

	// Check command line arguments
	if len(os.Args) > 1 {
		switch os.Args[1] {
		case "--api":
			startAPIServer(logger)
			return
		case "--cli":
			startCLI(logger)
			return
		case "--help", "-h":
			showHelp()
			return
		}
	}

	// If no arguments, show help
	showHelp()
}

func checkForUpdates(logger *utils.Logger) {
	logger.Info("Checking for updates...")

	// Create updater
	upd, err := updater.NewUpdater(AppVersion, UpdateAPIURL)
	if err != nil {
		logger.Error("Failed to create updater: %v", err)
		return
	}

	// Check for update
	release, err := upd.CheckForUpdate()
	if err != nil {
		logger.Error("Failed to check for updates: %v", err)
		return
	}

	if release != nil {
		logger.Info("New version available: %s", release.Version)
		
		// In a real implementation, you would:
		// 1. Notify the user about the update
		// 2. Download the update
		// 3. Apply the update
		// 4. Restart the application
		
		fmt.Printf("New version %s is available!\n", release.Version)
		fmt.Printf("Release notes: %s\n", release.Notes)
		fmt.Println("Please visit our website to download the latest version.")
	} else {
		logger.Info("Application is up to date")
	}
}

func startAPIServer(logger *utils.Logger) {
	logger.Info("Starting API server")

	// Initialize managers
	configManager := managers.NewConfigManager("./config/app.json")
	serverManager := managers.NewServerManager()
	connManager := managers.NewConnectionManager()
	
	// Initialize protocol factory
	protocolFactory := protocols.NewProtocolFactory()
	
	// Register protocol factory with connection manager
	connManager.SetProtocolFactory(protocolFactory)
	
	// Create API server
	apiServer := api.NewAPIServer(serverManager, connManager, configManager)
	
	// Get API port from config or use default
	config := configManager.GetConfig()
	port := config.APIPort
	if port == 0 {
		port = 8080
	}
	
	logger.Info("API server starting on port %d", port)
	
	// Start the server
	if err := apiServer.Start(fmt.Sprintf(":%d", port)); err != nil {
		logger.Error("Failed to start API server: %v", err)
		log.Fatal(err)
	}
}

func startCLI(logger *utils.Logger) {
	logger.Info("Starting CLI mode")
	
	// Initialize managers
	configManager := managers.NewConfigManager("./config/app.json")
	serverManager := managers.NewServerManager()
	connManager := managers.NewConnectionManager()
	
	// Initialize protocol factory
	protocolFactory := protocols.NewProtocolFactory()
	
	// Register protocol factory with connection manager
	connManager.SetProtocolFactory(protocolFactory)
	
	// Create CLI
	cliApp := cli.NewCLI(serverManager, connManager, configManager)
	
	// Run CLI
	cliApp.Run()
}

func showHelp() {
	execName := filepath.Base(os.Args[0])
	
	fmt.Printf("Usage: %s [option]\n\n", execName)
	fmt.Println("Options:")
	fmt.Println("  --api     Start API server")
	fmt.Println("  --cli     Start command-line interface")
	fmt.Println("  --help    Show this help message")
	fmt.Println()
	fmt.Println("Examples:")
	fmt.Printf("  %s --api\n", execName)
	fmt.Printf("  %s --cli\n", execName)
}

// GetAppVersion returns the application version
func GetAppVersion() string {
	return AppVersion
}

// GetUpdateAPIURL returns the update API URL
func GetUpdateAPIURL() string {
	return UpdateAPIURL
}