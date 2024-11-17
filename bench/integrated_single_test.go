package bench

import (
	"context"
	"testing"
	"time"
)

// need optimizing, over 1 allocs (logfmt, json), or 53 allocas (color)
//
//	go test -cpuprofile cpu.prof -memprofile mem.prof -benchmem -run=^$ -tags hzstudio,hzwork -bench ^BenchmarkWithoutFieldsSpecial$ github.com/hedzr/logg/bench -v
//	go test -cpuprofile cpu.prof -benchmem -run=^$ -tags hzstudio,hzwork -bench ^BenchmarkWithoutFieldsSpecial$ github.com/hedzr/logg/bench -v
//	go test -memprofile mem.prof -benchmem -run=^$ -tags hzstudio,hzwork -bench ^BenchmarkWithoutFieldsSpecial$ github.com/hedzr/logg/bench -v
//	go tool pprof -http=:6060 cpu.prof
//
// run all benchmarks
//
//	go test -benchmem -bench . -run=^$ github.com/hedzr/logg/bench -v
//	go test -benchmem -bench . -run=^$ github.com/hedzr/go-zag/benchmarks -v
func BenchmarkWithoutFieldsSpecial(b *testing.B) {
	b.Logf("Logging without any structured context. [BenchmarkWithoutFields]")
	elapsedTimes := make(map[string]time.Duration)
	ctx := context.Background()

	b.Run("hedzr/logg/slog", func(b *testing.B) {
		logger := newLoggTextMode()
		b.ResetTimer()
		b.RunParallel(func(pb *testing.PB) {
			for pb.Next() {
				logger.InfoContext(ctx, getMessage(0))
			}
		})
		elapsedTimes[b.Name()] = b.Elapsed()
	})
}

// need optimizing, over 2 allocs (logfmt, json), or 144 allocas (color)
func BenchmarkAccumulatedContextSpecial(b *testing.B) {
	b.Logf("Logging with some accumulated context. [BenchmarkAccumulatedContext]")
	elapsedTimes := make(map[string]time.Duration)
	ctx := context.Background()

	b.Run("hedzr/logg/slog TEXT", func(b *testing.B) {
		logger := newLoggTextMode().Set(fakeLoggArgs()...)
		b.ResetTimer()
		b.RunParallel(func(pb *testing.PB) {
			for pb.Next() {
				logger.InfoContext(ctx, getMessage(0))
			}
		})
		elapsedTimes[b.Name()] = b.Elapsed()
	})
}

// need optimizing, over 2 allocs (logfmt, json), or 144 allocas (color)
func BenchmarkAccumulatedContextSpecial1(b *testing.B) {
	b.Logf("Logging with some accumulated context. [BenchmarkAccumulatedContextLoggOnly]")
	elapsedTimes := make(map[string]time.Duration)
	ctx := context.Background()

	loggerWarmUp := newLoggTextMode()
	_ = loggerWarmUp
	loggerWarmUp = newLoggTextMode()
	_ = loggerWarmUp

	b.Run("hedzr/logg/slog", func(b *testing.B) {
		logger := newLoggTextMode().Set(fakeLoggArgs()...)
		b.ResetTimer()
		b.RunParallel(func(pb *testing.PB) {
			for pb.Next() {
				logger.InfoContext(ctx, getMessage(1))
			}
		})
		elapsedTimes[b.Name()] = b.Elapsed()
	})
	dumpElapsedTimes(b, elapsedTimes)
}

// need optimizing, over 12 allocs (logfmt, json), or 154 allocas (color)
func BenchmarkAddingFieldsSpecial(b *testing.B) {
	b.Logf("Logging with additional context at each log site. [BenchmarkAddingFields]")
	ctx := context.Background()
	b.Run("hedzr/logg/slog", func(b *testing.B) {
		logger := newLogg()
		b.ResetTimer()
		b.RunParallel(func(pb *testing.PB) {
			for pb.Next() {
				logger.InfoContext(ctx, getMessage(1), fakeLoggArgs()...)
			}
		})
	})
}
