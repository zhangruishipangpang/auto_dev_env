package mac

import (
	"auto_dev_env/src/util"
	"os/exec"
	"path/filepath"
)

const (
	PsEnvSet = "macos_bash_env_set.sh"
	PsEnvGet = "macos_bash_env_get.sh"
)

type MacCmd struct {
}

func NewMacCmd() MacCmd {
	mac := MacCmd{}
	return mac
}

func (w MacCmd) Cmd(name string, arg ...string) ([]byte, error) {
	// 创建命令
	cmd := exec.Command(name, arg...)

	// 执行命令并捕获输出
	output, err := cmd.CombinedOutput()
	if err != nil {
		return nil, err
	}

	return output, nil
}

func (w MacCmd) SetEnv(key, value string) error {

	shellDirPath := util.FindShellDirPath()

	scriptPath := filepath.Join(shellDirPath, PsEnvSet)

	// 创建一个新的 Cmd 对象
	cmd := exec.Command(scriptPath, key, value)

	// 执行命令
	err := cmd.Run()
	if err != nil {
		panic(err)
	}

	return nil
}

func (w MacCmd) GetEnv(key string) string {

	shellDirPath := util.FindShellDirPath()

	// 创建一个新的 Cmd 对象
	cmd := exec.Command(shellDirPath, key)

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
