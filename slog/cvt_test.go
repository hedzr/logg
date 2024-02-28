package slog

import (
	"testing"
	"time"
)

func TestCvt(t *testing.T) {
	t.Run("intToString", func(t *testing.T) {
		for i, c := range []struct {
			from   int
			expect string
		}{
			{-100, "-100"},
			{-1, "-1"},
			{0, "0"},
			{1, "1"},
			{101, "101"},
		} {
			actual := intToString(c.from)
			if actual != c.expect {
				t.Fatalf("%5d. expecting %q, but got %q", i, c.expect, actual)
			}
			t.Logf("s: %q", actual)
		}
	})
	t.Run("intToStringEx", func(t *testing.T) {
		for i, c := range []struct {
			from   int
			expect string
			base   int
		}{
			{-100, "-100", 10},
			{-1, "-1", 10},
			{0, "0", 10},
			{1, "1", 10},
			{101, "101", 10},

			{-100, "-64", 16},
			{-1, "-1", 16},
			{0, "0", 16},
			{1, "1", 16},
			{101, "65", 16},
		} {
			actual := intToStringEx(c.from, c.base)
			if actual != c.expect {
				t.Fatalf("%5d. expecting %q, but got %q", i, c.expect, actual)
			}
			t.Logf("s: %q", actual)
		}
	})

	t.Run("uintToString", func(t *testing.T) {
		for i, c := range []struct {
			from   uint
			expect string
		}{
			{0, "0"},
			{1, "1"},
			{101, "101"},
		} {
			actual := uintToString(c.from)
			if actual != c.expect {
				t.Fatalf("%5d. expecting %q, but got %q", i, c.expect, actual)
			}
			t.Logf("s: %q", actual)
		}
	})
	t.Run("uintToStringEx", func(t *testing.T) {
		for i, c := range []struct {
			from   uint
			expect string
			base   int
		}{
			{0, "0", 10},
			{1, "1", 10},
			{101, "101", 10},

			{0, "0", 16},
			{1, "1", 16},
			{101, "65", 16},
		} {
			actual := uintToStringEx(c.from, c.base)
			if actual != c.expect {
				t.Fatalf("%5d. expecting %q, but got %q", i, c.expect, actual)
			}
			t.Logf("s: %q", actual)
		}
	})

	t.Run("floatToString", func(t *testing.T) {
		for i, c := range []struct {
			from   float64
			expect string
		}{
			{0, "0"},
			{1, "1"},
			{101, "101"},
		} {
			actual := floatToString(c.from)
			if actual != c.expect {
				t.Fatalf("%5d. expecting %q, but got %q", i, c.expect, actual)
			}
			t.Logf("s: %q", actual)
		}
	})
	t.Run("floatToStringEx", func(t *testing.T) {
		for i, c := range []struct {
			from   float32
			expect string
			format byte
		}{
			{0, "0e+00", 'e'},
			{1, "1e+00", 'e'},
			{101, "1.01e+02", 'e'},

			{0, "0", 'f'},
			{1, "1", 'f'},
			{101.12345678901, "101.12345886230469", 'f'},
		} {
			actual := floatToStringEx(c.from, c.format, -1, 64)
			if actual != c.expect {
				t.Fatalf("%5d. expecting %q, but got %q", i, c.expect, actual)
			}
			t.Logf("s: %q", actual)
		}
	})

	t.Run("complexToString", func(t *testing.T) {
		for i, c := range []struct {
			from   complex128
			expect string
		}{
			{0, "(0+0i)"},
			{1, "(1+0i)"},
			{101, "(101+0i)"},
		} {
			actual := complexToString(c.from)
			if actual != c.expect {
				t.Fatalf("%5d. expecting %q, but got %q", i, c.expect, actual)
			}
			t.Logf("s: %q", actual)
		}
	})
	t.Run("complexToStringEx", func(t *testing.T) {
		for i, c := range []struct {
			from   complex128
			expect string
			format byte
		}{
			{0, "(0e+00+0e+00i)", 'e'},
			{1, "(1e+00+0e+00i)", 'e'},
			{101, "(1.01e+02+0e+00i)", 'e'},

			{0, "(0+0i)", 'f'},
			{1, "(1+0i)", 'f'},
			{101.12345678901, "(101.12346+0i)", 'f'},
		} {
			actual := complexToStringEx(c.from, c.format, -1, 64)
			if actual != c.expect {
				t.Fatalf("%5d. expecting %q, but got %q", i, c.expect, actual)
			}
			t.Logf("s: %q", actual)
		}
	})

	t.Run("boolToString", func(t *testing.T) {
		for i, c := range []struct {
			from   bool
			expect string
		}{
			{false, "false"},
			{true, "true"},
		} {
			actual := boolToString(c.from)
			if actual != c.expect {
				t.Fatalf("%5d. expecting %q, but got %q", i, c.expect, actual)
			}
			t.Logf("s: %q", actual)
		}
	})

	t.Run("intSliceToString", func(t *testing.T) {
		for i, c := range []struct {
			from   []int
			expect string
		}{
			{[]int{-100, 0, 101}, "[-100,0,101]"},
			{[]int{}, "[]"},
			{nil, "[]"},
		} {
			actual := intSliceToString(c.from)
			if actual != c.expect {
				t.Fatalf("%5d. expecting %q, but got %q", i, c.expect, actual)
			}
			t.Logf("s: %q", actual)
		}
	})

	t.Run("uintSliceToString", func(t *testing.T) {
		for i, c := range []struct {
			from   []uint
			expect string
		}{
			{[]uint{100, 0, 101}, "[100,0,101]"},
			{[]uint{}, "[]"},
			{nil, "[]"},
		} {
			actual := uintSliceToString(c.from)
			if actual != c.expect {
				t.Fatalf("%5d. expecting %q, but got %q", i, c.expect, actual)
			}
			t.Logf("s: %q", actual)
		}
	})

	t.Run("floatSliceToString", func(t *testing.T) {
		for i, c := range []struct {
			from   []float64
			expect string
		}{
			{[]float64{-100, 0, 101}, "[-100,0,101]"},
			{[]float64{}, "[]"},
			{nil, "[]"},
			{[]float64{-100.123, 0, 101}, "[-100.123,0,101]"},
		} {
			actual := floatSliceToString(c.from)
			if actual != c.expect {
				t.Fatalf("%5d. expecting %q, but got %q", i, c.expect, actual)
			}
			t.Logf("s: %q", actual)
		}
	})

	t.Run("complexSliceToString", func(t *testing.T) {
		for i, c := range []struct {
			from   []complex128
			expect string
		}{
			{[]complex128{-100, 0, 101}, "[(-100+0i),(0+0i),(101+0i)]"},
			{[]complex128{}, "[]"},
			{nil, "[]"},
			{[]complex128{-100.123, 0, 101}, "[(-100.123+0i),(0+0i),(101+0i)]"},
		} {
			actual := complexSliceToString(c.from)
			if actual != c.expect {
				t.Fatalf("%5d. expecting %q, but got %q", i, c.expect, actual)
			}
			t.Logf("s: %q", actual)
		}
	})

	t.Run("stringSliceToString", func(t *testing.T) {
		for i, c := range []struct {
			from   []string
			expect string
		}{
			{[]string{"-100", "0", "101"}, "[\"-100\",\"0\",\"101\"]"},
			{[]string{}, "[]"},
			{nil, "[]"},
			{[]string{"-100.123", "", "101"}, "[\"-100.123\",\"\",\"101\"]"},
		} {
			actual := stringSliceToString(c.from)
			if actual != c.expect {
				t.Fatalf("%5d. expecting %q, but got %q", i, c.expect, actual)
			}
			t.Logf("s: %q", actual)
		}
	})

	t.Run("boolSliceToString", func(t *testing.T) {
		for i, c := range []struct {
			from   []bool
			expect string
		}{
			{[]bool{false, true, true}, "[false,true,true]"},
			{[]bool{}, "[]"},
			{nil, "[]"},
			{[]bool{true, false, true}, "[true,false,true]"},
		} {
			actual := boolSliceToString(c.from)
			if actual != c.expect {
				t.Fatalf("%5d. expecting %q, but got %q", i, c.expect, actual)
			}
			t.Logf("s: %q", actual)
		}
	})

	t.Run("timeSliceToString", func(t *testing.T) {
		loc, _ := time.LoadLocation("Asia/Shanghai")

		for i, c := range []struct {
			from   []time.Time
			expect string
		}{
			{[]time.Time{time.Unix(0, 0).In(loc), time.Unix(1, 1).In(loc)},
				"[\"1970-01-01T08:00:00+08:00\",\"1970-01-01T08:00:01.000000001+08:00\"]"},
			{[]time.Time{}, "[]"},
			{nil, "[]"},
		} {
			actual := timeSliceToString(c.from)
			if actual != c.expect {
				t.Fatalf("%5d. expecting %q, but got %q", i, c.expect, actual)
			}
			t.Logf("s: %q", actual)
		}
	})

	t.Run("durationSliceToString", func(t *testing.T) {
		for i, c := range []struct {
			from   []time.Duration
			expect string
		}{
			{[]time.Duration{time.Second, time.Hour}, "[\"1s\",\"1h0m0s\"]"},
			{[]time.Duration{}, "[]"},
			{nil, "[]"},
		} {
			actual := durationSliceToString(c.from)
			if actual != c.expect {
				t.Fatalf("%5d. expecting %q, but got %q", i, c.expect, actual)
			}
			t.Logf("s: %q", actual)
		}
	})
}
