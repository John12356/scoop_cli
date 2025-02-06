package commands

import (
	"crypto/tls"
	"crypto/x509"
	"encoding/json"
	"fmt"
	"net/http"
	"net/url"
	"os"
	"path/filepath"
	"strconv"
	"strings"

	"github.com/spf13/cobra"
)

type Config struct {
	URL        string `json:"url"`
	AuthToken  string `json:"authtoken"`
	VerifyCert bool   `json:"verifyCert"`
}

func saveConfig(config Config) error {
    configPath, err := getConfigFilePath()
    if err != nil {
        return err
    }

    if err := os.MkdirAll(filepath.Dir(configPath), 0755); err != nil {
        return err
    }

    file, err := os.Create(configPath)
    if err != nil {
        return err
    }
    defer file.Close()

    encoder := json.NewEncoder(file)
    return encoder.Encode(config)
}
func getConfigFilePath() (string, error) {
    configDir, err := os.UserConfigDir()
    if err != nil {
        return "", err
    }
    return filepath.Join(configDir, "securden-cli", "config.json"), nil
}

func loadConfig() (Config, error) {
    configPath, err := getConfigFilePath()
    if err != nil {
        return Config{}, err
    }

    file, err := os.Open(configPath)
    if err != nil {
        return Config{}, err
    }
    defer file.Close()

    var config Config
    decoder := json.NewDecoder(file)
    err = decoder.Decode(&config)
    return config, err
}

func getHostFromURL(URL string) (string, error){ 
	parsedURL, err := url.Parse(URL)
	if err != nil {
		return "", fmt.Errorf("invalid BaseURL: %v", err)
	}

	host := parsedURL.Host
	if !strings.Contains(host, ":") {
		host = fmt.Sprintf("%s:443", host)
	}
	return host, nil
}

func getHTTPClient(config Config) (*http.Client, error) {
	host, err := getHostFromURL(config.URL)
	if err != nil {
		return nil, fmt.Errorf("invalid BaseURL: %v", err)
	}
	
	if !config.VerifyCert {
		transport := &http.Transport{
			TLSClientConfig: &tls.Config{InsecureSkipVerify: true},
		}
		return &http.Client{Transport: transport}, nil
	}

	// Fetch SSL certificate from the server and use it for verification
	conn, err := tls.Dial("tcp", host, &tls.Config{InsecureSkipVerify: true})
	if err != nil {
		return nil, fmt.Errorf("failed to fetch certificate from server: %v", err)
	}
	defer conn.Close()

	if len(conn.ConnectionState().PeerCertificates) == 0 {
		return nil, fmt.Errorf("no certificates found on the server")
	}

	serverCert := conn.ConnectionState().PeerCertificates[0]
	certPool := x509.NewCertPool()
	certPool.AddCert(serverCert)

	transport := &http.Transport{
		TLSClientConfig: &tls.Config{RootCAs: certPool},
	}
	return &http.Client{Transport: transport}, nil
}

var configCmd = &cobra.Command{
	Use:   "config",
	Short: "Configure server details",
	Run: func(cmd *cobra.Command, args []string) {
		url, _ := cmd.Flags().GetString("url")
		authToken, _ := cmd.Flags().GetString("authtoken")
		verifyCertStr, _ := cmd.Flags().GetString("verifyCert")
		
		verifyCert, parseErr := strconv.ParseBool(verifyCertStr)
		if parseErr != nil {
			fmt.Println("Invalid value for verifyCert. Use 'true' or 'false'.")
			return
		}

		if url == "" || authToken == "" {
			fmt.Println("Both --url and --authtoken are required.")
			return
		}

		config := Config{URL: url, AuthToken: authToken, VerifyCert: verifyCert}

		err := saveConfig(config)
		if err != nil {
			fmt.Printf("error saving configuration: %v\n", err)
			return
		}
		fmt.Println("Securden server configured with the CLI successfully.")
	},
}

func init() {
	configCmd.Flags().String("url", "", "Server URL")
	configCmd.Flags().String("authtoken", "", "Authentication token")
	configCmd.Flags().String("verifyCert", "true", "Enable SSL certificate verification (optional)")
	configCmd.MarkFlagRequired("url")
	configCmd.MarkFlagRequired("authtoken")
	RootCmd.AddCommand(configCmd)
}
