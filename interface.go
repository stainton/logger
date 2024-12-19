package logger

type Logger interface {
	Debug(msg string)
	Debugf(template string, args ...interface{})
	Info(msg string)
	Infof(template string, args ...interface{})
	Warn(msg string)
	Warnf(template string, args ...interface{})
	Error(msg string)
	Errorf(template string, args ...interface{})
	Fatal(msg string)
	Fatalf(template string, args ...interface{})
	Panic(msg string)
	Panicf(template string, args ...interface{})
	Sync()
}
