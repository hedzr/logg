# hedzr/logg

![Go](https://github.com/hedzr/logg/workflows/Go/badge.svg)
[![GitHub tag (latest SemVer)](https://img.shields.io/github/tag/hedzr/logg.svg?label=release)](https://github.com/hedzr/logg/releases)
[![Go Dev](https://img.shields.io/badge/go-dev-green)](https://pkg.go.dev/github.com/hedzr/logg) <!--
[![Go Report Card](https://goreportcard.com/badge/github.com/hedzr/logg)](https://goreportcard.com/report/github.com/hedzr/logg)
[![Coverage Status](https://coveralls.io/repos/github/hedzr/logg/badge.svg?branch=master&.9)](https://coveralls.io/github/hedzr/logg?branch=master)
--> [![deps.dev](https://img.shields.io/badge/deps-dev-green)](https://deps.dev/go/github.com%2Fhedzr%2Flogg)

A golang logging library, to provide colorful output at terminal.

## Features

It is pre-releasing currently. Some abilities are:

- fast enough: performance is not our unique aim, and this one is enough quick.
- colorful console output by default.
- switch to logfmt or json format dynamically.
- interfaces and abilities similar with log/slog.
- adapted into log/slog to enable color logging, here some our verbs (such as Fatal, Panic) cannot work directly.
- cascade child logger and dump attrs of parent recursively (need enable `LattrsR` to avoid taking more cpu usages).
- very lite child loggers.
- user-defined levels, writer, and value stringer.
- privacy enough: harden filepath, shorten package name (by `Lprivacypath` and `Lprivacypathregexp`); and implement `LogObjectMarshaller` or `LogArrayMarshaller` to security sensitive fields.
- better multiline outputs.

![image-20231107091609707](https://cdn.jsdelivr.net/gh/hzimg/blog-pics@master/uPic/image-20231107091609707.png)

See also [CHANGELOG](CHANGELOG).

> Since v0.5.7, `logg/slog` enables privacy hardening flags by default now.

## Motivation

As an opt-in copy of `log/slog`, we provide an out-of-box colored text outputting logger with more verbs.

And an auto-optimized Verb: `Verbose(msg, args...)`
would print logging contents only if build tag `verbose` defined.
The point is it will be optimized completely in a default build.

## Guide

### Basics

`logg/slog` provides package-level functions for logging, `Info`, `Error`, and so on They will be mapped to the builtin default logger.

`Default()` returns the default logger and `SetDefult()` replaces it.

The basic usages are:

```go
import "github.com/hedzr/logg/slog"

msg := "A message"
args := []any{
    "attr1", 0,
    "attr2", false,
    "attr3", 3.13,
    "attr4", errors.New("simple"),
    // ,,,
    slog.Group("group1",
        "attr1", 0,
        "attr2", false,
        slog.NewAttr("attr3", "any styles what u prefer"),
        "attrn", // unpaired key can work here 
        slog.Group("more", "group", []byte("more subgroup here")),
    ),
    // ...
    slog.Int("id", 23123),
    slog.Group("properties",
        slog.Int("width", 4000),
        slog.Int("height", 3000),
        slog.String("format", "jpeg"),
    ),
}

slog.Print("")    // logging a clean newline without decorations
slog.Println("")  // logging a clean newline without decorations
slog.Println()    // logging a clean newline without decorations
slog.Print(msg, args...)
slog.Println(msg) // synosym of Print
slog.Fatal(msg, args...)
slog.Panic(msg, args...)
slog.Error(msg, args...)
slog.Warn(msg, args...)
slog.Info(msg, args...)
slog.Debug(msg, args...)
slog.Trace(msg, args...)

// only print the logging contents while built with `-tags verbose`
slog.Verbose(msg, args...) // 

// some verbs with more meanings
slog.OK("ok")
slog.OK("ok", args...)
slog.Success(msg, args...)
slog.Fail(msg, args...)

// Contextual logging
slog.InfoContext(context.TODO(), "info msg", args...)

// No-featured predicate here
slog.LogAttrs(context.TODO(), slog.Debug, "debug msg", args...)
```

The sub-loggers are also supported, see [SubLogger](#sublogger).

#### Println

A `Println(args ...any)` has slight differences with other verbs like `Print(msg string, args ...any)`. But you are using it with same form like others. That is, the first of the args passing to Println is indeed a msg string. Just a little benefit, you can pass nothing to Println. If you're doing with this way, a complete blank line prints by `slog.Println()`, no timestamp, no serverity, and no caller info.

```go
slog.Info("1")
slog.Println() // this makes a real blank line, for colorful and logfmt formats
slog.Info("2")
```

In a large long logging outputs, one and more blank line(s) maybe help your focus.

### Customizing Your Verbs

Like `OK`, `Success` and `Fail`, you could wrap logg/slogg package-level predicates/verbs for giving more features.

```go
func Tip(msg string, args...any) {
  if myConditionSatisfied {
    slog.WithSkip(1).OK(msg, args...)
  }
}
```

`WithSkip(n)` tells log/slog to strip extra 1 level(s) from caller stack frames so the logged message can hook to the caller of Tip(), rather than Line 3 in Tip().

Also, using a custom level is a not-bad idea. See the [Customizing the Level](#customizing-the-level).

By creating and managing a sublogger, make your own logger might be dead simple:

```go
func newMyLogger2() *mylogger2 {
    l := slog.New("mylogger").WithLevel(slog.InfoLevel)
    s := &mylogger2{
        l, // Provides basic Logger interface such as Info, Debug, etc.
    true, // enable Infof()
        l.WithSkip(1), // A sublogger created here. specially for Infof
    }
    return s
}

type mylogger2 struct {
    slog.Logger
    SprintfLikeLoggingIsEnabled bool
    sl                          slog.Logger
}

func (s *mylogger2) Close() { s.Logger.Close() }

func (s *mylogger2) Infof(msg string, args ...any) {
    if s.SprintfLikeLoggingIsEnabled {
        if len(args) > 0 {
            // msg = fmt.Sprintf(msg, args...)

            var data []byte
            data = fmt.Appendf(data, msg, args...)
            s.sl.Info(string(data))
        }
        s.sl.Info(msg)
    }
}

func TestSlogBasic4(t *testing.T) {
    l := newMyLogger2()
    l.Infof("what's wrong with %v", "him")
    l.Info("no matter")
    // l.Infof("what's wrong with %v, %v, %v, %v, %v, %v", m1AttrsAsAnySlice()...)
}
```

That is it.

### Builtin Output Formats

logg/slog has tree builtin optput formats: logfmt, json and colorful mode.

The default output is colorful to fit for debug console. But you can switch to the other two easily:

```go
slog.WithJSONMode()       // to get JSON format
slog.WithColorMode(false) // to get logfmt format
slog.WithColorMode()      // return to colorful mode
```

The above settings modify and apply effects to all loggers globally.

`logg/slog` has no way to modify formats for certain a logger or sub-logger. We do believe that's normal action to keep a uniform output format.

The outputs:

```bash
{"time":"14:53:10.907238+08:00","level":"debug","msg":"Debug message","source":{"function":"github.com/hedzr/logg/slog_test.TestSlogJSON","file":"./i_test.go","line":42}}

time="14:52:50.083343+08:00" level="debug" msg="Debug message" source.function="github.com/hedzr/logg/slog_test.TestSlogLogfmt" source.file="./i_test.go" source.line=68

```

and,

![image-20231028145416343](https://cdn.jsdelivr.net/gh/hzimg/blog-pics@master/uPic/image-20231028145416343.png)

### Set Level

The default is `WarnLevel` for a released app. If debugger detected the `DebugLevel` will be assumed.

Also the executable path will be tested for looking up if runing in test mode.

```go
slog.SetLevel(slog.InfoLevel)
println(slog.Level)
```

Each sub-logger can hold its own level different with their parents or default logger.

To restore old level on Default logger, `SaveLevelAndSet` is available:

```go
func TestSlogBasic2(t *testing.T) {
    defer slog.SaveLevelAndSet(slog.TraceLevel)()
    slog.Debug("Debug message") // should be enabled now
    slog.Info("Info message")   // should be enabled now
    slog.Warn("Warning message")
    slog.Error("Error message")
}
```

### Sublogger

Creating a detached logger is possible. Different from default logger, the logger's methods are used to log your message.

```go
logger := slog.New() // colorful logger
logger := slog.New().WithJSONMode() // json format logger
logger := slog.New().WithColorMode(false) // logfmt logger
logger := slog.New().With("attr1", v1, "attr2", v2))

logger := slog.New("name")
logger := slog.New("name", slog.WithAttrs(args...))
logger := slog.New("name", slog.NewAttr("attr1", v1))
logger := slog.New("name", slog.Int("attr1", i1))
logger := slog.New("name", slog.Group("group1", slog.Int("attr1", i1)))
logger := slog.New("name", "attr1", v1, "attr2", v2).WithAttrs(args...)
```

Sublogger should derive from Default() or a detached logger.

```go
logger := slog.New(args...).WithLevel(slog.InfoLevel)

sl1 := logger.New().WithLevel(slog.TraceLevel)
sl2 := Default().New() // keep the parent's level
```

By default, parent shares his features (level and other settings) to children, so `sl2` get `InfoLevel` same with `logger`.

If `LattrsR` is set, the parent's attributes will be inherited to. For performance reason, it's unset by default.

```go
logger := slog.New("parent-logger").With("attr", "parent").WithLevel(slog.InfoLevel)
sl := logger.New("child-logger")

slog.AddFlags(slog.LattrsR)
sl.Info("info", "attr1", 1)
slog.RemoveFlags(slog.LattrsR)
sl.Info("info", "attr1", 1)
```

The outputs looks like:

```bash
13:35:31.524263+08:00 child-logger [INF] info                                                    attr=parent attr1=1 /Volumes/VolHack/work/godev/cmdr.v2/libs.logg/slog/i_test.go:454 slog_test.TestSlogAttrsR
13:35:31.524276+08:00 child-logger [INF] info                                                    attr1=1 /Volumes/VolHack/work/godev/cmdr.v2/libs.logg/slog/i_test.go:456 slog_test.TestSlogAttrsR
```

A logger or a sublogger could be identify by a unique name. Passing a string as first parameter to `New(...)`, it's assumed the logger name.

And attributes and WithOpts can follow the logger name. `New(...)` parses all of them and process them. For examples:

```go
    l := slog.New("standalone-logger-for-app",
        slog.NewAttr("attr1", 2),
        slog.NewAttrs("attr2", 3, "attr3", 4.1),
        "attr4", true, "attr3", "string",
        slog.WithLevel(slog.AlwaysLevel),
    )
    defer l.Close()

    sub1 := l.New("sub1").With("logger", "sub1")
    sub2 := l.New("sub2").With("logger", "sub2").WithLevel(slog.InfoLevel)

    sub1.Debug("hi debug", "AA", 1.23456789)
```

Making children loggers is possible.

### Logging with Attributes

The attributes can be prepared or passed in several forms.

#### Plain form

```go
    logger.Info("info message",
        "attr1", 0,
        "attr2", false,
        "attr3", 3.13,
    ) // plain key and value pairs
```

#### NewAttr and NewGroupedAttr

```go
logger.Info("info message",
        "attr1", 0,
        "attr4", errors.New("simple"),
        "attr3", 3.13,
        slog.NewAttr("attr3", "any styles what u prefer"),
        slog.NewGroupedAttrEasy("group1", "attr1", 13, "attr2", false),
    ) // use NewAttr, NewGroupedAttrs
```

#### Int, Float, String, Any, ..., Group

logg/slog supports strong typed attributes:

```go
logger.Info("image uploaded",
        slog.Int("id", 23123),
        slog.Group("properties",
            slog.Int("width", 4000),
            slog.Int("height", 3000),
            slog.String("format", "jpeg"),
        ),
    ) // use Int, Float, String, Any, ..., and Group
```

These interfaces are very similar with standard log/slog.

#### Mixes all above forms

The above forms can be mixed in any order together.

```go
logger.Info("image uploaded",
        "attr1", 0,
        "attr2", false,
        "attr3", 3.13,
        "attr4", errors.New("simple"),
        // ,,,
        slog.Group("group1", 
            "attr1", 0,
            "attr2", false,
            slog.NewAttr("attr3", "any styles what u prefer"),
            slog.Group("more", "group", []byte("more sub-attrs")),
            "attrN", // unpaired key can work here
        ),
        // ...
        slog.Int("id", 23123),
        slog.Group("properties",
            slog.Int("width", 4000),
            slog.Int("height", 3000),
            slog.String("format", "jpeg"),
        ),
    )
```

#### Work with common Attributes

While creating a sublogger, you could specify some common attributes.
They are no more effects for performance reason by default.
`log.Info` and others printers will check out and print all of
parents' common attributes while `LattrsR` is set.

Both of the following forms are valid:

```go
// available forms
logger := slog.New("logger-name")
logger := slog.New("logger-name", slog.WithAttrs(args...))
logger := slog.New("logger-name", slog.NewAttr("attr1", v1))
logger := slog.New("logger-name", slog.Int("attr1", i1))
logger := slog.New("logger-name", slog.Group("group1", slog.Int("attr1", i1)))
logger := slog.New("logger-name", "attr1", v1, "attr2", v2).WithAttrs(args...)

// the attributes will be printed out if LattrsR set,
slog.AddFlags(slog.LattrsR)
logger.Info("message", "type", "int") // Out: ,,, A message    attr1=v1 attr2=v2 ... type=int
```

#### Grouping attributes

See above of above.

### Logging contextual attrs

Same to standard `log/slog`, `logg/slog` has LogAttrs() to log attributes contextually.

```go
    logger := slog.New().WithAttrs(slog.String("app-version", "v0.0.1-beta"))
    ctx := context.Background()
    logger.InfoContext(ctx, "info msg", "attr1", 111333,
        slog.Group("memory",
            slog.Int("current", 50),
            slog.Int("min", 20),
            slog.Int("max", 80)),
        slog.Int("cpu", 10),
    )
```

#### Extracting attrs from context

Sometimes the attributes can be extracted from context.Context.

```go
func TestSlogWithContext(t *testing.T) {
    logger := slog.New().WithAttrs(slog.String("app-version", "v0.0.1-beta"))
    ctx := context.WithValue(context.Background(), "ctx", "oh,oh,oh")
    logger.WithContextKeys("ctx").InfoContext(ctx, "info msg",
        "attr1", 111333,
        slog.Group("memory",
            slog.Int("current", 50),
            slog.Int("min", 20),
            slog.Int("max", 80)),
        slog.Int("cpu", 10),
    )
}
```

The result:

![image-20231106074712683](https://cdn.jsdelivr.net/gh/hzimg/blog-pics@master/uPic/image-20231106074712683.png)

As you seen, the value in context was been extracted and printed out.

### Set Writer

logg/slog uses a `dualWriter` to serialize the logging contents.

A `dualWriter` sends contents to stdout or stderr in accord to the requesting logging level. For example, a `Info(...)` calling will be dispatched to stdout and a `Warn`, `Error`, `Panic`, or `Fatal` to stderr.

Not only for those, the dualWriter allows you stack many writers as its `Normal` or `Error` output devices. That means, a console `os.Stdout` and a file writer can be bundled into `Normal` at once. How to do it? Simple:

```go
logger := New("tty+file").AddWriter(slog.NewFileWriter("/tmp/app-stdout.log"))
```

Using `WithWriter` to replace our stocked version, which is stdout+stderr by default.

Or `AddErrorWrite(w)` can append to the `Error` device.

Of course, `ResetWriters` works so we can always back to default state.

#### Leveled Writers

Also we provides sub-feature to allow you specify special writer for a special logging level:

```go
func TestAddLevelWriter1(t *testing.T) {
    logger := New().AddLevelWriter(InfoLevel, &decorated{os.Stdout})
    logger.Info("info msg")
    logger.Debug(getMessage(0))
}

type decorated struct {
    *os.File
}

func (s *decorated) Write(p []byte) (n int, err error) {
    if s.File != nil {
        if ni, e := s.File.WriteString("[decorated] "); e != nil {
            err = errors.Join(err, e)
        } else {
            n += ni
        }
        if ni, e := s.File.Write(p); e != nil {
            err = errors.Join(err, e)
        } else {
            n += ni
        }
    }
    return
}
```

The result is similar with:

```bash
[decorated] 10:57:22.986797+08:00 [INF] info msg                              ./new_test.go:170 slog.TestAddLevelWriter1
10:57:22.987031+08:00 [DBG] Test logging, but use a somewhat realistic message length. (#0)  ./new_test.go:171 slog.TestAddLevelWriter1

```

#### Close()

In most cases, Logger's `Close()` has no more attentions to you.

But adding defer close codes is a best practice, like this:

```go
package main

import "github.com/hedzr/logg/slog"

func main() {
  logger := slog.New("tty+file").AddWriter(slog.NewFileWriter("/tmp/app-stdout.log"))
  defer logger.Close()
  
  logger.Info("info msg")
}
```

The file writer will get a chance to shutdown itself gracefully.

### Set Handler

### Customizing the Level

In logg/slog, using your own logging level is enough simple:

```go
const (
    NoticeLevel = slog.Level(17) // A custom level must have a value greater than slog.MaxLevel
    HintLevel   = slog.Level(-8) // Or use a negative number
    SwellLevel  = slog.Level(12) // Sometimes, you may use the value equal with slog.MaxLevel
)

func TestSlogCustomizedLevel(t *testing.T) {
    checkerr(t, slog.RegisterLevel(NoticeLevel, "NOTICE",
        slog.RegWithShortTags([6]string{"", "N", "NT", "NTC", "NOTC", "NOTIC"}),
        slog.RegWithColor(color.FgWhite, color.BgUnderline),
        slog.RegWithTreatedAsLevel(slog.InfoLevel),
    ))

    checkerr(t, slog.RegisterLevel(HintLevel, "Hint",
        slog.RegWithShortTags([6]string{"", "H", "HT", "HNT", "HINT", "HINT "}),
        slog.RegWithColor(color.NoColor, color.BgInverse),
        slog.RegWithTreatedAsLevel(slog.InfoLevel),
    ))

    checkerr(t, slog.RegisterLevel(SwellLevel, "SWELL",
        slog.RegWithShortTags([6]string{"", "S", "SW", "SWL", "SWEL", "SWEEL"}),
        slog.RegWithColor(color.FgRed, color.BgBoldOrBright),
        slog.RegWithTreatedAsLevel(slog.ErrorLevel),
        slog.RegWithPrintToErrorDevice(),
    ))

    logger := slog.New()

    logger.Debug("Debug message")
    logger.Info("Info message")
    logger.Warn("Warning message")
    logger.Error("Error message")

    slog.SetLevelOutputWidth(5)

    ctx := context.Background()
    logger.LogAttrs(ctx, NoticeLevel, "Notice message")
    logger.LogAttrs(ctx, HintLevel, "Hint message")
    logger.LogAttrs(ctx, SwellLevel, "Swell level")
}

func checkerr(t *testing.T, err error) {
    if err != nil {
        t.Error(err)
    }
}
```

Its outputs looks like

![image-20231030144425448](https://cdn.jsdelivr.net/gh/hzimg/blog-pics@master/uPic/image-20231030144425448.png)

The ansi color representations relyes on your terminal settings.

`SetLevelOutputWidth(n)` lets you can control the level serverity's display width (from 1 to 5). At customizing you should pass a array (`[6]string`) for each levels. For example, HintLevel can be displayed as "HINT" when output width is 4, or as "H" when width is 1. You could change it globally at any time. Just like the snapshot above, it was changed to 5 at last time, so HintLevel has width 5.

### Customizing the Colors

Dislike logg/slog's console colors of each levels? No matter, setup with yours:

```go
import color "github.com/hedzr/is/term/color"

slog.SetLevelColors(slog.PanicLevel, color.FgWhite, color.NoColor)
```

These are standard ANSI escaped sequences, any of SGR fore-, background colors. A special `NoColor` is `-1`.

```go
const NoColor = color.Color(-1)
```

[`hedzr/is`](https://github.com/hedzr/is) provides a color Translator to format your html-like string to ansi colored text for terminal outputting. For more information see hedzr/is doc.

### Adapt into `log/slog`

If you are using unified `log/slog` interfaces, put logg/slog into it:

```go
import "log/slog"
import logslog "github.com/hedzr/logg/slog"

func TestSlogUsedForLogSlog(t *testing.T) {
    l := slog.New("standalone-logger-for-app",
        slog.NewAttr("attr1", 2),
        slog.NewAttrs("attr2", 3, "attr3", 4.1),
        "attr4", true, "attr3", "string",
        slog.WithLevel(slog.AlwaysLevel),
    )
    defer l.Close()

    sub1 := l.New("sub1").With("logger", "sub1")
    sub2 := l.New("sub2").With("logger", "sub2").WithLevel(slog.InfoLevel)

    // create a log/slog logger HERE
    logger := logslog.New(slog.NewSlogHandler(l, nil))

    t.Logf("logger: %v", logger)
    t.Logf("l: %v", l)

    // and logging with log/slog
    logger.Debug("hi debug", "AA", 1.23456789)
    logger.Info(
        "incoming request",
        logslog.String("method", "GET"),
        logslog.String("path", "/api/user"),
        logslog.Int("status", 200),
    )

    // now using our logg/slog interface
    sub1.Debug("hi debug", "AA", 1.23456789)
    sub2.Debug("hi debug", "AA", 1.23456789)
    sub2.Info("hi info", "AA", 1.23456789)
}
```

The result:

![image-20231029100643178](https://cdn.jsdelivr.net/gh/hzimg/blog-pics@master/uPic/image-20231029100643178.png)

Our testing keeps logg/slog native outputs as the last three lines.

In this case, some features of our logg/slog cannot be used via log/slog APIs but it's still colorful.

### Customizing `ValueStringer`

logg/slog allows you handle stringerize value with customizing `ValueStringer`. So you can pass yours.

Here is a sample to pretty print the values in attributes:

```go
import "github.com/alecthomas/repr"

func NewSpewPrinter() *prettyPrinter { //nolint:revive // just a test
    return &prettyPrinter{
        repr.New(os.Stdout, repr.Indent("  ")),
    }
}

type prettyPrinter struct {
    *repr.Printer
}

func (p *prettyPrinter) SetWriter(w io.Writer) {
    p.Printer = repr.New(w, repr.Indent("  "))
}

func (p *prettyPrinter) WriteValue(value any) {
    // p.reprValue(map[reflect.Value]bool{}, reflect.ValueOf(value), "", true, false)
    p.Print(value)
}

func TestSlogCustomValueStringer(t *testing.T) {
    // slog.SetLevel(slog.InfoLevel)
    slog.AddFlags(slog.Lprivacypathregexp | slog.Lprivacypath)

    defer slog.SaveFlagsAnd(func() { //nolint:revive // ok
        slog.AddFlags(slog.Lprettyprint)
    })()

    printer := NewSpewPrinter()

    for _, logger := range []slog.Logger{
        slog.New(" spew ").WithValueStringer(printer),
        slog.New("normal"),
    } {
        logger.Println()
        logger.LogAttrs(
            context.Background(),
            slog.InfoLevel,
            "image uploaded",
            slog.Int("id", 23123),
            slog.Group("properties",
                slog.Int("width", 4000),
                slog.Int("height", 3000),
                slog.String("format", "jpeg"),
                slog.Any("Map", map[int][]float64{3: {3.14, 2.72}, 5: {0.717, 1.732}}),
            ),
        )
    }
}
```

The outputs are:

![image-20231029170419385](https://cdn.jsdelivr.net/gh/hzimg/blog-pics@master/uPic/image-20231029170419385.png)

To integrete `go-spew` is similar and simple.



### Hide the sensitive fields

Your struct can implement `LogObjectMashaller` or `LogArrayMashaller` so that the sensitive fields can be harden.

```go
type users []*user

func (uu users) MarshalLogArray(enc *slog.PrintCtx) (err error) {
    for i := range uu {
        if i > 0 {
            enc.WriteRune(',')
        }
        if e := uu[i].MarshalLogObject(enc); e != nil {
            err = errors.Join(err, e)
        }
    }
    return
}

type user struct {
    Name      string    `json:"name"`
    Email     string    `json:"email"`
    CreatedAt time.Time `json:"created_at"`
}

func (u *user) MarshalLogObject(enc *slog.PrintCtx) (err error) {
    enc.AddString("name", u.Name)
    enc.AddRune(',')
    enc.AddString("email", u.Email)
    enc.AddRune(',')
    enc.AddInt64("createdAt", u.CreatedAt.UnixNano())
    return
}
```

`*slog.PrintCtx` is our value encoder.



### Harden filepath, shorten package name

The caller information could leak the user's name, disk volumes, directory structure and others sensitive contents.

#### Harden filepath

For this test case:

```go
func TestLogOneTwoThree(t *testing.T) {
    l := slogg.New().WithLevel(slogg.InfoLevel)
    l.Info("info msg", "Aa", 1, "Bbb", "a string", "Cc", 3.732, "D", 2.71828+5.3571i)
}
```

It may print out:

```go
20:20:01.697577+08:00 [INF] info msg                             Aa=1 Bbb=a string Cc=3.732 D=(2.71828+5.3571i) /Volumes/VolHack/work/go.work/libs.log/bench/logg_test.go:17 bench.TestLogOneTwoThree
```

Now we enable privacy flags, a builtin regexp rule (`/Volumes/.*/(.*)` -> `~$1`) will take effects:

`````
20:22:33.688847+08:00 [INF] info msg                             Aa=1 Bbb=a string Cc=3.732 D=(2.71828+5.3571i) ~work/go.work/libs.log/bench/logg_test.go:17 bench.TestLogOneTwoThree
`````

The codes looks like:

```go
func init() {
    slog.AddFlags(slog.Lprivacypathregexp | slog.Lprivacypath)
}
```

The builtin rules includes truncate homdir to `~`, disable absolute pathname, and using relative path, and so on.

You may make calls to `AddKnownPathMapping(path,repl)` and `AddKnownPathRegexpMapping(expr, repl)` to setup
them.



#### Shorten package name

In outputs, package name in caller information has form `github.com/user/repo/offset.object.function`, By enabling `Lcallerpackagename`, `github.com` will be shortened to `GH`. The other well-known code-hosting providers are converted too:

- "github.com" -> "GH"
- "gitlab.com" -> "GL"
- "gitee.com" -> "GT"
- "bitbucker.com" -> "BB"

```go
func init() {
    slog.AddFlags(slog.Lcallerpackagename)
}
```

If `Lcallerpackagename` is not present (this is default behavior), the package name will be truncated simply.

`AddCodeHostingProviders(provide, repl)` API can add more rules for shortening. 



#### More Rules

You can always append yours with `AddKnownPathMapping(pathname, repl string)` and `AddKnownPathRegexpMapping(pathnameRegexpExpr, repl string)`.

If `Lprivacypathregexp` and `Lprivacypath` is not present (this is default behavior), we try to truncate the pathname as possible as we can.

And, `AddCodeHostingProviders(provider, repl string)` do similar things and need `Lcallerpackagename` is enabled.



### Other Helpers

Here are two saver so that you can write codes easiler:

```go
// Save global Level and set to new, and restore original
defer slog.SaveLevelAndSet(slog.WarnLevel)

// Save global Flags and modify it for local logic, and restore it after going back to up level
defer slog.SaveFlagsAnd(func() {
    slog.AddFlags(slog.LattrsR) // add, remove, or set flags
})()

```



## Others

### What meaning the Multiline print?

In colorful print mode, one logging output starts with timestamp and serverity, and ends with attributes and caller info. So the message is put at center of one line with limited width.

A very long message will be wrapped to next line(s), right?

Thinking about another case, a logging line has title and details, what form output is better.

```go
func TestLogLongLine(t *testing.T) {
    l := New()
    l.Info("/ERROR/ Started at "+time.Now().String()+", this is a multi-line test\nVersion: 2.0.1\nPackage: hedzr/hilog", "Aa", 1, "Bbb", "a string", "Cc", 3.732, "D", 2.71828+5.3571i)
}
```

It prints:

![image-20231106234956844](https://cdn.jsdelivr.net/gh/hzimg/blog-pics@master/uPic/image-20231106234956844.png)

As you seen, a message will be splitted to first line (as a title) and rest lines (as a details text). So you can avoid write a long message. Instead, you write a title with some details to describe a logging line better clearly. By the way, its json or logfmt format is still a message field with key name "msg".



## LICENSE

Apache 2.0
