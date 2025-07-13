package taskrunner

import (
	"fmt"
	"github.com/ocelot-cloud/task-runner/platform"
	"io"
	"os"
	"os/exec"
	"os/signal"
	"syscall"
)

func (t *TaskRunner) StartDaemon(dir, commandStr string, envs ...string) {
	cmd := platform.BuildCommand(dir, commandStr)
	appendEnvsToCommand(cmd, envs)
	platform.SetProcessGroup(cmd)

	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr

	err := cmd.Start()

	if cmd.Process == nil {
		t.Log.Error("error - the process was not able to start properly.")
		t.exitWithError()
		return
	}

	t.Config.idsOfDaemonProcessesCreated = append(t.Config.idsOfDaemonProcessesCreated, cmd.Process.Pid)

	if err != nil {
		t.Log.Error("Command: '%s' -> failed with error: %v", commandStr, err)
		t.exitWithError()
		return
	}

	t.Log.Info("started daemon with ID '%v' using command '%s'", cmd.Process.Pid, commandStr)

	go func() {
		if err := cmd.Wait(); err != nil {
			t.Log.Error("command: '%s' -> reason of stopping: %v", commandStr, err)
		} else {
			t.Log.Info("command: '%s' -> stopped through casual termination", commandStr)
		}
	}()
}

func (t *TaskRunner) Cleanup() {
	if t.Config.ShowCleanupOutput {
		t.Log.Info("Cleanup method called.")
	}
	t.killDaemonProcessesCreateDuringThisRun()
	if t.Config.CleanupFunc != nil {
		t.Config.CleanupFunc()
	}
	t.ResetCursor()
}

func (t *TaskRunner) exitWithError() {
	if t.Config.CleanupOnFailure && t.Config.CleanupFunc != nil {
		t.Cleanup()
	} else {
		t.killDaemonProcessesCreateDuringThisRun()
	}
	t.ResetCursor()
	os.Exit(1)
}

func (t *TaskRunner) ResetCursor() {
	fmt.Print("\x1b[?25h") // Shows the terminal cursor again if it was hidden.
	fmt.Print("\x1b[0m")   // Resets all terminal text attributes (color, bold, underline) back to default.
}

func (t *TaskRunner) killDaemonProcessesCreateDuringThisRun() {
	if len(t.Config.idsOfDaemonProcessesCreated) == 0 {
		return
	}
	t.Log.Info("Killing daemon processes")
	for _, pid := range t.Config.idsOfDaemonProcessesCreated {
		t.Log.Info("  Killing process with ID '%v'", pid)
		if err := platform.KillProcessGroup(pid); err != nil {
			t.Log.Error("Failed to kill process with ID '%v' because of error: %v", pid, err)
		}
	}
	t.Config.idsOfDaemonProcessesCreated = nil
}

func appendEnvsToCommand(cmd *exec.Cmd, envs []string) {
	envsWithLogLevel := append(envs, DefaultEnvs...)
	cmd.Env = append(os.Environ(), envsWithLogLevel...)
}

func (t *TaskRunner) ExitWithError() {
	t.Cleanup()
	os.Exit(1)
}

// TODO replace OS references with this interface
type OperatingSystem interface {
	BuildCommand(dir, commandStr string) *exec.Cmd
	KillProcessGroup(pid int) error
	SetProcessGroup(cmd *exec.Cmd)
	GetOsOutputs() (stdout, stderr io.Writer)
	GetOsEnvs() []string
}

var DefaultEnvs []string

func (t *TaskRunner) EnableAbortForKeystrokeControlPlusC() {
	go func() {
		sigChan := make(chan os.Signal, 1)
		signal.Notify(sigChan, syscall.SIGINT, syscall.SIGTERM)
		sig := <-sigChan
		t.Log.Info("Received signal: %v. Initiating graceful shutdown...", sig)
		t.Cleanup()
		os.Exit(1)
	}()
}
