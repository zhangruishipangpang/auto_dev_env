package util

import (
	"github.com/fatih/color"
	"os"
	"path/filepath"
)

var cpr = color.New(color.BgRed).Add(color.Bold)

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
		_, _ = cpr.Printf("Error getting the current working directory: %v", err)
		panic("Error getting the current working directory")
	}

	// 打印当前工作目录
	return currentDir
}
