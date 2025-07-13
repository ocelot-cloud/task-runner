package tr

import (
	"bytes"
	"log"
	"os/exec"
	"strings"
	"testing"
	"time"
)

func TestCleanup(t *testing.T) {
	StartDaemon(".", "sleep 10")
	Cleanup()
	assertThatNoProcessesSurvived([]string{"sleep 10"})
	idsOfDaemonProcessesCreatedDuringThisRun = []int{}
}

func assertThatNoProcessesSurvived(processes []string) {
	cmd := exec.Command("ps", "aux")
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
