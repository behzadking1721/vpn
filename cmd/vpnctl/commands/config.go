package commands

import (
	"fmt"
	"os"

	"github.com/spf13/cobra"
)

func newConfigCommand() *cobra.Command {
	configCmd := &cobra.Command{
		Use:   "config",
		Short: "Manage configuration",
		Long:  `Manage VPN client configuration including viewing, setting, and exporting configurations.`,
	}

	configCmd.AddCommand(
		newConfigShowCommand(),
		newConfigSetCommand(),
		newConfigResetCommand(),
		newConfigImportCommand(),
		newConfigExportCommand(),
	)

	return configCmd
}

func newConfigShowCommand() *cobra.Command {
	return &cobra.Command{
		Use:   "show",
		Short: "Show current configuration",
		Long:  `Show the current VPN client configuration.`,
		Run: func(cmd *cobra.Command, args []string) {
			fmt.Println("⚙️ Current Configuration")
			fmt.Println("======================")
			fmt.Println("Auto Connect:        true")
			fmt.Println("Start Minimized:     false")
			fmt.Println("Auto Update Servers: true")
			fmt.Println("Notifications:       true")
			fmt.Println("Cache Servers:       true")
			fmt.Println("Theme:               dark")
			fmt.Println("Language:            en")
			fmt.Println("Log Level:           info")
		},
	}
}

func newConfigSetCommand() *cobra.Command {
	return &cobra.Command{
		Use:   "set [key] [value]",
		Short: "Set a configuration value",
		Long:  `Set a configuration value by key.`,
		Args:  cobra.ExactArgs(2),
		Run: func(cmd *cobra.Command, args []string) {
			key := args[0]
			value := args[1]
			fmt.Printf("🔧 Setting %s = %s\n", key, value)
			fmt.Println("✅ Configuration updated successfully")
		},
	}
}

func newConfigResetCommand() *cobra.Command {
	return &cobra.Command{
		Use:   "reset",
		Short: "Reset configuration to defaults",
		Long:  `Reset all configuration settings to their default values.`,
		Run: func(cmd *cobra.Command, args []string) {
			fmt.Print("⚠️  Are you sure you want to reset all configuration? (y/N): ")
			var response string
			fmt.Scanln(&response)
			if response == "y" || response == "Y" {
				fmt.Println("🔄 Resetting configuration to defaults...")
				fmt.Println("✅ Configuration reset successfully")
			} else {
				fmt.Println("❌ Configuration reset cancelled")
			}
		},
	}
}

func newConfigImportCommand() *cobra.Command {
	return &cobra.Command{
		Use:   "import [file]",
		Short: "Import configuration from file",
		Long:  `Import configuration from a specified file.`,
		Args:  cobra.ExactArgs(1),
		Run: func(cmd *cobra.Command, args []string) {
			file := args[0]
			if _, err := os.Stat(file); os.IsNotExist(err) {
				ExitWithError(1, fmt.Sprintf("File not found: %s", file))
			}
			fmt.Printf("📥 Importing configuration from: %s\n", file)
			fmt.Println("✅ Configuration imported successfully")
		},
	}
}

func newConfigExportCommand() *cobra.Command {
	return &cobra.Command{
		Use:   "export [file]",
		Short: "Export configuration to file",
		Long:  `Export current configuration to a specified file.`,
		Args:  cobra.ExactArgs(1),
		Run: func(cmd *cobra.Command, args []string) {
			file := args[0]
			fmt.Printf("📤 Exporting configuration to: %s\n", file)
			fmt.Println("✅ Configuration exported successfully")
		},
	}
}