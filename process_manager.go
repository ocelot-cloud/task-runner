package tr

import (
	"fmt"
	platform2 "github.com/ocelot-cloud/task-runner/platform"
	"log"
	"os"
	"os/exec"
	"os/signal"
	"syscall"
)

var CustomCleanupFunc = func() {
	// Do nothing. This function is meant to be overridden by the user if needed.
}

var idsOfDaemonProcessesCreatedDuringThisRun []int

func StartDaemon(dir string, commandStr string, envs ...string) {
	var cmd *exec.Cmd
	cmd = platform2.BuildCommand(dir, commandStr)
	appendEnvsToCommand(cmd, envs)

	platform2.SetProcessGroup(cmd)

	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr

	err := cmd.Start()

	if cmd.Process == nil {
		fmt.Printf("Error: The process was not able to start properly.\n")
		CleanupAndExitWithError()
	} else {
		idsOfDaemonProcessesCreatedDuringThisRun = append(idsOfDaemonProcessesCreatedDuringThisRun, cmd.Process.Pid)
	}

	if err != nil {
		fmt.Printf("Command: '%s' -> failed with error: %v\n", commandStr, err)
		CleanupAndExitWithError()
	} else {
		ColoredPrintln("Started daemon with ID '%v' using command '%s'\n", cmd.Process.Pid, commandStr)
		go func() {
			err := cmd.Wait()
			if err != nil {
				fmt.Printf("Command: '%s' -> reason of stopping: %v\n", commandStr, err)
			} else {
				fmt.Printf("Command: '%s' -> stopped through casual termination\n", commandStr)
			}
		}()
	}
}

func Cleanup() {
	ColoredPrintln("Cleanup method called.\n")
	killDaemonProcessesCreateDuringThisRun()
	CustomCleanupFunc()
	fmt.Print("\x1b[?25h") // Ensure CLI cursor is visible
	fmt.Print("\x1b[0m")   // Resets all CLI cursor attributes such as color
}

func killDaemonProcessesCreateDuringThisRun() {
	println("Killing daemon processes")
	if len(idsOfDaemonProcessesCreatedDuringThisRun) == 0 {
		fmt.Println("  No daemon processes to kill.")
		return
	}

	for _, processID := range idsOfDaemonProcessesCreatedDuringThisRun {
		fmt.Printf("  Killing process with ID '%v'\n", processID)
		if err := platform2.KillProcessGroup(processID); err != nil {
			log.Fatalf("Failed to kill process with ID '%v' because of error: %v\n", processID, err)
		}
	}
	idsOfDaemonProcessesCreatedDuringThisRun = make([]int, 0)
}

func CleanupAndExitWithError() {
	Cleanup()
	os.Exit(1)
}

func appendEnvsToCommand(cmd *exec.Cmd, envs []string) {
	envsWithLogLevel := append(envs, DefaultEnvs...)
	cmd.Env = append(os.Environ(), envsWithLogLevel...)
}

var DefaultEnvs []string

func HandleSignals() {
	go func() {
		sigChan := make(chan os.Signal, 1)
		signal.Notify(sigChan, syscall.SIGINT, syscall.SIGTERM)
		sig := <-sigChan
		ColoredPrintln("\nReceived signal: %v. Initiating graceful shutdown...\n", sig)
		Cleanup()
		os.Exit(1)
	}()
}
