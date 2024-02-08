package strings

import (
	"testing"
)

func TestAddPrefix(t *testing.T) {
	str := AddPrefix('/', "leaf", "p1", "p2")
	if str != "p1/p2/leaf" {
		t.Fatalf("wrong AddPrefix(...)")
	}

	str = AddPrefixFaster('/', "leaf", "p1")
	if str != "p1/leaf" {
		t.Fatalf("wrong AddPrefixFaster(...)")
	}

	str = DotPrefix("leaf", "p1", "p2")
	if str != "p1.p2.leaf" {
		t.Fatalf("wrong DotPrefix(...)")
	}

	str = DotPrefix("leaf", "p1")
	if str != "p1.leaf" {
		t.Fatalf("wrong DotPrefix(...) #2")
	}

	str = DotPrefix("leaf", "")
	if str != "leaf" {
		t.Fatalf("wrong DotPrefix(...) #2")
	}

	str = DotPrefix("", "")
	if str != "" {
		t.Fatalf("wrong DotPrefix(...) #2")
	}
}

func TestStringToBool(t *testing.T) {
	v := StringToBool("1", false, false)
	if v != true {
		t.Fatalf("wrong, expect true")
	}

	v = StringToBool("off", true, false)
	if v != false {
		t.Fatalf("wrong, expect false")
	}

	v = StringToBool("", true, false)
	if v != false {
		t.Fatalf("wrong, expect false")
	}
}
