package args_ask

import (
	"auto_dev_env/src/platform"
	"github.com/AlecAivazis/survey/v2"
	"github.com/fatih/color"
	"strings"
)

var cpy = color.New(color.FgYellow).Add(color.Bold)

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
			Message: "选择 os ?",
			Options: platform.OsStore,
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
			Options: platform.EnvStore,
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

	cpy.Printf(" 待配置信息：  ")
	cpy.Printf("\n         os    : %s  ", answers.OsName)
	cpy.Printf("\n         config: %s  ", answers.ConfigFilePath)
	cpy.Printf("\n         envs  : %s  ", strings.Join(answers.Envs, ","))

	return answers, nil

}
