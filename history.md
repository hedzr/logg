# HISTORY

initial commit

- **Structured**: plain output (logfmt or colorful) or JSON format.
- **Leveled/Contextual**: sub-`Log` and attributes.
- Privacy has high order:
  - Harden filepath, shorten package name, etc.
  - Implements MarshalObjectValue and/or MarshalObjectArray to security sensitive fields.
- Efficient caller stack calculating
- Adapted into `log/slog` via `NewSlogHandler`
- Customizing value stringer/formatter by `ValueStringer`
- More...

state:

cpu: Intel(R) Core(TM) i9-9880H CPU @ 2.30GHz
BenchmarkAccumulatedContext
    integrated_test.go:140: Logging with some accumulated context. [BenchmarkAccumulatedContext]
BenchmarkAccumulatedContext/hedzr/logg/slog_TEXT
BenchmarkAccumulatedContext/hedzr/logg/slog_TEXT-16         	  220863	      5650 ns/op	    7825 B/op	      71 allocs/op
BenchmarkAccumulatedContext/hedzr/logg/slog_COLOR
BenchmarkAccumulatedContext/hedzr/logg/slog_COLOR-16        	  198340	      5981 ns/op	    8232 B/op	     135 allocs/op
BenchmarkAccumulatedContext/hedzr/logg/slog_JSON
BenchmarkAccumulatedContext/hedzr/logg/slog_JSON-16         	  198314	      6010 ns/op	    7824 B/op	      71 allocs/op
BenchmarkAccumulatedContext/slog
BenchmarkAccumulatedContext/slog-16                         	 7217726	       169.2 ns/op	       0 B/op	       0 allocs/op
BenchmarkAccumulatedContext/slog.LogAttrs
BenchmarkAccumulatedContext/slog.LogAttrs-16                	 6900033	       169.8 ns/op	       0 B/op	       0 allocs/op



old state:

BenchmarkAccumulatedContextLoggOnly/hedzr/logg/slog-16            166149              7509 ns/op            8785 B/op         87 allocs/op
BenchmarkAccumulatedContextLoggOnly/hedzr/logg/slog-16            183168              7900 ns/op            8785 B/op         87 allocs/op

> Because log/slog outputs json without caller info by default,
> So we updated logg/slog newLogg() to disable caller info
> at testing. Hence the result get a little improvement.
BenchmarkAccumulatedContextLoggOnly/hedzr/logg/slog-16            239474              5841 ns/op            7809 B/op         71 allocs/op
> Also collecting parent's attrs is disabled:
BenchmarkAccumulatedContextLoggOnly/hedzr/logg/slog-16            257452              6409 ns/op            7809 B/op         71 allocs/op

BenchmarkAccumulatedContextLoggOnly/hedzr/logg/slog-16         	  196978	      6551 ns/op	    7825 B/op	      71 allocs/op
