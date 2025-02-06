package main

import (
	"fmt"
	"os"

	"securden_cli/commands"
)

func init() {
	installPath, exists := os.LookupEnv("SECURDEN_CLI_PATH")
	if !exists {
		fmt.Println("Warning: Install path not found, using current directory")
		return
	}

	err := os.Chdir(installPath)
	if err != nil {
		fmt.Println("Warning: Failed to set working directory:", err)
	}
}


func main() {
	if err := commands.RootCmd.Execute(); err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
}





