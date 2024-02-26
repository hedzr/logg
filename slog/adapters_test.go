package slog

import (
	logslog "log/slog"
	"testing"
)

func TestConvertLevelToLogSlog(t *testing.T) {
	for i, c := range []struct {
		src    Level
		expect logslog.Level
	}{
		{PanicLevel, logslog.LevelError},
		{FatalLevel, logslog.LevelError},
		{ErrorLevel, logslog.LevelError},
		{WarnLevel, logslog.LevelWarn},
		{InfoLevel, logslog.LevelInfo},
		{DebugLevel, logslog.LevelDebug},
		{TraceLevel, logslog.LevelDebug},
		{OffLevel, logslog.LevelInfo},
		{AlwaysLevel, logslog.LevelInfo},
		{OKLevel, logslog.LevelInfo},
		{SuccessLevel, logslog.LevelInfo},
		{FailLevel, logslog.LevelInfo},
		{MaxLevel, logslog.LevelInfo},
	} {
		if actual := convertLevelToLogSlog(c.src); actual != c.expect {
			t.Fatalf("%5d. convertLogSlogLevel(%v) => expect %v, but got %v", i, c.src, c.expect, actual)
		}
	}
}

func TestConvertLogSlogToLevel(t *testing.T) {
	for i, c := range []struct {
		expect Level
		src    logslog.Level
	}{
		{ErrorLevel, logslog.LevelError},
		{WarnLevel, logslog.LevelWarn},
		{InfoLevel, logslog.LevelInfo},
		{DebugLevel, logslog.LevelDebug},
		{AlwaysLevel, logslog.LevelInfo - 1},
	} {
		if actual := convertLogSlogLevel(c.src); actual != c.expect {
			t.Fatalf("%5d. convertLogSlogLevel(%v) => expect %v, but got %v", i, c.src, c.expect, actual)
		}
	}
}
