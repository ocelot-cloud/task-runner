//go:build integration

package tr

import (
	"github.com/stretchr/testify/assert"
	"os"
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

	srcDirInfo, _ := os.Stat(tmpDir)
	dstDirInfo, _ := os.Stat("temp2/" + tmpDir)
	assert.Equal(t, srcDirInfo.Mode(), dstDirInfo.Mode())

	srcFileInfo, _ := os.Stat(tmpDir + "/test.txt")
	dstFileInfo, _ := os.Stat("temp2/" + tmpDir + "/test.txt")
	assert.Equal(t, srcFileInfo.Mode(), dstFileInfo.Mode())
}

func TestDirMove(t *testing.T) {
	defer Remove(tmpDir)
	MakeDir(tmpDir)
	assert.True(t, checkIfExists(tmpDir))
	ExecuteInDir(tmpDir, "touch test.txt")
	assert.True(t, checkIfExists(tmpDir+"/test.txt"))

	srcInfo, _ := os.Stat(tmpDir)
	srcFileInfo, _ := os.Stat(tmpDir + "/test.txt")

	defer Remove(tmpDir2)
	Move(tmpDir, tmpDir2)
	assert.False(t, checkIfExists(tmpDir))
	assert.True(t, checkIfExists("temp2"))
	assert.True(t, checkIfExists("temp2/test.txt"))

	dstDirInfo, _ := os.Stat(tmpDir2)
	dstFileInfo, _ := os.Stat(tmpDir2 + "/test.txt")
	assert.Equal(t, srcInfo.Mode(), dstDirInfo.Mode())
	assert.Equal(t, srcFileInfo.Mode(), dstFileInfo.Mode())
}

func TestRename(t *testing.T) {
	defer Remove(tmpDir)
	MakeDir(tmpDir)
	ExecuteInDir(tmpDir, "touch test.txt")
	assert.True(t, checkIfExists(tmpDir+"/test.txt"))
	assert.False(t, checkIfExists(tmpDir+"/test2.txt"))

	origInfo, _ := os.Stat(tmpDir + "/test.txt")
	origMode := origInfo.Mode()

	Rename(tmpDir, "test.txt", "test2.txt")
	assert.False(t, checkIfExists(tmpDir+"/test.txt"))
	assert.True(t, checkIfExists(tmpDir+"/test2.txt"))

	srcFileInfo, _ := os.Stat(tmpDir + "/test2.txt")
	assert.Equal(t, origMode, srcFileInfo.Mode())
}
