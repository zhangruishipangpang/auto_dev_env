package platform

import (
	"auto_dev_env/src/inter"
	"github.com/fatih/color"
)

// 定义一个机器的平台
// 		1、根据不同平台对文件与环境变量的处理方式不同，标明实现类
//      2、将处理器注册到程序中

var cpb = color.New(color.FgBlue).Add(color.Bold)

type ProcessorPlatform struct {
	OsName string
	CP     inter.CmdProcessor
	FP     inter.FileProcessor
	OG     inter.GenOsGeneral
}

type lazyInitPlatformProcessor = func() ProcessorPlatform

var store map[string]lazyInitPlatformProcessor = make(map[string]lazyInitPlatformProcessor)

func Register(osName string, fnc lazyInitPlatformProcessor) {

	_, exist := store[osName]
	if exist {
		_, _ = cpb.Printf("\n 该类型操作系统 %s 已经存在配置", osName)
	}

	store[osName] = fnc
	_, _ = cpb.Printf("\n 操作系统 %s 配置完成", osName)
}

func GetPlatformProcessor(osName string) ProcessorPlatform {

	fnc, exist := store[osName]
	if !exist {
		_, _ = cpb.Printf("\n 不支持该操作系统 [%s]", osName)
		panic("[platform.platform_store.go#GetPlatformProcessor] 不支持该操作系统:" + osName)
	}

	return fnc()
}
