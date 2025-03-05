package logger

import (
	"fmt"
	"path"
	"time"
)

const (
	_DEBUG = "DEBUG"
	_INFO  = "INFO"
	_WARN  = "WARN"
	_ERROR = "ERROR"
	_FATAL = "FATAL"
)

func logFormater(level, msg string) string {
	now := time.Now().Format("2006-01-02 15:04:05.999")
	file, _, line, ok := caller(3)
	if !ok {
		return fmt.Sprintf("[%s][%s] %s\n", now, level, msg)
	}
	_, fileName := path.Split(file)
	return fmt.Sprintf("[%s][%s][%s:%d] %s\n", now, level, fileName, line, msg)
}

func (l *BaseLogger) Debug(msg string) {
	l.logChan <- logFormater(_DEBUG, msg)
}

func (l *BaseLogger) Debugf(template string, args ...any) {
	msg := fmt.Sprintf(template, args...)
	l.logChan <- logFormater(_DEBUG, msg)
}

func (l *BaseLogger) Info(msg string) {
	l.logChan <- logFormater(_INFO, msg)
}

func (l *BaseLogger) Infof(template string, args ...any) {
	msg := fmt.Sprintf(template, args...)
	l.logChan <- logFormater(_INFO, msg)
}

func (l *BaseLogger) Warn(msg string) {
	l.logChan <- logFormater(_WARN, msg)
}

func (l *BaseLogger) Warnf(template string, args ...any) {
	msg := fmt.Sprintf(template, args...)
	l.logChan <- logFormater(_WARN, msg)
}

func (l *BaseLogger) Error(msg string) {
	l.logChan <- logFormater(_ERROR, msg)
}

func (l *BaseLogger) Errorf(template string, args ...any) {
	msg := fmt.Sprintf(template, args...)
	l.logChan <- logFormater(_ERROR, msg)
}

func (l *BaseLogger) Fatal(msg string) {
	l.logChan <- logFormater(_FATAL, msg)
}

func (l *BaseLogger) Fatalf(template string, args ...any) {
	msg := fmt.Sprintf(template, args...)
	l.logChan <- logFormater(_FATAL, msg)
}
