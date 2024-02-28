package slog

import (
	"context"
	"os"
	"testing"
)

func TestIsTty(t *testing.T) {
	t.Logf("IsTty:         %v", IsTty(os.Stdout))
	t.Logf("IsColoredTty:  %v", IsColoredTty(os.Stdout))
	t.Logf("IsTtyEscaped:  %v", IsTtyEscaped(ct.wrapDimColor("hello")))
	t.Logf("IsAnsiEscaped: %v", IsAnsiEscaped(ct.wrapDimColor("hello")))
	t.Logf("StripEscapes:  %v", StripEscapes(ct.wrapDimColor("hello")))
	c, r := GetTtySize()
	t.Logf("cols, rows = %v, %v", c, r)
}

func TestPanic(t *testing.T) {
	defer SaveFlagsAndMod(LnoInterrupt)()
	Panic("panic msg")
	Fatal("fatal msg")
	PanicContext(context.TODO(), "panic msg")
	FatalContext(context.TODO(), "fatal msg")
}

func TestLogContext(t *testing.T) {
	defer SaveFlagsAndMod(LnoInterrupt)()
	defer SaveLevelAndSet(TraceLevel)()

	ErrorContext(context.TODO(), "error msg")
	WarnContext(context.TODO(), "warn msg")
	InfoContext(context.TODO(), "info msg")
	DebugContext(context.TODO(), "debug msg")
	TraceContext(context.TODO(), "trace msg")
	PrintContext(context.TODO(), "print msg")
	PrintlnContext(context.TODO(), "print msg")
	OKContext(context.TODO(), "print msg")
	SuccessContext(context.TODO(), "print msg")
	FailContext(context.TODO(), "print msg")

	OKContext(context.TODO(), "print msg",
		Int("int", 3),
		Int8("int8", 3),
		Int16("int16", 3),
		Int32("int32", 3),
		Int64("int64", 3),
		Uint("uint", 3),
		Uint8("uint8", 3),
		Uint16("uint16", 3),
		Uint32("uint32", 3),
		Uint64("uint64", 3),
		Complex64("complex64", 1+2i),
		Complex128("complex128", -1-2i),
		Any("any", "music"),
		Numeric("num", 1-3i),
	)

	OKContext(context.TODO(), "print msg",
		Int("int", 3),
		Group("gorup",
			"attrs", []Attr{NewAttr("a1", 1)},
			[]Attr{NewAttr("a2", -2)},
			"empty-struct", struct{}{},
			struct{}{},
		),
		Numeric("num", 3.14159),
	)

	a1 := buildAttrs(
		Int("int", 3),
		Group("gorup",
			"attrs", []Attr{NewAttr("a1", 1)},
			[]Attr{NewAttr("a2", -2)},
			"empty-struct", struct{}{},
			struct{}{},
		),
		Numeric("num", 3.14159),
	)
	keysKnown := make(map[string]bool)
	a2 := buildUniqueAttrs(keysKnown,
		Int("int", 3),
		Group("gorup",
			"attrs", []Attr{NewAttr("a1", 1)},
			[]Attr{NewAttr("a2", -2)},
			"empty-struct", struct{}{},
			struct{}{},
		),
		Numeric("num", 3.14159),
		"str", "hello",
		[]Attr{NewAttr("a3", -2)},
		Attrs{NewAttr("a5", -2)},
		struct{}{},
	)

	t.Log(a1, a2)

	setUniqueKvp(keysKnown, a1, "a7", 9)
	setUniqueKvp(keysKnown, a2, "str", 9)
}

func TestPrintlnEmpty(t *testing.T) {
	Println()
}

func TestLogCtxCtx(t *testing.T) {
	def := Default()
	defer func() {
		SetDefault(def)
	}()

	l := newentry(nil)
	SetDefault(l)
	logctxctx(context.TODO(), WarnLevel,
		"logctxctx",
	)
}

func TestNewLogLogger(t *testing.T) {
	l := NewLogLogger(New(), DebugLevel)
	l.Print("debug")

	lw := &handlerWriter{New(), DebugLevel, true, 0}
	_, _ = lw.Write([]byte("string"))

	raiseerror("e")
}
