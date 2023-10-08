package internal

type TriBool int

const (
	TriNoState TriBool = iota
	TriFalse
	TriTrue
)

func (s TriBool) Equal(b TriBool) bool {
	if s == TriNoState {
		return false
	}
	return s == b
}

func (s TriBool) And(b TriBool) TriBool {
	if s == TriNoState {
		if b == TriFalse {
			return b
		}
		return s
	}
	if s == TriTrue && b == TriTrue {
		return TriTrue
	}
	return TriFalse
}

func (s TriBool) Or(b TriBool) TriBool {
	if s == TriNoState {
		if b == TriTrue {
			return b
		}
		return s
	}
	if s == TriTrue || b == TriTrue {
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
