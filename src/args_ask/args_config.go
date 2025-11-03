package args_ask

import (
	"auto_dev_env/src/platform"
	"auto_dev_env/src/util"
	"bufio"
	"os"
	"path/filepath"
	"strings"

	"github.com/AlecAivazis/survey/v2"
)

type Args struct {
	OsName         string   `survey:"OsName"`         // 操作系统类型
	ConfigFilePath string   `survey:"configFilePath"` // 配置文件路径
	Envs           []string `survey:"envs"`           // 需要配置的环境变量列表
}

// the questions to ask
var qs = []*survey.Question{
	{
		Name: "OsName",
		Prompt: &survey.Select{
			Message: "选择 os ? (默认当前系统)",
			Options: platform.OsStore,
			Default: util.GetCurrentOs(),
		},
		Validate: survey.Required,
	},
	{
		Name:      "ConfigFilePath",
		Prompt:    &survey.Input{Message: "输入配置文件路径：（默认/config/config.json）"},
		Transform: survey.Title,
	},
	{
		Name: "envs",
		Prompt: &survey.MultiSelect{
			Message: "选择需要初始化的环境变量",
			Options: findEnvs(),
		},
	},
}

func Ask() (Args, error) {
	answers := Args{}

	// perform the questions
	err := survey.Ask(qs, &answers)
	if err != nil {
		return Args{}, err
	}

	util.Info(" 待配置信息：  ")
	util.Info("         os    : %s  ", answers.OsName)
	util.Info("         config: %s  ", answers.ConfigFilePath)
	util.Info("         envs  : %s  ", strings.Join(answers.Envs, ","))

	return answers, nil

}

func findEnvs() []string {
	choosesPath := "envs_chooses.txt"

	abPath := filepath.Join("./config", choosesPath)
	util.Debug("查找可选环境变量文件: %s", abPath)

	// 打开文件
	file, err := os.Open(abPath)
	if err != nil {
		util.Error("无法打开环境变量选择文件: %v", err)
		util.Warn("使用空环境变量列表")
		return []string{}
	}
	defer file.Close()

	// 创建一个新的缓冲区读取器
	scanner := bufio.NewScanner(file)

	// 按行读取文件内容
	var lines []string
	for scanner.Scan() {
		line := scanner.Text()
		if strings.TrimSpace(line) != "" {
			lines = append(lines, line)
		}
	}

	// 检查读取过程中是否发生错误
	if err := scanner.Err(); err != nil {
		util.Error("读取环境变量文件时出错: %v", err)
		util.Warn("使用已读取的环境变量列表")
	}

	util.Debug("成功读取 %d 个可选环境变量", len(lines))

	return lines
}
