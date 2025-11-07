package commands

import (
	"encoding/json"
	"fmt"
	"time"

	"github.com/fatih/color"
	"github.com/spf13/cobra"
)

func newStatusCommand() *cobra.Command {
	return &cobra.Command{
		Use:   "status",
		Short: "Show VPN connection status",
		Long:  `Show the current status of the VPN connection including server info, duration, and data usage.`,
		Example: `  vpnctl status      # Show current VPN status
  vpnctl status -j   # Show current VPN status in JSON format`,
		Run: statusRun,
	}
}

func statusRun(cmd *cobra.Command, args []string) {
	// Check for JSON flag
	jsonOutput, _ := cmd.Flags().GetBool("json")

	// Create status data
	status := CLIStatus{
		Connected: true,
		Server:    "USA - New York",
		Duration:  "12:34:56",
		Upload:    "125 MB",
		Download:  "45 MB",
		IP:        "192.168.1.100",
	}

	if jsonOutput {
		// Output in JSON format
		jsonData, err := json.MarshalIndent(status, "", "  ")
		if err != nil {
			ExitWithError(1, fmt.Sprintf("Failed to marshal status to JSON: %v", err))
		}
		fmt.Println(string(jsonData))
		return
	}

	// Output in human-readable format
	green := color.New(color.FgGreen).SprintFunc()
	blue := color.New(color.FgBlue).SprintFunc()
	yellow := color.New(color.FgYellow).SprintFunc()

	fmt.Println(green("🟢 Connected to:"), status.Server)
	fmt.Println(yellow("⏱️  Duration:"), status.Duration)
	fmt.Printf("%s Download: %s ↑ | Upload: %s ↓\n",
		blue("📊"),
		status.Upload,
		status.Download)
	fmt.Println(green("🌐 IP:"), status.IP)
	fmt.Println(time.Now().Format("2006-01-02 15:04:05"))
}
