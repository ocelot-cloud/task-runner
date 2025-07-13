package tr

import (
	"fmt"
	"github.com/ocelot-cloud/task-runner/platform"
	"io"
	"os"
	"os/exec"
	"os/signal"
	"syscall"
)

// Do nothing. This function is meant to be overridden by the user if needed.
var CustomCleanupFunc func()

var idsOfDaemonProcessesCreatedDuringThisRun []int

func StartDaemon(dir, commandStr string, envs ...string) {
	cmd := platform.BuildCommand(dir, commandStr)
	appendEnvsToCommand(cmd, envs)
	platform.SetProcessGroup(cmd)

	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr

	err := cmd.Start()

	if cmd.Process == nil {
		Log.Error("error - the process was not able to start properly.")
		exitWithError()
		return
	}

	idsOfDaemonProcessesCreatedDuringThisRun = append(idsOfDaemonProcessesCreatedDuringThisRun, cmd.Process.Pid)

	if err != nil {
		Log.Error("Command: '%s' -> failed with error: %v", commandStr, err)
		exitWithError()
		return
	}

	Log.Info("started daemon with ID '%v' using command '%s'", cmd.Process.Pid, commandStr)

	go func() {
		if err := cmd.Wait(); err != nil {
			Log.Error("command: '%s' -> reason of stopping: %v", commandStr, err)
		} else {
			Log.Info("command: '%s' -> stopped through casual termination", commandStr)
		}
	}()
}

type config struct {
	CleanupOnFailure  bool
	ShowCleanupOutput bool
}

var cfg = config{
	CleanupOnFailure:  true,
	ShowCleanupOutput: true,
}

func HideCleanupOutput() { cfg.ShowCleanupOutput = false }

func Cleanup() {
	if cfg.ShowCleanupOutput {
		Log.Info("Cleanup method called.")
	}
	killDaemonProcessesCreateDuringThisRun()
	if CustomCleanupFunc != nil {
		CustomCleanupFunc()
	}
	ResetCursor()
}

func exitWithError() {
	if cfg.CleanupOnFailure && CustomCleanupFunc != nil {
		Cleanup()
	} else {
		killDaemonProcessesCreateDuringThisRun()
	}
	ResetCursor()
	os.Exit(1)
}

func ResetCursor() {
	fmt.Print("\x1b[?25h") // Shows the terminal cursor again if it was hidden.
	fmt.Print("\x1b[0m")   // Resets all terminal text attributes (color, bold, underline) back to default.
}

func killDaemonProcessesCreateDuringThisRun() {
	if len(idsOfDaemonProcessesCreatedDuringThisRun) == 0 {
		return
	}
	Log.Info("Killing daemon processes")
	for _, pid := range idsOfDaemonProcessesCreatedDuringThisRun {
		Log.Info("  Killing process with ID '%v'", pid)
		if err := platform.KillProcessGroup(pid); err != nil {
			Log.Error("Failed to kill process with ID '%v' because of error: %v", pid, err)
		}
	}
	idsOfDaemonProcessesCreatedDuringThisRun = nil
}

func appendEnvsToCommand(cmd *exec.Cmd, envs []string) {
	envsWithLogLevel := append(envs, DefaultEnvs...)
	cmd.Env = append(os.Environ(), envsWithLogLevel...)
}

func ExitWithError() {
	Cleanup()
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

func HandleSignals() {
	go func() {
		sigChan := make(chan os.Signal, 1)
		signal.Notify(sigChan, syscall.SIGINT, syscall.SIGTERM)
		sig := <-sigChan
		Log.Info("Received signal: %v. Initiating graceful shutdown...", sig)
		Cleanup()
		os.Exit(1)
	}()
}
