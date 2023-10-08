//go:build verbose
// +build verbose

package slog

import (
	"context"
)

// VerboseContext implements Logger.
func (s *entry) VerboseContext(ctx context.Context, msg string, args ...any) {
	// if s.EnabledContext(ctx, TraceLevel) {
	pc := getpc(2, s.extraFrames)
	s.logContext(ctx, TraceLevel, pc, msg, args...)
	// }
}

// Verbose implements Logger.
func (s *entry) Verbose(msg string, args ...any) {
	// if s.EnabledContext(ctx, TraceLevel) {
	pc := getpc(2, s.extraFrames)
	s.logContext(context.Background(), TraceLevel, pc, msg, args...)
	// }
}

func vlogctx(ctx context.Context, msg string, args ...any) {
	// if s.EnabledContext(ctx, TraceLevel) {
	switch s := defaultLog.(type) {
	case *logimp:
		pc := getpc(3, s.extraFrames)
		s.logContext(ctx, TraceLevel, pc, msg, args...)
	case *entry:
		pc := getpc(3, s.extraFrames)
		s.logContext(ctx, TraceLevel, pc, msg, args...)
	}
	// }
}
