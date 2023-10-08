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

	l1
	l2
	l3
	l4

	Lprivacypath       //
	Lprivacypathregexp //

	l5
	l6
	l7
	l8

	LsmartJSONMode // enable JSON mode if the current output device is not terminal/tty

	LnoInterrupt     // don't interrupt app running when Fatal or Panic
	Linterruptalways // raise panic or os.Exit always even if in testing mode

	LstdFlags = Ltime | Lmicroseconds | LlocalTime | Llineno | Lcaller | Lattrs

	Ldatetimeflags = Ldate | Ltime | Lmicroseconds

	Lempty Flags = 0
)

// 	Lprettyprint                         // pretty print the object with ValueStringer.
