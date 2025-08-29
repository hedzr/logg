//go:build !verbose
// +build !verbose

package slog

import (
	"context"
)

// VerboseContext implements Logger.
func (s *Entry) VerboseContext(ctx context.Context, msg string, args ...any) {}

// Verbose implements Logger.
func (s *Entry) Verbose(msg string, args ...any) {}

func vlogctx(ctx context.Context, isformat bool, msg string, args ...any) {}
