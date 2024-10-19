package slog

import (
	"context"
	logslog "log/slog"
	"testing"
	"time"
)

func TestHandler4LogSlog_Log(t *testing.T) { //nolint:revive
	l := New()

	h := NewSlogHandler(l, &HandlerOptions{
		NoColor:  false,
		NoSource: false,
		JSON:     false,
		Level:    DebugLevel,
	})

	ctx := context.Background()

	for i, c := range []struct {
		holding    Level
		requesting logslog.Level
		expect     bool
	}{
		{PanicLevel, logslog.LevelError, false},
		{FatalLevel, logslog.LevelError, false},
		{ErrorLevel, logslog.LevelError, true},
		{WarnLevel, logslog.LevelError, true},
		{InfoLevel, logslog.LevelError, true},
		{DebugLevel, logslog.LevelError, true},
		{TraceLevel, logslog.LevelError, true},

		{OffLevel, logslog.LevelError, false},
		{AlwaysLevel, logslog.LevelError, true},
		{OKLevel, logslog.LevelError, true},
		{SuccessLevel, logslog.LevelError, true},
		{FailLevel, logslog.LevelError, true},
		{MaxLevel, logslog.LevelError, true},

		{PanicLevel, logslog.LevelWarn, false},
		{FatalLevel, logslog.LevelWarn, false},
		{ErrorLevel, logslog.LevelWarn, false},
		{WarnLevel, logslog.LevelWarn, true},
		{InfoLevel, logslog.LevelWarn, true},
		{DebugLevel, logslog.LevelWarn, true},
		{TraceLevel, logslog.LevelWarn, true},

		{OffLevel, logslog.LevelWarn, false},
		{AlwaysLevel, logslog.LevelWarn, true},
		{OKLevel, logslog.LevelWarn, true},
		{SuccessLevel, logslog.LevelWarn, true},
		{FailLevel, logslog.LevelWarn, true},
		{MaxLevel, logslog.LevelWarn, true},

		{PanicLevel, logslog.LevelInfo, false},
		{FatalLevel, logslog.LevelInfo, false},
		{ErrorLevel, logslog.LevelInfo, false},
		{WarnLevel, logslog.LevelInfo, false},
		{InfoLevel, logslog.LevelInfo, true},
		{DebugLevel, logslog.LevelInfo, true},
		{TraceLevel, logslog.LevelInfo, true},

		{OffLevel, logslog.LevelInfo, false},
		{AlwaysLevel, logslog.LevelInfo, true},
		{OKLevel, logslog.LevelInfo, true},
		{SuccessLevel, logslog.LevelInfo, true},
		{FailLevel, logslog.LevelInfo, true},
		{MaxLevel, logslog.LevelInfo, true},

		{PanicLevel, logslog.LevelDebug, false},
		{FatalLevel, logslog.LevelDebug, false},
		{ErrorLevel, logslog.LevelDebug, false},
		{WarnLevel, logslog.LevelDebug, false},
		{InfoLevel, logslog.LevelDebug, false},
		{DebugLevel, logslog.LevelDebug, true},
		{TraceLevel, logslog.LevelDebug, true},

		{OffLevel, logslog.LevelDebug, false},
		{AlwaysLevel, logslog.LevelDebug, true},
		{OKLevel, logslog.LevelDebug, true},
		{SuccessLevel, logslog.LevelDebug, true},
		{FailLevel, logslog.LevelDebug, true},
		{MaxLevel, logslog.LevelDebug, true},
	} {
		if ll, ok := l.(interface{ WithLevel(l Level) *Entry }); ok {
			ll.WithLevel(c.holding)
		}
		if actual := h.Enabled(ctx, c.requesting); actual {
			t.Logf("%5d. h[%v].Enable(ctx, %v) => expect %v, actual is %v, passed", i+1, l.Level(), c.requesting, c.expect, actual)
			_ = h.WithAttrs([]logslog.Attr{
				logslog.Int("k1", 1),
				logslog.Bool("b1", true),
				logslog.Time("time1", time.Now()),
				logslog.Duration("duration1", time.Second),
				logslog.Uint64("u1", 2),
			}).
				WithGroup("g1").
				WithGroup("g2").
				Handle(ctx, logslog.NewRecord(time.Now(), c.requesting, "test msg", 0))
		}
	}
}

func TestHandler4LogSlog_Enabled(t *testing.T) { //nolint:revive
	l := New()

	h := NewSlogHandler(l, &HandlerOptions{
		NoColor:  false,
		NoSource: false,
		JSON:     false,
		Level:    DebugLevel,
	})

	ctx := context.Background()

	for i, c := range []struct {
		holding    Level
		requesting logslog.Level
		expect     bool
	}{
		{PanicLevel, logslog.LevelError, false},
		{FatalLevel, logslog.LevelError, false},
		{ErrorLevel, logslog.LevelError, true},
		{WarnLevel, logslog.LevelError, true},
		{InfoLevel, logslog.LevelError, true},
		{DebugLevel, logslog.LevelError, true},
		{TraceLevel, logslog.LevelError, true},

		{OffLevel, logslog.LevelError, false},
		{AlwaysLevel, logslog.LevelError, true},
		{OKLevel, logslog.LevelError, true},
		{SuccessLevel, logslog.LevelError, true},
		{FailLevel, logslog.LevelError, true},
		{MaxLevel, logslog.LevelError, true},

		{PanicLevel, logslog.LevelWarn, false},
		{FatalLevel, logslog.LevelWarn, false},
		{ErrorLevel, logslog.LevelWarn, false},
		{WarnLevel, logslog.LevelWarn, true},
		{InfoLevel, logslog.LevelWarn, true},
		{DebugLevel, logslog.LevelWarn, true},
		{TraceLevel, logslog.LevelWarn, true},

		{OffLevel, logslog.LevelWarn, false},
		{AlwaysLevel, logslog.LevelWarn, true},
		{OKLevel, logslog.LevelWarn, true},
		{SuccessLevel, logslog.LevelWarn, true},
		{FailLevel, logslog.LevelWarn, true},
		{MaxLevel, logslog.LevelWarn, true},

		{PanicLevel, logslog.LevelInfo, false},
		{FatalLevel, logslog.LevelInfo, false},
		{ErrorLevel, logslog.LevelInfo, false},
		{WarnLevel, logslog.LevelInfo, false},
		{InfoLevel, logslog.LevelInfo, true},
		{DebugLevel, logslog.LevelInfo, true},
		{TraceLevel, logslog.LevelInfo, true},

		{OffLevel, logslog.LevelInfo, false},
		{AlwaysLevel, logslog.LevelInfo, true},
		{OKLevel, logslog.LevelInfo, true},
		{SuccessLevel, logslog.LevelInfo, true},
		{FailLevel, logslog.LevelInfo, true},
		{MaxLevel, logslog.LevelInfo, true},

		{PanicLevel, logslog.LevelDebug, true},
		{FatalLevel, logslog.LevelDebug, true},
		{ErrorLevel, logslog.LevelDebug, true},
		{WarnLevel, logslog.LevelDebug, true},
		{InfoLevel, logslog.LevelDebug, true},
		{DebugLevel, logslog.LevelDebug, true},
		{TraceLevel, logslog.LevelDebug, true},

		{OffLevel, logslog.LevelDebug, false},
		{AlwaysLevel, logslog.LevelDebug, true},
		{OKLevel, logslog.LevelDebug, true},
		{SuccessLevel, logslog.LevelDebug, true},
		{FailLevel, logslog.LevelDebug, true},
		{MaxLevel, logslog.LevelDebug, true},
	} {
		// 1. hold OffLevel: any request levels are denied
		// 2. hold AlwaysLevel: any request levels are allowed
		// 3. in testing/debugging mode, requesting DebugLevel are always allowed.
		if ll, ok := l.(interface{ SetLevel(l Level) *Entry }); ok {
			ll.SetLevel(c.holding)
		}
		if actual := h.Enabled(ctx, c.requesting); actual != c.expect {
			t.Fatalf("%5d. h[%v].Enable(ctx, %v) => expect %v, but got %v, FAILED!", i+1, l.Level(), c.requesting, c.expect, actual)
		} else {
			t.Logf("%5d. h[%v].Enable(ctx, %v) => expect %v, actual is %v, passed", i+1, l.Level(), c.requesting, c.expect, actual)
		}
	}
}

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
