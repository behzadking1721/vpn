package commands

import (
	"encoding/json"
	"fmt"
	"math/rand"
	"time"

	"github.com/cheggaaa/pb/v3"
	"github.com/spf13/cobra"
)

func newStatsCommand() *cobra.Command {
	return &cobra.Command{
		Use:   "stats",
		Short: "Show data usage statistics",
		Long:  `Show detailed data usage statistics including total data transferred, session history, and more.`,
		Run: func(cmd *cobra.Command, args []string) {
			// Check for JSON flag
			jsonOutput, _ := cmd.Flags().GetBool("json")

			if jsonOutput {
				// Simulate JSON output
				stats := map[string]interface{}{
					"total_data_sent":     "256.4 MB",
					"total_data_received": "1.2 GB",
					"current_session": map[string]string{
						"duration":  "2:34:15",
						"data_sent": "45.2 MB",
						"data_recv": "128.7 MB",
						"server":    "USA-NY",
					},
					"last_24_hours": map[string]string{
						"data_sent":     "128.1 MB",
						"data_received": "756.3 MB",
					},
				}

				jsonData, err := json.MarshalIndent(stats, "", "  ")
				if err != nil {
					ExitWithError(1, fmt.Sprintf("Failed to marshal stats to JSON: %v", err))
				}
				fmt.Println(string(jsonData))
				return
			}

			// Show progress bar
			fmt.Println("📊 Fetching statistics...")
			bar := pb.StartNew(100)
			for i := 0; i < 100; i++ {
				bar.Increment()
				time.Sleep(20 * time.Millisecond)
			}
			bar.Finish()

			// Display stats
			fmt.Println("\n📈 Data Usage Statistics")
			fmt.Println("========================")
			fmt.Printf("Total Data Sent:     %s\n", formatBytes(rand.Int63n(1000000000)))
			fmt.Printf("Total Data Received: %s\n", formatBytes(rand.Int63n(5000000000)))
			fmt.Printf("Current Session:     %s\n", formatDuration(time.Now().Add(-time.Hour).Unix()))
			fmt.Printf("Active Connections:  %d\n", rand.Intn(5)+1)
		},
	}
}

// formatBytes formats bytes to human readable format
func formatBytes(bytes int64) string {
	const (
		KB = 1024
		MB = KB * 1024
		GB = MB * 1024
	)

	switch {
	case bytes >= GB:
		return fmt.Sprintf("%.2f GB", float64(bytes)/GB)
	case bytes >= MB:
		return fmt.Sprintf("%.2f MB", float64(bytes)/MB)
	case bytes >= KB:
		return fmt.Sprintf("%.2f KB", float64(bytes)/KB)
	default:
		return fmt.Sprintf("%d B", bytes)
	}
}

// formatDuration formats duration to human readable format
func formatDuration(seconds int64) string {
	hours := seconds / 3600
	minutes := (seconds % 3600) / 60
	secs := seconds % 60
	return fmt.Sprintf("%02d:%02d:%02d", hours, minutes, secs)
}
