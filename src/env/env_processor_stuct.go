package env

import (
	"auto_dev_env/src"
	"strconv"
)

type AllConfig struct {
	DefaultZipDir string      `json:"default_zip_dir"`
	ConfigEnvs    []ConfigEnv `json:"configs"`
}

// ConfigEnv 环境配置信息，从Json文件中读取的自定义配置
type ConfigEnv struct {
	EnvName        string    `json:"env_name"`
	EnvCode        string    `json:"env_code"`
	EnvSourcePath  string    `json:"env_source_path"`
	EnvTargetPath  string    `json:"env_target_path"`
	DelSource      bool      `json:"del_source"`
	UseDefault     bool      `json:"use_default"`
	EnvSourceCheck []Checker `json:"env_source_check"`
	EnvConfig      []Config  `json:"env_config"`
}

func (c ConfigEnv) PrintString() string {

	return "EnvName:" + c.EnvName + " | " +
		"EnvCode:" + c.EnvCode + " | " +
		"EnvSourcePath:" + c.EnvSourcePath + " | " +
		"EnvTargetPath:" + c.EnvTargetPath + " | "
}

type Checker struct {
	Name string       `json:"name"`
	Type src.FileType `json:"type"`
	Path string       `json:"path"`
}

func (c Checker) PrintString() string {

	return "Name:" + c.Name + " | " +
		"Type:" + string(c.Type) + " | " +
		"Path:" + c.Path + " | "
}

type Config struct {
	Key        string   `json:"key"`
	Value      string   `json:"value"`
	Cover      bool     `json:"cover"`
	Suffix     []string `json:"suffix_path"`
	AppendPath bool     `json:"append_path"`
}

func (c Config) PrintString() string {

	return "Key:" + c.Key + " | " +
		"Value:" + c.Value + " | " +
		"Cover:" + strconv.FormatBool(c.Cover) + " | " +
		"AppendPath:" + strconv.FormatBool(c.AppendPath) + " | "
}
