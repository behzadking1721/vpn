package main

import (
	"fmt"
	"os"

	"vpnclient/cmd/vpnctl/commands"
)

// Version information set at build time
var (
	version   = "dev"
	buildTime = "unknown"
	gitCommit = "unknown"
)

func main() {
	// Handle version flag
	for _, arg := range os.Args[1:] {
		if arg == "--version" || arg == "-v" {
			fmt.Printf("VPN CLI Tool %s\nBuild: %s\nCommit: %s\n", version, buildTime, gitCommit)
			return
		}
	}

	rootCmd := commands.NewRootCommand()

	// Add version to the root command
	rootCmd.Version = fmt.Sprintf("%s (build: %s, commit: %s)", version, buildTime, gitCommit)

	if err := rootCmd.Execute(); err != nil {
		fmt.Fprintln(os.Stderr, err)
		os.Exit(1)
	}
}
