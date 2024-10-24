package cmd

// Processor 命令行执行
type Processor interface {
	// Cmd 执行命令行
	Cmd(name string, arg ...string) ([]byte, error)
	// SetEnv 添加全局环境变量信息
	SetEnv(key, value string) error
	// GetEnv 获取全局环境变量信息
	GetEnv(key string) string
}
