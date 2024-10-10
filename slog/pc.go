package slog

import (
	"bytes"
	"encoding"
	"errors"
	"fmt"
	"io"
	"strconv"
	"sync"
	"time"
	"unicode/utf8"

	"github.com/hedzr/is/term/color"
)

var printCtxPool = sync.Pool{New: func() any {
	return newPrintCtx()

	// return &PrintCtx{
	// 	buf:      make([]byte, 0, 1024),
	// 	noQuoted: true,
	// 	clr:      clrBasic,
	// 	bg:       clrNone,
	// }
}}

// func newPrintCtxAsAny() any { return newPrintCtx() }

func newPrintCtx() *PrintCtx {
	return &PrintCtx{
		buf:      make([]byte, 0, 1024),
		noQuoted: true,
		clr:      clrBasic,
		bg:       clrNone,
	}
}

// PrintCtx when formatting logging line in text logger
type PrintCtx struct {
	buf      []byte // contents are the bytes buf[off : len(buf)]
	off      int    // read at &buf[off], write at &buf[len(buf)]
	lastRead readOp // last read operation, so that Unread* can work correctly.

	noQuoted bool   // should quote the string values? default is YES
	jsonMode bool   // should print out the logging with JSON format? default is NO.
	noColor  bool   // use ansi escape sequences in console/terminal? default is ON.
	layout   string // time layout for formatting
	utcTime  int    // non-set(0), local(1) or utc(2) time? default is local time mode.

	lvl        Level
	msg        string
	firstLine  string
	restLines  string
	eol        bool
	kvps       Attrs
	clr, bg    color.Color
	now        time.Time
	stackFrame uintptr

	prefix        string
	inGroupedMode bool

	// curdir string

	valueStringer ValueStringer
}

func (s *PrintCtx) source() Source { return getpcsource(s.stackFrame) }

func (s *PrintCtx) setentry(e *Entry) {
	s.buf = s.buf[:0]

	s.jsonMode = e.useJSON
	useColor := e.useColor
	if e.useJSON && useColor {
		useColor = false
	}
	s.noColor = !useColor

	s.layout = e.timeLayout
	s.utcTime = e.modeUTC
	s.valueStringer = e.valueStringer

	s.lvl = e.level
	s.kvps = e.attrs
}

func (s *PrintCtx) set(e *Entry, lvl Level, timestamp time.Time, stackFrame uintptr, msg string, kvps Attrs) {
	s.setentry(e)

	s.lvl = lvl
	s.now = timestamp
	s.stackFrame = stackFrame
	s.msg = msg
	s.kvps = kvps
}

//
//
//

// smallBufferSize is an initial allocation minimal capacity.
const smallBufferSize = 64

// The readOp constants describe the last action performed on
// the buffer, so that UnreadRune and UnreadByte can check for
// invalid usage. opReadRuneX constants are chosen such that
// converted to int they correspond to the rune size that was read.
type readOp int8

// Don't use iota for these, as the values need to correspond with the
// names and comments, which is easier to see when being explicit.
const (
	opRead      readOp = -1 // Any other read operation.
	opInvalid   readOp = 0  // Non-read operation.
	opReadRune1 readOp = 1  // Read rune of size 1.
	opReadRune2 readOp = 2  // Read rune of size 2.
	opReadRune3 readOp = 3  // Read rune of size 3.
	opReadRune4 readOp = 4  // Read rune of size 4.
)

// ErrTooLarge is passed to panic if memory cannot be allocated to store data in a buffer.
var ErrTooLarge = errors.New("logg/slog.PrintCtx: too large")
var errNegativeRead = errors.New("logg/slog.PrintCtx: reader returned negative count from Read")

const maxInt = int(^uint(0) >> 1)

// Bytes returns a slice of length b.Len() holding the unread portion of the buffer.
// The slice is valid for use only until the next buffer modification (that is,
// only until the next call to a method like Read, Write, Reset, or Truncate).
// The slice aliases the buffer content at least until the next buffer modification,
// so immediate changes to the slice will affect the result of future reads.
func (s *PrintCtx) Bytes() []byte { return s.buf[s.off:] }

// AvailableBuffer returns an empty buffer with b.Available() capacity.
// This buffer is intended to be appended to and
// passed to an immediately succeeding Write call.
// The buffer is only valid until the next write operation on b.
func (s *PrintCtx) AvailableBuffer() []byte { return s.buf[len(s.buf):] }

// String returns the contents of the unread portion of the buffer
// as a string. If the Buffer is a nil pointer, it returns "<nil>".
//
// To build strings more efficiently, see the strings.Builder type.
func (s *PrintCtx) String() string {
	if s == nil {
		// Special case, useful in debugging.
		return "<nil>"
	}
	return string(s.buf[s.off:])
}

// empty reports whether the unread portion of the buffer is empty.
func (s *PrintCtx) empty() bool { return len(s.buf) <= s.off }

// Len returns the number of bytes of the unread portion of the buffer;
// b.Len() == len(b.Bytes()).
func (s *PrintCtx) Len() int { return len(s.buf) - s.off }

// Cap returns the capacity of the buffer's underlying byte slice, that is, the
// total space allocated for the buffer's data.
func (s *PrintCtx) Cap() int { return cap(s.buf) }

// Available returns how many bytes are unused in the buffer.
func (s *PrintCtx) Available() int { return cap(s.buf) - len(s.buf) }

// Truncate discards all but the first n unread bytes from the buffer
// but continues to use the same allocated storage.
// It panics if n is negative or greater than the length of the buffer.
func (s *PrintCtx) Truncate(n int) {
	if n == 0 {
		s.Reset()
		return
	}
	s.lastRead = opInvalid
	if n < 0 || n > s.Len() {
		panic("logg/slog.PrintCtx: truncation out of range")
	}
	s.buf = s.buf[:s.off+n]
}

// Reset resets the buffer to be empty,
// but it retains the underlying storage for use by future writes.
// Reset is the same as Truncate(0).
func (s *PrintCtx) Reset() {
	s.buf = s.buf[:0]
	s.off = 0
	s.lastRead = opInvalid
}

// tryGrowByReslice is an inlineable version of grow for the fast-case where the
// internal buffer only needs to be resliced.
// It returns the index where bytes should be written and whether it succeeded.
func (s *PrintCtx) tryGrowByReslice(n int) (int, bool) {
	if l := len(s.buf); n <= cap(s.buf)-l {
		s.buf = s.buf[:l+n]
		return l, true
	}
	return 0, false
}

// grow grows the buffer to guarantee space for n more bytes.
// It returns the index where bytes should be written.
// If the buffer can't grow it will panic with ErrTooLarge.
func (s *PrintCtx) grow(n int) int {
	m := s.Len()
	// If buffer is empty, reset to recover space.
	if m == 0 && s.off != 0 {
		s.Reset()
	}
	// Try to grow by means of a reslice.
	if i, ok := s.tryGrowByReslice(n); ok {
		return i
	}
	if s.buf == nil && n <= smallBufferSize {
		s.buf = make([]byte, n, smallBufferSize)
		return 0
	}
	c := cap(s.buf)
	if n <= c/2-m {
		// We can slide things down instead of allocating a new
		// slice. We only need m+n <= c to slide, but
		// we instead let capacity get twice as large so we
		// don't spend all our time copying.
		copy(s.buf, s.buf[s.off:])
	} else if c > maxInt-c-n {
		panic(ErrTooLarge)
	} else {
		// Add s.off to account for s.buf[:s.off] being sliced off the front.
		s.buf = growSlice(s.buf[s.off:], s.off+n)
	}
	// Restore b.off and len(b.buf).
	s.off = 0
	s.buf = s.buf[:m+n]
	return m
}

// Grow grows the buffer's capacity, if necessary, to guarantee space for
// another n bytes. After Grow(n), at least n bytes can be written to the
// buffer without another allocation.
// If n is negative, Grow will panic.
// If the buffer can't grow it will panic with ErrTooLarge.
func (s *PrintCtx) Grow(n int) {
	if n < 0 {
		panic("logg/slog.PrintCtx.Grow: negative count")
	}
	m := s.grow(n)
	s.buf = s.buf[:m]
}

// Write appends the contents of p to the buffer, growing the buffer as
// needed. The return value n is the length of p; err is always nil. If the
// buffer becomes too large, Write will panic with ErrTooLarge.
func (s *PrintCtx) Write(p []byte) (n int, err error) {
	s.lastRead = opInvalid
	m, ok := s.tryGrowByReslice(len(p))
	if !ok {
		m = s.grow(len(p))
	}
	return copy(s.buf[m:], p), nil
}

// WriteString appends the contents of s to the buffer, growing the buffer as
// needed. The return value n is the length of s; err is always nil. If the
// buffer becomes too large, WriteString will panic with ErrTooLarge.
func (s *PrintCtx) WriteString(str string) (n int, err error) {
	s.lastRead = opInvalid
	m, ok := s.tryGrowByReslice(len(str))
	if !ok {
		m = s.grow(len(str))
	}
	return copy(s.buf[m:], str), nil
}

// MinRead is the minimum slice size passed to a Read call by
// Buffer.ReadFrom. As long as the Buffer has at least MinRead bytes beyond
// what is required to hold the contents of r, ReadFrom will not grow the
// underlying buffer.
const MinRead = 512

// ReadFrom reads data from r until EOF and appends it to the buffer, growing
// the buffer as needed. The return value n is the number of bytes read. Any
// error except io.EOF encountered during the read is also returned. If the
// buffer becomes too large, ReadFrom will panic with ErrTooLarge.
func (s *PrintCtx) ReadFrom(r io.Reader) (n int64, err error) {
	s.lastRead = opInvalid
	for {
		i := s.grow(MinRead)
		s.buf = s.buf[:i]
		m, e := r.Read(s.buf[i:cap(s.buf)])
		if m < 0 {
			panic(errNegativeRead)
		}

		s.buf = s.buf[:i+m]
		n += int64(m)
		if e == io.EOF {
			return n, nil // e is EOF, so return nil explicitly
		}
		if e != nil {
			return n, e
		}
	}
}

// growSlice grows b by n, preserving the original content of b.
// If the allocation fails, it panics with ErrTooLarge.
func growSlice(b []byte, n int) []byte {
	defer func() {
		if recover() != nil {
			panic(ErrTooLarge)
		}
	}()
	// TODO(http://golang.org/issue/51462): We should rely on the append-make
	// pattern so that the compiler can call runtime.growslice. For example:
	//	return append(b, make([]byte, n)...)
	// This avoids unnecessary zero-ing of the first len(b) bytes of the
	// allocated slice, but this pattern causes b to escape onto the heap.
	//
	// Instead use the append-make pattern with a nil slice to ensure that
	// we allocate buffers rounded up to the closest size class.
	c := len(b) + n // ensure enough space for n elements
	if c < 2*cap(b) {
		// The growth rate has historically always been 2x. In the future,
		// we could rely purely on append to determine the growth rate.
		c = 2 * cap(b)
	}
	b2 := append([]byte(nil), make([]byte, c)...)
	copy(b2, b)
	return b2[:len(b)]
}

// WriteTo writes data to w until the buffer is drained or an error occurs.
// The return value n is the number of bytes written; it always fits into an
// int, but it is int64 to match the io.WriterTo interface. Any error
// encountered during the write is also returned.
func (s *PrintCtx) WriteTo(w io.Writer) (n int64, err error) {
	s.lastRead = opInvalid
	if nBytes := s.Len(); nBytes > 0 {
		m, e := w.Write(s.buf[s.off:])
		if m > nBytes {
			panic("logg/slog.PrintCtx.WriteTo: invalid Write count")
		}
		s.off += m
		n = int64(m)
		if e != nil {
			return n, e
		}
		// all bytes should have been written, by definition of
		// Write method in io.Writer
		if m != nBytes {
			return n, io.ErrShortWrite
		}
	}
	// Buffer is now empty; reset.
	s.Reset()
	return n, nil
}

// WriteByte appends the byte c to the buffer, growing the buffer as needed.
// The returned error is always nil, but is included to match bufio.Writer's
// WriteByte. If the buffer becomes too large, WriteByte will panic with
// ErrTooLarge.
func (s *PrintCtx) WriteByte(c byte) error {
	s.lastRead = opInvalid
	m, ok := s.tryGrowByReslice(1)
	if !ok {
		m = s.grow(1)
	}
	s.buf[m] = c
	return nil
}

// WriteRune appends the UTF-8 encoding of Unicode code point r to the
// buffer, returning its length and an error, which is always nil but is
// included to match bufio.Writer's WriteRune. The buffer is grown as needed;
// if it becomes too large, WriteRune will panic with ErrTooLarge.
func (s *PrintCtx) WriteRune(r rune) (n int, err error) {
	// Compare as uint32 to correctly handle negative runes.
	if uint32(r) < utf8.RuneSelf {
		s.WriteByte(byte(r))
		return 1, nil
	}
	s.lastRead = opInvalid
	m, ok := s.tryGrowByReslice(utf8.UTFMax)
	if !ok {
		m = s.grow(utf8.UTFMax)
	}
	s.buf = utf8.AppendRune(s.buf[:m], r)
	return len(s.buf) - m, nil

	// s.buf= append(s.buf, []byte(r)...)
}

// Read reads the next len(p) bytes from the buffer or until the buffer
// is drained. The return value n is the number of bytes read. If the
// buffer has no data to return, err is io.EOF (unless len(p) is zero);
// otherwise it is nil.
func (s *PrintCtx) Read(p []byte) (n int, err error) {
	s.lastRead = opInvalid
	if s.empty() {
		// Buffer is empty, reset to recover space.
		s.Reset()
		if len(p) == 0 {
			return 0, nil
		}
		return 0, io.EOF
	}
	n = copy(p, s.buf[s.off:])
	s.off += n
	if n > 0 {
		s.lastRead = opRead
	}
	return n, nil
}

// Next returns a slice containing the next n bytes from the buffer,
// advancing the buffer as if the bytes had been returned by Read.
// If there are fewer than n bytes in the buffer, Next returns the entire buffer.
// The slice is only valid until the next call to a read or write method.
func (s *PrintCtx) Next(n int) []byte {
	s.lastRead = opInvalid
	m := s.Len()
	if n > m {
		n = m
	}
	data := s.buf[s.off : s.off+n]
	s.off += n
	if n > 0 {
		s.lastRead = opRead
	}
	return data
}

// ReadByte reads and returns the next byte from the buffer.
// If no byte is available, it returns error io.EOF.
func (s *PrintCtx) ReadByte() (byte, error) {
	if s.empty() {
		// Buffer is empty, reset to recover space.
		s.Reset()
		return 0, io.EOF
	}
	c := s.buf[s.off]
	s.off++
	s.lastRead = opRead
	return c, nil
}

// ReadRune reads and returns the next UTF-8-encoded
// Unicode code point from the buffer.
// If no bytes are available, the error returned is io.EOF.
// If the bytes are an erroneous UTF-8 encoding, it
// consumes one byte and returns U+FFFD, 1.
func (s *PrintCtx) ReadRune() (r rune, size int, err error) {
	if s.empty() {
		// Buffer is empty, reset to recover space.
		s.Reset()
		return 0, 0, io.EOF
	}
	c := s.buf[s.off]
	if c < utf8.RuneSelf {
		s.off++
		s.lastRead = opReadRune1
		return rune(c), 1, nil
	}
	r, n := utf8.DecodeRune(s.buf[s.off:])
	s.off += n
	s.lastRead = readOp(n)
	return r, n, nil
}

// UnreadRune unreads the last rune returned by ReadRune.
// If the most recent read or write operation on the buffer was
// not a successful ReadRune, UnreadRune returns an error.  (In this regard
// it is stricter than UnreadByte, which will unread the last byte
// from any read operation.)
func (s *PrintCtx) UnreadRune() error {
	if s.lastRead <= opInvalid {
		return errors.New("logg/slog.PrintCtx: UnreadRune: previous operation was not a successful ReadRune")
	}
	if s.off >= int(s.lastRead) {
		s.off -= int(s.lastRead)
	}
	s.lastRead = opInvalid
	return nil
}

var errUnreadByte = errors.New("logg/slog.PrintCtx: UnreadByte: previous operation was not a successful read")

// UnreadByte unreads the last byte returned by the most recent successful
// read operation that read at least one byte. If a write has happened since
// the last read, if the last read returned an error, or if the read read zero
// bytes, UnreadByte returns an error.
func (s *PrintCtx) UnreadByte() error {
	if s.lastRead == opInvalid {
		return errUnreadByte
	}
	s.lastRead = opInvalid
	if s.off > 0 {
		s.off--
	}
	return nil
}

// ReadBytes reads until the first occurrence of delim in the input,
// returning a slice containing the data up to and including the delimiter.
// If ReadBytes encounters an error before finding a delimiter,
// it returns the data read before the error and the error itself (often io.EOF).
// ReadBytes returns err != nil if and only if the returned data does not end in
// delim.
func (s *PrintCtx) ReadBytes(delim byte) (line []byte, err error) {
	slice, err := s.readSlice(delim)
	// return a copy of slice. The buffer's backing array may
	// be overwritten by later calls.
	line = append(line, slice...)
	return line, err
}

// readSlice is like ReadBytes but returns a reference to internal buffer data.
func (s *PrintCtx) readSlice(delim byte) (line []byte, err error) {
	i := bytes.IndexByte(s.buf[s.off:], delim)
	end := s.off + i + 1
	if i < 0 {
		end = len(s.buf)
		err = io.EOF
	}
	line = s.buf[s.off:end]
	s.off = end
	s.lastRead = opRead
	return line, err
}

// ReadString reads until the first occurrence of delim in the input,
// returning a string containing the data up to and including the delimiter.
// If ReadString encounters an error before finding a delimiter,
// it returns the data read before the error and the error itself (often io.EOF).
// ReadString returns err != nil if and only if the returned data does not end
// in delim.
func (s *PrintCtx) ReadString(delim byte) (line string, err error) {
	slice, err := s.readSlice(delim)
	return string(slice), err
}

// NewPrintCtx creates and initializes a new Buffer using buf as its
// initial contents. The new Buffer takes ownership of buf, and the
// caller should not use buf after this call. NewBuffer is intended to
// prepare a Buffer to read existing data. It can also be used to set
// the initial size of the internal buffer for writing. To do that,
// buf should have the desired capacity but a length of zero.
//
// In most cases, new(Buffer) (or just declaring a Buffer variable) is
// sufficient to initialize a Buffer.
func NewPrintCtx(buf []byte) *PrintCtx { return &PrintCtx{buf: buf} }

// NewPrintCtxString creates and initializes a new Buffer using string s as its
// initial contents. It is intended to prepare a buffer to read an existing
// string.
//
// In most cases, new(Buffer) (or just declaring a Buffer variable) is
// sufficient to initialize a Buffer.
func NewPrintCtxString(s string) *PrintCtx {
	return &PrintCtx{buf: []byte(s)}
}

//
//
//

// func (s *PrintCtx) addRune(r rune) { s.WriteRune(r) }

// func (s *PrintCtx) addString(name string, value string) {
// 	// s.Grow(len(name)*3 + 1 + len(value)*3)
// 	s.pcAppendStringKey(name)
// 	s.pcAppendColon()
// 	s.pcAppendQuotedStringValue(value)
// }

func (s *PrintCtx) AddRune(r rune) {
	s.pcAppendRune(r)
}

func (s *PrintCtx) AddInt64(name string, value int64) {
	// s.Grow(len(name)*3 + 1 + 10)
	s.pcAppendStringKey(name)
	s.pcAppendColon()
	// s.pcAppendStringValue(intToString(value))
	itoaS(s, value)
}

func (s *PrintCtx) AddInt32(name string, value int32) {
	// s.Grow(len(name)*3 + 1 + 10)
	s.pcAppendStringKey(name)
	s.pcAppendColon()
	// s.pcAppendStringValue(intToString(value))
	itoaS(s, value)
}

func (s *PrintCtx) AddInt16(name string, value int16) {
	// s.Grow(len(name)*3 + 1 + 10)
	s.pcAppendStringKey(name)
	s.pcAppendColon()
	// s.pcAppendStringValue(intToString(value))
	itoaS(s, value)
}

func (s *PrintCtx) AddInt8(name string, value int8) {
	// s.Grow(len(name)*3 + 1 + 10)
	s.pcAppendStringKey(name)
	s.pcAppendColon()
	// s.pcAppendStringValue(intToString(value))
	itoaS(s, value)
}

func (s *PrintCtx) AddInt(name string, value int) {
	// s.Grow(len(name)*3 + 1 + 10)
	s.pcAppendStringKey(name)
	s.pcAppendColon()
	// s.pcAppendStringValue(intToString(value))
	itoaS(s, value)
}

func (s *PrintCtx) AddPrefixedInt(prefix, name string, value int) {
	// s.Grow(len(name)*3 + 1 + 10)
	s.pcAppendStringKeyPrefixed(name, prefix)
	s.pcAppendColon()
	// s.pcAppendStringValue(intToString(value))
	itoaS(s, value)
}

func (s *PrintCtx) AddUint64(name string, value uint64) {
	// s.Grow(len(name)*3 + 1 + 10)
	s.pcAppendStringKey(name)
	s.pcAppendColon()
	// s.pcAppendStringValue(uintToString(value))
	utoaS(s, value)
}

func (s *PrintCtx) AddUint32(name string, value uint32) {
	// s.Grow(len(name)*3 + 1 + 10)
	s.pcAppendStringKey(name)
	s.pcAppendColon()
	// s.pcAppendStringValue(uintToString(value))
	utoaS(s, value)
}

func (s *PrintCtx) AddUint16(name string, value uint16) {
	// s.Grow(len(name)*3 + 1 + 10)
	s.pcAppendStringKey(name)
	s.pcAppendColon()
	// s.pcAppendStringValue(uintToString(value))
	utoaS(s, value)
}

func (s *PrintCtx) AddUint8(name string, value uint8) {
	// s.Grow(len(name)*3 + 1 + 10)
	s.pcAppendStringKey(name)
	s.pcAppendColon()
	// s.pcAppendStringValue(uintToString(value))
	utoaS(s, value)
}

func (s *PrintCtx) AddUint(name string, value uint) {
	// s.Grow(len(name)*3 + 1 + 10)
	s.pcAppendStringKey(name)
	s.pcAppendColon()
	// s.pcAppendStringValue(uintToString(value))
	utoaS(s, value)
}

func (s *PrintCtx) AddFloat64(name string, value float64) {
	// s.Grow(len(name)*3 + 1 + 10)
	s.pcAppendStringKey(name)
	s.pcAppendColon()
	// s.pcAppendStringValue(floatToString(value))
	ftoaS(s, value)
}

func (s *PrintCtx) AddFloat32(name string, value float32) {
	// s.Grow(len(name)*3 + 1 + 10)
	s.pcAppendStringKey(name)
	s.pcAppendColon()
	// s.pcAppendStringValue(floatToString(value))
	ftoaS(s, value)
}

func (s *PrintCtx) AddComplex128(name string, value complex128) {
	// s.Grow(len(name)*3 + 1 + 10)
	s.pcAppendStringKey(name)
	s.pcAppendColon()
	// s.pcAppendStringValue(complexToString(value))
	ctoaS(s, value)
}

func (s *PrintCtx) AddComplex64(name string, value complex64) {
	// s.Grow(len(name)*3 + 1 + 32)
	s.pcAppendStringKey(name)
	s.pcAppendColon()
	// s.pcAppendStringValue(complexToString(value))
	ctoaS(s, value)
}

func (s *PrintCtx) AddBool(name string, value bool) {
	// s.Grow(len(name)*3 + 1 + 5)
	s.pcAppendStringKey(name)
	s.pcAppendColon()
	// s.pcAppendStringValue(boolToString(value))
	btoaS(s, value)
}

func (s *PrintCtx) AddString(name string, value string) {
	// s.Grow(len(name)*3 + 1 + 10)
	s.pcAppendStringKey(name)
	s.pcAppendColon()
	// s.pcAppendStringValue(intToString(value))
	if s.noColor {
		s.pcAppendQuotedStringValue(value)
	} else {
		s.pcAppendString(value)
	}
}

func (s *PrintCtx) AddPrefixedString(prefix, name string, value string) {
	// s.Grow(len(name)*3 + 1 + 10)
	s.pcAppendStringKeyPrefixed(name, prefix)
	s.pcAppendColon()
	// s.pcAppendStringValue(intToString(value))
	if s.noColor {
		s.pcAppendQuotedStringValue(value)
	} else {
		s.pcAppendString(value)
	}
}

func (s *PrintCtx) AppendRune(value rune) {
	s.pcAppendRune(value)
}

func (s *PrintCtx) AppendByte(value byte) {
	s.pcAppendByte(value)
}

func (s *PrintCtx) AppendBytes(value []byte) {
	_, err := s.Write(value)
	if err != nil {
		hintInternal(err, "PrintCtx.AppendBytes failed")
	}
}

func (s *PrintCtx) AppendRunes(value []rune) {
	s.buf = append(s.buf, []byte(string(value))...)
}

//

//

//

func (s *PrintCtx) AppendInt(val int) {
	itoaS(s, val)
	// s.WriteString(intToStringEx(val, 10))
}

//

func (s *PrintCtx) preCheck() {
	// if s.jsonMode {
	// 	s.WriteRune('{')
	// }
}

// func (s *PrintCtx) postCheck() {
// 	// if s.jsonMode {
// 	// 	s.WriteRune('}')
// 	// }
// }

//

func (s *PrintCtx) pcAppendByte(b byte) {
	s.checkerr(s.WriteByte(b))
}

func (s *PrintCtx) pcAppendRune(r rune) {
	s.preCheck()
	_, err := s.WriteRune(r)
	if err != nil {
		hintInternal(err, "PrintCtx.pcAppendRune failed")
	}
}

func (s *PrintCtx) pcTryQuoteValue(val string) {
	s.preCheck()
	if s.noColor {
		// return strconv.Quote(val)
		s.pcAppendByte('"')
		s.appendEscapedJSONString(val)
		s.pcAppendByte('"')
	} else {
		s.pcAppendStringValue(val)
	}
}

func (s *PrintCtx) pcQuoteValue(val string) {
	s.pcAppendByte('"')
	s.appendEscapedJSONString(val)
	s.pcAppendByte('"')
}

func (s *PrintCtx) pcAppendColon() {
	s.preCheck()
	if s.jsonMode {
		s.pcAppendByte(':')
	} else {
		s.pcAppendByte('=')
	}
}

func (s *PrintCtx) pcAppendComma() {
	if s.jsonMode {
		s.pcAppendByte(',')
	} else {
		s.pcAppendByte(' ')
	}
}

func (s *PrintCtx) pcAppendString(str string) {
	s.preCheck()
	_, err := s.WriteString(str)
	if err != nil {
		hintInternal(err, "PrintCtx.pcAppendString failed")
	}
}

// pcAppendStringValue append string without quotes, the string represents a value
func (s *PrintCtx) pcAppendStringValue(str string) {
	s.preCheck()
	s.WriteString(str)
}

// pcAppendQuotedStringValue append string with quotes always, the string represents a value
func (s *PrintCtx) pcAppendQuotedStringValue(str string) {
	s.preCheck()
	s.WriteRune('"')
	s.appendEscapedJSONString(str)
	s.WriteRune('"')
}

// pcAppendStringKey appends string with quotes (in json mode), the string represents a key.
//
// If a PrintCtx is in non json mode, the key shouldn't wrapped with quotes.
func (s *PrintCtx) pcAppendStringKey(str string) {
	s.preCheck()
	if s.jsonMode {
		// s.WriteString(strconv.Quote(str))
		// s.Grow(2 + len([]byte(str)))
		s.checkerr(s.WriteByte('"'))
		s.WriteString(str)
		s.checkerr(s.WriteByte('"'))
	} else {
		s.WriteString(str)
	}
}

func (s *PrintCtx) pcAppendStringKeyPrefixed(str, prefix string) {
	s.preCheck()
	if s.jsonMode {
		// s.WriteString(strconv.Quote(str))
		// s.Grow(2 + len([]byte(str)))
		s.checkerr(s.WriteByte('"'))
		s.WriteString(prefix)
		s.checkerr(s.WriteByte('.'))
		s.WriteString(str)
		s.checkerr(s.WriteByte('"'))
	} else {
		s.WriteString(prefix)
		s.checkerr(s.WriteByte('.'))
		s.WriteString(str)
	}
}

// func (s *PrintCtx) appendRune(val rune) {
// 	s.preCheck()
// 	s.WriteRune(val)
// }

// func (s *PrintCtx) appendRunes(val []rune) {
// 	s.pcAppendString(string(val))
// }

// func (s *PrintCtx) appendRuneValue(val rune) {
// 	s.pcAppendStringValue(string(val))
// }

// func (s *PrintCtx) appendRunesValue(val []rune) {
// 	s.pcAppendStringValue(string(val))
// }

func (s *PrintCtx) appendValue(val any) {
	switch z := val.(type) {
	case nil:
		s.pcAppendStringValue("<nil>")

	case ObjectSerializer:
		// pc.useColor = !s.noColor
		// pc.clr = color.FgDarkColor
		// pc.bg = clrNone
		z.SerializeValueTo(s)
	// z.SerializeValueTo(pc, prefix, inGrouping, !s.noColor, color.FgDarkColor, clrNone)

	case ArrayMarshaller:
		if err := z.MarshalSlogArray(s); err != nil {
			hintInternal(err, "MarshalLogArray failed")
			break
		}

	case ObjectMarshaller:
		if err := z.MarshalSlogObject(s); err != nil {
			hintInternal(err, "MarshalLogObject failed")
			break
		}

	case time.Duration:
		s.appendDuration(z)
	case time.Time:
		s.appendTime(z)
	case []time.Time:
		s.appendTimeSlice(z)
	case []time.Duration:
		s.appendDurationSlice(z)

	case Level:
		s.pcQuoteValue(z.String())

	case error:
		s.pcTryQuoteValue(z.Error())

	case ToString:
		s.pcQuoteValue(z.ToString())

	case Stringer:
		s.pcQuoteValue(z.String())

	case string:
		s.pcQuoteValue(z)

	case bool:
		btoaS(s, z)

	case []byte:
		s.appendBytes(z)

	case []string:
		s.appendStringSlice(z)

	case []bool:
		s.appendBoolSlice(z)

	case []int:
		intSliceTo(s, z)
	case []int8:
		intSliceTo(s, z)
	case []int16:
		intSliceTo(s, z)
	case []int32:
		intSliceTo(s, z)
	case []int64:
		intSliceTo(s, z)

	case int:
		itoaS(s, z)
	case int8:
		itoaS(s, z)
	case int16:
		itoaS(s, z)
	case int32:
		itoaS(s, z)
	case int64:
		itoaS(s, z)

	case []uint:
		uintSliceTo(s, z)
	// case []uint8: // = []byte
	case []uint16:
		uintSliceTo(s, z)
	case []uint32:
		uintSliceTo(s, z)
	case []uint64:
		uintSliceTo(s, z)

	case uint:
		utoaS(s, z)
	case uint8:
		utoaS(s, z)
	case uint16:
		utoaS(s, z)
	case uint32:
		utoaS(s, z)
	case uint64:
		utoaS(s, z)

	case []float32:
		floatSliceTo(s, z)
	case []float64:
		floatSliceTo(s, z)

	case float32:
		ftoaS(s, z)
	case float64:
		ftoaS(s, z)

	case []complex64:
		complexSliceTo(s, z)
	case []complex128:
		complexSliceTo(s, z)

	case complex64:
		ctoaS(s, z)
	case complex128:
		ctoaS(s, z)

	default:
		if s.jsonMode {
			if m, ok := val.(interface{ MarshalJSON() ([]byte, error) }); ok {
				data, err := m.MarshalJSON()
				if err != nil {
					hintInternal(err, "MarshalJSON failed")
					break
				}
				s.pcAppendStringValue(string(data))
				break
			}
		} else {
			if m, ok := val.(encoding.TextMarshaler); ok {
				data, err := m.MarshalText()
				if err != nil {
					hintInternal(err, "MarshalText failed")
					break
				}
				s.pcAppendStringValue(string(data))
				break
			}
		}

		// var typ = reflect.TypeOf(z)
		// if kind := typ.Kind(); kind == reflect.Slice {
		// 	if vax, ok := z.(MarshalLogArray); ok {
		// 		if err := vax.MarshalLogArray(pc.SB); err != nil {
		// 			break
		// 		}
		// 	}
		// 	s.expandSlice(val, typ, pc)
		// 	break
		// }

		// TODO remove usage to fmt.Sprintf
		s.pcTryQuoteValue(fmt.Sprintf("{{%v}}", z))
	}
}

// func (s *PrintCtx) expandSlice(val any, valtyp reflect.Type) {
// 	v := reflect.ValueOf(val)
// 	// kind := valtyp.Elem().Kind()
//
// 	s.pcAppendRune('[')
// 	for i := 0; i < v.Len(); i++ {
// 		if i > 0 {
// 			s.pcAppendRune(',')
// 			// s.appendRune(' ')
// 		}
// 		ve := v.Index(i)
// 		switch {
// 		case ve.CanInterface():
// 			vev := ve.Interface()
// 			s.appendValue(vev)
//
// 		// switch kind {
// 		// case reflect.Bool:
// 		// 	s.appendValue(vev, pc)
// 		// case reflect.Int:
// 		// 	s.appendValue(vev, pc)
// 		// case reflect.Int8:
// 		// 	s.appendValue(vev, pc)
// 		// case reflect.Int16:
// 		// 	s.appendValue(vev, pc)
// 		// case reflect.Int32:
// 		// 	s.appendValue(vev, pc)
// 		// case reflect.Int64:
// 		// 	s.appendValue(vev, pc)
// 		// case reflect.Uint:
// 		// 	s.appendValue(vev, pc)
// 		// case reflect.Uint8:
// 		// 	s.appendValue(vev, pc)
// 		// case reflect.Uint16:
// 		// 	s.appendValue(vev, pc)
// 		// case reflect.Uint32:
// 		// 	s.appendValue(vev, pc)
// 		// case reflect.Uint64:
// 		// 	s.appendValue(vev, pc)
// 		// case reflect.Uintptr:
// 		// case reflect.Float32:
// 		// 	s.appendValue(vev, pc)
// 		// case reflect.Float64:
// 		// 	s.appendValue(vev, pc)
// 		// case reflect.Complex64:
// 		// 	s.appendValue(vev, pc)
// 		// case reflect.Complex128:
// 		// 	s.appendValue(vev, pc)
// 		// case reflect.Array:
// 		// case reflect.Chan:
// 		// case reflect.Func:
// 		// case reflect.Interface:
// 		// case reflect.Map:
// 		// case reflect.Pointer:
// 		// case reflect.Slice:
// 		// case reflect.String:
// 		// 	s.appendValue(vev, pc)
// 		// case reflect.Struct:
// 		// case reflect.UnsafePointer:
// 		// default:
// 		// 	s.pcAppendStringValue(s.pcTryQuoteValue(fmt.Sprintf("{{%v}}", z)))
// 		// }
//
// 		case ve.IsZero():
// 			s.pcAppendString("<zero>")
// 		case ve.IsNil():
// 			s.pcAppendString("<nil>")
// 		}
// 	}
// 	s.pcAppendRune(']')
// }

//

//

//

func (s *PrintCtx) checkerr(err error) {
	if err != nil {
		hintInternal(err, "PrintCtx: some formats failed")
	}
}

// appendEscapedJSONString escapes s for JSON and appends it to buf.
// It does not surround the string in quotation marks.
//
// Modified from encoding/json/encode.go:encodeState.string,
// with escapeHTML set to false.
func (s *PrintCtx) appendEscapedJSONString(val string) {
	char := func(b byte) { s.pcAppendByte(b) /*buf = append(buf, b)*/ }
	strz := func(str string) { s.pcAppendString(str) /*buf = append(buf, s...) */ }

	start := 0
	for i := 0; i < len(val); {
		if b := val[i]; b < utf8.RuneSelf {
			if safeSet[b] {
				i++
				continue
			}
			if start < i {
				strz(val[start:i])
			}
			char('\\')
			switch b {
			case '\\', '"':
				char(b)
			case '\n':
				char('n')
			case '\r':
				char('r')
			case '\t':
				char('t')
			default:
				// This encodes bytes < 0x20 except for \t, \n and \r.
				strz(`u00`)
				char(hex[b>>4])
				char(hex[b&0xF])
			}
			i++
			start = i
			continue
		}
		c, size := utf8.DecodeRuneInString(val[i:])
		if c == utf8.RuneError && size == 1 {
			if start < i {
				strz(val[start:i])
			}
			strz(`\ufffd`)
			i += size
			start = i
			continue
		}
		// U+2028 is LINE SEPARATOR.
		// U+2029 is PARAGRAPH SEPARATOR.
		// They are both technically valid characters in JSON strings,
		// but don't work in JSONP, which has to be evaluated as JavaScript,
		// and can lead to security holes there. It is valid JSON to
		// escape them, so we do so unconditionally.
		// See http://timelessrepo.com/json-isnt-a-javascript-subset for discussion.
		if c == '\u2028' || c == '\u2029' {
			if start < i {
				strz(val[start:i])
			}
			strz(`\u202`)
			char(hex[c&0xF])
			i += size
			start = i
			continue
		}
		i += size
	}
	if start < len(val) {
		strz(val[start:])
	}
}

var hex = "0123456789abcdef"

// Copied from encoding/json/tables.go.
//
// safeSet holds the value true if the ASCII character with the given array
// position can be represented inside a JSON string without any further
// escaping.
//
// All values are true except for the ASCII control characters (0-31), the
// double quote ("), and the backslash character ("\").
var safeSet = [utf8.RuneSelf]bool{
	' ':      true,
	'!':      true,
	'"':      false,
	'#':      true,
	'$':      true,
	'%':      true,
	'&':      true,
	'\'':     true,
	'(':      true,
	')':      true,
	'*':      true,
	'+':      true,
	',':      true,
	'-':      true,
	'.':      true,
	'/':      true,
	'0':      true,
	'1':      true,
	'2':      true,
	'3':      true,
	'4':      true,
	'5':      true,
	'6':      true,
	'7':      true,
	'8':      true,
	'9':      true,
	':':      true,
	';':      true,
	'<':      true,
	'=':      true,
	'>':      true,
	'?':      true,
	'@':      true,
	'A':      true,
	'B':      true,
	'C':      true,
	'D':      true,
	'E':      true,
	'F':      true,
	'G':      true,
	'H':      true,
	'I':      true,
	'J':      true,
	'K':      true,
	'L':      true,
	'M':      true,
	'N':      true,
	'O':      true,
	'P':      true,
	'Q':      true,
	'R':      true,
	'S':      true,
	'T':      true,
	'U':      true,
	'V':      true,
	'W':      true,
	'X':      true,
	'Y':      true,
	'Z':      true,
	'[':      true,
	'\\':     false,
	']':      true,
	'^':      true,
	'_':      true,
	'`':      true,
	'a':      true,
	'b':      true,
	'c':      true,
	'd':      true,
	'e':      true,
	'f':      true,
	'g':      true,
	'h':      true,
	'i':      true,
	'j':      true,
	'k':      true,
	'l':      true,
	'm':      true,
	'n':      true,
	'o':      true,
	'p':      true,
	'q':      true,
	'r':      true,
	's':      true,
	't':      true,
	'u':      true,
	'v':      true,
	'w':      true,
	'x':      true,
	'y':      true,
	'z':      true,
	'{':      true,
	'|':      true,
	'}':      true,
	'~':      true,
	'\u007f': true,
}

func (s *PrintCtx) appendBytes(z []byte) {
	_, err := s.Write(z)
	if err != nil {
		hintInternal(err, "PrintCtx: appendBytes failed")
	}
}

func (s *PrintCtx) appendStringSlice(val []string) {
	s.buf = append(s.buf, '[')
	for i := range val {
		if i > 0 {
			s.buf = append(s.buf, ',')
		}
		// s.buf = strconv.AppendQuote(s.buf, val[i])
		s.buf = append(s.buf, '"')
		s.appendEscapedJSONString(val[i])
		s.buf = append(s.buf, '"')
	}
	s.buf = append(s.buf, ']')
}

func (s *PrintCtx) appendBoolSlice(val []bool) {
	s.buf = append(s.buf, '[')
	for i := range val {
		if i > 0 {
			s.buf = append(s.buf, ',')
		}
		s.buf = strconv.AppendBool(s.buf, val[i])
	}
	s.buf = append(s.buf, ']')
}

func intSliceTo[T Integers](s *PrintCtx, val IntSlice[T]) {
	s.buf = append(s.buf, '[')
	for i := range val {
		if i > 0 {
			s.buf = append(s.buf, ',')
		}
		s.buf = strconv.AppendInt(s.buf, int64(val[i]), 10)
	}
	s.buf = append(s.buf, ']')
}

func uintSliceTo[T Uintegers](s *PrintCtx, val UintSlice[T]) {
	s.buf = append(s.buf, '[')
	for i := range val {
		if i > 0 {
			s.buf = append(s.buf, ',')
		}
		s.buf = strconv.AppendUint(s.buf, uint64(val[i]), 10)
	}
	s.buf = append(s.buf, ']')
}

func floatSliceTo[T Floats](s *PrintCtx, val FloatSlice[T]) {
	s.buf = append(s.buf, '[')
	for i := range val {
		if i > 0 {
			s.buf = append(s.buf, ',')
		}
		ftoaS(s, val[i])
	}
	s.buf = append(s.buf, ']')
}

func complexSliceTo[T Complexes](s *PrintCtx, val ComplexSlice[T]) {
	s.buf = append(s.buf, '[')
	for i := range val {
		if i > 0 {
			s.buf = append(s.buf, ',')
		}
		ctoaS(s, val[i])
	}
	s.buf = append(s.buf, ']')
}

func (s *PrintCtx) appendDuration(z time.Duration) {
	s.pcAppendByte('"')
	s.appendEscapedJSONString(z.String())
	s.pcAppendByte('"')
}

func (s *PrintCtx) appendDurationSlice(z []time.Duration) {
	s.pcAppendByte('[')
	for i, dur := range z {
		if i > 0 {
			s.pcAppendByte(',')
		}
		s.appendDuration(dur)
	}
	s.pcAppendByte(']')
}

func (s *PrintCtx) appendTime(z time.Time) {
	const layout = time.RFC3339Nano
	if s.jsonMode || s.noColor {
		s.pcAppendByte('"')
		s.buf = z.AppendFormat(s.buf, layout)
		s.pcAppendByte('"')
	} else {
		s.buf = z.AppendFormat(s.buf, layout)
	}
}

func (s *PrintCtx) appendTimeSlice(z []time.Time) {
	s.pcAppendByte('[')
	for i, tm := range z {
		if i > 0 {
			s.pcAppendByte(',')
		}
		s.appendTime(tm)
	}
	s.pcAppendByte(']')
}

// appendTimestamp is specially for printing logging timestamp
func (s *PrintCtx) appendTimestamp(z time.Time) {
	var tm time.Time

	if s.utcTime == 2 || (s.utcTime == 0 && flags&LlocalTime == 0) {
		tm = z.UTC()
	} else {
		tm = z
	}

	var layout string
	if s.layout != "" {
		layout = s.layout
	} else {
		var ok bool
		if layout, ok = defaultLayouts[flags&Ldatetimeflags]; !ok {
			layout = TimeNano
		}
	}

	// return tm.Format(layout)

	if s.jsonMode || s.noColor {
		s.pcAppendByte('"')
		s.buf = tm.AppendFormat(s.buf, layout)
		// t := tm.Format(layout)
		// s.pcAppendString(t)
		s.pcAppendByte('"')
		// s.pcAppendByte(' ')
	} else {
		s.buf = tm.AppendFormat(s.buf, layout)
		s.pcAppendByte('|')
		// t := tm.Format(layout)
		// s.pcAppendString(t)
	}
}

func itoaS[T Integers](s *PrintCtx, val T) {
	// return intToStringEx(val, 10)

	// s.pcAppendStringValue(intToString(value))

	if s.jsonMode {
		s.checkerr(s.WriteByte('"'))
		itoasimple[T](s, val, 10)
		s.checkerr(s.WriteByte('"'))
	} else {
		itoasimple[T](s, val, 10)
	}
}

func itoasimple[T Integers](s *PrintCtx, val T, base int) {
	s.buf = strconv.AppendInt(s.buf, int64(val), base)
}

func utoaS[T Uintegers](s *PrintCtx, val T) {
	// return intToStringEx(val, 10)

	// s.pcAppendStringValue(intToString(value))

	if s.jsonMode {
		s.checkerr(s.WriteByte('"'))
		utoasimple[T](s, val, 10)
		s.checkerr(s.WriteByte('"'))
	} else {
		utoasimple[T](s, val, 10)
	}
}

func utoasimple[T Uintegers](s *PrintCtx, val T, base int) {
	s.buf = strconv.AppendUint(s.buf, uint64(val), base)
}

func ftoaS[T Floats](s *PrintCtx, val T) {
	// return intToStringEx(val, 10)

	// s.pcAppendStringValue(intToString(value))

	if s.jsonMode {
		s.checkerr(s.WriteByte('"'))
		ftoasimple[T](s, val, 'f', -1, 64)
		s.checkerr(s.WriteByte('"'))
	} else {
		ftoasimple[T](s, val, 'f', -1, 64)
	}
}

func ftoasimple[T Floats](s *PrintCtx, val T, format byte, prec, bitSize int) {
	s.buf = strconv.AppendFloat(s.buf, float64(val), format, prec, bitSize)
}

func ctoaS[T Complexes](s *PrintCtx, val T) {
	// return intToStringEx(val, 10)

	// s.pcAppendStringValue(intToString(value))

	if s.jsonMode {
		s.checkerr(s.WriteByte('"'))
		ctoasimple[T](s, val, 'f', -1, 64)
		s.checkerr(s.WriteByte('"'))
	} else {
		ctoasimple[T](s, val, 'f', -1, 64)
	}
}

func ctoasimple[T Complexes](s *PrintCtx, val T, format byte, prec, bitSize int) {
	s.checkerr(s.WriteByte('('))
	s.buf = strconv.AppendFloat(s.buf, real(complex128(val)), format, prec, bitSize)
	istart := len(s.buf)
	s.checkerr(s.WriteByte('+'))
	ix := len(s.buf)
	s.buf = strconv.AppendFloat(s.buf, imag(complex128(val)), format, prec, bitSize)
	if s.buf[ix] == '+' || s.buf[ix] == '-' {
		end := len(s.buf)
		copy(s.buf[istart:], s.buf[ix:])
		s.buf[end-1] = 'i'
		s.checkerr(s.WriteByte(')'))
	} else {
		s.checkerr(s.WriteByte('i'))
		s.checkerr(s.WriteByte(')'))
	}
}

func btoaS(s *PrintCtx, val bool) {
	s.buf = strconv.AppendBool(s.buf, val)
}
