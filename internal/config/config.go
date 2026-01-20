package config

import (
	"bufio"
	"os"
	"path/filepath"
	"strings"
)

// Config holds the application configuration
type Config struct {
	WorkNotesLocation string
	WorkplaceName     string   // Default workplace (for backward compatibility)
	Workplaces        []string // List of available workplaces
	OpenCodeServer    string
	AIProvider        string
	AIModel           string
}

// Load reads the configuration from ~/.config/worklog/config
func Load() (*Config, error) {
	// Load config from ~/.config/worklog/config
	configPath := getConfigPath()
	loadConfigFile(configPath)

	workplaceName := getEnv("WORKPLACE_NAME", "Work")
	workplacesStr := getEnv("WORKPLACES", "")

	// Parse workplaces list
	var workplaces []string
	if workplacesStr != "" {
		for _, wp := range strings.Split(workplacesStr, ",") {
			wp = strings.TrimSpace(wp)
			if wp != "" {
				workplaces = append(workplaces, wp)
			}
		}
	}

	// If no workplaces defined, use the default workplace name
	if len(workplaces) == 0 {
		workplaces = []string{workplaceName}
	}

	cfg := &Config{
		WorkNotesLocation: getEnv("WORK_NOTES_LOCATION", "~/Documents/obsidian-notes/Inbox/work"),
		WorkplaceName:     workplaceName,
		Workplaces:        workplaces,
		OpenCodeServer:    getEnv("OPENCODE_SERVER", "http://127.0.0.1:4096"),
		AIProvider:        getEnv("AI_PROVIDER", "github-copilot"),
		AIModel:           getEnv("AI_MODEL", "claude-sonnet-4"),
	}

	// Expand ~ in the path
	cfg.WorkNotesLocation = expandPath(cfg.WorkNotesLocation)

	return cfg, nil
}

// getConfigPath returns the path to the config file
func getConfigPath() string {
	home, err := os.UserHomeDir()
	if err != nil {
		return ""
	}
	return filepath.Join(home, ".config", "worklog", "config")
}

// loadConfigFile reads a key=value config file and sets environment variables
func loadConfigFile(path string) {
	file, err := os.Open(path)
	if err != nil {
		return // Config file doesn't exist, use defaults
	}
	defer file.Close()

	scanner := bufio.NewScanner(file)
	for scanner.Scan() {
		line := strings.TrimSpace(scanner.Text())

		// Skip empty lines and comments
		if line == "" || strings.HasPrefix(line, "#") {
			continue
		}

		// Parse key=value
		parts := strings.SplitN(line, "=", 2)
		if len(parts) != 2 {
			continue
		}

		key := strings.TrimSpace(parts[0])
		value := strings.TrimSpace(parts[1])

		// Only set if not already set in environment
		if _, exists := os.LookupEnv(key); !exists {
			os.Setenv(key, value)
		}
	}
}

// getEnv retrieves an environment variable or returns a default value
func getEnv(key, defaultValue string) string {
	if value, exists := os.LookupEnv(key); exists {
		return value
	}
	return defaultValue
}

// expandPath expands ~ to the user's home directory
func expandPath(path string) string {
	if strings.HasPrefix(path, "~/") {
		home, err := os.UserHomeDir()
		if err != nil {
			return path
		}
		return filepath.Join(home, path[2:])
	}
	return path
}

// EnsureNotesDirectory creates the notes directory if it doesn't exist
func (c *Config) EnsureNotesDirectory() error {
	return os.MkdirAll(c.WorkNotesLocation, 0755)
}
