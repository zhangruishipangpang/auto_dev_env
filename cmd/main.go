package main

import (
	"auto_dev_env/src/args_ask"
	"auto_dev_env/src/env"
	_ "auto_dev_env/src/platform"
	"auto_dev_env/src/util"
	"path/filepath"
)

func main() {
	// 初始化日志系统
	util.InitLogger(util.LogLevelInfo)

	util.Info("======= 开发环境配置工具 =======")
	util.Info("Version: 1.0.0")
	util.Info("支持: Windows, macOS")
	util.Info("")

	util.Info("开始获取用户配置...")
	answers, err := args_ask.Ask()
	if err != nil {
		util.Error("获取用户配置失败: %v", err)
		panic(err)
	}

	util.Info("用户配置获取完成")
	util.Info("\n\n ##################开始执行开发环境环境变量配置程序################## ")

	// 定义命令行参数
	configPath := answers.ConfigFilePath
	osName := answers.OsName

	if configPath == "" {
		configPath = filepath.Join(util.FindCurrentDir(), "config", "config.json")
	}

	util.Debug("操作系统: %s", osName)
	util.Debug("配置文件路径: %s", configPath)

	util.Info("配置文件：%s", configPath)

	// 初始化平台
	util.Info("初始化平台处理器...")

	// 创建处理器
	util.Info("创建环境处理器...")
	processor := env.NewEnvProcessorByCurrentOsName(osName, configPath, answers.Envs)
	processor.Process()

	util.Info("\n环境变量配置完成, 请重启终端或重新登录系统!")
}
