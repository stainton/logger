package logger

type Logger interface {
	Debug(msg string)
	Debugf(template string, args ...any)
	Info(msg string)
	Infof(template string, args ...any)
	Warn(msg string)
	Warnf(template string, args ...any)
	Error(msg string)
	Errorf(template string, args ...any)
	Fatal(msg string)
	Fatalf(template string, args ...any)
}
