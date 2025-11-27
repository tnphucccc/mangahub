package config

import (
	"fmt"
	"os"

	"gopkg.in/yaml.v3"
)

// Config holds all application configuration
type Config struct {
	Server   ServerConfig   `yaml:"server"`
	Database DatabaseConfig `yaml:"database"`
	JWT      JWTConfig      `yaml:"jwt"`
}

// ServerConfig holds server-specific configuration
type ServerConfig struct {
	Host          string `yaml:"host"`
	HTTPPort      string `yaml:"http_port"`
	TCPPort       string `yaml:"tcp_port"`
	UDPPort       string `yaml:"udp_port"`
	GRPCPort      string `yaml:"grpc_port"`
	WebSocketPort string `yaml:"websocket_port"`
}

// DatabaseConfig holds database configuration
type DatabaseConfig struct {
	Path string `yaml:"path"`
}

// JWTConfig holds JWT authentication configuration
type JWTConfig struct {
	Secret     string `yaml:"secret"`
	ExpiryDays int    `yaml:"expiry_days"`
}

// Load reads configuration from a YAML file
func Load(configPath string) (*Config, error) {
	// Read config file
	data, err := os.ReadFile(configPath)
	if err != nil {
		return nil, fmt.Errorf("failed to read config file: %w", err)
	}

	// Parse YAML
	var config Config
	if err := yaml.Unmarshal(data, &config); err != nil {
		return nil, fmt.Errorf("failed to parse config file: %w", err)
	}

	// Validate configuration
	if err := config.Validate(); err != nil {
		return nil, fmt.Errorf("invalid configuration: %w", err)
	}

	return &config, nil
}

// LoadFromEnv loads configuration with environment variable overrides
func LoadFromEnv(configPath string) (*Config, error) {
	config, err := Load(configPath)
	if err != nil {
		return nil, err
	}

	// Override with environment variables if present
	if host := os.Getenv("SERVER_HOST"); host != "" {
		config.Server.Host = host
	}
	if httpPort := os.Getenv("HTTP_PORT"); httpPort != "" {
		config.Server.HTTPPort = httpPort
	}
	if tcpPort := os.Getenv("TCP_PORT"); tcpPort != "" {
		config.Server.TCPPort = tcpPort
	}
	if udpPort := os.Getenv("UDP_PORT"); udpPort != "" {
		config.Server.UDPPort = udpPort
	}
	if grpcPort := os.Getenv("GRPC_PORT"); grpcPort != "" {
		config.Server.GRPCPort = grpcPort
	}
	if wsPort := os.Getenv("WEBSOCKET_PORT"); wsPort != "" {
		config.Server.WebSocketPort = wsPort
	}
	if dbPath := os.Getenv("DB_PATH"); dbPath != "" {
		config.Database.Path = dbPath
	}
	if jwtSecret := os.Getenv("JWT_SECRET"); jwtSecret != "" {
		config.JWT.Secret = jwtSecret
	}

	return config, nil
}

// Validate checks if the configuration is valid
func (c *Config) Validate() error {
	// Validate server config
	if c.Server.Host == "" {
		return fmt.Errorf("server host cannot be empty")
	}
	if c.Server.HTTPPort == "" {
		return fmt.Errorf("HTTP port cannot be empty")
	}
	if c.Server.TCPPort == "" {
		return fmt.Errorf("TCP port cannot be empty")
	}
	if c.Server.UDPPort == "" {
		return fmt.Errorf("UDP port cannot be empty")
	}
	if c.Server.GRPCPort == "" {
		return fmt.Errorf("gRPC port cannot be empty")
	}
	if c.Server.WebSocketPort == "" {
		return fmt.Errorf("WebSocket port cannot be empty")
	}

	// Validate database config
	if c.Database.Path == "" {
		return fmt.Errorf("database path cannot be empty")
	}

	// Validate JWT config
	if c.JWT.Secret == "" {
		return fmt.Errorf("JWT secret cannot be empty")
	}
	if c.JWT.ExpiryDays <= 0 {
		return fmt.Errorf("JWT expiry days must be positive")
	}

	return nil
}

// GetHTTPAddress returns the full HTTP server address
func (c *Config) GetHTTPAddress() string {
	return fmt.Sprintf("%s:%s", c.Server.Host, c.Server.HTTPPort)
}

// GetTCPAddress returns the full TCP server address
func (c *Config) GetTCPAddress() string {
	return fmt.Sprintf("%s:%s", c.Server.Host, c.Server.TCPPort)
}

// GetUDPAddress returns the full UDP server address
func (c *Config) GetUDPAddress() string {
	return fmt.Sprintf("%s:%s", c.Server.Host, c.Server.UDPPort)
}

// GetGRPCAddress returns the full gRPC server address
func (c *Config) GetGRPCAddress() string {
	return fmt.Sprintf("%s:%s", c.Server.Host, c.Server.GRPCPort)
}

// GetWebSocketAddress returns the full WebSocket server address
func (c *Config) GetWebSocketAddress() string {
	return fmt.Sprintf("%s:%s", c.Server.Host, c.Server.WebSocketPort)
}
