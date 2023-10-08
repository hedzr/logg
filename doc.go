// Package logg is a golang library for providing logging functions with more effectively, friendly programmatic API.
//
// To use this logging library, import it as:
//
//	import "github.com/hedzr/logg/slog"
//
//	var logger = slog.New().WithLevel(slog.Debug).WithJSONMode()
//	logger.Info("info message here", "attr1", 3, "attr2", false, "attr3", "text details")
//	logger.Println()
//	logger.Debug("debug message")
//
// For more detail, please take a look at:
//
//   - https://github.com/hedzr/logg
//   - https://pkg.go.dev/github.com/hedzr/logg
//   - https://deps.dev/go/github.com%2Fhedzr%2Flogg
package logg
