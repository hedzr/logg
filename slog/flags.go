package slog

type Flags int64

const (
	Ldate              Flags = 1 << iota // do print date part
	Ltime                                // do print time part
	Lmicroseconds                        // do print microseconds part
	LlocalTime                           // use local time instead of UTC
	Lattrs                               // do print Attr key-value pairs
	LattrsR                              // collect Attr along with entry and Logger chains
	Llineno                              // do print caller info (file:line)
	Lcaller                              // do print Caller information (function)
	Lcallerpackagename                   // do print the package name of caller, such as: GH/hedzr/logg/slog_test.TestSlogLogfmt

	l1 //nolint:unused
	l2 //nolint:unused
	l3 //nolint:unused
	l4 //nolint:unused

	Lprivacypath       // Privacy hardening flag. A string slice will be used for hiding disk pathname.
	Lprivacypathregexp // Privacy hardening flag. A regexp pattern slice will be used.

	l5 //nolint:unused
	l6 //nolint:unused
	l7 //nolint:unused
	l8 //nolint:unused

	LsmartJSONMode // enable JSON mode if the current output device is not terminal/tty

	LnoInterrupt     // don't interrupt app running when Fatal or Panic
	Linterruptalways // raise panic or os.Exit always even if in testing mode

	// LstdFlags is the default flags when unboxed
	LstdFlags = Ltime | Lmicroseconds | LlocalTime | Llineno | Lcaller | Lattrs |
		Lprivacypath | Lprivacypathregexp

	Ldatetimeflags = Ldate | Ltime | Lmicroseconds // for timestamp formatting

	Lempty Flags = 0 // for algor
)

// 	Lprettyprint                         // pretty print the object with ValueStringer.
