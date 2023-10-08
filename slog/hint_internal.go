//go:build hint
// +build hint

package slog

// hintInternal prints the internal errors which has been
// ignored since they are unused for the logging flow.
//
// In general, the error in formatting a invalid number
// and more similar cases are ignored by default.
//
// But you can trace these error by build with tags 'hint'
func hintInternal(err error, msg string) {
	logctx(ErrorLevel, "internal error", "hint", msg, "error", err)
}

func raiseerror(msg string) { panic(msg) }
