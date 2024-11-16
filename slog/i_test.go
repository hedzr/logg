package slog_test

import (
	"context"
	"errors"
	"fmt"
	logslog "log/slog"
	"net/http"
	"os"
	"os/signal"
	"runtime/debug"
	"testing"
	"time"

	"github.com/hedzr/is/term/color"

	"github.com/hedzr/logg/slog"
)

// func init() {
// 	slog.AddFlags(slog.Lprivacypath)
// }

func TestSlogBasic(t *testing.T) {
	slog.Debug("Debug message") // should be disabled
	slog.Info("Info message")   // should be disabled because defaultLog is in slog.WarnLevel by default
	// We had added new feature to detect if running in testing
	// or debugging mode, or is a debug build (-tags=delve), and
	// set the initial global logging level to DebugLevel.
	// So the above codes will print lines now.

	slog.Warn("Warning message")
	slog.Error("Error message")

	slog.Print("print")
	slog.Println("println")
	slog.OK("ok")
	slog.Success("success")
	slog.Fail("fail")

	t.Log("ok")
}

func TestSlogBasic1(t *testing.T) {
	defer slog.SaveLevelAndSet(slog.InfoLevel)()

	slog.Debug("Debug message") // should be disabled
	slog.Info("Info message")   // should be disabled because defaultLog is in slog.WarnLevel by default
	// We had added new feature to detect if running in testing
	// or debugging mode, or is a debug build (-tags=delve), and
	// set the initial global logging level to DebugLevel.
	// So the above codes will print lines now.

	slog.Warn("Warning message")
	slog.Error("Error message")

	slog.Print("print")
	slog.Println("println")
	slog.OK("ok")
	slog.Success("success")
	slog.Fail("fail")
	t.Log("ok")
}

func TestSlogBasic2(t *testing.T) {
	t.Logf("1. the level is %v", slog.GetLevel())
	done := make(chan struct{})
	fn := func() {
		defer slog.SaveLevelAndSet(slog.WarnLevel)()

		t.Logf("2.the level is %v", slog.GetLevel())
		slog.Debug("Debug message") // should be enabled now
		slog.Info("Info message")   // should be enabled now
		slog.Warn("Warning message", "level", slog.GetLevel())
		slog.Error("Error message")

		close(done)
	}
	fn()
	<-done
	t.Logf("3. the level is %v", slog.GetLevel())
}

func newMyLogger() *mylogger {
	s := &mylogger{
		slog.New("mylogger").SetLevel(slog.InfoLevel),
		true,
		1, // for ur own Infof, another 1 frame need to be ignored.
	}
	return s
}

type mylogger struct {
	slog.Logger
	SprintfLikeLoggingIsEnabled bool
	mySkip                      int
}

func (s *mylogger) Close() { s.Logger.Close() }

func (s *mylogger) Infof(msg string, args ...any) {
	if s.SprintfLikeLoggingIsEnabled {
		if len(args) > 0 {
			// msg = fmt.Sprintf(msg, args...)

			var data []byte
			data = fmt.Appendf(data, msg, args...)
			s.Logger.SetSkip(s.mySkip)
			s.Logger.Info(string(data))
		} else {
			s.Logger.SetSkip(s.mySkip)
			s.Info(msg)
		}
	}
}

func TestSlogBasic3CustomVerbsInUrOwnLogger(t *testing.T) {
	l := newMyLogger()
	l.Infof("what's wrong with %v", "him")
	l.Info("no matter")
	l.Infof("what's wrong with %v, %v, %v, %v, %v, %v", m1AttrsAsAnySlice()...)
}

func newMyLogger2() *mylogger2 {
	l := slog.New("mylogger").SetLevel(slog.InfoLevel)
	s := &mylogger2{
		l,
		true,
		l.WithSkip(1), // a sublogger here with different skip-frames
	}
	return s
}

// use a standalone sub-logger
type mylogger2 struct {
	slog.Logger
	SprintfLikeLoggingIsEnabled bool
	sl                          slog.Logger
}

func (s *mylogger2) Close() { s.Logger.Close() }

func (s *mylogger2) Infof(msg string, args ...any) {
	if s.SprintfLikeLoggingIsEnabled {
		if len(args) > 0 {
			// msg = fmt.Sprintf(msg, args...)

			var data []byte
			data = fmt.Appendf(data, msg, args...)
			s.sl.Info(string(data))
		} else {
			s.sl.Info(msg)
		}
	}
}

func TestSlogBasic4(t *testing.T) {
	// TestSlogBasic4CustomVerbsInUrOwnLogger here

	l := newMyLogger2()
	l.Infof("what's wrong with %v", "him")
	l.Info("no matter")
	l.Infof("what's wrong with %v, %v, %v, %v, %v, %v", m1AttrsAsAnySlice()...)
}

func TestSlogBasic5(t *testing.T) {
	defer slog.SaveFlagsAndMod(slog.Linterruptalways)()
	defer func() {
		if e := recover(); e != nil {
			t.Logf(`panic error caught and recovered: %v`, e)
		}
	}()

	l := slog.New()
	l.Panic("test a panic logging")
}

func TestSlogJSON(t *testing.T) {
	logger := slog.New().SetJSONMode().SetLevel(slog.DebugLevel)

	logger.Debug("Debug message") //
	logger.Info("Info message")   //
	logger.Warn("Warning message")
	logger.Error("Error message")

	logger.Info(
		"incoming request",
		"method", "GET",
		"time_taken_ms", 158,
		"path", "/hello/world?q=search",
		"status", 200,
		"user_agent", "Googlebot/2.1 (+https://www.google.com/bot.html)",
	)

	// the following codes should work fine

	logger.Info(
		"incoming request",
		"method", "GET",
		"time_taken_ms", // the value for this key is missing
	)

	logger1 := slog.New().SetJSONMode(true).SetLevel(slog.DebugLevel)

	logger1.Debug("Debug message") //
	logger1.Info("Info message")   //
	logger1.Warn("Warning message")
	logger1.Error("Error message")

}

func TestSlogLogfmt(t *testing.T) {
	logger := slog.New().SetLevel(slog.TraceLevel).SetColorMode(false)

	logger.Debug("Debug message")
	logger.Info("Info message")
	logger.Warn("Warning message")
	logger.Error("Error message")

	logger.Info(
		"incoming request",
		slog.String("method", "GET"),
		slog.Int("time_taken_ms", 158),
		slog.String("path", "/hello/world?q=search"),
		slog.Int("status", 200),
		slog.String(
			"user_agent",
			"Googlebot/2.1 (+https://www.google.com/bot.html)",
		),
		slog.Group("memory",
			slog.Int("current", 50),
			slog.Int("min", 20),
			slog.Int("max", 80)),
	)
	logger.Debug("debug, incoming request", m1AttrsAsAnySlice()...)
	logger.Trace("trace, incoming request", m2AttrsAsAnySlice()...)

	logger.TraceContext(context.TODO(), "trace, incoming request", m2AttrsAsAnySlice()...)

	logger.Println()
	logger.VerboseContext(context.TODO(), "verbose, incoming request")
	logger.Verbose("verbose, incoming request")
	slog.VerboseContext(context.TODO(), "verbose, incoming request")
	slog.Verbose("verbose, incoming request")

	logger.Println()
	logger.Debug("debug, incoming request")
	logger.Trace("trace, incoming request")
	slog.Debug("debug, incoming request")
	slog.Trace("trace, incoming request")
	logger.DebugContext(context.TODO(), "debug, incoming request")
	logger.TraceContext(context.TODO(), "trace, incoming request")
	slog.DebugContext(context.TODO(), "debug, incoming request")
	slog.TraceContext(context.TODO(), "trace, incoming request")
}

func TestSlogSetDefault(t *testing.T) {
	// slog.SetLevel(slog.WarnLevel)
	defer slog.SaveLevelAndSet(slog.WarnLevel)
	defer slog.SaveFlagsAndMod(slog.LattrsR)() // add, remove, or set flags

	logger := slog.New().SetJSONMode().SetLevel(slog.InfoLevel).Set(m2AttrsAsAnySlice()...)

	slog.SetDefault(logger)

	slog.Info("Info message") // JSON mode here

	l := logger.WithColorMode() // make a child logger,
	l.Error("Error message")    // and apply to color format

	logger.SetColorMode()
	slog.Error("Error message") // now it's in colorful mode.
}

func TestSlogWithAttrs(t *testing.T) {
	logger := slog.New()

	logger.Info("info message",
		"attr1", 0,
		"attr2", false,
		"attr3", 3.13,
	) // plain key and value pairs
	logger.Info("info message",
		"attr1", 0,
		"attr4", errors.New("simple"),
		"attr3", 3.13,
		slog.NewAttr("attr3", "any styles what u prefer"),
		slog.NewGroupedAttrEasy("group1", "attr1", 13, "attr2", false),
	) // use NewAttr, NewGroupedAttrs
	logger.Info("image uploaded",
		slog.Int("id", 23123),
		slog.Group("properties",
			slog.Int("width", 4000),
			slog.Int("height", 3000),
			slog.String("format", "jpeg"),
		),
	) // use Int, Float, String, Any, ..., and Group

	// mixes all above forms
	logger.Info("image uploaded",
		"attr1", 0,
		"attr2", false,
		"attr3", 3.13,
		"attr4", errors.New("simple"),
		// ,,,
		slog.Group("group1",
			"attr1", 0,
			"attr2", false,
			slog.NewAttr("attr3", "any styles what u prefer"),
			slog.Group("more", "group", []byte("more sub-attrs")),
			"attrN", // unpaired key can work here
		),
		// ...
		slog.Int("id", 23123),
		slog.Group("properties",
			slog.Int("width", 4000),
			slog.Int("height", 3000),
			slog.String("format", "jpeg"),
		),
	)
}

// another sample:

func testSlogAndHTTPServer(t testing.TB) { //nolint:unused
	demoWorkingWithLegacyCodes := func() {
		var srv http.Server
		var sigint chan os.Signal

		idleConnsClosed := make(chan struct{})

		logger := slog.New().SetJSONMode()
		srv.ErrorLog = slog.NewLogLogger(logger, slog.ErrorLevel)

		go func() {
			sigint = make(chan os.Signal, 1)
			signal.Notify(sigint, os.Interrupt)
			<-sigint

			// We received an interrupt signal, shut down.
			if err := srv.Shutdown(context.Background()); err != nil {
				// Error from closing listeners, or context timeout:
				t.Logf("HTTP server Shutdown: %v", err)
			}
			close(idleConnsClosed)
		}()

		go func() {
			time.Sleep(time.Second * 3)
			sigint <- os.Interrupt // trigger go routine to shut down the http server after 3 seconds
			t.Log("shutdown signal triggered.")
		}()

		if err := srv.ListenAndServe(); !errors.Is(err, http.ErrServerClosed) {
			// Error starting or closing listener:
			t.Fatalf("HTTP server ListenAndServe: %v", err)
		}
		<-idleConnsClosed
		t.Log("shutdown ok.")
	}

	demoWorkingWithLegacyCodes()

	t.Log("demoWorkingWithLegacyCodes completed.")
}

func TestSlogStrongTypedAttrs(t *testing.T) {
	// slog.SetLevel(slog.InfoLevel)
	for _, logger := range []slog.Logger{
		slog.New(),                     // color mode
		slog.New().SetColorMode(false), // logfmt
		slog.New().SetJSONMode(),       // json mode
	} {
		logger.Println()
		logger.Info(
			"incoming request",
			slog.String("method", "GET"),
			slog.Int("time_taken_ms", 158),
			slog.String("path", "/hello/world?q=search"),
			slog.Int("status", 200),
			slog.String(
				"user_agent",
				"Googlebot/2.1 (+https://www.google.com/bot.html)",
			),
		)
		logger.Info(
			"incoming request",
			"method", "GET",
			slog.Int("time_taken_ms", 158),
			slog.String("path", "/hello/world?q=search"),
			"status", 200,
			slog.String(
				"user_agent",
				"Googlebot/2.1 (+https://www.google.com/bot.html)",
			),
		)
		logger.InfoContext(
			context.Background(),
			"incoming request",
			slog.String("method", "GET"),
			slog.Int("time_taken_ms", 158),
			slog.String("path", "/hello/world?q=search"),
			slog.Int("status", 200),
			slog.String(
				"user_agent",
				"Googlebot/2.1 (+https://www.google.com/bot.html)",
			),
		)
	}
}

func TestSlogGrouping(t *testing.T) {
	for _, logger := range []slog.Logger{
		slog.New(),                     // color mode
		slog.New().SetColorMode(false), // logfmt
		slog.New().SetJSONMode(),       // json mode
	} {
		logger.LogAttrs(
			context.Background(),
			slog.InfoLevel,
			"image uploaded",
			slog.Int("id", 23123),
			slog.Group("properties",
				slog.Int("width", 4000),
				slog.Int("height", 3000),
				slog.String("format", "jpeg"),

				slog.NewAttr("attr3", "any styles what u prefer"),

				slog.Group("more", "group", []byte("more sub-attrs")),

				"attrN", // unpaired key can work here
			),
		)
	}
}

func TestSlogChildLogger(t *testing.T) {
	for _, logger := range []slog.Logger{
		slog.New(),                     // color mode
		slog.New().SetColorMode(false), // logfmt
		slog.New().SetJSONMode(),       // json mode
	} {
		buildInfo, _ := debug.ReadBuildInfo()

		child := logger.New("child").SetAttrs(
			slog.Group("program_info",
				slog.Int("pid", os.Getpid()),
				slog.String("go_version", buildInfo.GoVersion),
			),
		)

		child.Info("image upload successful", slog.String("image_id", "39ud88"))
		child.Warn(
			"storage is 90% full",
			slog.String("available_space", "900.1 mb"),
		)

		child.LogAttrs(
			context.Background(),
			slog.InfoLevel,
			"image uploaded",
			slog.Int("id", 23123),
			slog.Group("properties",
				slog.Int("width", 4000),
				slog.Int("height", 3000),
				slog.String("format", "jpeg"),
			),
		)
	}
}

func TestSlogNewWithAttrs(t *testing.T) {
	logger := slog.New("", "app-version", "v0.0.1-beta")
	ctx := context.WithValue(context.Background(), "ctx", "oh,oh,oh") //nolint:staticcheck
	logger.InfoContext(ctx, "info msg",
		"attr1", 111333,
		slog.Group("memory",
			slog.Int("current", 50),
			slog.Int("min", 20),
			slog.Int("max", 80)),
		slog.Int("cpu", 10),
	)
}

func TestSlogWithContext(t *testing.T) {
	logger := slog.New().SetAttrs(slog.String("app-version", "v0.0.1-beta"))
	ctx := context.WithValue(context.Background(), "ctx", "oh,oh,oh") //nolint:staticcheck
	logger.SetContextKeys("ctx").InfoContext(ctx, "info msg",
		"attr1", 111333,
		slog.Group("memory",
			slog.Int("current", 50),
			slog.Int("min", 20),
			slog.Int("max", 80)),
		slog.Int("cpu", 10),
	)
}

func TestSlogPassLoggerWithContext(t *testing.T) {
	logger := slog.New().SetAttrs(slog.String("app-version", "v0.0.1-beta"))
	ctx := context.WithValue(context.Background(), LoggerKey, logger) //nolint:staticcheck // 👈 context containing logger
	sendUsageStatus(ctx)
}

const LoggerKey = "slogChildLoggerGlobal"

func sendUsageStatus(ctx context.Context) {
	logger := ctx.Value(LoggerKey).(slog.Logger)
	logger.InfoContext(ctx, "info msg", "attr1", 111333,
		slog.Group("memory",
			slog.Int("current", 50),
			slog.Int("min", 20),
			slog.Int("max", 80)),
		slog.Int("cpu", 10),
	)
}

func TestSlogSetLevel(t *testing.T) {
	slog.SetLevel(slog.InfoLevel)
	for _, parent := range []slog.Logger{
		slog.New(),                      // color mode
		slog.New().WithColorMode(false), // logfmt
		slog.New().WithJSONMode(),       // json mode
	} {
		// you can change the level anytime like this
		// parent.SetLevel(slog.TraceLevel)

		parent.Warn("Warn message --------", "lvl", parent.Level()) // this line is invisible in logging outputs

		// or, create a child logger with different level
		logger := parent.New("child").SetLevel(slog.DebugLevel)

		logger.Trace("Trace message/c") // invisible
		logger.Debug("Debug message/c")
		logger.Info("Info message/c")
		logger.Warn("Warning message/c")
		logger.Error("Error message/c")

		logger.Close()
	}
}

const (
	NoticeLevel = slog.Level(17) // A custom level must have a value larger than slog.MaxLevel
	HintLevel   = slog.Level(-8) // Or use a negative number
	SwellLevel  = slog.Level(12) // Sometimes, you may use the value equal with slog.MaxLevel
)

func TestSlogCustomizedLevel(t *testing.T) {
	checkerr(t, slog.RegisterLevel(NoticeLevel, "NOTICE",
		slog.RegWithShortTags([6]string{"", "N", "NT", "NTC", "NOTC", "NOTIC"}),
		slog.RegWithColor(color.FgWhite, color.BgUnderline),
		slog.RegWithTreatedAsLevel(slog.InfoLevel),
	))

	checkerr(t, slog.RegisterLevel(HintLevel, "Hint",
		slog.RegWithShortTags([6]string{"", "H", "HT", "HNT", "HINT", "HINT "}),
		slog.RegWithColor(color.NoColor, color.BgInverse),
		slog.RegWithTreatedAsLevel(slog.InfoLevel),
	))

	checkerr(t, slog.RegisterLevel(SwellLevel, "SWELL",
		slog.RegWithShortTags([6]string{"", "S", "SW", "SWL", "SWEL", "SWEEL"}),
		slog.RegWithColor(color.FgRed, color.BgBoldOrBright),
		slog.RegWithTreatedAsLevel(slog.ErrorLevel),
		slog.RegWithPrintToErrorDevice(),
	))

	logger := slog.New()

	slog.SetLevelColors(slog.DebugLevel, color.FgCyan, color.BgInverse)

	logger.Debug("Debug message")
	logger.Info("Info message")
	logger.Warn("Warning message")
	logger.Error("Error message")

	slog.SetLevelOutputWidth(5)

	ctx := context.Background()
	logger.LogAttrs(ctx, NoticeLevel, "Notice message")
	logger.LogAttrs(ctx, HintLevel, "Hint message")
	logger.LogAttrs(ctx, SwellLevel, "Swell level")
}

func checkerr(t *testing.T, err error) {
	if err != nil {
		t.Error(err)
	}
}

func TestSlogSwitchHandlers(t *testing.T) {
	// No, we don't need any backends so-called handlers.
}

func TestSlogCustomizedHandler(t *testing.T) {
	// No, we don't need any backends so-called handlers.
}

func TestSlogPassValuesToCustomizedHandler(t *testing.T) {
	// No, we don't need any backends so-called handlers.
}

func TestSlogAdapter(t *testing.T) {
	t.Log("slog")

	l := slog.New("standalone-logger-for-app",
		slog.NewAttr("attr1", 2),
		slog.NewAttrs("attr2", 3, "attr3", 4.1),
		"attr4", true, "attr3", "string",
		slog.WithLevel(slog.AlwaysLevel),
	)
	defer l.Close()

	sub1 := l.New("sub1").Set("logger", "sub1")
	sub2 := l.New("sub2").Set("logger", "sub2").SetLevel(slog.InfoLevel)

	// create a log/slog logger HERE
	logger := logslog.New(slog.NewSlogHandler(l, &slog.HandlerOptions{
		NoColor:  false,
		NoSource: true,
		JSON:     false,
		Level:    slog.DebugLevel,
	}))

	t.Logf("logger: %v", logger)
	t.Logf("l: %v", l)

	// and logging with log/slog
	logger.Debug("hi debug", "AA", 1.23456789)
	logger.Info(
		"incoming request",
		logslog.String("method", "GET"),
		logslog.String("path", "/api/user"),
		logslog.Int("status", 200),
	)

	// now using our logg/slog interface
	sub1.Debug("hi debug", "AA", 1.23456789)
	sub2.Debug("hi debug", "AA", 1.23456789)
	sub2.Info("hi info", "AA", 1.23456789)
}

func TestSlogUsedForLogSlog(t *testing.T) {
	l := slog.New("standalone-logger-for-app",
		slog.NewAttr("attr1", 2),
		slog.NewAttrs("attr2", 3, "attr3", 4.1),
		"attr4", true, "attr3", "string",
		slog.WithLevel(slog.AlwaysLevel),
	)
	defer l.Close()

	sub1 := l.New("sub1").Set("logger", "sub1")
	sub2 := l.New("sub2").Set("logger", "sub2").SetLevel(slog.InfoLevel)

	// create a log/slog logger HERE
	logger := logslog.New(slog.NewSlogHandler(l, nil))

	t.Logf("logger: %v", logger)
	t.Logf("l: %v", l)

	// and logging with log/slog
	logger.Debug("hi debug", "AA", 1.23456789)
	logger.Info(
		"incoming request",
		logslog.String("method", "GET"),
		logslog.String("path", "/api/user"),
		logslog.Int("status", 200),
	)

	// now using our logg/slog interface
	sub1.Debug("hi debug", "AA", 1.23456789)
	sub2.Debug("hi debug", "AA", 1.23456789)
	sub2.Info("hi info", "AA", 1.23456789)
}

func Tip(msg string, args ...any) {
	if myConditionSatisfied {
		slog.WithSkip(1).OK(msg, args...)
	}
}

var myConditionSatisfied bool

func TestSlogCustomizedVerbs(t *testing.T) {
	slog.SetLevel(slog.DebugLevel)
	myConditionSatisfied = true
	Tip("hi debug", "AA", 1.23456789)
}

func TestDisabledLog(t *testing.T) {
	// To avoid global level was set to DebugLevel in testing/debugging
	// context, here we force it to WarnLevel
	slog.SetLevel(slog.WarnLevel)

	// And, create a detached logger with WarnLevel
	logger := slog.New().Set("attr", "parent").SetLevel(slog.WarnLevel)

	// Now the Info request should not print out anything
	logger.Info("info msg", "root-level", slog.GetLevel())
}

func TestSlogAttrsR(t *testing.T) {
	logger := slog.New("parent-logger").Set("attr", "parent").SetLevel(slog.InfoLevel)
	sl := logger.New("child-logger")

	logger.Info("info", "attr1", 1)
	slog.AddFlags(slog.LattrsR)
	sl.Info("info", "attr1", 1)
	slog.RemoveFlags(slog.LattrsR)
	sl.Info("info", "attr1", 1)
}

//

//

//

func m1AttrsAsAnySlice() []any {
	return []any{
		slog.String("method", "GET"),
		slog.Int("time_taken_ms", 158),
		slog.String("path", "/hello/world?q=search"),
		slog.Int("status", 200),
		slog.String(
			"user_agent",
			"Googlebot/2.1 (+https://www.google.com/bot.html)",
		),
		slog.Group("memory",
			slog.Int("current", 50),
			slog.Int("min", 20),
			slog.Int("max", 80)),
	}
}

func m2AttrsAsAnySlice() []any {
	return []any{
		slog.String("method", "GET"),
		slog.Int("time_taken_ms", 158),
		slog.String("path", "/hello/world?q=search"),
		slog.Int("status", 200),
		slog.String(
			"user_agent",
			"Googlebot/2.1 (+https://www.google.com/bot.html)",
		),
		slog.Group("memory",
			slog.Group("stacked",
				slog.Float32("Wired", 11),
				slog.Float32("Compressed", 23.21),
				slog.Float32("Cache", 39.79),
			),
			slog.Int("current", 50),
			slog.Int("min", 20),
			slog.Int("max", 80)),
	}
}
