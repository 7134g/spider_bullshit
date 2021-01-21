package logs

import (
	"fmt"
	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
	"sync"
)

var (
	Logger   *zap.Logger
	once     sync.Once
	Level    = zap.DebugLevel
	FilePath = `./log.log`
	Encode   = `console`
)

func InitLoggerGlobal() {
	once.Do(func() {
		encoderConfig := zapcore.EncoderConfig{
			TimeKey:        "time",
			LevelKey:       "level",
			NameKey:        "logger",
			CallerKey:      "caller",
			MessageKey:     "msg",
			StacktraceKey:  "stacktrace",
			LineEnding:     zapcore.DefaultLineEnding,
			EncodeLevel:    zapcore.LowercaseLevelEncoder, // 小写编码器
			EncodeTime:     zapcore.ISO8601TimeEncoder,    // ISO8601 UTC 时间格式
			EncodeDuration: zapcore.SecondsDurationEncoder,
			EncodeCaller:   zapcore.FullCallerEncoder, // 全路径编码器
		}

		// 设置日志级别
		atom := zap.NewAtomicLevelAt(Level)

		config := zap.Config{
			Level:            atom,                                                // 日志级别
			Development:      true,                                                // 开发模式，堆栈跟踪
			Encoding:         Encode,                                              // 输出格式 console 或 json
			EncoderConfig:    encoderConfig,                                       // 编码器配置
			InitialFields:    map[string]interface{}{"serviceName": "spikeProxy"}, // 初始化字段，如：添加一个服务器名称
			OutputPaths:      []string{"stdout", FilePath},                        // 输出到指定文件 stdout（标准输出，正常颜色） stderr（错误输出，红色）
			ErrorOutputPaths: []string{"stderr"},
		}

		// 构建日志
		logger, err := config.Build()
		if err != nil {
			panic(fmt.Sprintf("log 初始化失败: %v", err))
		}
		Logger = logger
	})
}
