//go:generate stringer --type=Mode --trimprefix=Mode --linecomment

package slog

import (
	"fmt"
	"strings"
)

type Mode int

const (
	ModeJSON      Mode       = iota // json
	ModeLogFmt                      // logfmt
	ModeColorful                    // color
	ModePlain                       // plain
	ModeUndefined = Mode(99)        // undefined
)

// func (a Mode) String() string {
// 	t, _ := a.MarshalText()
// 	return string(t)
// }

func (a Mode) MarshalText() (text []byte, err error) {
	text, err = []byte(a.String()), nil
	return
}

func (a *Mode) UnmarshalText(text []byte) error {
	str := string(text)
	for _, piece := range []struct {
		start int
		name  string
		index []uint8
	}{
		{0, _Mode_name_0, _Mode_index_0[:]},
	} {
		if idx := strings.Index(piece.name, str); idx >= 0 {
			for i, ix := range piece.index {
				if int(ix) == idx && len(str)+idx == int(piece.index[i+1]) {
					*a = Mode(i + piece.start)
					return nil
				}
			}
		}
	}

	lastPiece := _Mode_name_1
	lastMode := ModeUndefined
	if str == lastPiece {
		*a = lastMode
		return nil
	}

	return fmt.Errorf("Unknown Mode text %q", string(text))
}

// func (a Mode) MarshalText() (text []byte, err error) {
// 	var str string
// 	switch a {
// 	case ModeJSON:
// 		str = "json"
// 	case ModeLogFmt:
// 		str = "logfmt"
// 	case ModeColorful:
// 		str = "color"
// 	case ModePlain:
// 		str = "plain"
// 	case ModeUndefined:
// 		str = "undefined"
// 	default:
// 		str = fmt.Sprintf("Mode(%d)", int(a))
// 	}
// 	return []byte(str), nil
// }

// func (a *Mode) UnmarshalText(text []byte) error {
// 	switch strings.ToLower(string(text)) {
// 	case "json":
// 		*a = ModeJSON
// 	case "logfmt":
// 		*a = ModeLogFmt
// 	case "text", "plain":
// 		*a = ModePlain
// 	case "color":
// 		*a = ModeColorful
// 	case "undefined", "unk", "unknown":
// 		*a = ModeUndefined
// 	}
// 	return fmt.Errorf("Unknown Mode text %q", string(text))
// }
