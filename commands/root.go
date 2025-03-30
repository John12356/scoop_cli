// package commands

// import (
// 	"bufio"
// 	"fmt"
// 	"os"
// 	"strings"

// 	"github.com/spf13/cobra"
// 	"github.com/spf13/pflag"
// )

// var RootCmd = &cobra.Command{
// 	Use:   "securden-cli",
// 	Short: "Securden CLI for API interactions",
// 	Long:  "A command-line interface to interact with Securden APIs.",
// 	Run: func(cmd *cobra.Command, args []string) {
// 		fmt.Println("Welcome to Securden CLI. Type 'help' to see available commands or 'exit' to quit.")
// 		securdenTemplate(cmd)
// 	},
// }

// func securdenTemplate(cmd *cobra.Command) {
// 	reader := bufio.NewReader(os.Stdin)
// 	for {
// 		fmt.Print(">>> ")
// 		input, _ := reader.ReadString('\n')
// 		input = strings.TrimSpace(input)

// 		// Exit the REPL
// 		if input == "exit" || input == "quit" {
// 			fmt.Println("Exiting Securden CLI.")
// 			break
// 		}

// 		// Execute the command
// 		if input != "" {
// 			args := strings.Split(input, " ")

// 			// Check if the -h flag is present
// 			if containsHelpFlag(args) {
// 				cmd.SetArgs(args)
// 				cmd.Help()
// 				continue // Skip further processing
// 			}

// 			// Create a fresh command instance for each execution
// 			newCmd := *cmd // Copy the root command
// 			newCmd.SetArgs(args)

// 			// Reset flags and command state
// 			newCmd.Flags().VisitAll(func(f *pflag.Flag) {
// 				f.Changed = false // Reset "changed" state
// 				if err := f.Value.Set(f.DefValue); err != nil {
// 					fmt.Printf("Error resetting flag %s: %v\n", f.Name, err)
// 				}
// 			})

// 			// Execute the command
// 			if err := newCmd.Execute(); err != nil {
// 				fmt.Printf("Error: %s\n", err)
// 			}
// 		}
// 	}
// }

// // Helper function to check if the -h flag is present
// func containsHelpFlag(args []string) bool {
// 	for _, arg := range args {
// 		if arg == "-h" || arg == "--help" {
// 			return true
// 		}
// 	}
// 	return false
// }

package commands

import (
	"fmt"
	"os"
	"strings"

	"github.com/c-bata/go-prompt"
	"github.com/spf13/cobra"
	"github.com/spf13/pflag"
)

var RootCmd = &cobra.Command{
	Use:   "securden-cli",
	Short: "Securden CLI for API interactions",
	Long:  "A command-line interface to interact with Securden APIs.",
	Run: func(cmd *cobra.Command, args []string) {
		fmt.Println("Welcome to Securden CLI. Type 'help' to see available commands or 'exit' to quit.")
		startREPL(cmd)
	},
}

func startREPL(cmd *cobra.Command) {
	// Define the executor function for go-prompt
	executor := func(input string) {
		input = strings.TrimSpace(input)

		// Exit the REPL
		if input == "exit" || input == "quit" {
			fmt.Println("Exiting Securden CLI.")
			os.Exit(0) // Terminate the program
		}

		// Execute the command
		if input != "" {
			args := strings.Split(input, " ")

			// Check if the -h flag is present
			if containsHelpFlag(args) {
				cmd.SetArgs(args)
				cmd.Help()
				return
			}

			// Create a fresh command instance for each execution
			newCmd := *cmd // Copy the root command
			newCmd.SetArgs(args)

			// Reset flags and command state
			newCmd.Flags().VisitAll(func(f *pflag.Flag) {
				f.Changed = false // Reset "changed" state
				if err := f.Value.Set(f.DefValue); err != nil {
					fmt.Printf("Error resetting flag %s: %v\n", f.Name, err)
				}
			})

			// Execute the command
			if err := newCmd.Execute(); err != nil {
				fmt.Printf("Error: %s\n", err)
			}
		}
	}

	// Define the completer function for go-prompt (optional)
	completer := func(d prompt.Document) []prompt.Suggest {
		return []prompt.Suggest{} // Add autocomplete suggestions here
	}

	// Start the REPL with go-prompt
	p := prompt.New(
		executor,
		completer,
		prompt.OptionPrefix(">>> "),
		prompt.OptionTitle("securden-cli"),
	)
	p.Run()
}

// Helper function to check if the -h flag is present
func containsHelpFlag(args []string) bool {
	for _, arg := range args {
		if arg == "-h" || arg == "--help" {
			return true
		}
	}
	return false
}