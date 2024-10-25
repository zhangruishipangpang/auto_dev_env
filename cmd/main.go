package main

import (
	"auto_dev_env/src/cmd"
	"auto_dev_env/src/env"
	"auto_dev_env/src/file"
	"auto_dev_env/src/platform"
	"auto_dev_env/src/util"
	"flag"
	"log"
	"path/filepath"
)

func main() {

	// 定义命令行参数
	configPath := flag.String("config", "", "config path")

	// 解析命令行参数
	flag.Parse()

	if configPath == nil || *configPath == "" {
		*configPath = filepath.Join(util.FindCurrentDir(), "config", "config.json")
	}

	log.Printf("=====> path " + *configPath)

	platform.Register("win", func() platform.ProcessorPlatform {
		return platform.ProcessorPlatform{
			OsName: "win",
			CP:     cmd.WinCmd{},
			FP:     file.CommonFileProcessor{},
		}
	})

	platformProcessor := platform.GetPlatformProcessor("win")

	processor := env.NewEnvProcessor("win", *configPath, platformProcessor.CP, platformProcessor.FP)
	processor.Process()
}
