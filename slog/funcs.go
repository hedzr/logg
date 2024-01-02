package slog

import (
	"context"
	"errors"
	"io"
	"log"
	"time"

	"github.com/hedzr/is"
)

// IsTty detects a writer if it is abstracting from a tty (console, terminal) device.
func IsTty(w io.Writer) bool { return is.Tty(w) }

// IsColoredTty detects a writer if it is a colorful tty device.
//
// A colorful tty device can receive ANSI escaped sequences and draw its.
func IsColoredTty(w io.Writer) bool { return is.ColoredTty(w) }

// IsTtyEscaped detects a string if it contains ansi color escaped sequences
func IsTtyEscaped(s string) bool { return is.AnsiEscaped(s) }

// IsAnsiEscaped detects a string if it contains ansi color escaped sequences
func IsAnsiEscaped(s string) bool { return is.AnsiEscaped(s) }

// StripEscapes removes any ansi color escaped sequences from a string
func StripEscapes(str string) (strCleaned string) { return is.StripEscapes(str) }

// ReadPassword reads the password from stdin with safe protection
func ReadPassword() (text string, err error) { return is.ReadPassword() }

// GetTtySize returns the window size in columns and rows in the active console window.
// The return value of this function is in the order of cols, rows.
func GetTtySize() (cols, rows int) { return is.GetTtySize() }

//

//

//

func Panic(msg string, args ...any)   { logctx(PanicLevel, msg, args...) }   // Panic with Default Logger.
func Fatal(msg string, args ...any)   { logctx(FatalLevel, msg, args...) }   // Fatal with Default Logger.
func Error(msg string, args ...any)   { logctx(ErrorLevel, msg, args...) }   // Error with Default Logger.
func Warn(msg string, args ...any)    { logctx(WarnLevel, msg, args...) }    // Warn with Default Logger.
func Info(msg string, args ...any)    { logctx(InfoLevel, msg, args...) }    // Info with Default Logger.
func Debug(msg string, args ...any)   { logctx(DebugLevel, msg, args...) }   // Debug with Default Logger.
func Trace(msg string, args ...any)   { logctx(TraceLevel, msg, args...) }   // Trace with Default Logger.
func Print(msg string, args ...any)   { logctx(AlwaysLevel, msg, args...) }  // Print with Default Logger.
func OK(msg string, args ...any)      { logctx(OKLevel, msg, args...) }      // OK with Default Logger.
func Success(msg string, args ...any) { logctx(SuccessLevel, msg, args...) } // Success with Default Logger.
func Fail(msg string, args ...any)    { logctx(FailLevel, msg, args...) }    // Fail with Default Logger.

func Verbose(msg string, args ...any)                             { vlogctx(context.TODO(), msg, args...) } // Verbose with Default Logger.
func VerboseContext(ctx context.Context, msg string, args ...any) { vlogctx(ctx, msg, args...) }            // Verbose with Default Logger.

// Println with Default Logger.
func Println(args ...any) {
	if len(args) == 0 {
		logctx(AlwaysLevel, "")
		return
	}
	var msg string
	msg, args = args[0].(string), args[1:]
	logctx(AlwaysLevel, msg, args...)
}

func logctx(lvl Level, msg string, args ...any) {
	ctx := context.Background()
	switch s := defaultLog.(type) {
	case *logimp:
		if s.EnabledContext(ctx, lvl) {
			pc := getpc(3, s.extraFrames) // caller -> slog.Info -> logctx (this func)
			s.logContext(ctx, lvl, pc, msg, args...)
		}
	case *entry:
		if s.EnabledContext(ctx, lvl) {
			pc := getpc(3, s.extraFrames) // caller -> slog.Info -> logctx (this func)
			s.logContext(ctx, lvl, pc, msg, args...)
		}
	}
}

func logctxctx(ctx context.Context, lvl Level, msg string, args ...any) {
	switch s := defaultLog.(type) {
	case *logimp:
		if s.EnabledContext(ctx, lvl) {
			pc := getpc(3, s.extraFrames) // caller -> slog.Info -> logctx (this func)
			s.logContext(ctx, lvl, pc, msg, args...)
		}
	case *entry:
		if s.EnabledContext(ctx, lvl) {
			pc := getpc(3, s.extraFrames) // caller -> slog.Info -> logctx (this func)
			s.logContext(ctx, lvl, pc, msg, args...)
		}
	}
}

//

// PanicContext with Default Logger.
func PanicContext(ctx context.Context, msg string, args ...any) {
	logctxctx(ctx, PanicLevel, msg, args...)
}

// FatalContext with Default Logger.
func FatalContext(ctx context.Context, msg string, args ...any) {
	logctxctx(ctx, FatalLevel, msg, args...)
}

// ErrorContext with Default Logger.
func ErrorContext(ctx context.Context, msg string, args ...any) {
	logctxctx(ctx, ErrorLevel, msg, args...)
}

// WarnContext with Default Logger.
func WarnContext(ctx context.Context, msg string, args ...any) {
	logctxctx(ctx, WarnLevel, msg, args...)
}

// InfoContext with Default Logger.
func InfoContext(ctx context.Context, msg string, args ...any) {
	logctxctx(ctx, InfoLevel, msg, args...)
}

// DebugContext with Default Logger.
func DebugContext(ctx context.Context, msg string, args ...any) {
	logctxctx(ctx, DebugLevel, msg, args...)
}

// TraceContext with Default Logger.
func TraceContext(ctx context.Context, msg string, args ...any) {
	logctxctx(ctx, TraceLevel, msg, args...)
}

// PrintContext with Default Logger.
func PrintContext(ctx context.Context, msg string, args ...any) {
	logctxctx(ctx, AlwaysLevel, msg, args...)
}

// PrintlnContext with Default Logger.
func PrintlnContext(ctx context.Context, msg string, args ...any) {
	logctxctx(ctx, AlwaysLevel, msg, args...)
}

// OKContext with Default Logger.
func OKContext(ctx context.Context, msg string, args ...any) {
	logctxctx(ctx, OKLevel, msg, args...)
}

// SuccessContext with Default Logger.
func SuccessContext(ctx context.Context, msg string, args ...any) {
	logctxctx(ctx, SuccessLevel, msg, args...)
}

// FailContext with Default Logger.
func FailContext(ctx context.Context, msg string, args ...any) {
	logctxctx(ctx, FailLevel, msg, args...)
}

//

//

//

// String constructs a key-value pair with string value like log/slog.
//
// The different is we have not optimized these functions (String, Int, ...) for
// performance and memory allocations. So they are just compatible with log/slog.
//
// For performance, using With(attrs...) / WithAttrs(...) to get prefer effects.
func String(key string, val string) Attr          { return &kvp{key, val} }
func Bool(key string, val bool) Attr              { return &kvp{key, val} } // constructs boolean k-v pair. see String for performance tip.
func Int(key string, val int) Attr                { return &kvp{key, val} } // constructs Int k-v pair. see String for performance tip.
func Int8(key string, val int8) Attr              { return &kvp{key, val} } // constructs Int8 k-v pair. see String for performance tip.
func Int16(key string, val int16) Attr            { return &kvp{key, val} } // constructs Int16 k-v pair. see String for performance tip.
func Int32(key string, val int32) Attr            { return &kvp{key, val} } // constructs Int32 k-v pair. see String for performance tip.
func Int64(key string, val int64) Attr            { return &kvp{key, val} } // constructs Int64 k-v pair. see String for performance tip.
func Uint(key string, val uint) Attr              { return &kvp{key, val} } // constructs Uint k-v pair. see String for performance tip.
func Uint8(key string, val uint8) Attr            { return &kvp{key, val} } // constructs Uint8 k-v pair. see String for performance tip.
func Uint16(key string, val uint16) Attr          { return &kvp{key, val} } // constructs Uint16 k-v pair. see String for performance tip.
func Uint32(key string, val uint32) Attr          { return &kvp{key, val} } // constructs Uint32 k-v pair. see String for performance tip.
func Uint64(key string, val uint64) Attr          { return &kvp{key, val} } // constructs Uint64 k-v pair. see String for performance tip.
func Float32(key string, val float32) Attr        { return &kvp{key, val} } // constructs Float32 k-v pair. see String for performance tip.
func Float64(key string, val float64) Attr        { return &kvp{key, val} } // constructs Float64 k-v pair. see String for performance tip.
func Complex64(key string, val complex64) Attr    { return &kvp{key, val} } // constructs Complex64 k-v pair. see String for performance tip.
func Complex128(key string, val complex128) Attr  { return &kvp{key, val} } // constructs Complex128 k-v pair. see String for performance tip.
func Time(key string, val time.Time) Attr         { return &kvp{key, val} } // constructs Time k-v pair. see String for performance tip.
func Duration(key string, val time.Duration) Attr { return &kvp{key, val} } // constructs Duration k-v pair. see String for performance tip.
func Any(key string, val any) Attr                { return &kvp{key, val} } // constructs Any k-v pair. see String for performance tip.

func Numeric[T Numerics](key string, val T) Attr { return &kvp{key, val} } // constructs Numeric k-v pair. see String for performance tip.

// Group constructs grouped k-v pair container, which can hold a set of normal attrs.
//
// See String for performance tip.
//
// For example:
//
//	g := Group("source",
//	   String("file", filename),
//	   Int("line", lineno),
//	)
func Group(key string, args ...any) Attr {
	var g = &gkvp{key: key, items: argsToAttrs(nil, args...)}
	return g
}

func buildAttrs(as ...any) (kvps Attrs)                             { return argsToAttrs(nil, as...) }
func buildUniqueAttrs(keys map[string]bool, as ...any) (kvps Attrs) { return argsToAttrs(keys, as...) }

func argsToAttrs(keysKnown map[string]bool, args ...any) (kvps Attrs) {
	var key string
	if keysKnown == nil {
		// keysKnown = make(map[string]bool)
		for _, it := range args {
			if key == "" {
				switch k := it.(type) {
				case string:
					key = k
				case Attr:
					kvps = append(kvps, k)
					key = ""
				case []Attr:
					for _, el := range k {
						kvps = append(kvps, el)
					}
					key = ""
				case Attrs:
					for _, el := range k {
						kvps = append(kvps, el)
					}
					key = ""
				default:
					// raiseerror(`bad sequences. The right list should be:
					// NewGroupedAttrEasy("key", "attr1", 1, "attr2", false)`)
					hintInternal(errUnmatchedPair, "expecting 'key' and 'value' pair in 'args' list, but unmatched 'key' found") // args must be key and value pair, key should be a string
				}
			} else {
				kvps = append(kvps, NewAttr(key, it))
				key = ""
			}
		}
		return
	}

	for _, it := range args {
		if key == "" {
			switch k := it.(type) {
			case string:
				key = k
			case Attr:
				if _, ok := keysKnown[k.Key()]; !ok {
					kvps = append(kvps, k)
					keysKnown[k.Key()] = true
				}
				key = ""
			case []Attr:
				for _, el := range k {
					if _, ok := keysKnown[el.Key()]; !ok {
						kvps = append(kvps, el)
						keysKnown[el.Key()] = true
					}
				}
				key = ""
			case Attrs:
				for _, el := range k {
					if _, ok := keysKnown[el.Key()]; !ok {
						kvps = append(kvps, el)
						keysKnown[el.Key()] = true
					}
				}
				key = ""
			default:
				// raiseerror(`bad sequences. The right list should be:
				// NewGroupedAttrEasy("key", "attr1", 1, "attr2", false)`)
				hintInternal(errUnmatchedPair, "expecting 'key' and 'value' pair in 'args' list, but unmatched 'key' found") // args must be key and value pair, key should be a string
			}
		} else {
			kvps = setUniqueKvp(keysKnown, kvps, key, it)
			keysKnown[key] = true
			key = ""
		}
	}
	return
}

func setUniqueKvp(keys map[string]bool, kvps []Attr, key string, val any) []Attr {
	if _, ok := keys[key]; ok {
		for ix, iv := range kvps {
			if iv.Key() == key {
				kvps[ix].SetValue(val)
				break
			}
		}
	} else {
		kvps = append(kvps, NewAttr(key, val))
		keys[key] = true
	}
	return kvps
}

// NewLogLogger returns a new log.Logger such that each call to its Output method
// dispatches a Record to the specified handler. The logger acts as a bridge from
// the older log API to newer structured logging handlers.
func NewLogLogger(h Logger, lvl Level) *log.Logger {
	return log.New(&handlerWriter{h, lvl, true, 0}, "", 0)
}

type handlerWriter struct {
	l           Logger
	lvl         Level
	capturePC   bool
	extraFrames int
}

func (s *handlerWriter) Write(buf []byte) (n int, err error) {
	if s.lvl >= s.l.Level() {
		var pc uintptr
		if s.capturePC {
			// skip [runtime.Callers, s.Write, Logger.Output, log.Print]
			pc = getpc(4, s.extraFrames)
		}
		if h, ok := s.l.(LogLoggerAware); ok {
			n, err = h.WriteInternal(context.Background(), s.lvl, pc, buf)
		}
	}
	return
}

var errUnmatchedPair = errors.New("unmatched (key,value) pair")

// var err // "args must be key and value pair, key should be a string"
