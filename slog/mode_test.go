package slog

import "testing"

func TestMode(t *testing.T) {
	var strctr []string

	modes := []Mode{ModeJSON, ModeLogFmt, ModeColorful, ModePlain, ModeUndefined}

	for _, m := range modes {
		t.Logf("%s, ", m)
		strctr = append(strctr, m.String())
	}

	for _, str := range strctr {
		var mode Mode
		if err := (&mode).UnmarshalText([]byte(str)); err != nil {
			t.Fatalf("error in Mode.UnmarshalText: %+v", err)
		}
		t.Logf("unmarshalled mode is: %v", mode)
	}

	t.Logf("\n")
}
