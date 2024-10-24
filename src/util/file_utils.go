package util

import (
	"log"
	"os"
	"path/filepath"
)

const (
	AbsoluteShellDir = "shell"
)

func FindShellDirPath() string {

	workDir := FindCurrentDir()

	return filepath.Join(workDir, AbsoluteShellDir)
}

func FindCurrentDir() string {
	currentDir, err := os.Getwd()
	if err != nil {
		log.Fatalf("Error getting the current working directory: %v", err)
		panic("Error getting the current working directory")
	}

	// 打印当前工作目录
	return currentDir
}
