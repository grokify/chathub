// Package config provides configuration loading for ChatHub.
package config

import (
	"errors"
	"fmt"
	"os"
	"strconv"
)

// Backend types
const (
	BackendGitHub  = "github"
	BackendS3      = "s3"
	BackendDropbox = "dropbox"
	BackendFile    = "file"
	BackendMemory  = "memory"
)

// Transport types
const (
	TransportStdio = "stdio"
	TransportHTTP  = "http"
)

// Config holds the ChatHub configuration.
type Config struct {
	Backend           string
	BackendConfig     map[string]string
	Folder            string
	Transport         string
	Port              int
	NgrokEnabled      bool
	NgrokAuthtoken    string
	NgrokDomain       string
	OAuthEnabled      bool
	OAuthClientID     string
	OAuthClientSecret string

	// OAuth 2.1 with PKCE (for ChatGPT.com and other clients requiring DCR)
	OAuth2Enabled      bool
	OAuth2Username     string
	OAuth2Password     string
	OAuth2ClientID     string
	OAuth2ClientSecret string
}

// Load loads configuration from environment variables.
func Load() (*Config, error) {
	ngrokEnabled := getEnvBool("CHATHUB_NGROK", false)

	oauth2Enabled := getEnvBool("CHATHUB_OAUTH2", ngrokEnabled) // defaults to true when ngrok is enabled

	cfg := &Config{
		Backend:            getEnv("CHATHUB_BACKEND", BackendGitHub),
		Folder:             getEnv("CHATHUB_FOLDER", "conversations"),
		Transport:          getEnv("CHATHUB_TRANSPORT", TransportStdio),
		Port:               getEnvInt("CHATHUB_PORT", 8080),
		NgrokEnabled:       ngrokEnabled,
		NgrokAuthtoken:     getEnv("CHATHUB_NGROK_AUTHTOKEN", ""),
		NgrokDomain:        getEnv("CHATHUB_NGROK_DOMAIN", ""),
		OAuthEnabled:       getEnvBool("CHATHUB_OAUTH", false), // legacy OAuth client_credentials
		OAuthClientID:      getEnv("CHATHUB_OAUTH_CLIENT_ID", ""),
		OAuthClientSecret:  getEnv("CHATHUB_OAUTH_CLIENT_SECRET", ""),
		OAuth2Enabled:      oauth2Enabled,
		OAuth2Username:     getEnv("CHATHUB_OAUTH2_USERNAME", "admin"),
		OAuth2Password:     getEnv("CHATHUB_OAUTH2_PASSWORD", ""),
		OAuth2ClientID:     getEnv("CHATHUB_OAUTH2_CLIENT_ID", ""),
		OAuth2ClientSecret: getEnv("CHATHUB_OAUTH2_CLIENT_SECRET", ""),
	}

	cfg.BackendConfig = loadBackendConfig(cfg.Backend)

	if err := cfg.Validate(); err != nil {
		return nil, err
	}

	return cfg, nil
}

// Validate validates the configuration.
func (c *Config) Validate() error {
	switch c.Backend {
	case BackendGitHub, BackendS3, BackendDropbox, BackendFile, BackendMemory:
		// valid
	default:
		return fmt.Errorf("invalid backend: %s", c.Backend)
	}

	switch c.Transport {
	case TransportStdio, TransportHTTP:
		// valid
	default:
		return fmt.Errorf("invalid transport: %s", c.Transport)
	}

	// Validate backend-specific config
	switch c.Backend {
	case BackendGitHub:
		if c.BackendConfig["token"] == "" {
			return errors.New("GITHUB_TOKEN is required for GitHub backend")
		}
		if c.BackendConfig["owner"] == "" {
			return errors.New("GITHUB_OWNER is required for GitHub backend")
		}
		if c.BackendConfig["repo"] == "" {
			return errors.New("GITHUB_REPO is required for GitHub backend")
		}
	case BackendS3:
		if c.BackendConfig["bucket"] == "" {
			return errors.New("S3_BUCKET is required for S3 backend")
		}
	case BackendFile:
		if c.BackendConfig["root"] == "" {
			return errors.New("FILE_ROOT is required for file backend")
		}
	}

	return nil
}

func loadBackendConfig(backend string) map[string]string {
	switch backend {
	case BackendGitHub:
		return map[string]string{
			"token":  getEnv("GITHUB_TOKEN", ""),
			"owner":  getEnv("GITHUB_OWNER", ""),
			"repo":   getEnv("GITHUB_REPO", ""),
			"branch": getEnv("GITHUB_BRANCH", "main"),
		}
	case BackendS3:
		return map[string]string{
			"bucket":   getEnv("S3_BUCKET", ""),
			"region":   getEnv("S3_REGION", "us-east-1"),
			"endpoint": getEnv("S3_ENDPOINT", ""),
		}
	case BackendDropbox:
		return map[string]string{
			"token": getEnv("DROPBOX_TOKEN", ""),
		}
	case BackendFile:
		return map[string]string{
			"root": getEnv("FILE_ROOT", ""),
		}
	case BackendMemory:
		return map[string]string{}
	default:
		return map[string]string{}
	}
}

func getEnv(key, defaultValue string) string {
	if value := os.Getenv(key); value != "" {
		return value
	}
	return defaultValue
}

func getEnvInt(key string, defaultValue int) int {
	if value := os.Getenv(key); value != "" {
		if i, err := strconv.Atoi(value); err == nil {
			return i
		}
	}
	return defaultValue
}

func getEnvBool(key string, defaultValue bool) bool {
	if value := os.Getenv(key); value != "" {
		if b, err := strconv.ParseBool(value); err == nil {
			return b
		}
	}
	return defaultValue
}
