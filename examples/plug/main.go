package main

import (
	"context"
	"errors"

	"github.com/hedzr/logg/examples/plug/sub"
	logz "github.com/hedzr/logg/slog"
)

func main() {
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	logz.SetLevel(logz.TraceLevel)

	forLogz(ctx)

	testLogz1(ctx)

	testSubLoggers(ctx)
}

func testSubLoggers(ctx context.Context) {
	cool := sub.RawLogger().New("cool")

	sl := cool.New("cool.sub1")
	sl2 := cool.New("cool.sub2")
	ssl := sl2.New("cool.sub2.sub")

	cool.InfoContext(ctx, "Hello, world!")
	sl.InfoContext(ctx, "Hello, world!")
	sl2.InfoContext(ctx, "Hello, world!")
	ssl.InfoContext(ctx, "Hello, world!")

	ssl.Parent().InfoContext(ctx, "Hello, world!")
	ssl.Parent().Parent().InfoContext(ctx, "Hello, world!")
	atop := ssl.Parent().Parent().Parent()
	atop.InfoContext(ctx, "Hello, world!")
	atop.Parent().InfoContext(ctx, "Hello, world!")
	logz.Default().InfoContext(ctx, "Hello, world!")

	println(cool.Root().DumpSubloggers())
	println(logz.Default().Sublogger("").DumpSubloggers())
}

func forLogz(ctx context.Context) {
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
