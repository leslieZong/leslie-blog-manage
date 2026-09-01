// Package logger 基于 zap 封装的统一日志器
package logger

import (
	"os"
	"strings"

	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
)

var (
	log     *zap.Logger
	sugar   *zap.SugaredLogger
)

// Init 初始化全局 logger。debug=true 时输出 console 彩色，否则 JSON。
func Init(debug bool) {
	var cfg zap.Config
	if debug {
		cfg = zap.NewDevelopmentConfig()
		cfg.EncoderConfig.EncodeLevel = zapcore.CapitalColorLevel
	} else {
		cfg = zap.NewProductionConfig()
		cfg.EncoderConfig.TimeKey = "ts"
		cfg.EncoderConfig.EncodeTime = zapcore.ISO8601TimeEncoder
	}
	// 开发模式输出到 stderr，生产模式同
	cfg.OutputPaths = []string{"stderr"}
	var err error
	log, err = cfg.Build(zap.AddCallerSkip(0))
	if err != nil {
		// 兜底：直接 panic，配置错误必须暴露
		panic("init zap: " + err.Error())
	}
	sugar = log.Sugar()
}

// L 返回底层 zap.Logger
func L() *zap.Logger {
	if log == nil {
		Init(true)
	}
	return log
}

// S 返回 SugaredLogger，便于格式化日志
func S() *zap.SugaredLogger {
	if sugar == nil {
		Init(true)
	}
	return sugar
}

// Sync 在退出前调用，刷新缓冲
func Sync() {
	if log != nil {
		_ = log.Sync()
	}
}

// Level 解析字符串为 zapcore.Level
func Level(s string) zapcore.Level {
	switch strings.ToLower(s) {
	case "debug":
		return zapcore.DebugLevel
	case "info":
		return zapcore.InfoLevel
	case "error":
		return zapcore.ErrorLevel
	default:
		return zapcore.WarnLevel
	}
}

// Stderr 返回一个 *os.File 用于 zap 输出（便于测试替换）
func Stderr() *os.File { return os.Stderr }
