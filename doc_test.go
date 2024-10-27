package logg

import (
	"testing"

	logz "github.com/hedzr/logg/slog"
)

func TestSlog(t *testing.T) {
	t.Log("OK")

	logz.SetLevel(logz.InfoLevel)

	logz.Info("Hello", "target", "world")
	logz.Default().SetColorMode(false)
	logz.Info("Hello", "target", "world")
	logz.Default().SetJSONMode(true)
	logz.Info("Hello", "target", "world")
}
