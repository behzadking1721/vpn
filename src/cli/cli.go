package cli

import (
	"bufio"
	"c:/Users/behza/OneDrive/Documents/vpn/src/core"
	"c:/Users/behza/OneDrive/Documents/vpn/src/managers"
	"c:/Users/behza/OneDrive/Documents/vpn/src/protocols"
	"c:/Users/behza/OneDrive/Documents/vpn/src/utils"
	"fmt"
	"os"
	"strconv"
	"strings"
)

// CLI represents the command line interface
type CLI struct {
	serverManager *managers.ServerManager
	connManager   *managers.ConnectionManager
	configManager *managers.ConfigManager
	scanner       *bufio.Scanner
}

// NewCLI creates a new CLI instance
func NewCLI(
	serverMgr *managers.ServerManager,
	connMgr *managers.ConnectionManager,
	configMgr *managers.ConfigManager) *CLI {
	
	return &CLI{
		serverManager: serverMgr,
		connManager:   connMgr,
		configManager: configMgr,
		scanner:       bufio.NewScanner(os.Stdin),
	}
}

// Run starts the CLI interface
func (c *CLI) Run() {
	fmt.Println("VPN Client CLI")
	fmt.Println("==============")
	
	for {
		c.showMenu()
		fmt.Print("Enter your choice: ")
		
		if !c.scanner.Scan() {
			break
		}
		
		choice := strings.TrimSpace(c.scanner.Text())
		
		switch choice {
		case "1":
			c.listServers()
		case "2":
			c.addServer()
		case "3":
			c.connectToServer()
		case "4":
			c.disconnect()
		case "5":
			c.showStatus()
		case "6":
			c.pingServers()
		case "7":
			c.findFastestServer()
		case "8":
			fmt.Println("Goodbye!")
			return
		default:
			fmt.Println("Invalid choice. Please try again.")
		}
		
		fmt.Println()
	}
}

// showMenu displays the main menu
func (c *CLI) showMenu() {
	fmt.Println("\n--- Main Menu ---")
	fmt.Println("1. List servers")
	fmt.Println("2. Add server")
	fmt.Println("3. Connect to server")
	fmt.Println("4. Disconnect")
	fmt.Println("5. Show status")
	fmt.Println("6. Ping servers")
	fmt.Println("7. Find fastest server")
	fmt.Println("8. Exit")
}

// listServers displays all servers
func (c *CLI) listServers() {
	servers := c.serverManager.GetAllServers()
	
	if len(servers) == 0 {
		fmt.Println("No servers configured.")
		return
	}
	
	fmt.Println("\nConfigured Servers:")
	fmt.Println("===================")
	for i, server := range servers {
		status := "Disabled"
		if server.Enabled {
			status = "Enabled"
		}
		
		ping := "N/A"
		if server.Ping > 0 {
			ping = fmt.Sprintf("%d ms", server.Ping)
		}
		
		fmt.Printf("%d. %s (%s:%d) [%s] - %s - Ping: %s\n", 
			i+1, server.Name, server.Host, server.Port, server.Protocol, status, ping)
	}
}

// addServer allows user to add a new server
func (c *CLI) addServer() {
	fmt.Println("\n--- Add New Server ---")
	
	server := core.Server{
		ID: utils.GenerateID(),
	}
	
	fmt.Print("Server name: ")
	if c.scanner.Scan() {
		server.Name = strings.TrimSpace(c.scanner.Text())
	}
	
	fmt.Print("Host: ")
	if c.scanner.Scan() {
		server.Host = strings.TrimSpace(c.scanner.Text())
	}
	
	fmt.Print("Port: ")
	if c.scanner.Scan() {
		portStr := strings.TrimSpace(c.scanner.Text())
		port, err := strconv.Atoi(portStr)
		if err != nil {
			fmt.Println("Invalid port number. Using default port 443.")
			server.Port = 443
		} else {
			server.Port = port
		}
	} else {
		server.Port = 443
	}
	
	// Protocol selection
	fmt.Println("Available protocols:")
	protocols := []core.ProtocolType{
		core.ProtocolVMess,
		core.ProtocolVLESS,
		core.ProtocolShadowsocks,
		core.ProtocolTrojan,
		core.ProtocolReality,
		core.ProtocolHysteria,
		core.ProtocolTUIC,
		core.ProtocolSSH,
	}
	
	for i, proto := range protocols {
		fmt.Printf("%d. %s\n", i+1, proto)
	}
	
	fmt.Print("Select protocol (1-8): ")
	if c.scanner.Scan() {
		choiceStr := strings.TrimSpace(c.scanner.Text())
		choice, err := strconv.Atoi(choiceStr)
		if err != nil || choice < 1 || choice > len(protocols) {
			fmt.Println("Invalid choice. Using VMess as default.")
			server.Protocol = core.ProtocolVMess
		} else {
			server.Protocol = protocols[choice-1]
		}
	} else {
		server.Protocol = core.ProtocolVMess
	}
	
	// For Shadowsocks, ask for method and password
	if server.Protocol == core.ProtocolShadowsocks {
		fmt.Print("Encryption method: ")
		if c.scanner.Scan() {
			server.Method = strings.TrimSpace(c.scanner.Text())
		}
		
		fmt.Print("Password: ")
		if c.scanner.Scan() {
			server.Password = strings.TrimSpace(c.scanner.Text())
		}
	} else {
		// For other protocols, ask for general encryption settings
		fmt.Print("Encryption (optional): ")
		if c.scanner.Scan() {
			server.Encryption = strings.TrimSpace(c.scanner.Text())
		}
		
		fmt.Print("Password/UUID (optional): ")
		if c.scanner.Scan() {
			server.Password = strings.TrimSpace(c.scanner.Text())
		}
	}
	
	// Enable by default
	server.Enabled = true
	
	// Add to server manager
	err := c.serverManager.AddServer(server)
	if err != nil {
		fmt.Printf("Failed to add server: %v\n", err)
		return
	}
	
	fmt.Println("Server added successfully!")
}

// connectToServer allows user to connect to a server
func (c *CLI) connectToServer() {
	servers := c.serverManager.GetAllServers()
	if len(servers) == 0 {
		fmt.Println("No servers configured. Please add a server first.")
		return
	}
	
	fmt.Println("\n--- Connect to Server ---")
	c.listServers()
	
	fmt.Print("Enter server number: ")
	if !c.scanner.Scan() {
		return
	}
	
	choiceStr := strings.TrimSpace(c.scanner.Text())
	choice, err := strconv.Atoi(choiceStr)
	if err != nil || choice < 1 || choice > len(servers) {
		fmt.Println("Invalid server number.")
		return
	}
	
	server := servers[choice-1]
	
	fmt.Printf("Connecting to %s...\n", server.Name)
	err = c.connManager.Connect(server)
	if err != nil {
		fmt.Printf("Failed to connect: %v\n", err)
		return
	}
	
	fmt.Println("Connected successfully!")
}

// disconnect disconnects from the current server
func (c *CLI) disconnect() {
	fmt.Println("\n--- Disconnect ---")
	
	err := c.connManager.Disconnect()
	if err != nil {
		fmt.Printf("Failed to disconnect: %v\n", err)
		return
	}
	
	fmt.Println("Disconnected successfully!")
}

// showStatus displays the current connection status
func (c *CLI) showStatus() {
	fmt.Println("\n--- Connection Status ---")
	
	status := c.connManager.GetStatus()
	fmt.Printf("Status: %s\n", status)
	
	if status == core.StatusConnected {
		info := c.connManager.GetConnectionInfo()
		fmt.Printf("Connected since: %s\n", info.StartTime.Format("2006-01-02 15:04:05"))
		
		if info.DataSent > 0 || info.DataReceived > 0 {
			fmt.Printf("Data sent: %s\n", utils.FormatBytes(info.DataSent))
			fmt.Printf("Data received: %s\n", utils.FormatBytes(info.DataReceived))
		}
		
		if currentServer := c.connManager.GetCurrentServer(); currentServer != nil {
			fmt.Printf("Connected to: %s (%s:%d)\n", 
				currentServer.Name, currentServer.Host, currentServer.Port)
		}
	}
}

// pingServers pings all enabled servers
func (c *CLI) pingServers() {
	fmt.Println("\n--- Ping Servers ---")
	
	servers := c.serverManager.GetAllServers()
	if len(servers) == 0 {
		fmt.Println("No servers configured.")
		return
	}
	
	fmt.Println("Pinging servers...")
	results := utils.PingServers(servers)
	
	for _, result := range results {
		if result.Error == nil {
			c.serverManager.UpdateServerPing(result.ServerID, result.Ping)
			server, _ := c.serverManager.GetServer(result.ServerID)
			fmt.Printf("- %s: %d ms\n", server.Name, result.Ping)
		} else {
			server, _ := c.serverManager.GetServer(result.ServerID)
			fmt.Printf("- %s: Error (%v)\n", server.Name, result.Error)
		}
	}
	
	fmt.Println("Ping completed!")
}

// findFastestServer finds and displays the fastest server
func (c *CLI) findFastestServer() {
	fmt.Println("\n--- Find Fastest Server ---")
	
	server, err := c.serverManager.FindFastestServer()
	if err != nil {
		fmt.Printf("Error finding fastest server: %v\n", err)
		return
	}
	
	fmt.Printf("Fastest server: %s (%d ms)\n", server.Name, server.Ping)
}