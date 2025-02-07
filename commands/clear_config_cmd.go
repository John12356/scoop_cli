package commands

import (
	"fmt"
	"os"
	"path/filepath"

	"github.com/spf13/cobra"
)

var clearConfigCmd = &cobra.Command{
	Use:   "clear-config",
	Short: "Clear the saved server configuration",
	Run: func(cmd *cobra.Command, args []string) {
		configPath, err := getConfigFilePath()
		if err != nil {
			fmt.Println("Error getting config path:", err)
			return
		}

		err = os.Remove(configPath)
		if err != nil {
			if os.IsNotExist(err) {
				fmt.Println("No configuration found to clear.")
			} else {
				fmt.Printf("Error clearing configuration: %v\n", err)
			}
			return
		}
		
		// Optional: Remove the directory if empty
		configDir := filepath.Dir(configPath)
		if err := os.Remove(configDir); err != nil && !os.IsNotExist(err) {
			fmt.Printf("Note: Could not remove config directory: %v\n", err)
		}

		fmt.Println("Configuration cleared successfully.")
	},
}

func init() {
	RootCmd.AddCommand(clearConfigCmd)
}