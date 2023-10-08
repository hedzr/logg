//go:build !verbose
// +build !verbose

package slog

import (
	"context"
)

// VerboseContext implements Logger.
func (s *entry) VerboseContext(ctx context.Context, msg string, args ...any) {}

// Verbose implements Logger.
func (s *entry) Verbose(msg string, args ...any) {}

func vlogctx(ctx context.Context, msg string, args ...any) {}
