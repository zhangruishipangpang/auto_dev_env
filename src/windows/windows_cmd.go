package windows

import (
	"auto_dev_env/src/util"
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
	// 创建命令
	cmd := exec.Command(name, arg...)

	// 执行命令并捕获输出
	output, err := cmd.CombinedOutput()
	if err != nil {
		return nil, err
	}

	return output, nil
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

	// 执行命令
	err := cmd.Run()
	if err != nil {
		return err
	}

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

	if len(outputBytes) == 0 {
		return ""
	}
	return string(outputBytes)
}
