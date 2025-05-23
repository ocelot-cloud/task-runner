//go:build darwin
// +build darwin

package tr

import (
	"fmt"
	"os/exec"
	"syscall"
)

func setProcessGroup(cmd *exec.Cmd) {
	cmd.SysProcAttr = &syscall.SysProcAttr{Setpgid: true}
}

func killProcessGroup(processID int) error {
	pgid, err := syscall.Getpgid(processID)
	if err != nil {
		return fmt.Errorf("Failed to get process group ID of process ID '%v' because of error: %v", processID, err)
	}
	if err := syscall.Kill(-pgid, syscall.SIGKILL); err != nil {
		return fmt.Errorf("Failed to kill process group ID '%v' because of error: %v", pgid, err)
	}
	return nil
}
