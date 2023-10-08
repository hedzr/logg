package slog

import (
	"io"
	"strconv"
	"strings"

	"github.com/hedzr/is/term/color"
)

var ct colorizeToolS

type colorizeToolS struct{}

func (colorizeToolS) wrapRuneTo(sb *strings.Builder, text string, pre, post rune) {
	sb.WriteRune(pre)
	sb.WriteString(text)
	sb.WriteRune(post)
}

func (colorizeToolS) wrapColorTo(w io.Writer, clr color.Color, text string) {
	color.WrapColorTo(w, clr, text)
}

func (colorizeToolS) wrapColorAndBgTo(w io.Writer, clr, bg color.Color, text string) {
	color.WrapColorAndBgTo(w, clr, bg, text)
}

func (colorizeToolS) wrapDimColorTo(w io.Writer, text string) {
	color.WrapDimToLite(w, text)
}

func (colorizeToolS) wrapHighlightColorTo(w io.Writer, text string) {
	color.WrapHighlightTo(w, text)
}

func (colorizeToolS) wrapRune(text string, pre, post rune) string {
	var buf = make([]rune, 0, len(text)+2)
	buf = append(buf, pre)
	buf = append(buf, []rune(text)...)
	buf = append(buf, post)
	return string(buf)
	// var sb strings.Builder
	// sb.WriteRune(pre)
	// sb.WriteString(text)
	// sb.WriteRune(post)
	// return sb.String()
}

func (colorizeToolS) wrapColor(text string, clr color.Color) string {
	// return color.ToColor(clr, text)

	var sb strings.Builder
	color.WrapColorTo(&sb, clr, text)
	return sb.String()
}

func (s colorizeToolS) wrapColorAndBg(text string, clr, bg color.Color) string {
	// return color.ToColor(clr, text)

	var sb strings.Builder
	if bg != clrNone {
		s.echoBgColor(&sb, bg)
	}
	color.WrapColorTo(&sb, clr, text)
	return sb.String()
}

func (colorizeToolS) wrapDimColor(text string) string {
	// return color.ToDim(text)

	var sb strings.Builder
	color.WrapDimTo(&sb, text)
	return sb.String()
}

func (colorizeToolS) wrapHighlightColor(text string) string {
	// return color.ToHighlight(text)

	var sb strings.Builder
	color.WrapHighlightTo(&sb, text)
	return sb.String()
}

func (colorizeToolS) echoBgColor(out io.Writer, clr color.Color) {
	if clr != clrNone {
		// _, _ = fmt.Fprintf(os.Stdout, "\x1b[%dm", c)
		_, _ = out.Write([]byte("\x1b["))
		_, _ = out.Write([]byte(strconv.Itoa(int(clr))))
		_, _ = out.Write([]byte{'m'})
	}
}

func (colorizeToolS) echoColor(out io.Writer, clr color.Color) {
	if clr != clrNone {
		// _, _ = fmt.Fprintf(os.Stdout, "\x1b[%dm", c)
		_, _ = out.Write([]byte("\x1b["))
		_, _ = out.Write([]byte(strconv.Itoa(int(clr))))
		_, _ = out.Write([]byte{'m'})
	}
}

func (colorizeToolS) echoColorAndBg(out io.Writer, clr, bg color.Color) {
	// _, _ = fmt.Fprintf(os.Stdout, "\x1b[%dm", c)

	if clr != clrNone {
		_, _ = out.Write([]byte("\x1b["))
		_, _ = out.Write([]byte(strconv.Itoa(int(clr))))
		_, _ = out.Write([]byte{'m'})
	}
	if bg != clrNone {
		_, _ = out.Write([]byte("\x1b["))
		_, _ = out.Write([]byte(strconv.Itoa(int(bg))))
		_, _ = out.Write([]byte{'m'})
	}
}

func (colorizeToolS) echoResetColor(out io.Writer) { //nolint:unused //no
	// _, _ = fmt.Fprint(os.Stdout, "\x1b[0m")
	_, _ = out.Write([]byte("\x1b[0m"))
}

//
//

func (colorizeToolS) translate(str string, initialColor ...color.Color) string {
	var clr = color.FgDefault
	for _, c := range initialColor {
		clr = c
	}
	return color.GetCPT().Translate(str, clr)
}

func (colorizeToolS) rightPad(str string, padChar string, minw int) string {
	l := minw - len(str)
	if l > 0 {
		return str + strings.Repeat(padChar, l)
	}
	return str
}

func (colorizeToolS) pad(str string, padChar string, count int) string {
	lead := strings.Repeat(padChar, count)
	if strings.Contains(str, "\n") {
		lines := strings.Split(str, "\n")
		for i := 0; i < len(lines); i++ {
			lines[i] = lead + lines[i]
		}
		return strings.Join(lines, "\n")
	}
	return lead + str
}

func (colorizeToolS) padFunc(str string, padChar string, count int, fn func(i int, line string) string) string {
	lead := strings.Repeat(padChar, count)
	if strings.Contains(str, "\n") {
		lines := strings.Split(str, "\n")
		for i := 0; i < len(lines); i++ {
			lines[i] = fn(i, lead+lines[i])
		}
		return strings.Join(lines, "\n")
	}
	return lead + str
}

func (colorizeToolS) firstLine(str string) string {
	a := strings.Split(str, "\n")
	return a[0]
}

func (colorizeToolS) restLines(str string) (ret string, eol bool) {
	if len(str) > 0 {
		eol = str[len(str)-1] == '\n'
		if eol {
			str = strings.TrimRight(str, "\n\r")
		}
		a := strings.Split(str, "\n")
		if len(a) > 1 {
			ret = strings.Join(a[1:], "\n")
		}
	}
	return
}

func (colorizeToolS) splitFirstAndRestLines(str string) (firstLine, restLines string, eol bool) {
	if len(str) > 0 {
		eol = str[len(str)-1] == '\n'
		if eol {
			str = strings.TrimRight(str, "\n\r")
		}
		ix := strings.IndexRune(str, '\n')
		if ix >= 0 {
			firstLine, restLines = str[:ix], str[ix+1:]
		} else {
			firstLine = str
		}
	}
	return
}
