package slog

import (
	"fmt"
	"io"
	"runtime/debug"
	"strings"
	"testing"
	"time"

	"github.com/hedzr/logg/slog/internal/times"
)

func TestPrintCtx_TruncateIllegal(t *testing.T) {
	defer func() {
		if e := recover(); e != nil {
			t.Log("OK")
		} else {
			t.Fatal("expect Truncate(-1) raised a panic but it was absent")
		}
	}()

	var pc PrintCtx
	pc.Truncate(-1)
}

func TestPrintCtx_GrowIllegal(t *testing.T) {
	defer func() {
		if e := recover(); e != nil {
			t.Log("OK")
		} else {
			t.Fatal("expect Grow(-1) raised a panic but it was absent")
		}
	}()

	var pc PrintCtx
	pc.Grow(-1)
}

func TestPrintCtx_pcAppendString(t *testing.T) {
	var pc PrintCtx

	pc = *newPrintCtx()

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

	pc.SetMode(ModeLogFmt)
	pc.pcAppendColon()
	if str := pc.String(); str != "()[]=" {
		t.Fatalf("pc is %q now, but expecting '()[]='", str)
	}

	// pc.jsonMode = true
	pc.SetMode(ModeJSON)
	pc.pcAppendColon()
	if str := pc.String(); str != "()[]=:" {
		t.Fatalf("pc is %q now, but expecting '()[]=:'", str)
	}

	t.Logf("Available bytes: %d", pc.Available())
}

func TestPrintCtx_pcTryQuoteValue(t *testing.T) {
	var pc PrintCtx

	t.Log("json mode")
	pc.SetMode(ModeJSON)
	// pc.jsonMode = true
	// pc.noColor = true
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
	pc.SetMode(ModeLogFmt)
	// pc.jsonMode = false
	// pc.noColor = true
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
	pc.SetMode(ModeColorful)
	// pc.jsonMode = false
	// pc.noColor = false
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

	for i, c := range []struct {
		mode Mode
		// json, noColor bool
		src, expect string
	}{
		{ModeJSON, "msg", `"msg"`},
		{ModeJSON, `msg"pc"`, `"msg\"pc\""`},
		{ModeLogFmt, "msg", `"msg"`},
		{ModeLogFmt, `msg"pc"`, `"msg\"pc\""`},
		{ModeColorful, "msg", `"msg"`},
		{ModeColorful, `msg"pc"`, `"msg\"pc\""`},
		{ModeUndefined, "msg", `"msg"`},
		{ModeUndefined, `msg"pc"`, `"msg\"pc\""`},
		// {true, false, "msg", `"msg"`},
		// {true, false, `msg"pc"`, `"msg\"pc\""`},
	} {
		// pc.jsonMode, pc.noColor = c.json, c.noColor
		pc.SetMode(c.mode)
		pc.AppendQuotedStringValue(c.src)
		if str := pc.String(); str != c.expect {
			t.Fatalf("#%5d. pc is %q now, but expecting %q", i, str, c.expect)
		}
		pc.Reset()
	}

	for i, c := range []struct {
		mode Mode
		// json, noColor bool
		src, expect string
	}{
		{ModeJSON, "msg", `msg`},
		{ModeJSON, `msg"pc"`, `msg"pc"`},
		{ModeLogFmt, "msg", `msg`},
		{ModeLogFmt, `msg"pc"`, `msg"pc"`},
		{ModeColorful, "msg", `msg`},
		{ModeColorful, `msg"pc"`, `msg"pc"`},
		{ModeUndefined, "msg", `msg`},
		{ModeUndefined, `msg"pc"`, `msg"pc"`},
		// {true, false, "msg", `"msg"`},
		// {true, false, `msg"pc"`, `"msg\"pc\""`},
	} {
		// pc.jsonMode, pc.noColor = c.json, c.noColor
		pc.SetMode(c.mode)
		pc.AppendStringValue(c.src)
		if str := pc.String(); str != c.expect {
			t.Fatalf(">%5d. pc is %q now, but expecting %q", i, str, c.expect)
		}
		pc.Reset()
	}

	// for _, c := range []struct {
	// 	json, noColor bool
	// 	src, expect   string
	// }{
	// 	{true, false, "msg", `msg`},
	// 	{true, false, `msg"pc"`, `msg"pc"`},
	// } {
	// 	pc.jsonMode, pc.noColor = c.json, c.noColor
	// 	pc.pcAppendStringValue(c.src)
	// 	if str := pc.String(); str != c.expect {
	// 		t.Fatalf("pc is %q now, but expecting %q", str, c.expect)
	// 	}
	// 	pc.Reset()
	// }
}

func TestPrintCtx_appendDuration(t *testing.T) {
	var pc PrintCtx
	pc.SetMode(ModePlain)

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
		pc.AppendDuration(c.src)
		if str := pc.String(); str != c.expect {
			t.Fatalf("pc is %q now, but expecting %q", str, c.expect)
		}
		pc.Reset()
	}
}

func TestPrintCtx_appendDurationSlice(t *testing.T) {
	var pc PrintCtx
	pc.SetMode(ModePlain)

	for _, c := range []struct {
		src    []time.Duration
		expect string
	}{
		{[]time.Duration{279001 * time.Nanosecond, 312 * time.Microsecond}, `["279.001µs","312µs"]`},
		// {4*time.Minute + 39*time.Second + 1*time.Millisecond, `"4m39.001s"`},
		// {29*time.Hour + 39*time.Second + 1*time.Millisecond, `"29h0m39.001s"`},
	} {
		pc.AppendDurationSlice(c.src)
		if str := pc.String(); str != c.expect {
			t.Fatalf("pc is %q now, but expecting %q", str, c.expect)
		}
		pc.Reset()
	}
}

func TestPrintCtx_appendTime(t *testing.T) {
	var pc PrintCtx
	pc.SetMode(ModePlain)

	for _, c := range []struct {
		src    time.Time
		expect string
	}{
		{times.MustSmartParseTime("2000-1-1 3:0:59.001059"), `2000-01-01T03:00:59.001059Z`},
		{times.MustSmartParseTime("2023-11-01 3:0:59.001059"), `2023-11-01T03:00:59.001059Z`},
	} {
		pc.SetMode(ModeColorful)
		pc.AppendTime(c.src.UTC())
		if str := pc.String(); str != c.expect {
			t.Fatalf("pc is %q now, but expecting %q", str, c.expect)
		}
		pc.Reset()
	}
}

func TestPrintCtx_appendTimestamp(t *testing.T) {
	var pc PrintCtx

	defer SaveFlagsAndMod(Lempty, Ldatetimeflags|LlocalTime)() //nolint:revive // ok

	pc.SetMode(ModeJSON)
	// pc.jsonMode = true
	// pc.noColor = false
	for _, c := range []struct {
		src    time.Time
		expect string
	}{
		{times.MustSmartParseTime("2000-1-1 3:0:59.001059"), `"03:00:59.001059Z"`},
		{times.MustSmartParseTime("2023-11-01 3:0:59.001059"), `"03:00:59.001059Z"`},
	} {
		pc.AppendTimestamp(c.src)
		if str := pc.String(); str != c.expect {
			t.Fatalf("pc is %q now, but expecting %q", str, c.expect)
		}
		pc.Reset()
	}

	pc.SetMode(ModeColorful)
	// pc.jsonMode = false
	// pc.noColor = false
	for _, c := range []struct {
		src    time.Time
		expect string
	}{
		{times.MustSmartParseTime("2000-1-1 3:0:59.001059"), `03:00:59.001059Z|`},
		{times.MustSmartParseTime("2023-11-01 3:0:59.001059"), `03:00:59.001059Z|`},
	} {
		pc.AppendTimestamp(c.src)
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
		pc.SetMode(ModeColorful)
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
		pc.SetMode(ModeColorful)
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
		pc.SetMode(ModeColorful)
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

func TestAddXXX(t *testing.T) {
	var pc PrintCtx

	subTest := func(t *testing.T, pc *PrintCtx) {
		pc.AddInt64("int64", int64(-101))
		pc.AddRune(',')
		pc.AddInt32("int32", int32(-32))
		pc.AddRune(',')
		pc.AddInt16("int16", int16(-16))
		pc.AddRune(',')
		pc.AddInt8("int8", int8(-8))
		pc.AddRune(',')

		pc.AddUint64("uint64", uint64(101))
		pc.AddRune(',')
		pc.AddUint32("uint32", uint32(32))
		pc.AddRune(',')
		pc.AddUint16("uint16", uint16(16))
		pc.AddRune(',')
		pc.AddUint8("uint8", uint8(8))
		pc.AddRune(',')
		pc.AddUint("uint", uint(65536))
		pc.AddRune(',')

		pc.AddFloat32("float32", float32(65536.32768))
		pc.AddRune(',')
		pc.AddFloat64("float64", float64(3.13))
		pc.AddRune(',')
		pc.AddFloat64("float64", float64(65536.32768))
		pc.AddRune(',')

		pc.AddComplex64("complex64", 65536.32768+2.718i)
		pc.AddRune(',')
		pc.AddComplex128("complex128", 3.13-2.718i)
		pc.AddRune(',')
		pc.AddComplex128("complex128", 65536.32768+2.718i)
		pc.AddRune(',')

		pc.AddBool("bool", false)
		pc.AddRune(',')

		pc.AddString("ending-string", "modified")
		pc.AddRune(',')

		t.Logf("pc: %v", pc.String())
	}

	l := newentry(nil, WithJSONMode(true))
	pc.internalsetentry(l)
	subTest(t, &pc)

	l = newentry(nil, WithJSONMode(true))
	pc.internalsetentry(l)
	subTest(t, &pc)

	l = newentry(nil, WithJSONMode(true))
	pc.internalsetentry(l)
	pc.appendValue(struct{}{})
	subTest(t, &pc)

	l = newentry(nil, WithColorMode(true), WithJSONMode(true))
	// l.useColor = true
	pc.internalsetentry(l)
	// pc.noColor = false
	pc.AddPrefixedString("pre-", "want", "vn")
	pc.AddPrefixedInt("pre-", "want", 92)
	pc.appendValue(nil)
	pc.appendValue(&amS{})
	pc.appendValue(&omS{})
	pc.appendValue([]time.Time{time.Now()})
	pc.appendValue([]time.Duration{time.Second})
	pc.appendValue(io.ErrNoProgress)
	var ss strings.Builder
	pc.appendValue(&ss)
	pc.appendValue([]byte("hello"))
	pc.appendValue([]string{"hello"})
	pc.appendValue([]bool{true})
	pc.appendValue([]int{23})
	pc.appendValue([]int8{23})
	pc.appendValue([]int16{23})
	pc.appendValue([]int32{23})
	pc.appendValue([]int64{23})
	pc.appendValue([]uint{23})
	pc.appendValue([]uint8{23})
	pc.appendValue([]uint16{23})
	pc.appendValue([]uint32{23})
	pc.appendValue([]uint64{23})
	pc.appendValue([]float64{23})
	pc.appendValue([]float32{23})
	pc.appendValue([]complex128{23})
	pc.appendValue([]complex64{23})
	subTest(t, &pc)
}

type amS struct{}

func (s *amS) MarshalSlogArray(enc *PrintCtx) error { return nil }

type omS struct{}

func (s *omS) MarshalSlogObject(enc *PrintCtx) error { return nil }

func TestAppendXXX(t *testing.T) {
	var pc PrintCtx

	pc.AppendByte('^')
	pc.AppendRune(',')

	pc.AppendBytes([]byte("yes"))
	pc.AppendRune(',')
	pc.AppendRunes([]rune("no"))
	pc.AppendRune('$')

	t.Logf("pc: %v", pc.String())
}

func TestSource_ToGroup(t *testing.T) {
	var pc PrintCtx
	t.Log(pc.source().toGroup())

	t.Log(stack(0, 0))
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
