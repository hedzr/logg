package slog

// New creates a new detached logger and you can make it
// default by SetDefault.
//
// You also call the package-level logging functions directly.
// Such as: Info, Debug, Trace, Warn, Error, Fatal, Panic, ...
//
// There are some special severities by calling OK, Success and Fail.
//
// From a logger, you can make new child logger cascaded with its parent.
// It has different logging context and also share the
// parent's context like common attributes.
//
// First of args must be a string to identify this logger, i.e.,
// it's the logger name.
//
// The rest of args can be these sequences:
//  1. one or more element(s) with type Attr or Attrs
//  2. one or more element(s) with type Opt
//  3. one or more key-value-pair(s)
//
// For example:
//
//	logger := slog.New("standalone-logger-for-app",
//	    slog.NewAttr("attr1", 2),
//	    slog.NewAttrs("attr2", 3, "attr3", 4.1),
//	    "attr4", true, "attr3", "string",
//	    slog.WithLevel(slog.DebugLevel), // an Opt here is allowed
//	)
//
// The logger name is a unique name. Reusing a used name will
// pick the exact child.
//
// Passing an empty name is allowed, a random name will be generated.
//
// For example:
//
//	logger := slog.New("my-app")
func New(args ...any) Logger {
	return newDetachedLogger(args...)
}

type Opt func(s *Entry) // can be passed to New as args

func newDetachedLogger(args ...any) *logimp { return &logimp{newentry(nil, args...)} }

type logimp struct {
	*Entry
}
