package commands

import (
	"fmt"
	"os"

	"github.com/spf13/cobra"
)

// NewRootCommand creates the root CLI command
func NewRootCommand() *cobra.Command {
	rootCmd := &cobra.Command{
		Use:   "vpnctl",
		Short: "VPN Client CLI",
		Long: `vpnctl is a command line interface for managing your VPN client.
It allows you to connect to servers, manage configurations, view statistics, and more.`,
	}

	// Add subcommands
	rootCmd.AddCommand(
		newConnectCommand(),
		newDisconnectCommand(),
		newStatusCommand(),
		newServerCommand(),
		newStatsCommand(),
		newLogsCommand(),
		newConfigCommand(),
		newSpeedTestCommand(),
		newTrafficCommand(),
		newRestartCommand(),
	)

	// Global flags
	rootCmd.PersistentFlags().BoolP("json", "j", false, "Output in JSON format")
	rootCmd.PersistentFlags().BoolP("verbose", "v", false, "Verbose output")
	rootCmd.PersistentFlags().BoolP("quiet", "q", false, "Quiet mode")
	rootCmd.PersistentFlags().String("config", "", "Config file path")
	rootCmd.PersistentFlags().Int("timeout", 30, "Timeout in seconds")

	// Version command
	rootCmd.Version = "1.0.0"
	rootCmd.SetVersionTemplate("vpnctl version {{.Version}}\n")

	return rootCmd
}

// CLIStatus represents the status of the VPN connection
type CLIStatus struct {
	Connected bool   `json:"connected"`
	Server    string `json:"server"`
	Duration  string `json:"duration"`
	Upload    string `json:"upload"`
	Download  string `json:"download"`
	IP        string `json:"ip"`
}

// ServerInfo represents information about a VPN server
type ServerInfo struct {
	Name     string  `json:"name"`
	Country  string  `json:"country"`
	Ping     int     `json:"ping"`
	Speed    float64 `json:"speed"`
	Protocol string  `json:"protocol"`
	Status   string  `json:"status"`
}

// ExitWithError exits with error code and message
func ExitWithError(code int, message string) {
	fmt.Fprintln(os.Stderr, message)
	os.Exit(code)
}