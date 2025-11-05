package commands

import (
	"fmt"
	"time"

	"github.com/spf13/cobra"
)

func newDisconnectCommand() *cobra.Command {
	return &cobra.Command{
		Use:   "disconnect",
		Short: "Disconnect from the VPN",
		Long:  `Disconnect from the currently connected VPN server.`,
		Example: `  vpnctl disconnect    # Disconnect from current VPN server`,
		Run: func(cmd *cobra.Command, args []string) {
			// Show progress
			fmt.Print("Disconnecting from VPN")
			for i := 0; i < 3; i++ {
				fmt.Print(".")
				time.Sleep(500 * time.Millisecond)
			}
			fmt.Println()

			// Simulate disconnection
			fmt.Println("✅ Disconnected from VPN server")
		},
	}
}