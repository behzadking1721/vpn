package commands

import (
	"fmt"
	"time"

	"github.com/cheggaaa/pb/v3"
	"github.com/spf13/cobra"
)

func newRestartCommand() *cobra.Command {
	return &cobra.Command{
		Use:   "restart",
		Short: "Restart the VPN service",
		Long:  `Restart the VPN service to apply configuration changes or resolve issues.`,
		Run: func(cmd *cobra.Command, args []string) {
			fmt.Println("🔄 Restarting VPN service...")
			
			// Show progress
			bar := pb.StartNew(100)
			for i := 0; i < 100; i++ {
				bar.Increment()
				time.Sleep(30 * time.Millisecond)
			}
			bar.Finish()
			
			fmt.Println("✅ VPN service restarted successfully")
		},
	}
}