package slog

import (
	"encoding"
	"fmt"
	"sync/atomic"
	"time"

	"github.com/hedzr/is/term/color"
	errorsv3 "gopkg.in/hedzr/errors.v3"
)

type Painter interface {
	Colorful() bool

	Begin(pc *PrintCtx)
	End(pc *PrintCtx, newline bool)
	BeginArray(pc *PrintCtx)
	EndArray(pc *PrintCtx, newline bool)

	AppendColon(pc *PrintCtx)
	AppendComma(pc *PrintCtx)
	AppendTime(pc *PrintCtx, z time.Time)
	AppendTimeSlice(pc *PrintCtx, z []time.Time)
	AppendTimestamp(pc *PrintCtx, tm time.Time, layout string)
	AppendDuration(pc *PrintCtx, z time.Duration)
	AppendDurationSlice(pc *PrintCtx, z []time.Duration)
	Append(pc *PrintCtx, data []byte)
	AppendStringKeyPrefixed(pc *PrintCtx, str, prefix string)
	AppendStringKey(pc *PrintCtx, str string)
	AppendKey(pc *PrintCtx, key string, clr, bg color.Color)
	AppendError(pc *PrintCtx, err error)
	AppendErrorAfterPrinted(pc *PrintCtx, err error)

	AddTimestampField(pc *PrintCtx, tm time.Time)
	AddLoggerNameField(pc *PrintCtx, name string)
	AddSeverity(pc *PrintCtx, lvl Level)
	AddPCField(pc *PrintCtx, source *Source)
	AddMsgField(pc *PrintCtx, msg string)
	AddMsgFieldFirstLine(pc *PrintCtx, firstLine string)
	AddMsgFieldRestLines(pc *PrintCtx, restLines string, eol bool)

	AddPrefixedString(pc *PrintCtx, prefix, name string, value string)
	TryQuoteValue(pc *PrintCtx, val string)
	MarshalValue(pc *PrintCtx, val any) (handled bool, err error)
}

var _ Painter = (*colorfulPainter)(nil)
var _ Painter = (*logfmtPainter)(nil)
var _ Painter = (*jsonPainter)(nil)

type colorfulPainter struct{ colorful bool }

func (s *colorfulPainter) Colorful() bool { return s.colorful }

func (s *colorfulPainter) AddMsgField(pc *PrintCtx, msg string) {
	pc.AddString(messageFieldName, ct.translate(pc.msg))
}

func (s *logfmtPainter) AddMsgField(pc *PrintCtx, msg string) {
	pc.AddString(messageFieldName, pc.msg)
}

func (s *jsonPainter) AddMsgField(pc *PrintCtx, msg string) {
	pc.AddString(messageFieldName, pc.msg)
}

func (s *colorfulPainter) AddMsgFieldFirstLine(pc *PrintCtx, firstLine string) {
	if minimalMessageWidth > 0 {
		str := ct.rightPad(firstLine, " ", minimalMessageWidth)
		if s.colorful {
			str = ct.translate(str)
			_, _ = pc.WriteString(ct.wrapColorAndBg(str, pc.clr, pc.bg))
		} else {
			_, _ = pc.WriteString(str)
		}
	} else {
		if s.colorful {
			str := ct.translate(firstLine)
			_, _ = pc.WriteString(ct.wrapColorAndBg(str, pc.clr, pc.bg))
		} else {
			_, _ = pc.WriteString(firstLine)
		}
	}
}

func (s *colorfulPainter) AddMsgFieldRestLines(pc *PrintCtx, restLines string, eol bool) {
	if restLines != "" {
		pc.AppendByte('\n')
		pc.AppendString(ct.padFunc(pc.restLines, " ", 4, func(i int, line string) string {
			if s.colorful {
				return ct.wrapColorAndBg(line, pc.clr, pc.bg)
			} else {
				return line
			}
		}))
		if eol {
			pc.pcAppendByte('\n')
		}
	}
}

func (s *logfmtPainter) AddMsgFieldFirstLine(pc *PrintCtx, firstLine string) {
	// pc.AppendByte('{')
}

func (s *logfmtPainter) AddMsgFieldRestLines(pc *PrintCtx, restLines string, eol bool) {
	// pc.AppendByte('{')
}

func (s *jsonPainter) AddMsgFieldFirstLine(pc *PrintCtx, firstLine string) {
	// pc.AppendByte('{')
}

func (s *jsonPainter) AddMsgFieldRestLines(pc *PrintCtx, restLines string, eol bool) {
	// pc.AppendByte('{')
}

func (s *colorfulPainter) Begin(pc *PrintCtx) {
	// pc.AppendByte('{')
}

func (s *colorfulPainter) End(pc *PrintCtx, newline bool) {
	// pc.AppendByte('}')
	if newline {
		pc.AppendByte('\n')
	}
}

func (s *colorfulPainter) BeginArray(pc *PrintCtx) {
	// pc.AppendByte('{')
}

func (s *colorfulPainter) EndArray(pc *PrintCtx, newline bool) {
	// pc.AppendByte('}')
	if newline {
		pc.AppendByte('\n')
	}
}

func (s *colorfulPainter) AppendColon(pc *PrintCtx) {
	pc.AppendByte('=')

	// s.preCheck()
	// switch s.mode {
	// case ModeJSON:
	// 	s.pcAppendByte(':')
	// default:
	// 	s.pcAppendByte('=')
	// }

	// if s.jsonMode {
	// 	s.pcAppendByte(':')
	// } else {
	// 	s.pcAppendByte('=')
	// }
}

func (s *colorfulPainter) AppendComma(pc *PrintCtx) {
	pc.AppendByte(' ')

	// switch s.mode {
	// case ModeJSON:
	// 	s.pcAppendByte(',')
	// case ModeLogFmt:
	// 	s.pcAppendByte(' ')
	// }

	// if s.jsonMode {
	// 	s.pcAppendByte(',')
	// } else {
	// 	s.pcAppendByte(' ')
	// }
}

func (s *colorfulPainter) AppendTime(pc *PrintCtx, z time.Time) {
	const layout = time.RFC3339Nano
	pc.Buf(func(buf []byte) []byte {
		return z.AppendFormat(buf, layout)
	})

	// if s.mode != ModeColorful {
	// 	// if s.jsonMode || s.noColor {
	// 	s.pcAppendByte('"')
	// 	s.buf = z.AppendFormat(s.buf, layout)
	// 	s.pcAppendByte('"')
	// } else {
	// 	s.buf = z.AppendFormat(s.buf, layout)
	// }
}

func (s *colorfulPainter) AppendTimeSlice(pc *PrintCtx, z []time.Time) {
	pc.AppendByte('[')
	if l := len(z); l > 0 {
		s.AppendTime(pc, z[0])
		for i := 1; i < len(z); i++ {
			dur := z[i]
			pc.AppendByte(',')
			s.AppendTime(pc, dur)
		}
	}
	pc.AppendByte(']')
}

func (s *colorfulPainter) AppendTimestamp(pc *PrintCtx, tm time.Time, layout string) {
	pc.Buf(func(buf []byte) []byte {
		return tm.AppendFormat(buf, layout)
	})
	pc.AppendByte('|')
}

func (s *colorfulPainter) AppendDuration(pc *PrintCtx, z time.Duration) {
	// s.pcAppendByte('"')
	// s.appendEscapedJSONString(z.String())
	// s.pcAppendByte('"')
	pc.AppendQuotedString(z.String())
}

func (s *colorfulPainter) AppendDurationSlice(pc *PrintCtx, z []time.Duration) {
	pc.AppendByte('[')
	if l := len(z); l > 0 {
		s.AppendDuration(pc, z[0])
		for i := 1; i < len(z); i++ {
			dur := z[i]
			pc.AppendByte(',')
			s.AppendDuration(pc, dur)
		}
	}
	pc.AppendByte(']')
}

func (s *colorfulPainter) Append(pc *PrintCtx, data []byte) {
	pc.AppendBytes(data)
}

func (s *colorfulPainter) AppendStringKeyPrefixed(pc *PrintCtx, str, prefix string) {
	_, _ = pc.WriteString(prefix)
	_ = pc.WriteByte('.')
	_, _ = pc.WriteString(str)

	// switch s.mode {
	// case ModeJSON:
	// 	// s.WriteString(strconv.Quote(str))
	// 	// s.Grow(2 + len([]byte(str)))
	// 	s.checkerr(s.WriteByte('"'))
	// 	_, _ = s.WriteString(prefix)
	// 	s.checkerr(s.WriteByte('.'))
	// 	_, _ = s.WriteString(str)
	// 	s.checkerr(s.WriteByte('"'))
	// // case ModeLogFmt:
	// default:
	// 	_, _ = s.WriteString(prefix)
	// 	s.checkerr(s.WriteByte('.'))
	// 	_, _ = s.WriteString(str)
	// }

	// if s.jsonMode {
	// 	// s.WriteString(strconv.Quote(str))
	// 	// s.Grow(2 + len([]byte(str)))
	// 	s.checkerr(s.WriteByte('"'))
	// 	_, _ = s.WriteString(prefix)
	// 	s.checkerr(s.WriteByte('.'))
	// 	_, _ = s.WriteString(str)
	// 	s.checkerr(s.WriteByte('"'))
	// } else {
	// 	_, _ = s.WriteString(prefix)
	// 	s.checkerr(s.WriteByte('.'))
	// 	_, _ = s.WriteString(str)
	// }
}

func (s *colorfulPainter) AppendStringKey(pc *PrintCtx, str string) {
	_, _ = pc.WriteString(str)

	// s.preCheck()

	// switch s.mode {
	// case ModeJSON:
	// 	// s.WriteString(strconv.Quote(str))
	// 	// s.Grow(2 + len([]byte(str)))
	// 	s.checkerr(s.WriteByte('"'))
	// 	_, _ = s.WriteString(str)
	// 	s.checkerr(s.WriteByte('"'))
	// // case ModeLogFmt:
	// default:
	// 	_, _ = s.WriteString(str)
	// }

	// if s.jsonMode {
	// 	// s.WriteString(strconv.Quote(str))
	// 	// s.Grow(2 + len([]byte(str)))
	// 	s.checkerr(s.WriteByte('"'))
	// 	_, _ = s.WriteString(str)
	// 	s.checkerr(s.WriteByte('"'))
	// } else {
	// 	_, _ = s.WriteString(str)
	// }
}

func (s *colorfulPainter) AppendKey(pc *PrintCtx, key string, clr, bg color.Color) {
	if pc.IsColorfulStyle() {
		ct.echoColorAndBg(pc, clrAttrKey, clrAttrKeyBg)
		s.AppendStringKey(pc, key)
		ct.echoColorAndBg(pc, clr, bg)
	} else {
		s.AppendStringKey(pc, key)
	}
	s.AppendColon(pc)
}

func (s *colorfulPainter) AppendError(pc *PrintCtx, err error) {
	ct.echoColor(pc, clrError)
	s.TryQuoteValue(pc, err.Error())
	ct.echoResetColor(pc)
}

func (s *colorfulPainter) AppendErrorAfterPrinted(pc *PrintCtx, err error) {
	if err != nil && (inTesting || isDebuggingOrBuild || isDebug()) && !inBenching {
		// the following job must follow the normal line, so it can't be committed
		// at PrintCtx.appendValue.

		if f, st := pc.getStackTrace(err); st != nil {
			pc.AppendByte('\n')

			frame := st[0]
			pc.cachedSource.Extract(uintptr(frame))
			s.AppendStringKey(pc, "       error: ")
			if pc.IsColorfulStyle() {
				ct.wrapColorAndBgTo(pc, clrError, clrNone, f.Error())
			} else {
				pc.AppendString(f.Error())
			}
			pc.AppendByte('\n')
			s.AppendStringKey(pc, "   file/line: ")
			pc.AppendString(pc.cachedSource.File)
			pc.AppendRune(':')
			pc.AppendInt(pc.cachedSource.Line)
			pc.AppendByte('\n')
			s.AppendStringKey(pc, "    function: ")
			if pc.IsColorfulStyle() {
				ct.wrapColorAndBgTo(pc, clrFuncName, clrNone, pc.cachedSource.Function)
			} else {
				pc.AppendString(pc.cachedSource.Function)
			}
		}

		var stackInfo string
		stackInfo = fmt.Sprintf("%+v", err)
		// if _, ok := err.(interface{ Format(s fmt.State, verb rune) }); ok {
		// 	stackInfo = fmt.Sprintf("%+v", err)
		// } else {
		// 	var x Stringer
		// 	if x, ok = err.(Stringer); ok {
		// 		stackInfo = x.String() // special for those illegal error impl like toml.DecodeError, which have no Format
		// 	} else {
		// 		stackInfo = fmt.Sprintf("%v", err)
		// 	}
		// }
		pc.AppendByte('\n')
		txt := ct.pad(stackInfo, "    ", 1)
		if pc.IsColorfulStyle() {
			ct.wrapColorAndBgTo(pc, clrError, clrLoggerNameBg, txt)
		} else {
			pc.AppendString(txt)
		}
	}
}

func (s *colorfulPainter) AddTimestampField(pc *PrintCtx, tm time.Time) {
	if pc.IsColorfulStyle() {
		ct.echoColor(pc, clrTimestamp)
	}
	pc.AppendTimestamp(pc.now)
	pc.AppendByte(' ')
}

func (s *colorfulPainter) AddLoggerNameField(pc *PrintCtx, name string) {
	l, lnw := len(name), int(atomic.LoadInt32(&longestNameWidth))
	if r := lnw - l; r > 0 {
		pc.AppendRuneTimes(' ', r)
	}

	if s.colorful {
		ct.wrapColorAndBgTo(pc, clrLoggerName, clrLoggerNameBg, name)
	} else {
		pc.AppendString(name)
	}

	pc.AppendByte(' ')
}

func (s *colorfulPainter) AddSeverity(pc *PrintCtx, lvl Level) {
	if s.colorful {
		ct.wrapColorAndBgTo(pc, pc.clr, pc.bg, ct.wrapRune(lvl.ShortTag(levelOutputWidth), '[', ']'))
	} else {
		pc.AppendString(ct.wrapRune(lvl.ShortTag(levelOutputWidth), '[', ']'))
	}
	pc.AppendByte(' ')
}

func (s *colorfulPainter) AddPCField(pc *PrintCtx, source *Source) {
	pc.AppendByte(' ')
	// pc.appendRune('(')
	pc.AppendString(source.File)
	pc.AppendByte(':')
	pc.AppendInt(source.Line)
	// pc.appendRune(')')
	pc.AppendByte(' ')
	// ct.wrapDimColorTo(pc.SB, source.checkedfuncname()) // clion p-term in run panel cannot support dim color.
	if s.colorful {
		ct.wrapColorTo(pc, clrFuncName, checkedfuncname(source.Function))
		ct.echoResetColor(pc)
	} else {
		pc.AppendString(checkedfuncname(source.Function))
	}
}

func (s *colorfulPainter) AddPrefixedString(pc *PrintCtx, prefix, name string, value string) {
	// s.Grow(len(name)*3 + 1 + 10)
	pc.AppendStringKeyPrefixed(name, prefix)
	s.AppendColon(pc)
	// s.pcAppendStringValue(intToString(value))
	pc.AppendString(value)
	// if s.noColor {
	// 	s.pcAppendQuotedStringValue(value)
	// } else {
	// 	s.pcAppendString(value)
	// }
}

func (s *colorfulPainter) TryQuoteValue(pc *PrintCtx, val string) {
	pc.AppendStringValue(val)
}

func (s *colorfulPainter) MarshalValue(pc *PrintCtx, val any) (handled bool, err error) {
	if m, ok := val.(encoding.TextMarshaler); ok {
		var data []byte
		data, err = m.MarshalText()
		if err != nil {
			hintInternal(err, "MarshalText failed")
			return
		}
		pc.AppendStringValue(string(data))
		handled = true
	}
	return
}

type logfmtPainter struct{}

func (s *logfmtPainter) Colorful() bool { return false }

func (s *logfmtPainter) Begin(pc *PrintCtx) {
	// pc.AppendByte('{')
}

func (s *logfmtPainter) End(pc *PrintCtx, newline bool) {
	// pc.AppendByte('}')
	if newline {
		pc.AppendByte('\n')
	}
}

func (s *logfmtPainter) BeginArray(pc *PrintCtx) {
	// pc.AppendByte('{')
}

func (s *logfmtPainter) EndArray(pc *PrintCtx, newline bool) {
	// pc.AppendByte('}')
	if newline {
		pc.AppendByte('\n')
	}
}

func (s *logfmtPainter) AppendColon(pc *PrintCtx) {
	pc.AppendByte('=')
}

func (s *logfmtPainter) AppendComma(pc *PrintCtx) {
	pc.AppendByte(',')
}

func (s *logfmtPainter) AppendTime(pc *PrintCtx, z time.Time) {
	const layout = time.RFC3339Nano
	pc.AppendByte('"')
	pc.Buf(func(buf []byte) []byte {
		return z.AppendFormat(buf, layout)
	})
	pc.AppendByte('"')
}

func (s *logfmtPainter) AppendTimeSlice(pc *PrintCtx, z []time.Time) {
	pc.AppendByte('[')
	if l := len(z); l > 0 {
		s.AppendTime(pc, z[0])
		for i := 1; i < len(z); i++ {
			dur := z[i]
			pc.AppendByte(',')
			s.AppendTime(pc, dur)
		}
	}
	pc.AppendByte(']')
}

func (s *logfmtPainter) AppendTimestamp(pc *PrintCtx, tm time.Time, layout string) {
	pc.AppendByte('"')
	pc.Buf(func(buf []byte) []byte {
		return tm.AppendFormat(buf, layout)
	})
	pc.AppendByte('"')
}

func (s *logfmtPainter) AppendDuration(pc *PrintCtx, z time.Duration) {
	// s.pcAppendByte('"')
	// s.appendEscapedJSONString(z.String())
	// s.pcAppendByte('"')
	pc.AppendQuotedString(z.String())
}

func (s *logfmtPainter) AppendDurationSlice(pc *PrintCtx, z []time.Duration) {
	pc.AppendByte('[')
	if l := len(z); l > 0 {
		s.AppendDuration(pc, z[0])
		for i := 1; i < len(z); i++ {
			dur := z[i]
			pc.AppendByte(',')
			s.AppendDuration(pc, dur)
		}
	}
	pc.AppendByte(']')
}

func (s *logfmtPainter) Append(pc *PrintCtx, data []byte) {
	pc.AppendBytes(data)
}

func (s *logfmtPainter) AppendStringKeyPrefixed(pc *PrintCtx, str, prefix string) {
	_, _ = pc.WriteString(prefix)
	_ = pc.WriteByte('.')
	_, _ = pc.WriteString(str)
}

func (s *logfmtPainter) AppendStringKey(pc *PrintCtx, str string) {
	_, _ = pc.WriteString(str)
}

func (s *logfmtPainter) AppendKey(pc *PrintCtx, key string, clr, bg color.Color) {
	if pc.IsColorfulStyle() {
		ct.echoColorAndBg(pc, clrAttrKey, clrAttrKeyBg)
		s.AppendStringKey(pc, key)
		ct.echoColorAndBg(pc, clr, bg)
	} else {
		s.AppendStringKey(pc, key)
	}
	s.AppendColon(pc)
}

func (s *logfmtPainter) AppendError(pc *PrintCtx, err error) {
	s.TryQuoteValue(pc, err.Error())
}

func (s *logfmtPainter) AppendErrorAfterPrinted(pc *PrintCtx, err error) {
}

func (s *logfmtPainter) AddTimestampField(pc *PrintCtx, tm time.Time) {
	pc.AppendStringKey(timestampFieldName)
	pc.AddColon()
	// pc.pcAppendByte('"')
	pc.AppendTimestamp(pc.now)
	pc.AddComma()
}

func (s *logfmtPainter) AddLoggerNameField(pc *PrintCtx, name string) {
	pc.AddString("logger", name)
	pc.AddComma()
}

func (s *logfmtPainter) AddSeverity(pc *PrintCtx, lvl Level) {
	pc.AppendString(ct.wrapRune(lvl.ShortTag(levelOutputWidth), '[', ']'))
	pc.AppendByte(' ')
}

func (s *logfmtPainter) AddPCField(pc *PrintCtx, source *Source) {
	pc.AddComma()

	pc.AppendStringKey(callerFieldName)
	pc.AddColon()
	pc.AppendByte('{')

	pc.AddString("file", source.File)
	pc.AddComma()
	pc.AddInt("line", source.Line)
	pc.AddComma()
	pc.AddString("function", source.Function)

	pc.AppendByte('}')
}

func (s *logfmtPainter) AddPrefixedString(pc *PrintCtx, prefix, name string, value string) {
	// s.Grow(len(name)*3 + 1 + 10)
	pc.AppendStringKeyPrefixed(name, prefix)
	s.AppendColon(pc)
	// s.pcAppendStringValue(intToString(value))
	pc.AppendQuotedStringValue(value)
}

func (s *logfmtPainter) TryQuoteValue(pc *PrintCtx, val string) {
	pc.AppendQuotedString(val)
}

func (s *logfmtPainter) MarshalValue(pc *PrintCtx, val any) (handled bool, err error) {
	if m, ok := val.(encoding.TextMarshaler); ok {
		var data []byte
		data, err = m.MarshalText()
		if err != nil {
			hintInternal(err, "MarshalText failed")
			return
		}
		pc.AppendStringValue(string(data))
		handled = true
	}
	return
}

type jsonPainter struct{}

func (s *jsonPainter) Colorful() bool { return false }

func (s *jsonPainter) Begin(pc *PrintCtx) {
	pc.AppendByte('{')
}

func (s *jsonPainter) End(pc *PrintCtx, newline bool) {
	pc.AppendByte('}')
	if newline {
		pc.AppendByte('\n')
	}
}

func (s *jsonPainter) BeginArray(pc *PrintCtx) {
	pc.AppendByte('[')
}

func (s *jsonPainter) EndArray(pc *PrintCtx, newline bool) {
	pc.AppendByte(']')
	if newline {
		pc.AppendByte('\n')
	}
}

func (s *jsonPainter) AppendColon(pc *PrintCtx) {
	pc.AppendByte(':')
}

func (s *jsonPainter) AppendComma(pc *PrintCtx) {
	pc.AppendByte(',')
}

func (s *jsonPainter) AppendTime(pc *PrintCtx, z time.Time) {
	const layout = time.RFC3339Nano
	pc.AppendByte('"')
	pc.Buf(func(buf []byte) []byte {
		return z.AppendFormat(buf, layout)
	})
	pc.AppendByte('"')
}

func (s *jsonPainter) AppendTimeSlice(pc *PrintCtx, z []time.Time) {
	pc.AppendByte('[')
	if l := len(z); l > 0 {
		s.AppendTime(pc, z[0])
		for i := 1; i < len(z); i++ {
			dur := z[i]
			pc.AppendByte(',')
			s.AppendTime(pc, dur)
		}
	}
	pc.AppendByte(']')
}

func (s *jsonPainter) AppendTimestamp(pc *PrintCtx, tm time.Time, layout string) {
	pc.AppendByte('"')
	pc.Buf(func(buf []byte) []byte {
		return tm.AppendFormat(buf, layout)
	})
	pc.AppendByte('"')
}

func (s *jsonPainter) AppendDuration(pc *PrintCtx, z time.Duration) {
	// s.pcAppendByte('"')
	// s.appendEscapedJSONString(z.String())
	// s.pcAppendByte('"')
	pc.AppendQuotedString(z.String())
}

func (s *jsonPainter) AppendDurationSlice(pc *PrintCtx, z []time.Duration) {
	pc.AppendByte('[')
	if l := len(z); l > 0 {
		s.AppendDuration(pc, z[0])
		for i := 1; i < len(z); i++ {
			dur := z[i]
			pc.AppendByte(',')
			s.AppendDuration(pc, dur)
		}
	}
	pc.AppendByte(']')
}

func (s *jsonPainter) Append(pc *PrintCtx, data []byte) {
	pc.AppendBytes(data)
}

func (s *jsonPainter) AppendStringKeyPrefixed(pc *PrintCtx, str, prefix string) {
	_ = pc.WriteByte('"')
	_, _ = pc.WriteString(prefix)
	_ = pc.WriteByte('.')
	_, _ = pc.WriteString(str)
	_ = pc.WriteByte('"')
}

func (s *jsonPainter) AppendStringKey(pc *PrintCtx, str string) {
	_ = pc.WriteByte('"')
	_, _ = pc.WriteString(str)
	_ = pc.WriteByte('"')
}

func (s *jsonPainter) AppendKey(pc *PrintCtx, key string, clr, bg color.Color) {
	s.AppendStringKey(pc, key)
	s.AppendColon(pc)
}

func (s *jsonPainter) AppendError(pc *PrintCtx, err error) {
	pc.Begin()
	pc.AppendStringKey("message")
	pc.AddColon()
	pc.AppendQuotedString(err.Error())
	// ee := errorsv3.New("")
	// var e3 errorsv3.Error
	// if errors.As(err, &e3) {
	if f, ok := err.(*errorsv3.WithStackInfo); ok {
		if st := f.StackTrace(); st != nil {
			pc.AddComma()
			pc.AppendStringKey("trace")
			pc.AddColon()

			frame := st[0]
			pc.cachedSource.Extract(uintptr(frame))
			pc.Begin()
			pc.AppendStringKey("file")
			pc.AddColon()
			pc.AppendQuotedStringValue(pc.cachedSource.File)
			pc.AddComma()
			pc.AppendStringKey("line")
			pc.AddColon()
			pc.AppendInt(pc.cachedSource.Line)
			pc.AddComma()
			pc.AppendStringKey("func")
			pc.AddColon()
			pc.AppendQuotedStringValue(pc.cachedSource.Function)
			pc.End(false)
		}
	}
	// }
	pc.End(false)
}

func (s *jsonPainter) AppendErrorAfterPrinted(pc *PrintCtx, err error) {
}

func (s *jsonPainter) AddTimestampField(pc *PrintCtx, tm time.Time) {
	pc.AppendStringKey(timestampFieldName)
	pc.AddColon()
	// pc.pcAppendByte('"')
	pc.AppendTimestamp(pc.now)
	pc.AddComma()
}

func (s *jsonPainter) AddLoggerNameField(pc *PrintCtx, name string) {
	pc.AppendStringKey("logger")
	pc.AddColon()
	pc.AppendByte('"')
	pc.AppendStringValue(name)
	pc.AppendByte('"')
	pc.AddComma()
}

func (s *jsonPainter) AddSeverity(pc *PrintCtx, lvl Level) {
	pc.AddString(levelFieldName, lvl.String())
	// pc.pcAppendStringKey(levelFieldName)
	// pc.pcAppendColon()
	// pc.pcAppendByte('"')
	// pc.pcAppendStringValue(pc.lvl.String())
	// pc.pcAppendByte('"')
	pc.AddComma()
}

func (s *jsonPainter) AddPCField(pc *PrintCtx, source *Source) {
	pc.AddComma()

	pc.AppendStringKey(callerFieldName)
	pc.AddColon()
	pc.AppendByte('{')

	pc.AddString("file", source.File)
	pc.AddComma()
	pc.AddInt("line", source.Line)
	pc.AddComma()
	pc.AddString("function", source.Function)

	pc.AppendByte('}')
}

func (s *jsonPainter) AddPrefixedString(pc *PrintCtx, prefix, name string, value string) {
	// s.Grow(len(name)*3 + 1 + 10)
	pc.AppendStringKeyPrefixed(name, prefix)
	pc.AddColon()
	// s.pcAppendStringValue(intToString(value))
	pc.AppendQuotedStringValue(value)
}

func (s *jsonPainter) TryQuoteValue(pc *PrintCtx, val string) {
	pc.AppendQuotedString(val)
}

func (s *jsonPainter) MarshalValue(pc *PrintCtx, val any) (handled bool, err error) {
	if m, ok := val.(interface{ MarshalJSON() ([]byte, error) }); ok {
		var data []byte
		data, err = m.MarshalJSON()
		if err != nil {
			hintInternal(err, "MarshalJSON failed")
			return
		}
		pc.AppendStringValue(string(data))
		handled = true
	}
	return
}
