package env

import (
	"auto_dev_env/src/inter"
	"auto_dev_env/src/platform"
	"auto_dev_env/src/util"
	"encoding/json"
	"errors"
	"fmt"
	"path/filepath"
	"strings"

	_ "github.com/fatih/color"
)

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
func NewEnvProcessor(osName, configPath string, interestedEnv []string, cmdProcessor inter.CmdProcessor, fileProcessor inter.FileProcessor, osGeneral inter.GenOsGeneral) Processor {

	if cmdProcessor == nil {
		panic("cmdProcessor is nil")
	}
	if fileProcessor == nil {
		panic("fileProcessor is nil")
	}

	config := readConfig(configPath, interestedEnv, fileProcessor)

	return Processor{
		OsName:     osName,
		CP:         cmdProcessor,
		FP:         fileProcessor,
		OG:         osGeneral,
		AllConfigs: config,
	}
}

// NewEnvProcessorByCurrentOsName 创建一个环境处理器
func NewEnvProcessorByCurrentOsName(osNameArg, configPath string, interestedEnv []string) Processor {

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

	config := readConfig(configPath, interestedEnv, processorPlatform.FP)

	return Processor{
		OsName:     osName,
		CP:         processorPlatform.CP,
		FP:         processorPlatform.FP,
		OG:         processorPlatform.OG,
		AllConfigs: config,
	}
}

func readConfig(configPath string, interestedEnv []string, fp inter.FileProcessor) AllConfig {
	fileBytes, err := fp.ReadFile(configPath)
	if err != nil {
		panic(err)
	}

	var config AllConfig

	err = json.Unmarshal(fileBytes, &config)
	if err != nil {
		panic(err)
	}

	ConfigEnvs := make([]ConfigEnv, 0)
	tmpStore := make(map[string]bool)

	for _, ie := range interestedEnv {
		tmpStore[ie] = true
	}

	for _, c := range config.ConfigEnvs {

		abS, err := filepath.Abs(c.EnvSourcePath)
		if err != nil {
			panic(abS)
		}
		c.EnvSourcePath = abS

		abT, err := filepath.Abs(c.EnvTargetPath)
		if err != nil {
			return AllConfig{}
		}
		c.EnvTargetPath = abT

		if b := tmpStore[c.EnvCode]; b {
			ConfigEnvs = append(ConfigEnvs, c)
		}
	}

	abp, err := filepath.Abs(config.DefaultZipDir)
	if err != nil {
		panic(err)
	}

	var config0 = AllConfig{
		DefaultZipDir: abp,
		ConfigEnvs:    ConfigEnvs,
	}

	return config0
}

// Process 处理环境变量配置，包含进度反馈
func (p Processor) Process() {
	// 初始化日志系统
	util.InitLogger(util.LogLevelDebug)

	// 定义处理步骤
	steps := []struct {
		name string
		fn   func() error
	}{
		{"检查文件", p.checkAndCopy},
		{"创建环境变量", p.createEnvs},
		{"添加路径", p.addPaths},
	}

	util.Info("开始环境变量配置...")
	util.Info("操作系统: %s", p.OsName)

	// 执行所有步骤并显示进度
	for i, step := range steps {
		util.StepProgress(i+1, len(steps), step.name)

		if err := step.fn(); err != nil {
			util.Error("步骤 '%s' 执行失败: %v", step.name, err)
			return
		}

		util.Info("✓ 步骤 '%s' 执行完成", step.name)
	}

	util.Info("\n🎉 环境变量配置全部完成！")
}

// checkAndCopy 检查文件是否齐全并复制
func (p Processor) checkAndCopy() error {
	var errorMsg []error
	defaultZipDir := p.AllConfigs.DefaultZipDir

	util.Debug("默认解压目录: %s", defaultZipDir)
	util.Debug("需要处理的环境配置数量: %d", len(p.AllConfigs.ConfigEnvs))

	for _, config := range p.AllConfigs.ConfigEnvs {
		util.Debug("处理配置: %s", config.EnvName)

		// 检查是否需要配置该配置文件

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
			//_, _ = cpb.Printf("\n check source: %s", config.PrintString())

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
				util.Warn("---->[%s]文件检查未通过", path)
				errorMsg = append(errorMsg, fmt.Errorf("检查配置：%s %s不存在，请检查路径", name, string(fileType)))
				continue
			}
			util.Debug("===>[%s]文件检查通过", path)
		}

		// copy
		if !checkSuccess {
			continue
		}

		if targetPath == "" || targetPath == sourcePath {
			util.Debug("目标路径为空或与源路径相同，无需复制")
			continue
		}

		util.Info("复制文件: 从 %s 到 %s", sourcePath, targetPath)
		copyR, err := p.FP.Copy(sourcePath, targetPath, config.DelSource)
		if err != nil {
			return fmt.Errorf("复制文件失败: %w", err)
		}

		if !copyR {
			return errors.New("文件复制失败")
		}
		util.Info("✓ 文件复制成功")

	}

	fmt.Println()

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

	util.Info("查找默认配置包: %s", envZipName)
	exist, err := p.FP.Exist(envZipName)
	if err != nil {
		return err
	}

	if !exist {
		return fmt.Errorf("%s 不存在配置", envZipName)
	}

	util.Info("✓ 查找到待解压文件：[%s]", envZipName)

	util.Info("开始解压文件...")
	err = p.FP.UnZip(envZipName, defaultZipDir)
	if err != nil {
		return fmt.Errorf("解压文件错误: %w", err)
	}
	util.Info("✓ 文件解压成功")

	targetCopyPath := filepath.Join(env.EnvSourcePath, env.EnvCode)

	util.Info("复制解压后的文件到目标路径...")
	_, err = p.FP.Copy(envName, targetCopyPath, true)
	if err != nil {
		return fmt.Errorf("复制解压文件失败: %w", err)
	}
	util.Info("✓ 配置文件复制完成")

	return nil
}

// createEnvs 创建环境变量
func (p Processor) createEnvs() error {
	util.Debug("开始创建环境变量...")

	for _, config := range p.AllConfigs.ConfigEnvs {
		util.Debug("处理环境配置: %s", config.EnvName)
		placeholder := filepath.Join(config.EnvTargetPath, config.EnvCode)

		for _, ec := range config.EnvConfig {
			//_, _ = cpb.Printf("\n env:    %s", config.PrintString())

			existEnv := p.CP.GetEnv(ec.Key)

			if existEnv != "" {
				if !ec.Cover {
					util.Info("===>变量[%s]已经存在，并且 cover=false，跳过...", ec.Key)
					continue
				}
				util.Warn("变量[%s]已存在，将被覆盖", ec.Key)
			}

			value := ec.Value

			// 处理占位符 - sourcePath
			if strings.HasPrefix(value, "$") {
				value = filepath.Join(placeholder, value[1:])
			}

			util.Info("设置环境变量: %s = %s", ec.Key, value)
			err := p.CP.SetEnv(ec.Key, value)
			if err != nil {
				return fmt.Errorf("设置环境变量失败: %w", err)
			}

			util.Info("✓ 变量[%s]配置完成", ec.Key)

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

	fmt.Println()

	return nil
}

// addPaths 添加路径到PATH环境变量
func (p Processor) addPaths() error {
	util.Debug("开始添加路径到PATH...")

	needAddPaths := getNeedAddPaths()

	if needAddPaths == nil || len(needAddPaths) == 0 {
		util.Info("===> 不需要配置 path")
		return nil
	}

	util.Info("需要添加的路径数量: %d", len(needAddPaths))

	path := p.CP.GetEnv("PATH")

	util.Info("备份当前PATH环境变量")
	err := p.CP.SetEnv("PATH_BAK", path)
	if err != nil {
		return fmt.Errorf("备份PATH失败: %w", err)
	}

	for _, newPath := range needAddPaths {

		path = p.OG.PathGeneral(path, newPath)
	}

	err = p.CP.SetEnv("PATH", path)
	if err != nil {
		return err
	}

	util.Info("===> path 配置完成，已配置的路径：%s", strings.Join(needAddPaths, ","))

	//_, _ = cpf.Printf(" 执行 addPaths 节点完成 ")
	return nil
}
