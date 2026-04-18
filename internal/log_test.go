package internal

import (
	"bytes"
	"strings"
	"testing"
)

func TestInfoWritesAtInfoLevel(t *testing.T) {
	var buf bytes.Buffer
	l := newLogger(&buf, &buf, LevelInfo)

	l.Info("hello %s", "world")

	if !strings.Contains(buf.String(), "hello world") {
		t.Errorf("Info() = %q, want to contain %q", buf.String(), "hello world")
	}
}

func TestInfoSilentAtErrorLevel(t *testing.T) {
	var buf bytes.Buffer
	l := newLogger(&buf, &buf, LevelError)

	l.Info("should not appear")

	if buf.Len() != 0 {
		t.Errorf("Info() wrote %q at LevelError", buf.String())
	}
}

func TestDebugWritesAtDebugLevel(t *testing.T) {
	var buf bytes.Buffer
	l := newLogger(&buf, &buf, LevelDebug)

	l.Debug("debug %d", 42)

	if !strings.Contains(buf.String(), "debug 42") {
		t.Errorf("Debug() = %q, want to contain %q", buf.String(), "debug 42")
	}
}

func TestDebugSilentAtInfoLevel(t *testing.T) {
	var buf bytes.Buffer
	l := newLogger(&buf, &buf, LevelInfo)

	l.Debug("should not appear")

	if buf.Len() != 0 {
		t.Errorf("Debug() wrote %q at LevelInfo", buf.String())
	}
}

func TestTraceWritesAtTraceLevel(t *testing.T) {
	var buf bytes.Buffer
	l := newLogger(&buf, &buf, LevelTrace)

	l.Trace("trace detail %s", "here")

	if !strings.Contains(buf.String(), "trace detail here") {
		t.Errorf("Trace() = %q, want to contain %q", buf.String(), "trace detail here")
	}
}

func TestTraceSilentAtDebugLevel(t *testing.T) {
	var buf bytes.Buffer
	l := newLogger(&buf, &buf, LevelDebug)

	l.Trace("should not appear")

	if buf.Len() != 0 {
		t.Errorf("Trace() wrote %q at LevelDebug", buf.String())
	}
}

func TestErrorWritesToStderr(t *testing.T) {
	var stdout, stderr bytes.Buffer
	l := newLogger(&stdout, &stderr, LevelError)

	l.Error("fail: %s", "oops")

	if stdout.Len() != 0 {
		t.Errorf("Error() wrote to stdout: %q", stdout.String())
	}
	if !strings.Contains(stderr.String(), "fail: oops") {
		t.Errorf("Error() = %q, want to contain %q", stderr.String(), "fail: oops")
	}
}

func TestErrorHasPrefix(t *testing.T) {
	var stdout, stderr bytes.Buffer
	l := newLogger(&stdout, &stderr, LevelError)

	l.Error("something broke")

	if !strings.Contains(stderr.String(), "[ERROR]") {
		t.Errorf("Error() = %q, want to contain [ERROR] prefix", stderr.String())
	}
}

func TestNewLoggerLevels(t *testing.T) {
	tests := []struct {
		verbosity int
		expected  LogLevel
	}{
		{0, LevelError},
		{1, LevelInfo},
		{2, LevelDebug},
		{3, LevelTrace},
		{4, LevelTrace},
	}

	for _, tt := range tests {
		l := NewLogger(tt.verbosity).(*logger)
		if l.level != tt.expected {
			t.Errorf("NewLogger(%d).level = %d, want %d", tt.verbosity, l.level, tt.expected)
		}
	}
}
