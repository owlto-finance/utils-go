package log

import (
	"context"
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

	// Get project name
	projectName := getProjectName()

	// Get user's home directory
	homeDir, err := os.UserHomeDir()
	if err != nil {
		panic("Failed to get user home directory: " + err.Error())
	}

	// Create project log directory
	logDir := filepath.Join(homeDir, "logs", projectName)
	if err = os.MkdirAll(logDir, os.ModePerm); err != nil {
		panic("Failed to create log directory: " + err.Error())
	}

	logFilePath := filepath.Join(logDir, "app.log")

	// Configure Lumberjack
	lumberjackLogger := &lumberjack.Logger{
		Filename:   logFilePath,
		MaxSize:    100, // Max size in MB
		MaxBackups: 3,
		MaxAge:     28, // Max age in days
	}

	// Use custom formatter
	customFormatter := &CustomFormatter{EnableColors: true}

	// Set log output to console with colors
	log.SetOutput(os.Stdout)
	log.SetFormatter(customFormatter)
	log.SetLevel(logrus.InfoLevel)

	// Set timezone to Beijing time
	log.AddHook(&TimezoneHook{Location: "Asia/Shanghai"})

	// Add file output hook with custom formatter
	fileFormatter := &CustomFormatter{EnableColors: false}
	log.AddHook(&Hook{Writer: lumberjackLogger, Formatter: fileFormatter, LogLevels: logrus.AllLevels})
}

func getProjectName() string {
	// Get executable path
	exePath, err := os.Executable()
	if err != nil {
		panic("Failed to get executable path: " + err.Error())
	}

	// Extract project name from executable path
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
	// Get caller file and line number
	file, line := findCaller()
	file = filepath.Base(file)

	// Format timestamp, filename, line number, and log level
	timestamp := entry.Time.Format("2006-01-02 15:04:05.000000")
	level := strings.ToUpper(entry.Level.String())
	if level == "WARNING" {
		level = "WARN"
	}

	logId := entry.Data["logId"]
	logIdStr := ""
	if logId != nil && logId != "unknown" {
		logIdStr = fmt.Sprintf("%v ", logId)
	}

	if f.EnableColors {
		levelColor := getColorByLevel(entry.Level)
		msg := fmt.Sprintf("%s %s:%d %s[%s%s%s] %s\n", timestamp, file, line, logIdStr, levelColor, level, "\033[0m", entry.Message)
		return []byte(msg), nil
	}

	msg := fmt.Sprintf("%s %s:%d %s[%s] %s\n", timestamp, file, line, logIdStr, level, entry.Message)
	return []byte(msg), nil
}

func getColorByLevel(level logrus.Level) string {
	switch level {
	case logrus.DebugLevel:
		return "\033[37m" // White
	case logrus.InfoLevel:
		return "\033[32m" // Green
	case logrus.WarnLevel:
		return "\033[33m" // Yellow
	case logrus.ErrorLevel:
		return "\033[31m" // Red
	case logrus.FatalLevel, logrus.PanicLevel:
		return "\033[35m" // Purple
	default:
		return "\033[0m" // Default
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

func CtxWithFields(ctx context.Context) *logrus.Entry {
	requestId := ctx.Value("logId")
	if requestId == nil {
		requestId = "unknown"
	}
	return log.WithField("logId", requestId)
}

// Wrappers for logrus functions
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

func Panic(args ...interface{}) {
	log.Panic(args...)
}

func Panicf(format string, args ...interface{}) {
	log.Panicf(format, args...)
}

// CtxInfo logs an info message with context fields
func CtxInfo(ctx context.Context, message string) {
	CtxWithFields(ctx).Info(message)
}

// CtxInfof logs a formatted info message with context fields
func CtxInfof(ctx context.Context, format string, args ...interface{}) {
	message := fmt.Sprintf(format, args...)
	CtxWithFields(ctx).Info(message)
}

// CtxError logs an error message with context fields
func CtxError(ctx context.Context, message string) {
	CtxWithFields(ctx).Error(message)
}

// CtxErrorf logs a formatted error message with context fields
func CtxErrorf(ctx context.Context, format string, args ...interface{}) {
	message := fmt.Sprintf(format, args...)
	CtxWithFields(ctx).Error(message)
}

// CtxDebug logs a debug message with context fields
func CtxDebug(ctx context.Context, message string) {
	CtxWithFields(ctx).Debug(message)
}

// CtxDebugf logs a formatted debug message with context fields
func CtxDebugf(ctx context.Context, format string, args ...interface{}) {
	message := fmt.Sprintf(format, args...)
	CtxWithFields(ctx).Debug(message)
}

// CtxWarn logs a warning message with context fields
func CtxWarn(ctx context.Context, message string) {
	CtxWithFields(ctx).Warn(message)
}

// CtxWarnf logs a formatted warning message with context fields
func CtxWarnf(ctx context.Context, format string, args ...interface{}) {
	message := fmt.Sprintf(format, args...)
	CtxWithFields(ctx).Warn(message)
}
