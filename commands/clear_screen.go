package commands

import (
	"fmt"
	"os"
	"os/exec"
	"runtime"

	"github.com/spf13/cobra"
)

var clearCmd = &cobra.Command{
	Use:   "clear",
	Short: "Clears the terminal screen",
	Run: func(cmd *cobra.Command, args []string) {
		ClearScreen()
	},
}

func ClearScreen() {
	var cmd *exec.Cmd

	switch runtime.GOOS {
	case "windows":
		cmd = exec.Command("cmd", "/c", "cls")
	default: // Linux, macOS, etc.
		cmd = exec.Command("clear")
	}

	cmd.Stdout = os.Stdout
	err := cmd.Run()
	if err != nil {
		fmt.Println("Failed to clear the screen:", err)
	}
}

func init() {
	RootCmd.AddCommand(clearCmd)
}
