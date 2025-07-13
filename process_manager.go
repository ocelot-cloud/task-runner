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

// TODO name should be casual "cleanup" function?
var CustomCleanupFunc = func() {
	// Do nothing. This function is meant to be overridden by the user if needed.
}

var idsOfDaemonProcessesCreatedDuringThisRun []int

func StartDaemon(dir, commandStr string, envs ...string) {
	cmd := platform.BuildCommand(dir, commandStr)
	appendEnvsToCommand(cmd, envs)
	platform.SetProcessGroup(cmd)

	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr

	err := cmd.Start()

	if cmd.Process == nil {
		Log.Error("Error: The process was not able to start properly.\n")
		exitWithError()
		return
	}

	idsOfDaemonProcessesCreatedDuringThisRun = append(idsOfDaemonProcessesCreatedDuringThisRun, cmd.Process.Pid)

	if err != nil {
		Log.Error("Command: '%s' -> failed with error: %v\n", commandStr, err)
		exitWithError()
		return
	}

	Log.Info("Started daemon with ID '%v' using command '%s'\n", cmd.Process.Pid, commandStr)

	go func() {
		if err := cmd.Wait(); err != nil {
			Log.Error("Command: '%s' -> reason of stopping: %v\n", commandStr, err)
		} else {
			Log.Info("Command: '%s' -> stopped through casual termination\n", commandStr)
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
		Log.Info("Cleanup method called.\n")
	}
	killDaemonProcessesCreateDuringThisRun()
	if CustomCleanupFunc != nil {
		CustomCleanupFunc()
	}

	// TODO if that is needed, explain why -> I thikn this should always be executed at the end of the program, but how to ensure that? -> maybe tell the user to call "defer ResetCursor()" in their main function?
	fmt.Print("\x1b[?25h") // Shows the terminal cursor again if it was hidden.
	fmt.Print("\x1b[0m")   // Resets all terminal text attributes (color, bold, underline) back to default.
}

func exitWithError() {
	if cfg.CleanupOnFailure && CustomCleanupFunc != nil {
		Cleanup()
	} else {
		killDaemonProcessesCreateDuringThisRun()
		// TODO explain why this is needed? Maybe abstract duplication
		fmt.Print("\x1b[?25h")
		fmt.Print("\x1b[0m")
	}
	os.Exit(1)
}

func killDaemonProcessesCreateDuringThisRun() {
	if len(idsOfDaemonProcessesCreatedDuringThisRun) == 0 {
		return
	}
	Log.Info("Killing daemon processes\n")
	for _, pid := range idsOfDaemonProcessesCreatedDuringThisRun {
		Log.Info("  Killing process with ID '%v'\n", pid)
		if err := platform.KillProcessGroup(pid); err != nil {
			Log.Error("Failed to kill process with ID '%v' because of error: %v\n", pid, err)
		}
	}
	idsOfDaemonProcessesCreatedDuringThisRun = nil
}

func appendEnvsToCommand(cmd *exec.Cmd, envs []string) {
	envsWithLogLevel := append(envs, DefaultEnvs...)
	cmd.Env = append(os.Environ(), envsWithLogLevel...)
}

// TODO maybe just "ExitWithError"? and cleanup is implicitly called?
func CleanupAndExitWithError() {
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
		Log.Info("\nReceived signal: %v. Initiating graceful shutdown...\n", sig)
		Cleanup()
		os.Exit(1)
	}()
}
