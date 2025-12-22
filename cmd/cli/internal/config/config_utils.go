package config

import (
	"fmt"
	"os"
	"os/user"
	"path/filepath"

	"gopkg.in/yaml.v3"
)

// CLIConfig holds all application configuration for the CLI
type CLIConfig struct {
	Server        ServerConfig        `yaml:"server"`
	Database      DatabaseConfig      `yaml:"database"`
	User          UserConfig          `yaml:"user"`
	Sync          SyncConfig          `yaml:"sync"`
	Notifications NotificationsConfig `yaml:"notifications"`
	Logging       LoggingConfig       `yaml:"logging"`
}

// ServerConfig holds server-specific configuration for the CLI to connect to
type ServerConfig struct {
	Host          string `yaml:"host"`
	HTTPPort      int    `yaml:"http_port"`
	TCPPort       int    `yaml:"tcp_port"`
	UDPPort       int    `yaml:"udp_port"`
	GRPCPort      int    `yaml:"grpc_port"`
	WebSocketPort int    `yaml:"websocket_port"`
}

// DatabaseConfig holds database configuration for the CLI (if it manages its own local DB)
type DatabaseConfig struct {
	Path string `yaml:"path"`
}

// UserConfig holds user-specific configuration (e.g., current authenticated user token)
type UserConfig struct {
	Username string `yaml:"username"`
	Token    string `yaml:"token"`
}

// SyncConfig holds synchronization settings
type SyncConfig struct {
	AutoSync           bool   `yaml:"auto_sync"`
	ConflictResolution string `yaml:"conflict_resolution"`
}

// NotificationsConfig holds notification settings
type NotificationsConfig struct {
	Enabled bool `yaml:"enabled"`
	Sound   bool `yaml:"sound"`
}

// LoggingConfig holds logging settings
type LoggingConfig struct {
	Level string `yaml:"level"`
	Path  string `yaml:"path"`
}

func GetCLIConfigPath() (string, string, error) {
	usr, err := user.Current()
	if err != nil {
		return "", "", fmt.Errorf("error getting current user home directory: %w", err)
	}
	configDir := filepath.Join(usr.HomeDir, ".mangahub")
	configPath := filepath.Join(configDir, "config.yaml")
	return configDir, configPath, nil
}

func LoadCLIConfig() (*CLIConfig, error) {
	_, configPath, err := GetCLIConfigPath()
	if err != nil {
		return nil, err
	}

	yamlFile, err := os.ReadFile(configPath)
	if err != nil {
		return nil, fmt.Errorf("error reading config file: %w", err)
	}

	var config CLIConfig
	err = yaml.Unmarshal(yamlFile, &config)
	if err != nil {
		return nil, fmt.Errorf("error unmarshaling config file: %w", err)
	}

	return &config, nil
}

func SaveCLIConfig(config *CLIConfig) error {
	_, configPath, err := GetCLIConfigPath()
	if err != nil {
		return err
	}

	yamlData, err := yaml.Marshal(config)
	if err != nil {
		return fmt.Errorf("error marshaling config to YAML: %w", err)
	}

	err = os.WriteFile(configPath, yamlData, 0644)
	if err != nil {
		return fmt.Errorf("error writing config file: %w", err)
	}
	return nil
}

// InitConfig initializes the CLI configuration file with default values.
func InitConfig() {
	fmt.Println("Initializing MangaHub configuration...")

	configDir, configPath, err := GetCLIConfigPath()
	if err != nil {
		fmt.Printf("Error: %v\n", err)
		os.Exit(1)
	}

	// Create config directory if it doesn't exist
	if _, err := os.Stat(configDir); os.IsNotExist(err) {
		err = os.MkdirAll(configDir, 0755)
		if err != nil {
			fmt.Printf("Error creating config directory %s: %v\n", configDir, err)
			os.Exit(1)
		}
	}

	// Check if config file already exists
	if _, err := os.Stat(configPath); err == nil {
		fmt.Printf("Configuration file already exists at %s. To re-initialize, please remove it first.\n", configPath)
		os.Exit(0)
	}

	// Create default config
	defaultConfig := CLIConfig{
		Server: ServerConfig{
			Host:          "localhost",
			HTTPPort:      8080,
			TCPPort:       9090,
			UDPPort:       9091,
			GRPCPort:      9092,
			WebSocketPort: 9093,
		},
		Database: DatabaseConfig{
			Path: filepath.Join(configDir, "data.db"),
		},
		User: UserConfig{
			Username: "",
			Token:    "",
		},
		Sync: SyncConfig{
			AutoSync:           true,
			ConflictResolution: "last_write_wins",
		},
		Notifications: NotificationsConfig{
			Enabled: true,
			Sound:   false,
		},
		Logging: LoggingConfig{
			Level: "info",
			Path:  filepath.Join(configDir, "logs"),
		},
	}

	// Marshal to YAML
	yamlData, err := yaml.Marshal(&defaultConfig)
	if err != nil {
		fmt.Printf("Error marshaling default config to YAML: %v\n", err)
		os.Exit(1)
	}

	// Write to file
	err = os.WriteFile(configPath, yamlData, 0644)
	if err != nil {
		fmt.Printf("Error writing config file to %s: %v\n", configPath, err)
		os.Exit(1)
	}

	fmt.Printf("Configuration initialized successfully at %s\n", configPath)
}
