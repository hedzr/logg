package slog

import (
	"errors"
	"os"
	"regexp"
	"sync"

	"github.com/hedzr/is"
	"github.com/hedzr/is/term/color"
	"github.com/hedzr/logg/slog/internal/strings"
)

func A() int        { return severity() }
func severity() int { return 2 }

type X struct {
	IsAllBitsSet bool `json:"is-all-bits-set,omitempty"`
	BattleEnd    int  `json:"battle-end,omitempty"`
}

// AddCodeHostingProviders appends more provider, repl pair to reduce the caller text width.
// The builtin providers includes:
//   - "github.com" -> "GH"
//   - "gitlab.com" -> "GL"
//   - "gitee.com" -> "GT"
//   - "bitbucket.com" -> "BB"
//   - ...
func AddCodeHostingProviders(provider, repl string) { codeHostingProvidersMap[provider] = repl }

// AddKnownPathMapping appends more pathname, repl pair to reduce the caller filepath width.
// Such as:
//   - "$HOME" -> "~"
//   - pwd -> "." (current directory -> '.', that means any abspath will be converted to relpath)
func AddKnownPathMapping(pathname, repl string) { knownPathMap[pathname] = repl }

// AddKnownPathRegexpMapping adds regexp pattern, repl pair to reduce the called filepath width.
func AddKnownPathRegexpMapping(pathnameRegexpExpr, repl string) {
	knownPathRegexpMap = append(knownPathRegexpMap, regRepl{
		expr: regexp.MustCompile(pathnameRegexpExpr),
		repl: repl,
	})
}

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

func GetFlags() Flags           { return flags }
func SetFlags(f Flags)          { flags = f }
func AddFlags(f Flags)          { flags |= f }
func RemoveFlags(f Flags)       { flags &= ^f }
func ResetFlags()               { flags = LstdFlags }
func IsAnyBitsSet(f Flags) bool { return flags&f != 0 }
func IsAllBitsSet(f Flags) bool { return flags&f == f }

// SaveFlagsAndMod saves old flags, modify it, and restore the old at defer time
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

func Reset() {
	ResetLevel()
	ResetFlags()
}

func SetDefault(l Logger) { defaultLog = l }
func Default() Logger     { return defaultLog }

var (
	defaultWriter *dualWriter
	defaultLog    Logger // builtin logger as default device, see func init()
	onceInit      sync.Once
)

const (
	STDLOG = "std"
	GOLOG  = "go/log"
)

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

var (
	codeHostingProvidersMap map[string]string

	knownPathMap map[string]string

	knownPathRegexpMap []regRepl
)

type regRepl struct {
	expr *regexp.Regexp
	repl string
}

var homeDir, currDir string

var errNotReady = errors.New("not ready") // here we just need a very simple error message object

func init() {
	onceInit.Do(func() {
		lvlCurrent = WarnLevel

		codeHostingProvidersMap = map[string]string{
			"github.com":    "GH",
			"gitlab.com":    "GL",
			"gitee.com":     "GT",
			"bitbucket.com": "BB",
		}

		homeDir, _ = os.UserHomeDir()
		currDir, _ = os.Getwd()
		knownPathMap = map[string]string{
			homeDir: "~",
			currDir: ".",
		}
		knownPathRegexpMap = append(knownPathRegexpMap, regRepl{
			expr: regexp.MustCompile(`/Volumes/[^/]+/`),
			repl: `~`,
		})

		if is.InDebugging() || is.InTesting() || strings.StringToBool(os.Getenv("DEBUG")) {
			lvlCurrent = DebugLevel
		}

		defaultWriter = newDualWriter()
		defaultLog = newDetachedLogger()
	})
}
