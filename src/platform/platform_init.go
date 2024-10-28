package platform

import (
	"auto_dev_env/src/cmd"
	"auto_dev_env/src/file"
	"auto_dev_env/src/general"
	"log"
)

func init() {
	log.Printf("plarform init")
	defaultPlatform()

	// TODO: 待添加读取 plugin

}

func defaultPlatform() {
	Register("windows", func() ProcessorPlatform {
		return ProcessorPlatform{
			OsName: "windows",
			CP:     cmd.WinCmd{},
			FP:     file.CommonFileProcessor{},
			OG:     general.WindowsGeneral{},
		}
	})
}
