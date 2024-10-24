package env

import (
	"auto_dev_env/src"
	"strconv"
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

func (c ConfigEnv) PrintString() string {

	return "EnvName:" + c.EnvName + "\t\r" +
		"EnvCode:" + c.EnvCode + "\t\r" +
		"EnvSourcePath:" + c.EnvSourcePath + "\t\r" +
		"EnvTargetPath:" + c.EnvTargetPath + "\t\r"
}

type Checker struct {
	Name string       `json:"name"`
	Type src.FileType `json:"type"`
	Path string       `json:"path"`
}

func (c Checker) PrintString() string {

	return "Name:" + c.Name + "\t\r" +
		"Type:" + string(c.Type) + "\t\r" +
		"Path:" + c.Path + "\t\r"
}

type Config struct {
	Key        string `json:"key"`
	Value      string `json:"value"`
	Cover      bool   `json:"cover"`
	AppendPath bool   `json:"append_path"`
}

func (c Config) PrintString() string {

	return "Key:" + c.Key + "\t\r" +
		"Value:" + c.Value + "\t\r" +
		"Cover:" + strconv.FormatBool(c.Cover) + "\t\r" +
		"AppendPath:" + strconv.FormatBool(c.AppendPath) + "\t\r"
}
