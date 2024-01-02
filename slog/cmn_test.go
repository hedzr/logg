package slog

import (
	"testing"
)

func TestSetFlags(t *testing.T) {
	SetFlags(LstdFlags)

	A()
	severity()

	t.Logf("flags: %v", GetFlags())
}

func A() int        { return severity() }
func severity() int { return 2 }

type X struct {
	IsAllBitsSet bool `json:"is-all-bits-set,omitempty"`
	BattleEnd    int  `json:"battle-end,omitempty"`
}
