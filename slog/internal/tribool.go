package internal

import (
	"strings"
)

// TriBool wraps a tri-state boolean value.
//
// See: https://en.wikipedia.org/wiki/Three-valued_logic
type TriBool int

const (
	TriNoState TriBool = iota
	TriFalse
	TriTrue
)

func (s *TriBool) UnmarshalText(b []byte) error {
	switch strings.ToLower(string(b)) {
	case "true":
		*s = TriTrue
	case "false":
		*s = TriFalse
	default:
		*s = TriNoState
	}
	return nil
}

func (s TriBool) MarshalText() (b []byte, err error) {
	switch s {
	case TriTrue:
		b = []byte("true")
	case TriFalse:
		b = []byte("false")
	default:
		b = []byte("unset")
	}
	return
}

func (s TriBool) String() string {
	switch s {
	case TriTrue:
		return "true"
	case TriFalse:
		return "false"
	default:
		return "unset"
	}
}

func (s TriBool) Equal(b TriBool) bool {
	if s == TriNoState {
		return false
	}
	return s == b
}

func (s TriBool) And(b TriBool) TriBool {
	if s == TriFalse || b == TriFalse {
		return TriFalse
	}
	if s == TriNoState || b == TriNoState {
		return TriNoState
	}
	return TriTrue
}

func (s TriBool) Or(b TriBool) TriBool {
	if s == TriTrue || b == TriTrue {
		return TriTrue
	}
	if s == TriNoState || b == TriNoState {
		return TriNoState
	}
	return TriFalse
}

func (s TriBool) Xor(b TriBool) TriBool {
	if s == TriNoState || b == TriNoState {
		return TriNoState
	}
	if s != b {
		return TriTrue
	}
	return TriFalse
}

func (s TriBool) Not() TriBool {
	if s == TriNoState {
		return s
	}
	if s == TriTrue {
		return TriFalse
	}
	return TriTrue
}
