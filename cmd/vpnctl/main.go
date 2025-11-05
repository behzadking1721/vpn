package main

import (
	"fmt"
	"os"

	"vpnclient/cmd/vpnctl/commands"
)

func main() {
	rootCmd := commands.NewRootCommand()
	
	if err := rootCmd.Execute(); err != nil {
		fmt.Fprintln(os.Stderr, err)
		os.Exit(1)
	}
}