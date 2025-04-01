// package commands

// import (
// 	"encoding/json"
// 	"fmt"
// 	"io"
// 	"net/http"
// 	"net/url"
// 	"strconv"

// 	"github.com/spf13/cobra"
// )

// func init() {
// 	apiCmd.SetHelpCommand(&cobra.Command{Hidden: true})
// 	apiCmd.Flags().BoolP("help", "h", false, "Help for get_password")
// 	apiCmd.Flags().StringP("account-id", "i", "", "Account ID for the API call")
// 	apiCmd.Flags().StringP("account-name", "n", "", "Account Name for the API call")
// 	apiCmd.Flags().StringP("account-type", "t", "", "Account Type for the API call")
// 	apiCmd.Flags().StringP("account-title", "l", "", "Account Title for the API call")
// 	apiCmd.Flags().String("url", "", "Server URL (if not configured)")
// 	apiCmd.Flags().String("authtoken", "", "Authentication token (if not configured)")
// 	apiCmd.Flags().String("verifycert", "true", "Enable SSL certificate verification (if not configured)")
// 	RootCmd.AddCommand(apiCmd)
// }

// type PasswordResponse struct {
// 	Password string `json:"password"`
// 	Message  string `json:"message"`
// }

// var apiCmd = &cobra.Command{
// 	Use:   "get-password",
// 	Short: "Retrieve the password from the Securden server",
// 	Run: func(cmd *cobra.Command, args []string) {
// 		// Defer flag clearing to ensure it happens after execution
// 		defer func() {
// 			cmd.Flags().Set("account-id", "")
// 			cmd.Flags().Set("account-name", "")
// 			cmd.Flags().Set("account-type", "")
// 			cmd.Flags().Set("account-title", "")
// 		}()

// 		// Load existing configuration. If conf has no url, you should pass it out on commands
// 		cfg, err := loadConfig()
// 		if err != nil {
// 			cfg = Config{}
// 			cfg.VerifyCert = true
// 		}

// 		// Get fresh flag values
// 		accountID, _ := cmd.Flags().GetString("account-id")
// 		accountName, _ := cmd.Flags().GetString("account-name")
// 		accountType, _ := cmd.Flags().GetString("account-type")
// 		accountTitle, _ := cmd.Flags().GetString("account-title")

// 		if accountID == "" && accountName == "" && accountTitle == "" {
// 			fmt.Println("Error: At least one of --account-id, --account-name or --account-title must be provided")
// 			return
// 		}

// 		if cmd.Flags().Changed("url") {
// 			cfg.URL, _ = cmd.Flags().GetString("url")
// 		}
// 		AuthToken, _ := cmd.Flags().GetString("authtoken")
// 		if cmd.Flags().Changed("verifycert") {
// 			verifyCertStr, _ := cmd.Flags().GetString("verifycert")
// 			verifyCert, parseErr := strconv.ParseBool(verifyCertStr)
// 			if parseErr != nil {
// 				fmt.Println("Invalid value for --verifycert. Use 'true' or 'false'.")
// 				return
// 			}
// 			cfg.VerifyCert = verifyCert
// 		}

// 		if cfg.URL == "" || AuthToken == "" {
// 			if cfg.URL == "" {
// 				fmt.Println("URL field required. Provide it via flags or configure it using the config command.")
// 				return
// 			}
// 			fmt.Println("Please provide the authentication token via flags at least one time for a session.")
// 			return
// 		}

// 		client, err := getHTTPClient(cfg)
// 		if err != nil {
// 			fmt.Println("Error setting up HTTP client:", err)
// 			return
// 		}

// 		u, err := url.Parse(cfg.URL + "/secretsmanagement/get_password_via_tools")
// 		if err != nil {
// 			fmt.Println("Error parsing URL:", err)
// 			return
// 		}

// 		query := u.Query()
// 		if accountID != "" {
// 			query.Set("account_id", accountID)
// 		}
// 		if accountName != "" {
// 			query.Set("account_name", accountName)
// 		}
// 		if accountType != "" {
// 			query.Set("account_type", accountType)
// 		}
// 		if accountTitle != "" {
// 			query.Set("account_title", accountTitle)
// 		}
// 		u.RawQuery = query.Encode()

// 		req, err := http.NewRequest("GET", u.String(), nil)
// 		if err != nil {
// 			fmt.Println("Error creating request:", err)
// 			return
// 		}
// 		req.Header.Add("authtoken", AuthToken)
// 		response, err := client.Do(req)
// 		if err != nil {
// 			fmt.Println("Error calling API:", err)
// 			return
// 		}
// 		defer response.Body.Close()
// 		body, err := io.ReadAll(response.Body)
// 		if err != nil {
// 			fmt.Println("Error reading response:", err)
// 			return
// 		}

// 		if response.StatusCode == http.StatusForbidden {
// 			var errorResp PasswordResponse
// 			err = json.Unmarshal(body, &errorResp)
// 			if err != nil {
// 				fmt.Println("Error parsing error response:", err)
// 				return
// 			}
// 			fmt.Println("Error:", errorResp.Message)
// 			return
// 		}

// 		if response.StatusCode != http.StatusOK {
// 			fmt.Printf("Error: Received status code %d from server\n", response.StatusCode)
// 			return
// 		}

// 		var passwordResp PasswordResponse
// 		err = json.Unmarshal(body, &passwordResp)
// 		if err != nil {
// 			fmt.Println("Error parsing JSON response:", err)
// 			return
// 		}
// 		fmt.Println(passwordResp.Password)
// 	},
// }

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
	// apiCmd.SetHelpCommand(&cobra.Command{Hidden: true})
	// apiCmd.Flags().BoolP("help", "h", false, "Help for get_password")
	apiCmd.Flags().StringP("account-id", "i", "", "Account ID for the API call")
	apiCmd.Flags().StringP("account-name", "n", "", "Account Name for the API call")
	apiCmd.Flags().StringP("account-type", "t", "", "Account Type for the API call")
	apiCmd.Flags().StringP("account-title", "l", "", "Account Title for the API call")
	apiCmd.Flags().String("url", "", "Server URL (if not configured)")
	apiCmd.Flags().String("authtoken", "", "Authentication token (if not configured)")
	apiCmd.Flags().String("verifycert", "true", "Enable SSL certificate verification (if not configured)")
	apiCmd.MarkFlagsOneRequired("account-id", "account-name", "account-title")
	RootCmd.AddCommand(apiCmd)
}

type PasswordResponse struct {
	Password string `json:"password"`
	Message  string `json:"message"`
}

var apiCmd = &cobra.Command{
	Use:   "get-password",
	Short: "Retrieves the password from the Securden server",
	Run: func(cmd *cobra.Command, args []string) {
		if len(args) == 1 && args[0] == "help" {
			cmd.Help()
			return
		}

		//clear the flags after execution to avoid confusion in future commands
		defer func() {
			cmd.Flags().Set("account-id", "")
			cmd.Flags().Set("account-name", "")
			cmd.Flags().Set("account-type", "")
			cmd.Flags().Set("account-title", "")
		}()
		// Load existing configuration, if not found, we have pass it manually in evrey command
		cfg, err := loadConfig()
		if err != nil {
			cfg = Config{}
			cfg.VerifyCert = true
		}

		// Get flag values
		accountID, _ := cmd.Flags().GetString("account-id")
		accountName, _ := cmd.Flags().GetString("account-name")
		accountType, _ := cmd.Flags().GetString("account-type")
		accountTitle, _ := cmd.Flags().GetString("account-title")

		if cmd.Flags().Changed("url") {
			cfg.URL, _ = cmd.Flags().GetString("url")
		}
		AuthToken, _ := cmd.Flags().GetString("authtoken")
		if cmd.Flags().Changed("verifycert") {
			verifyCertStr, _ := cmd.Flags().GetString("verifycert")
			verifyCert, parseErr := strconv.ParseBool(verifyCertStr)
			if parseErr != nil {
				fmt.Println("Invalid value for --verifycert. The value must be either 'true' or 'false'.")
				return
			}
			cfg.VerifyCert = verifyCert
		}

		if cfg.URL == "" || AuthToken == "" {
			if cfg.URL == "" {
				fmt.Println("URL field cannot be left empty. Specify URL via flags or configure it using the config command.")
				return
			}
			fmt.Println("Please provide the authentication token via flags at least one time per session.")
			return
		}

		client, err := getHTTPClient(cfg)
		if err != nil {
			fmt.Println("Error setting up HTTP client:", err)
			return
		}

		u, err := url.Parse(cfg.URL + "/secretsmanagement/get_password_via_tools")
		if err != nil {
			fmt.Println("Error parsing URL:", err)
			return
		}

		query := u.Query()
		if accountID != "" {
			query.Set("account_id", accountID)
		}
		if accountName != "" {
			query.Set("account_name", accountName)
		}
		if accountType != "" {
			query.Set("account_type", accountType)
		}
		if accountTitle != "" {
			query.Set("account_title", accountTitle)
		}
		u.RawQuery = query.Encode()

		req, err := http.NewRequest("GET", u.String(), nil)
		if err != nil {
			fmt.Println("Error creating request:", err)
			return
		}
		req.Header.Add("authtoken", AuthToken)
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

		if response.StatusCode == http.StatusForbidden {
			var errorResp PasswordResponse
			err = json.Unmarshal(body, &errorResp)
			if err != nil {
				fmt.Println("Error parsing error response:", err)
				return
			}
			fmt.Println("Error:", errorResp.Message)
			return
		}

		if response.StatusCode != http.StatusOK {
			fmt.Printf("Error: Received status code %d from server\n", response.StatusCode)
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