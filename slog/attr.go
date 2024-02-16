package slog

import (
	"fmt"
	"time"

	"github.com/hedzr/logg/slog/internal/strings"
)

func NewAttr(key string, val any) Attr              { return &kvp{key, val} }                            // create an attribute
func NewAttrs(args ...any) Attrs                    { return buildUniqueAttrs(nil, args...) }            //nolint:revive,lll // freeform args here, used by WithAttrs1. See also New or With for the usage.
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
	serializeAttrs(pc, s.items)
}

func (s Attrs) SerializeValueTo(pc *PrintCtx) {
	serializeAttrs(pc, s)
}

func serializeAttrs(pc *PrintCtx, kvps Attrs) { //nolint:revive
	prefix := pc.prefix
	for _, v := range kvps {
		if pc.noColor {
			pc.pcAppendComma()
		} else {
			pc.pcAppendByte(' ')
			ct.echoColorAndBg(pc, pc.clr, pc.bg)
		}

		inGroupedMode := pc.inGroupedMode
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
		}
		pc.prefix = prefix
	}

	if !pc.noColor {
		ct.echoResetColor(pc)
	}
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
