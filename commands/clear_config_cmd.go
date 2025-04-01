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
				fmt.Println("Connectivity details not found.")
			} else {
				fmt.Printf("Error deleting connectivity details: %v\n", err)
			}
			return
		}
		
		// Removing the directory if empty
		configDir := filepath.Dir(configPath)
		if err := os.Remove(configDir); err != nil && !os.IsNotExist(err) {
			fmt.Printf("Note: Could not remove config directory: %v\n", err)
		}

		fmt.Println("Connectivity details deleted successfully.")
	},
}

func init() {
	RootCmd.AddCommand(clearConfigCmd)
}