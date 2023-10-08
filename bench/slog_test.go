package bench

import (
	"testing"

	slogg "github.com/hedzr/logg/slog"
)

func Test1(t *testing.T) {
	logger := newLogg()
	msg, attrs := getMessage(0), fakeLoggArgs()
	logger.Info(msg, attrs...)
}

func Test2(t *testing.T) {
	logger := newLogg().WithJSONMode()
	msg, attrs := getMessage(0), fakeLoggArgs()
	logger.Info(msg, attrs...)
}

func TestUseJSON(t *testing.T) {
	logger := slogg.New().WithLevel(slogg.DebugLevel).WithJSONMode()
	msg, attrs := getMessage(0), fakeLoggArgs()
	logger.Info(msg, attrs...)

	logger.WithJSONMode(false)
	logger.Info(msg, attrs...)
}
