package tr

import (
	"bytes"
	"fmt"
	"log"
	"os/exec"
	"strings"
	"time"
)

func KillProcesses(processes []string) {
	processKillCommandTemplate := "pgrep -f %s | xargs -I %% kill -9 %%" // TODO get rid of linux dependency
	var processKillCommands []string
	for _, process := range processes {
		command := fmt.Sprintf(processKillCommandTemplate, process)
		processKillCommands = append(processKillCommands, command)
	}
	runShellCommands(processKillCommands)
	assertThatNoProcessesSurvived(processes)
}

func runShellCommands(commands []string) {
	for _, command := range commands {
		_ = exec.Command("/bin/sh", "-c", command).Run()
	}
}

func CleanupDockerArtifacts() {
	dockerPruningCommands := []string{
		"docker rm -f $(docker ps -a -q)",
		"docker network prune -f",
		"docker volume prune -a -f",
		"docker image prune -f"}
	runShellCommands(dockerPruningCommands)
}

func assertThatNoProcessesSurvived(processes []string) {
	cmd := exec.Command("ps", "aux") // TODO get rid of linux dependency
	var out bytes.Buffer
	cmd.Stdout = &out
	err := cmd.Run()
	if err != nil {
		log.Fatalf("Command finished with error: %v", err)
	}
	for _, process := range processes {
		for i := 0; i < 5; i++ {
			if !strings.Contains(out.String(), process) {
				break
			}
			if i == 4 {
				ColoredPrintln("The backend daemon process was not killed after tests were completed.\n")
				CleanupAndExitWithError()
			}
			time.Sleep(1 * time.Second)
		}
	}
}
