package tr

import (
	"testing"
)

func TestCleanup(t *testing.T) {
	StartDaemon(".", "sleep 10")
	KillProcesses([]string{"sleep"})
	assertThatNoProcessesSurvived([]string{"sleep 10"})
	idsOfDaemonProcessesCreatedDuringThisRun = []int{}
}
