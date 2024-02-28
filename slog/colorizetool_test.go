package slog

import (
	"strings"
	"testing"

	"github.com/hedzr/is/term/color"
)

func TestColorizeToolS_WrapRuneTo(t *testing.T) {
	s := ct.wrapRune("hello", '<', '>')
	t.Logf("s: %q", s)

	var sb strings.Builder
	ct.wrapRuneTo(&sb, "hello", '<', '>')
	t.Logf("sb: %q", sb.String())

	if s != sb.String() {
		t.Fail()
	}

	s = ct.wrapHighlightColor("highlight text")
	t.Logf("s: %q", s)
	ct.wrapHighlightColorTo(&sb, "highlight text")
	t.Logf("sb: %q", sb.String())

	s = ct.wrapColor("green text", color.FgGreen)
	t.Logf("s: %q", s)

	s = ct.wrapDimColor("dim text")
	t.Logf("dim: %q", s)

	sb.Reset()
	ct.wrapDimColorTo(&sb, "dim text") // dim lite
	t.Logf(" sb: %q", sb.String())

	s = ct.pad("hello", " ", 10)
	t.Logf("s: %q", s)
	if s != "          hello" {
		t.Fail()
	}

	s = ct.translate("hello", color.FgDefault)
	t.Logf("s: %q", s)
	if s != "hello" {
		t.Fail()
	}

	for _, s = range []string{"中字 Right 对齐", "中字 Right 对齐\nAlign Width"} {
		for i := -2 + len(s); i < len(s)*2; i++ {
			s1 := ct.alignr(s, " ", i)
			t.Logf("align-right: %s", s1)
		}
	}

	s1 := ct.pad(s, " ", 8)
	t.Logf("pad: %q", s1)
	s1 = ct.rightPad(s, " ", 8)
	t.Logf("pad-right: %q", s1)

	s = "中字 Right 对齐\nAlign Width\n"
	s1 = ct.pad(s, " ", 8)
	t.Logf("pad: \n%s", s1)

	fl := ct.firstLine(s)
	rl, eol := ct.restLines(s)
	t.Logf("%v, %v, %v", fl, rl, eol)

	fl, rl, eol = ct.splitFirstAndRestLines(s)
	t.Logf("%v, %v, %v", fl, rl, eol)
}
