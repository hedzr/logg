package sub

import (
	"context"
	"fmt"
	stdlog "log/slog"
	"sync"

	logz "github.com/hedzr/logg/slog"
)

// NewChildLogger make a new log/slog logger associated with an
// underlying logz logger (hedzr/logg/slog as logz).
func NewChildLogger(name string) *stdlog.Logger {
	log00 := logz.New(name).SetLevel(logz.DebugLevel)
	// log00.Verbose("init dbg-log")

	log = log00.
		WithSkip(1) // extra stack frame(s) shall be ignored for dbglog.Info/...

	const addSource = true
	sll := logz.NewSlogHandler(log, &logz.HandlerOptions{
		NoColor:  false,
		NoSource: !addSource,
		JSON:     useJSON,
		Level:    logz.DebugLevel,
	})

	// lvl := new(stdlog.LevelVar)
	// lvl.Set(stdlog.LevelInfo)
	// stdlogger1 := stdlog.New(stdlog.NewTextHandler(os.Stderr, &stdlog.HandlerOptions{
	// 	Level:       lvl,
	// 	AddSource:   addSource,
	// 	ReplaceAttr: nil,
	// }))

	stdLogger := stdlog.New(sll)

	mStdToLogz[stdLogger] = log

	return stdLogger
}

var mStdToLogz map[*stdlog.Logger]logz.Logger

func init() {
	mStdToLogz = make(map[*stdlog.Logger]logz.Logger)

	log00 := logz.New(applogname).SetLevel(logz.InfoLevel)
	log00.Verbose("init dbg-log")
	// log00.Warn("applog(ger) created")

	log = log00.
		WithSkip(1) // extra stack frame(s) shall be ignored for dbglog.Info/...
	Verbose("[cmdr.service] applog(ger)-chld created")

	wrs = &wrS{log00}

	// attach 'log' into log/slog
	sll := logz.NewSlogHandler(log, &logz.HandlerOptions{
		NoColor:  false,
		NoSource: false,
		JSON:     useJSON,
		Level:    logz.InfoLevel,
	})
	Logger = stdlog.New(sll)

	mStdToLogz[Logger] = log

	// sync output format to all these loggers
	SetJSONMode(useJSON)

	// ctx := context.Background()
	// InfoContext(ctx, "hello, world")
	// DebugContext(ctx, "hello, world")
}

type wrS struct{ logz.Logger }

func (s *wrS) Write(data []byte) (n int, err error) {
	ctx := context.Background()
	s.Logit(ctx, logz.InfoLevel, string(data))
	return
}

//

//

func Infof(msg string, args ...any) {
	// Logger.Info(msg, args...) // NOTE, std log/slog cannot ignore extra stack frame(s)
	log.Infof(msg, args...)
}

func Warnf(msg string, args ...any) {
	log.Warnf(msg, args...)
}

func Errorf(msg string, args ...any) {
	log.Errorf(msg, args...)
}

func Debugf(msg string, args ...any) {
	log.Debug(fmt.Sprintf(msg, args...))
}

func Tracef(msg string, args ...any) {
	log.Trace(fmt.Sprintf(msg, args...))
}

func Fatalf(msg string, args ...any) {
	log.Fatal(fmt.Sprintf(msg, args...))
}

func Panicf(msg string, args ...any) {
	log.Panic(fmt.Sprintf(msg, args...))
}

func Info(msg string, args ...any) {
	// Logger.Info(msg, args...) // NOTE, std log/slog cannot ignore extra stack frame(s)
	log.Info(msg, args...)
}

func Warn(msg string, args ...any) {
	log.Warn(msg, args...)
}

func Error(msg string, args ...any) {
	log.Error(msg, args...)
}

func Debug(msg string, args ...any) {
	log.Debug(msg, args...)
}

func Trace(msg string, args ...any) {
	log.Trace(msg, args...)
}

func Verbose(msg string, args ...any) {
	log.Verbose(msg, args...)
}

func Panic(msg string, args ...any) {
	log.Panic(msg, args...)
}

func Fatal(msg string, args ...any) {
	log.Fatal(msg, args...)
}

func Print(msg string, args ...any) {
	log.Print(msg, args...)
}

func Println(args ...any) {
	log.Println(args...)
}

func Printf(msg string, args ...any) {
	log.Println(fmt.Sprintf(msg, args...))
}

func OK(msg string, args ...any) {
	log.OK(msg, args...)
}

func Fail(msg string, args ...any) {
	log.Fail(msg, args...)
}

func Success(msg string, args ...any) {
	log.Success(msg, args...)
}

func InfoContext(ctx context.Context, msg string, args ...any) {
	log.InfoContext(ctx, msg, args...)
}

func WarnContext(ctx context.Context, msg string, args ...any) {
	log.WarnContext(ctx, msg, args...)
}

func ErrorContext(ctx context.Context, msg string, args ...any) {
	log.ErrorContext(ctx, msg, args...)
}

func DebugContext(ctx context.Context, msg string, args ...any) {
	log.DebugContext(ctx, msg, args...)
}

func TraceContext(ctx context.Context, msg string, args ...any) {
	// if is.Tracing() {
	// 	log.DebugContext(ctx, msg, args...)
	// }
	log.TraceContext(ctx, msg, args...)
}

func VerboseContext(ctx context.Context, msg string, args ...any) {
	// if is.VerboseBuild() {
	// 	log.DebugContext(ctx, msg, args...)
	// }
	log.VerboseContext(ctx, msg, args...)
}

func PanicContext(ctx context.Context, msg string, args ...any) {
	log.PanicContext(ctx, msg, args...)
}

func FatalContext(ctx context.Context, msg string, args ...any) {
	log.FatalContext(ctx, msg, args...)
}

func PrintContext(ctx context.Context, msg string, args ...any) {
	log.PrintContext(ctx, msg, args...)
}

func PrintlnContext(ctx context.Context, msg string, args ...any) {
	log.PrintlnContext(ctx, msg, args...)
}

func OKContext(ctx context.Context, msg string, args ...any) {
	log.OKContext(ctx, msg, args...)
}

func FailContext(ctx context.Context, msg string, args ...any) {
	log.FailContext(ctx, msg, args...)
}

func SuccessContext(ctx context.Context, msg string, args ...any) {
	log.SuccessContext(ctx, msg, args...)
}

func SetLevel(level logz.Level) {
	log.SetLevel(level)
}

func GetLevel() logz.Level { return log.Level() }

func SetColorMode(mode bool) {
	if p := log.Parent(); p != nil {
		p.SetColorMode(mode)
	}

	if mode == false {
		log.SetJSONMode(true)
	}
	log.SetColorMode(mode)
	ZLogger().SetColorMode(mode)
}

func SetJSONMode(mode bool) {
	if p := log.Parent(); p != nil {
		p.SetJSONMode(mode)
	}

	log.SetJSONMode(mode)
	ZLogger().SetJSONMode(mode)
	// // sync cmdr's internal logger to json mode
	// cmdrlogz.SetJSONMode(mode)
}

func SetSkip(skip int) {
	log.SetSkip(skip)
}

func AddSkip(delta int) {
	log.SetSkip(delta + log.Skip())
	sll := logz.NewSlogHandler(log, &logz.HandlerOptions{
		NoColor:  false,
		NoSource: false,
		JSON:     useJSON,
		Level:    logz.InfoLevel,
	})
	Logger = stdlog.New(sll)
}

//

//

// WrappedLogger returns a reference to *slog.Logger which was
// wrapped to hedzr/logg/slog.
//
// In most cases, you'd better use dbglog.Info/... directly because
// these forms can locate the preferred stack frame(s) of the caller.
func WrappedLogger() *stdlog.Logger { return Logger }

func RawLogger() logz.Logger { return log }

//

//

func ZLogger() *wrS {
	return wrs
}

//

//

const applogname = "sub.log"
const useJSON = false

var Logger *stdlog.Logger
var log logz.Logger
var wrs *wrS
var onceLog sync.Once
