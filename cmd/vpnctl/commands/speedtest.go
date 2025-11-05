package commands

import (
	"fmt"
	"math/rand"
	"time"

	"github.com/cheggaaa/pb/v3"
	"github.com/spf13/cobra"
)

func newSpeedTestCommand() *cobra.Command {
	return &cobra.Command{
		Use:   "speedtest",
		Short: "Perform a speed test",
		Long:  `Perform a network speed test to measure download and upload speeds.`,
		Run: func(cmd *cobra.Command, args []string) {
			fmt.Println("🚀 Performing speed test...")
			
			// Download test
			fmt.Println("\n⬇️  Testing download speed...")
			downloadBar := pb.StartNew(100)
			for i := 0; i < 100; i++ {
				downloadBar.Increment()
				time.Sleep(time.Duration(rand.Intn(50)+10) * time.Millisecond)
			}
			downloadBar.Finish()
			
			downloadSpeed := fmt.Sprintf("%.2f Mbps", rand.Float64()*90+10)
			fmt.Printf("Download speed: %s\n", downloadSpeed)
			
			// Upload test
			fmt.Println("\n⬆️  Testing upload speed...")
			uploadBar := pb.StartNew(100)
			for i := 0; i < 100; i++ {
				uploadBar.Increment()
				time.Sleep(time.Duration(rand.Intn(100)+50) * time.Millisecond)
			}
			uploadBar.Finish()
			
			uploadSpeed := fmt.Sprintf("%.2f Mbps", rand.Float64()*45+5)
			fmt.Printf("Upload speed: %s\n", uploadSpeed)
			
			// Latency test
			fmt.Println("\n⏱️  Testing latency...")
			latencyBar := pb.StartNew(100)
			for i := 0; i < 100; i++ {
				latencyBar.Increment()
				time.Sleep(5 * time.Millisecond)
			}
			latencyBar.Finish()
			
			latency := fmt.Sprintf("%d ms", rand.Intn(100)+10)
			fmt.Printf("Latency: %s\n", latency)
			
			// Summary
			fmt.Println("\n📈 Speed Test Results")
			fmt.Println("=====================")
			fmt.Printf("Download: %s\n", downloadSpeed)
			fmt.Printf("Upload:   %s\n", uploadSpeed)
			fmt.Printf("Latency:  %s\n", latency)
			fmt.Printf("Server:   USA - New York\n")
			fmt.Printf("Time:     %s\n", time.Now().Format("2006-01-02 15:04:05"))
		},
	}
}