package commands

import (
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"net/url"
	"strconv"

	"github.com/spf13/cobra"
)

func init() {
	// Disabling the help command for this command since it's causing conflicts & manual help flag is added
	apiCmd.SetHelpCommand(&cobra.Command{Hidden: true})
	apiCmd.Flags().BoolP("help", "h", false, "Help for get_password")
	apiCmd.Flags().StringP("account-id", "a", "", "Account ID for the API call")
	apiCmd.Flags().String("url", "", "Server URL (if not configured)")
	apiCmd.Flags().String("authtoken", "", "Authentication token (if not configured)")
	apiCmd.Flags().String("verifycert", "true", "Enable SSL certificate verification (if not configured)")
	apiCmd.MarkFlagRequired("account-id")
	RootCmd.AddCommand(apiCmd)
}

type PasswordResponse struct {
	Password string `json:"password"`
}

var apiCmd = &cobra.Command{
	Use:   "get-password",
	Short: "Retrieve the password from the Securden server",
	Run: func(cmd *cobra.Command, args []string) {
		// Load existing configuration
		cfg, err := loadConfig()
		if err != nil {
			cfg = Config{} // Initialize an empty config if the file doesn't exist
		}

		// Overriding configured values with command-line flags if provided
		if cmd.Flags().Changed("url") {
			cfg.URL, _ = cmd.Flags().GetString("url")
		}
		if cmd.Flags().Changed("authtoken") {
			cfg.AuthToken, _ = cmd.Flags().GetString("authtoken")
		}
		if cmd.Flags().Changed("verifycert") {
			// Taking verifyCert as string and converting explicitly it to a bool 
			verifyCertStr, _ := cmd.Flags().GetString("verifycert")
			verifyCert, parseErr := strconv.ParseBool(verifyCertStr)
			if parseErr != nil {
				fmt.Println("Invalid value for --verifycert. Use 'true' or 'false'.")
				return
			}
			cfg.VerifyCert = verifyCert
		}

		if cfg.URL == "" || cfg.AuthToken == "" {
			fmt.Println("Both URL and AuthToken are required. Provide them via flags or configure them using the config command.")
			return
		}
		accountID, _ := cmd.Flags().GetString("account-id")

		// Setting up the HTTP client
		client, err := getHTTPClient(cfg)
		if err != nil {
			fmt.Println("Error setting up HTTP client:", err)
			return
		}

		u, err := url.Parse(cfg.URL + "/api/get_password")
		if err != nil {
			fmt.Println("Error parsing URL:", err)
			return
		}

		query := u.Query()
		query.Set("account_id", accountID)
		u.RawQuery = query.Encode()

		// Create the HTTP request
		req, err := http.NewRequest("GET", u.String(), nil)
		if err != nil {
			fmt.Println("Error creating request:", err)
			return
		}
		req.Header.Add("authtoken", cfg.AuthToken)
		response, err := client.Do(req)
		if err != nil {
			fmt.Println("Error calling API:", err)
			return
		}
		defer response.Body.Close()
		body, err := io.ReadAll(response.Body)
		if err != nil {
			fmt.Println("Error reading response:", err)
			return
		}
		var passwordResp PasswordResponse
		err = json.Unmarshal(body, &passwordResp)
		if err != nil {
			fmt.Println("Error parsing JSON response:", err)
			return
		}
		fmt.Println(passwordResp.Password)
	},
}