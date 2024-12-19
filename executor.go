package tr

import (
	"bytes"
	"fmt"
	"github.com/mattn/go-shellwords"
	"io"
	"log"
	"net"
	"net/http"
	"os"
	"os/exec"
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
	ColoredPrintln(elapsedTimeStr)
	if err != nil {
		ColoredPrintln(" => %v", err)
		CleanupAndExitWithError()
	}
}

func executeInDir(dir string, commandStr string, envs ...string) (string, error) {
	shortDir := strings.Replace(dir, parentDir, "", -1)
	ColoredPrintln("\nIn directory '.%s', executing '%s'\n", shortDir, commandStr)

	cmd := buildCommand(dir, commandStr)
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

	output := stdoutBuf.String() + stderrBuf.String()
	elapsedTimeSummary := fmt.Sprintf("Time taken: %s seconds.", elapsedStr)
	if err != nil {
		return elapsedTimeSummary, fmt.Errorf("command failed with error: %v", err)
	} else {
		if strings.Contains(commandStr, "go test") {
			return elapsedTimeSummary, checkGoTestOutput(output)
		} else {
			ColoredPrintln(" => Command successful.")
			return elapsedTimeSummary, nil
		}
	}
}

func checkGoTestOutput(output string) error {
	if strings.Contains(output, "no test files") {
		return fmt.Errorf("testing failed because no tests were found")
	} else if strings.Contains(output, "no tests to run") {
		return fmt.Errorf("testing failed because no tests were in test file")
	} else if !strings.Contains(output, "PASS:") && !containsOkLine(output) {
		return fmt.Errorf("testing failed because no tests were actually executed; all tests were either skipped or not included")
	} else {
		return nil
	}
}

func containsOkLine(output string) bool {
	lines := strings.Split(output, "\n")
	for _, line := range lines {
		if strings.HasPrefix(strings.TrimSpace(line), "ok") {
			return true
		}
	}
	return false
}

func ColoredPrintln(format string, a ...interface{}) {
	colorReset := "\033[0m"
	colorCode := "\033[31m"
	fmt.Printf(colorCode+format+"\n"+colorReset, a...)
}

func buildCommand(dir string, commandStr string) *exec.Cmd {
	parts, err := parseCommand(commandStr)
	if err != nil {
		ColoredPrintln("error parsing command: %v", err)
		CleanupAndExitWithError()
	}
	if len(parts) == 0 {
		ColoredPrintln("error, no command found in commandStr: %v", err)
		CleanupAndExitWithError()
	}
	command := parts[0]
	args := parts[1:]

	cmd := exec.Command(command, args...)
	cmd.Dir = dir

	return cmd
}

func parseCommand(command string) ([]string, error) {
	parser := shellwords.NewParser()
	return parser.Parse(command)
}

func WaitUntilPortIsReady(port string) {
	retryOperation(func() (bool, error) {
		conn, err := net.DialTimeout("tcp", "localhost:"+port, 1*time.Second)
		if err == nil {
			conn.Close()
			return true, nil
		}
		return false, err
	}, "port", "localhost:"+port, 30)
}

func retryOperation(operation func() (bool, error), description, target string, maxAttempts int) {
	attempt := 0
	for attempt < maxAttempts {
		success, err := operation()
		if success && err == nil {
			fmt.Printf("%s was requested successfully at %s\n", description, target)
			return
		} else {
			if attempt%5 == 0 {
				fmt.Printf("Attempt %v/%v: %s is not yet reachable at %s. error: %v. Trying again...\n", attempt, maxAttempts, description, target, err)
			}
			attempt++
			time.Sleep(1 * time.Second)
		}
	}
	fmt.Printf("Error: %s could not be reached in time at %s. Cleanup and exit...\n", description, target)
	CleanupAndExitWithError()
}

func WaitForWebPageToBeReady(url string) {
	retryOperation(func() (bool, error) {
		response, err := http.Get(url)
		if err == nil && response.StatusCode == 200 {
			return true, nil
		}
		return false, err
	}, "Index page", url, 30)
}

func PrintTaskDescription(text string) {
	ColoredPrintln("\n=== %s ===\n", text)
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
