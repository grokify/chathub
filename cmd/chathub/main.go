// Package main provides the ChatHub MCP server entry point.
package main

import (
	"context"
	"crypto/rand"
	"encoding/base64"
	"fmt"
	"log"
	"os"
	"os/signal"
	"syscall"

	"github.com/agentplexus/mcpkit/runtime"
	"github.com/grokify/chathub/internal/config"
	"github.com/grokify/chathub/internal/storage"
	"github.com/grokify/chathub/internal/tools"
	"github.com/modelcontextprotocol/go-sdk/mcp"
)

const (
	appName    = "chathub"
	appVersion = "0.1.1"
)

func main() {
	if err := run(); err != nil {
		log.Fatal(err)
	}
}

func run() error {
	// Load configuration
	cfg, err := config.Load()
	if err != nil {
		return fmt.Errorf("failed to load config: %w", err)
	}

	// Initialize storage backend
	store, err := storage.NewFromConfig(cfg.Backend, cfg.BackendConfig, cfg.Folder)
	if err != nil {
		return fmt.Errorf("failed to initialize storage: %w", err)
	}
	defer store.Close()

	// Create MCP runtime
	rt := runtime.New(&mcp.Implementation{
		Name:    appName,
		Version: appVersion,
	}, nil)

	// Register tools
	tools.RegisterAll(rt, store)

	// Set up context with cancellation
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	// Handle signals for graceful shutdown
	sigChan := make(chan os.Signal, 1)
	signal.Notify(sigChan, syscall.SIGINT, syscall.SIGTERM)
	go func() {
		<-sigChan
		cancel()
	}()

	// Serve based on transport
	switch cfg.Transport {
	case config.TransportStdio:
		log.Printf("Starting %s v%s (stdio transport)", appName, appVersion)
		return rt.ServeStdio(ctx)

	case config.TransportHTTP:
		opts := &runtime.HTTPServerOptions{
			Addr: fmt.Sprintf(":%d", cfg.Port),
			OnReady: func(result *runtime.HTTPServerResult) {
				if result.PublicURL != "" {
					log.Printf("MCP server available at: %s", result.PublicURL)
				} else {
					log.Printf("MCP server available at: %s", result.LocalURL)
				}

				// Log OAuth 2.1 info if enabled
				if result.OAuth2 != nil {
					log.Println("")
					log.Println("OAuth 2.1 with PKCE enabled (for ChatGPT.com):")
					log.Printf("  Client ID:     %s", result.OAuth2.ClientID)
					log.Printf("  Client Secret: %s", result.OAuth2.ClientSecret)
					log.Printf("  Authorization: %s", result.OAuth2.AuthorizationEndpoint)
					log.Printf("  Token:         %s", result.OAuth2.TokenEndpoint)
					log.Printf("  Login User:    %v", result.OAuth2.Users)
					log.Println("")
				}

				// Log legacy OAuth credentials if enabled
				if result.OAuth != nil {
					log.Println("")
					log.Println("OAuth 2.0 Client Credentials (legacy):")
					log.Printf("  Client ID:      %s", result.OAuth.ClientID)
					log.Printf("  Client Secret:  %s", result.OAuth.ClientSecret)
					log.Printf("  Token Endpoint: %s", result.OAuth.TokenEndpoint)
					log.Println("")
				}
			},
		}

		if cfg.NgrokEnabled {
			opts.Ngrok = &runtime.NgrokOptions{
				Authtoken: cfg.NgrokAuthtoken,
				Domain:    cfg.NgrokDomain,
			}
		}

		// Prefer OAuth 2.1 (for ChatGPT.com support) over legacy OAuth
		if cfg.OAuth2Enabled {
			password := cfg.OAuth2Password
			if password == "" {
				// Auto-generate password if not set
				password = generateRandomPassword()
				log.Printf("Generated OAuth2 login password: %s", password)
			}
			opts.OAuth2 = &runtime.OAuth2Options{
				Users: map[string]string{
					cfg.OAuth2Username: password,
				},
				ClientID:     cfg.OAuth2ClientID,
				ClientSecret: cfg.OAuth2ClientSecret,
				Debug:        true, // Enable OAuth debug logging
			}
		} else if cfg.OAuthEnabled {
			opts.OAuth = &runtime.OAuthOptions{
				ClientID:     cfg.OAuthClientID,
				ClientSecret: cfg.OAuthClientSecret,
			}
		}

		log.Printf("Starting %s v%s (HTTP transport)", appName, appVersion)

		_, err := rt.ServeHTTP(ctx, opts)
		if err != nil {
			return fmt.Errorf("failed to start HTTP server: %w", err)
		}

		return nil

	default:
		return fmt.Errorf("unknown transport: %s", cfg.Transport)
	}
}

// generateRandomPassword generates a secure random password.
func generateRandomPassword() string {
	bytes := make([]byte, 16)
	if _, err := rand.Read(bytes); err != nil {
		// Fallback to a default if crypto/rand fails (shouldn't happen)
		return "changeme"
	}
	return base64.RawURLEncoding.EncodeToString(bytes)
}
