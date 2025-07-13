package tr

import (
	"io"
	"os"
	"path/filepath"
)

func Copy(srcDir, fileName, destFolder string) {
	src := filepath.Join(srcDir, fileName)
	dest := filepath.Join(destFolder, fileName)

	info, err := os.Stat(src)
	if err != nil {
		Log.Error("error stating %s: %v", src, err)
		ExitWithError()
	}

	if info.IsDir() {
		copyDir(src, dest)
	} else {
		copyFile(src, dest)
	}
}

func copyFile(src, dest string) {
	info, err := os.Stat(src)
	if err != nil {
		Log.Error("error stating %s: %v", src, err)
		ExitWithError()
	}

	input, err := os.Open(src)
	if err != nil {
		Log.Error("error opening %s: %v", src, err)
		ExitWithError()
	}
	defer input.Close()

	destDir := filepath.Dir(dest)
	err = os.MkdirAll(destDir, 0700)
	if err != nil {
		Log.Error("error creating directory %s: %v", destDir, err)
		ExitWithError()
	}

	output, err := os.OpenFile(dest, os.O_CREATE|os.O_WRONLY|os.O_TRUNC, info.Mode())
	if err != nil {
		Log.Error("error creating %s: %v", dest, err)
		ExitWithError()
	}
	defer output.Close()

	_, err = io.Copy(output, input)
	if err != nil {
		Log.Error("error copying from %s to %s: %v", src, dest, err)
		ExitWithError()
	}

	_ = os.Chmod(dest, info.Mode())
}

func copyDir(srcDir, destDir string) {
	srcInfo, err := os.Stat(srcDir)
	if err != nil {
		Log.Error("error stating %s: %v", srcDir, err)
		ExitWithError()
	}

	err = os.MkdirAll(destDir, srcInfo.Mode())
	if err != nil {
		Log.Error("error creating directory %s: %v", destDir, err)
		ExitWithError()
	}

	entries, err := os.ReadDir(srcDir)
	if err != nil {
		Log.Error("error reading directory %s: %v", srcDir, err)
		ExitWithError()
	}

	for _, entry := range entries {
		srcPath := filepath.Join(srcDir, entry.Name())
		destPath := filepath.Join(destDir, entry.Name())

		if entry.IsDir() {
			copyDir(srcPath, destPath)
		} else {
			copyFile(srcPath, destPath)
		}
	}

	_ = os.Chmod(destDir, srcInfo.Mode())
}

func Remove(paths ...string) {
	for _, path := range paths {
		if _, err := os.Stat(path); os.IsNotExist(err) {
			continue
		}
		err := os.RemoveAll(path)
		if err != nil {
			Log.Error("Error removing %s: %v", path, err)
			ExitWithError()
		}
	}
}

func MakeDir(dir string) {
	err := os.MkdirAll(dir, 0700)
	if err != nil {
		Log.Error("Error creating %s: %v", dir, err)
		ExitWithError()
	}
}

func Move(src, dest string) {
	err := os.Rename(src, dest)
	if err != nil {
		Log.Error("error renaming %s to %s: %v", src, dest, err)
		ExitWithError()
	}
}

func Rename(dir, currentFileName, newFileName string) {
	src := filepath.Join(dir, currentFileName)
	dest := filepath.Join(dir, newFileName)
	err := os.Rename(src, dest)
	if err != nil {
		Log.Error("error renaming %s to %s: %v", src, dest, err)
		ExitWithError()
	}
}
