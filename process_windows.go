//go:build windows
// +build windows

package tr

import (
	"fmt"
	"os"
	"os/exec"
	"syscall"
)

func setProcessGroup(cmd *exec.Cmd) {
	cmd.SysProcAttr = &syscall.SysProcAttr{
		CreationFlags: syscall.CREATE_NEW_PROCESS_GROUP,
	}
}

func killProcessGroup(processID int) error {
	process, err := os.FindProcess(processID)
	if err != nil {
		return fmt.Errorf("Failed to find process with ID '%v' because of error: %v", processID, err)
	}
	if err := process.Kill(); err != nil {
		return fmt.Errorf("Failed to kill process with ID '%v' because of error: %v", processID, err)
	}
	return nil
}
