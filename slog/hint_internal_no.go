//go:build !hint
// +build !hint

package slog

// hintInternal prints the internal errors which has been
// ignored since they are unused for the logging flow.
//
// In general, the error in formatting a invalid number
// and more similar cases are ignored by default.
//
// But you can trace these error by build with tags 'hint'
func hintInternal(err error, msg string) { //nolint:revive,unparam
	// this function will be removed completely when building a release copy
}

func raiseerror(msg string) {} //nolint:revive,unparam,unused
