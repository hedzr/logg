package slog

import (
	"context"
	"fmt"
	"io"
	logslog "log/slog"
	"os"
	"strconv"
	"strings"
	"time"
	"unsafe"

	"github.com/hedzr/is"
	"github.com/hedzr/is/stringtool"
)

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

type Opt func(s *entry) // can be passed to New as args

func newDetachedLogger(args ...any) *logimp { return &logimp{newentry(nil, args...)} }

func newentry(parent *entry, args ...any) *entry {
	level := GetLevel()
	if parent != nil {
		level = parent.Level()
	}

	s := &entry{
		owner:    parent,
		useColor: true,
		level:    level,
	}

	var todo []any
	for i, o := range args {
		if i == 0 {
			var ok bool
			if s.name, ok = o.(string); ok {
				continue
			}
		}

		if opt, ok := o.(Opt); ok {
			opt(s)
			continue
		}

		if h, ok := o.(logslog.Handler); ok {
			s.handlerOpt = h
		}

		if i > 0 {
			todo = append(todo, o)
		}
	}

	if len(todo) > 0 {
		s.attrs = argsToAttrs(nil, todo...)
	}

	return s
}

type logimp struct {
	*entry
}

type entry struct {
	name  string
	owner *entry
	items map[string]*entry

	// keys []string

	// msg string
	// kvp map[string]any

	useJSON       bool
	useColor      bool
	timeLayout    string
	modeUTC       int // non-set(0), local-time(1), and utc-time(2)
	level         Level
	attrs         []Attr
	writer        *dualWriter
	valueStringer ValueStringer
	handlerOpt    logslog.Handler
	extraFrames   int
	contextKeys   []any
}

// New make a new child Logger, which is identified by a unique name.
//
// In each logger, using same name will pick the exact child.
//
// Passing an empty name is allowed, a random name will be generated.
//
// For example:
//
//	logger := slog.New("my-app")
func (s *entry) New(args ...any) BasicLogger { return s.newChildLogger(args...) }

func (s *entry) newChildLogger(args ...any) *entry {
	if s.items == nil {
		s.items = make(map[string]*entry)
	}

	var name string
	if len(args) == 0 {
		name = stringtool.RandomStringPure(6)
	} else {
		name = args[0].(string)
		if name == "" {
			name = stringtool.RandomStringPure(6)
		}
	}
	if l, ok := s.items[name]; ok {
		return l
	}

	s.items[name] = newentry(s, args...)
	return s.items[name]
}

// func (s *entry) children() map[string]*entry { return s.items } //nolint:unused // to be
// func (s *entry) parent() *entry              { return s.owner } //nolint:unused // to be

func (s *entry) Parent() Entry { return s.owner } //nolint:unused // to be
func (s *entry) Root() Entry {
	p := s
	for p.owner != nil {
		p = p.owner
	}
	return p
}

func (s *entry) Close() {
	// reserved for future
}

func (s *entry) Name() string { return s.name } //nolint:unused // to be

// String implements Logger.
func (s *entry) String() string {
	var sb strings.Builder
	sb.WriteString("*entry{")
	if s.name != "" {
		sb.WriteString(strconv.Quote(s.name))
	}
	sb.WriteString("} (0x")
	sb.WriteString(strconv.FormatUint(uint64(uintptr(unsafe.Pointer(s))), 16))
	sb.WriteRune(',')
	sb.WriteRune(' ')
	sb.WriteString(s.Level().String())
	sb.WriteRune(',')
	sb.WriteRune(' ')
	if s.useJSON {
		sb.WriteString("json")
	} else if s.useColor {
		sb.WriteString("colorful")
	} else {
		sb.WriteString("logfmt")
	}
	if len(s.items) > 0 {
		sb.WriteRune(',')
		sb.WriteRune(' ')
		sb.WriteRune('[')
		for k, v := range s.items {
			sb.WriteRune('\n')
			sb.WriteString(strconv.Quote(k))
			sb.WriteRune('=')
			sb.WriteRune('>')
			sb.WriteString(v.String())
		}
		sb.WriteRune('\n')
		sb.WriteRune(']')
	}
	sb.WriteRune(')')
	return sb.String()
}

//
//
//

func WithJSONMode(b ...bool) Opt {
	return func(s *entry) {
		s.WithJSONMode(b...)
	}
}

func (s *entry) WithJSONMode(b ...bool) Entry {
	mode := true
	for _, bb := range b {
		mode = bb
	}
	if mode {
		s.useColor = false
	}
	s.useJSON = mode
	return s
}

func WithColorMode(b ...bool) Opt {
	return func(s *entry) {
		s.WithColorMode(b...)
	}
}

func (s *entry) WithColorMode(b ...bool) Entry {
	mode := true
	for _, bb := range b {
		mode = bb
	}
	if mode {
		s.useJSON = false
	}
	s.useColor = mode
	return s
}

func WithUTCMode(b ...bool) Opt {
	return func(s *entry) {
		s.WithUTCMode(b...)
	}
}

func (s *entry) WithUTCMode(b ...bool) Entry {
	mode := 2
	for _, bb := range b {
		if bb {
			mode = 2
		} else {
			mode = 1
		}
	}
	s.modeUTC = mode
	return s
}

func WithTimeFormat(layout ...string) Opt {
	return func(s *entry) {
		s.WithTimeFormat(layout...)
	}
}

func (s *entry) WithTimeFormat(layout ...string) Entry {
	var lay = time.RFC3339Nano
	for _, ll := range layout {
		if ll != "" {
			lay = ll
		}
	}
	s.timeLayout = lay
	return s
}

func WithLevel(lvl Level) Opt {
	return func(s *entry) {
		s.WithLevel(lvl)
	}
}

func (s *entry) WithLevel(lvl Level) Entry {
	s.level = lvl
	return s
}

func (s *entry) Level() (lvl Level) {
	// if s.level >= lvlCurrent {
	// 	return s.level
	// }
	// return lvlCurrent
	return s.level
}

// WithAttrs declares some common attributes bound to the
// logger.
//
// When logging, attributes of the logger and its parents
// will be merged together. If duplicated attr found, the
// parent's will be overwritten.
//
//	lc1 := l.New("c1").WithAttrs(NewAttr("lc1", true))
//	lc3 := lc1.New("c3").WithAttrs(NewAttr("lc3", true), NewAttr("lc1", 1))
//	lc3.Warn("lc3 warn msg", "local", false)
//
// In above case, attr 'lc1' will be rewritten while lc3.Warn, it looks like:
//
//	17:47:24.765422+08:00 [WRN] lc3 warn msg          lc1=1 lc3=true local=false
//
// You can initialize attributes with different forms, try
// using WithAttrs1(attrs Attrs) or With(args ...any) for
// instead.
func WithAttrs(attrs ...Attr) Opt {
	return func(s *entry) {
		s.WithAttrs(attrs...)
	}
}

// WithAttrs declares some common attributes bound to the
// logger.
//
// When logging, attributes of the logger and its parents
// will be merged together. If duplicated attr found, the
// parent's will be overwritten.
//
//	lc1 := l.New("c1").WithAttrs(NewAttr("lc1", true))
//	lc3 := lc1.New("c3").WithAttrs(NewAttr("lc3", true), NewAttr("lc1", 1))
//	lc3.Warn("lc3 warn msg", "local", false)
//
// In above case, attr 'lc1' will be rewritten while lc3.Warn, it looks like:
//
//	17:47:24.765422+08:00 [WRN] lc3 warn msg          lc1=1 lc3=true local=false
//
// You can initialize attributes with different forms, try
// using WithAttrs1(attrs Attrs) or With(args ...any) for
// instead.
func (s *entry) WithAttrs(attrs ...Attr) Entry {
	s.attrs = append(s.attrs, attrs...)
	return s
}

// WithAttrs1 allows an Attrs passed into New. Sample is:
//
//	lc1 := l.New("c1").WithAttrs1(NewAttrs("a1", 1, "a2", 2.7, NewAttr("a3", "string")))
//
// NewAttrs receives a freeform args list.
//
// You can use With(...) to simplify WithAttrs1+NewAttrs1 calling.
func WithAttrs1(attrs Attrs) Opt {
	return func(s *entry) {
		s.WithAttrs1(attrs)
	}
}

// WithAttrs1 allows an Attrs passed into New. Sample is:
//
//	lc1 := l.New("c1").WithAttrs1(NewAttrs("a1", 1, "a2", 2.7, NewAttr("a3", "string")))
//
// NewAttrs receives a freeform args list.
//
// You can use With(...) to simplify WithAttrs1+NewAttrs1 calling.
func (s *entry) WithAttrs1(attrs Attrs) Entry {
	s.attrs = append(s.attrs, attrs...)
	return s
}

// With allows an freeform arg list passed into New. Sample is:
//
//	lc1 := l.New("c1").With("a1", 1, "a2", 2.7, NewAttr("a3", "string"))
//
// More samples can be found at New.
func With(args ...any) Opt {
	return func(s *entry) {
		s.With(args...)
	}
}

// With allows an freeform arg list passed into New. Sample is:
//
//	lc1 := l.New("c1").With("a1", 1, "a2", 2.7, NewAttr("a3", "string"))
//
// More samples can be found at New.
func (s *entry) With(args ...any) Entry { // key1,val1,key2,val2,.... Of course, Attr, Attrs in args will be recognized as is
	s.attrs = append(s.attrs, argsToAttrs(nil, args...)...)
	return s
}

type ValueStringer interface {
	SetWriter(w io.Writer)
	WriteValue(value any)
}

func (s *entry) WithValueStringer(vs ValueStringer) Entry {
	s.valueStringer = vs
	return s
}

func GetDefaultWriter() (wr io.Writer)        { return defaultWriter }          // return package-level default writer
func GetDefaultLoggersWriter() (wr io.Writer) { return defaultLog.GetWriter() } // return package-level default logger's writer

// GetWriter looks up the best writer for current level.
func (s *entry) GetWriter() (wr LogWriter) {
	return s.findWriter(s.level)
}

// GetWriterBy returns the leveled writer.
func (s *entry) GetWriterBy(level Level) (wr LogWriter) {
	return s.findWriter(level)
}

// WithWriter sets a std writer to Default logger, the
// original std writers will be cleared.
// It is a Opt functor so you have to invoke it at New(,,,).
//
// For each child loggers, uses their method [entry.WithWriter],
func WithWriter(wr io.Writer) Opt {
	return func(s *entry) {
		s.WithWriter(wr)
	}
}

func (s *entry) WithWriter(wr io.Writer) Entry {
	if s.writer == nil {
		s.writer = newDualWriter()
	}
	s.writer.SetWriter(wr)
	// s.writer.SetLogWriter(wr)
	return s
}

// AddWriter adds a stdout writers to Default logger.
// It is a Opt functor so you have to invoke it at New(,,,).
//
// For each child loggers, uses their method [entry.AddWriter],
func AddWriter(wr io.Writer) Opt {
	return func(s *entry) {
		s.AddWriter(wr)
	}
}

func (s *entry) AddWriter(wr io.Writer) Entry {
	if s.writer == nil {
		s.writer = newDualWriter()
	}
	s.writer.Add(wr)
	// s.writer.SetLogWriter(wr)
	return s
}

// AddErrorWriter adds a stderr writers to Default logger.
// It is a Opt functor so you have to invoke it at New(,,,).
//
// For each child loggers, uses their method [entry.AddErrorWriter],
func AddErrorWriter(wr io.Writer) Opt {
	return func(s *entry) {
		s.AddErrorWriter(wr)
	}
}

func (s *entry) AddErrorWriter(wr io.Writer) Entry {
	if s.writer == nil {
		s.writer = newDualWriter()
	}
	s.writer.SetErrorWriter(wr)
	// s.writer.SetLogWriter(wr)
	return s
}

// ResetWriters clear all stdout and stderr writers in Default logger.
// It is a Opt functor so you have to invoke it at New(,,,).
//
// For each child loggers, uses their method [entry.ResetWriters],
func ResetWriters() Opt {
	return func(s *entry) {
		s.ResetWriters()
	}
}

func (s *entry) ResetWriters() Entry {
	if s.writer == nil {
		s.writer = newDualWriter()
	}
	s.writer.Reset()
	return s
}

// AddLevelWriter add a leveled writer in Default logger.
// It is a Opt functor so you have to invoke it at New(,,,).
//
// For each child loggers, uses their method [entry.AddLevelWriter],
//
// A leveled writer has higher priorities than normal writers,
// see AddWriter and AddErrorWriter.
func AddLevelWriter(lvl Level, w io.Writer) Opt {
	return func(s *entry) {
		s.AddLevelWriter(lvl, w)
	}
}

func (s *entry) AddLevelWriter(lvl Level, w io.Writer) Entry {
	if s.writer == nil {
		s.writer = newDualWriter()
	}
	s.writer.AddLevelWriter(lvl, w)
	return s
}

// RemoveLevelWriter remove a leveled writer in Default logger.
// It is a Opt functor so you have to invoke it at New(,,,).
//
// For each child loggers, uses their method [entry.RemoveLevelWriter],
func RemoveLevelWriter(lvl Level, w io.Writer) Opt {
	return func(s *entry) {
		s.RemoveLevelWriter(lvl, w)
	}
}

func (s *entry) RemoveLevelWriter(lvl Level, w io.Writer) Entry {
	if s.writer == nil {
		s.writer = newDualWriter()
	}
	s.writer.RemoveLevelWriter(lvl, w)
	return s
}

// ResetLevelWriter clear anu leveled writers for a specified
// Level in Default logger.
// It is a Opt functor so you have to invoke it at New(,,,).
//
// For each child loggers, uses their method [entry.ResetLevelWriter],
func ResetLevelWriter(lvl Level) Opt {
	return func(s *entry) {
		s.ResetLevelWriter(lvl)
	}
}

func (s *entry) ResetLevelWriter(lvl Level) Entry {
	if s.writer == nil {
		s.writer = newDualWriter()
	}
	s.writer.ResetLevelWriter(lvl)
	return s
}

// ResetLevelWriters clear all leveled writers in Default logger.
// It is a Opt functor so you have to invoke it at New(,,,).
//
// For each child loggers, uses their method [entry.ResetLevelWriters],
func ResetLevelWriters() Opt {
	return func(s *entry) {
		s.ResetLevelWriters()
	}
}

func (s *entry) ResetLevelWriters() Entry {
	if s.writer == nil {
		s.writer = newDualWriter()
	}
	s.writer.ResetLevelWriters()
	return s
}

// WithSkip make a child logger from Default and set the extra ignored frames with given value.
//
// By default, LattrsR is not enabled. So the new child logger cannot
// print the parent's attrs. You could have to AddFlags(LattrsR) to
// give WithSkip a better behavior.
//
// If you dislike to make another one new child Logger instance, using SetSkip pls.
func WithSkip(extraFrames int) Entry { return defaultLog.WithSkip(extraFrames) }
func SetSkip(extraFrames int)        { defaultLog.SetSkip(extraFrames) } // set extra frames ignored

// WithSkip make a child of itself and set the extra ignored frames with given integer.
//
// By default, LattrsR is not enabled. So the new child logger cannot
// print the parent's attrs. You could have to AddFlags(LattrsR) to
// make WithSkip a better behavior.
//
// If you dislike to make another one new child Logger instance, using SetSkip pls.
func (s *entry) WithSkip(extraFrames int) Entry {
	return s.newChildLogger(fmt.Sprintf("c/%s[%d]", s.name, extraFrames)).withSkip(extraFrames)
}

func (s *entry) withSkip(extraFrames int) Entry {
	s.extraFrames = extraFrames
	return s
}

func (s *entry) SetSkip(extraFrames int) { s.extraFrames = extraFrames } // set extra frames ignored
func (s *entry) Skip() int               { return s.extraFrames }        // return extra frames ignored

//
//
//

func (s *entry) Enabled(lvl Level) bool { return s.Level().Enabled(context.TODO(), lvl) }
func (s *entry) EnabledContext(ctx context.Context, lvl Level) bool {
	return s.Level().Enabled(ctx, lvl)
}

//
//
//

func (s *entry) Panic(msg string, args ...any)   { s.log1(PanicLevel, msg, args...) }   // Panic implements Logger.
func (s *entry) Fatal(msg string, args ...any)   { s.log1(FatalLevel, msg, args...) }   // Fatal implements Logger.
func (s *entry) Error(msg string, args ...any)   { s.log1(ErrorLevel, msg, args...) }   // Error implements Logger.
func (s *entry) Warn(msg string, args ...any)    { s.log1(WarnLevel, msg, args...) }    // Warn implements Logger.
func (s *entry) Info(msg string, args ...any)    { s.log1(InfoLevel, msg, args...) }    // Info implements Logger.
func (s *entry) Debug(msg string, args ...any)   { s.log1(DebugLevel, msg, args...) }   // Debug implements Logger.
func (s *entry) Trace(msg string, args ...any)   { s.log1(TraceLevel, msg, args...) }   // Trace implements Logger.
func (s *entry) Print(msg string, args ...any)   { s.log1(AlwaysLevel, msg, args...) }  // Print implements Logger.
func (s *entry) OK(msg string, args ...any)      { s.log1(OKLevel, msg, args...) }      // OK implements Logger.
func (s *entry) Success(msg string, args ...any) { s.log1(SuccessLevel, msg, args...) } // Success implements Logger.
func (s *entry) Fail(msg string, args ...any)    { s.log1(FailLevel, msg, args...) }    // Fail implements Logger.
func (s *entry) Println(args ...any) {
	if len(args) == 0 {
		s.log1(AlwaysLevel, "")
		return
	}
	s.log1(AlwaysLevel, args[0].(string), args[1:]...)
} // Println implements Logger.

//
//

//

// PanicContext implements Logger.
func (s *entry) PanicContext(ctx context.Context, msg string, args ...any) {
	if s.EnabledContext(ctx, PanicLevel) {
		pc := getpc(2, s.extraFrames)
		s.logContext(ctx, PanicLevel, pc, msg, args...)
	}
}

// FatalContext implements Logger.
func (s *entry) FatalContext(ctx context.Context, msg string, args ...any) {
	if s.EnabledContext(ctx, FatalLevel) {
		pc := getpc(2, s.extraFrames)
		s.logContext(ctx, FatalLevel, pc, msg, args...)
	}
}

// ErrorContext implements Logger.
func (s *entry) ErrorContext(ctx context.Context, msg string, args ...any) {
	if s.EnabledContext(ctx, ErrorLevel) {
		pc := getpc(2, s.extraFrames)
		s.logContext(ctx, ErrorLevel, pc, msg, args...)
	}
}

// WarnContext implements Logger.
func (s *entry) WarnContext(ctx context.Context, msg string, args ...any) {
	if s.EnabledContext(ctx, WarnLevel) {
		pc := getpc(2, s.extraFrames)
		s.logContext(ctx, WarnLevel, pc, msg, args...)
	}
}

// InfoContext implements Logger.
func (s *entry) InfoContext(ctx context.Context, msg string, args ...any) {
	if s.EnabledContext(ctx, InfoLevel) {
		pc := getpc(2, s.extraFrames)
		s.logContext(ctx, InfoLevel, pc, msg, args...)
	}
}

// DebugContext implements Logger.
func (s *entry) DebugContext(ctx context.Context, msg string, args ...any) {
	if s.EnabledContext(ctx, DebugLevel) {
		pc := getpc(2, s.extraFrames)
		s.logContext(ctx, DebugLevel, pc, msg, args...)
	}
}

// TraceContext implements Logger.
func (s *entry) TraceContext(ctx context.Context, msg string, args ...any) {
	if s.EnabledContext(ctx, TraceLevel) {
		pc := getpc(2, s.extraFrames)
		s.logContext(ctx, TraceLevel, pc, msg, args...)
	}
}

// PrintContext implements Logger.
func (s *entry) PrintContext(ctx context.Context, msg string, args ...any) {
	pc := getpc(2, s.extraFrames)
	s.logContext(ctx, AlwaysLevel, pc, msg, args...)
	// panic("unimplemented")
}

// PrintlnContext implements Logger.
func (s *entry) PrintlnContext(ctx context.Context, msg string, args ...any) {
	pc := getpc(2, s.extraFrames)
	s.logContext(ctx, AlwaysLevel, pc, msg, args...)
}

// OKContext implements Logger.
func (s *entry) OKContext(ctx context.Context, msg string, args ...any) {
	pc := getpc(2, s.extraFrames)
	s.logContext(ctx, OKLevel, pc, msg, args...)
}

// SuccessContext implements Logger.
func (s *entry) SuccessContext(ctx context.Context, msg string, args ...any) {
	pc := getpc(2, s.extraFrames)
	s.logContext(ctx, SuccessLevel, pc, msg, args...)
}

// FailContext implements Logger.
func (s *entry) FailContext(ctx context.Context, msg string, args ...any) {
	pc := getpc(2, s.extraFrames)
	s.logContext(ctx, FailLevel, pc, msg, args...)
}

// LogAttrs implements Logger.
func (s *entry) LogAttrs(ctx context.Context, level Level, msg string, args ...any) {
	if s.EnabledContext(ctx, level) {
		pc := getpc(2, s.extraFrames)
		s.logContext(ctx, level, pc, msg, args...)
	}
}

// Log implements Logger.
func (s *entry) Log(ctx context.Context, level Level, msg string, args ...any) {
	if s.EnabledContext(ctx, level) {
		pc := getpc(2, s.extraFrames)
		s.logContext(ctx, level, pc, msg, args...)
	}
}

//

//

func (s *entry) WriteThru(ctx context.Context, lvl Level, timestamp time.Time, stackFrame uintptr, msg string, attrs Attrs) {
	s.print(ctx, lvl, timestamp, stackFrame, msg, attrs)
}

func (s *entry) WriteInternal(ctx context.Context, lvl Level, stackFrame uintptr, buf []byte) (n int, err error) {
	return s.writeInternal(ctx, lvl, stackFrame, buf)
}

func (s *entry) writeInternal(ctx context.Context, lvl Level, stackFrame uintptr, buf []byte) (n int, err error) {
	// Remove final newline.
	origLen := len(buf) // Report that the entire buf was written.
	if len(buf) > 0 && buf[len(buf)-1] == '\n' {
		buf = buf[:len(buf)-1]
	}
	n = origLen
	now := time.Now()
	s.print(ctx, lvl, now, stackFrame, string(buf), nil)
	return
}

func (s *entry) parseArgs(ctx context.Context, lvl Level, stackFrame uintptr, msg string, args ...any) (kvps Attrs) {
	var roughSize = 4 + len(s.attrs) + len(args)
	// var key string
	// var keys = make(map[string]bool, roughSize)

	kvps = s.leadingTags(roughSize, lvl, stackFrame, msg)

	if s.ctxKeysWanted() {
		kvps = append(kvps, s.fromCtx(ctx)...)
	}
	if len(s.attrs) > 0 {
		kvps = append(kvps, s.walkParentAttrs(ctx, lvl, s, nil)...)
	}
	if len(args) > 0 {
		kvps = append(kvps, argsToAttrs(nil, args...)...)
	}

	// if key != "" {
	// 	kvps = setUniqueKvp(keys, kvps, BADKEY, key)
	// }

	// for _, it := range s.attrs {
	// 	kvps = setUniqueKvp(keys, kvps, it.Key(), it.Value())
	// }

	return
}

func (s *entry) walkParentAttrs(ctx context.Context, lvl Level, e *entry, keysKnown map[string]bool) (kvps []Attr) {
	if e == nil {
		return
	}

	if keysKnown == nil {
		return e.attrs
	}

	// try appending unique attributes and walk all parents

	var roughlen = len(e.attrs)
	if roughlen == 0 && !IsAnyBitsSet(LattrsR) {
		return
	}

	if roughlen < 8 {
		roughlen = 8
	}

	kvps = make([]Attr, 0, roughlen)

	if IsAnyBitsSet(LattrsR) {
		// lookup parents
		if p := e.owner; p != nil {
			parentTags := s.walkParentAttrs(ctx, lvl, p, keysKnown)
			kvps = append(kvps, parentTags...)
		}
	}

	for _, attr := range e.attrs {
		key := attr.Key()
		if _, ok := keysKnown[key]; ok {
			for ix, iv := range kvps {
				if iv.Key() == key {
					kvps[ix].SetValue(attr.Value())
					break
				}
			}
		} else {
			kvps = append(kvps, attr)
			keysKnown[key] = true
		}
	}
	return
}

func (s *entry) leadingTags(roughSize int, lvl Level, stackFrame uintptr, msg string) (kvps Attrs) {
	// if s.useJSON || !s.useColor {
	// 	kvps = make(Attrs, 0, roughSize+4) // pre-alloc slice spaces roughly
	//
	// 	// simulate logfmt format here, while non-JSON and non-colorful mode.
	// 	// these key-value-pairs fit for serializing in JSON mode too.
	// 	kvps = append(kvps,
	// 		&kvp{timestampFieldName, time.Now()},
	// 		&kvp{levelFieldName, lvl},
	// 		&kvp{messageFieldName, msg})
	//
	// 	if IsAnyBitsSet(Lcaller) {
	// 		source := getpcsource(stackFrame)
	// 		kvps = append(kvps, source.toGroup())
	// 	}
	// }
	if kvps == nil && roughSize > 0 {
		kvps = make(Attrs, 0, roughSize) // pre-alloc slice spaces roughly
	}
	return
}

func (s *entry) WithContextKeys(keys ...any) Entry {
	s.contextKeys = append(s.contextKeys, keys...)
	return s
}

func (s *entry) ctxKeysWanted() bool { return len(s.contextKeys) > 0 }
func (s *entry) ctxKeys() []any      { return s.contextKeys }
func (s *entry) fromCtx(ctx context.Context) (ret Attrs) {
	mu := make(map[string]struct{})
	for _, k := range s.ctxKeys() {
		if v := ctx.Value(k); v != nil {
			switch key := k.(type) {
			case Stringer:
				kk := key.String()
				if _, ok := mu[kk]; !ok {
					ret = append(ret, &kvp{kk, v})
					mu[kk] = struct{}{}
				}
			case string:
				if _, ok := mu[key]; !ok {
					ret = append(ret, &kvp{key, v})
					mu[key] = struct{}{}
				}
			}
		}
	}
	return
}

// var poolHelper = Pool(
// 	newPrintCtx,
// 	func(ctx *PrintCtx, i int) []byte {
// 		return nil
// 	})
// ret = poolHelper(1)
//
// func(pc *PrintCtx, cb func(s *entry, ctx context.Context, pc *PrintCtx) (ret []byte)) (ret []byte) {
// 	pc.set(s, lvl, timestamp, stackFrame, msg, kvps)
// 	// pc := newPrintCtx(s, lvl, timestamp, stackFrame, msg, kvps)
// 	return cb(ctx, pc)
// })

func (s *entry) print(ctx context.Context, lvl Level, timestamp time.Time, stackFrame uintptr, msg string, kvps Attrs) (ret []byte) {
	pc := printCtxPool.Get().(*PrintCtx)
	defer func() { printCtxPool.Put(pc) }()

	// pc := newPrintCtx(s, lvl, timestamp, stackFrame, msg, kvps)

	// pc.set will truncate internal buffer and reset all states for
	// this current session. So, don't worry about a reused buffer
	// takes wasted bytes.
	pc.set(s, lvl, timestamp, stackFrame, msg, kvps)

	return s.printImpl(ctx, pc)
}

func (s *entry) printImpl(ctx context.Context, pc *PrintCtx) (ret []byte) {
	if pc.lvl == AlwaysLevel && strings.Trim(pc.msg, "\n\r \t") == "" {
		pc.pcAppendByte('\n')
		ret = pc.Bytes()
		s.printOut(pc.lvl, ret)
		return
	}

	if pc.jsonMode {
		pc.pcAppendByte('{')
	}

	if pc.noColor { // json or logfmt
		s.printTimestamp(ctx, pc)
		s.printLoggerName(ctx, pc)
		s.printSeverity(ctx, pc)
		s.printMsg(ctx, pc)
	} else {
		if aa, ok := mLevelColors[pc.lvl]; ok {
			pc.clr = aa[0]
			if len(aa) > 1 {
				pc.bg = aa[1]
			}
		}

		s.printTimestamp(ctx, pc)
		s.printLoggerName(ctx, pc)
		s.printSeverity(ctx, pc)
		s.printFirstLineOfMsg(ctx, pc)
	}

	serializeAttrs(pc, pc.kvps)

	if IsAnyBitsSet(Lcaller) {
		s.printPC(ctx, pc)
	}

	s.printRestLinesOfMsg(ctx, pc)

	if pc.jsonMode {
		pc.pcAppendByte('}')
	}

	pc.pcAppendByte('\n')

	// ret = pc.String()
	// s.printOut(pc.lvl, []byte(ret))
	ret = pc.Bytes()
	s.printOut(pc.lvl, ret)
	return
}

func (s *entry) printTimestamp(ctx context.Context, pc *PrintCtx) {
	if pc.noColor { // json or logfmt
		pc.pcAppendStringKey(timestampFieldName)
		pc.pcAppendColon()
		// pc.pcAppendByte('"')
		pc.appendTimestamp(pc.now)
		pc.pcAppendComma()
	} else {
		ct.echoColor(pc, clrTimestamp)
		pc.appendTimestamp(pc.now)
		pc.pcAppendByte(' ')
	}
}

func (s *entry) printLoggerName(ctx context.Context, pc *PrintCtx) {
	if s.name != "" {
		if pc.noColor { // json or logfmt
			pc.AddString("logger", s.name)
			// pc.pcAppendStringKey("logger")
			// pc.pcAppendColon()
			// pc.pcAppendByte('"')
			// pc.pcAppendStringValue(s.name)
			// pc.pcAppendByte('"')
			pc.pcAppendComma()
		} else {
			ct.wrapColorAndBgTo(pc, clrLoggerName, clrLoggerNameBg, s.name)
			pc.pcAppendByte(' ')
		}
	}
}

func (s *entry) printSeverity(ctx context.Context, pc *PrintCtx) {
	if pc.noColor { // json or logfmt
		pc.AddString(levelFieldName, pc.lvl.String())
		// pc.pcAppendStringKey(levelFieldName)
		// pc.pcAppendColon()
		// pc.pcAppendByte('"')
		// pc.pcAppendStringValue(pc.lvl.String())
		// pc.pcAppendByte('"')
		pc.pcAppendComma()
	} else {
		ct.wrapColorAndBgTo(pc, pc.clr, pc.bg, ct.wrapRune(pc.lvl.ShortTag(levelOutputWidth), '[', ']'))
		pc.pcAppendByte(' ')
	}
}

func (s *entry) printPC(ctx context.Context, pc *PrintCtx) {
	if pc.noColor {
		pc.pcAppendComma()

		source := pc.source()
		if pc.jsonMode {
			pc.pcAppendStringKey(callerFieldName)
			pc.pcAppendColon()
			pc.pcAppendByte('{')

			pc.AddString("file", source.File)
			pc.pcAppendComma()
			pc.AddInt("line", source.Line)
			pc.pcAppendComma()
			pc.AddString("function", source.Function)

			pc.pcAppendByte('}')
		} else {
			pc.AddPrefixedString(callerFieldName, "file", source.File)
			pc.pcAppendComma()

			pc.AddPrefixedInt(callerFieldName, "line", source.Line)
			pc.pcAppendComma()

			pc.AddPrefixedString(callerFieldName, "function", source.Function)
		}
		// pc.pcAppendComma()
		return
	}

	source := pc.source()
	pc.pcAppendByte(' ')
	// pc.appendRune('(')
	pc.pcAppendString(source.File)
	pc.pcAppendByte(':')
	pc.AppendInt(source.Line)
	// pc.appendRune(')')
	pc.pcAppendByte(' ')
	// ct.wrapDimColorTo(pc.SB, source.checkedfuncname()) // clion p-term in run panel cannot support dim color.
	ct.wrapColorTo(pc, clrFuncName, checkedfuncname(source.Function))
	ct.echoResetColor(pc)
}

func (s *entry) printMsg(ctx context.Context, pc *PrintCtx) {
	pc.AddString(messageFieldName, ct.translate(pc.msg))
	// pc.pcAppendComma()
}

func (s *entry) printFirstLineOfMsg(ctx context.Context, pc *PrintCtx) {
	var firstLine string
	firstLine, pc.restLines, pc.eol = ct.splitFirstAndRestLines(pc.msg)
	if minimalMessageWidth > 0 {
		str := ct.rightPad(firstLine, " ", minimalMessageWidth)
		str = ct.translate(str)
		_, _ = pc.WriteString(ct.wrapColorAndBg(str, pc.clr, pc.bg))
	} else {
		str := ct.translate(firstLine)
		_, _ = pc.WriteString(ct.wrapColorAndBg(str, pc.clr, pc.bg))
	}
	// pc.pcAppendByte(' ')
	// pc.pcAppendByte('|')
}

func (s *entry) printRestLinesOfMsg(ctx context.Context, pc *PrintCtx) {
	if !pc.noColor && pc.restLines != "" {
		pc.pcAppendByte('\n')
		pc.pcAppendString(ct.padFunc(pc.restLines, " ", 4, func(i int, line string) string {
			return ct.wrapColorAndBg(line, pc.clr, pc.bg)
		}))
		if pc.eol {
			pc.pcAppendByte('\n')
		}
	}
}

func (s *entry) printOut(lvl Level, ret []byte) {
	if w := s.findWriter(lvl); w != nil {
		_, err := w.Write(ret)
		if err != nil && lvl != WarnLevel { // don't warn on warning to avoid infinite calls
			s.Warn("slog print log failed", "error", err)
		}
	}
}

func (s *entry) log1(lvl Level, msg string, args ...any) {
	ctx := context.Background()
	if s.EnabledContext(ctx, lvl) {
		stackFrame := getpc(3, s.extraFrames)
		s.logContext(ctx, lvl, stackFrame, msg, args...)
	}
}

// func (s *entry) logc1(ctx context.Context, lvl Level, msg string, args ...any) {
// 	if s.Enabled(lvl) {
// 		pc := getpc(3)
// 		s.logContext(ctx, lvl, pc, msg, args...)
// 	}
// }

func (s *entry) logContext(ctx context.Context, lvl Level, stackFrame uintptr, msg string, args ...any) {
	if hh := s.handlerOpt; hh != nil {
		level := convertLevelToLogSlog(lvl)
		if hh.Enabled(ctx, level) {
			// todo handler
		}
		return
	}

	if ctx == nil {
		ctx = context.TODO()
	}

	// if s.EnabledContext(ctx, lvl) {
	now := time.Now()
	kvps := s.parseArgs(ctx, lvl, stackFrame, msg, args...)
	s.print(ctx, lvl, now, stackFrame, msg, kvps)

	if !is.InTesting() || IsAnyBitsSet(Linterruptalways) {
		if !IsAllBitsSet(LnoInterrupt) {
			if lvl == PanicLevel {
				panic(msg)
			}
			if lvl == FatalLevel {
				os.Exit(-3)
			}
		}
	}
	// }
}

func (s *entry) findWriter(lvl Level) LogWriter {
	if s.writer != nil {
		if r := s.writer.Get(lvl); r != nil {
			return r
		}
	}
	return defaultWriter.Get(lvl)
}

const BADKEY = "!BADKEY"
