package logger

type EncodingEnum string

const (
	EncodingJSON    EncodingEnum = "json"
	EncodingConsole EncodingEnum = "console"
)

type LogLevelEnum int8

const (
	Debug LogLevelEnum = iota - 1
	Info
	Warn
	Error
	DPanic
	Panic
	Fatal
)

type Options struct {
	outputPaths []string
	encoding    EncodingEnum
	minLevel    LogLevelEnum
}

type Option func(*Options)

func NewOptions(opts ...Option) *Options {
	o := &Options{
		outputPaths: []string{"stdout"},
		encoding:    EncodingConsole,
		minLevel:    Info,
	}
	for _, opt := range opts {
		opt(o)
	}
	return o
}

func (o *Options) OutputPaths() []string {
	return o.outputPaths
}
func (o *Options) Encoding() EncodingEnum {
	return o.encoding
}
func (o *Options) MinLevel() LogLevelEnum {
	return o.minLevel
}

func WithOutputPaths(paths []string) Option {
	return func(o *Options) {
		o.outputPaths = paths
	}
}

func WithEncoding(encoding EncodingEnum) Option {
	return func(o *Options) {
		o.encoding = encoding
	}
}

func WithMinLevel(level LogLevelEnum) Option {
	return func(o *Options) {
		o.minLevel = level
	}
}
