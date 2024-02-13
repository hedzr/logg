package slog

import (
	"io"
	"testing"
)

func TestDualWriter(t *testing.T) {
	testDualWriter(t)
	t.Log(" TestDualWriter OK")
}

func testDualWriter(t *testing.T) {
	dw := newDualWriter()

	defer func(dw *dualWriter) {
		err := dw.Close()
		if err != nil {
			t.Logf("has error in closing: %v", err)
		}
	}(dw)

	_, _ = dw.Write([]byte("hello\n"))

	var d1 discard
	_, _ = d1.Write([]byte("hello"))

	dw.Set(d1)
	dw.SetWriter(d1)
	dw.SetErrorWriter(d1)

	dw.Add(d1)
	dw.AddErrorWriter(d1)
	dw.AddLevelWriter(ErrorLevel, d1)
	dw.RemoveLevelWriter(ErrorLevel, d1)

	dw.ResetLevelWriter(ErrorLevel)
	dw.ResetLevelWriters()
	dw.Clear()
	dw.Reset()
}

func TestLWs(t *testing.T) {
	var dw LWs
	defer dw.Close()

	_, _ = dw.Write([]byte("hello"))

	t.Log(" TestLWs OK")
}

func TestLogwr(t *testing.T) {
	dw := &logwr{io.Discard}
	defer dw.Close()

	_, _ = dw.Write([]byte("hello"))

	t.Log(" TestLogwr OK")
}

// func TestFilewr(t *testing.T) {
// 	dw := &filewr{os.Stdout}
// 	defer dw.Close()
//
// 	_, _ = dw.WriteString("hello\n")
//
// 	t.Log(" TestFilewr OK")
// }

func TestNewLogWriter(t *testing.T) {
	dw := NewLogWriter(io.Discard)
	defer dw.Close()

	_, _ = dw.Write([]byte("hello"))

	t.Log(" TestNewLogWriter OK")
}

func TestNewFileWriter(t *testing.T) {
	dw := NewFileWriter("/dev/null")
	defer dw.Close()

	_, _ = dw.WriteString("hello")

	t.Log(" TestNewFileWriter OK")
}
