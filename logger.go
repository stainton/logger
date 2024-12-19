package logger

import (
	"context"
	"fmt"
	"os"
	"strings"
	"sync"
	"time"

	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
)

type zapWrapper struct {
	index     []zapcore.Field
	zapLogger *zap.Logger
}

func (l *zapWrapper) Debug(msg string) {
	l.zapLogger.Debug(msg, l.index...)
}

func (l *zapWrapper) Debugf(template string, args ...interface{}) {
	msg := fmt.Sprintf(template, args...)
	l.zapLogger.Debug(msg, l.index...)
}

func (l *zapWrapper) Info(msg string) {
	l.zapLogger.Info(msg, l.index...)
}

func (l *zapWrapper) Infof(template string, args ...interface{}) {
	msg := fmt.Sprintf(template, args...)
	l.zapLogger.Info(msg, l.index...)
}

func (l *zapWrapper) Warn(msg string) {
	l.zapLogger.Warn(msg, l.index...)
}

func (l *zapWrapper) Warnf(template string, args ...interface{}) {
	msg := fmt.Sprintf(template, args...)
	l.zapLogger.Warn(msg, l.index...)
}

func (l *zapWrapper) Error(msg string) {
	l.zapLogger.Error(msg, l.index...)
}

func (l *zapWrapper) Errorf(template string, args ...interface{}) {
	msg := fmt.Sprintf(template, args...)
	l.zapLogger.Error(msg, l.index...)
}

func (l *zapWrapper) Fatal(msg string) {
	l.zapLogger.Fatal(msg, l.index...)
}

func (l *zapWrapper) Fatalf(template string, args ...interface{}) {
	msg := fmt.Sprintf(template, args...)
	l.zapLogger.Fatal(msg, l.index...)
}

func (l *zapWrapper) Panic(msg string) {
	l.zapLogger.Panic(msg, l.index...)
}

func (l *zapWrapper) Panicf(template string, args ...interface{}) {
	msg := fmt.Sprintf(template, args...)
	l.zapLogger.Panic(msg, l.index...)
}

func (l *zapWrapper) Sync() {
	l.zapLogger.Sync()
}

// NewLoggerWithWG creates a new logger with periodic log flushing.
//
// It initializes a zap logger wrapped in a zapWrapper struct, which implements the Logger interface.
// The function also starts a goroutine that periodically flushes logs and handles graceful shutdown.
//
// Parameters:
//   - ctx: A context.Context for cancellation and shutdown signaling.
//   - wg: A sync.WaitGroup for coordinating goroutine completion.
//   - servicename: A string identifying the service using this logger.
//   - connectionString: A string specifying the connection details for log output.
//   - filePath: A string specifying the file path for log output.
//
// Returns:
//   - Logger: An interface that provides logging methods.
func NewLoggerWithWG(ctx context.Context, wg *sync.WaitGroup, servicename, connectionString, filePath string) Logger {
	hostName, err := os.Hostname()
	if err != nil {
		hostName = "unknown host"
	}
	hostName = strings.ToLower(hostName)
	core := zapcore.NewCore(
		zapcore.NewJSONEncoder(zap.NewProductionEncoderConfig()),
		zapcore.AddSync(Output(connectionString, filePath)),
		zap.InfoLevel,
	)

	l := &zapWrapper{
		index: []zapcore.Field{
			zap.String("servicename", servicename),
			zap.String("hostname", hostName),
		},
		zapLogger: zap.New(core),
	}
	wg.Add(1)
	go func() {
		ticker := time.NewTicker(5 * time.Second)
		for {
			select {
			case <-ctx.Done():
				l.Sync()
				ticker.Stop()
				wg.Done()
				return
			case <-ticker.C:
				l.Sync()
			}
		}
	}()
	return l
}
