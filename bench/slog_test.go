package bench

import (
	"context"
	"os"
	"testing"

	slogg "github.com/hedzr/logg/slog"
)

func Test1(t *testing.T) {
	logger := newLogg().SetWriter(os.Stdout).SetColorMode().SetLevel(slogg.InfoLevel)
	msg, attrs := getMessage(0), fakeLoggArgs()
	logger.Info(msg, attrs...)
}

func Test2(t *testing.T) {
	logger := newLogg().SetWriter(os.Stdout).SetColorMode(false).SetLevel(slogg.InfoLevel)
	msg, attrs := getMessage(0), fakeLoggArgs()
	logger.Info(msg, attrs...)
}

func TestUseJSON(t *testing.T) {
	defer slogg.SaveFlagsAndMod(slogg.Lcaller)()

	logger := slogg.New(slogg.With("int", 3)).
		SetWriter(os.Stdout).
		SetLevel(slogg.DebugLevel).
		SetJSONMode()

	ctx, msg, attrs := context.Background(), getMessage(0), fakeLoggArgs()
	logger.InfoContext(ctx, msg, attrs...)

	logger.SetJSONMode(false)
	logger.InfoContext(ctx, msg, attrs...)

	logger.SetColorMode()
	logger.InfoContext(ctx, msg, attrs...)
}

func Test3(t *testing.T) {
	ctx, msg, attrs := context.Background(), getMessage(0), fakeLoggArgs()
	logger := newLoggTextMode().Set(attrs...).
		SetWriter(os.Stdout).SetLevel(slogg.InfoLevel)
	logger.InfoContext(ctx, msg)
}

func Test4(t *testing.T) {
	ctx, msg, attrs := context.Background(), getMessage(0), fakeLoggArgs()
	logger := newLoggTextMode().Set(attrs...).
		SetWriter(os.Stdout).SetLevel(slogg.InfoLevel)
	logger.InfoContext(ctx, msg)
}
