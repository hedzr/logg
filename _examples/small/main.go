package main

import (
	"context"
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
		WithOnLoop(dbStarter, cacheStarter, mqStarter).
		WithOnSignalCaught(func(sig os.Signal, wg *sync.WaitGroup) {
			println()
			slog.Info("signal caught", "sig", sig)
			cancel() // cancel user's loop, see Wait(...)
		}).
		Wait(func(stopChan chan<- os.Signal, wgDone *sync.WaitGroup) {
			slog.Debug("entering looper's loop...")
			go func() {
				// to terminate this app after a while automatically:
				time.Sleep(10 * time.Second)
				stopChan <- os.Interrupt
			}()
			<-ctx.Done()  // waiting until any os signal caught
			wgDone.Done() // and complete myself
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

func dbStarter(stopChan chan<- os.Signal, wgDone *sync.WaitGroup) {
	// initializing database connections...
	// ...
	wgDone.Done()
}

func cacheStarter(stopChan chan<- os.Signal, wgDone *sync.WaitGroup) {
	// initializing redis cache connections...
	// ...
	wgDone.Done()
}

func mqStarter(stopChan chan<- os.Signal, wgDone *sync.WaitGroup) {
	// initializing message queue connections...
	// ...
	wgDone.Done()
}
