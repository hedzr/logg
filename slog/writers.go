package slog

import (
	"errors"
	"io"
	"os"
)

type dualWriter struct {
	Normal  LWs
	Error   LWs
	leveled map[Level]LWs
}

func newDualWriter() *dualWriter {
	return (&dualWriter{}).Reset()
}

func (s *dualWriter) Close() (err error) {
	if e := s.Normal.Close(); e != nil {
		err = errors.Join(err, e)
	}
	if e := s.Error.Close(); e != nil {
		err = errors.Join(err, e)
	}

	// var ec = errors.New("dualWriter.close")
	// defer ec.Defer(&err)
	// ec.Attach(s.Normal.Close())
	// ec.Attach(s.Error.Close())
	return
}

func (s *dualWriter) Write(p []byte) (n int, err error) {
	// TO/DO implement me
	// /panic("implement me")

	if ni, e := s.Normal.Write(p); e != nil {
		err = errors.Join(err, e)
	} else {
		n += ni
	}

	// var ec = errors.New("dualWriter.Write")
	// var e error
	// defer ec.Defer(&err)
	// n, e = s.Normal.Write(p)
	// ec.Attach(e)
	return
}

// var discardWriter = LWs{&logwr{io.Discard}}

type discard struct{}

func (discard) Write(b []byte) (n int, err error) { return len(b), nil }

var discardWriter = LWs{&logwr{discard{}}}

func (s *dualWriter) Get(lvl Level) (w LWs) {
	if lvl == OffLevel {
		return discardWriter
	}
	if s.leveled != nil {
		if ed, ok := s.leveled[lvl]; ok && len(ed) > 0 {
			return ed
		}
	}
	if _, ok := mLevelUseErrorDevice[lvl]; ok {
		return s.Error
	}
	return s.Normal
}

func (s *dualWriter) Set(w io.Writer) {
	if w != nil {
		s.Normal = nil
		s.Add(w)
	}
}

func (s *dualWriter) SetWriter(w io.Writer) {
	if w != nil {
		s.Normal = nil
		if lw, ok := w.(LogWriter); ok {
			s.Normal = append(s.Normal, lw)
			return
		}
		s.Normal = append(s.Normal, &logwr{w})
	}
}

func (s *dualWriter) SetErrorWriter(w io.Writer) {
	if w != nil {
		s.Error = nil
		if lw, ok := w.(LogWriter); ok {
			s.Error = append(s.Error, lw)
			return
		}
		s.Error = append(s.Error, &logwr{w})
	}
}

func (s *dualWriter) Add(w io.Writer) {
	if w != nil {
		if lw, ok := w.(LogWriter); ok {
			s.Normal = append(s.Normal, lw)
			return
		}
		s.Normal = append(s.Normal, &logwr{w})
	}
}

func (s *dualWriter) AddErrorWriter(w io.Writer) {
	if w != nil {
		if lw, ok := w.(LogWriter); ok {
			s.Error = append(s.Error, lw)
			return
		}
		s.Error = append(s.Error, &logwr{w})
	}
}

func (s *dualWriter) AddLevelWriter(lvl Level, w io.Writer) {
	if w != nil {
		if s.leveled == nil {
			s.leveled = make(map[Level]LWs)
		}
		if lw, ok := w.(LogWriter); ok {
			s.leveled[lvl] = append(s.leveled[lvl], lw)
			return
		}
		s.leveled[lvl] = append(s.leveled[lvl], &logwr{w})
	}
}

func (s *dualWriter) RemoveLevelWriter(lvl Level, w io.Writer) {
	if w != nil {
		if s.leveled == nil {
			s.leveled = make(map[Level]LWs)
		}
		if lw, ok := s.leveled[lvl]; ok {
			for i, wr := range lw {
				if wr == w {
					s.leveled[lvl] = append(s.leveled[lvl][:i], s.leveled[lvl][i+1:]...)
					break
				}
			}
		}
	}
}

func (s *dualWriter) ResetLevelWriter(lvl Level) {
	if s.leveled == nil {
		s.leveled = make(map[Level]LWs)
	}
	if _, ok := s.leveled[lvl]; ok {
		s.leveled[lvl] = nil
	}
}

func (s *dualWriter) ResetLevelWriters() {
	s.leveled = nil
}

func (s *dualWriter) Clear() {
	s.Normal = nil
	s.Error = nil
}

func (s *dualWriter) Reset() *dualWriter {
	s.Normal = []LogWriter{&filewr{os.Stdout}}
	s.Error = []LogWriter{&filewr{os.Stderr}}
	s.leveled = nil
	return s
}

//

//

//

type LWs []LogWriter

func (s LWs) Close() (err error) {
	for _, w := range s {
		if w != nil {
			if e := w.Close(); e != nil {
				err = errors.Join(err, e)
			}
		}
	}

	// var ec = errors.New("[]LogWriter.Close")
	// defer ec.Defer(&err)
	// for _, w := range s {
	// 	if w != nil {
	// 		ec.Attach(w.Close())
	// 	}
	// }
	return
}

func (s LWs) Write(p []byte) (n int, err error) {
	// TO/DO implement me
	// /panic("implement me")

	for _, w := range s {
		if ni, e := w.Write(p); e != nil {
			err = errors.Join(err, e)
		} else {
			n += ni
		}
	}

	// var ec = errors.New("[]LogWriter.Write")
	// var e error
	// defer ec.Defer(&err)
	// for _, w := range s {
	// 	n, e = w.Write(p)
	// 	ec.Attach(e)
	// }
	return
}

//

//

//

type logwr struct {
	io.Writer
}

func (s *logwr) Close() error {
	if c, ok := s.Writer.(io.Closer); ok {
		return c.Close()
	}
	if c, ok := s.Writer.(interface{ Close() }); ok {
		c.Close()
	}
	return nil
}

type filewr struct {
	*os.File
}

func (s *filewr) Close() (err error) {
	if s.File != nil {
		err = s.File.Close()
		s.File = nil
	}
	return
}

func NewLogWriter(w io.Writer) *logwr {
	return &logwr{w}
}

func NewFileWriter(pathname string) *filewr {
	f, err := os.OpenFile(pathname, os.O_RDWR|os.O_CREATE, 0755)
	if err != nil {
		Fatal("cannot create logging file", "error", err, "pathname", pathname)
	}
	s := &filewr{f}
	return s
}
