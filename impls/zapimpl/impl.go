package zapimpl

import (
	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"

	"github.com/stainton/logger"
)

type zapWrapper struct {
	l *zap.Logger
}

func NewDevelopLogger(opts logger.Options) (logger.Logger, error) {
	zapConfig := zap.NewDevelopmentConfig()
	return newLogger(zapConfig, opts)
}
func NewProductionLogger(opts logger.Options) (logger.Logger, error) {
	zapConfig := zap.NewProductionConfig()
	return newLogger(zapConfig, opts)
}

func newLogger(zapConfig zap.Config, opts logger.Options) (logger.Logger, error) {
	zapConfig.OutputPaths = opts.OutputPaths()
	zapConfig.Encoding = string(opts.Encoding())
	zapConfig.Level = zap.NewAtomicLevelAt(zapcore.Level(opts.MinLevel()))
	logger, err := zapConfig.Build(zap.WithCaller(true))
	if err != nil {
		return nil, err
	}
	return &zapWrapper{l: logger.Sugar()}, nil
}

func (z *zapWrapper) Debug() {

}
