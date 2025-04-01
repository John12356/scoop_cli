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

// package commands

// import (
// 	"fmt"
// 	"os"
// 	"strings"

// 	"github.com/c-bata/go-prompt"
// 	"github.com/spf13/cobra"
// 	"github.com/spf13/pflag"
// )

// var RootCmd = &cobra.Command{
// 	Use:   "securden-cli",
// 	Short: "Securden CLI for API interactions",
// 	Long:  "A command-line interface to interact with Securden APIs.",
// 	Run: func(cmd *cobra.Command, args []string) {
// 		fmt.Println("Welcome to Securden CLI. Type 'help' to see available commands or 'exit' to quit.")
// 		startREPL(cmd)
// 	},
// }

// func startREPL(cmd *cobra.Command) {
// 	// Define the executor function for go-prompt
// 	executor := func(input string) {
// 		input = strings.TrimSpace(input)

// 		// Exit the REPL
// 		if input == "exit" || input == "quit" {
// 			fmt.Println("Exiting Securden CLI.")
// 			os.Exit(0)
// 		}

// 		// Execute the command
// 		if input != "" {
// 			args := strings.Split(input, " ")

// 			// Check if the -h flag is present
// 			if containsHelpFlag(args) {
// 				cmd.SetArgs(args)
// 				cmd.Help()
// 				return
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

// 	// Define the completer function for go-prompt (optional)
// 	completer := func(d prompt.Document) []prompt.Suggest {
// 		return []prompt.Suggest{} // Add autocomplete suggestions here
// 	}

// 	// Start the REPL with go-prompt
// 	p := prompt.New(
// 		executor,
// 		completer,
// 		prompt.OptionPrefix(">>> "),
// 		prompt.OptionTitle("securden-cli"),
// 	)
// 	p.Run()
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

// package commands

// import (
// 	"fmt"
// 	"os"
// 	"os/signal"
// 	"strings"
// 	"syscall"

// 	"github.com/c-bata/go-prompt"
// 	"github.com/spf13/cobra"
// 	"github.com/spf13/pflag"
// 	"golang.org/x/term"
// )

// var RootCmd = &cobra.Command{
// 	Use:   "securden-cli",
// 	Short: "Securden CLI for API interactions",
// 	Long:  "A command-line interface to interact with Securden APIs.",
// 	Run: func(cmd *cobra.Command, args []string) {
// 		fmt.Println("Welcome to Securden CLI. Type 'help' to see available commands or 'exit' to quit.")
// 		startREPL(cmd)
// 	},
// }

// func startREPL(cmd *cobra.Command) {
// 	// Save original terminal state
// 	originalState, err := term.GetState(int(os.Stdin.Fd()))
// 	if err != nil {
// 		fmt.Printf("Warning: Couldn't save terminal state: %v\n", err)
// 	}

// 	// Setup signal handling for proper cleanup
// 	sigChan := make(chan os.Signal, 1)
// 	signal.Notify(sigChan, syscall.SIGINT, syscall.SIGTERM)

// 	// Channel to signal when we want to exit
// 	exitChan := make(chan struct{})

// 	// Define the executor function
// 	executor := func(input string) {
// 		input = strings.TrimSpace(input)

// 		if input == "exit" || input == "quit" {
// 			fmt.Println("Exiting Securden CLI.")
// 			close(exitChan)
// 			return
// 		}

// 		if input != "" {
// 			args := strings.Split(input, " ")

// 			if containsHelpFlag(args) {
// 				cmd.SetArgs(args)
// 				cmd.Help()
// 				return
// 			}

// 			newCmd := *cmd
// 			newCmd.SetArgs(args)

// 			newCmd.Flags().VisitAll(func(f *pflag.Flag) {
// 				f.Changed = false
// 				if err := f.Value.Set(f.DefValue); err != nil {
// 					fmt.Printf("Error resetting flag %s: %v\n", f.Name, err)
// 				}
// 			})

// 			if err := newCmd.Execute(); err != nil {
// 				fmt.Printf("Error: %s\n", err)
// 			}
// 		}
// 	}

// 	// Start the prompt in a goroutine
// 	go func() {
// 		p := prompt.New(
// 			executor,
// 			func(d prompt.Document) []prompt.Suggest { return nil },
// 			prompt.OptionPrefix(">>> "),
// 			prompt.OptionTitle("securden-cli"),
// 			prompt.OptionAddKeyBind(prompt.KeyBind{
// 				Key: prompt.ControlC,
// 				Fn: func(*prompt.Buffer) {
// 					close(exitChan)
// 				},
// 			}),
// 		)
// 		p.Run()
// 	}()

// 	// Wait for exit signal
// 	select {
// 	case <-exitChan:
// 	case <-sigChan:
// 		fmt.Println("\nExiting Securden CLI.")
// 	}

// 	// Restore terminal state
// 	if originalState != nil {
// 		if err := term.Restore(int(os.Stdin.Fd()), originalState); err != nil {
// 			fmt.Printf("Warning: Couldn't restore terminal state: %v\n", err)
// 		}
// 	}
// }

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
	"os/signal"
	"strings"
	"syscall"

	"github.com/c-bata/go-prompt"
	"github.com/spf13/cobra"
	"github.com/spf13/pflag"
	"golang.org/x/term"
)

var RootCmd = &cobra.Command{
	Use:   "securden-cli",
	Short: "CLI for Securden APIs",
	Long:  "A command-line interface to interact with Securden APIs.",
	Run: func(cmd *cobra.Command, args []string) {
		fmt.Println("Welcome to Securden CLI. Type '--help' to see the list of available commands or 'exit' to quit.")
		startREPL(cmd)
	},
}

func startREPL(cmd *cobra.Command) {
	// Save original terminal state
	originalState, err := term.GetState(int(os.Stdin.Fd()))
	if err != nil {
		fmt.Printf("Warning: Couldn't save terminal state: %v\n", err)
	}

	// Setup signal handling for proper cleanup
	sigChan := make(chan os.Signal, 1)
	signal.Notify(sigChan, syscall.SIGINT, syscall.SIGTERM)

	// Channel to signal when we want to exit
	exitChan := make(chan struct{})

	// Define the executor function with proper argument parsing
	executor := func(input string) {
		input = strings.TrimSpace(input)
		if input == "exit" || input == "quit" {
			fmt.Println("Exiting Securden CLI.")
			close(exitChan)
			return
		}

		if input != "" {
			// Parse the input with proper quote handling
			args, err := parseCommandLine(input)
			if err != nil {
				fmt.Printf("Error parsing command: %v\n", err)
				return
			}

			if containsHelpFlag(args) {
				cmd.SetArgs(args)
				cmd.Help()
				return
			}

			// Create a new command instance to avoid flag pollution
			newCmd := *cmd
			newCmd.SetArgs(args)

			// Reset all flags to their default values
			newCmd.Flags().VisitAll(func(f *pflag.Flag) {
				f.Changed = false
				if err := f.Value.Set(f.DefValue); err != nil {
					fmt.Printf("Error resetting flag %s: %v\n", f.Name, err)
				}
			})

			if err := newCmd.Execute(); err != nil {
				fmt.Printf("Error: %v\n", err)
			}
		}
	}

	// Start the prompt in a goroutine
	go func() {
		p := prompt.New(
			executor,
			func(d prompt.Document) []prompt.Suggest { return nil },
			prompt.OptionPrefix(">>> "),
			prompt.OptionTitle("securden-cli"),
			prompt.OptionAddKeyBind(prompt.KeyBind{
				Key: prompt.ControlC,
				Fn: func(*prompt.Buffer) {
					close(exitChan)
				},
			}),
		)
		p.Run()
	}()

	// Wait for exit signal
	select {
	case <-exitChan:
	case <-sigChan:
		fmt.Println("\nExiting Securden CLI.")
	}

	// Restore terminal state
	if originalState != nil {
		if err := term.Restore(int(os.Stdin.Fd()), originalState); err != nil {
			fmt.Printf("Warning: Couldn't restore terminal state: %v\n", err)
		}
	}
}

// parseCommandLine properly splits a command string into arguments, handling quotes
func parseCommandLine(input string) ([]string, error) {
	var args []string
	var currentArg strings.Builder
	inQuotes := false
	quoteChar := rune(0)
	escapeNext := false

	for _, r := range input {
		if escapeNext {
			currentArg.WriteRune(r)
			escapeNext = false
			continue
		}

		switch r {
		case '\\':
			escapeNext = true
		case '"', '\'':
			if inQuotes {
				if r == quoteChar {
					inQuotes = false
					quoteChar = 0
				} else {
					currentArg.WriteRune(r)
				}
			} else {
				inQuotes = true
				quoteChar = r
			}
		case ' ':
			if inQuotes {
				currentArg.WriteRune(r)
			} else if currentArg.Len() > 0 {
				args = append(args, currentArg.String())
				currentArg.Reset()
			}
		default:
			currentArg.WriteRune(r)
		}
	}

	if currentArg.Len() > 0 {
		args = append(args, currentArg.String())
	}

	if inQuotes {
		return nil, fmt.Errorf("unclosed quotes in command")
	}

	return args, nil
}

func containsHelpFlag(args []string) bool {
	for _, arg := range args {
		if arg == "-h" || arg == "--help" {
			return true
		}
	}
	return false
}