// ideas and some code from https://github.com/asim/go-micro/blob/v4.5.0/logger
package logger

import (
	"fmt"
	"os"
	"strings"
)

type Level int8

const (
	TraceLevel = iota - 2
	DebugLevel
	InfoLevel
	WarnLevel
	ErrorLevel
	FatalLevel
)

func (l Level) Enabled(lvl Level) bool {
	return lvl >= l
}

func GetLevel(lvlString string) (Level, error) {
	switch strings.ToLower(lvlString) {
	case "trace":
		return TraceLevel, nil
	case "debug":
		return DebugLevel, nil
	case "info":
		return InfoLevel, nil
	case "warn":
		return WarnLevel, nil
	case "error":
		return ErrorLevel, nil
	case "fatal":
		return FatalLevel, nil
	}
	return InfoLevel, fmt.Errorf("unknown level %s setting InfoLevel by default", lvlString)
}

func (l Level) String() string {
	switch l {
	case TraceLevel:
		return "trace"
	case DebugLevel:
		return "debug"
	case InfoLevel:
		return "info"
	case WarnLevel:
		return "warn"
	case ErrorLevel:
		return "error"
	case FatalLevel:
		return "fatal"
	}
	return ""
}

type PluggedLogger interface {
	Init(opts ...Option) error
	Log(level Level, args ...interface{})
	Logf(level Level, format string, args ...interface{})
	String() string
	// The Logger options
	Options() Options
}

type Logger struct {
	logger PluggedLogger
}

func New(level string, logger PluggedLogger) *Logger {
	lvl, _ := GetLevel(level)
	logger.Init(WithLevel(lvl))
	return &Logger{logger: logger}
}

func (l *Logger) Info(msg ...interface{}) {
	if !l.logger.Options().Level.Enabled(InfoLevel) {
		return
	}
	l.logger.Log(InfoLevel, msg...)
}

func (l *Logger) Infof(template string, msg ...interface{}) {
	if !l.logger.Options().Level.Enabled(InfoLevel) {
		return
	}
	l.logger.Logf(InfoLevel, template, msg...)
}

func (l *Logger) Trace(args ...interface{}) {
	if !l.logger.Options().Level.Enabled(TraceLevel) {
		return
	}
	l.logger.Log(TraceLevel, args...)
}

func (l *Logger) Tracef(template string, args ...interface{}) {
	if !l.logger.Options().Level.Enabled(TraceLevel) {
		return
	}
	l.logger.Logf(TraceLevel, template, args...)
}

func (l *Logger) Debug(args ...interface{}) {
	if !l.logger.Options().Level.Enabled(DebugLevel) {
		return
	}
	l.logger.Log(DebugLevel, args...)
}

func (l *Logger) Debugf(template string, args ...interface{}) {
	if !l.logger.Options().Level.Enabled(DebugLevel) {
		return
	}
	l.logger.Logf(DebugLevel, template, args...)
}

func (l *Logger) Warn(args ...interface{}) {
	if !l.logger.Options().Level.Enabled(WarnLevel) {
		return
	}
	l.logger.Log(WarnLevel, args...)
}

func (l *Logger) Warnf(template string, args ...interface{}) {
	if !l.logger.Options().Level.Enabled(WarnLevel) {
		return
	}
	l.logger.Logf(WarnLevel, template, args...)
}

func (l *Logger) Error(args ...interface{}) {
	if !l.logger.Options().Level.Enabled(ErrorLevel) {
		return
	}
	l.logger.Log(ErrorLevel, args...)
}

func (l *Logger) Errorf(template string, args ...interface{}) {
	if !l.logger.Options().Level.Enabled(ErrorLevel) {
		return
	}
	l.logger.Logf(ErrorLevel, template, args...)
}

func (l *Logger) Fatal(args ...interface{}) {
	if !l.logger.Options().Level.Enabled(FatalLevel) {
		return
	}
	l.logger.Log(FatalLevel, args...)
	os.Exit(1)
}

func (l *Logger) Fatalf(template string, args ...interface{}) {
	if !l.logger.Options().Level.Enabled(FatalLevel) {
		return
	}
	l.logger.Logf(FatalLevel, template, args...)
	os.Exit(1)
}
