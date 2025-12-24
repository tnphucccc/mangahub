package server

import (
	"fmt"
	"os"
	"os/exec"
	"strconv"

	"github.com/tnphucccc/mangahub/cmd/cli/internal/config"
)

func HandleServerCommand() {
	if len(os.Args) < 3 {
		printServerUsage()
		os.Exit(1)
	}

	subcommand := os.Args[2]

	switch subcommand {
	case "start":
		serverStart()
	case "stop":
		serverStop()
	case "status":
		serverStatus()
	default:
		fmt.Printf("Unknown server subcommand: %s\n", subcommand)
		printServerUsage()
		os.Exit(1)
	}
}

func printServerUsage() {
	fmt.Println("Usage: mangahub server <subcommand> [options]")
	fmt.Println("\nSubcommands:")
	fmt.Println("  start [api|tcp|udp|grpc|all]  Start a server (default: all)")
	fmt.Println("  stop [api|tcp|udp|grpc|all]   Stop a server (default: all)")
	fmt.Println("  status                        Show the status of all servers")
}

func serverStart() {
	server := "all"
	if len(os.Args) > 3 {
		server = os.Args[3]
	}

	switch server {
	case "api":
		startServer("api-server")
	case "tcp":
		startServer("tcp-server")
	case "udp":
		startServer("udp-server")
	case "grpc":
		startServer("grpc-server")
	case "all":
		startServer("api-server")
		startServer("tcp-server")
		startServer("udp-server")
		startServer("grpc-server")
	default:
		fmt.Printf("Unknown server: %s\n", server)
		printServerUsage()
		os.Exit(1)
	}
}

func serverStop() {
	server := "all"
	if len(os.Args) > 3 {
		server = os.Args[3]
	}

	switch server {
	case "api":
		stopServer("api-server")
	case "tcp":
		stopServer("tcp-server")
	case "udp":
		stopServer("udp-server")
	case "grpc":
		stopServer("grpc-server")
	case "all":
		stopServer("api-server")
		stopServer("tcp-server")
		stopServer("udp-server")
		stopServer("grpc-server")
	default:
		fmt.Printf("Unknown server: %s\n", server)
		printServerUsage()
		os.Exit(1)
	}
}

func serverStatus() {
	checkServerStatus("api-server")
	checkServerStatus("tcp-server")
	checkServerStatus("udp-server")
	checkServerStatus("grpc-server")
}

func startServer(name string) {
	pidFile := getPidFile(name)
	if _, err := os.Stat(pidFile); err == nil {
		fmt.Printf("Server '%s' is already running.\n", name)
		return
	}

	// Build the server binary
	binaryPath := fmt.Sprintf("./bin/%s", name)
	buildCmd := exec.Command("go", "build", "-o", binaryPath, fmt.Sprintf("./cmd/%s", name))
	if err := buildCmd.Run(); err != nil {
		fmt.Printf("Failed to build server '%s': %v\n", name, err)
		return
	}

	// Start the compiled binary
	cmd := exec.Command(binaryPath)
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr
	setProcessGroup(cmd)

	if err := cmd.Start(); err != nil {
		fmt.Printf("Failed to start server '%s': %v\n", name, err)
		return
	}

	if err := os.WriteFile(pidFile, []byte(strconv.Itoa(cmd.Process.Pid)), 0644); err != nil {
		fmt.Printf("Failed to write pid file for server '%s': %v\n", name, err)
		return
	}

	fmt.Printf("Server '%s' started with PID %d.\n", name, cmd.Process.Pid)
}

func stopServer(name string) {
	pidFile := getPidFile(name)
	pidBytes, err := os.ReadFile(pidFile)
	if err != nil {
		fmt.Printf("Server '%s' is not running.\n", name)
		return
	}

	pid, err := strconv.Atoi(string(pidBytes))
	if err != nil {
		fmt.Printf("Invalid PID in pid file for server '%s': %v\n", name, err)
		return
	}

	// Attempt to kill the process group.
	if err := killProcessGroup(pid); err != nil {
		fmt.Printf("Warning: could not kill process group %d (it may have already exited): %v\n", pid, err)
	}

	// Always attempt to remove the pid file.
	if err := os.Remove(pidFile); err != nil {
		fmt.Printf("Warning: failed to remove pid file for server '%s': %v\n", name, err)
	}

	fmt.Printf("Server '%s' stopped.\n", name)
}

func checkServerStatus(name string) {
	pidFile := getPidFile(name)
	_, err := os.Stat(pidFile)
	if err == nil {
		fmt.Printf("Server '%s' is running.\n", name)
	} else {
		fmt.Printf("Server '%s' is not running.\n", name)
	}
}

func getPidFile(name string) string {
	configDir, _, _ := config.GetCLIConfigPath()
	return fmt.Sprintf("%s/%s.pid", configDir, name)
}
