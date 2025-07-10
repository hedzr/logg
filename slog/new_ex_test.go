package slog_test

import (
	"errors"
	"testing"
	"time"

	"github.com/hedzr/logg/slog"
)

func TestNewSubLoggers(t *testing.T) {
	logger1 := slog.Default().New()                               // colorful logger, detached
	logger2 := slog.Default().New().SetJSONMode()                 // json format logger, detached
	logger3 := slog.Default().New().SetColorMode(false)           // logfmt logger, detached
	logger4 := slog.Default().New().Set("attr1", v1, "attr2", v2) // detached
	_, _, _, _ = logger1, logger2, logger3, logger4

	logger5 := slog.New("detached-logger")
	_ = logger5
	logger5 = slog.New("name", slog.WithAttrs(fakeLoggFields()...))
	_ = logger5
	logger5 = slog.New("name", slog.NewAttr("attr1", v1))
	_ = logger5
	logger5 = slog.New("name", slog.Group("group1", slog.Int("attr1", _tenInts[0])))
	_ = logger5
	logger5 = slog.New("name", "attr1", v1, "attr2", v2).SetAttrs(fakeLoggFields()...)
	_ = logger5
	logger5 = slog.New(slog.Int("attr1", _tenInts[0]))
	_ = logger5

	// sublogger name is optional:
	logger := slog.New("child")
	_ = logger

	// child of child
	sl := logger.New()
	sl1 := logger.New("grandson")
	sl2 := logger.WithLevel(slog.InfoLevel) // another child, since v0.7,x and v1

	// since v0.7 and v1, WithXXX can get a new sublogger, SetXXX not.

	// // as a compasion, children of Default()
	// logger = slog.Default().New()
	// sl1 := logger.New()
	// sl2 := logger.WithLevel(slog.InfoLevel)

	logs(sl)
	logs(sl1)
	logs(sl2)

	// logs(logger1)
	// logs(logger)

	t.Logf("\n%v", logger.(interface{ DumpSubloggers() string }).DumpSubloggers())
	t.Logf("\n%v", slog.Default().(interface{ DumpSubloggers() string }).DumpSubloggers())
	t.Logf("sl2: %v", sl2)
}

func logs(l slog.Logger) {
	attrs := smallPackagedAttrs()
	attrsDeep := deepPackagedAttrs()

	println(l.Name(), " =========== logs (new_ex_test.go)")
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
}

func smallPackagedAttrs() []any {
	return []any{
		"Attr1", 3.13,
		slog.NewGroupedAttrEasy("group1", "Aa", 1, "Bbb", "a string", "Cc", 3.732, "D", 2.71828+5.3571i),
		"Time1", time.Now(),
		"Dur1", 2*time.Second + 65*time.Millisecond + 4589*time.Microsecond,
		"Bool", true,
		"BoolFalse", false,
	}
}

func deepPackagedAttrs() []any {
	return []any{
		"Attr1", 3.13,
		slog.Group("group1", "Aa", 1, "Bbb", "a string", "Cc", 3.732,
			slog.Group("subgroup1",
				"Zzz", false, "Y", true, "xXx", "\"string",
			),
			"D", 2.71828+5.3571i),
		"Time1", time.Now(),
		"Dur1", 2*time.Second + 65*time.Millisecond + 4589*time.Microsecond,
		"Bool", true,
		"BoolFalse", false,
	}
}

func fakeLoggFields() []slog.Attr { return loggAttrs }

var loggAttrs = []slog.Attr{
	slog.NewAttr("int", _tenInts[0]),
	slog.NewAttr("ints", _tenInts),
	slog.NewAttr("string", _tenStrings[0]),
	slog.NewAttr("strings", _tenStrings),
	slog.NewAttr("time_", _tenTimes[0]),
	slog.NewAttr("times", _tenTimes),
	// slog..NewAttr("user1", _oneUser),
	// slog..NewAttr("user2", _oneUser),
	// slog..NewAttr("users", _tenUsers),
	slog.NewAttr("error", errExample),
}

func fakeLoggArgs() []any { return loggArgs }

var loggArgs = []any{
	"int", _tenInts[0],
	"ints", _tenInts,
	"string", _tenStrings[0],
	"strings", _tenStrings,
	"time", _tenTimes[0],
	"times", _tenTimes,
	// "user1", _oneUser,
	// "user2", _oneUser,
	// "users", _tenUsers,
	"error", errExample,
}
var (
	errExample = errors.New("fail")

	// _messages   = fakeMessages(1000)
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
	// _oneUser = &user{
	// 	Name:      "Jane Doe",
	// 	Email:     "jane@test.com",
	// 	CreatedAt: time.Date(1980, 1, 1, 12, 0, 0, 0, time.UTC),
	// }
	// _tenUsers = users{
	// 	_oneUser,
	// 	_oneUser,
	// 	_oneUser,
	// 	_oneUser,
	// 	_oneUser,
	// 	_oneUser,
	// 	_oneUser,
	// 	_oneUser,
	// 	_oneUser,
	// 	_oneUser,
	// }
)

var (
	v1 = 3
	v2 = 2.718
)
