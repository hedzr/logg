package slog

import (
	"context"
	logslog "log/slog"
	"runtime"
)

// NewSlogHandler makes a log/slog Handler to adapt into standard slog.
//
// Demo code:
//
//		import logz "github.com/hedzr/logg/slog"
//	 import "log/slog"
//
//		sll := logz.NewSlogHandler(logz.Default(), &logz.HandlerOptions{
//			NoColor:  false,
//			NoSource: false,
//			JSON:     false,
//			Level:    logz.InfoLevel,
//		})
//
//		logger = slog.New(sll)
//		logger.Info("log/slog: app starting...", "args", os.Args)
//		defer func() {
//			logger.Info("log/slog: app terminated.", "args", os.Args)
//		}()
func NewSlogHandler(logger Logger, config *HandlerOptions) logslog.Handler {
	if config == nil {
		config = &HandlerOptions{} //nolint:revive
	}
	if config.NoSource {
		RemoveFlags(Lcaller)
	} else {
		AddFlags(Lcaller)
	}

	if config.Level != PanicLevel {
		logger.SetLevel(config.Level)
	}

	return &handler4LogSlog{logger.SetColorMode(!config.NoColor).SetJSONMode(config.JSON)}
}

// HandlerOptions is used for our log/slog Handler.
//
// The HandlerOptions.Level field can be used to initial log/slog.Level.
// Due the slog.Level is less than Level, so only InfoLevel, WarnLevel,
// DebugLevel and ErrorLevel can be applied, the others have no effect.
type HandlerOptions struct {
	NoColor  bool  // is colorful outputting?
	NoSource bool  // has caller info?
	JSON     bool  // logging as JSON format?
	Level    Level // zero value means no setup level. Note that zero value represents indeed PanicLevel, so it cannot be used for SetLevel.
}

type handler4LogSlog struct {
	Logger
}

func convertLevelToLogSlog(lvl Level) logslog.Level {
	if l, ok := mLevelToLogSlog[lvl]; ok {
		return l
	}
	return logslog.LevelInfo
}

func convertLogSlogLevel(lvl logslog.Level) Level {
	if l, ok := mLogSlogLevelToLevel[lvl]; ok {
		return l
	}
	return AlwaysLevel
}

func convertLogSlogRecordAttrs(rec logslog.Record) Attrs {
	fields := make([]Attr, 0, rec.NumAttrs())
	rec.Attrs(func(attr logslog.Attr) bool {
		fields = append(fields, convertAttrToField(attr))
		return true
	})
	return fields
}

// Enabled reports whether the handler handles records at the given level.
func (s *handler4LogSlog) Enabled(ctx context.Context, lvl logslog.Level) bool {
	if l, ok := mLogSlogLevelToLevel[lvl]; ok {
		return s.Logger.EnabledContext(ctx, l)
	}
	return true
}

// Handle handles the Record.
func (s *handler4LogSlog) Handle(ctx context.Context, rec logslog.Record) error {
	lvl := convertLogSlogLevel(rec.Level)
	if wi, ok := s.Logger.(LogSlogAware); ok {
		fields := convertLogSlogRecordAttrs(rec)

		// rec.PC would be abandoned because we want skip the extra frames
		ei := 0
		if sa, ok := s.Logger.(interface{ Skip() int }); ok {
			ei = sa.Skip()
		}
		var pcs [1]uintptr
		runtime.Callers(3+1+ei, pcs[:])
		rec.PC = pcs[0]

		wi.WriteThru(ctx, lvl, rec.Time, rec.PC, rec.Message, fields)
	} else {
		fields := convertLogSlogRecordAttrs(rec)
		s.LogAttrs(ctx, lvl, rec.Message, fields)
	}
	return nil
}

// WithAttrs returns a new Handler whose attributes consist of
// both the receiver's attributes and the arguments.
func (s *handler4LogSlog) WithAttrs(attrs []logslog.Attr) logslog.Handler {
	fields := make([]Attr, len(attrs))
	for i, attr := range attrs {
		fields[i] = convertAttrToField(attr)
	}
	return s.withFields(fields...)
}

// WithGroup returns a new Handler with the given group appended to
// the receiver's existing groups.
func (s *handler4LogSlog) WithGroup(name string) logslog.Handler {
	return s.withFields(Group(name))
}

// withFields returns a cloned Handler with the given fields.
func (s *handler4LogSlog) withFields(fields ...Attr) *handler4LogSlog {
	cloned := &handler4LogSlog{
		New().SetAttrs(fields...),
	}
	return cloned
}

var _ logslog.Handler = (*handler4LogSlog)(nil)

func convertGroupToFields(attrs []logslog.Attr) (ret Attrs) {
	for _, a := range attrs {
		ret = append(ret, convertAttrToField(a))
	}
	return
}

func convertAttrToField(attr logslog.Attr) Attr { //nolint:revive
	switch attr.Value.Kind() {
	case logslog.KindBool:
		return Bool(attr.Key, attr.Value.Bool())
	case logslog.KindTime:
		return Time(attr.Key, attr.Value.Time())
	case logslog.KindDuration:
		return Duration(attr.Key, attr.Value.Duration())
	case logslog.KindFloat64:
		return Float64(attr.Key, attr.Value.Float64())
	case logslog.KindInt64:
		return Int64(attr.Key, attr.Value.Int64())
	case logslog.KindString:
		return String(attr.Key, attr.Value.String())
	case logslog.KindUint64:
		return Uint64(attr.Key, attr.Value.Uint64())
	case logslog.KindGroup:
		return Group(attr.Key, convertGroupToFields(attr.Value.Group()))
	case logslog.KindLogValuer:
		return convertAttrToField(logslog.Attr{
			Key: attr.Key,
			// TODO: resolve the value in a lazy way.
			// This probably needs a new Zap field type
			// that can be resolved lazily.
			Value: attr.Value.Resolve(),
		})
	default:
		return Any(attr.Key, attr.Value.Any())
	}
}
