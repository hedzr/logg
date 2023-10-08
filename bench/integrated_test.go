package bench

import (
	"context"
	"errors"
	"fmt"
	"io"
	"log/slog"
	"slices"
	"strings"
	"testing"
	"time"

	slogg "github.com/hedzr/logg/slog"
)

func BenchmarkDisabledWithoutFields(b *testing.B) {
	b.Logf("Logging at a disabled level without any structured context.")
	elapsedTimes := make(map[string]time.Duration)

	b.Run("hedzr/logg/slog", func(b *testing.B) {
		logger := newDisabledLogg()
		b.ResetTimer()
		b.RunParallel(func(pb *testing.PB) {
			for pb.Next() {
				logger.Info(getMessage(0))
			}
		})
		elapsedTimes[b.Name()] = b.Elapsed()
	})
	b.Run("slog", func(b *testing.B) {
		logger := newDisabledSlog()
		b.ResetTimer()
		b.RunParallel(func(pb *testing.PB) {
			for pb.Next() {
				logger.Info(getMessage(0))
			}
		})
		elapsedTimes[b.Name()] = b.Elapsed()
	})
	b.Run("slog.LogAttrs", func(b *testing.B) {
		logger := newDisabledSlog()
		b.ResetTimer()
		b.RunParallel(func(pb *testing.PB) {
			for pb.Next() {
				logger.LogAttrs(context.Background(), slog.LevelInfo, getMessage(0))
			}
		})
		elapsedTimes[b.Name()] = b.Elapsed()
	})
	dumpElapsedTimes(b, elapsedTimes)
}

func BenchmarkDisabledAccumulatedContext(b *testing.B) {
	b.Logf("Logging at a disabled level with some accumulated context.")
	elapsedTimes := make(map[string]time.Duration)

	b.Run("hedzr/logg/slog", func(b *testing.B) {
		logger := newDisabledLogg(fakeLoggFields()...)
		b.ResetTimer()
		b.RunParallel(func(pb *testing.PB) {
			for pb.Next() {
				logger.Info(getMessage(0))
			}
		})
		elapsedTimes[b.Name()] = b.Elapsed()
	})
	b.Run("slog", func(b *testing.B) {
		logger := newDisabledSlog(fakeSlogFields()...)
		b.ResetTimer()
		b.RunParallel(func(pb *testing.PB) {
			for pb.Next() {
				logger.Info(getMessage(0))
			}
		})
		elapsedTimes[b.Name()] = b.Elapsed()
	})
	b.Run("slog.LogAttrs", func(b *testing.B) {
		logger := newDisabledSlog(fakeSlogFields()...)
		b.ResetTimer()
		b.RunParallel(func(pb *testing.PB) {
			for pb.Next() {
				logger.LogAttrs(context.Background(), slog.LevelInfo, getMessage(0))
			}
		})
		elapsedTimes[b.Name()] = b.Elapsed()
	})
	dumpElapsedTimes(b, elapsedTimes)
}

func BenchmarkDisabledAddingFields(b *testing.B) {
	b.Logf("Logging at a disabled level, adding context at each log site.")
	elapsedTimes := make(map[string]time.Duration)

	b.Run("hedzr/logg/slog", func(b *testing.B) {
		logger := newDisabledLogg()
		b.ResetTimer()
		b.RunParallel(func(pb *testing.PB) {
			for pb.Next() {
				logger.Info(getMessage(0), fakeLoggArgs()...)
			}
		})
		elapsedTimes[b.Name()] = b.Elapsed()
	})
	b.Run("slog", func(b *testing.B) {
		logger := newDisabledSlog()
		b.ResetTimer()
		b.RunParallel(func(pb *testing.PB) {
			for pb.Next() {
				logger.Info(getMessage(0), fakeSlogArgs()...)
			}
		})
		elapsedTimes[b.Name()] = b.Elapsed()
	})
	b.Run("slog.LogAttrs", func(b *testing.B) {
		logger := newDisabledSlog()
		b.ResetTimer()
		b.RunParallel(func(pb *testing.PB) {
			for pb.Next() {
				logger.LogAttrs(context.Background(), slog.LevelInfo, getMessage(0), fakeSlogFields()...)
			}
		})
		elapsedTimes[b.Name()] = b.Elapsed()
	})
	dumpElapsedTimes(b, elapsedTimes)
}

func BenchmarkWithoutFields(b *testing.B) {
	b.Logf("Logging without any structured context. [BenchmarkWithoutFields]")
	elapsedTimes := make(map[string]time.Duration)

	b.Run("hedzr/logg/slog", func(b *testing.B) {
		logger := newLoggTextMode()
		b.ResetTimer()
		b.RunParallel(func(pb *testing.PB) {
			for pb.Next() {
				logger.Info(getMessage(0))
			}
		})
		elapsedTimes[b.Name()] = b.Elapsed()
	})
	b.Run("slog", func(b *testing.B) {
		logger := newSlog()
		b.ResetTimer()
		b.RunParallel(func(pb *testing.PB) {
			for pb.Next() {
				logger.Info(getMessage(0))
			}
		})
		elapsedTimes[b.Name()] = b.Elapsed()
	})
	b.Run("slog.LogAttrs", func(b *testing.B) {
		logger := newSlog()
		b.ResetTimer()
		b.RunParallel(func(pb *testing.PB) {
			for pb.Next() {
				logger.LogAttrs(context.Background(), slog.LevelInfo, getMessage(0))
			}
		})
		elapsedTimes[b.Name()] = b.Elapsed()
	})
	dumpElapsedTimes(b, elapsedTimes)
}

func BenchmarkAccumulatedContext(b *testing.B) {
	b.Logf("Logging with some accumulated context. [BenchmarkAccumulatedContext]")
	elapsedTimes := make(map[string]time.Duration)

	b.Run("hedzr/logg/slog TEXT", func(b *testing.B) {
		logger := newLoggTextMode().With(fakeLoggArgs()...)
		b.ResetTimer()
		b.RunParallel(func(pb *testing.PB) {
			for pb.Next() {
				logger.Info(getMessage(0))
			}
		})
		elapsedTimes[b.Name()] = b.Elapsed()
	})
	b.Run("hedzr/logg/slog COLOR", func(b *testing.B) {
		logger := newLoggColorMode().With(fakeLoggArgs()...)
		b.ResetTimer()
		b.RunParallel(func(pb *testing.PB) {
			for pb.Next() {
				logger.Info(getMessage(0))
			}
		})
		elapsedTimes[b.Name()] = b.Elapsed()
	})
	b.Run("hedzr/logg/slog JSON", func(b *testing.B) {
		logger := newLogg().With(fakeLoggArgs()...)
		b.ResetTimer()
		b.RunParallel(func(pb *testing.PB) {
			for pb.Next() {
				logger.Info(getMessage(0))
			}
		})
		elapsedTimes[b.Name()] = b.Elapsed()
	})
	b.Run("slog", func(b *testing.B) {
		logger := newSlog(fakeSlogFields()...)
		b.ResetTimer()
		b.RunParallel(func(pb *testing.PB) {
			for pb.Next() {
				logger.Info(getMessage(0))
			}
		})
		elapsedTimes[b.Name()] = b.Elapsed()
	})
	b.Run("slog.LogAttrs", func(b *testing.B) {
		logger := newSlog(fakeSlogFields()...)
		b.ResetTimer()
		b.RunParallel(func(pb *testing.PB) {
			for pb.Next() {
				logger.LogAttrs(context.Background(), slog.LevelInfo, getMessage(0))
			}
		})
		elapsedTimes[b.Name()] = b.Elapsed()
	})
	dumpElapsedTimes(b, elapsedTimes)
}

func BenchmarkAccumulatedContextLoggOnly(b *testing.B) {
	b.Logf("Logging with some accumulated context. [BenchmarkAccumulatedContextLoggOnly]")
	elapsedTimes := make(map[string]time.Duration)
	loggerWarmUp := newLoggTextMode()
	_ = loggerWarmUp
	loggerWarmUp = newLoggTextMode()
	_ = loggerWarmUp
	b.Run("hedzr/logg/slog", func(b *testing.B) {
		logger := newLoggTextMode().With(fakeLoggArgs()...)
		b.ResetTimer()
		b.RunParallel(func(pb *testing.PB) {
			for pb.Next() {
				logger.Info(getMessage(1))
			}
		})
		elapsedTimes[b.Name()] = b.Elapsed()
	})
	dumpElapsedTimes(b, elapsedTimes)
}

func BenchmarkAddingFields(b *testing.B) {
	b.Logf("Logging with additional context at each log site. [BenchmarkAddingFields]")
	elapsedTimes := make(map[string]time.Duration)

	b.Run("hedzr/logg/slog TEXT", func(b *testing.B) {
		logger := newLoggTextMode()
		b.ResetTimer()
		b.RunParallel(func(pb *testing.PB) {
			for pb.Next() {
				logger.Info(getMessage(0), fakeLoggArgs()...)
			}
		})
		elapsedTimes[b.Name()] = b.Elapsed()
	})
	b.Run("hedzr/logg/slog COLOR", func(b *testing.B) {
		logger := newLoggColorMode()
		b.ResetTimer()
		b.RunParallel(func(pb *testing.PB) {
			for pb.Next() {
				logger.Info(getMessage(0), fakeLoggArgs()...)
			}
		})
		elapsedTimes[b.Name()] = b.Elapsed()
	})
	b.Run("hedzr/logg/slog JSON", func(b *testing.B) {
		logger := newLogg()
		b.ResetTimer()
		b.RunParallel(func(pb *testing.PB) {
			for pb.Next() {
				logger.Info(getMessage(0), fakeLoggArgs()...)
			}
		})
		elapsedTimes[b.Name()] = b.Elapsed()
	})
	b.Run("slog", func(b *testing.B) {
		logger := newSlog()
		b.ResetTimer()
		b.RunParallel(func(pb *testing.PB) {
			for pb.Next() {
				logger.Info(getMessage(0), fakeSlogArgs()...)
			}
		})
		elapsedTimes[b.Name()] = b.Elapsed()
	})
	b.Run("slog.LogAttrs", func(b *testing.B) {
		logger := newSlog()
		b.ResetTimer()
		b.RunParallel(func(pb *testing.PB) {
			for pb.Next() {
				logger.LogAttrs(context.Background(), slog.LevelInfo, getMessage(0), fakeSlogFields()...)
			}
		})
		elapsedTimes[b.Name()] = b.Elapsed()
	})
	// b.Logf("%v", elapsedTimes)
	dumpElapsedTimes(b, elapsedTimes)
}

func dumpElapsedTimes(b testing.TB, m map[string]time.Duration) {
	// b.Logf("%v", m)

	keys := make([]string, 0, len(m))
	for k := range m {
		keys = append(keys, k)
	}

	// sort.Strings(keys)
	slices.Sort(keys)

	lastOne := float64(m[keys[len(keys)-1]].Nanoseconds())
	for _, k := range keys {
		v := m[k]
		b.Logf("%-56s\t%16s\t%v", rightPad(k, 56), v.String(), float64(v.Nanoseconds())/lastOne)
	}
}

func rightPad(s string, pad int) string {
	if len(s) < pad {
		str := s + strings.Repeat(" ", pad)
		return str[:pad]
	}
	return s
}

func BenchmarkAddingFieldsOnly(b *testing.B) {
	b.Logf("Logging with additional context at each log site. [BenchmarkAddingFields]")
	b.Run("hedzr/logg/slog", func(b *testing.B) {
		logger := newLogg()
		b.ResetTimer()
		b.RunParallel(func(pb *testing.PB) {
			for pb.Next() {
				logger.Info(getMessage(1), fakeLoggArgs()...)
			}
		})
	})
}

func init() {
	// no source:
	slogg.RemoveFlags(slogg.Lcaller | slogg.LattrsR)
}

func newLogg(fields ...slogg.Attr) slogg.Logger {
	return slogg.New().WithLevel(slogg.DebugLevel).WithAttrs(fields...).
		WithJSONMode().
		// WithLevel(slogg.OffLevel)
		WithWriter(io.Discard)
	// .EnableJSON(useJSON).EnableSearchKnownPackage(searchingKnownPackages).Build()
	// return hilog.New(hilog.NewJSONHandler(io.Discard, nil).WithAttrs(fields))
}

func newLoggColorMode(fields ...slogg.Attr) slogg.Logger {
	return slogg.New().WithLevel(slogg.DebugLevel).WithAttrs(fields...).
		WithColorMode().
		// WithLevel(slogg.OffLevel)
		WithWriter(io.Discard)
}

func newLoggTextMode(fields ...slogg.Attr) slogg.Logger {
	return slogg.New().WithLevel(slogg.DebugLevel).WithAttrs(fields...).
		WithColorMode(false).
		// WithLevel(slogg.OffLevel)
		WithWriter(io.Discard)
}

func newDisabledLogg(fields ...slogg.Attr) slogg.Logger {
	slogg.SetLevel(slogg.WarnLevel)
	return slogg.New().WithLevel(slogg.WarnLevel).WithAttrs(fields...).WithColorMode(false).WithWriter(io.Discard)
	// .EnableJSON(useJSON).EnableSearchKnownPackage(searchingKnownPackages).Build()
	// return hilog.New(hilog.NewJSONHandler(io.Discard, &slog.HandlerOptions{Level: slog.LevelError}).WithAttrs(fields))
}

func fakeLoggFields() []slogg.Attr { return loggAttrs }

var loggAttrs = []slogg.Attr{
	slogg.NewAttr("int", _tenInts[0]),
	slogg.NewAttr("ints", _tenInts),
	slogg.NewAttr("string", _tenStrings[0]),
	slogg.NewAttr("strings", _tenStrings),
	slogg.NewAttr("time", _tenTimes[0]),
	slogg.NewAttr("times", _tenTimes),
	slogg.NewAttr("user1", _oneUser),
	slogg.NewAttr("user2", _oneUser),
	slogg.NewAttr("users", _tenUsers),
	slogg.NewAttr("error", errExample),
}

func fakeLoggArgs() []any { return loggArgs }

var loggArgs = []any{
	"int", _tenInts[0],
	"ints", _tenInts,
	"string", _tenStrings[0],
	"strings", _tenStrings,
	"time", _tenTimes[0],
	"times", _tenTimes,
	"user1", _oneUser,
	"user2", _oneUser,
	"users", _tenUsers,
	"error", errExample,
}

func newSlog(fields ...slog.Attr) *slog.Logger {
	return slog.New(slog.NewJSONHandler(io.Discard, nil).WithAttrs(fields))
}

func newDisabledSlog(fields ...slog.Attr) *slog.Logger {
	return slog.New(slog.NewJSONHandler(io.Discard, &slog.HandlerOptions{Level: slog.LevelError}).WithAttrs(fields))
}

func fakeSlogFields() []slog.Attr { return slogAttrs }

func fakeSlogArgs() []any { return slogArgs }

var slogAttrs = []slog.Attr{
	slog.Int("int", _tenInts[0]),
	slog.Any("ints", _tenInts),
	slog.String("string", _tenStrings[0]),
	slog.Any("strings", _tenStrings),
	slog.Time("time", _tenTimes[0]),
	slog.Any("times", _tenTimes),
	slog.Any("user1", _oneUser),
	slog.Any("user2", _oneUser),
	slog.Any("users", _tenUsers),
	slog.Any("error", errExample),
}

var slogArgs = []any{
	"int", _tenInts[0],
	"ints", _tenInts,
	"string", _tenStrings[0],
	"strings", _tenStrings,
	"time", _tenTimes[0],
	"times", _tenTimes,
	"user1", _oneUser,
	"user2", _oneUser,
	"users", _tenUsers,
	"error", errExample,
}

var (
	errExample = errors.New("fail")

	_messages   = fakeMessages(1000)
	_tenInts    = []int{1, 2, 3, 4, 5, 6, 7, 8, 9, 0}
	_tenStrings = []string{"a", "b", "c", "d", "e", "f", "g", "h", "i", "j"}
	_tenTimes   = []time.Time{
		time.Unix(0, 0),
		time.Unix(1, 0),
		time.Unix(2, 0),
		time.Unix(3, 0),
		time.Unix(4, 0),
		time.Unix(5, 0),
		time.Unix(6, 0),
		time.Unix(7, 0),
		time.Unix(8, 0),
		time.Unix(9, 0),
	}
	_oneUser = &user{
		Name:      "Jane Doe",
		Email:     "jane@test.com",
		CreatedAt: time.Date(1980, 1, 1, 12, 0, 0, 0, time.UTC),
	}
	_tenUsers = users{
		_oneUser,
		_oneUser,
		_oneUser,
		_oneUser,
		_oneUser,
		_oneUser,
		_oneUser,
		_oneUser,
		_oneUser,
		_oneUser,
	}
)

func fakeMessages(n int) []string {
	messages := make([]string, n)
	for i := range messages {
		messages[i] = fmt.Sprintf("Test logging, but use a somewhat realistic message length. (#%v)", i)
	}
	return messages
}

func getMessage(iter int) string {
	return _messages[iter%1000]
}

type users []*user

func (uu users) MarshalSlogArray(enc *slogg.PrintCtx) error {
	var err error
	for i := range uu {
		if e := uu[i].MarshalSlogObject(enc); e != nil {
			err = errors.Join(err, e)
		}
	}
	return err
}

type user struct {
	Name      string    `json:"name"`
	Email     string    `json:"email"`
	CreatedAt time.Time `json:"created_at"`
}

func (u *user) MarshalSlogObject(enc *slogg.PrintCtx) error {
	enc.AddString("name", u.Name)
	enc.AddString("email", u.Email)
	enc.AddInt64("createdAt", u.CreatedAt.UnixNano())
	return nil
}
