package tr

import (
	"bytes"
	"fmt"
	"github.com/ocelot-cloud/task-runner/platform"
	"io"
	"log"
	"os"
	"path/filepath"
	"strings"
	"time"
)

var parentDir = getParentDir()

func getParentDir() string {
	currentDir, err := os.Getwd()
	if err != nil {
		log.Fatalf("failed to get current dir: %v", err)
	}
	return filepath.Dir(currentDir)
}

func ExecuteInDir(dir string, commandStr string, envs ...string) {
	elapsedTimeStr, err := executeInDir(dir, commandStr, envs...)
	Log.Info(elapsedTimeStr)
	if err != nil {
		Log.Error(" => %v", err)
		ExitWithError()
	}
}

func executeInDir(dir string, commandStr string, envs ...string) (string, error) {
	shortDir := strings.Replace(dir, parentDir, "", -1)
	Log.Info("in directory '.%s', executing '%s'", shortDir, commandStr)

	cmd := platform.BuildCommand(dir, commandStr)
	appendEnvsToCommand(cmd, envs)

	var stdoutBuf, stderrBuf bytes.Buffer
	stdoutMulti := io.MultiWriter(os.Stdout, &stdoutBuf)
	stderrMulti := io.MultiWriter(os.Stderr, &stderrBuf)
	cmd.Stdout = stdoutMulti
	cmd.Stderr = stderrMulti

	startTime := time.Now()
	err := cmd.Run()
	elapsed := time.Since(startTime)
	elapsedStr := fmt.Sprintf("%.3f", elapsed.Seconds())

	elapsedTimeSummary := fmt.Sprintf("Time taken: %s seconds.", elapsedStr)
	if err != nil {
		return elapsedTimeSummary, fmt.Errorf("command failed with error: %v", err)
	} else {
		Log.Info(" => Command successful.")
		return elapsedTimeSummary, nil
	}
}

func Execute(commandStr string, envs ...string) {
	ExecuteInDir(".", commandStr, envs...)
}

func PromptForContinuation(prompt string) {
	fmt.Printf("%s (y/N): ", prompt)
	var response string
	fmt.Scanln(&response)
	if response != "y" && response != "Y" {
		fmt.Println("Command aborted.")
		os.Exit(0)
	}
}
