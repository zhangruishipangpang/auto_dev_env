package env

import (
	"auto_dev_env/src"
)

// ConfigEnv 环境配置信息，从Json文件中读取的自定义配置
type ConfigEnv struct {
	EnvName        string    `json:"env_name"`
	EnvCode        string    `json:"env_code"`
	EnvSourcePath  string    `json:"env_source_path"`
	EnvTargetPath  string    `json:"env_target_path"`
	EnvSourceCheck []Checker `json:"env_source_check"`
	EnvConfig      []Config  `json:"env_config"`
}

type Checker struct {
	Name string       `json:"name"`
	Type src.FileType `json:"type"`
	Path string       `json:"path"`
}

type Config struct {
	Key        string `json:"key"`
	Value      string `json:"value"`
	Cover      bool   `json:"cover"`
	AppendPath bool   `json:"append_path"`
}
