//go:build windows

package tr

import (
	"os"
	"os/exec"
)

func buildCommand(dir, commandStr string) *exec.Cmd {
	shell := os.Getenv("COMSPEC")
	if shell == "" {
		shell = "powershell"
	}
	cmd := exec.Command(shell, "-Command", commandStr)
	cmd.Dir = dir
	return cmd
}
