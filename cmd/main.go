package main

import (
	"auto_dev_env/src/args_ask"
	"auto_dev_env/src/env"
	_ "auto_dev_env/src/platform"
	"auto_dev_env/src/util"
	"fmt"
	"github.com/fatih/color"
	"path/filepath"
)

var cp = color.New(color.FgCyan).Add(color.Bold)

func main() {

	answers, err := args_ask.Ask()
	if err != nil {
		panic(err)
	}

	_, _ = cp.Printf("\n\n ##################开始执行开发环境环境变量配置程序################## ")

	// 定义命令行参数
	configPath := answers.ConfigFilePath
	osName := answers.OsName

	if configPath == "" {
		configPath = filepath.Join(util.FindCurrentDir(), "config", "config.json")
	}

	fmt.Println()
	cp.Printf(" 配置文件：%s \n", configPath)

	processor := env.NewEnvProcessorByCurrentOsName(osName, configPath, answers.Envs)
	processor.Process()
}
