package log

import (
	"os"
	"path/filepath"
	"sync"

	"github.com/HoronLee/EchoHub/internal/config"
	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
	"gopkg.in/natefinch/lumberjack.v2"
)

// Logger 日志记录器接口
type Logger struct {
	*zap.Logger
}

var (
	globalLogger *Logger
	mu           sync.RWMutex
)

// GetLogger 获取全局 logger（用于工具函数）
func GetLogger() *Logger {
	mu.RLock()
	defer mu.RUnlock()
	if globalLogger == nil {
		// 返回默认 logger
		return newDefaultLogger()
	}
	return globalLogger
}

// SetGlobalLogger 设置全局 logger（由 DI 系统调用）
func SetGlobalLogger(logger *Logger) {
	mu.Lock()
	defer mu.Unlock()
	globalLogger = logger
}

func newDefaultLogger() *Logger {
	encoderConfig := zapcore.EncoderConfig{
		TimeKey:        "time",
		LevelKey:       "level",
		CallerKey:      "caller",
		MessageKey:     "msg",
		EncodeTime:     zapcore.ISO8601TimeEncoder,
		EncodeDuration: zapcore.SecondsDurationEncoder,
		EncodeCaller:   zapcore.ShortCallerEncoder,
		EncodeLevel:    zapcore.CapitalColorLevelEncoder,
	}
	core := zapcore.NewCore(
		zapcore.NewConsoleEncoder(encoderConfig),
		zapcore.AddSync(os.Stdout),
		zapcore.DebugLevel,
	)
	return &Logger{zap.New(core, zap.AddCaller())}
}

// NewLogger 创建日志记录器
func NewLogger(cfg *config.AppConfig) *Logger {
	mode := cfg.Server.Mode
	var cores []zapcore.Core
	encoderConfig := zapcore.EncoderConfig{
		TimeKey:        "time",
		LevelKey:       "level",
		CallerKey:      "caller",
		MessageKey:     "msg",
		StacktraceKey:  "stacktrace",
		LineEnding:     zapcore.DefaultLineEnding,
		EncodeTime:     zapcore.ISO8601TimeEncoder,
		EncodeDuration: zapcore.SecondsDurationEncoder,
		EncodeCaller:   zapcore.ShortCallerEncoder,
	}

	if mode == "debug" {
		// Debug模式：控制台输出，带颜色，非JSON
		encoderConfig.EncodeLevel = zapcore.CapitalColorLevelEncoder
		consoleEncoder := zapcore.NewConsoleEncoder(encoderConfig)
		cores = append(cores, zapcore.NewCore(
			consoleEncoder,
			zapcore.AddSync(os.Stdout),
			zapcore.DebugLevel,
		))
	} else {
		// Release模式：控制台+文件，JSON格式，无颜色
		encoderConfig.EncodeLevel = zapcore.LowercaseLevelEncoder
		jsonEncoder := zapcore.NewJSONEncoder(encoderConfig)

		// 控制台输出
		cores = append(cores, zapcore.NewCore(
			jsonEncoder,
			zapcore.AddSync(os.Stdout),
			zapcore.InfoLevel,
		))

		// 文件输出
		logDir := "logs"
		os.MkdirAll(logDir, 0755)
		writer := &lumberjack.Logger{
			Filename:   filepath.Join(logDir, "app.log"),
			MaxSize:    100,
			MaxBackups: 5,
			MaxAge:     30,
			Compress:   true,
		}
		cores = append(cores, zapcore.NewCore(
			jsonEncoder,
			zapcore.AddSync(writer),
			zapcore.InfoLevel,
		))
	}

	core := zapcore.NewTee(cores...)
	zapLogger := zap.New(core, zap.AddCaller(), zap.AddStacktrace(zapcore.ErrorLevel))
	logger := &Logger{zapLogger}
	SetGlobalLogger(logger) // 自动设置为全局 logger
	return logger
}
