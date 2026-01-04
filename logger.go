package logger

type Logger interface {
	DPanic(...any)
	DPanicf(string, ...any)
	DPanicln(...any)
	Fatal(...any)
	Fatalf(string, ...any)
	Fatalln(...any)
	Error(...any)
	Errorf(string, ...any)
	Errorln(...any)
	Warn(...any)
	Warnf(string, ...any)
	Warnln(...any)
	Info(...any)
	Infof(string, ...any)
	Infoln(...any)
	Debug(...any)
	Debugf(string, ...any)
	Debugln(...any)
	Sync() error
}
