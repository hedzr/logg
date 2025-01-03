//go:build logglock
// +build logglock

package slog

type writeLock struct {
	mu sync.Mutex
}

func (s *Entry) printOut(lvl Level, msg []byte) {
	if w := s.findWriter(lvl); w != nil {
		s.muWrite.Lock()
		defer s.muWrite.Unlock()

		// if a target user-defined writer can be SetLevel, set it before writing.
		if x, ok := w.(LevelSettable); ok {
			x.SetLevel(lvl)
		}

		n, err := w.Write(msg)
		collectWrittenBytes(n)

		if err != nil && lvl != WarnLevel { // don't warn on warning to avoid infinite calls
			s.Warn("slog print log failed", "error", err)
		}
	}
}
