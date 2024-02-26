package internal

import (
	"testing"
)

func TestTriBool_Equal(t *testing.T) {
	var b = TriTrue
	if !b.Equal(TriTrue) {
		t.Fail()
	}
	if b.Equal(TriFalse) {
		t.Fail()
	}
	if b.Equal(TriNoState) {
		t.Fail()
	}

	if TriNoState.Equal(TriFalse) {
		t.Fail()
	}
}

func TestTriBool_And(t *testing.T) {
	for i, c := range []struct {
		src, obj, expect TriBool
	}{
		{TriTrue, TriTrue, TriTrue},
		{TriTrue, TriFalse, TriFalse},
		{TriTrue, TriNoState, TriNoState},
		{TriFalse, TriTrue, TriFalse},
		{TriFalse, TriFalse, TriFalse},
		{TriFalse, TriNoState, TriFalse},
		{TriNoState, TriTrue, TriNoState},
		{TriNoState, TriFalse, TriFalse},
		{TriNoState, TriNoState, TriNoState},
	} {
		if actual := c.src.And(c.obj); actual != c.expect {
			t.Fatalf("%5d. %v & %v, expect %v, but got %v", i, c.src, c.obj, c.expect, actual)
		}
	}
}

func TestTriBool_Or(t *testing.T) {
	for i, c := range []struct {
		src, obj, expect TriBool
	}{
		{TriTrue, TriTrue, TriTrue},
		{TriTrue, TriFalse, TriTrue},
		{TriTrue, TriNoState, TriTrue},
		{TriFalse, TriTrue, TriTrue},
		{TriFalse, TriFalse, TriFalse},
		{TriFalse, TriNoState, TriNoState},
		{TriNoState, TriTrue, TriTrue},
		{TriNoState, TriFalse, TriNoState},
		{TriNoState, TriNoState, TriNoState},
	} {
		if actual := c.src.Or(c.obj); actual != c.expect {
			t.Fatalf("%5d. %v | %v, expect %v, but got %v", i, c.src, c.obj, c.expect, actual)
		}
		t.Logf("%5d. %v | %v = %v", i, c.src, c.obj, c.expect)
	}
}

func TestTriBool_Xor(t *testing.T) {
	for i, c := range []struct {
		src, obj, expect TriBool
	}{
		{TriTrue, TriTrue, TriFalse},
		{TriTrue, TriFalse, TriTrue},
		{TriTrue, TriNoState, TriNoState},
		{TriFalse, TriTrue, TriTrue},
		{TriFalse, TriFalse, TriFalse},
		{TriFalse, TriNoState, TriNoState},
		{TriNoState, TriTrue, TriNoState},
		{TriNoState, TriFalse, TriNoState},
		{TriNoState, TriNoState, TriNoState},
	} {
		if actual := c.src.Xor(c.obj); actual != c.expect {
			t.Fatalf("%5d. %v ^ %v (xor), expect %v, but got %v", i, c.src, c.obj, c.expect, actual)
		}
	}
}

func TestTriBool_Not(t *testing.T) {
	for i, c := range []struct {
		src, expect TriBool
	}{
		{TriTrue, TriFalse},
		{TriFalse, TriTrue},
		{TriNoState, TriNoState},
	} {
		if actual := c.src.Not(); actual != c.expect {
			t.Fatalf("%5d. ! %v, expect %v, but got %v", i, c.src, c.expect, actual)
		}
	}
}

func TestTriBool_Marshal(t *testing.T) {
	for i, c := range []struct {
		src    TriBool
		expect string
	}{
		{TriTrue, "true"},
		{TriFalse, "false"},
		{TriNoState, "unset"},
	} {
		if actual, err := c.src.MarshalText(); string(actual) != c.expect || err != nil {
			t.Fatalf("%5d. %v marshal to: expect %v, but got %v", i, c.src, c.expect, actual)
		}
	}
}

func TestTriBool_Unmarshal(t *testing.T) {
	for i, c := range []struct {
		expect TriBool
		src    string
	}{
		{TriTrue, "true"},
		{TriFalse, "false"},
		{TriNoState, "unset"},
	} {
		var tb TriBool
		if err := tb.UnmarshalText([]byte(c.src)); tb != c.expect || err != nil {
			t.Fatalf("%5d. %v unmarshal to: expect %v, but got %v", i, c.src, c.expect, tb)
		}
	}
}
