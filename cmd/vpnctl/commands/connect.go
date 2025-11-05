package commands

import (
	"fmt"
	"time"

	"github.com/spf13/cobra"
)

func newConnectCommand() *cobra.Command {
	connectCmd := &cobra.Command{
		Use:   "connect [server-name]",
		Short: "Connect to a VPN server",
		Long: `Connect to a VPN server by name, or connect to the best/fastest server.
If no server name is provided, you can use --best or --fastest flags to automatically
select the optimal server.`,
		Example: `  vpnctl connect               # Connect with default settings
  vpnctl connect server1       # Connect to a specific server
  vpnctl connect --best        # Connect to the best server
  vpnctl connect --fastest     # Connect to the fastest server`,
		Run: connectRun,
	}

	connectCmd.Flags().Bool("best", false, "Connect to the best server")
	connectCmd.Flags().Bool("fastest", false, "Connect to the fastest server")

	return connectCmd
}

func connectRun(cmd *cobra.Command, args []string) {
	// Parse flags
	best, _ := cmd.Flags().GetBool("best")
	fastest, _ := cmd.Flags().GetBool("fastest")

	// Validate arguments
	if len(args) > 1 {
		ExitWithError(1, "connect command accepts at most one argument")
	}

	if best && fastest {
		ExitWithError(1, "cannot use both --best and --fastest flags")
	}

	// Get server name if provided
	var serverName string
	if len(args) == 1 {
		serverName = args[0]
	}

	// Show progress
	fmt.Print("Connecting to VPN")
	for i := 0; i < 3; i++ {
		fmt.Print(".")
		time.Sleep(500 * time.Millisecond)
	}
	fmt.Println()

	// Simulate connection logic
	if serverName != "" {
		fmt.Printf("✅ Connected to server: %s\n", serverName)
	} else if best {
		fmt.Println("✅ Connected to the best available server")
	} else if fastest {
		fmt.Println("✅ Connected to the fastest available server")
	} else {
		fmt.Println("✅ Connected to default server")
	}

	// Show connection details
	fmt.Println("⏱️  Duration: 00:00:00")
	fmt.Println("📊 Download: 0 MB ↑ | Upload: 0 MB ↓")
	fmt.Println("🌐 IP: 192.168.1.100")
}