package main

import (
	"vpnclient/internal/api"
	"vpnclient/internal/database"
	"vpnclient/internal/managers"
)

var version = "dev"

func main() {
	// Initialize database
	store, err := database.NewDB("data")
	if err != nil {
		panic(err)
	}
	defer store.Close()

	// Initialize managers
	serverManager := managers.NewServerManager(store)
	connectionManager := managers.NewConnectionManager()

	// Start API server
	apiServer := api.NewServer(":8080", serverManager, connectionManager)
	if err := apiServer.Start(); err != nil {
		panic(err)
	}
}