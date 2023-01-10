package log

import (
	"encoding/json"
	"fmt"
	"io"
	"os"
	"strconv"
	"time"
)

var (
	defaultLogger  = NewLogger("main")
	defaultLevel   = LevelInfo
	defaultOutputs = []*Output{{Writer: os.Stdout, Formatter: BasicFormatter}}
)

func DefaultLogger() Logger              { return defaultLogger }
func SetDefaultLogger(logger Logger)     { defaultLogger = logger }
func SetDefaultLevel(level Level)        { defaultLevel = level }
func SetDefaultOutput(outputs []*Output) { defaultOutputs = outputs }

func Print(level Level, msg string, args ...any)     { DefaultLogger().Print(level, msg, args...) }
func Printf(level Level, format string, args ...any) { DefaultLogger().Printf(level, format, args...) }
func Fatalf(format string, args ...any)              { DefaultLogger().Fatalf(format, args...) }
func Fatal(msg string, args ...any)                  { DefaultLogger().Fatal(msg, args...) }
func Errorf(format string, args ...any)              { DefaultLogger().Errorf(format, args...) }
func Error(msg string, err error, args ...any)       { DefaultLogger().Error(msg, err, args...) }
func Warningf(format string, args ...any)            { DefaultLogger().Warningf(format, args...) }
func Warning(msg string, args ...any)                { DefaultLogger().Warning(msg, args...) }
func Infof(format string, args ...any)               { DefaultLogger().Infof(format, args...) }
func Info(msg string, args ...any)                   { DefaultLogger().Info(msg, args...) }
func Debugf(format string, args ...any)              { DefaultLogger().Debugf(format, args...) }
func Debug(msg string, args ...any)                  { DefaultLogger().Debug(msg, args...) }
func Tracef(format string, args ...any)              { DefaultLogger().Tracef(format, args...) }
func Trace(msg string, args ...any)                  { DefaultLogger().Trace(msg, args...) }

// ===============================================================================================

type slogger struct {
	module  string
	level   Level
	outputs []*Output
}

func NewLogger(module string) Logger {
	return &slogger{
		module: module,
		level:  -1,
	}
}

func (l *slogger) SetLevel(level Level) {
	l.level = level
}

func (l *slogger) Level() Level {
	if l.level >= LevelError {
		return l.level
	} else {
		return defaultLevel
	}
}

func (l *slogger) SetOutputs(outputs []*Output) {
	l.outputs = outputs
}

func (l *slogger) getOutputs() []*Output {
	if l.outputs == nil {
		return defaultOutputs
	} else {
		return l.outputs
	}
}

func (l *slogger) Print(level Level, msg string, args ...any) {
	if l.Level() >= level {
		for _, output := range l.getOutputs() {
			output.Formatter(output.Writer, time.Now(), level, l.module, msg, args)
		}
	}
}

func (l *slogger) Printf(level Level, format string, args ...any) {
	l.Print(level, fmt.Sprintf(format, args...))
}

func (l *slogger) Fatalf(format string, args ...any) {
	l.Printf(LevelFatal, format, args...)
}

func (l *slogger) Fatal(msg string, args ...any) {
	l.Print(LevelFatal, msg, args...)
}

func (l *slogger) Errorf(format string, args ...any) {
	l.Printf(LevelError, format, args...)
}

func (l *slogger) Error(msg string, err error, args ...any) {
	params := []any{"err", err.Error()}
	args = append(params, args...)
	l.Print(LevelError, msg, args...)
}

func (l *slogger) Warningf(format string, args ...any) {
	l.Printf(LevelWarning, format, args...)
}

func (l *slogger) Warning(msg string, args ...any) {
	l.Print(LevelWarning, msg, args...)
}

func (l *slogger) Infof(format string, args ...any) {
	l.Printf(LevelInfo, format, args...)
}

func (l *slogger) Info(msg string, args ...any) {
	l.Print(LevelInfo, msg, args...)
}

func (l *slogger) Debugf(format string, args ...any) {
	l.Printf(LevelDebug, format, args...)
}

func (l *slogger) Debug(msg string, args ...any) {
	l.Print(LevelDebug, msg, args...)
}

func (l *slogger) Tracef(format string, args ...any) {
	l.Printf(LevelTrace, format, args...)
}

func (l *slogger) Trace(msg string, args ...any) {
	l.Print(LevelTrace, msg, args...)
}

// ================================================================================================

func BasicFormatter(w io.Writer, tm time.Time, level Level, module string, msg string, args []any) {
	fmt.Fprintf(w, `%s %-5s [%s] %s`,
		tm.Format("2006-01-02 15:04:05"),
		level,
		module,
		msg,
	)
	if len(args) > 0 {
		fmt.Fprint(w, ":")
		for i := 0; i < len(args); i += 2 {
			fmt.Fprintf(w, ` %s=`, args[i])
			switch value := args[i+1].(type) {
			case *time.Time:
				fmt.Fprintf(w, `"%s"`, value.Format("2006-01-02 15:04:05"))
			case string:
				fmt.Fprintf(w, `"%s"`, value)
			case []byte:
				fmt.Fprintf(w, `"%s"`, string(value))
			case int, int32, int64:
				fmt.Fprintf(w, "%d", value)
			case float64, float32:
				fmt.Fprintf(w, "%f", value)
			case bool:
				fmt.Fprintf(w, "%v", value)
			case fmt.Stringer:
				fmt.Fprintf(w, `"%s"`, value.String())
			default:
				fmt.Fprintf(w, "%+v", value)
			}
		}
	}
	fmt.Fprintf(w, "\n")
}

func TextFormatter(w io.Writer, tm time.Time, level Level, module string, msg string, args []any) {
	fmt.Fprintf(w, `time="%s" level="%s" module="%s" msg=%s`,
		tm.Format("2006-01-02 15:04:05"),
		level,
		module,
		strconv.Quote(msg),
	)
	for i := 0; i < len(args); i += 2 {
		fmt.Fprintf(w, ` %s=`, args[i])
		switch value := args[i+1].(type) {
		case *time.Time:
			fmt.Fprintf(w, "%s", value.Format("2006-01-02 15:04:05"))
		case string:
			fmt.Fprintf(w, "%s", strconv.Quote(value))
		case []byte:
			fmt.Fprintf(w, "%s", strconv.Quote(string(value)))
		case int, int32, int64:
			fmt.Fprintf(w, "%d", value)
		case float64, float32:
			fmt.Fprintf(w, "%f", value)
		case bool:
			fmt.Fprintf(w, "%v", value)
		case fmt.Stringer:
			fmt.Fprintf(w, "%s", strconv.Quote(value.String()))
		default:
			fmt.Fprintf(w, `"%+v"`, value)
		}
	}
	fmt.Fprintf(w, "\n")
}

func JsonFormatter(w io.Writer, tm time.Time, level Level, module string, msg string, args []any) {
	data := make(map[string]any)
	data["time"] = tm.Format("2006-01-02 15:04:05")
	data["level"] = level
	data["module"] = module
	for i := 0; i < len(args); i += 2 {
		key := args[i].(string)
		switch value := args[i+1].(type) {
		case *time.Time:
			data[key] = value.Format("2006-01-02 15:04:05")
		case string:
			data[key] = value
		case []byte:
			data[key] = string(value)
		case int, int32, int64, float32, float64, bool:
			data[key] = value
		case fmt.Stringer:
			data[key] = value.String()
		default:
			data[key] = fmt.Sprintf(`%v`, value)
		}
	}
	json.NewEncoder(w).Encode(data)
}
