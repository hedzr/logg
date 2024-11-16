//go:build metrics
// +build metrics

package slog

import (
	"atomic"
)

func collectWrittenBytes(n int) {
	atomic.AddInt64(&writtenBytes, int64(n))
}

var writtenBytes int64
