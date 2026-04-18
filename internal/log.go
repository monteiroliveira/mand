package internal

import (
	"io"
	"log"
	"os"
)

type LogLevel int

const (
	LevelError LogLevel = iota
	LevelInfo
	LevelDebug
	LevelTrace
)

type Logger interface {
	Info(format string, args ...any)
	Debug(format string, args ...any)
	Trace(format string, args ...any)
	Error(format string, args ...any)
}

type logger struct {
	out   *log.Logger
	err   *log.Logger
	level LogLevel
}

func NewLogger(verbosity int) Logger {
	level := LevelError
	if verbosity >= 3 {
		level = LevelTrace
	} else if verbosity >= 2 {
		level = LevelDebug
	} else if verbosity >= 1 {
		level = LevelInfo
	}

	return newLogger(os.Stdout, os.Stderr, level)
}

func newLogger(out io.Writer, err io.Writer, level LogLevel) *logger {
	return &logger{
		out:   log.New(out, "", log.LstdFlags),
		err:   log.New(err, "[ERROR] ", log.LstdFlags),
		level: level,
	}
}

func (l *logger) Info(format string, args ...any) {
	if l.level >= LevelInfo {
		l.out.Printf(format, args...)
	}
}

func (l *logger) Debug(format string, args ...any) {
	if l.level >= LevelDebug {
		l.out.Printf(format, args...)
	}
}

func (l *logger) Trace(format string, args ...any) {
	if l.level >= LevelTrace {
		l.out.Printf(format, args...)
	}
}

func (l *logger) Error(format string, args ...any) {
	l.err.Printf(format, args...)
}
