package slog

import (
	"context"
	"io"
	"time"
)

type (
	// Logger interface can be used for others logging kits or apps.
	Logger interface {
		Entry
	}

	// Entry is a small and efficient tiny logger, which is the entity of real logger.
	Entry interface {
		BasicLogger

		Close() // Closeable interface

		String() string // Stringer interface

		Parent() Entry // parent logger of a sub-logger
		Root() Entry   // root logger (always is Default) of a sub-logger

		// Children() Entries

		Level() Level // logging level associated with this logger

		// writeInternal(ctx context.Context, lvl Level, pc uintptr, buf []byte) (n int, err error)
		// logContext(ctx context.Context, lvl Level, pc uintptr, msg string, args ...any)
	}

	// LogLoggerAware for external adapters
	LogLoggerAware interface {
		WriteInternal(ctx context.Context, lvl Level, pc uintptr, buf []byte) (n int, err error)
	}

	// LogSlogAware for external adapters
	LogSlogAware interface {
		WriteThru(ctx context.Context, lvl Level, timestamp time.Time, pc uintptr, msg string, attrs Attrs)
	}

	// BasicLogger supplies basic logging apis.
	BasicLogger interface {
		Printer
		PrinterWithContext
		Builder

		Enabled(requestingLevel Level) bool // to test the requesting logging level should be allowed.
		EnabledContext(ctx context.Context, requestingLevel Level) bool

		LogAttrs(ctx context.Context, level Level, msg string, args ...any) // Attr, Attrs in args will be recognized as is

		// WithSkip create a new child logger with specified extra
		// ignored stack frames, which will be plussed over the
		// internal stack frames stripping tool.
		//
		// A child logger is super lite commonly. It'll take a little
		// more resource usages only if you have LattrsR set globally.
		// In that case, child logger looks up all its parents for
		// collecting all attributes and logging them.
		WithSkip(extraFrames int) Entry

		// SetSkip is very similar with WithSkip but no child logger
		// created, it modifies THIS logger.
		//
		// Use it when you know all what u want.
		SetSkip(extraFrames int)
		Skip() int // return current frames count should be ignored in addition. 0 for most cases.

		Name() string // this logger's name
	}

	// Printer supplies the printable apis
	Printer interface {
		Panic(msg string, args ...any)   // error and panic
		Fatal(msg string, args ...any)   // error and os.Exit(-3)
		Error(msg string, args ...any)   // error
		Warn(msg string, args ...any)    // warning
		Info(msg string, args ...any)    // info. Attr, Attrs in args will be recognized as is
		Debug(msg string, args ...any)   // only for state.Env().InDebugging() or IsDebugBuild()
		Trace(msg string, args ...any)   // only for state.Env().InTracing()
		Verbose(msg string, args ...any) // only for -tags=verbose
		Print(msg string, args ...any)   // logging always
		Println(args ...any)             // synonym to Print, NOTE first elem of args decoded as msg here
		OK(msg string, args ...any)      // identify it is in OK mode
		Success(msg string, args ...any) // identify a successful operation done
		Fail(msg string, args ...any)    // identify a wrong occurs, default to stderr device
	}

	// PrinterWithContext supplies the printable apis with context.Context
	PrinterWithContext interface {
		PanicContext(ctx context.Context, msg string, args ...any)   // error and panic
		FatalContext(ctx context.Context, msg string, args ...any)   // error and os.Exit(-3)
		ErrorContext(ctx context.Context, msg string, args ...any)   // error
		WarnContext(ctx context.Context, msg string, args ...any)    // warning
		InfoContext(ctx context.Context, msg string, args ...any)    // info. Attr, Attrs in args will be recognized as is
		DebugContext(ctx context.Context, msg string, args ...any)   // only for state.Env().InDebugging() or IsDebugBuild()
		TraceContext(ctx context.Context, msg string, args ...any)   // only for state.Env().InTracing()
		VerboseContext(ctx context.Context, msg string, args ...any) // only for -tags=verbose
		PrintContext(ctx context.Context, msg string, args ...any)   // logging always
		PrintlnContext(ctx context.Context, msg string, args ...any) // synonym to Print
		OKContext(ctx context.Context, msg string, args ...any)      // identify it is in OK mode
		SuccessContext(ctx context.Context, msg string, args ...any) // identify a successful operation done
		FailContext(ctx context.Context, msg string, args ...any)    // identify a wrong occurs, default to stderr device
	}

	// Builder is used for building a new logger
	Builder interface {
		New(args ...any) BasicLogger // 1st of args is name, the rest are k, v pairs

		WithJSONMode(b ...bool) Entry          // entering JSON mode, the output are json format
		WithColorMode(b ...bool) Entry         // entering Colorful mode for the modern terminal. false means using logfmt format.
		WithUTCMode(b ...bool) Entry           // default is local mode, true will switch to UTC mode
		WithTimeFormat(layout ...string) Entry // specify your timestamp format layout string
		WithLevel(lvl Level) Entry             //
		WithAttrs(attrs ...Attr) Entry         //
		WithAttrs1(attrs Attrs) Entry          //
		With(args ...any) Entry                // key1,val1,key2,val2,.... Of course, Attr, Attrs in args will be recognized as is

		WithContextKeys(keys ...any) Entry // given keys will be tried extracting from context.Context automatically

		WithWriter(wr io.Writer) Entry          // use the given writer
		AddWriter(wr io.Writer) Entry           // append more writers via this interface
		AddErrorWriter(wr io.Writer) Entry      //
		ResetWriters() Entry                    //
		GetWriter() (wr LogWriter)              // return level-matched writer
		GetWriterBy(level Level) (wr LogWriter) // return writer matched given level

		AddLevelWriter(lvl Level, w io.Writer) Entry
		RemoveLevelWriter(lvl Level, w io.Writer) Entry
		ResetLevelWriter(lvl Level) Entry
		ResetLevelWriters() Entry

		WithValueStringer(vs ValueStringer) Entry
	}

	// Entries collects many entry objects as a map
	Entries map[string]Entry
)

// LogWriter for external adapters
type LogWriter interface {
	io.Writer
	io.Closer
}

// LogWriters for external adapters
type LogWriters interface{} //nolint:revive

type (
	// Attr for external adapters
	Attr interface {
		Key() string
		Value() any
		SetValue(v any) // to modify the value dynamically
	}
)

type (
	// LogValuer for external adapters
	LogValuer interface {
		Value() Attr
	}
)

const maxLogValues = 100 //nolint:unused

// func errorlog(err error) {
// 	// _, _ = fmt.Fprintf(lw.GetErrorOutput(), "[logg/slog] error occurs: %+v", err)
// 	if w := defaultWriter.Get(ErrorLevel); w != nil {
// 		_, _ = w.Write([]byte("[logg/slog] error occurs: "))
// 		_, _ = w.Write([]byte(err.Error()))
// 	}
// }
//
// func raiseerror(msg string) { panic(msg) }
