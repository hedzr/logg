package slog

import (
	"context"
	"errors"
	"fmt"
	"io"
	"os"
	"strings"
	"testing"
	"time"
)

func TestNewLogger(t *testing.T) {
	l := New(WithJSONMode(false, false),
		WithColorMode(false),
		WithUTCMode(false, true, false),
		WithTimeFormat("", "", time.RFC3339Nano),
		WithAttrs(Int("a", 1)),
		WithAttrs1(NewAttrs("a", 1)),
		With("b", 2),
	)
	ll := l.(*logimp)

	if ll.owner != nil {
		t.Error("ll.owner should be nil")
	}
	// assert.Nil(t, ll.owner)

	t.Logf("%v", l.GetWriter())
	t.Logf("%v", l.GetWriterBy(DebugLevel))
	t.Logf("%v", GetDefaultWriter())
	t.Logf("%v", GetDefaultLoggersWriter())
}

func TestNewChildLogger(t *testing.T) {
	defer SaveFlagsAndMod(LnoInterrupt | LattrsR)()
	defer SaveLevelAndSet(TraceLevel)()

	l := New()
	ll := l.(*logimp)

	t.Log(ll.owner, ll.Parent())
	if ll.owner != nil && ll.Parent() != nil {
		t.Error("ll.owner should be nil")
	}
	// assert.Nil(t, ll.owner)

	l.Warn("l warn msg", "local", false, "n", l.Name())

	lc1 := l.New("c1").WithAttrs(NewAttr("lc1", true))
	llc1 := lc1 // lc1.(*Entry)
	// assert.Equal(t, llc1.owner, ll.Entry)
	if llc1.owner != ll.Entry {
		t.Error("llc1.owner should equal with ll.Entry")
	}
	if llc1.Root() != ll.Entry {
		t.Errorf("llc1.Root() (%v) should equal with ll.Entry (%v)", llc1.Root(), ll.Entry)
	}

	lc1.Warn("lc1 warn msg", "local", false)

	lc2 := lc1.New("c2").WithAttrs(NewAttr("lc2", true))
	llc2 := lc2 // lc2.(*Entry)
	// assert.Equal(t, llc2.owner, lc1.Entry)
	if llc2.owner != llc1 {
		t.Error("llc2.owner should equal with lc1.Entry")
	}
	if llc2.Root() != ll.Entry {
		t.Errorf("llc2.Root() (%v) should equal with ll.Entry (%v)", llc2.Root(), ll.Entry)
	}

	lc2.Warn("lc2 warn msg", "local", false)
	lc1.Warn("lc1 warn msg again", "local", false)
	lc2.Warn("lc2 warn msg again", "local", false)

	lc3 := lc1.New("c3").WithAttrs(NewAttr("lc3", true), NewAttr("lc1", 1))
	llc3 := lc3 // lc3.(*Entry)
	if llc3.owner != llc1 {
		t.Error("llc3.owner should equal with lc1.Entry")
	}
	if llc3.Root() != ll.Entry {
		t.Errorf("llc3.Root() (%v) should equal with ll.Entry (%v)", llc3.Root(), ll.Entry)
	}

	lc3.Warn("lc3 warn msg", "local", false)
}

func TestLogOneTwoThree(t *testing.T) {
	l := New().WithLevel(InfoLevel)
	l.Info("info msg", "Aa", 1, "Bbb", "a string", "Cc", 3.732, "D", 2.71828+5.3571i)
}

func TestLogLongLine(t *testing.T) {
	l := New()
	l.Info("/ERROR/ Started at 2023-10-07 09:14:53.171853 +0800 CST m=+0.001131696, this is a multi-line test\nVersion: 2.0.1\nPackage: hedzr/hilog", "Aa", 1, "Bbb", "a string", "Cc", 3.732, "D", 2.71828+5.3571i)
}

func TestLogAllVerbs(t *testing.T) {
	defer SaveFlagsAndMod(LnoInterrupt | LattrsR)()
	defer SaveLevelAndSet(TraceLevel)()

	attrs := smallPackagedAttrs()
	attrsDeep := deepPackagedAttrs()

	for _, cas := range []struct {
		name   string
		logger Logger
	}{
		{"json", New(" JSON ").WithLevel(AlwaysLevel).WithJSONMode(true)},     // json format
		{"logfmt", New("LOGFMT").WithLevel(AlwaysLevel).WithColorMode(false)}, // logfmt
		{"colorful", New("COLOUR").WithLevel(AlwaysLevel).WithColorMode()},    // colorful mode
	} {
		t.Log("")
		t.Log("")
		t.Log("")
		t.Logf("Using %q logger...", cas.name)
		t.Log("")

		l := cas.logger
		l.Print("print msg", "Aa", 1, "Bbb", "a string", "Cc", 3.732, "D", 2.71828+5.3571i)
		l.Verbose("verbose msg", "Aa", 1, "Bbb", "a string", "Cc", 3.732, "D", 2.71828+5.3571i)
		l.Trace("trace msg", "Aa", 1, "Bbb", "a string", "Cc", 3.732, "D", 2.71828+5.3571i)
		l.Debug("debug msg", "Aa", 1, "Bbb", "a string", "Cc", 3.732, "D", 2.71828+5.3571i)
		l.Info("info msg", "Aa", 1, "Bbb", "a string", "Cc", 3.732, "D", 2.71828+5.3571i)
		l.Warn("warn msg", "Aa", 1, "Bbb", "a string", "Cc", 3.732, "D", 2.71828+5.3571i)
		l.Warn("very long warning msg As 256-color lookup tables became common on graphic cards, escape sequences were added to select from a pre-defined set of 256 colors.", "Aa", 1, "Bbb", "a string", "Cc", 3.732, "D", 2.71828+5.3571i)
		l.Warn("multi-line warning msg\nAs 256-color lookup tables became common on graphic cards, escape sequences were added to select from a pre-defined set of 256 colors.", "Aa", 1, "Bbb", "a string", "Cc", 3.732, "D", 2.71828+5.3571i)
		l.Info("/ERROR/ Started at 2023-10-07 09:14:53.171853 +0800 CST m=+0.001131696, this is a multi-line test\nVersion: 2.0.1\nPackage: hedzr/hilog", "Aa", 1, "Bbb", "a string", "Cc", 3.732, "D", 2.71828+5.3571i)
		l.Error("error msg", "Aa", 1, "Bbb", "a string", "Cc", 3.732, "D", 2.71828+5.3571i)
		l.Fatal("fatal msg", "Aa", 1, "Bbb", "a string", "Cc", 3.732, "D", 2.71828+5.3571i)
		l.Panic("panic msg", "Aa", 1, "Bbb", "a string", "Cc", 3.732, "D", 2.71828+5.3571i)
		l.Println()
		l.OK("OK msg", "Aa", 1, "Bbb", "a string", "Cc", 3.732, "D", 2.71828+5.3571i)
		l.Success("success msg", "Aa", 1, "Bbb", "a string", "Cc", 3.732, "D", 2.71828+5.3571i)
		l.Fail("fail msg", "Aa", 1, "Bbb", "a string", "Cc", 3.732, "D", 2.71828+5.3571i)
		l.Println()
		l.Warn("warn msg", attrs...)
		l.Info("info msg", attrs...)
		l.Debug("debug msg", attrs...)
		l.Println()
		l.Warn("DEEP warn msg", attrsDeep...)
		l.Info("DEEP info msg", attrsDeep...)
		l.Debug("DEEP debug msg", attrsDeep...)

		// l2 := New().WithLevel(AlwaysLevel).WithColorMode(false)
		// l2.Println()
		// l2.Warn("DEEP warn msg", attrsDeep...)
		// // time="09:59:47.464235+08:00" level="warning" msg="warn msg" Attr1=3.13 Aa=1 Bbb="a string" Cc=3.732 D=(2.71828+5.3571i) Time1="2023-10-22T09:59:47.46414+08:00" Dur1="2.069589s" Bool=true BoolFalse=false
		//
		// l3 := New().WithLevel(AlwaysLevel).WithJSONMode()
		// l3.Println()
		// l3.Warn("DEEP warn msg", attrsDeep...)
	}
}

func TestLogDefault(t *testing.T) {
	defer SaveFlagsAndMod(LnoInterrupt | LattrsR)()
	defer SaveLevelAndSet(TraceLevel)()

	l := New().WithLevel(AlwaysLevel)

	l.Print("print msg", "Aa", 1, "Bbb", "a string", "Cc", 3.732, "D", 2.71828+5.3571i)
	l.Verbose("verbose msg", "Aa", 1, "Bbb", "a string", "Cc", 3.732, "D", 2.71828+5.3571i)
	l.Trace("trace msg", "Aa", 1, "Bbb", "a string", "Cc", 3.732, "D", 2.71828+5.3571i)
	l.Debug("debug msg", "Aa", 1, "Bbb", "a string", "Cc", 3.732, "D", 2.71828+5.3571i)
	l.Info("info msg", "Aa", 1, "Bbb", "a string", "Cc", 3.732, "D", 2.71828+5.3571i)
	l.Warn("warn msg", "Aa", 1, "Bbb", "a string", "Cc", 3.732, "D", 2.71828+5.3571i)
	l.Warn("very long warning msg As 256-color lookup tables became common on graphic cards, escape sequences were added to select from a pre-defined set of 256 colors.", "Aa", 1, "Bbb", "a string", "Cc", 3.732, "D", 2.71828+5.3571i)
	l.Warn("multi-line warning msg\nAs 256-color lookup tables became common on graphic cards, escape sequences were added to select from a pre-defined set of 256 colors.", "Aa", 1, "Bbb", "a string", "Cc", 3.732, "D", 2.71828+5.3571i)
	l.Error("error msg", "Aa", 1, "Bbb", "a string", "Cc", 3.732, "D", 2.71828+5.3571i)
	l.Fatal("fatal msg", "Aa", 1, "Bbb", "a string", "Cc", 3.732, "D", 2.71828+5.3571i)
	l.Panic("panic msg", "Aa", 1, "Bbb", "a string", "Cc", 3.732, "D", 2.71828+5.3571i)
	l.Println()
	l.OK("OK msg", "Aa", 1, "Bbb", "a string", "Cc", 3.732, "D", 2.71828+5.3571i)
	l.Success("success msg", "Aa", 1, "Bbb", "a string", "Cc", 3.732, "D", 2.71828+5.3571i)
	l.Fail("fail msg", "Aa", 1, "Bbb", "a string", "Cc", 3.732, "D", 2.71828+5.3571i)
	l.Println()

	attrs := smallPackagedAttrs()

	l.Warn("warn msg", attrs...)
	l.Warn("warn msg 2", attrs...)
	l.Info("info msg", attrs...)
	l.Debug("debug msg", attrs...)
	// 2023-10-22T09:59:47.464221+08:00 [DBG] debug msg                                               Attr1=3.13 group1.Aa=1 group1.Bbb="a string" group1.Cc=3.732 group1.D=(2.71828+5.3571i) Time1="2023-10-22T09:59:47.46414+08:00" Dur1="2.069589s" Bool=true BoolFalse=false

	l2 := New().WithLevel(AlwaysLevel).WithColorMode(false)
	l2.Println()
	l2.Warn("warn msg", attrs...)
	// time="09:59:47.464235+08:00" level="warning" msg="warn msg" Attr1=3.13 Aa=1 Bbb="a string" Cc=3.732 D=(2.71828+5.3571i) Time1="2023-10-22T09:59:47.46414+08:00" Dur1="2.069589s" Bool=true BoolFalse=false

	l3 := New().WithLevel(AlwaysLevel).WithJSONMode()
	l3.Println()
	l3.Warn("warn msg", attrs...)
	// {"time":"09:59:47.464267+08:00","level":"warning","msg":"warn msg","Attr1":3.13,"group1":{"Aa":1,"Bbb":"a string","Cc":3.732,"D":"(2.71828+5.3571i)"},"Time1":"2023-10-22T09:59:47.46414+08:00","Dur1":"2.069589s","Bool":true,"BoolFalse":false}
}

func TestAddLevelWriter1(t *testing.T) {
	logger := New().AddLevelWriter(InfoLevel, &decorated{os.Stdout})
	logger.Info("info msg")
	logger.Debug(getMessage(0))
}

type decorated struct {
	*os.File
}

func (s *decorated) Write(p []byte) (n int, err error) {
	if s.File != nil {
		if ni, e := s.File.WriteString("[decorated] "); e != nil {
			err = errors.Join(err, e)
		} else {
			n += ni
		}
		if ni, e := s.File.Write(p); e != nil {
			err = errors.Join(err, e)
		} else {
			n += ni
		}
	}
	return
}

func TestDisabledWithoutFields(t *testing.T) {
	defer SaveLevelAndSet(WarnLevel)()
	logger := New().WithLevel(WarnLevel).WithColorMode(false)
	logger.Info(getMessage(0))
	logger.Debug(getMessage(0))
}

func TestWithoutFields(t *testing.T) {
	logger := New().WithLevel(DebugLevel).WithColorMode(false)
	logger.Info(getMessage(0))
	logger.Debug(getMessage(0))
}

func TestWithWriter(t *testing.T) {
	defer SaveFlagsAndMod(LnoInterrupt | LattrsR)()
	defer SaveLevelAndSet(TraceLevel)()

	l := New(WithJSONMode(false, false),
		WithColorMode(false),
		WithUTCMode(false, true, false),
		WithTimeFormat("", "", time.RFC3339Nano),
		WithAttrs(Int("a", 1)),
		WithAttrs1(NewAttrs("a", 1)),
		With("b", 2),

		WithWriter(io.Discard),
		AddWriter(io.Discard),
		WithErrorWriter(io.Discard),
		AddErrorWriter(io.Discard),
		ResetWriters(),

		AddLevelWriter(ErrorLevel, io.Discard),
		RemoveLevelWriter(ErrorLevel, io.Discard),
		ResetLevelWriter(ErrorLevel),
		ResetLevelWriters(),

		WithSkip(1),
	)
	ll := l.(*logimp)

	if ll.owner != nil {
		t.Error("ll.owner should be nil")
	}
	// assert.Nil(t, ll.owner)

	SetSkip(0)

	l.WithValueStringer(nil)
	l.Println("")

	t.Logf("%v", l.Skip())

	t.Logf("%v", l.Enabled(ErrorLevel))

	t.Logf("%v", l.GetWriter())
	t.Logf("%v", l.GetWriterBy(DebugLevel))
	t.Logf("%v", GetDefaultWriter())
	t.Logf("%v", GetDefaultLoggersWriter())

	var ss strings.Builder
	ss.WriteString("val")

	l.WithContextKeys(&ss, "from") // set two keys here: one is a Stringer, another is a string

	ctx := context.WithValue(
		context.WithValue(context.Background(), "from", "consul-center"),
		&ss, "test")

	l.PanicContext(ctx, "panic")
	l.FatalContext(ctx, "fatal")
	l.ErrorContext(ctx, "error")
	l.WarnContext(ctx, "warn")
	l.InfoContext(ctx, "info")
	l.DebugContext(ctx, "debug")
	l.TraceContext(ctx, "trace")
	l.PrintContext(ctx, "print")
	l.PrintlnContext(ctx, "println")
	l.OKContext(ctx, "ok")
	l.SuccessContext(ctx, "success")
	l.FailContext(ctx, "fail")
	l.Log(ctx, AlwaysLevel, "log")
}

//

//

//

func BenchmarkWithoutFields(b *testing.B) {
	b.Logf("Logging without any structured context. [BenchmarkWithoutFields]")
	b.Run("hedzr/logg/slog", func(b *testing.B) {
		logger := newLoggTextMode()
		b.ResetTimer()
		b.RunParallel(func(pb *testing.PB) {
			for pb.Next() {
				logger.Info(getMessage(0))
			}
		})
	})
}

func newLoggTextMode(fields ...Attr) Logger {
	return New().WithLevel(DebugLevel).WithAttrs(fields...).
		// WithLevel(slogg.OffLevel)
		// WithWriter(io.Discard).
		WithColorMode(false)
}

func TestAccumulatedContext(t *testing.T) {
	logger := New().WithAttrs(fakeLoggFields()...).WithColorMode(false)
	logger.Info(getMessage(0))
	logger.Debug(getMessage(0))
}

func TestAccumulatedContextAll(t *testing.T) {
	for _, cas := range []struct {
		name   string
		logger Logger
	}{
		{"json", New(" json ").WithAttrs(fakeLoggFields()...).WithJSONMode(true)},     // json format
		{"logfmt", New("logfmt").WithAttrs(fakeLoggFields()...).WithColorMode(false)}, // logfmt
		{"colorful", New("color ").WithAttrs(fakeLoggFields()...).WithColorMode()},    // colorful mode
	} {
		cas.logger.Info(getMessage(0))
		cas.logger.Println()
	}
}

func smallPackagedAttrs() []any {
	return []any{
		"Attr1", 3.13,
		NewGroupedAttrEasy("group1", "Aa", 1, "Bbb", "a string", "Cc", 3.732, "D", 2.71828+5.3571i),
		"Time1", time.Now(),
		"Dur1", 2*time.Second + 65*time.Millisecond + 4589*time.Microsecond,
		"Bool", true,
		"BoolFalse", false,
	}
}

func deepPackagedAttrs() []any {
	return []any{
		"Attr1", 3.13,
		Group("group1", "Aa", 1, "Bbb", "a string", "Cc", 3.732,
			Group("subgroup1",
				"Zzz", false, "Y", true, "xXx", "\"string",
			),
			"D", 2.71828+5.3571i),
		"Time1", time.Now(),
		"Dur1", 2*time.Second + 65*time.Millisecond + 4589*time.Microsecond,
		"Bool", true,
		"BoolFalse", false,
	}
}

func fakeLoggFields() []Attr { return loggAttrs }

var loggAttrs = []Attr{
	NewAttr("int", _tenInts[0]),
	NewAttr("ints", _tenInts),
	NewAttr("string", _tenStrings[0]),
	NewAttr("strings", _tenStrings),
	NewAttr("time_", _tenTimes[0]),
	NewAttr("times", _tenTimes),
	NewAttr("user1", _oneUser),
	NewAttr("user2", _oneUser),
	NewAttr("users", _tenUsers),
	NewAttr("error", errExample),
}

func fakeLoggArgs() []any { return loggArgs }

var loggArgs = []any{
	"int", _tenInts[0],
	"ints", _tenInts,
	"string", _tenStrings[0],
	"strings", _tenStrings,
	"time", _tenTimes[0],
	"times", _tenTimes,
	"user1", _oneUser,
	"user2", _oneUser,
	"users", _tenUsers,
	"error", errExample,
}

func fakeMessages(n int) []string {
	messages := make([]string, n)
	for i := range messages {
		messages[i] = fmt.Sprintf("Test logging, but use a somewhat realistic message length. (#%v)", i)
	}
	return messages
}

func getMessage(iter int) string {
	return _messages[iter%1000]
}

type users []*user

func (uu users) MarshalLogArray(enc *PrintCtx) (err error) {
	for i := range uu {
		if i > 0 {
			enc.WriteRune(',')
		}
		if e := uu[i].MarshalLogObject(enc); e != nil {
			err = errors.Join(err, e)
		}
	}
	return
}

type user struct {
	Name      string    `json:"name"`
	Email     string    `json:"email"`
	CreatedAt time.Time `json:"created_at"`
}

func (u *user) MarshalLogObject(enc *PrintCtx) (err error) {
	enc.AddString("name", u.Name)
	enc.AppendRune(',')
	enc.AddString("email", u.Email)
	enc.AppendRune(',')
	enc.AddInt64("createdAt", u.CreatedAt.UnixNano())
	return
}

var (
	errExample = errors.New("fail")

	_messages   = fakeMessages(1000)
	_tenInts    = []int{1, 2, 3, 4, 5, 6, 7, 8, 9, 0}
	_tenStrings = []string{"a", "b", "c", "d", "e", "f", "g", "h", "i", "j"}
	_tenTimes   = []time.Time{
		time.Unix(0, 0),
		time.Unix(1, 0),
		time.Unix(2, 0),
		time.Unix(3, 0),
		time.Unix(4, 0),
		time.Unix(5, 0),
		time.Unix(6, 0),
		time.Unix(7, 0),
		time.Unix(8, 0),
		time.Unix(9, 0),
	}
	_oneUser = &user{
		Name:      "Jane Doe",
		Email:     "jane@test.com",
		CreatedAt: time.Date(1980, 1, 1, 12, 0, 0, 0, time.UTC),
	}
	_tenUsers = users{
		_oneUser,
		_oneUser,
		_oneUser,
		_oneUser,
		_oneUser,
		_oneUser,
		_oneUser,
		_oneUser,
		_oneUser,
		_oneUser,
	}
)
