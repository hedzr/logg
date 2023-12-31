// Package logg is a golang library for providing logging functions
// with more effectively, friendly programmatic API.
//
// logg/slog is an opt-in copy from log/slog. Different with originals,
// logg has more verbs and features and colored text outputting
// with out-of-the-box. Have a see at https://github.com/hedzr/logg.
//
// To use this logging library, import it as:
//
//	import "github.com/hedzr/logg/slog"
//
//	var logger = slog.New("my-app").WithLevel(slog.Debug).WithJSONMode()
//	logger.Info("info message here", "attr1", 3, "attr2", false, "attr3", "text details")
//	logger.Println() // just an empty line
//	logger.Println("text message", attrs...)
//	logger.Print("text message", attrs...)
//	logger.Debug("debug message", attrs...)
//	logger.Trace("trace message", attrs...)
//	logger.Warn("warn message", attrs...)
//	logger.Error("error message", attrs...)
//	logger.Fatal("fatal message", attrs...)
//	logger.Panic("panic message", attrs...)
//	logger.OK("ok message", attrs...)
//	logger.Success("success message", attrs...)
//	logger.Fail("fail message", attrs...)
//	logger.Verbose("verbose message", attrs...) // only work for build tag 'verbose' defined
//
//	var subl = logger.New("child1").With(attrs...)
//	subl.Debug("debug")
//
// For more detail, please take a look at:
//
//   - https://github.com/hedzr/logg
//   - https://pkg.go.dev/github.com/hedzr/logg
//   - https://deps.dev/go/github.com%2Fhedzr%2Flogg
package logg
