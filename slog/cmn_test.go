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

func TestAddCodeHostingProviders(t *testing.T) {
	AddCodeHostingProviders("zetek.com", "ZK")
	Warn("warn for testing")
	t.Log("TestAddCodeHostingProviders")
}

func TestAddKnownPathMapping(t *testing.T) {
	AddKnownPathMapping("C:/", "~c/")
	Warn("warn for testing")
	t.Log("TestAddKnownPathMapping")
}

func TestAddKnownPathRegexpMapping(t *testing.T) {
	AddKnownPathRegexpMapping("C:/(.* )and( .*)", "~c/$1-n-$2")
	Warn("warn for testing")
	t.Log("TestAddKnownPathRegexpMapping")
}

func TestSetMessageMinimalWidth(t *testing.T) {
	SetMessageMinimalWidth(72)
	Warn("warn for testing", "right1", 1, "r2", true)
	t.Log("TestSetMessageMinimalWidth")

	Reset()

	Default().Warn("default . warn")
}
