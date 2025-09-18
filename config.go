package main

import (
	"fmt"
	"os"
	"path/filepath"

	"jeera/decorators"
	"jeera/decorators/foreground"
	"github.com/joho/godotenv"
)

// Config holds the JIRA configuration
type Config struct {
	BaseURL  string
	Username string
	APIToken string
	UsePAT   bool // indicates if we're using Personal Access Token (Bearer auth)
}

// LoadConfig loads configuration from .env file and environment variables
func LoadConfig() *Config {
	// Try to load .env file (ignore error if file doesn't exist)
	if err := loadEnvFile(); err != nil {
		fmt.Printf(foreground.LIGHT_BLUE + "[INFO] No .env file found, using environment variables only\n" + decorators.RESET_ALL)
	}

	config := &Config{
		BaseURL:  getEnvOrDefault("JIRA_BASE_URL", ""),
		Username: getEnvOrDefault("JIRA_USERNAME", ""),
		APIToken: getAPIToken(),
	}

	// Determine authentication method based on token format or explicit setting
	config.UsePAT = detectPATUsage()

	return config
}

// loadEnvFile attempts to load environment variables from .env file
func loadEnvFile() error {
	// Get current working directory
	cwd, err := os.Getwd()
	if err != nil {
		return err
	}

	// Try multiple possible locations for .env file
	envPaths := []string{
		filepath.Join(cwd, ".env"),
		filepath.Join(cwd, ".env.local"),
		".env",
		".env.local",
	}

	var lastErr error
	for _, path := range envPaths {
		if err := godotenv.Load(path); err == nil {
			fmt.Printf(foreground.LIGHT_BLUE + "[INFO] Loaded configuration from: %s\n" + decorators.RESET_ALL, path)
			return nil
		} else {
			lastErr = err
		}
	}

	return lastErr
}

// getEnvOrDefault returns the environment variable value or a default value
func getEnvOrDefault(key, defaultValue string) string {
	if value := os.Getenv(key); value != "" {
		return value
	}
	return defaultValue
}

// getAPIToken gets the API token from either JIRA_API_TOKEN or JIRA_PAT
func getAPIToken() string {
	// Try JIRA_PAT first (Personal Access Token), then JIRA_API_TOKEN for backward compatibility
	if token := os.Getenv("JIRA_PAT"); token != "" {
		return token
	}
	return os.Getenv("JIRA_API_TOKEN")
}

// detectPATUsage determines if we should use PAT (Bearer) authentication
func detectPATUsage() bool {
	// If JIRA_PAT is explicitly set, use PAT authentication
	if os.Getenv("JIRA_PAT") != "" {
		return true
	}
	
	// If JIRA_USE_PAT is explicitly set to true, use PAT authentication
	if getEnvOrDefault("JIRA_USE_PAT", "false") == "true" {
		return true
	}
	
	// Check if token looks like a PAT (typically longer and different format)
	token := getAPIToken()
	if len(token) > 50 { // PATs are typically much longer than API tokens
		return true
	}
	
	return false
}

// Validate checks if all required configuration values are present
func (c *Config) Validate() bool {
	return c.BaseURL != "" && c.Username != "" && c.APIToken != ""
}
