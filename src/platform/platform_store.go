package platform

import (
	"auto_dev_env/src/cmd"
	"auto_dev_env/src/file"
	"log"
)

// 定义一个机器的平台
// 		1、根据不同平台对文件与环境变量的处理方式不同，标明实现类
//      2、将处理器注册到程序中

type ProcessorPlatform struct {
	OsName string
	CP     cmd.Processor
	FP     file.Processor
}

type lazyInitPlatformProcessor = func() ProcessorPlatform

var store map[string]lazyInitPlatformProcessor = make(map[string]lazyInitPlatformProcessor)

func Register(osName string, fnc lazyInitPlatformProcessor) {

	_, exist := store[osName]
	if exist {
		log.Printf("[platform.platform_store.go#Register] %s 已经存在配置", osName)
	}

	store[osName] = fnc
	log.Printf("[platform.platform_store.go#Register] %s 配置完成", osName)
}

func GetPlatformProcessor(osName string) ProcessorPlatform {

	fnc, exist := store[osName]
	if !exist {
		log.Printf("[platform.platform_store.go#GetPlatformProcessor] 不支持该操作系统 [%s]", osName)
	}

	return fnc()
}
