package main

import (
	"fmt"
	"log"
	"os"
	"os/user"
	"path/filepath"
	"time"

	"github.com/golang-jwt/jwt/v4"
	"gopkg.in/yaml.v3"
)

type CLIConfig struct {
	Server struct {
		Host          string `yaml:"host"`
		HTTPPort      int    `yaml:"http_port"`
		TCPPort       int    `yaml:"tcp_port"`
		UDPPort       int    `yaml:"udp_port"`
		GRPCPort      int    `yaml:"grpc_port"`
		WebSocketPort int    `yaml:"websocket_port"`
	} `yaml:"server"`
	User struct {
		Username string `yaml:"username"`
		Token    string `yaml:"token"`
	} `yaml:"user"`
}

func main() {
	// 1. Generate a JWT token for 'user-alice'
	secret := "dev-secret-key-change-in-production"
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, jwt.MapClaims{
		"user_id":  "user-alice",
		"username": "alice",
		"exp":      time.Now().Add(time.Hour * 24).Unix(),
	})

	tokenString, err := token.SignedString([]byte(secret))
	if err != nil {
		log.Fatalf("Failed to sign token: %v", err)
	}

	// 2. Load/Create CLI config
	usr, _ := user.Current()
	configPath := filepath.Join(usr.HomeDir, ".mangahub", "config.yaml")
	
	configDir := filepath.Dir(configPath)
	os.MkdirAll(configDir, 0755)

	cfg := CLIConfig{}
	cfg.Server.Host = "localhost"
	cfg.Server.HTTPPort = 8080
	cfg.Server.TCPPort = 9090
	cfg.Server.UDPPort = 9091
	cfg.Server.GRPCPort = 9092
	cfg.Server.WebSocketPort = 9093
	cfg.User.Username = "alice"
	cfg.User.Token = tokenString

	data, _ := yaml.Marshal(&cfg)
	os.WriteFile(configPath, data, 0644)

	fmt.Printf("Successfully configured CLI for user 'alice' with token: %s...\n", tokenString[:10])
	fmt.Println("You can now run: go run cmd/cli/main.go stats view")
}
