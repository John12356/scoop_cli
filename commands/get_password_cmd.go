package commands

import (
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"net/url"
	"strconv"
	"strings"

	"github.com/spf13/cobra"
)

func createGetPasswordCmd() *cobra.Command {
	var cmd = &cobra.Command{
		Use:   "get-password",
		Short: "Retrieves the password from the Securden server",
		Run: func(cmd *cobra.Command, args []string) {
			fileCfg, err := loadConfig()
			fileExists := err == nil && fileCfg.URL != ""
			
			// URL Resolution
			var finalURL string
			urlProvided := cmd.Flags().Changed("url")
			
			if urlProvided {
				// 1. Command line has highest priority
				finalURL, _ = cmd.Flags().GetString("url")
				sessionConfig.URL = finalURL
				sessionConfig.URLSet = true
			} else if fileExists {
				// 2. Use file config if available
				finalURL = fileCfg.URL
			} else if sessionConfig.URLSet {
				// 3. Fallback to session if set
				finalURL = sessionConfig.URL
			} else {
				fmt.Println("URL field cannot be left empty. Specify URL via flags or configure it using the config command.")
				return
			}

			// VerifyCert Resolution (default true)
			var finalVerifyCert bool = true
			verifyCertProvided := cmd.Flags().Changed("verifycert")
			
			if verifyCertProvided {
				verifyCertStr, _ := cmd.Flags().GetString("verifycert")
				finalVerifyCert, parseErr := strconv.ParseBool(verifyCertStr)
				if parseErr != nil {
					fmt.Println("Invalid value for --verifycert. The value must be either 'true' or 'false'.")
					return
				}
				sessionConfig.VerifyCert = finalVerifyCert
				sessionConfig.VerifyCertSet = true
			} else if fileExists {
				finalVerifyCert = fileCfg.VerifyCert
			} else if sessionConfig.VerifyCertSet {
				finalVerifyCert = sessionConfig.VerifyCert
			}

			// AuthToken is always session-persistent
			var AuthToken string
			if cmd.Flags().Changed("authtoken") {
				AuthToken, _ = cmd.Flags().GetString("authtoken")
				sessionAuthToken = AuthToken
			} else if sessionAuthToken != "" {
				AuthToken = sessionAuthToken
			} else {
				fmt.Println("Please provide the authentication token via --authtoken flag at least once per session.")
				return
			}

			accountID, _ := cmd.Flags().GetString("account-id")
			accountName, _ := cmd.Flags().GetString("account-name")
			accountType, _ := cmd.Flags().GetString("account-type")
			accountTitle, _ := cmd.Flags().GetString("account-title")
			accountCategory, _ := cmd.Flags().GetString("account-category")
			ticketID, _ := cmd.Flags().GetString("ticket-id")
			reason, _ := cmd.Flags().GetString("reason")

			// account-category validation
			if accountCategory != "" {
				categoryValue := strings.ToLower(accountCategory)
				switch categoryValue {
				case "work", "personal":
					//if found to be valid value --> continue
				default:
					fmt.Println("Error: Invalid value for --account-category. The value must be either 'work' or 'personal'.")
					return
				}
			}

			cfg := Config{
                URL:        finalURL,
                VerifyCert: finalVerifyCert,
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
			if accountCategory != "" {
				categoryValue := strings.ToLower(accountCategory)
				switch categoryValue {
				case "work":
					query.Set("account_category", "1")
				case "personal":
					query.Set("account_category", "3")
				}
			}
			if ticketID != "" {
				query.Set("ticket_id", ticketID)
			}
			if reason != "" {
				query.Set("reason", reason)
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
				fmt.Println("Error calling APIi:", err)
				return
			}
			defer response.Body.Close()
			body, err := io.ReadAll(response.Body)
			if err != nil {
				fmt.Println("Error reading response:", err)
				return
			}
			
			if response.StatusCode == http.StatusForbidden {
				type ErrorResponse struct {
					StatusCode int `json:"status_code"`
					Error struct {
						Code    string `json:"code"`
						Message string `json:"message"`
					} `json:"error"`
				}

				// if response has message field inside error field, parse it
				var errorResp ErrorResponse
				err = json.Unmarshal(body, &errorResp)
				if err == nil && errorResp.Error.Message != "" {
					fmt.Println("Error:", errorResp.Error.Message)
					return
				}

				// if response has message field, print it
				var simpleResp PasswordResponse
				if err := json.Unmarshal(body, &simpleResp); err == nil && simpleResp.Message != "" {
					fmt.Println("Error:", simpleResp.Message)
					return
				}

				fmt.Println("Error: Unauthorized Access")
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

	cmd.Flags().StringP("account-id", "i", "", "Account ID for the API call")
	cmd.Flags().StringP("account-name", "n", "", "Account Name for the API call")
	cmd.Flags().String("account-type", "", "Account Type for the API call")
	cmd.Flags().StringP("account-title", "t", "", "Account Title for the API call")
	cmd.Flags().StringP("account-category", "c", "", "Account Category for the API call (work, personal)")
	cmd.Flags().String("ticket-id", "", "Ticket ID for the request")
	cmd.Flags().String("reason", "", "Reason for accessing the password")
	cmd.Flags().String("url", "", "Server URL (if not configured)")
	cmd.Flags().String("authtoken", "", "Authentication token (if not configured)")
	cmd.Flags().String("verifycert", "true", "Enable SSL certificate verification (if not configured)")
	cmd.MarkFlagsOneRequired("account-id", "account-name", "account-title")

	return cmd
}

type PasswordResponse struct {
	Password string `json:"password"`
	Message  string `json:"message"`
}