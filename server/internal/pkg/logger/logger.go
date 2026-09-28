package logger

import (
	"io"
	"log/slog"
	"os"
)

// Logger 是项目统一日志对象。
//
// 业务层不需要关心日志具体使用什么实现。
// 目前底层使用 Go 标准库 slog。
type Logger struct {
	*slog.Logger
}

// New 创建 Logger。
//
// development：开发环境使用更容易阅读的 TextHandler。
// production：生产环境使用 JSONHandler，方便日志系统解析。
func New(
	development bool,
) *Logger {

	var handler slog.Handler

	if development {

		// 开发环境：
		// 人直接阅读日志，因此使用 Text 格式。
		handler = slog.NewTextHandler(
			os.Stdout,
			&slog.HandlerOptions{
				Level: slog.LevelDebug,
			},
		)

	} else {

		// 生产环境：
		// 推荐 JSON 格式，方便 ELK、Loki 等日志系统处理。
		handler = slog.NewJSONHandler(
			os.Stdout,
			&slog.HandlerOptions{
				Level: slog.LevelInfo,
			},
		)
	}

	return &Logger{
		Logger: slog.New(handler),
	}
}

// NewWithWriter 主要用于测试。
//
// 这样单元测试可以把日志输出到 bytes.Buffer，
// 而不是直接打印到终端。
func NewWithWriter(
	writer io.Writer,
	development bool,
) *Logger {

	var handler slog.Handler

	if development {

		handler = slog.NewTextHandler(
			writer,
			&slog.HandlerOptions{
				Level: slog.LevelDebug,
			},
		)

	} else {

		handler = slog.NewJSONHandler(
			writer,
			&slog.HandlerOptions{
				Level: slog.LevelInfo,
			},
		)
	}

	return &Logger{
		Logger: slog.New(handler),
	}
}
