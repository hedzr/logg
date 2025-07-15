package slog

import (
	"context"
	"strings"
	"sync/atomic"
	"time"
)

func (s *Entry) print(ctx context.Context, lvl Level, timestamp time.Time, stackFrame uintptr, msg string, kvps Attrs) {
	pc := poolPrintCtx.Get().(*PrintCtx)

	// pc.set will truncate internal buffer and reset all states for
	// this current session. So, don't worry about a reused buffer
	// takes wasted bytes.
	pc.set(s, lvl, timestamp, stackFrame, msg, kvps)

	s.printImpl(ctx, pc)

	pc.putBack()
}

func (s *Entry) printImpl(ctx context.Context, pc *PrintCtx) {
	_ = ctx

	// s.Println() or s.Println("") will print out just an empty line,
	// without timestamp, loggername, and others decorated fields.
	if pc.lvl == AlwaysLevel && strings.Trim(pc.msg, "\n\r \t") == "" {
		s.printOut(pc.lvl, []byte{'\n'})
		return
	}

	colorStyle := pc.IsColorStyle()
	if colorStyle {
		pc.SetupColors()
	}

	pc.Begin()

	s.printTimestamp(pc)
	s.printLoggerName(pc)
	s.printSeverity(pc)

	if colorStyle {
		s.printFirstLineOfMsg(pc)
	} else { // json or logfmt
		s.printMsg(pc)
	}

	holdErrorValue := serializeAttrs(pc, pc.kvps)

	if IsAnyBitsSet(Lcaller) {
		s.printPC(pc)
	}

	s.printRestLinesOfMsg(pc)

	pc.appendErrorAfterPrinted(holdErrorValue)

	pc.End(true)

	// ret = pc.String()
	// s.printOut(pc.lvl, []byte(ret))
	msg := pc.Bytes()
	s.printOut(pc.lvl, msg)
	// return
}

func (s *Entry) printTimestamp(pc *PrintCtx) {
	pc.AddTimestampField()

	// if pc.IsColorStyle() {
	// 	if pc.colorful {
	// 		ct.echoColor(pc, clrTimestamp)
	// 	}
	// 	pc.AppendTimestamp(pc.now)
	// 	pc.pcAppendByte(' ')
	// } else {
	// 	pc.AppendStringKey(timestampFieldName)
	// 	pc.AddColon()
	// 	// pc.pcAppendByte('"')
	// 	pc.AppendTimestamp(pc.now)
	// 	pc.AddComma()
	// }

	// if pc.noColor { // json or logfmt
	// 	pc.pcAppendStringKey(timestampFieldName)
	// 	pc.pcAppendColon()
	// 	// pc.pcAppendByte('"')
	// 	pc.appendTimestamp(pc.now)
	// 	pc.pcAppendComma()
	// } else {
	// 	if pc.colorful {
	// 		ct.echoColor(pc, clrTimestamp)
	// 	}
	// 	pc.appendTimestamp(pc.now)
	// 	pc.pcAppendByte(' ')
	// }
}

func (s *Entry) printLoggerName(pc *PrintCtx) {
	if s.name != "" {
		pc.ip.AddLoggerNameField(pc, s.name)

		// switch s.mode {
		// case ModeJSON:
		// 	pc.AppendStringKey("logger")
		// 	pc.pcAppendColon()
		// 	pc.pcAppendByte('"')
		// 	pc.AppendStringValue(s.name)
		// 	pc.pcAppendByte('"')
		// 	pc.pcAppendComma()
		// case ModeLogFmt:
		// 	pc.AddString("logger", s.name)
		// 	pc.pcAppendComma()
		// default:
		// 	l, lnw := len(s.name), int(atomic.LoadInt32(&longestNameWidth))
		// 	if r := lnw - l; r > 0 {
		// 		pc.AppendRuneTimes(' ', r)
		// 	}
		// 	if pc.colorful {
		// 		ct.wrapColorAndBgTo(pc, clrLoggerName, clrLoggerNameBg, s.name)
		// 	} else {
		// 		pc.AppendString(s.name)
		// 	}
		// 	pc.pcAppendByte(' ')
		// }

		// if pc.noColor { // json or logfmt
		// 	if pc.jsonMode {
		// 		pc.pcAppendStringKey("logger")
		// 		pc.pcAppendColon()
		// 		pc.pcAppendByte('"')
		// 		pc.pcAppendStringValue(s.name)
		// 		pc.pcAppendByte('"')
		// 	} else {
		// 		pc.AddString("logger", s.name)
		// 	}
		// 	pc.pcAppendComma()
		// } else {
		// 	if pc.colorful {
		// 		ct.wrapColorAndBgTo(pc, clrLoggerName, clrLoggerNameBg, s.name)
		// 	} else {
		// 		pc.pcAppendString(s.name)
		// 	}
		// 	pc.pcAppendByte(' ')
		// }
	} else {
		lnw := int(atomic.LoadInt32(&longestNameWidth))
		pc.AppendRuneTimes(' ', lnw+1)
	}
}

func (s *Entry) printSeverity(pc *PrintCtx) {
	pc.ip.AddSeverity(pc, pc.lvl)

	// switch s.mode {
	// case ModeJSON, ModeLogFmt:
	// 	pc.AddString(levelFieldName, pc.lvl.String())
	// 	// pc.pcAppendStringKey(levelFieldName)
	// 	// pc.pcAppendColon()
	// 	// pc.pcAppendByte('"')
	// 	// pc.pcAppendStringValue(pc.lvl.String())
	// 	// pc.pcAppendByte('"')
	// 	pc.pcAppendComma()
	// default:
	// 	if pc.colorful {
	// 		ct.wrapColorAndBgTo(pc, pc.clr, pc.bg, ct.wrapRune(pc.lvl.ShortTag(levelOutputWidth), '[', ']'))
	// 	} else {
	// 		pc.AppendString(ct.wrapRune(pc.lvl.ShortTag(levelOutputWidth), '[', ']'))
	// 	}
	// 	pc.pcAppendByte(' ')
	// }

	// if pc.noColor { // json or logfmt
	// 	pc.AddString(levelFieldName, pc.lvl.String())
	// 	// pc.pcAppendStringKey(levelFieldName)
	// 	// pc.pcAppendColon()
	// 	// pc.pcAppendByte('"')
	// 	// pc.pcAppendStringValue(pc.lvl.String())
	// 	// pc.pcAppendByte('"')
	// 	pc.pcAppendComma()
	// } else {
	// 	if pc.colorful {
	// 		ct.wrapColorAndBgTo(pc, pc.clr, pc.bg, ct.wrapRune(pc.lvl.ShortTag(levelOutputWidth), '[', ']'))
	// 	} else {
	// 		pc.pcAppendString(ct.wrapRune(pc.lvl.ShortTag(levelOutputWidth), '[', ']'))
	// 	}
	// 	pc.pcAppendByte(' ')
	// }
}

func (s *Entry) printPC(pc *PrintCtx) {
	pc.ip.AddPCField(pc, pc.source())

	// switch s.mode {
	// case ModeJSON, ModeLogFmt:
	// 	pc.pcAppendComma()
	//
	// 	source := pc.source()
	// 	if s.mode == ModeJSON {
	// 		pc.AppendStringKey(callerFieldName)
	// 		pc.pcAppendColon()
	// 		pc.pcAppendByte('{')
	//
	// 		pc.AddString("file", source.File)
	// 		pc.pcAppendComma()
	// 		pc.AddInt("line", source.Line)
	// 		pc.pcAppendComma()
	// 		pc.AddString("function", source.Function)
	//
	// 		pc.pcAppendByte('}')
	// 	} else {
	// 		pc.AddPrefixedString(callerFieldName, "file", source.File)
	// 		pc.pcAppendComma()

	// 		pc.AddPrefixedInt(callerFieldName, "line", source.Line)
	// 		pc.pcAppendComma()
	//
	// 		pc.AddPrefixedString(callerFieldName, "function", source.Function)
	// 	}
	// 	// pc.pcAppendComma()
	// default:
	// 	source := pc.source()
	// 	pc.pcAppendByte(' ')
	// 	// pc.appendRune('(')
	// 	pc.AppendString(source.File)
	// 	pc.pcAppendByte(':')
	// 	pc.AppendInt(source.Line)
	// 	// pc.appendRune(')')
	// 	pc.pcAppendByte(' ')
	// 	// ct.wrapDimColorTo(pc.SB, source.checkedfuncname()) // clion p-term in run panel cannot support dim color.
	// 	if pc.colorful {
	// 		ct.wrapColorTo(pc, clrFuncName, checkedfuncname(source.Function))
	// 		ct.echoResetColor(pc)
	// 	} else {
	// 		pc.AppendString(checkedfuncname(source.Function))
	// 	}
	// }

	// if pc.noColor {
	// 	pc.pcAppendComma()
	//
	// 	source := pc.source()
	// 	if pc.jsonMode {
	// 		pc.pcAppendStringKey(callerFieldName)
	// 		pc.pcAppendColon()
	// 		pc.pcAppendByte('{')
	//
	// 		pc.AddString("file", source.File)
	// 		pc.pcAppendComma()
	// 		pc.AddInt("line", source.Line)
	// 		pc.pcAppendComma()
	// 		pc.AddString("function", source.Function)
	//
	// 		pc.pcAppendByte('}')
	// 	} else {
	// 		pc.AddPrefixedString(callerFieldName, "file", source.File)
	// 		pc.pcAppendComma()
	//
	// 		pc.AddPrefixedInt(callerFieldName, "line", source.Line)
	// 		pc.pcAppendComma()
	//
	// 		pc.AddPrefixedString(callerFieldName, "function", source.Function)
	// 	}
	// 	// pc.pcAppendComma()
	// 	return
	// }

	// source := pc.source()
	// pc.pcAppendByte(' ')
	// // pc.appendRune('(')
	// pc.pcAppendString(source.File)
	// pc.pcAppendByte(':')
	// pc.AppendInt(source.Line)
	// // pc.appendRune(')')
	// pc.pcAppendByte(' ')
	// // ct.wrapDimColorTo(pc.SB, source.checkedfuncname()) // clion p-term in run panel cannot support dim color.
	// if pc.colorful {
	// 	ct.wrapColorTo(pc, clrFuncName, checkedfuncname(source.Function))
	// 	ct.echoResetColor(pc)
	// } else {
	// 	pc.pcAppendString(checkedfuncname(source.Function))
	// }
}

func (s *Entry) printMsg(pc *PrintCtx) {
	pc.ip.AddMsgField(pc, pc.msg)

	// switch s.mode {
	// case ModeJSON, ModeLogFmt:
	// 	pc.AddString(messageFieldName, pc.msg)
	// 	// pc.pcAppendComma()
	// default:
	// 	pc.AddString(messageFieldName, ct.translate(pc.msg))
	// 	// pc.pcAppendByte(' ')
	// }

	// if pc.noColor {
	// 	pc.AddString(messageFieldName, pc.msg)
	// 	// pc.pcAppendComma()
	// } else {
	// 	pc.AddString(messageFieldName, ct.translate(pc.msg))
	// 	// pc.pcAppendByte(' ')
	// }
	// // NOTE: serializeAttrs() will supply a leading comma char.
}

func (s *Entry) printFirstLineOfMsg(pc *PrintCtx) {
	var firstLine string
	firstLine, pc.restLines, pc.eol = ct.splitFirstAndRestLines(pc.msg)
	pc.ip.AddMsgFieldFirstLine(pc, firstLine)

	// if minimalMessageWidth > 0 {
	// 	str := ct.rightPad(firstLine, " ", minimalMessageWidth)
	// 	if pc.colorful {
	// 		str = ct.translate(str)
	// 		_, _ = pc.WriteString(ct.wrapColorAndBg(str, pc.clr, pc.bg))
	// 	} else {
	// 		_, _ = pc.WriteString(str)
	// 	}
	// } else {
	// 	if pc.colorful {
	// 		str := ct.translate(firstLine)
	// 		_, _ = pc.WriteString(ct.wrapColorAndBg(str, pc.clr, pc.bg))
	// 	} else {
	// 		_, _ = pc.WriteString(firstLine)
	// 	}
	// }

	// pc.pcAppendByte(' ')
	// pc.pcAppendByte('|')
}

func (s *Entry) printRestLinesOfMsg(pc *PrintCtx) {
	pc.ip.AddMsgFieldRestLines(pc, pc.restLines, pc.eol)

	// switch s.mode {
	// case ModeJSON, ModeLogFmt:
	// default:
	// 	if pc.restLines != "" {
	// 		pc.pcAppendByte('\n')
	// 		pc.AppendString(ct.padFunc(pc.restLines, " ", 4, func(i int, line string) string {
	// 			if pc.colorful {
	// 				return ct.wrapColorAndBg(line, pc.clr, pc.bg)
	// 			} else {
	// 				return line
	// 			}
	// 		}))
	// 		if pc.eol {
	// 			pc.pcAppendByte('\n')
	// 		}
	// 	}
	// }

	// if !pc.noColor && pc.restLines != "" {
	// 	pc.pcAppendByte('\n')
	// 	pc.pcAppendString(ct.padFunc(pc.restLines, " ", 4, func(i int, line string) string {
	// 		if pc.colorful {
	// 			return ct.wrapColorAndBg(line, pc.clr, pc.bg)
	// 		} else {
	// 			return line
	// 		}
	// 	}))
	// 	if pc.eol {
	// 		pc.pcAppendByte('\n')
	// 	}
	// }
}
