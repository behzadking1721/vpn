package main

import (
	"c:/Users/behza/OneDrive/Documents/vpn/src/core"
	"c:/Users/behza/OneDrive/Documents/vpn/src/protocols"
	"fmt"
	"time"
)

func main() {
	fmt.Println("Testing Protocol Integration")
	fmt.Println("============================")

	// Test protocol factory
	factory := protocols.NewProtocolFactory()

	// Test all protocols
	protocolTypes := []core.ProtocolType{
		core.ProtocolVMess,
		core.ProtocolVLESS,
		core.ProtocolShadowsocks,
		core.ProtocolTrojan,
		core.ProtocolReality,
		core.ProtocolHysteria,
		core.ProtocolTUIC,
		core.ProtocolSSH,
	}

	for _, protoType := range protocolTypes {
		fmt.Printf("\nTesting %s protocol...\n", protoType)
		
		handler, err := factory.CreateHandler(protoType)
		if err != nil {
			fmt.Printf("  Error creating handler: %v\n", err)
			continue
		}

		fmt.Printf("  Handler created successfully: %s\n", handler.GetProtocol())

		// Test with a sample server
		server := createSampleServer(protoType)
		
		// Test connection (this will be simulated for now)
		fmt.Printf("  Connecting to %s server...\n", protoType)
		err = handler.Connect(server)
		if err != nil {
			fmt.Printf("  Connection error: %v\n", err)
		} else {
			fmt.Printf("  Connected successfully\n")
			
			// Test data usage
			sent, received, err := handler.GetDataUsage()
			if err != nil {
				fmt.Printf("  Data usage error: %v\n", err)
			} else {
				fmt.Printf("  Data usage - Sent: %d, Received: %d\n", sent, received)
			}
			
			// Test connection details
			details, err := handler.GetConnectionDetails()
			if err != nil {
				fmt.Printf("  Connection details error: %v\n", err)
			} else {
				fmt.Printf("  Connection details: %v\n", details)
			}
			
			// Simulate connection duration
			time.Sleep(1 * time.Second)
			
			// Test disconnection
			fmt.Printf("  Disconnecting...\n")
			err = handler.Disconnect()
			if err != nil {
				fmt.Printf("  Disconnection error: %v\n", err)
			} else {
				fmt.Printf("  Disconnected successfully\n")
			}
		}
	}

	fmt.Println("\nProtocol testing completed!")
}

func createSampleServer(protocol core.ProtocolType) core.Server {
	server := core.Server{
		ID:       "test-server",
		Name:     fmt.Sprintf("Test %s Server", protocol),
		Host:     "example.com",
		Port:     8388,
		Protocol: protocol,
		Enabled:  true,
	}

	switch protocol {
	case core.ProtocolShadowsocks:
		server.Method = "aes-256-gcm"
		server.Password = "test-password"
	case core.ProtocolVMess, core.ProtocolVLESS:
		server.Encryption = "auto"
		server.Password = "test-uuid"
		server.TLS = true
	default:
		server.Password = "test-password"
	}

	return server
}