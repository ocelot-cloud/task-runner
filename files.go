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
		ColoredPrintln("error stating %s: %v", src, err)
		CleanupAndExitWithError()
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
		ColoredPrintln("error stating %s: %v", src, err)
		CleanupAndExitWithError()
	}

	input, err := os.Open(src)
	if err != nil {
		ColoredPrintln("error opening %s: %v", src, err)
		CleanupAndExitWithError()
	}
	defer input.Close()

	destDir := filepath.Dir(dest)
	err = os.MkdirAll(destDir, 0700)
	if err != nil {
		ColoredPrintln("error creating directory %s: %v", destDir, err)
		CleanupAndExitWithError()
	}

	output, err := os.OpenFile(dest, os.O_CREATE|os.O_WRONLY|os.O_TRUNC, info.Mode())
	if err != nil {
		ColoredPrintln("error creating %s: %v", dest, err)
		CleanupAndExitWithError()
	}
	defer output.Close()

	_, err = io.Copy(output, input)
	if err != nil {
		ColoredPrintln("error copying from %s to %s: %v", src, dest, err)
		CleanupAndExitWithError()
	}

	_ = os.Chmod(dest, info.Mode())
}

func copyDir(srcDir, destDir string) {
	srcInfo, err := os.Stat(srcDir)
	if err != nil {
		ColoredPrintln("error stating %s: %v", srcDir, err)
		CleanupAndExitWithError()
	}

	err = os.MkdirAll(destDir, srcInfo.Mode())
	if err != nil {
		ColoredPrintln("error creating directory %s: %v", destDir, err)
		CleanupAndExitWithError()
	}

	entries, err := os.ReadDir(srcDir)
	if err != nil {
		ColoredPrintln("error reading directory %s: %v", srcDir, err)
		CleanupAndExitWithError()
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
			ColoredPrintln("Error removing %s: %v", path, err)
			CleanupAndExitWithError()
		}
	}
}

func MakeDir(dir string) {
	err := os.MkdirAll(dir, 0700)
	if err != nil {
		ColoredPrintln("Error creating %s: %v", dir, err)
		CleanupAndExitWithError()
	}
}

func Move(src, dest string) {
	err := os.Rename(src, dest)
	if err != nil {
		ColoredPrintln("error renaming %s to %s: %v", src, dest, err)
		CleanupAndExitWithError()
	}
}

func Rename(dir, currentFileName, newFileName string) {
	src := filepath.Join(dir, currentFileName)
	dest := filepath.Join(dir, newFileName)
	err := os.Rename(src, dest)
	if err != nil {
		ColoredPrintln("error renaming %s to %s: %v", src, dest, err)
		CleanupAndExitWithError()
	}
}
