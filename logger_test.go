package logger

import (
	"os"
	"testing"

	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
)

func TestMain(t *testing.T) {
	// withZapConfig()
	customLevel()
}

func customLevel() {
	// 创建自定义的zapcore配置
	encoderConfig := zapcore.EncoderConfig{
		TimeKey:        "timestamp",
		LevelKey:       "level",
		MessageKey:     "message",
		CallerKey:      "caller",
		StacktraceKey:  "stacktrace",
		LineEnding:     zapcore.DefaultLineEnding,
		EncodeLevel:    zapcore.CapitalLevelEncoder,
		EncodeTime:     zapcore.ISO8601TimeEncoder,
		EncodeDuration: zapcore.SecondsDurationEncoder,
		EncodeCaller:   zapcore.ShortCallerEncoder,
	}

	// 创建核心的日志级别控制器
	atomicLevel := zap.NewAtomicLevel()
	atomicLevel.SetLevel(zap.InfoLevel)

	// 创建日志核心，指定编码器和输出位置
	core := zapcore.NewCore(
		zapcore.NewJSONEncoder(encoderConfig),
		// zapcore.NewConsoleEncoder(encoderConfig),    // 输出为 console 格式
		zapcore.AddSync(zapcore.AddSync(os.Stdout)), // 输出到控制台
		atomicLevel, // 日志级别
	)

	// 创建 logger
	logger := zap.New(core, zap.AddStacktrace(zap.ErrorLevel))
	defer logger.Sync()

	// 记录日志
	logger.Debug("This is a debug message")
	logger.Info("This is an info message")
	logger.Warn("This is a warning")
	logger.Error("This is an error message")
}

func withZapConfig() {
	config := zap.Config{
		Encoding: "json",                              // 日志格式，可以是 "json" 或 "console"
		Level:    zap.NewAtomicLevelAt(zap.InfoLevel), // 设置日志级别
		// OutputPaths:      []string{"stdout", "/var/log/myapp.log"}, // 输出到控制台和文件
		OutputPaths:      []string{"stdout"},
		ErrorOutputPaths: []string{"stderr"}, // 错误日志输出路径
		EncoderConfig: zapcore.EncoderConfig{
			TimeKey:        "time",                        // 时间字段的键
			LevelKey:       "level",                       // 日志级别字段的键
			NameKey:        "logger",                      // 日志器名称字段的键
			CallerKey:      "caller",                      // 调用者字段的键
			MessageKey:     "msg",                         // 消息字段的键
			StacktraceKey:  "stacktrace",                  // 堆栈跟踪字段的键
			LineEnding:     zapcore.DefaultLineEnding,     // 行结束符
			EncodeLevel:    zapcore.LowercaseLevelEncoder, // 日志级别小写输出
			EncodeTime:     zapcore.ISO8601TimeEncoder,    // 时间格式为 ISO8601
			EncodeDuration: zapcore.StringDurationEncoder, // 持续时间的编码格式
			EncodeCaller:   zapcore.ShortCallerEncoder,    // 调用者信息缩短格式
		},
	}

	// 创建 logger
	l, err := config.Build()
	if err != nil {
		panic(err)
	}
	defer l.Sync()

	// 记录日志
	// l.Info("Custom zap logger started", zap.String("module", "main"), zap.Int("version", 1))
	// l.Error("Custom zap logger started", zap.String("module", "main"), zap.Int("version", 1))
	l.Info("Custom zap logger started")
	l.Error("Custom zap logger started")
}
