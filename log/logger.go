package log

import (
	"fmt"
	"io"
	"os"
	"path/filepath"
	"runtime"
	"strings"
	"time"

	"github.com/sirupsen/logrus"
	"gopkg.in/natefinch/lumberjack.v2"
)

var log *logrus.Logger

func init() {
	log = logrus.New()

	// 获取项目名称
	projectName := getProjectName()

	// 获取用户主目录
	homeDir, err := os.UserHomeDir()
	if err != nil {
		panic("Failed to get user home directory: " + err.Error())
	}

	// 创建项目日志目录
	logDir := filepath.Join(homeDir, "logs", projectName)
	if err = os.MkdirAll(logDir, os.ModePerm); err != nil {
		panic("Failed to create log directory: " + err.Error())
	}

	logFilePath := filepath.Join(logDir, "app.log")

	// 配置 Lumberjack
	lumberjackLogger := &lumberjack.Logger{
		Filename:   logFilePath,
		MaxSize:    100, // 以MB为单位
		MaxBackups: 3,
		MaxAge:     28, // 以天为单位
	}

	// 使用自定义格式化器
	customFormatter := &CustomFormatter{EnableColors: true}

	// 设置日志输出到控制台，并使用颜色
	log.SetOutput(os.Stdout)
	log.SetFormatter(customFormatter)
	log.SetLevel(logrus.InfoLevel)

	// 设置时区为北京时区
	log.AddHook(&TimezoneHook{Location: "Asia/Shanghai"})

	// 添加文件输出的钩子，使用自定义格式化器
	fileFormatter := &CustomFormatter{EnableColors: false}
	log.AddHook(&Hook{Writer: lumberjackLogger, Formatter: fileFormatter, LogLevels: logrus.AllLevels})
}

func getProjectName() string {
	// 获取可执行文件路径
	exePath, err := os.Executable()
	if err != nil {
		panic("Failed to get executable path: " + err.Error())
	}

	// 从可执行文件路径中提取项目名称
	projectName := filepath.Base(exePath)
	if runtime.GOOS == "windows" {
		projectName = strings.TrimSuffix(projectName, ".exe")
	}
	return projectName
}

type Hook struct {
	Writer    io.Writer
	Formatter logrus.Formatter
	LogLevels []logrus.Level
}

func (hook *Hook) Levels() []logrus.Level {
	return hook.LogLevels
}

func (hook *Hook) Fire(entry *logrus.Entry) error {
	line, err := hook.Formatter.Format(entry)
	if err != nil {
		return err
	}
	_, err = hook.Writer.Write(line)
	return err
}

type TimezoneHook struct {
	Location string
}

func (hook *TimezoneHook) Levels() []logrus.Level {
	return logrus.AllLevels
}

func (hook *TimezoneHook) Fire(entry *logrus.Entry) error {
	loc, err := time.LoadLocation(hook.Location)
	if err != nil {
		return err
	}
	entry.Time = entry.Time.In(loc)
	return nil
}

type CustomFormatter struct {
	EnableColors bool
}

func (f *CustomFormatter) Format(entry *logrus.Entry) ([]byte, error) {
	// 获取调用者文件和行号
	file, line := findCaller()
	file = filepath.Base(file)

	// 格式化时间戳、文件名、行号和日志级别
	timestamp := entry.Time.Format("2006-01-02 15:04:05.000000")
	level := strings.ToUpper(entry.Level.String())
	if level == "WARNING" {
		level = "WARN"
	}

	if f.EnableColors {
		levelColor := getColorByLevel(entry.Level)
		msg := fmt.Sprintf("%s %s:%d [%s%s%s] %s\n", timestamp, file, line, levelColor, level, "\033[0m", entry.Message)
		return []byte(msg), nil
	}

	msg := fmt.Sprintf("%s %s:%d [%s] %s\n", timestamp, file, line, level, entry.Message)
	return []byte(msg), nil
}

func getColorByLevel(level logrus.Level) string {
	switch level {
	case logrus.DebugLevel:
		return "\033[37m" // 白色
	case logrus.InfoLevel:
		return "\033[32m" // 绿色
	case logrus.WarnLevel:
		return "\033[33m" // 黄色
	case logrus.ErrorLevel:
		return "\033[31m" // 红色
	case logrus.FatalLevel, logrus.PanicLevel:
		return "\033[35m" // 紫色
	default:
		return "\033[0m" // 默认
	}
}

func findCaller() (string, int) {
	for skip := 8; ; skip++ {
		pc, file, line, ok := runtime.Caller(skip)
		if !ok {
			return "unknown", 0
		}
		funcName := runtime.FuncForPC(pc).Name()
		if !strings.Contains(funcName, "log.") && !strings.Contains(funcName, "logrus") {
			return file, line
		}
	}
}

// logrus 的包装函数
func Info(args ...interface{}) {
	log.Info(args...)
}

func Infof(format string, args ...interface{}) {
	log.Infof(format, args...)
}

func Error(args ...interface{}) {
	log.Error(args...)
}

func Errorf(format string, args ...interface{}) {
	log.Errorf(format, args...)
}

func Debug(args ...interface{}) {
	log.Debug(args...)
}

func Debugf(format string, args ...interface{}) {
	log.Debugf(format, args...)
}

func Warn(args ...interface{}) {
	log.Warn(args...)
}

func Warnf(format string, args ...interface{}) {
	log.Warnf(format, args...)
}

func Fatal(args ...interface{}) {
	log.Fatal(args...)
}

func Fatalf(format string, args ...interface{}) {
	log.Fatalf(format, args...)
}
