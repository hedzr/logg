package slog

import (
	"context"
	"io"
	stdlog "log/slog"
	"time"
)

type (
	// Logger interface can be used for others logging kits or apps.
	Logger interface {
		EntryI
		BuilderI

		// Log(ctx context.Context, level Level, msg string, args ...any)

		// Log to log/slog
		Log(ctx context.Context, level stdlog.Level, msg string, args ...any)

		ExtraPrintersI
	}

	// ExtraPrintersI to make compatible with classical log.
	//
	// These are not efficient apis, we just have a simple
	// implementations concretely.
	ExtraPrintersI interface {
		Infof(format string, a ...interface{}) error
		Warnf(format string, a ...interface{}) error
		Errorf(format string, a ...interface{}) error
	}

	// EntryI is a small and efficient tiny logger, which is the entity of real logger.
	EntryI interface {
		BasicLogger

		Close() // Closeable interface

		String() string // Stringer interface

		Parent() *Entry // parent logger of a sub-logger
		Root() *Entry   // root logger (always is Default) of a sub-logger

		// WithSkip create a new child logger with specified extra
		// ignored stack frames, which will be plussed over the
		// internal stack frames stripping tool.
		//
		// A child logger is super lite commonly. It'll take a little
		// more resource usages only if you have LattrsR set globally.
		// In that case, child logger looks up all its parents for
		// collecting all attributes and logging them.
		WithSkip(extraFrames int) *Entry

		// Children() Entries

		Level() Level // logging level associated with this logger

		// writeInternal(ctx context.Context, lvl Level, pc uintptr, buf []byte) (n int, err error)
		// logContext(ctx context.Context, lvl Level, pc uintptr, msg string, args ...any)
	}

	// BuilderI is used for building a new logger
	//
	// // WithXXX apis: make a child logger and apply the new settings.
	// _
	// // SetXXX apis: apply the new settings on this logger
	// _
	BuilderI interface {
		New(args ...any) *Entry // 1st of args is name, the rest are k, v pairs

		// WithJSONMode and WithColorMode sets output format to json or logfmt(+color).
		WithJSONMode(b ...bool) *Entry          // entering JSON mode, the output are json format
		WithColorMode(b ...bool) *Entry         // entering Colorful mode for the modern terminal. false means using logfmt format.
		WithUTCMode(b ...bool) *Entry           // default is local mode, true will switch to UTC mode
		WithTimeFormat(layout ...string) *Entry // specify your timestamp format layout string
		WithLevel(lvl Level) *Entry             //
		WithAttrs(attrs ...Attr) *Entry         //
		WithAttrs1(attrs Attrs) *Entry          //
		With(args ...any) *Entry                // key1,val1,key2,val2,.... Of course, Attr, Attrs in args will be recognized as is

		SetJSONMode(b ...bool) *Entry          // entering JSON mode, the output are json format
		SetColorMode(b ...bool) *Entry         // entering Colorful mode for the modern terminal. false means using logfmt format.
		SetUTCMode(b ...bool) *Entry           // default is local mode, true will switch to UTC mode
		SetTimeFormat(layout ...string) *Entry // specify your timestamp format layout string
		SetLevel(lvl Level) *Entry             //
		SetAttrs(attrs ...Attr) *Entry         //
		SetAttrs1(attrs Attrs) *Entry          //
		Set(args ...any) *Entry                // key1,val1,key2,val2,.... Of course, Attr, Attrs in args will be recognized as is

		SetContextKeys(keys ...any) *Entry  // given keys will be tried extracting from context.Context automatically
		WithContextKeys(keys ...any) *Entry // given keys will be tried extracting from context.Context automatically

		SetWriter(wr io.Writer) *Entry    // use the given writer
		WithWriter(wr io.Writer) *Entry   // use the given writer
		AddWriter(wr io.Writer) *Entry    // append more writers via this interface
		RemoveWriter(wr io.Writer) *Entry // remove a writer

		SetErrorWriter(wr io.Writer) *Entry    //
		WithErrorWriter(wr io.Writer) *Entry   //
		AddErrorWriter(wr io.Writer) *Entry    //
		RemoveErrorWriter(wr io.Writer) *Entry //

		ResetWriters() *Entry // reset std and error writers

		GetWriter() (wr LogWriter)              // return level-matched writer
		GetWriterBy(level Level) (wr LogWriter) // return writer matched given level

		AddLevelWriter(lvl Level, w io.Writer) *Entry    //
		RemoveLevelWriter(lvl Level, w io.Writer) *Entry //
		ResetLevelWriter(lvl Level) *Entry               // reset the writers in a level
		ResetLevelWriters() *Entry                       // reset all leveled writers

		SetValueStringer(vs ValueStringer) *Entry  //
		WithValueStringer(vs ValueStringer) *Entry //
	}

	// Entries collects many Entry objects as a map
	Entries map[string]*Entry

	// BasicLogger supplies basic logging apis.
	BasicLogger interface {
		Printer
		PrinterWithContext

		Enabled(requestingLevel Level) bool // to test the requesting logging level should be allowed.
		EnabledContext(ctx context.Context, requestingLevel Level) bool

		LogAttrs(ctx context.Context, level Level, msg string, args ...any) // Attr, Attrs in args will be recognized as is
		Logit(ctx context.Context, level Level, msg string, args ...any)    // Attr, Attrs in args will be recognized as is

		// SetSkip is very similar with WithSkip but no child logger
		// created, it modifies THIS logger.
		//
		// Use it when you know all what u want.
		SetSkip(extraFrames int)
		Skip() int // return current frames count should be ignored in addition. 0 for most cases.

		Name() string    // this logger's name
		JSONMode() bool  // return the mode
		ColorMode() bool // return the mode
	}

	// LogLoggerAware for external adapters
	LogLoggerAware interface {
		WriteInternal(ctx context.Context, lvl Level, pc uintptr, buf []byte) (n int, err error)
	}

	// LogSlogAware for external adapters
	LogSlogAware interface {
		WriteThru(ctx context.Context, lvl Level, timestamp time.Time, pc uintptr, msg string, attrs Attrs)
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
