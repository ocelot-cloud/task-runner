//go:build !windows

package tr

import (
	"os"
	"os/exec"
)

func buildCommand(dir, commandStr string) *exec.Cmd {
	shell := os.Getenv("SHELL")
	if shell == "" {
		shell = "bash"
	}
	cmd := exec.Command(shell, "-c", commandStr)
	cmd.Dir = dir
	return cmd
}
