package commands

import (
	"crypto/tls"
	"crypto/x509"
	"encoding/json"
	"fmt"
	"net"
	"net/http"
	"net/url"
	"os"
	"path/filepath"
	"strconv"
	"strings"
	"time"

	"github.com/spf13/cobra"
)

func createConfigCmd() *cobra.Command {
	var cmd = &cobra.Command{
		Use:   "config",
		Short: "To save server connectivity details",
		Run: func(cmd *cobra.Command, args []string) {
			url, _ := cmd.Flags().GetString("url")
			verifyCertStr, _ := cmd.Flags().GetString("verifycert")
			
			verifyCert, err := strconv.ParseBool(verifyCertStr)
			if err != nil {
				fmt.Println("Invalid value for --verifycert. The value must be either 'true' or 'false'.")
				return
			}
			
			if url == "" {
				fmt.Println("--url field is required")
				return
			}
			
			config := Config{
				URL:        url,
				VerifyCert: verifyCert,
			}
			
			if err := saveConfig(config); err != nil {
				fmt.Printf("Error saving config: %v\n", err)
				return
			}
			
			fmt.Println("Securden server connectivity details saved successfully.")    
		},
	}
	// Configure flags
	cmd.Flags().String("url", "", "Server URL")
	cmd.Flags().String("verifycert", "true", "Enable SSL certificate verification (optional)")
	cmd.MarkFlagRequired("url")

	return cmd
}

type Config struct {
	URL        string `json:"url"`
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

func getHostFromURL(URL string) (string, error) { 
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

// func getHTTPClient(config Config) (*http.Client, error) {
// 	host, err := getHostFromURL(config.URL)
// 	if err != nil {
// 		return nil, fmt.Errorf("invalid BaseURL: %v", err)
// 	}
	
// 	if !config.VerifyCert {
// 		transport := &http.Transport{
// 			TLSClientConfig: &tls.Config{InsecureSkipVerify: true},
// 		}
// 		return &http.Client{Transport: transport}, nil
// 	}

// 	conn, err := tls.Dial("tcp", host, &tls.Config{InsecureSkipVerify: true})
// 	if err != nil {
// 		return nil, fmt.Errorf("failed to fetch certificate from server: %v", err)
// 	}
// 	defer conn.Close()

// 	if len(conn.ConnectionState().PeerCertificates) == 0 {
// 		return nil, fmt.Errorf("no certificates found on the server")
// 	}

// 	serverCert := conn.ConnectionState().PeerCertificates[0]
// 	certPool := x509.NewCertPool()
// 	certPool.AddCert(serverCert)

// 	transport := &http.Transport{
// 		TLSClientConfig: &tls.Config{RootCAs: certPool},
// 	}
// 	return &http.Client{Transport: transport}, nil
// }

func getHTTPClient(config Config) (*http.Client, error) {
	host, err := getHostFromURL(config.URL)
	if err != nil {
		return nil, fmt.Errorf("invalid BaseURL: %v", err)
	}

	// Configure timeouts for all cases
	timeoutConfig := &http.Client{
		Timeout: 30 * time.Second, // Total request timeout
		Transport: &http.Transport{
			DialContext: (&net.Dialer{
				Timeout:   10 * time.Second, // Connection timeout
			}).DialContext,
			TLSHandshakeTimeout:   10 * time.Second,
			ResponseHeaderTimeout: 10 * time.Second,
		},
	}

	if !config.VerifyCert {
		transport := &http.Transport{
			TLSClientConfig:     &tls.Config{InsecureSkipVerify: true},
			DialContext:        timeoutConfig.Transport.(*http.Transport).DialContext,
			TLSHandshakeTimeout: timeoutConfig.Transport.(*http.Transport).TLSHandshakeTimeout,
			ResponseHeaderTimeout: timeoutConfig.Transport.(*http.Transport).ResponseHeaderTimeout,
		}
		return &http.Client{
			Transport: transport,
			Timeout:   timeoutConfig.Timeout,
		}, nil
	}

	// For certificate verification dialer with timeout
	dialer := &net.Dialer{Timeout: 5 * time.Second}
	conn, err := tls.DialWithDialer(dialer, "tcp", host, &tls.Config{
		InsecureSkipVerify: true, // just for fetching the cert.
	})
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
		TLSClientConfig: &tls.Config{
			RootCAs: certPool,
		},
		DialContext:           timeoutConfig.Transport.(*http.Transport).DialContext,
		TLSHandshakeTimeout:   timeoutConfig.Transport.(*http.Transport).TLSHandshakeTimeout,
		ResponseHeaderTimeout: timeoutConfig.Transport.(*http.Transport).ResponseHeaderTimeout,
	}

	return &http.Client{
		Transport: transport,
		Timeout:   timeoutConfig.Timeout,
	}, nil
}