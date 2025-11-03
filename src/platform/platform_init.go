package platform

import (
	"auto_dev_env/src/common"
	"auto_dev_env/src/linux"
	"auto_dev_env/src/mac"
	"auto_dev_env/src/windows"
	"github.com/fatih/color"
)

var cpg = color.New(color.FgGreen).Add(color.Bold)

var OsStore []string = make([]string, 0)

func init() {
	//_, _ = cpg.Println("\n platform init")
	defaultPlatform()
	macPlatform()
	linuxPlatform()
}

func defaultPlatform() {
	OsStore = append(OsStore, "windows")
	Register("windows", func() ProcessorPlatform {
		return ProcessorPlatform{
			OsName: "windows",
			CP:     windows.WinCmd{},
			FP:     common.CommonFileProcessor{},
			OG:     windows.WindowsGeneral{},
		}
	})
}

func macPlatform() {
	OsStore = append(OsStore, "macOS")
	Register("macOS", func() ProcessorPlatform {
		return ProcessorPlatform{
			OsName: "macOS",
			CP:     mac.MacCmd{},
			FP:     common.CommonFileProcessor{},
			OG:     mac.MacGeneral{},
		}
	})
}

func linuxPlatform() {
	OsStore = append(OsStore, "linux")
	Register("linux", func() ProcessorPlatform {
		return ProcessorPlatform{
			OsName: "linux",
			CP:     mac.MacCmd{}, // 使用bash
			FP:     common.CommonFileProcessor{},
			OG:     linux.LinuxGeneral{},
		}
	})
}
