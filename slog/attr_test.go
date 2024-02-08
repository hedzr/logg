package slog

import (
	"testing"
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
