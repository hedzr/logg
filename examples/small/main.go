package main

import (
	"context"
	"errors"
	"fmt"
	"log/slog"
	"os"
	"sync"
	"time"

	"github.com/hedzr/is"
	"github.com/hedzr/is/basics"
	"github.com/hedzr/is/term"
	"github.com/hedzr/is/term/color"

	logz "github.com/hedzr/logg/slog"
)

func main() {
	defer basics.Close()

	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	testIs(ctx)
	testLogz(ctx)

	catcher := is.Signals().Catch()
	catcher.
		WithPrompt("Press CTRL-C to quit...").
		WithOnLoopFunc(dbStarter, cacheStarter, mqStarter).
		WithOnSignalCaught(func(sig os.Signal, wg *sync.WaitGroup) {
			println()
			slog.Info("signal caught", "sig", sig)
			cancel() // cancel user's loop, see Wait(...)
		}).
		WaitFor(func(closer func()) {
			slog.Debug("entering looper's loop...")
			go func() {
				// to terminate this app after a while automatically:
				time.Sleep(10 * time.Second)
				// stopChan <- os.Interrupt
				closer()
			}()
			<-ctx.Done() // waiting until any os signal caught
			// wgDone.Done() // and complete myself
		})
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

func testLogz(ctx context.Context) {
	logz.SetLevel(logz.DebugLevel)

	logz.InfoContext(ctx, "Hello, world!")

	logz.InfoContext(ctx, "Hello", "target", "world")
	logz.Default().SetColorMode(false)
	logz.InfoContext(ctx, "Hello", "target", "world")
	logz.Default().SetJSONMode(true)
	logz.InfoContext(ctx, "Hello", "target", "world")

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

func dbStarter(closer func()) {
	defer closer()
	// initializing database connections...
	// ...
}

func cacheStarter(closer func()) {
	defer closer()
	// initializing redis cache connections...
	// ...
}

func mqStarter(closer func()) {
	defer closer()
	// initializing message queue connections...
	// ...
}

const (
	AppNameExample = "small"     // appName for the current demo app
	appName        = "logg/slog" // appName of hedzr/logg package
	version        = "v0.8.1"    // version of hedzr/logg package | update it while bumping hedzr/logg' version
	Version        = version     // version name you can check it
)
