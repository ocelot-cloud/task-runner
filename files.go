package taskrunner

import (
	"io"
	"os"
	"path/filepath"
)

func (t *TaskRunner) Copy(srcDir, fileName, destFolder string) {
	src := filepath.Join(srcDir, fileName)
	dest := filepath.Join(destFolder, fileName)

	info, err := os.Stat(src)
	if err != nil {
		t.Log.Error("error stating %s: %v", src, err)
		t.ExitWithError()
	}

	if info.IsDir() {
		t.copyDir(src, dest)
	} else {
		t.copyFile(src, dest)
	}
}

func (t *TaskRunner) copyFile(src, dest string) {
	info, err := os.Stat(src)
	if err != nil {
		t.Log.Error("error stating %s: %v", src, err)
		t.ExitWithError()
	}

	input, err := os.Open(src)
	if err != nil {
		t.Log.Error("error opening %s: %v", src, err)
		t.ExitWithError()
	}
	defer input.Close()

	destDir := filepath.Dir(dest)
	err = os.MkdirAll(destDir, 0700)
	if err != nil {
		t.Log.Error("error creating directory %s: %v", destDir, err)
		t.ExitWithError()
	}

	output, err := os.OpenFile(dest, os.O_CREATE|os.O_WRONLY|os.O_TRUNC, info.Mode())
	if err != nil {
		t.Log.Error("error creating %s: %v", dest, err)
		t.ExitWithError()
	}
	defer output.Close()

	_, err = io.Copy(output, input)
	if err != nil {
		t.Log.Error("error copying from %s to %s: %v", src, dest, err)
		t.ExitWithError()
	}

	_ = os.Chmod(dest, info.Mode())
}

func (t *TaskRunner) copyDir(srcDir, destDir string) {
	srcInfo, err := os.Stat(srcDir)
	if err != nil {
		t.Log.Error("error stating %s: %v", srcDir, err)
		t.ExitWithError()
	}

	err = os.MkdirAll(destDir, srcInfo.Mode())
	if err != nil {
		t.Log.Error("error creating directory %s: %v", destDir, err)
		t.ExitWithError()
	}

	entries, err := os.ReadDir(srcDir)
	if err != nil {
		t.Log.Error("error reading directory %s: %v", srcDir, err)
		t.ExitWithError()
	}

	for _, entry := range entries {
		srcPath := filepath.Join(srcDir, entry.Name())
		destPath := filepath.Join(destDir, entry.Name())

		if entry.IsDir() {
			t.copyDir(srcPath, destPath)
		} else {
			t.copyFile(srcPath, destPath)
		}
	}

	_ = os.Chmod(destDir, srcInfo.Mode())
}

func (t *TaskRunner) Remove(paths ...string) {
	for _, path := range paths {
		if _, err := os.Stat(path); os.IsNotExist(err) {
			continue
		}
		err := os.RemoveAll(path)
		if err != nil {
			t.Log.Error("Error removing %s: %v", path, err)
			t.ExitWithError()
		}
	}
}

func (t *TaskRunner) MakeDir(dir string) {
	err := os.MkdirAll(dir, 0700)
	if err != nil {
		t.Log.Error("Error creating %s: %v", dir, err)
		t.ExitWithError()
	}
}

func (t *TaskRunner) Move(src, dest string) {
	err := os.Rename(src, dest)
	if err != nil {
		t.Log.Error("error renaming %s to %s: %v", src, dest, err)
		t.ExitWithError()
	}
}

func (t *TaskRunner) Rename(dir, currentFileName, newFileName string) {
	src := filepath.Join(dir, currentFileName)
	dest := filepath.Join(dir, newFileName)
	err := os.Rename(src, dest)
	if err != nil {
		t.Log.Error("error renaming %s to %s: %v", src, dest, err)
		t.ExitWithError()
	}
}
