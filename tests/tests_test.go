package tests

import (
	logslog "log/slog"
	"testing"

	"github.com/hedzr/logg/slog"
)

func TestSlogAdapter(t *testing.T) {
	// t.Log("slog")

	l := slog.New("standalone-logger-for-app",
		slog.NewAttr("attr1", 2),
		slog.NewAttrs("attr2", 3, "attr3", 4.1),
		"attr4", true, "attr3", "string",
		slog.WithLevel(slog.AlwaysLevel),
	)
	defer l.Close()

	sub1 := l.New("sub1").With("logger", "sub1")
	sub2 := l.New("sub2").With("logger", "sub2").WithLevel(slog.InfoLevel)

	// create a log/slog logger HERE
	logger := logslog.New(slog.NewSlogHandler(l, &slog.HandlerOptions{
		NoColor:  false,
		NoSource: true,
		JSON:     false,
		Level:    slog.DebugLevel,
	}))

	t.Logf("logger: %v", logger)
	t.Logf("l: %v", l)

	// and logging with log/slog
	logger.Debug("hi debug", "AA", 1.23456789)
	logger.Info(
		"incoming request",
		logslog.String("method", "GET"),
		logslog.String("path", "/api/user"),
		logslog.Int("status", 200),
	)

	// now using our logg/slog interface
	sub1.Debug("hi debug", "AA", 1.23456789)
	sub2.Debug("hi debug", "AA", 1.23456789)
	sub2.Info("hi info", "AA", 1.23456789)
}
