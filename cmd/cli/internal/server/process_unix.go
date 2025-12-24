//go:build !windows

package server

import (
	"os/exec"
	"syscall"
)

func setProcessGroup(cmd *exec.Cmd) {
	cmd.SysProcAttr = &syscall.SysProcAttr{Setpgid: true}
}

func killProcessGroup(pid int) error {
	// Negative PID kills the process group
	return syscall.Kill(-pid, syscall.SIGKILL)
}
