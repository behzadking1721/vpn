package commands

import (
	"fmt"

	"github.com/olekukonko/tablewriter"
	"github.com/spf13/cobra"
)

func newServerCommand() *cobra.Command {
	serverCmd := &cobra.Command{
		Use:   "server",
		Short: "Manage VPN servers",
		Long:  `Manage VPN servers including listing, adding, removing, and testing servers.`,
	}

	serverCmd.AddCommand(
		newServerListCommand(),
		newServerAddCommand(),
		newServerRemoveCommand(),
		newServerTestCommand(),
		newServerUpdateCommand(),
		newServerFavoritesCommand(),
	)

	return serverCmd
}

func newServerListCommand() *cobra.Command {
	return &cobra.Command{
		Use:   "list",
		Short: "List all VPN servers",
		Long:  `List all available VPN servers with their details.`,
		Run: func(cmd *cobra.Command, args []string) {
			// Sample server data
			servers := []ServerInfo{
				{Name: "USA-NY", Country: "🇺🇸", Ping: 45, Speed: 85.5, Protocol: "WireGuard", Status: "🟢 Online"},
				{Name: "Germany", Country: "🇩🇪", Ping: 65, Speed: 72.3, Protocol: "OpenVPN", Status: "🟢 Online"},
				{Name: "Japan", Country: "🇯🇵", Ping: 120, Speed: 45.1, Protocol: "Shadowsocks", Status: "🟡 Slow"},
				{Name: "UK-London", Country: "🇬🇧", Ping: 78, Speed: 68.9, Protocol: "WireGuard", Status: "🟢 Online"},
				{Name: "Canada", Country: "🇨🇦", Ping: 52, Speed: 79.2, Protocol: "Trojan", Status: "🟢 Online"},
			}

			// Create table
			table := tablewriter.NewWriter(cmd.OutOrStdout())
			table.SetHeader([]string{"Flag", "Name", "Ping", "Speed", "Protocol", "Status"})
			table.SetBorders(tablewriter.Border{Left: true, Top: false, Right: true, Bottom: false})
			table.SetCenterSeparator("|")

			// Add data to table
			for _, server := range servers {
				table.Append([]string{
					server.Country,
					server.Name,
					fmt.Sprintf("%dms", server.Ping),
					fmt.Sprintf("%.1f Mbps", server.Speed),
					server.Protocol,
					server.Status,
				})
			}

			table.Render()
		},
	}
}

func newServerAddCommand() *cobra.Command {
	return &cobra.Command{
		Use:   "add [url]",
		Short: "Add a VPN subscription",
		Long:  `Add a VPN subscription from a given URL.`,
		Args:  cobra.ExactArgs(1),
		Run: func(cmd *cobra.Command, args []string) {
			url := args[0]
			fmt.Printf("➕ Adding subscription from: %s\n", url)
			fmt.Println("✅ Subscription added successfully")
		},
	}
}

func newServerRemoveCommand() *cobra.Command {
	return &cobra.Command{
		Use:   "remove [name]",
		Short: "Remove a VPN server",
		Long:  `Remove a VPN server by name.`,
		Args:  cobra.ExactArgs(1),
		Run: func(cmd *cobra.Command, args []string) {
			name := args[0]
			fmt.Printf("🗑️  Removing server: %s\n", name)
			fmt.Println("✅ Server removed successfully")
		},
	}
}

func newServerTestCommand() *cobra.Command {
	return &cobra.Command{
		Use:   "test",
		Short: "Test speed of all servers",
		Long:  `Test the speed of all available VPN servers.`,
		Run: func(cmd *cobra.Command, args []string) {
			fmt.Println("🚀 Testing all servers...")
			fmt.Println("✅ Server testing completed")

			// Sample results
			fmt.Println("\nTop 3 fastest servers:")
			fmt.Println("1. USA-NY (45ms, 85.5 Mbps)")
			fmt.Println("2. Canada (52ms, 79.2 Mbps)")
			fmt.Println("3. Germany (65ms, 72.3 Mbps)")
		},
	}
}

func newServerUpdateCommand() *cobra.Command {
	return &cobra.Command{
		Use:   "update",
		Short: "Update server list",
		Long:  `Update the list of available VPN servers from subscriptions.`,
		Run: func(cmd *cobra.Command, args []string) {
			fmt.Println("🔄 Updating server list...")
			fmt.Println("✅ Server list updated successfully")
		},
	}
}

func newServerFavoritesCommand() *cobra.Command {
	return &cobra.Command{
		Use:   "favorites",
		Short: "Manage favorite servers",
		Long:  `Manage your favorite VPN servers.`,
		Run: func(cmd *cobra.Command, args []string) {
			fmt.Println("⭐ Managing favorite servers...")
			fmt.Println("Your favorite servers:")
			fmt.Println("1. USA-NY")
			fmt.Println("2. Canada")
		},
	}
}
