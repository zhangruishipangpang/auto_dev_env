package env

import (
	"auto_dev_env/src/inter"
	"auto_dev_env/src/platform"
	"auto_dev_env/src/util"
	"encoding/json"
	"errors"
	"github.com/fatih/color"
	"path/filepath"
	"strings"
)

var cpf = color.New(color.FgBlue).Add(color.Bold)
var cpb = color.New(color.FgCyan)
var cpr = color.New(color.BgRed).Add(color.Bold).Add(color.Underline)

// Processor 环境处理器
// 该处理包含了命令处理器与文件处理器，对环境变量的处理操作都在该结构体中实现
type Processor struct {
	OsName     string
	CP         inter.CmdProcessor
	FP         inter.FileProcessor
	OG         inter.GenOsGeneral
	AllConfigs AllConfig
}

// NewEnvProcessor 创建一个环境处理器
func NewEnvProcessor(osName, configPath string, cmdProcessor inter.CmdProcessor, fileProcessor inter.FileProcessor, osGeneral inter.GenOsGeneral) Processor {

	if cmdProcessor == nil {
		panic("cmdProcessor is nil")
	}
	if fileProcessor == nil {
		panic("fileProcessor is nil")
	}

	config := readConfig(configPath, fileProcessor)

	return Processor{
		OsName:     osName,
		CP:         cmdProcessor,
		FP:         fileProcessor,
		OG:         osGeneral,
		AllConfigs: config,
	}
}

// NewEnvProcessorByCurrentOsName 创建一个环境处理器
func NewEnvProcessorByCurrentOsName(osNameArg, configPath string) Processor {

	osName := osNameArg

	if strings.TrimSpace(osName) == "" {
		osName = util.GetCurrentOs()
	}

	if osName == "" {
		panic("cmdProcessor is nil")
	}

	if configPath == "" {
		panic("fileProcessor is nil")
	}

	processorPlatform := platform.GetPlatformProcessor(osName)

	config := readConfig(configPath, processorPlatform.FP)

	return Processor{
		OsName:     osName,
		CP:         processorPlatform.CP,
		FP:         processorPlatform.FP,
		OG:         processorPlatform.OG,
		AllConfigs: config,
	}
}

func readConfig(configPath string, fp inter.FileProcessor) AllConfig {
	fileBytes, err := fp.ReadFile(configPath)
	if err != nil {
		panic(err)
	}

	var config AllConfig

	err = json.Unmarshal(fileBytes, &config)
	if err != nil {
		panic(err)
	}

	return config
}

func (p Processor) Process() {

	err := p.checkAndCopy()
	if err != nil {
		_, _ = cpr.Println(" [env.processor#Process.err(checkAndCopy)]" + err.Error())
		return
	}

	err = p.createEnvs()
	if err != nil {
		_, _ = cpr.Println(" [env.processor#Process.err(createEnvs)]" + err.Error())
		return
	}

	err = p.addPaths()
	if err != nil {
		_, _ = cpr.Println(" [env.processor#Process.err(addPaths)]" + err.Error())
		return
	}
}

// check 检查文件是否齐全
func (p Processor) checkAndCopy() error {

	_, _ = cpf.Printf("\n\n 开始执行 checkAndCopy 节点 ")

	var errorMsg []error

	defaultZipDir := p.AllConfigs.DefaultZipDir

	for _, config := range p.AllConfigs.ConfigEnvs {
		_, _ = cpb.Printf("\n config -> %s", config.PrintString())

		envCode := config.EnvCode
		sourcePath := filepath.Join(config.EnvSourcePath, envCode)
		targetPath := filepath.Join(config.EnvTargetPath, envCode)

		// 如果开启了使用默认配置，则直接覆盖sourcePath配置
		if config.UseDefault {
			err := p.readDefaultZip(defaultZipDir, config)
			if err != nil {
				errorMsg = append(errorMsg, err)
				continue
			}
		}

		checkSuccess := true

		// check
		for _, checkSource := range config.EnvSourceCheck {
			_, _ = cpb.Printf("\n check source: %s", config.PrintString())

			name := checkSource.Name
			path := checkSource.Path
			fileType := checkSource.Type
			if strings.HasPrefix(path, "$") {
				path = filepath.Join(sourcePath, path[1:])
			}

			exist, err := p.FP.Exist(path)
			if err != nil {
				return err
			}
			if !exist {
				checkSuccess = false
				_, _ = cpb.Printf("\n ---->[%s]文件检查未通过", path)
				errorMsg = append(errorMsg, errors.New("检查配置："+name+" "+string(fileType)+"不存在，请检查路径"))
				continue
			}
			_, _ = cpb.Printf("\n ---->[%s]文件检查通过", path)
		}

		// copy
		if !checkSuccess {
			continue
		}

		if targetPath == "" || targetPath == sourcePath {
			_, _ = cpb.Printf("\n\n 无须copy")
			continue
		}

		copyR, err := p.FP.Copy(sourcePath, targetPath, config.DelSource)
		if err != nil {
			return err
		}

		if !copyR {
			return errors.New("----->[env.processor#check] file copy fail")
		}

	}

	if len(errorMsg) == 0 {
		return nil
	}

	return errors.Join(errorMsg...)
}

// readDefaultZip 读取默认的zip配置文件，解压到配置的env_source_path中
func (p Processor) readDefaultZip(defaultZipDir string, env ConfigEnv) error {

	if defaultZipDir == "" {
		return errors.New("default_zip_dir 未配置")
	}

	envName := filepath.Join(defaultZipDir, env.EnvCode)
	envZipName := envName + ".zip"

	exist, err := p.FP.Exist(envZipName)
	if err != nil {
		return err
	}

	if !exist {
		return errors.New(envZipName + " 不存在配置")
	}

	_, _ = cpb.Printf("\n 查找到待解压文件：[%s]", envZipName)

	err = p.FP.UnZip(envZipName, defaultZipDir)
	if err != nil {
		return errors.New("解压文件错误" + err.Error())
	}

	targetCopyPath := filepath.Join(env.EnvTargetPath, env.EnvCode)

	_, err = p.FP.Copy(envName, targetCopyPath, true)
	if err != nil {
		return err
	}

	return nil
}

func (p Processor) createEnvs() error {

	_, _ = cpf.Printf("\n\n 开始执行 createEnvs 节点 ")

	for _, config := range p.AllConfigs.ConfigEnvs {
		_, _ = cpb.Printf("\n config env:    %s", config.PrintString())
		placeholder := config.EnvTargetPath

		for _, ec := range config.EnvConfig {
			_, _ = cpb.Printf("\n env:    %s", config.PrintString())

			existEnv := p.CP.GetEnv(ec.Key)

			if existEnv != "" {
				if !ec.Cover {
					_, _ = cpb.Printf("\n --->已经存在，skip...")
					continue
				}
			}

			value := ec.Value

			// 处理占位符 - sourcePath
			if strings.HasPrefix(value, "$") {
				value = filepath.Join(placeholder, value[1:])
			}

			err := p.CP.SetEnv(ec.Key, value)
			if err != nil {
				return err
			}

			_, _ = cpb.Printf("\n --->[%s] 配置完成", ec.Key)

			// 如果需要添加path，则添加到待添加path列表
			if ec.AppendPath {

				newPath := ec.Key

				// 处理添加到 path 中的后置
				if ec.Suffix != nil && len(ec.Suffix) > 0 {
					newPath = filepath.Join(p.OG.PathMapping(newPath), filepath.Join(ec.Suffix...))
				}

				addPathStore(newPath)
			}
		}
	}

	return nil
}

func (p Processor) addPaths() error {

	_, _ = cpf.Printf("\n\n 开始执行 addPaths 节点 ")

	needAddPaths := getNeedAddPaths()

	if needAddPaths == nil {
		return nil
	}

	path := p.CP.GetEnv("PATH")

	err := p.CP.SetEnv("PATH_BAK", path)
	if err != nil {
		return err
	}

	for _, newPath := range needAddPaths {

		path = p.OG.PathGeneral(path, newPath)
	}

	err = p.CP.SetEnv("PATH", path)
	if err != nil {
		return err
	}
	return nil
}
