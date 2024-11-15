package slog

import (
	"regexp"
	"sync"

	"github.com/hedzr/is/term/color"
)

// AddCodeHostingProviders appends more provider, repl pair to reduce the caller text width.
//
// The builtin providers are:
//   - "github.com" -> "GH"
//   - "gitlab.com" -> "GL"
//   - "gitee.com" -> "GT"
//   - "bitbucket.com" -> "BB"
//   - ...
func AddCodeHostingProviders(provider, repl string) { codeHostingProvidersMap[provider] = repl }

// AddKnownPathMapping appends more (pathname, repl) pair to reduce the caller filepath width.
//
// Such as:
//   - "$HOME" -> "~"
//   - pwd -> "." (current directory -> '.', that means any abs-path will be converted to rel-path)
func AddKnownPathMapping(pathname, repl string) { knownPathMap[pathname] = repl }

// RemoveKnownPathMapping _
func RemoveKnownPathMapping(pathname string) {
	delete(knownPathMap, pathname)
}

// ResetKnownPathMapping _
func ResetKnownPathMapping() {
	clear(knownPathMap) // just for go1.21+
}

// AddKnownPathRegexpMapping adds regexp pattern, repl pair to reduce the called filepath width.
func AddKnownPathRegexpMapping(pathnameRegexpExpr, repl string) {
	knownPathRegexpMap = append(knownPathRegexpMap, regRepl{
		expr: regexp.MustCompile(pathnameRegexpExpr),
		repl: repl,
	})
}

// RemoveKnownPathRegexpMapping _
func RemoveKnownPathRegexpMapping(pathnameRegexpExpr string) {
	for i, vv := range knownPathRegexpMap {
		if vv.expr.String() == pathnameRegexpExpr {
			knownPathRegexpMap = append(knownPathRegexpMap[:i], knownPathRegexpMap[i+1:]...)
			return
		}
	}
}

// ResetKnownPathRegexpMapping clears all rules in knownPathRegexpMap.
func ResetKnownPathRegexpMapping() {
	knownPathRegexpMap = nil
}

// SetLevelOutputWidth sets how many characters of level string should
// be formatted and output to logging lines.
//
// While you are customizing your level, a 1..5 characters array is
// required for formatting purpose.
//
// For example:
//
//	const NoticeLevel = slog.Level(17) // A custom level must have a value greater than slog.MaxLevel
//	slog.RegisterLevel(NoticeLevel, "NOTICE",
//	    slog.RegWithShortTags([6]string{"", "N", "NT", "NTC", "NOTC", "NOTIC"}),
//	    slog.RegWithColor(color.FgWhite, color.BgUnderline),
//	    slog.RegWithTreatedAsLevel(slog.InfoLevel),
//	))
func SetLevelOutputWidth(width int) {
	if width >= 0 && width <= 5 {
		levelOutputWidth = width
	}
}

// SetMessageMinimalWidth modify the minimal width between
// message and attributes. It works for only colorful mode.
//
// The default width is 36, that means a message will be left
// padding to take 36 columns, filled by space char (' ').
func SetMessageMinimalWidth(w int) {
	if w >= 16 {
		minimalMessageWidth = w
	}
}

func GetFlags() Flags           { return flags }        // returns logg/slog Flags
func SetFlags(f Flags)          { flags = f }           // sets logg/slog Flags
func ResetFlags()               { flags = LstdFlags }   // resets logg/slog Flags to factory settings
func IsAnyBitsSet(f Flags) bool { return flags&f != 0 } // detects if any of some Flags are set
func IsAllBitsSet(f Flags) bool { return flags&f == f } // detects if all of given Flags are both set

// AddFlags adds some Flags (bitwise Or operation).
//
// These ways can merge flags into internal settings:
//
//	AddFlags(Lprivacypathregexp | Lprivacypath)
//	AddFlags(Lprivacypathregexp, Lprivacypath)
func AddFlags(flagsToAdd ...Flags) {
	for _, f := range flagsToAdd {
		flags |= f
		Verbose("add a flag", "flag", f)
	}
}

// RemoveFlags removes some Flags (bitwise And negative operation).
//
// These ways can strip flags off from internal settings:
//
//	RemoveFlags(Lprivacypathregexp | Lprivacypath)
//	RemoveFlags(Lprivacypathregexp, Lprivacypath)
func RemoveFlags(flagsToRemove ...Flags) {
	for _, f := range flagsToRemove {
		flags &= ^f
	}
}

// SaveFlagsAndMod saves old flags, modify it, and restore the old at defer time.
//
// A typical usage might be:
//
//	// Inside a test case, you wanna add date part to timestamp output,
//	// and disable panic (by Panic or Fatal) breaking the testing flow.
//	// So this line will make those temporary modifications:
//	defer SaveFlagsAndMod(Ldate | LnoInterrupt)()
//	// ...
//	// concrete testing codes here
func SaveFlagsAndMod(addingFlags Flags, removingFlags ...Flags) (deferFn func()) {
	save := flags
	AddFlags(addingFlags)
	for _, f := range removingFlags {
		RemoveFlags(f)
	}
	return func() {
		flags = save
	}
}

// func SaveFlagsAndMod(actions ...func()) (deferFn func()) {
// 	save := flags
// 	for _, action := range actions {
// 		action()
// 	}
// 	return func() {
// 		flags = save
// 	}
// }

// Reset clear user settings and restore Default to default.
func Reset() {
	ResetLevel()
	ResetFlags()
}

func SetDefault(l Logger) { defaultLog = l }    // sets user-defined logger as Default
func Default() Logger     { return defaultLog } // return native default logger

var (
	defaultWriter *dualWriter
	defaultLog    Logger // builtin logger as default device, see func init()
	onceInit      sync.Once
)

// const (
// 	STDLOG = "std"
// 	GOLOG  = "go/log"
// )

var (
	lvlCurrent          Level
	flags               = LstdFlags
	minimalMessageWidth = 36
	levelOutputWidth    = 3
)

const (
	clrNone  = color.NoColor
	clrBasic = color.FgLightMagenta

	darkGray  = color.FgDarkGray
	lightGray = color.FgLightGray
	cyan      = color.FgCyan
	hiRed     = color.FgLightRed
	red       = color.FgRed
	yellow    = color.FgYellow

	clrTimestamp    = color.FgGreen
	clrFuncName     = darkGray
	clrAttrKey      = darkGray
	clrAttrKeyBg    = clrNone
	clrLoggerName   = color.FgLightGray
	clrLoggerNameBg = clrNone
)

const (
	timestampFieldName = "time"
	levelFieldName     = "level"
	callerFieldName    = "caller"
	messageFieldName   = "msg"
)

// var errNotReady = errors.New("not ready") // here we just need a very simple error message object
