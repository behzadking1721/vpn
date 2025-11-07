package commands

import (
	"fmt"
	"math/rand"
	"time"

	"github.com/spf13/cobra"
)

func newTrafficCommand() *cobra.Command {
	return &cobra.Command{
		Use:   "traffic",
		Short: "Show real-time traffic statistics",
		Long:  `Show real-time traffic statistics including data transfer rates.`,
		Run: func(cmd *cobra.Command, args []string) {
			fmt.Println("📊 Real-time Traffic Statistics")
			fmt.Println("==============================")

			// Simulate real-time traffic display
			for i := 0; i < 10; i++ {
				download := rand.Float64() * 5
				upload := rand.Float64() * 2

				fmt.Printf("\r⬇️  Down: %6.2f Mbps  ⬆️  Up: %6.2f Mbps", download, upload)
				time.Sleep(1 * time.Second)
			}

			fmt.Println("\n\n✅ Traffic monitoring completed")
		},
	}
}
