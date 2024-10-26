package slog

import (
	"context"
	"fmt"
	logslog "log/slog"
	"strings"

	"github.com/hedzr/is/states"
	"github.com/hedzr/is/term/color"
)

type Level int // logging level

// LevelSettable interface can be used when a user-defined writer
// needs update its internal level lively before writing.
type LevelSettable interface {
	SetLevel(level Level)
}

// SaveLevelAndSet sets Default logger's level and default logger level.
//
// SaveLevelAndSet saves old level and return a functor to restore it. So
// a typical usage is:
//
//	func A() {
//	    defer slog.SaveLevelAndSet(slog.PanicLevel)()
//	    l := slog.New().SetLevel(slog.WarnLevel)
//	    ...
//	}
//
// GetLevel and SetLevel can access the default logger level.
func SaveLevelAndSet(lvl Level) func() {
	save := GetLevel()
	SetLevel(lvl)
	return func() {
		SetLevel(save)
	}
}

// GetLevel returns the default logger level
func GetLevel() Level { return lvlCurrent }

// SetLevel sets the default logger level
func SetLevel(lvl Level) {
	lvlCurrent = lvl
	defaultLog.SetLevel(lvl)
	// if l, ok := defaultLog.(interface{ SetLevel(lvl Level) *Entry }); ok {
	// 	l.SetLevel(lvl)
	// } else if l, ok := defaultLog.(interface{ SetLevel(lvl Level) }); ok {
	// 	l.SetLevel(lvl)
	// } else if l, ok := defaultLog.(interface{ WithLevel(lvl Level) *Entry }); ok {
	// 	l.WithLevel(lvl)
	// }
}

// ResetLevel restore the default logger level to factory value (WarnLevel).
func ResetLevel() { SetLevel(WarnLevel) }

// RegisterLevel records the given Level value as stocked.
func RegisterLevel(levelValue Level, title string, opts ...RegOpt) error {
	for _, v := range allLevels {
		if int(v) == int(levelValue) {
			// return errorsv2.New("the given level %q is duplicated with %q", shortTags[5], v)
			return fmt.Errorf("the given level %q is duplicated with %q", shortTagMap[5][v], v)
		}
	}
	if l, ok := stringToLevel[title]; ok {
		// return errorsv2.New("the title %q has been used for %q", title, l)
		return fmt.Errorf("the title %q has been used for %q", title, l)
	}

	var pack = regPack{clr: color.NoColor,
		bg:      color.NoColor,
		treatAs: MaxLevel,
	}

	for _, opt := range opts {
		opt(&pack)
	}

	allLevels = append(allLevels, levelValue)
	levelToString[levelValue] = title
	stringToLevel[title] = levelValue

	for i := 0; i < MaxLengthShortTag; i++ {
		if str := pack.shortTags[i]; str != "" {
			shortTagMap[i][levelValue] = str
		}
	}
	if pack.clr != color.NoColor {
		if pack.bg != color.NoColor {
			mLevelColors[levelValue] = []color.Color{pack.clr, pack.bg}
		} else {
			mLevelColors[levelValue] = []color.Color{pack.clr}
		}
	}
	if pack.treatAs < MaxLevel {
		mLevelIsEnabledAs[levelValue] = pack.treatAs
	}
	if pack.printOutToErrorDevice {
		mLevelUseErrorDevice[levelValue] = true
	}

	return nil
}

type regPack struct {
	shortTags             [MaxLengthShortTag]string
	clr, bg               color.Color
	treatAs               Level
	printOutToErrorDevice bool
}

type RegOpt func(pack *regPack) // used by RegisterLevel

// RegWithShortTags associates short tag with the new level
func RegWithShortTags(shortTags [MaxLengthShortTag]string) RegOpt {
	return func(pack *regPack) {
		pack.shortTags = shortTags
	}
}

// RegWithColor associates terminal color with the new level
func RegWithColor(clr color.Color, bgColor ...color.Color) RegOpt {
	return func(pack *regPack) {
		pack.clr = clr
		for _, bg := range bgColor {
			pack.bg = bg
		}
	}
}

// RegWithTreatedAsLevel associates the underlying level with the new level.
//
// It means the new level acts as treatAs level.
//
// For instance, see RegWithPrintToErrorDevice.
// After registered,
//
//	slog.Log(ctx, SwellLevel, "xxx")
//
// will redirect the logging line to stderr device just like ErrorLevel.
func RegWithTreatedAsLevel(treatAs Level) RegOpt {
	return func(pack *regPack) {
		pack.treatAs = treatAs
	}
}

// RegWithPrintToErrorDevice declares the logging text should be
// redirected to stderr device.
//
// For instance, you're declaring Swell Level, which should be
// treated as ErrorLevel, so the corresponding codes are:
//
//	const SwellLevel  = slog.Level(12) // Sometimes, you may use the value equal with slog.MaxLevel
//	slog.RegisterLevel(SwellLevel, "SWELL",
//	    slog.RegWithShortTags([6]string{"", "S", "SW", "SWL", "SWEL", "SWEEL"}),
//	    slog.RegWithColor(color.FgRed, color.BgBoldOrBright),
//	    slog.RegWithTreatedAsLevel(slog.ErrorLevel),
//	    slog.RegWithPrintToErrorDevice(),
//	)
//
// After registered, slog.Log(ctx, SwellLevel, "xxx") will redirect the logging
// line to stderr device just like ErrorLevel.
func RegWithPrintToErrorDevice(b ...bool) RegOpt {
	return func(pack *regPack) {
		for _, v := range b {
			pack.printOutToErrorDevice = v
		}
	}
}

//

//

//

// Enabled tests the requesting level is denied or allowed.
//
// 1. hold OffLevel: any request levels are denied
// 2. hold AlwaysLevel: any request levels are allowed
// 3. in testing/debugging mode, requesting DebugLevel are always allowed.
// 4. an internal table built by the RegWithTreatedAsLevel RegOpt
// in RegisterLevel() will be checked and mapped on the requesting level.
//
// Note that the levels are ordinal: PanicLevel (0) .. DebugLevel (5)
// .. OffLevel (7), AlwaysLevel (8), OKLevel (9) .. MaxLevel (dyn).
func (level Level) Enabled(ctx context.Context, testingLevel Level) bool {
	if level == OffLevel || testingLevel == OffLevel {
		return false
	}
	if level == AlwaysLevel || testingLevel == AlwaysLevel {
		return true
	}
	if states.Env().GetDebugMode() && testingLevel == DebugLevel {
		return true
	}
	if l, ok := mLevelIsEnabledAs[testingLevel]; ok {
		testingLevel = l
	}
	return level >= testingLevel
}

// Convert the Level to a string. eg. PanicLevel becomes "panic".
func (level Level) String() string {
	if b, err := level.MarshalText(); err == nil {
		return string(b)
	}
	return fmt.Sprintf("L#%d", int(level))
}

// ShortTag convert Level to a short tag string. eg. PanicLevel becomes "P".
func (level Level) ShortTag(length int) string {
	if length <= 0 || length >= MaxLengthShortTag {
		panic("invalid length. the valid range: 1-5")
	}
	if va, ok := shortTagMap[length]; ok {
		if ix, ok := va[level]; ok {
			return ix
		}
	}

	if t := level.String(); len(t) > 0 {
		switch l := len(t); {
		case l == length:
			return t
		case l < length:
			t = t + strings.Repeat(" ", length)
			fallthrough
		default:
			return t[:length]
		}
	}
	return strings.Repeat("?", length)
}

func (level *Level) UnmarshalJSON(text []byte) error {
	return level.UnmarshalText(text)
}

func (level Level) MarshalJSON() ([]byte, error) {
	b, err := level.MarshalText()
	return []byte(fmt.Sprintf("%q", string(b))), err
}

// UnmarshalText implements encoding.TextUnmarshaler.
func (level *Level) UnmarshalText(text []byte) error {
	l, err := ParseLevel(string(text))
	if err != nil {
		return err
	}

	*level = Level(l)

	return nil
}

// MarshalText convert Level to string and []byte.
//
// Available level names are:
//   - "disable"
//   - "fatal"
//   - "error"
//   - "warn"
//   - "info"
//   - "debug"
//   - ...
func (level Level) MarshalText() ([]byte, error) {
	if str, ok := levelToString[level]; ok {
		return []byte(str), nil
	}

	return nil, fmt.Errorf("not a valid error level %d", level)
}

// ParseLevel takes a string level and returns the Logrus log level constant.
func ParseLevel(lvl string) (Level, error) {
	if l, ok := stringToLevel[strings.ToLower(lvl)]; ok {
		return l, nil
	}

	defaultLog.Warn("unknown logging level string", "string-value", lvl)

	var l Level
	return l, fmt.Errorf("not a valid logging Level: %q", lvl)
}

// AllLevels is a constant exposing all logging levels
func AllLevels() []Level { return allLevels }

// SetLevelColors defines your fore-/background color for printing
// the leveled text to terminal.
func SetLevelColors(lvl Level, fg, bg color.Color) {
	mLevelColors[lvl] = []color.Color{fg, bg}
}

// mLevelIsEnabledAs is a replacement table of two levels.
//
// The given level is treated as another one (generally
// this means a builtin level, such as InfoLevel, ...).
//
// The feature is used for interpreting a user-defined
// Level to a internal builtin level so that we can
// judge its priority.
//
// What's the priorities?
//
// The builtin Level(s) are an integer so the smaller
// level has higher priority. If logger holding WarnLevel
// and user is requesting by Info(),
// since WarnLevel (3) < InfoLevel (4),
// so the lower priority requesting (for this case, a InfoLevel)
// will be revoked. It prints nothing to output device.
//
// But on the contrary, A Warn() will prints out for a
// logger holding InfoLevel level.
var mLevelIsEnabledAs = map[Level]Level{
	OKLevel:      InfoLevel,
	SuccessLevel: InfoLevel,
	FailLevel:    ErrorLevel,
}

// mLevelUseErrorDevice maps Level to boolean flag which decides
// if the contents of that level will be written to error device
// or not.
var mLevelUseErrorDevice = map[Level]bool{
	PanicLevel: true,
	FatalLevel: true,
	ErrorLevel: true,
	WarnLevel:  true,
	FailLevel:  true,
}

var mLevelColors = map[Level][]color.Color{
	PanicLevel:   {hiRed, clrNone},                    //
	FatalLevel:   {hiRed, clrNone},                    //
	ErrorLevel:   {red, clrNone},                      //
	WarnLevel:    {yellow, clrNone},                   //
	InfoLevel:    {cyan, clrNone},                     //
	DebugLevel:   {color.FgMagenta, clrNone},          //
	TraceLevel:   {yellow, color.BgDim},               //
	OffLevel:     {color.FgBlack, color.BgDim},        // never used.
	AlwaysLevel:  {lightGray, color.BgBlink},          // blink color relies on concrete terminal but it commonly takes no effect.
	OKLevel:      {color.FgLightCyan, color.BgBlink},  //
	SuccessLevel: {color.FgGreen, color.BgBlink},      //
	FailLevel:    {color.FgRed, color.BgBoldOrBright}, //
}

var shortTagMap = map[int]map[Level]string{
	0: {},
	1: {PanicLevel: "P", FatalLevel: "F", ErrorLevel: "E", WarnLevel: "W", InfoLevel: "I", DebugLevel: "D", TraceLevel: "T", OffLevel: " ", AlwaysLevel: "A", OKLevel: "o", SuccessLevel: "s", FailLevel: "f"},
	2: {PanicLevel: "PC", FatalLevel: "FL", ErrorLevel: "ER", WarnLevel: "WN", InfoLevel: "IF", DebugLevel: "DG", TraceLevel: "TC", OffLevel: "  ", AlwaysLevel: "AA", OKLevel: "OK", SuccessLevel: "SU", FailLevel: "FA"},
	3: {PanicLevel: "PNC", FatalLevel: "FTL", ErrorLevel: "ERR", WarnLevel: "WRN", InfoLevel: "INF", DebugLevel: "DBG", TraceLevel: "TRC", OffLevel: "   ", AlwaysLevel: " A ", OKLevel: " OK", SuccessLevel: "SUC", FailLevel: "FAI"},
	4: {PanicLevel: "PNIC", FatalLevel: "FTAL", ErrorLevel: "ERRO", WarnLevel: "WARN", InfoLevel: "INFO", DebugLevel: "DBUG", TraceLevel: "TRAC", OffLevel: "    ", AlwaysLevel: " AA ", OKLevel: " OK ", SuccessLevel: "SUCC", FailLevel: "FAIL"},
	5: {PanicLevel: "PANIC", FatalLevel: "FATAL", ErrorLevel: "ERROR", WarnLevel: "WARNI", InfoLevel: "INFOR", DebugLevel: "DEBUG", TraceLevel: "TRACE", OffLevel: "     ", AlwaysLevel: "  A  ", OKLevel: "  OK ", SuccessLevel: "SUCCS", FailLevel: " FAIL"},
}

const MaxLengthShortTag = 6 // Level string length while formatting and printing

var stringToLevel = map[string]Level{
	"fail":     FailLevel,
	"success":  SuccessLevel,
	"ok":       OKLevel,
	"always":   AlwaysLevel,
	"off":      OffLevel,
	"no":       OffLevel,
	"disabled": OffLevel,
	"trace":    TraceLevel,
	"debug":    DebugLevel,
	"devel":    DebugLevel,
	"dev":      DebugLevel,
	"develop":  DebugLevel,
	"info":     InfoLevel,
	"warn":     WarnLevel,
	"warning":  WarnLevel,
	"error":    ErrorLevel,
	"fatal":    FatalLevel,
	"panic":    PanicLevel,
}

var levelToString = map[Level]string{
	FailLevel:    "fail",
	SuccessLevel: "success",
	OKLevel:      "ok",
	AlwaysLevel:  "always",
	OffLevel:     "off",
	TraceLevel:   "trace",
	DebugLevel:   "debug",
	InfoLevel:    "info",
	WarnLevel:    "warning",
	ErrorLevel:   "error",
	FatalLevel:   "fatal",
	PanicLevel:   "panic",
}

var mLevelToLogSlog = map[Level]logslog.Level{
	PanicLevel:   logslog.LevelError, //
	FatalLevel:   logslog.LevelError, //
	ErrorLevel:   logslog.LevelError, //
	WarnLevel:    logslog.LevelWarn,  //
	InfoLevel:    logslog.LevelInfo,  //
	DebugLevel:   logslog.LevelDebug, //
	TraceLevel:   logslog.LevelDebug, //
	OffLevel:     logslog.LevelInfo,  //
	AlwaysLevel:  logslog.LevelInfo,  //
	OKLevel:      logslog.LevelInfo,  //
	SuccessLevel: logslog.LevelInfo,  //
	FailLevel:    logslog.LevelInfo,  //
}

var mLogSlogLevelToLevel = map[logslog.Level]Level{
	logslog.LevelDebug: DebugLevel, // Level = -4
	logslog.LevelInfo:  InfoLevel,  // Level = 0
	logslog.LevelWarn:  WarnLevel,  // Level = 4
	logslog.LevelError: ErrorLevel, // Level = 8
}

var allLevels = []Level{
	PanicLevel,
	FatalLevel,
	ErrorLevel,
	WarnLevel,
	InfoLevel,
	DebugLevel,
	TraceLevel,
	OffLevel,
	AlwaysLevel,
	OKLevel,
	SuccessLevel,
	FailLevel,
}

const (

	// PanicLevel level, highest level of severity. Logs and then calls panic with the
	// message passed to Debug, Info, ...
	PanicLevel Level = iota
	// FatalLevel level. Logs and then calls `os.Exit(-9)`. It will exit even if the
	// logging level is set to PanicLevel.
	FatalLevel
	// ErrorLevel level. Logs. Used for errors that should definitely be noted.
	// Commonly used for hooks to send errors to an error tracking service.
	ErrorLevel
	// WarnLevel level. Non-critical entries that deserve eyes.
	WarnLevel
	// InfoLevel level. General operational entries about what's going on inside the
	// application.
	InfoLevel
	// DebugLevel level. Usually only enabled when debugging. Very verbose logging.
	DebugLevel
	// TraceLevel level. Designates finer-grained informational events than the DebugLevel.
	TraceLevel

	// OffLevel level. The logger will be shutdown.
	OffLevel
	// AlwaysLevel level. Used for Print, Printf, Println, OK, Success and Fail (use ErrorLevel).
	AlwaysLevel

	OKLevel      // OKLevel for operation okay.
	SuccessLevel // SuccessLevel for operation successfully.
	FailLevel    // FailLevel for operation failed.

	MaxLevel // maximal level value for algorithm usage
)
