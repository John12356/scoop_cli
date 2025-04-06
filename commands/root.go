package commands

import (
	"bufio"
	"fmt"
	"os"
	"os/signal"
	"runtime"
	"strings"
	"syscall"

	"github.com/c-bata/go-prompt"
	"github.com/spf13/cobra"
	"golang.org/x/term"
)

const (
	Version = "1.0.0"
)
type SessionConfig struct {
	URL           string
	VerifyCert    bool
	URLSet        bool
	VerifyCertSet bool
}

var (
	sessionAuthToken string
	sessionConfig    = &SessionConfig{}
)

var RootCmd = &cobra.Command{
	Use:   "securden-cli",
	Short: "A Command-Line-Interface for Securden APIs",
	Long:  "A command-line interface to interact with Securden APIs.",
	Version: Version,
	Run: func(cmd *cobra.Command, args []string) {
		startREPL()
	},
}

func init() {
	RootCmd.AddCommand(createGetPasswordCmd())
	RootCmd.AddCommand(createConfigCmd())
	RootCmd.AddCommand(createClearCmd())
	RootCmd.AddCommand(createClearConfigCmd())
	templateFunc := func() string {
		return fmt.Sprintf("Securden CLI version %s\nGo version: %s\nOS/Arch: %s/%s\n",
			Version,
			runtime.Version(),
			runtime.GOOS,
			runtime.GOARCH)
	}
	RootCmd.SetVersionTemplate(templateFunc()) //version template as a flag
}

func startREPL() {
	// Save and restore terminal state
	originalState, err := term.GetState(int(os.Stdin.Fd()))
	if err != nil {
		fmt.Printf("Warning: Couldn't save terminal state: %v\n", err)
	}
	defer func() {
		if originalState != nil {
			term.Restore(int(os.Stdin.Fd()), originalState)
		}
	}()

	cleanup := func() {
		sessionAuthToken = "" // Clear the session token
		sessionConfig = &SessionConfig{}
	}

	// Handle interrupts
	sigChan := make(chan os.Signal, 1)
	signal.Notify(sigChan, syscall.SIGINT, syscall.SIGTERM)

	fmt.Println("Welcome to Securden CLI (1.0.0). Type '--help' for available commands or 'exit' to quit.")

	if runtime.GOOS == "windows" {
		startWindowsREPL(cleanup, sigChan)
	} else {
		startUnixREPL(cleanup, sigChan)
	}
}

func startWindowsREPL(cleanup func(), sigChan chan os.Signal) {
	reader := bufio.NewReader(os.Stdin)
	exitChan := make(chan struct{})

	go func() {
		<-sigChan
		fmt.Println("\nExiting Securden CLI.")
		cleanup()
		close(exitChan)
	}()

	for {
		select {
		case <-exitChan:
			return
		default:
			fmt.Print(">>> ")
			line, err := reader.ReadString('\n')
			if err != nil {
				fmt.Println("\nExiting Securden CLI.")
				return
			}

			line = strings.TrimSpace(line)
			if line == "" {
				continue
			}

			if line == "exit" || line == "quit" {
				fmt.Println("Exiting Securden CLI.")
				cleanup()
				return
			}

			args, err := parseCommandLine(line)
			if err != nil {
				fmt.Printf("Error parsing command: %v\n", err)
				continue
			}

			cmd := createRootCmd()
			cmd.SetArgs(args)

			if err := cmd.Execute(); err != nil {
				fmt.Printf("Error: %v\n", err)
			}
		}
	}
}

func startUnixREPL(cleanup func(), sigChan chan os.Signal) {
	exitChan := make(chan struct{})

	executor := func(input string) {
		input = strings.TrimSpace(input)
		if input == "exit" || input == "quit" {
			fmt.Println("Exiting Securden CLI.")
			cleanup()
			close(exitChan)
			return
		}

		if input != "" {
			args, err := parseCommandLine(input)
			if err != nil {
				fmt.Printf("Error parsing command: %v\n", err)
				return
			}

			cmd := createRootCmd()
			cmd.SetArgs(args)

			if err := cmd.Execute(); err != nil {
				fmt.Printf("Error: %v\n", err)
			}
		}
	}

	go func() {
		p := prompt.New(
			executor,
			func(d prompt.Document) []prompt.Suggest { return nil }, // We can add suggestion logic here if needed in future
			prompt.OptionPrefix(">>> "),
			prompt.OptionAddKeyBind(prompt.KeyBind{
				Key: prompt.ControlC,
				Fn: func(*prompt.Buffer) {
					fmt.Println("\nExiting Securden CLI.")
					cleanup()
					close(exitChan)
				},
			}),
			prompt.OptionInputTextColor(prompt.Yellow),
            prompt.OptionPrefixTextColor(prompt.Blue),
		)
		p.Run()
	}()

	select {
		case <-exitChan:
		case <-sigChan:
			fmt.Println("\nExiting Securden CLI.")
			cleanup()
	}
}

func createRootCmd() *cobra.Command {
	rootCmd := &cobra.Command{
		Use:   "securden-cli",
		Short: "A Command-Line-Interface for Securden APIs",
	}
	rootCmd.AddCommand(createGetPasswordCmd())
	rootCmd.AddCommand(createConfigCmd())
	rootCmd.AddCommand(createClearCmd())
	rootCmd.AddCommand(createClearConfigCmd())
	templateFunc := func() string {
		return fmt.Sprintf("Securden CLI version %s\nGo version: %s\nOS/Arch: %s/%s\n",
			Version,
			runtime.Version(),
			runtime.GOOS,
			runtime.GOARCH)
	}
	rootCmd.SetVersionTemplate(templateFunc())
	return rootCmd
}

// for input values that has spaces or special characters and wrapped in quotes
// or escaped with backslash, this function will parse the command line input
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