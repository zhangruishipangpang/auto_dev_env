package platform

import (
	"auto_dev_env/src/common"
	"auto_dev_env/src/windows"
	"github.com/fatih/color"
)

var cpg = color.New(color.FgGreen).Add(color.Bold)

var OsStore []string = make([]string, 0)
var EnvStore []string = []string{"jdk", "maven"}

func init() {
	//_, _ = cpg.Println("\n platform init")
	defaultPlatform()
}

func defaultPlatform() {
	OsStore = append(OsStore, "windows")
	OsStore = append(OsStore, "linux")
	Register("windows", func() ProcessorPlatform {
		return ProcessorPlatform{
			OsName: "windows",
			CP:     windows.WinCmd{},
			FP:     common.CommonFileProcessor{},
			OG:     windows.WindowsGeneral{},
		}
	})
}
