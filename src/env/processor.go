package env

import (
	"auto_dev_env/src/cmd"
	"auto_dev_env/src/file"
	"auto_dev_env/src/general"
	"auto_dev_env/src/platform"
	"auto_dev_env/src/util"
	"encoding/json"
	"errors"
	"log"
	"path/filepath"
	"strings"
)

// Processor 环境处理器
// 该处理包含了命令处理器与文件处理器，对环境变量的处理操作都在该结构体中实现
type Processor struct {
	OsName     string
	CP         cmd.Processor
	FP         file.Processor
	OG         general.OsGeneral
	AllConfigs AllConfig
}

// NewEnvProcessor 创建一个环境处理器
func NewEnvProcessor(osName, configPath string, cmdProcessor cmd.Processor, fileProcessor file.Processor, osGeneral general.OsGeneral) Processor {

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

func readConfig(configPath string, fp file.Processor) AllConfig {
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
		log.Println("[env.processor#Process.err(checkAndCopy)]" + err.Error())
		return
	}

	err = p.createEnvs()
	if err != nil {
		log.Println("[env.processor#Process.err(createEnvs)]" + err.Error())
		return
	}

	err = p.addPaths()
	if err != nil {
		log.Println("[env.processor#Process.err(addPaths)]" + err.Error())
		return
	}
}

// check 检查文件是否齐全
func (p Processor) checkAndCopy() error {

	var errorMsg []error

	defaultZipDir := p.AllConfigs.DefaultZipDir

	for _, config := range p.AllConfigs.ConfigEnvs {
		log.Println("->[env.processor#check]" + config.PrintString())

		sourcePath := config.EnvSourcePath

		if sourcePath != "" {
			exist, err := p.FP.Exist(sourcePath)
			if err != nil {
				errorMsg = append(errorMsg, err)
				continue
			}

			// 资源文件不存在的情况，需要考虑是否将资源文件映射到默认的环境变量配置
			if !exist {
				if config.UseDefault {
					err := p.readDefaultZip(defaultZipDir, config)
					if err != nil {
						errorMsg = append(errorMsg, err)
						continue
					}
				}
			}
		} else {
			// 资源文件不存在的情况，需要考虑是否将资源文件映射到默认的环境变量配置
			if config.UseDefault {
				err := p.readDefaultZip(defaultZipDir, config)
				if err != nil {
					errorMsg = append(errorMsg, err)
					continue
				}
			}
		}

		checkSuccess := true

		// check
		for _, checkSource := range config.EnvSourceCheck {
			log.Println("--->[env.processor#check]" + checkSource.PrintString())

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
				log.Println("----->[env.processor#check]" + path + "文件检查未通过")
				errorMsg = append(errorMsg, errors.New("检查配置："+name+" "+string(fileType)+"不存在，请检查路径"))
				continue
			}
			log.Println("----->[env.processor#check]文件检查通过")
		}

		// copy
		if !checkSuccess {
			continue
		}

		targetPath := config.EnvTargetPath
		if targetPath == "" || targetPath == sourcePath {
			log.Println("----->[env.processor#check] 无须copy")
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

	log.Println("->[env.processor#readDefaultZip]" + envZipName + " - " + env.EnvSourcePath)
	err = p.FP.UnZip(envZipName, defaultZipDir)
	if err != nil {
		return errors.New("解压文件错误" + err.Error())
	}

	_, err = p.FP.Copy(envName, env.EnvSourcePath, true)
	if err != nil {
		return err
	}

	return nil
}

func (p Processor) createEnvs() error {

	for _, config := range p.AllConfigs.ConfigEnvs {
		log.Println("->[env.processor#createEnvs]" + config.PrintString())

		sourcePath := config.EnvTargetPath

		for _, ec := range config.EnvConfig {
			log.Println("--->[env.processor#createEnvs]" + ec.PrintString())

			existEnv := p.CP.GetEnv(ec.Key)

			if existEnv != "" {
				if !ec.Cover {
					log.Println("----->[env.processor#createEnvs] 已经存在，skip...")
					continue
				}
			}

			value := ec.Value
			if strings.HasPrefix(value, "$") {
				value = filepath.Join(sourcePath, value[1:])
			}

			err := p.CP.SetEnv(ec.Key, value)
			if err != nil {
				return err
			}

			log.Printf("----->[env.processor#createEnvs] %s 配置完成\n", ec.Key)

			if ec.AppendPath {
				addPathStore(ec.Key)
			}
		}
	}

	return nil
}

func (p Processor) addPaths() error {

	if true {
		return errors.New("未开启配置PATH")
	}

	needAddPaths := getNeedAddPaths()
	if needAddPaths == nil {
		return nil
	}

	for _, newPath := range needAddPaths {

		path := p.CP.GetEnv("PATH")

		err := p.CP.SetEnv("PATH_BAK", path)
		if err != nil {
			return err
		}

		path = p.OG.PathGeneral(path, newPath)

		err = p.CP.SetEnv("PATH", path)
		if err != nil {
			return err
		}
	}

	return nil
}
