package logger

import (
	"fmt"
	"io"
	"os"
	"runtime"
	"strings"
)

type defaultEnvGetter struct{}

func (de *defaultEnvGetter) Env() string {
	return "ENV"
}

func (de *defaultEnvGetter) Level() string {
	return "LOG_LEVEL"
}

var envNames interface {
	Env() string
	Level() string
} = &defaultEnvGetter{}

type EnvLevel int

func NewEnvLevel(s string) EnvLevel {
	switch s {
	case Prod.String():
		return Prod
	case Stg.String():
		return Stg
	case Dev.String():
		return Dev
	}
	return Dev
}

func (e EnvLevel) String() string {
	switch e {
	case Prod:
		return "production"
	case Stg:
		return "staging"
	case Dev:
		return "development"
	}
	return "development"
}

const (
	Dev EnvLevel = iota
	Stg
	Prod
)

type LogLevel int

func NewLogLevel(s string) LogLevel {
	switch s {
	case "error":
		return Err
	case "warn":
		return Warn
	case "info":
		return Info
	}
	return Err
}

func (l LogLevel) String() string {
	switch l {
	case Err:
		return "ERROR"
	case Warn:
		return "WARN"
	case Info:
		return "INFO"
	case Disable:
		return "Disable"
	}
	return "Unknown"
}

const (
	Disable LogLevel = iota
	Info
	Warn
	Err
)

type Logger struct {
	writetr io.Writer
}

func (_ *Logger) prefix(l LogLevel) string {
	var b strings.Builder
	fmt.Fprint(&b, "[")
	fmt.Fprint(&b, l.String())
	fmt.Fprint(&b, "] ")
	return b.String()
}

func NewLogger() *Logger {
	return &Logger{
		writetr: os.Stdout,
	}
}

func (l *Logger) Infof(format string, v ...interface{}) {
	if NewLogLevel(os.Getenv(envNames.Level())) < Info {
		return
	}
	if Stg <= NewEnvLevel(os.Getenv(envNames.Env())) {
		return
	}
	_, fn, line, _ := runtime.Caller(1)
	var b strings.Builder
	fmt.Fprint(&b, l.prefix(Info))
	fmt.Fprint(&b, fmt.Sprintf("%s:%d -> ", fn, line))
	fmt.Fprint(&b, format)

	fmt.Fprintf(l.writetr, b.String(), v...)
}

func (l *Logger) Warnf(format string, v ...interface{}) {
	if NewLogLevel(os.Getenv(envNames.Level())) < Warn {
		return
	}
	if Prod <= NewEnvLevel(os.Getenv(envNames.Env())) {
		return
	}
	_, fn, line, _ := runtime.Caller(1)
	var b strings.Builder
	fmt.Fprint(&b, l.prefix(Warn))
	fmt.Fprint(&b, fmt.Sprintf("%s:%d -> ", fn, line))
	fmt.Fprint(&b, format)

	fmt.Fprintf(l.writetr, b.String(), v...)
}

func (l *Logger) Errof(format string, v ...interface{}) {
	if NewLogLevel(os.Getenv(envNames.Level())) < Err {
		return
	}
	_, fn, line, _ := runtime.Caller(1)
	var b strings.Builder
	fmt.Fprint(&b, l.prefix(Err))
	fmt.Fprint(&b, fmt.Sprintf("%s:%d -> ", fn, line))
	fmt.Fprint(&b, format)

	fmt.Fprintf(l.writetr, b.String(), v...)
}
