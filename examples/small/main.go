package main

import (
	"context"
	"errors"
	"fmt"
	"log/slog"
	"os"
	"time"

	"github.com/hedzr/is"
	"github.com/hedzr/is/term"
	"github.com/hedzr/is/term/color"
	"github.com/hedzr/is/timing"

	logz "github.com/hedzr/logg/slog"
)

func main() {
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	testIs(ctx)
	testIs2(ctx)
	testLogz(ctx)

	p := timing.New()
	defer p.CalcNow()

	// go func() {
	// 	time.Sleep(3 * time.Second)
	// 	cancel() // stop after 4s instead of waiting for 6s later.
	// }()

	is.SignalsEnh().WaitForSeconds(ctx, cancel, 6*time.Second,
		// is.WithCatcherCloser(cancel),
		is.WithCatcherMsg("Press CTRL-C to quit, or waiting for 6s..."),
	)
}

func testIs(ctx context.Context) {
	is.RegisterStateGetter("custom", func() bool { return is.InVscodeTerminal() })

	println("state.InTesting:   ", is.InTesting())
	println("state.in-testing:  ", is.State("in-testing"))
	println("state.custom:      ", is.State("custom")) // detects a state with custom detector
	println("env.GetDebugLevel: ", is.Env().GetDebugLevel())
	if is.InDebugMode() {
		slog.SetDefault(slog.New(slog.NewTextHandler(os.Stdout, &slog.HandlerOptions{AddSource: true, Level: slog.LevelDebug})))
	}

	fmt.Printf("\n%v", color.GetCPT().Translate(`<code>code</code> | <kbd>CTRL</kbd>
		<b>bold / strong / em</b>
		<i>italic / cite</i>
		<u>underline</u>
		<mark>inverse mark</mark>
		<del>strike / del </del>
		<font color="green">green text</font>
`, color.FgDefault))

	println("term.IsTerminal:               ", term.IsTerminal(int(os.Stdout.Fd())))
	println("term.IsAnsiEscaped:            ", term.IsAnsiEscaped(color.GetCPT().Translate(`<code>code</code>`, color.FgDefault)))
	println("term.IsCharDevice(stdout):     ", term.IsCharDevice(os.Stdout))
	rows, cols, err := term.GetFdSize(os.Stdout.Fd())
	println("term.GetFdSize(stdout):        ", rows, cols, err)
	rows, cols, err = term.GetTtySizeByFd(os.Stdout.Fd())
	println("term.GetTtySizeByFd(stdout):   ", rows, cols, err)
	rows, cols, err = term.GetTtySizeByFile(os.Stdout)
	println("term.GetTtySizeByFile(stdout): ", rows, cols, err)
	println("term.IsStartupByDoubleClick:   ", term.IsStartupByDoubleClick())

	logz.InfoContext(ctx, "pre-forecasting") // by default this line shall not be displayed since logz (logg/slog) is in WarnLevel.
}

func testIs2(ctx context.Context) {
	// start a color text builder
	var c = color.New()
	var pos color.CursorPos

	// paint and get the result (with ansi-color-seq ready)
	var result = c.Println().
		Color16(color.FgRed).Printf("[1st] hello, %s.", "world").
		Println().
		SavePosNow().
		Println("XX").
		Color16(color.FgGreen).Printf("hello, %s.\n", "world").
		Color256(160).Printf("[160] hello, %s.\n", "world").
		Color256(161).Printf("[161] hello, %s.\n", "world").
		Color256(162).Printf("[162] hello, %s.\n", "world").
		Color256(163).Printf("[163] hello, %s.\n", "world").
		Color256(164).Printf("[164] hello, %s.\n", "world").
		Color256(165).Printf("[165] hello, %s.\n", "world").
		UpNow(4).Echo(" ERASED ").
		RightNow(11).
		CursorGet(ctx, &pos).
		RGB(211, 211, 33).Printf("[16m] hello, %s. pos=%+v", "world", pos).
		Println().
		RestorePosNow().
		Println("ZZ").
		DownNow(8).
		Println("DONE").
		Build()

		// and render the result
	fmt.Println(result)

	// another colorful builfer
	c = color.New()
	fmt.Println(c.Color16(color.FgRed).
		Printf("[2nd] hello, %s.", "world").Println().Build())

	// cursor operations
	c = color.New()
	c.SavePosNow()
	// fmt.Println(c.CursorSavePos().Build())

	fmt.Print(c.
		Printf("[3rd] hello, %s.", "world").
		Println().
		Color256(163).Printf("[163] hello, %s.\n", "world").
		Color256(164).Printf("[164] hello, %s.\n", "world").
		Color256(165).Printf("[165] hello, %s.\n", "world").
		Build())

	fmt.Print("0")         // now, col = 1
	c.UpNow(2)             //
	fmt.Print("ABC")       // embedded "ABC" into "[]"
	c.CursorGet(ctx, &pos) //
	c.RightNow(2)          // to be overwrite "hello"
	fmt.Print("HELLO")     //

	c.RestorePosNow()
	c.DownNow(1)
	fmt.Print("T") // write "T" to beginning of "[163]" line

	c.DownNow(4)

	// color.Down(4)
	// color.Left(1)
	fmt.Printf("\nEND (pos = %+v)\n", pos)
}

func testLogz(ctx context.Context) {
	logz.SetLevel(logz.TraceLevel)

	logz.Default().SetMode(logz.ModeColorful)

	logz.InfoContext(ctx, "Hello, world!")

	println("Colorful")
	logz.InfoContext(ctx, "Hello", "target", "world")
	println("Plain")
	logz.Default().SetMode(logz.ModePlain)
	logz.InfoContext(ctx, "Hello", "target", "world")
	println("Logfmt")
	// logz.Default().SetColorMode(false)
	logz.Default().SetMode(logz.ModeLogFmt)
	logz.InfoContext(ctx, "Hello", "target", "world")
	println("JSON")
	// logz.Default().SetJSONMode(true)
	logz.Default().SetMode(logz.ModeJSON)
	logz.InfoContext(ctx, "Hello", "target", "world")

	logz.Default().SetMode(logz.ModeColorful)

	testLogzAdapter(ctx)

	testLogz1(ctx)
}

func testLogz1(ctx context.Context) {
	msg := "A message"
	args := []any{
		"attr1", 0,
		"attr2", false,
		"attr3", 3.13,
		"attr4", errors.New("simple"),
		// ,,,
		logz.Group("group1",
			"attr1", 0,
			"attr2", false,
			logz.NewAttr("attr3", "any styles what u prefer"),
			"attrn", // unpaired key can work here
			logz.Group("more", "group", []byte("more subgroup here")),
		),
		// ...
		logz.Int("id", 23123),
		logz.Group("properties",
			logz.Int("width", 4000),
			logz.Int("height", 3000),
			logz.String("format", "jpeg"),
		),
	}

	// disable unconditional termination inside logz.Panic/Fatal() calls.
	logz.AddFlags(logz.LnoInterrupt)

	logz.Print("")   // logging a clean newline without decorations
	logz.Println("") // logging a clean newline without decorations
	logz.Println()   // logging a clean newline without decorations
	logz.Print(msg, args...)
	logz.Println(msg) // synosym of Print
	logz.Fatal(msg, args...)
	logz.Panic(msg, args...)
	logz.Error(msg, args...)
	logz.Warn(msg, args...)
	logz.Info(msg, args...)
	logz.Debug(msg, args...)
	logz.Trace(msg, args...)

	// only print the logging contents while built with `-tags verbose`
	logz.Verbose(msg, args...) //

	// some verbs with more meanings
	logz.OK("ok")
	logz.OK("ok", args...)
	logz.Success(msg, args...)
	logz.Fail(msg, args...)

	// Contextual logging
	logz.InfoContext(ctx, "info msg", args...)

	logName := "child1"
	log := logz.New(logName)
	defer log.Close() // when you added file writer into `log`

	log.Print("")   // logging a clean newline without decorations
	log.Println("") // logging a clean newline without decorations
	log.Println()   // logging a clean newline without decorations
	log.Print(msg, args...)
	log.Println(msg) // synosym of Print
	log.Fatal(msg, args...)
	log.Panic(msg, args...)
	log.Error(msg, args...)
	log.Warn(msg, args...)
	log.Info(msg, args...)
	log.Debug(msg, args...)
	log.Trace(msg, args...)

	log.LogAttrs(ctx, logz.DebugLevel, "debug msg", args...)
}

func testLogzAdapter(ctx context.Context) {
	l := logz.New("standalone-logger-for-app",
		logz.NewAttr("attr1", 2),
		logz.NewAttrs("attr2", 3, "attr3", 4.1),
		"attr4", true, "attr3", "string",
		logz.WithLevel(logz.AlwaysLevel),
	)
	defer l.Close()

	sub1 := l.New("sub1").With("logger", "sub1")
	sub2 := l.New("sub2").With("logger", "sub2").
		WithLevel(logz.InfoLevel) // a new child instance here

	// create a log/slog logger HERE
	logger := slog.New(logz.NewSlogHandler(l, &logz.HandlerOptions{
		NoColor:  false,
		NoSource: true,
		JSON:     false,
		Level:    logz.DebugLevel,
	}))

	l.Infof("logger: %v", logger)
	l.Infof("l: %v", l)

	// and logging with log/slog
	logger.DebugContext(ctx, "hi debug", "AA", 1.23456789)
	logger.InfoContext(ctx,
		"incoming request",
		slog.String("method", "GET"),
		slog.String("path", "/api/user"),
		slog.Int("status", 200),
	)

	// now using our logg/slog interface
	sub1.DebugContext(ctx, "hi debug", "AA", 1.23456789)
	sub2.DebugContext(ctx, "hi debug", "AA", 1.23456789)
	sub2.InfoContext(ctx, "hi info", "AA", 1.23456789)
}

const (
	AppNameExample = "small"     // appName for the current demo app
	appName        = "logg/slog" // appName of hedzr/logg package
	version        = "v0.8.1"    // version of hedzr/logg package | update it while bumping hedzr/logg' version
	Version        = version     // version name you can check it
)
