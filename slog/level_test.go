package slog

import (
	"context"
	"testing"

	"github.com/hedzr/is/term/color"
)

func TestSlogLevelEnabled(t *testing.T) {
	ctx := context.Background()

	for ix, c := range []struct {
		currentLevel Level
		testingLevel Level
		enabled      bool
	}{
		{PanicLevel, DebugLevel, true},

		{PanicLevel, PanicLevel, true},
		{PanicLevel, FatalLevel, false},
		{PanicLevel, ErrorLevel, false},
		{PanicLevel, InfoLevel, false},
		{PanicLevel, DebugLevel, true},
		{PanicLevel, OKLevel, false},
		{PanicLevel, SuccessLevel, false},
		{PanicLevel, FailLevel, false},

		{FatalLevel, PanicLevel, true},
		{FatalLevel, FatalLevel, true},
		{FatalLevel, ErrorLevel, false},
		{FatalLevel, InfoLevel, false},
		{FatalLevel, DebugLevel, true},
		{FatalLevel, OKLevel, false},
		{FatalLevel, SuccessLevel, false},
		{FatalLevel, FailLevel, false},

		{ErrorLevel, PanicLevel, true},
		{ErrorLevel, FatalLevel, true},
		{ErrorLevel, ErrorLevel, true},
		{ErrorLevel, InfoLevel, false},
		{ErrorLevel, DebugLevel, true},
		{ErrorLevel, OKLevel, false},
		{ErrorLevel, SuccessLevel, false},
		{ErrorLevel, FailLevel, true},

		{WarnLevel, PanicLevel, true},
		{WarnLevel, FatalLevel, true},
		{WarnLevel, ErrorLevel, true},
		{WarnLevel, WarnLevel, true},
		{WarnLevel, InfoLevel, false},
		{WarnLevel, DebugLevel, true},
		{WarnLevel, OKLevel, false},
		{WarnLevel, SuccessLevel, false},
		{WarnLevel, FailLevel, true},

		{InfoLevel, PanicLevel, true},
		{InfoLevel, FatalLevel, true},
		{InfoLevel, ErrorLevel, true},
		{InfoLevel, WarnLevel, true},
		{InfoLevel, InfoLevel, true},
		{InfoLevel, DebugLevel, true},
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

func TestAllLevels(t *testing.T) {
	t.Log(AllLevels())

	for _, l := range AllLevels() {
		if b, err := l.MarshalText(); err == nil {
			var x Level
			xl := x.UnmarshalText(b)
			t.Logf("MarshalText: %v - %v/%v", string(b), x, xl)
		}
		if b, err := l.MarshalJSON(); err == nil {
			var x Level
			xl := x.UnmarshalJSON(b)
			t.Logf("MarshalJSON: %v - %v/%v", string(b), x, xl)
		}
		t.Log(l)
	}

	for i, c := range []struct {
		from   string
		expect Level
	}{
		{"panic", PanicLevel},
		{"p", PanicLevel},
	} {
		actual, err := ParseLevel(c.from)
		if actual != c.expect || err != nil {
			if actual != PanicLevel {
				t.Fatalf("%5d. expect %v, but got %v (err=%v)", i, c.expect, actual, err)
			}
		}
	}
}

const (
	NoticeLevel1 = Level(18) // A custom level must have a value larger than slog.MaxLevel
	NoticeLevel2 = Level(19)
	// HintLevel   = Level(-8) // Or use a negative number
	// SwellLevel  = Level(12) // Sometimes, you may use the value equal with slog.MaxLevel
)

func TestLevel_ShortTag(t *testing.T) {
	RegisterLevel(NoticeLevel1, "NOTICE1",
		RegWithShortTags([6]string{"", "1", "1T", "1TC", "1OTC", "1OTIC"}),
		RegWithColor(color.FgWhite, color.BgUnderline),
		RegWithTreatedAsLevel(InfoLevel),
		RegWithPrintToErrorDevice(false, true),
	)

	for _, l := range AllLevels() {
		for i := 1; i < MaxLengthShortTag; i++ {
			s := l.ShortTag(i)
			t.Logf("Level %v: %v", l, s)
		}
	}

	for i := 1; i < MaxLengthShortTag; i++ {
		x := Level(123)
		s := x.ShortTag(i)
		t.Logf("Level %v: %v", x, s)
	}

	RegisterLevel(NoticeLevel2, "NOTICE2",
		RegWithShortTags([6]string{"", "2", "2T", "2TC", "2OTC", "2OTIC"}),
		RegWithColor(color.FgWhite, color.NoColor),
		RegWithTreatedAsLevel(InfoLevel),
		RegWithPrintToErrorDevice(false),
	)

	for _, l := range AllLevels() {
		for i := 1; i < MaxLengthShortTag; i++ {
			s := l.ShortTag(i)
			t.Logf("Level %v: %v", l, s)
		}
	}

	RegisterLevel(NoticeLevel2, "NOTICE2",
		RegWithShortTags([6]string{"", "2", "2T", "2TC", "2OTC", "2OTIC"}),
		RegWithColor(color.FgWhite, color.NoColor),
		RegWithTreatedAsLevel(InfoLevel),
		RegWithPrintToErrorDevice(false),
	)
}
