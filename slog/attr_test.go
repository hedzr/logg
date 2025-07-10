package slog

import (
	"testing"
	"time"
)

func TestNewGroupedAttr(t *testing.T) {
	ga := NewGroupedAttr("grp",
		NewAttr("i1", 1),
		NewAttr("i2", 2),
	)
	ga.SetValue(NewAttr("i3", 3))
	ga.SetValue(NewAttrs("i3", 3))
	t.Logf("ga: %v", ga)
}

var testBytes []byte // test data; same as testString but as a slice.

func TestKvp_SerializeValueTo(t *testing.T) {
	pc := NewPrintCtx(testBytes)
	k := &kvp{"k1", 1}

	k.SerializeValueTo(pc)

	t.Logf("%v", pc.String())
}

func TestGkvp_SetValue(t *testing.T) {
	pc := NewPrintCtx(testBytes)

	g := gkvp{
		key:   "g1",
		items: nil,
	}

	// g.SetValue(1)
	// g.SerializeValueTo(pc)
	// t.Logf("%v", pc.String())

	g.Add(String("s1", "hello"))
	t.Logf("%v", pc.String())
	testBytes = testBytes[:0]
	pc = NewPrintCtx(testBytes)
	g.SerializeValueTo(pc)
	t.Logf("%v", pc.String())

	g.SetValue(String("s1", "hello"))
	testBytes = testBytes[:0]
	pc = NewPrintCtx(testBytes)
	g.SerializeValueTo(pc)
	t.Logf("%v", pc.String())

	g.SetValue([]Attr{
		String("s1", "world"),
	})
	testBytes = testBytes[:0]
	pc = NewPrintCtx(testBytes)
	g.SerializeValueTo(pc)
	t.Logf("%v", pc.String())

	g.SetValue(NewAttrs(
		String("s1", "hello"),
		Time("time", time.Now()),
	))
	testBytes = testBytes[:0]
	pc = NewPrintCtx(testBytes)
	g.SerializeValueTo(pc)
	t.Logf("%v", pc.String())

	// json mode
	testBytes = testBytes[:0]
	pc = NewPrintCtx(testBytes)
	// pc.jsonMode = true
	pc.mode = ModeJSON
	g.SerializeValueTo(pc)
	t.Logf("%v", pc.String())

	// and no color
	testBytes = testBytes[:0]
	pc = NewPrintCtx(testBytes)
	// pc.jsonMode = true
	// pc.noColor = true
	pc.mode = ModeLogFmt
	pc.colorful = false
	g.SerializeValueTo(pc)
	t.Logf("%v", pc.String())
}
