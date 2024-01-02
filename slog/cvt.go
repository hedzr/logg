package slog

import (
	"strconv"
	"time"
)

// type objwr struct{ SB }

const (
	TimeNoNano      = "15:04:05Z07:00"                      // text-logging timestamp format: time only, without nano second part
	TimeNano        = "15:04:05.000000Z07:00"               // text-logging timestamp format: time only
	DateTime        = "2006-01-0215:04:05Z07:00"            // text-logging timestamp format: date and time, with timezone
	RFC3339Nano     = "2006-01-02T15:04:05.000000Z07:00"    // text-logging timestamp format: RFC3339Nano
	RFC3339NanoOrig = "2006-01-02T15:04:05.999999999Z07:00" // text-logging timestamp format: RFC3339Nano with 9 bits nano seconds
)

var defaultLayouts = map[Flags]string{
	Ldate:                         "2006-01-02",
	Ltime:                         TimeNoNano,
	Ltime | Lmicroseconds:         TimeNano,
	Ldate | Ltime:                 DateTime,
	Ldate | Lmicroseconds:         RFC3339Nano,
	Ldate | Ltime | Lmicroseconds: RFC3339Nano,
}

func intToString[T Integers](val T) string {
	return intToStringEx(val, 10)
}

func intToStringEx[T Integers](val T, base int) string {
	return strconv.FormatInt(int64(val), base)
}

func uintToString[T Uintegers](val T) string {
	return uintToStringEx(val, 10)
}

func uintToStringEx[T Uintegers](val T, base int) string {
	return strconv.FormatUint(uint64(val), base)
}

func floatToString[T Floats](val T) string {
	return floatToStringEx(float64(val), 'f', -1, 64)
}

func floatToStringEx[T Floats](val T, format byte, prec, bitSize int) string {
	return strconv.FormatFloat(float64(val), format, prec, bitSize)
}

func complexToString[T Complexes](val T) string {
	return complexToStringEx(val, 'f', -1, 128)
}

func complexToStringEx[T Complexes](val T, format byte, prec, bitSize int) string {
	return strconv.FormatComplex(complex128(val), format, prec, bitSize)
}

func boolToString(b bool) string {
	return strconv.FormatBool(b)
}

//

func intSliceToString[T Integers](val IntSlice[T]) string {
	var b = make([]byte, 0, len(val)*8) // 8: assume integer need 8 runes
	b = append(b, []byte("[")...)
	for i := range val {
		if i > 0 {
			b = append(b, []byte(",")...)
		}
		b = strconv.AppendInt(b, int64(val[i]), 10)
	}
	b = append(b, []byte("]")...)
	return string(b)
}

func uintSliceToString[T Uintegers](val UintSlice[T]) string {
	var b = make([]byte, 0, len(val)*8) // 8: assume unsigned integer need 8 runes
	b = append(b, []byte("[")...)
	for i := range val {
		if i > 0 {
			b = append(b, []byte(",")...)
		}
		b = strconv.AppendUint(b, uint64(val[i]), 10)
	}
	b = append(b, []byte("]")...)
	return string(b)
}

func floatSliceToString[T Floats](val FloatSlice[T]) string {
	var b = make([]byte, 0, len(val)*16+2) // 8: assume floats need 16 runes
	b = append(b, []byte("[")...)
	for i := range val {
		if i > 0 {
			b = append(b, []byte(",")...)
		}
		b = strconv.AppendFloat(b, float64(val[i]), 'f', -1, 64)
	}
	b = append(b, []byte("]")...)
	return string(b)
}

func complexSliceToString[T Complexes](val ComplexSlice[T]) string {
	var b = make([]byte, 0, len(val)*32+2) // 8: assume complex need 32 runes
	b = append(b, []byte("[")...)
	for i := range val {
		if i > 0 {
			b = append(b, []byte(",")...)
		}
		num := strconv.FormatComplex(complex128(val[i]), 'f', -1, 128)
		b = append(b, []byte(num)...)
	}
	b = append(b, []byte("]")...)
	return string(b)
}

// func complexSliceToString[T Complexes](val ComplexSlice[T]) string {
// 	var b = make([]byte, 0, len(val)*32+2) // 8: assume complex need 32 runes
// 	b = append(b, []byte("[")...)
// 	for i := range val {
// 		if i > 0 {
// 			b = append(b, []byte(",")...)
// 		}
// 		b=strconv.AppendComplex(b, complex128(val[i]), 'f', -1, 128)
// 	}
// 	b = append(b, []byte("]")...)
// 	return string(b)
// }

func stringSliceToString(val []string) string {
	var b = make([]byte, 0, len(val)*32+2) // 8: assume integer need 32 runes
	b = append(b, []byte("[")...)
	for i := range val {
		if i > 0 {
			b = append(b, []byte(",")...)
		}
		b = strconv.AppendQuote(b, val[i])
	}
	b = append(b, []byte("]")...)
	return string(b)
}

func boolSliceToString(val []bool) string {
	var b = make([]byte, 0, len(val)*8) // 8: assume bool need 5 runes
	b = append(b, []byte("[")...)
	for i := range val {
		if i > 0 {
			b = append(b, []byte(",")...)
		}
		b = strconv.AppendBool(b, val[i])
	}
	b = append(b, []byte("]")...)
	return string(b)
}

func timeSliceToString(val []time.Time) string {
	var b = make([]byte, 0, len(val)*32+2) // 8: assume time need 24 runes
	b = append(b, []byte("[")...)
	for i := range val {
		if i > 0 {
			b = append(b, []byte(",")...)
		}
		b = strconv.AppendQuote(b, val[i].Format(time.RFC3339Nano))
	}
	b = append(b, []byte("]")...)
	return string(b)
}

func durationSliceToString(val []time.Duration) string {
	var b = make([]byte, 0, len(val)*16+2) // 8: assume duration need 16 runes
	b = append(b, []byte("[")...)
	for i := range val {
		if i > 0 {
			b = append(b, []byte(",")...)
		}
		b = strconv.AppendQuote(b, val[i].String())
	}
	b = append(b, []byte("]")...)
	return string(b)
}

//

// Stringer is a synonym to fmt.Stringer
type Stringer interface {
	String() string
}

// ToString interface for some object
type ToString interface {
	ToString(args ...any) string
}

// Integers declares signed integers generic type
type Integers interface {
	int | int8 | int16 | int32 | int64
}

// Uintegers declares unsigned integers generic type
type Uintegers interface {
	uint | uint8 | uint16 | uint32 | uint64
}

// Floats declares float number generic type
type Floats interface {
	float32 | float64
}

// Complexes declares complex number generic type
type Complexes interface {
	complex64 | complex128
}

// Numerics declares numeric generic type
type Numerics interface {
	Integers | Uintegers | Floats | Complexes
}

type IntSlice[T Integers] []T      // IntSlice declares slice of signed integers generic type
type UintSlice[T Uintegers] []T    // UintSlice declares slice of unsigned integers generic type
type FloatSlice[T Floats] []T      // FloatSlice declares slice of float number generic type
type ComplexSlice[T Complexes] []T // ComplexSlice declares slice of complex number generic type
type StringSlice[T string] []T     // StringSlice declares slice of string generic type
type BoolSlice[T bool] []T         // BoolSlice declares slice of boolean generic type

type Slice[T Integers | Uintegers | Floats] []T // Slice declares slice of numeric generic type
