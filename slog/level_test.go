package slog

import (
	"context"
	"testing"
)

func TestSlogLevelEnabled(t *testing.T) {
	ctx := context.Background()

	for ix, c := range []struct {
		currentLevel Level
		testingLevel Level
		enabled      bool
	}{
		{PanicLevel, PanicLevel, true},
		{PanicLevel, FatalLevel, false},
		{PanicLevel, ErrorLevel, false},
		{PanicLevel, InfoLevel, false},
		{PanicLevel, DebugLevel, false},
		{PanicLevel, OKLevel, false},
		{PanicLevel, SuccessLevel, false},
		{PanicLevel, FailLevel, false},

		{FatalLevel, PanicLevel, true},
		{FatalLevel, FatalLevel, true},
		{FatalLevel, ErrorLevel, false},
		{FatalLevel, InfoLevel, false},
		{FatalLevel, DebugLevel, false},
		{FatalLevel, OKLevel, false},
		{FatalLevel, SuccessLevel, false},
		{FatalLevel, FailLevel, false},

		{ErrorLevel, PanicLevel, true},
		{ErrorLevel, FatalLevel, true},
		{ErrorLevel, ErrorLevel, true},
		{ErrorLevel, InfoLevel, false},
		{ErrorLevel, DebugLevel, false},
		{ErrorLevel, OKLevel, false},
		{ErrorLevel, SuccessLevel, false},
		{ErrorLevel, FailLevel, true},

		{WarnLevel, PanicLevel, true},
		{WarnLevel, FatalLevel, true},
		{WarnLevel, ErrorLevel, true},
		{WarnLevel, WarnLevel, true},
		{WarnLevel, InfoLevel, false},
		{WarnLevel, DebugLevel, false},
		{WarnLevel, OKLevel, false},
		{WarnLevel, SuccessLevel, false},
		{WarnLevel, FailLevel, true},

		{InfoLevel, PanicLevel, true},
		{InfoLevel, FatalLevel, true},
		{InfoLevel, ErrorLevel, true},
		{InfoLevel, WarnLevel, true},
		{InfoLevel, InfoLevel, true},
		{InfoLevel, DebugLevel, false},
		{InfoLevel, OKLevel, true},
		{InfoLevel, SuccessLevel, true},
		{InfoLevel, FailLevel, true},

		{DebugLevel, PanicLevel, true},
		{DebugLevel, FatalLevel, true},
		{DebugLevel, ErrorLevel, true},
		{DebugLevel, WarnLevel, true},
		{DebugLevel, InfoLevel, true},
		{DebugLevel, DebugLevel, true},
		{DebugLevel, TraceLevel, false},
		{DebugLevel, OKLevel, true},
		{DebugLevel, SuccessLevel, true},
		{DebugLevel, FailLevel, true},

		{TraceLevel, PanicLevel, true},
		{TraceLevel, FatalLevel, true},
		{TraceLevel, ErrorLevel, true},
		{TraceLevel, WarnLevel, true},
		{TraceLevel, InfoLevel, true},
		{TraceLevel, DebugLevel, true},
		{TraceLevel, TraceLevel, true},
		{TraceLevel, OKLevel, true},
		{TraceLevel, SuccessLevel, true},
		{TraceLevel, FailLevel, true},

		{OffLevel, PanicLevel, false},
		{OffLevel, FatalLevel, false},
		{OffLevel, ErrorLevel, false},
		{OffLevel, WarnLevel, false},
		{OffLevel, InfoLevel, false},
		{OffLevel, DebugLevel, false},
		{OffLevel, TraceLevel, false},
		{OffLevel, OKLevel, false},
		{OffLevel, SuccessLevel, false},
		{OffLevel, FailLevel, false},

		{AlwaysLevel, PanicLevel, true},
		{AlwaysLevel, FatalLevel, true},
		{AlwaysLevel, ErrorLevel, true},
		{AlwaysLevel, WarnLevel, true},
		{AlwaysLevel, InfoLevel, true},
		{AlwaysLevel, DebugLevel, true},
		{AlwaysLevel, TraceLevel, true},
		{AlwaysLevel, OKLevel, true},
		{AlwaysLevel, SuccessLevel, true},
		{AlwaysLevel, FailLevel, true},

		{PanicLevel, OffLevel, false},
		{PanicLevel, AlwaysLevel, true},
		{FatalLevel, OffLevel, false},
		{FatalLevel, AlwaysLevel, true},
		{ErrorLevel, OffLevel, false},
		{ErrorLevel, AlwaysLevel, true},
		{WarnLevel, OffLevel, false},
		{WarnLevel, AlwaysLevel, true},
		{InfoLevel, OffLevel, false},
		{InfoLevel, AlwaysLevel, true},
		{DebugLevel, OffLevel, false},
		{DebugLevel, AlwaysLevel, true},
		{TraceLevel, OffLevel, false},
		{TraceLevel, AlwaysLevel, true},
		{OKLevel, OffLevel, false},
		{OKLevel, AlwaysLevel, true},
		{SuccessLevel, OffLevel, false},
		{SuccessLevel, AlwaysLevel, true},
		{FailLevel, OffLevel, false},
		{FailLevel, AlwaysLevel, true},
	} {
		ret := c.currentLevel.Enabled(ctx, c.testingLevel)
		if ret != c.enabled {
			t.Fatalf("%5d. current level %q, testing level %q should return '%v', but got '%v'", ix, c.currentLevel, c.testingLevel, c.enabled, ret)
		} else {
			t.Logf("%5d,  current level %q, testing level %q | passed.", ix, c.currentLevel, c.testingLevel)
		}
	}
	// Debug("hi debug", "AA", 1.23456789)

}
