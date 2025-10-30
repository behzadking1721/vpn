package main

import (
	"c:/Users/behza/OneDrive/Documents/vpn/src/api"
	"c:/Users/behza/OneDrive/Documents/vpn/src/cli"
	"c:/Users/behza/OneDrive/Documents/vpn/src/core"
	"c:/Users/behza/OneDrive/Documents/vpn/src/managers"
	"c:/Users/behza/OneDrive/Documents/vpn/src/protocols"
	"c:/Users/behza/OneDrive/Documents/vpn/src/utils"
	"fmt"
	"log"
	"os"
)

func main() {
	fmt.Println("VPN Client Application")
	fmt.Println("======================")

	// Check command line arguments
	if len(os.Args) > 1 {
		switch os.Args[1] {
		case "--api":
			startAPIServer()
			return
		case "--cli":
			startCLI()
			return
		case "--help", "-h":
			showHelp()
			return
		}
	}

	// If no arguments, show help
	showHelp()
}

func showHelp() {
	fmt.Println("VPN Client - Multi-platform VPN application")
	fmt.Println()
	fmt.Println("Usage:")
	fmt.Println("  go run main.go --api    Start API server")
	fmt.Println("  go run main.go --cli    Start command-line interface")
	fmt.Println("  go run main.go --help   Show this help message")
	fmt.Println()
	fmt.Println("Examples:")
	fmt.Println("  go run main.go --api")
	fmt.Println("  go run main.go --cli")
}

func startAPIServer() {
	fmt.Println("Starting VPN Client API Server...")
	
	// Initialize managers
	serverManager := managers.NewServerManager()
	connectionManager := managers.NewConnectionManager()
	subscriptionManager := managers.NewSubscriptionManager(serverManager)
	configManager := managers.NewConfigManager("./config/app.json")

	// Add sample servers for testing
	server1 := core.Server{
		ID:         utils.GenerateID(),
		Name:       "🇯🇵 Japan Server",
		Host:       "jp.example.com",
		Port:       443,
		Protocol:   core.ProtocolVMess,
		Encryption: "auto",
		TLS:        true,
		Remark:     "Tokyo Server",
		Enabled:    true,
		Ping:       22,
	}

	server2 := core.Server{
		ID:       utils.GenerateID(),
		Name:     "🇺🇸 USA Server",
		Host:     "us.example.com",
		Port:     8388,
		Protocol: core.ProtocolShadowsocks,
		Method:   "aes-256-gcm",
		Password: "test-password",
		Remark:   "New York Server",
		Enabled:  true,
		Ping:     45,
	}

	serverManager.AddServer(server1)
	serverManager.AddServer(server2)

	// Create and start API server
	apiServer := api.NewAPIServer(serverManager, connectionManager, configManager)
	fmt.Println("API Server listening on http://localhost:8080")
	log.Fatal(apiServer.Start(":8080"))
}

func startCLI() {
	fmt.Println("Starting VPN Client CLI...")
	
	// Initialize managers
	serverManager := managers.NewServerManager()
	connectionManager := managers.NewConnectionManager()
	configManager := managers.NewConfigManager("./config/app.json")

	// Add sample servers for testing
	server1 := core.Server{
		ID:         utils.GenerateID(),
		Name:       "🇯🇵 Japan Server",
		Host:       "jp.example.com",
		Port:       443,
		Protocol:   core.ProtocolVMess,
		Encryption: "auto",
		TLS:        true,
		Remark:     "Tokyo Server",
		Enabled:    true,
		Ping:       22,
	}

	server2 := core.Server{
		ID:       utils.GenerateID(),
		Name:     "🇺🇸 USA Server",
		Host:     "us.example.com",
		Port:     8388,
		Protocol: core.ProtocolShadowsocks,
		Method:   "aes-256-gcm",
		Password: "test-password",
		Remark:   "New York Server",
		Enabled:  true,
		Ping:     45,
	}

	serverManager.AddServer(server1)
	serverManager.AddServer(server2)

	// Create and start CLI
	cliInterface := cli.NewCLI(serverManager, connectionManager, configManager)
	cliInterface.Run()
}