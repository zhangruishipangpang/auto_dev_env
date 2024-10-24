package cmd

import (
	"auto_dev_env/src/util"
	"bytes"
	"fmt"
	"log"
	"os/exec"
	"path/filepath"
)

const (
	WindowPsEnvSet = "windows_ps_env_set.ps1"
	WindowPsEnvGet = "windows_ps_env_get.ps1"
)

type WinCmd struct {
}

func NewWinCmd() WinCmd {
	win := WinCmd{}
	return win
}

func (w WinCmd) Cmd(name string, arg ...string) ([]byte, error) {
	panic("unSupport operation")
}

func (w WinCmd) SetEnv(key, value string) error {

	shellDirPath := util.FindShellDirPath()

	scriptPath := filepath.Join(shellDirPath, WindowPsEnvSet)
	params := []string{
		"-File", scriptPath,
		"-Key", key,
		"-Value", value,
	}

	// 创建一个新的 Cmd 对象
	cmd := exec.Command("powershell.exe", params...)

	// 捕获标准输出和标准错误
	var out bytes.Buffer
	var stderr bytes.Buffer
	cmd.Stdout = &out
	cmd.Stderr = &stderr

	// 执行命令
	err := cmd.Run()
	if err != nil {
		log.Fatalf("Error executing PowerShell script: %v\n%v", err, stderr.String())
		return err
	}

	// 打印输出
	fmt.Printf("Output: %s\n", out.String())
	return nil
}

func (w WinCmd) GetEnv(key string) string {

	shellDirPath := util.FindShellDirPath()

	scriptPath := filepath.Join(shellDirPath, WindowPsEnvGet)
	params := []string{
		"-File", scriptPath,
		"-Key", key,
	}

	// 创建一个新的 Cmd 对象
	cmd := exec.Command("powershell.exe", params...)

	// 执行命令
	outputBytes, err := cmd.CombinedOutput()
	if err != nil {
		panic(err)
	}

	return string(outputBytes)
}
