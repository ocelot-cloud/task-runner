//go:build windows

package platform

import (
	"os"
	"os/exec"
)

func BuildCommand(dir, commandStr string) *exec.Cmd {
	shell := os.Getenv("COMSPEC")
	if shell == "" {
		shell = "powershell"
	}
	cmd := exec.Command(shell, "-Command", commandStr)
	cmd.Dir = dir
	return cmd
}

func KillProcesses(processes []string) {
	for _, p := range processes {
		cmd := exec.Command("cmd", "/C", fmt.Sprintf("taskkill /F /IM %s.exe /T", p))
		_ = cmd.Run()
	}
}
