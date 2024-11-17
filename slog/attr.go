package slog

import (
	"fmt"
	"slices"
	"time"

	"github.com/hedzr/logg/slog/internal/strings"
)

func NewAttr(key string, val any) Attr              { return &kvp{key, val} }                            // create an attribute
func NewAttrs(args ...any) Attrs                    { return buildUniqueAttrs(args...) }                 //nolint:revive,lll // freeform args here, used by WithAttrs1. See also New or With for the usage.
func NewGroupedAttr(key string, as ...Attr) Attr    { return &gkvp{key, as} }                            // similar with Group
func NewGroupedAttrEasy(key string, as ...any) Attr { return &gkvp{key: key, items: buildAttrs(as...)} } // synonym to Group

// kvp is a tiny key-value pair
type kvp struct {
	key string
	val any
}

func (s *kvp) Key() string    { return s.key }
func (s *kvp) Value() any     { return s.val }
func (s *kvp) SetValue(v any) { s.val = v }

func (s *kvp) SerializeValueTo(pc *PrintCtx) {
	pc.pcAppendStringKey(s.key)
	pc.pcAppendByte('=')

	pc.prefix = ""
	pc.inGroupedMode = false
	pc.appendValue(s.val)
}

type Attrs []Attr // slice of Attr

// gkvp is a grouped kvp
type gkvp struct {
	key   string
	items Attrs
}

type groupedValue interface {
	Add(as ...Attr)
}

func (s *gkvp) Key() string { return s.key }
func (s *gkvp) Value() any  { return s.items }
func (s *gkvp) SetValue(v any) {
	if v1, ok := v.(Attrs); ok {
		s.items = v1
	} else if v1, ok = v.([]Attr); ok {
		s.items = v1
	} else if v2, ok := v.(Attr); ok {
		s.items = append(s.items, v2)
	} else {
		panic(fmt.Sprintf("unexpected value was set: %+v", v))
	}
}

func (s *gkvp) Add(as ...Attr) {
	s.items = append(s.items, as...)
}

// func (s *gkvp) LogValue() Attr {
// 	return s
// }

func (s *gkvp) SerializeValueTo(pc *PrintCtx) {
	if pc.jsonMode {
		if pc.noColor {
			pc.pcAppendStringKey(s.key)
			pc.pcAppendByte(':')
		} else {
			ct.wrapDimColorTo(pc, s.key)
			pc.pcAppendByte(':')
			ct.echoColorAndBg(pc, pc.clr, pc.bg)
		}
	}
	// if sb.jsonMode {
	// 	sb.appendRune('{')
	// }
	// for ix, attr := range s.items {
	// 	if ix > 0 {
	// 		sb.appendRune(',')
	// 	}
	// 	sb.appendStringKey(attr.Key())
	// 	sb.appendRune('=')
	// 	sb.appendValue(attr.Value())
	// }
	// if sb.jsonMode {
	// 	sb.appendRune('}')
	// }
	_ = serializeAttrs(pc, s.items)
}

func (s Attrs) SerializeValueTo(pc *PrintCtx) {
	_ = serializeAttrs(pc, s)
}

func dedupeSlice[S ~[]E, E any](x S, cmp func(a, b E) bool) S {
	if len(x) == 0 {
		return nil
	}

	j := 0
	for i := 1; i < len(x); i++ {
		if cmp(x[i], x[j]) {
			x[j] = x[i]
			continue
		}
		j++
		// preserve the original data
		// in[i], in[j] = in[j], in[i]
		// only set what is required
		x[j] = x[i]
	}
	result := x[:j+1]
	return result
}

// serializeAttrs returns an error object if it's found in the given Attrs.
// The caller can do something with the object, For instance, printImpl
// will dump the error's stack trace if necessary.
func serializeAttrs(pc *PrintCtx, kvps Attrs) (err error) { //nolint:revive
	prefix := pc.prefix
	inGroupedMode := pc.inGroupedMode

	if pc.dedupeAttrs {
		slices.SortFunc(kvps, func(a, b Attr) int {
			if a == nil {
				if b == nil {
					return 0
				}
				return -1
			}
			if b == nil {
				return 1
			}

			k1, k2 := a.Key(), b.Key()
			if k1 < k2 {
				return -1
			}
			if k1 == k2 {
				return 0
			}
			return 1
		})

		// sort.Slice(kvps, func(i, j int) bool {
		// 	return kvps[i].Key() < kvps[j].Key()
		// })
		kvps = dedupeSlice(kvps, func(a, b Attr) bool {
			if a == nil || b == nil {
				if b == a {
					return true
				}
				return false
			}
			return a.Key() == b.Key()
		})
	}

	for _, v := range kvps {
		if v == nil {
			continue
		}

		if pc.noColor {
			pc.pcAppendComma()
		} else {
			pc.pcAppendByte(' ')
			ct.echoColorAndBg(pc, pc.clr, pc.bg)
		}

		if !inGroupedMode {
			_, inGroupedMode = v.(groupedValue)
		}

		key := v.Key()
		if inGroupedMode && !pc.jsonMode && pc.valueStringer == nil {
			key = strings.DotPrefix(key, prefix)
		} else {
			if inGroupedMode && !pc.jsonMode && pc.valueStringer == nil {
				panic("impossible condition matched: inGroupedMode && !pc.jsonMode")
				// if inGroupedMode && !pc.jsonMode {
				// 	key = DotPrefix(key, prefix)
				// }
			}
			if !pc.jsonMode {
				key = strings.DotPrefix(key, prefix)
			}
			if pc.noColor {
				pc.pcAppendStringKey(key)
			} else {
				ct.echoColorAndBg(pc, clrAttrKey, clrAttrKeyBg)
				pc.pcAppendStringKey(key)
				ct.echoColorAndBg(pc, pc.clr, pc.bg)
			}

			pc.pcAppendColon()
		}

		if key == timestampFieldName {
			// we format timestamp in according to the setting in flags
			if z, ok := v.Value().(time.Time); ok {
				// if pc.jsonMode || pc.noColor {
				// 	pc.WriteRune('"')
				// 	pc.WriteString(z.Format(time.RFC3339Nano))
				// 	pc.WriteRune('"')
				// } else {
				// 	pc.appendTimestamp(z)
				// }
				pc.appendTimestamp(z)
				continue
			}
		}

		pc.prefix = key
		val := v.Value()
		if pc.valueStringer != nil { // && IsAnyBitsSet(Lprettyprint) {
			pc.valueStringer.WriteValue(val)
		} else {
			pc.appendValue(val)
			if e, ok := val.(error); ok && e != nil {
				err = e // just a error value
			}
		}
		pc.prefix = prefix
	}

	if !pc.noColor {
		ct.echoResetColor(pc)
	}
	return
}

// type canSerializeValue interface {
// 	SerializeValueTo(pc *printContextS)
// }
//
// type LogObjectMarshaller interface {
// 	MarshalSlogObject(enc *SB) error
// }
//
// type LogArrayMarshaller interface {
// 	MarshalSlogArray(enc *SB) error
// }

// ObjectSerializer to allow your object serialized by our PrintCtx
type ObjectSerializer interface {
	SerializeValueTo(pc *PrintCtx)
}

// ObjectMarshaller to allow your object serialized by our PrintCtx
type ObjectMarshaller interface {
	MarshalSlogObject(enc *PrintCtx) error
}

// ArrayMarshaller to allow your slice or array object serialized by our PrintCtx
type ArrayMarshaller interface {
	MarshalSlogArray(enc *PrintCtx) error
}
