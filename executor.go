package tr

import (
	"bytes"
	"fmt"
	"github.com/mattn/go-shellwords"
	"io"
	"log"
	"net"
	"net/http"
	"net/url"
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

	elapsedTimeSummary := fmt.Sprintf("Time taken: %s seconds.", elapsedStr)
	if err != nil {
		return elapsedTimeSummary, fmt.Errorf("command failed with error: %v", err)
	} else {
		ColoredPrintln(" => Command successful.")
		return elapsedTimeSummary, nil
	}
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

func WaitForWebPageToBeReady(targetUrl string) {
	retryOperation(func() (bool, error) {
		parsedURL, err := url.Parse(targetUrl)
		if err != nil {
			return false, err
		}

		req, err := http.NewRequest("GET", targetUrl, nil)
		if err != nil {
			return false, err
		}
		req.Header.Set("Origin", parsedURL.Scheme+"://"+parsedURL.Host)

		response, err := http.DefaultClient.Do(req)
		if err == nil && response.StatusCode == 200 {
			return true, nil
		}
		return false, err
	}, "Index page", targetUrl, 60)
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
