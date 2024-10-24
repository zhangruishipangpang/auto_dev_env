package env

import (
	"auto_dev_env/src/cmd"
	"auto_dev_env/src/file"
)

// Processor 环境处理器
// 该处理包含了命令处理器与文件处理器，对环境变量的处理操作都在该结构体中实现
type Processor struct {
	OsName  string
	CP      cmd.Processor
	FP      file.Processor
	Configs []ConfigEnv
}

// NewEnvProcessor 创建一个环境处理器
func NewEnvProcessor(osName string, cmdProcessor cmd.Processor, fileProcessor file.Processor) Processor {

	if cmdProcessor == nil {
		panic("cmdProcessor is nil")
	}
	if fileProcessor == nil {
		panic("fileProcessor is nil")
	}

	return Processor{
		OsName: osName,
		CP:     cmdProcessor,
		FP:     fileProcessor,
	}
}

// check 检查文件是否齐全
func (p Processor) check() error {
	return nil
}
