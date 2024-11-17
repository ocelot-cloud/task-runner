package tr

import (
	"github.com/stretchr/testify/assert"
	"testing"
)

func TestDirCopy(t *testing.T) {
	defer Remove(tmpDir)
	MakeDir(tmpDir)
	assert.True(t, checkIfExists(tmpDir))
	ExecuteInDir(tmpDir, "touch test.txt")
	assert.True(t, checkIfExists(tmpDir+"/test.txt"))

	defer Remove(tmpDir2)
	MakeDir(tmpDir2)
	Copy(".", "temp", tmpDir2)
	assert.True(t, checkIfExists(tmpDir))
	assert.True(t, checkIfExists("temp2/"+tmpDir))
	assert.True(t, checkIfExists("temp2/"+tmpDir+"/test.txt"))
}

func TestDirMove(t *testing.T) {
	defer Remove(tmpDir)
	MakeDir(tmpDir)
	assert.True(t, checkIfExists(tmpDir))
	ExecuteInDir(tmpDir, "touch test.txt")
	assert.True(t, checkIfExists(tmpDir+"/test.txt"))

	defer Remove(tmpDir2)
	Move(tmpDir, tmpDir2)
	assert.False(t, checkIfExists(tmpDir))
	assert.True(t, checkIfExists("temp2"))
	assert.True(t, checkIfExists("temp2/test.txt"))
}
