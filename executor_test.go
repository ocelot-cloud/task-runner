package tr

import (
	"fmt"
	"github.com/stretchr/testify/assert"
	"os"
	"testing"
)

var (
	sampleTestDir = "sample_tests"
	tmpDir        = "temp"
	tmpDir2       = "temp2"
)

func TestMain(m *testing.M) {
	Remove(tmpDir, tmpDir2)
	exitCode := m.Run()
	os.Exit(exitCode)
}

func TestCommandSuccessful(t *testing.T) {
	_, err := executeInDir(sampleTestDir, "go test success_test.go")
	assert.Nil(t, err)
}

func TestDirCreationAndDeletion(t *testing.T) {
	assert.False(t, checkIfExists(tmpDir))
	defer Remove(tmpDir)
	MakeDir(tmpDir)
	assert.True(t, checkIfExists(tmpDir))

	ExecuteInDir(tmpDir, "touch test.txt")
	assert.True(t, checkIfExists(tmpDir+"/test.txt"))
	Remove(tmpDir)
	assert.False(t, checkIfExists(tmpDir))
}

func checkIfExists(path string) bool {
	_, err := os.Stat(path)
	if err == nil {
		return true
	}
	if os.IsNotExist(err) {
		return false
	}
	panic("unexpected error: " + err.Error())
}

func TestDaemon(t *testing.T) {
	assert.Equal(t, 0, len(idsOfDaemonProcessesCreatedDuringThisRun))

	StartDaemon(".", "sleep 100")
	assert.Equal(t, 1, len(idsOfDaemonProcessesCreatedDuringThisRun))
	processId := idsOfDaemonProcessesCreatedDuringThisRun[0]
	command := fmt.Sprintf("bash -c 'ps -p %d -o cmd= | grep -q sleep'", processId)
	ExecuteInDir(".", command)

	Cleanup()
	assert.Equal(t, 0, len(idsOfDaemonProcessesCreatedDuringThisRun))
	command = fmt.Sprintf("bash -c '! ps -p %d'", processId)
	ExecuteInDir(".", command)
}

func TestCustomCleanupFunction(t *testing.T) {
	defer Remove(tmpDir)
	CustomCleanupFunc = func() {
		MakeDir(tmpDir)
	}
	assert.False(t, checkIfExists(tmpDir))
	Cleanup()
	assert.True(t, checkIfExists(tmpDir))
}
