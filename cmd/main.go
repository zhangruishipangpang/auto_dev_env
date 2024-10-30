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
