package env

import (
	"auto_dev_env/src/cmd"
	"auto_dev_env/src/file"
	"encoding/json"
	"errors"
	"log"
	"path/filepath"
	"strings"
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
func NewEnvProcessor(osName, configPath string, cmdProcessor cmd.Processor, fileProcessor file.Processor) Processor {

	if cmdProcessor == nil {
		panic("cmdProcessor is nil")
	}
	if fileProcessor == nil {
		panic("fileProcessor is nil")
	}

	fileBytes, err := fileProcessor.ReadFile(configPath)
	if err != nil {
		panic(err)
	}

	var config []ConfigEnv

	err = json.Unmarshal(fileBytes, &config)
	if err != nil {
		panic(err)
	}

	return Processor{
		OsName:  osName,
		CP:      cmdProcessor,
		FP:      fileProcessor,
		Configs: config,
	}
}

func (p Processor) Process() {

	err := p.checkAndCopy()
	if err != nil {
		log.Println("[env.processor#Process]" + err.Error())
		return
	}

	err = p.createEnvs()
	if err != nil {
		log.Println("[env.processor#Process]" + err.Error())
		return
	}
}

// check 检查文件是否齐全
func (p Processor) checkAndCopy() error {

	var errorMsg []error

	for _, config := range p.Configs {
		log.Println("->[env.processor#check]" + config.PrintString())

		sourcePath := config.EnvSourcePath

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
				break
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

func (p Processor) createEnvs() error {

	for _, config := range p.Configs {
		log.Println("->[env.processor#createEnvs]" + config.PrintString())

		sourcePath := config.EnvSourcePath

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

			// TODO: 暂时不处理 PATH 评估是否需要
			if ec.AppendPath {
				// do nothing
			}
		}
	}

	return nil
}
