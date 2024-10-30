package main

import (
	"auto_dev_env/src/env"
	_ "auto_dev_env/src/platform"
	"auto_dev_env/src/util"
	"flag"
	"github.com/fatih/color"
	"path/filepath"
)

var cp = color.New(color.FgCyan).Add(color.Bold)

func main() {

	_, _ = cp.Printf("\n\n ++++开始执行开发环境环境变量配置程序++++ ")

	// 定义命令行参数
	configPath := flag.String("config", "", "config path")
	osName := flag.String("os_name", "", "os name")

	// 解析命令行参数
	flag.Parse()

	if configPath == nil || *configPath == "" {
		*configPath = filepath.Join(util.FindCurrentDir(), "config", "config.json")
	}

	cp.Printf("\n 配置文件：%s \n", *configPath)

	processor := env.NewEnvProcessorByCurrentOsName(*osName, *configPath)
	processor.Process()
}

//
// the questions to ask
//var qs = []*survey.Question{
//	{
//		Name:      "name",
//		Prompt:    &survey.Input{Message: "What is your name?"},
//		Validate:  survey.Required,
//		Transform: survey.Title,
//	},
//	{
//		Name: "color",
//		Prompt: &survey.Select{
//			Message: "Choose a color:",
//			Options: []string{"red", "blue", "green"},
//			Default: "red",
//		},
//	},
//	{
//		Name:   "age",
//		Prompt: &survey.Input{Message: "How old are you?"},
//	},
//}
//
//func main() {
//	// the answers will be written to this struct
//	answers := struct {
//		Name          string // survey will match the question and field names
//		FavoriteColor string `survey:"color"` // or you can tag fields to match a specific name
//		Age           int    // if the types don't match, survey will convert it
//	}{}
//
//	// perform the questions
//	err := survey.Ask(qs, &answers)
//	if err != nil {
//		fmt.Println(err.Error())
//		return
//	}
//
//	fmt.Printf("%s chose %s.", answers.Name, answers.FavoriteColor)
//}
