package main

import (
	"auto_dev_env/src/cmd"
	"auto_dev_env/src/env"
	"auto_dev_env/src/file"
	"auto_dev_env/src/general"
	"auto_dev_env/src/platform"
	_ "auto_dev_env/src/platform"
	"auto_dev_env/src/util"
	"flag"
	"log"
	"path/filepath"
)

func main() {

	// 定义命令行参数
	configPath := flag.String("config", "", "config path")
	osName := flag.String("os_name", "", "os name")

	// 解析命令行参数
	flag.Parse()

	if configPath == nil || *configPath == "" {
		*configPath = filepath.Join(util.FindCurrentDir(), "config", "config.json")
	}

	log.Printf("=====> path " + *configPath)

	platform.Register("windows", func() platform.ProcessorPlatform {
		return platform.ProcessorPlatform{
			OsName: "win",
			CP:     cmd.WinCmd{},
			FP:     file.CommonFileProcessor{},
			OG:     general.WindowsGeneral{},
		}
	})

	processor := env.NewEnvProcessorByCurrentOsName(*osName, *configPath)
	processor.Process()
}
