package log

import (
	"io"
	"time"
)

type Level int

func (l Level) String() string {
	switch l {
	case LevelFatal:
		return "FATAL"
	case LevelError:
		return "ERROR"
	case LevelWarning:
		return "WARN"
	case LevelInfo:
		return "INFO"
	case LevelDebug:
		return "DEBUG"
	case LevelTrace:
		return "TRACE"
	default:
		return "?"
	}
}

const (
	LevelFatal Level = iota
	LevelError
	LevelWarning
	LevelInfo
	LevelDebug
	LevelTrace
)

type Logger interface {
	SetLevel(level Level)
	Level() Level
	SetWriter(io.Writer)
	SetFormatter(Formatter)

	Fatalf(format string, args ...any)
	Fatal(msg string, args ...any)
	Errorf(format string, args ...any)
	Error(msg string, err error, args ...any)
	Warningf(format string, args ...any)
	Warning(msg string, args ...any)
	Infof(format string, args ...any)
	Info(msg string, args ...any)
	Debugf(format string, args ...any)
	Debug(msg string, args ...any)
	Tracef(format string, args ...any)
	Trace(msg string, args ...any)
}

type Formatter func(w io.Writer, time time.Time, level Level, module string, msg string, args ...any)
