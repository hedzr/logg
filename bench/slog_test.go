package bench

import (
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
	logger := slogg.New().SetWriter(os.Stdout).SetLevel(slogg.DebugLevel).SetJSONMode()
	msg, attrs := getMessage(0), fakeLoggArgs()
	logger.Info(msg, attrs...)

	logger.WithJSONMode(false)
	logger.Info(msg, attrs...)
}
