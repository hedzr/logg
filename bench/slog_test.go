package bench

import (
	"context"
	"errors"
	"fmt"
	"os"
	"testing"

	"github.com/hedzr/is"

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
	logger.Println()

	logger.SetJSONMode(false)
	logger.InfoContext(ctx, msg, attrs...)

	logger.SetColorMode()
	logger.InfoContext(ctx, msg, attrs...)
}

func TestErrors(t *testing.T) {
	ctx, msg, attrs := context.Background(), getMessage(0), fakeLoggArgs()
	logger := newLoggTextMode().Set(attrs...).
		SetWriter(os.Stdout).SetLevel(slogg.InfoLevel)
	logger.InfoContext(ctx, msg)

	err := errors.New("test error")
	logger.Error("error occurred", "err", err)

	logger = newLoggTextMode().SetWriter(os.Stdout).SetLevel(slogg.InfoLevel)

	is.SetDebugMode(true)

	logger.SetColorMode(false)
	logger.Error("error occurred", "testing", is.InTesting(), "err", err)

	logger.SetColorMode()
	logger.Error("error occurred", "err", err)

	func() {
		defer func() {
			if err := recover(); err != nil {
				logger.SetColorMode(false)
				logger.Error("panic caught", "testing", is.InTesting(), "err", err)
				logger.SetColorMode()
				logger.Error("panic caught", "err", err)
			}
		}()

		panic("panic test")
	}()

	fmt.Printf("err: %+v\n", err)
}

func Test4(t *testing.T) {
	ctx, msg, attrs := context.Background(), getMessage(0), fakeLoggArgs()
	logger := newLoggTextMode().Set(attrs...).
		SetWriter(os.Stdout).SetLevel(slogg.InfoLevel)
	logger.InfoContext(ctx, msg)
}
