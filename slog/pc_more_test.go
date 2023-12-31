package slog

import (
	"fmt"
	"runtime/debug"
	"testing"
	"time"

	"github.com/hedzr/logg/slog/internal/times"
)

func TestPrintCtx_pcAppendString(t *testing.T) {
	var pc PrintCtx

	pc.pcAppendByte('(')
	pc.pcAppendByte(')')
	if str := pc.String(); str != "()" {
		t.Fatalf("pc is %q now, but expecting '()'", str)
	}

	pc.pcAppendRune('[')
	pc.pcAppendRune(']')
	if str := pc.String(); str != "()[]" {
		t.Fatalf("pc is %q now, but expecting '()[]'", str)
	}

	pc.pcAppendColon()
	if str := pc.String(); str != "()[]=" {
		t.Fatalf("pc is %q now, but expecting '()[]='", str)
	}

	pc.jsonMode = true
	pc.pcAppendColon()
	if str := pc.String(); str != "()[]=:" {
		t.Fatalf("pc is %q now, but expecting '()[]=:'", str)
	}
}

func TestPrintCtx_pcTryQuoteValue(t *testing.T) {
	var pc PrintCtx

	t.Log("json mode")
	pc.jsonMode = true
	pc.noColor = true
	for _, c := range []struct{ src, expect string }{
		{"msg", `"msg"`},
		{`msg"pc"`, `"msg\"pc\""`},
	} {
		pc.pcTryQuoteValue(c.src)
		if str := pc.String(); str != c.expect {
			t.Fatalf("pc is %q now, but expecting %q", str, c.expect)
		}
		pc.Reset()
	}

	t.Log("logfmt mode")
	pc.jsonMode = false
	pc.noColor = true
	for _, c := range []struct{ src, expect string }{
		{"msg", `"msg"`},
		{`msg"pc"`, `"msg\"pc\""`},
	} {
		pc.pcTryQuoteValue(c.src)
		if str := pc.String(); str != c.expect {
			t.Fatalf("pc is %q now, but expecting %q", str, c.expect)
		}
		pc.Reset()
	}

	t.Log("color mode")
	pc.jsonMode = false
	pc.noColor = false
	for _, c := range []struct{ src, expect string }{
		{"msg", `msg`},
		{`msg"pc"`, `msg"pc"`},
	} {
		pc.pcTryQuoteValue(c.src)
		if str := pc.String(); str != c.expect {
			t.Fatalf("pc is %q now, but expecting %q", str, c.expect)
		}
		pc.Reset()
	}
}

func TestPrintCtx_pcAppendStringValue(t *testing.T) {
	var pc PrintCtx

	for _, c := range []struct {
		json, noColor bool
		src, expect   string
	}{
		{true, false, "msg", `"msg"`},
		{true, false, `msg"pc"`, `"msg\"pc\""`},
	} {
		pc.jsonMode, pc.noColor = c.json, c.noColor
		pc.pcAppendQuotedStringValue(c.src)
		if str := pc.String(); str != c.expect {
			t.Fatalf("pc is %q now, but expecting %q", str, c.expect)
		}
		pc.Reset()
	}

	for _, c := range []struct {
		json, noColor bool
		src, expect   string
	}{
		{true, false, "msg", `msg`},
		{true, false, `msg"pc"`, `msg"pc"`},
	} {
		pc.jsonMode, pc.noColor = c.json, c.noColor
		pc.pcAppendStringValue(c.src)
		if str := pc.String(); str != c.expect {
			t.Fatalf("pc is %q now, but expecting %q", str, c.expect)
		}
		pc.Reset()
	}
}

func TestPrintCtx_appendDuration(t *testing.T) {
	var pc PrintCtx

	for _, c := range []struct {
		src    time.Duration
		expect string
	}{
		{279001 * time.Nanosecond, `"279.001µs"`},
		{312 * time.Microsecond, `"312µs"`},
		{279001 * time.Millisecond, `"4m39.001s"`},
		{0 * time.Second, `"0s"`},
		{4*time.Minute + 39*time.Second + 1*time.Millisecond, `"4m39.001s"`},
		{29*time.Hour + 39*time.Second + 1*time.Millisecond, `"29h0m39.001s"`},
	} {
		pc.appendDuration(c.src)
		if str := pc.String(); str != c.expect {
			t.Fatalf("pc is %q now, but expecting %q", str, c.expect)
		}
		pc.Reset()
	}
}

func TestPrintCtx_appendDurationSlice(t *testing.T) {
	var pc PrintCtx

	for _, c := range []struct {
		src    []time.Duration
		expect string
	}{
		{[]time.Duration{279001 * time.Nanosecond, 312 * time.Microsecond}, `["279.001µs","312µs"]`},
		// {4*time.Minute + 39*time.Second + 1*time.Millisecond, `"4m39.001s"`},
		// {29*time.Hour + 39*time.Second + 1*time.Millisecond, `"29h0m39.001s"`},
	} {
		pc.appendDurationSlice(c.src)
		if str := pc.String(); str != c.expect {
			t.Fatalf("pc is %q now, but expecting %q", str, c.expect)
		}
		pc.Reset()
	}
}

func TestPrintCtx_appendTime(t *testing.T) {
	var pc PrintCtx

	for _, c := range []struct {
		src    time.Time
		expect string
	}{
		{times.MustSmartParseTime("2000-1-1 3:0:59.001059"), `2000-01-01T03:00:59.001059Z`},
		{times.MustSmartParseTime("2023-11-01 3:0:59.001059"), `2023-11-01T03:00:59.001059Z`},
	} {
		pc.appendTime(c.src.UTC())
		if str := pc.String(); str != c.expect {
			t.Fatalf("pc is %q now, but expecting %q", str, c.expect)
		}
		pc.Reset()
	}
}

func TestPrintCtx_appendTimestamp(t *testing.T) {
	var pc PrintCtx

	defer SaveFlagsAndMod(Lempty, Ldatetimeflags|LlocalTime)() //nolint:revive // ok

	pc.jsonMode = true
	pc.noColor = false
	for _, c := range []struct {
		src    time.Time
		expect string
	}{
		{times.MustSmartParseTime("2000-1-1 3:0:59.001059"), `"03:00:59.001059Z"`},
		{times.MustSmartParseTime("2023-11-01 3:0:59.001059"), `"03:00:59.001059Z"`},
	} {
		pc.appendTimestamp(c.src)
		if str := pc.String(); str != c.expect {
			t.Fatalf("pc is %q now, but expecting %q", str, c.expect)
		}
		pc.Reset()
	}

	pc.jsonMode = false
	pc.noColor = false
	for _, c := range []struct {
		src    time.Time
		expect string
	}{
		{times.MustSmartParseTime("2000-1-1 3:0:59.001059"), `03:00:59.001059Z|`},
		{times.MustSmartParseTime("2023-11-01 3:0:59.001059"), `03:00:59.001059Z|`},
	} {
		pc.appendTimestamp(c.src)
		if str := pc.String(); str != c.expect {
			t.Fatalf("pc is %q now, but expecting %q", str, c.expect)
		}
		pc.Reset()
	}
}

func TestItoaS(t *testing.T) {
	var pc PrintCtx

	for _, c := range []struct {
		src    int64
		expect string
	}{
		{2147483647, `2147483647`},
		{-2147483648, `-2147483648`},
		{5, `5`},
		{0, `0`},
		{-1, `-1`},
		{-179, `-179`},
		{-9223372036854775808, `-9223372036854775808`},
		{9223372036854775807, `9223372036854775807`},
	} {
		itoaS(&pc, c.src)
		if str := pc.String(); str != c.expect {
			t.Fatalf("pc is %q now, but expecting %q", str, c.expect)
		}
		pc.Reset()
	}
}

func TestUtoaS(t *testing.T) {
	var pc PrintCtx

	for _, c := range []struct {
		src    uint64
		expect string
	}{
		{2147483647, `2147483647`},
		{2147483648, `2147483648`},
		{4294967295, `4294967295`},
		{0, `0`},
		{1, `1`},
		{123, `123`},
		{18446744073709551615, `18446744073709551615`},
	} {
		utoaS(&pc, c.src)
		if str := pc.String(); str != c.expect {
			t.Fatalf("pc is %q now, but expecting %q", str, c.expect)
		}
		pc.Reset()
	}
}

func TestFtoaS(t *testing.T) {
	var pc PrintCtx

	for _, c := range []struct {
		// json, noColor bool
		src    float64
		expect string
	}{
		{3.1234567890123, `3.1234567890123`},
		{3.1234567890123, `3.1234567890123`},
		{0.00001357, `0.00001357`},
		{-0.00001357, `-0.00001357`},
		{1280.00001357, `1280.00001357`},
		{-1280.00001357, `-1280.00001357`},
	} {
		ftoaS(&pc, c.src)
		if str := pc.String(); str != c.expect {
			t.Fatalf("pc is %q now, but expecting %q", str, c.expect)
		}
		pc.Reset()
	}
}

func TestCtoaS(t *testing.T) {
	var pc PrintCtx

	for _, c := range []struct {
		// json, noColor bool
		src    complex128
		expect string
	}{
		{3.1234567890123 + 0.9876543210917i, `(3.1234567890123+0.9876543210917i)`},
		{3.1234567890123 - 0.9876543210917i, `(3.1234567890123-0.9876543210917i)`},
		{0.00001357, `(0.00001357+0i)`},
		{-0.00001357, `(-0.00001357+0i)`},
		{0.00001357i, `(0+0.00001357i)`},
		{-0.00001357i, `(0-0.00001357i)`},
	} {
		ctoaS(&pc, c.src)
		if str := pc.String(); str != c.expect {
			t.Fatalf("pc is %q now, but expecting %q", str, c.expect)
		}
		pc.Reset()
	}
}

func TestBtoaS(t *testing.T) {
	var pc PrintCtx

	for _, c := range []struct {
		// json, noColor bool
		src    bool
		expect string
	}{
		{true, `true`},
		{false, `false`},
	} {
		btoaS(&pc, c.src)
		if str := pc.String(); str != c.expect {
			t.Fatalf("pc is %q now, but expecting %q", str, c.expect)
		}
		pc.Reset()
	}
}

func TestEnvTest(t *testing.T) {
	info, _ := debug.ReadBuildInfo()
	fmt.Println(info)
}

func BenchmarkPrintCtx_appendValue(b *testing.B) {
	// b.Logf("Logging with additional context at each log site. [BenchmarkAddingFields]")
	b.Run("logg/slog.PrintCtx.appendValue", func(b *testing.B) {
		b.ResetTimer()
		b.RunParallel(func(pb *testing.PB) {
			pc := newPrintCtx() // pc supports one writer for concurrency environ
			idx := 0
			for pb.Next() {
				for ix, vx := range loggArgs {
					if ix%2 == 1 {
						pc.appendValue(vx)
					}
				}
				idx++
				pc.Reset()
			}
		})
	})
	// b.Run("logg/slog.SB.appendValue", func(b *testing.B) {
	// 	b.ResetTimer()
	// 	b.RunParallel(func(pb *testing.PB) {
	// 		pc := printContextS{
	// 			SB: SB{jsonMode: true},
	// 		}
	// 		idx := 0
	// 		for pb.Next() {
	// 			for ix, vx := range loggArgs {
	// 				if ix%2 == 1 {
	// 					pc.appendValue(vx, &pc)
	// 				}
	// 			}
	// 			idx++
	// 			pc.Reset()
	// 		}
	// 	})
	// })
}

//

//

//
