//go:build verbose
// +build verbose

package slog

import (
	"context"
)

// VerboseContext implements Logger.
func (s *Entry) VerboseContext(ctx context.Context, msg string, args ...any) {
	// if s.EnabledContext(ctx, TraceLevel) {
	pc := getpc(2, s.extraFrames)
	s.logContext(ctx, TraceLevel, false, pc, msg, args...)
	// }
}

// Verbose implements Logger.
func (s *Entry) Verbose(msg string, args ...any) {
	// if s.EnabledContext(ctx, TraceLevel) {
	pc := getpc(2, s.extraFrames)
	s.logContext(context.Background(), TraceLevel, false, pc, msg, args...)
	// }
}

func vlogctx(ctx context.Context, isformat bool, msg string, args ...any) {
	// if s.EnabledContext(ctx, TraceLevel) {
	switch s := defaultLog.(type) {
	case *logimp:
		pc := getpc(3, s.extraFrames)
		s.logContext(ctx, TraceLevel, isformat, pc, msg, args...)
	case *Entry:
		pc := getpc(3, s.extraFrames)
		s.logContext(ctx, TraceLevel, isformat, pc, msg, args...)
	}
	// }
}
