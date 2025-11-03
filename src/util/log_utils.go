package util

import (
	"fmt"
	"log"
	"os"
	"time"

	"github.com/fatih/color"
)

// 日志级别
const (
	LogLevelDebug = "debug"
	LogLevelInfo  = "info"
	LogLevelWarn  = "warn"
	LogLevelError = "error"
)

var (
	logLevel     = LogLevelInfo
	logFile      *os.File
	logFormatter *log.Logger

	// 颜色定义
	colorDebug = color.New(color.FgCyan)
	colorInfo  = color.New(color.FgBlue)
	colorWarn  = color.New(color.FgYellow)
	colorError = color.New(color.FgRed)
)

// InitLogger 初始化日志系统
func InitLogger(level string) {
	logLevel = level
	
	// 创建控制台日志格式化器
	logFormatter = log.New(os.Stdout, "", 0)
	
	// 可以在这里添加文件日志功能
	// logFile, _ = os.OpenFile("app.log", os.O_CREATE|os.O_WRONLY|os.O_APPEND, 0666)
}

// formatLog 格式化日志信息
func formatLog(level, message string) string {
	timestamp := time.Now().Format("2006-01-02 15:04:05")
	return fmt.Sprintf("[%s] [%s] %s", timestamp, level, message)
}

// Debug 输出调试日志
func Debug(format string, args ...interface{}) {
	if logLevel == LogLevelDebug {
		message := fmt.Sprintf(format, args...)
		colored := colorDebug.Sprintf(formatLog("DEBUG", message))
		logFormatter.Println(colored)
	}
}

// Info 输出信息日志
func Info(format string, args ...interface{}) {
	if logLevel == LogLevelDebug || logLevel == LogLevelInfo {
		message := fmt.Sprintf(format, args...)
		colored := colorInfo.Sprintf(formatLog("INFO", message))
		logFormatter.Println(colored)
	}
}

// Warn 输出警告日志
func Warn(format string, args ...interface{}) {
	if logLevel == LogLevelDebug || logLevel == LogLevelInfo || logLevel == LogLevelWarn {
		message := fmt.Sprintf(format, args...)
		colored := colorWarn.Sprintf(formatLog("WARN", message))
		logFormatter.Println(colored)
	}
}

// Error 输出错误日志
func Error(format string, args ...interface{}) {
	message := fmt.Sprintf(format, args...)
	colored := colorError.Sprintf(formatLog("ERROR", message))
	logFormatter.Println(colored)
}

// StepProgress 步骤进度条
func StepProgress(current, total int, message string) {
	progress := float64(current) / float64(total) * 100
	stepStr := colorInfo.Sprintf("[Step %d/%d] %.1f%% - %s", current, total, progress, message)
	logFormatter.Println(stepStr)
}