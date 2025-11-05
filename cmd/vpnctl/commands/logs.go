package commands

import (
	"fmt"
	"math/rand"
	"time"

	"github.com/cheggaaa/pb/v3"
	"github.com/spf13/cobra"
)

func newLogsCommand() *cobra.Command {
	logsCmd := &cobra.Command{
		Use:   "logs",
		Short: "View application logs",
		Long:  `View application logs with optional real-time following.`,
		Run:   logsRun,
	}

	logsCmd.Flags().BoolP("follow", "f", false, "Follow logs in real-time")

	return logsCmd
}

func logsRun(cmd *cobra.Command, args []string) {
	follow, _ := cmd.Flags().GetBool("follow")

	// Show progress
	fmt.Println("🔍 Retrieving logs...")
	bar := pb.StartNew(100)
	for i := 0; i < 100; i++ {
		bar.Increment()
		time.Sleep(10 * time.Millisecond)
	}
	bar.Finish()

	// Sample log entries
	logEntries := []string{
		"[2023-01-01 10:00:00] [INFO] Application started",
		"[2023-01-01 10:00:01] [DEBUG] Initializing server manager",
		"[2023-01-01 10:00:02] [INFO] Loading configuration from file",
		"[2023-01-01 10:01:15] [INFO] Server USA-NY added successfully",
		"[2023-01-01 10:02:30] [DEBUG] Testing server connectivity",
		"[2023-01-01 10:02:31] [INFO] Server test completed - ping: 45ms",
		"[2023-01-01 10:05:22] [INFO] Connecting to server USA-NY",
		"[2023-01-01 10:05:23] [INFO] Successfully connected to USA-NY",
		"[2023-01-01 10:30:45] [DEBUG] Updating connection statistics",
		"[2023-01-01 11:15:10] [WARNING] High latency detected (120ms)",
		"[2023-01-01 12:00:00] [INFO] Disconnecting from server",
		"[2023-01-01 12:00:01] [INFO] Disconnected successfully",
	}

	// Display logs
	fmt.Println("\n📜 Application Logs")
	fmt.Println("===================")
	for _, entry := range logEntries {
		fmt.Println(entry)
	}

	// Follow logs if requested
	if follow {
		fmt.Println("\n🔄 Following logs (press Ctrl+C to stop)...")
		ticker := time.NewTicker(2 * time.Second)
		defer ticker.Stop()

		i := 0
		logTypes := []string{"INFO", "DEBUG", "WARNING", "ERROR"}
		messages := []string{
			"Background task running",
			"Statistics updated",
			"Checking server status",
			"Cache refreshed",
			"Memory usage: 45MB",
		}

		for {
			select {
			case <-ticker.C:
				timestamp := time.Now().Format("2006-01-02 15:04:05")
				logType := logTypes[rand.Intn(len(logTypes))]
				message := messages[rand.Intn(len(messages))]
				fmt.Printf("[%s] [%s] %s\n", timestamp, logType, message)
				i++
				if i > 5 {
					return
				}
			}
		}
	}
}