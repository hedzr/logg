package slog

import (
	"fmt"
	"os"
	"path/filepath"
	"runtime"
	"strings"
)

// Source describes the location of a line of source code.
type Source struct {
	// Function is the package path-qualified function name containing the
	// source line. If non-empty, this string uniquely identifies a single
	// function in the program. This may be the empty string if not known.
	Function string `json:"function"`
	// File and Line are the file name and line number (1-based) of the source
	// line. These may be the empty string and zero, respectively, if not known.
	File string `json:"file"`
	Line int    `json:"line"`
}

func (s Source) toGroup() (as Attr) {
	as = &gkvp{"source", Attrs{
		&kvp{"function", s.Function},
		&kvp{"file", s.File},
		&kvp{"line", s.Line},
	}}
	// as = Group("source",
	// 	NewAttr("function", s.Function),
	// 	NewAttr("file", s.File),
	// 	NewAttr("line", s.Line))
	return
}

func getpc(skip int, extra int) (pc uintptr) {
	var pcs [1]uintptr
	runtime.Callers(skip+extra+1, pcs[:])
	pc = pcs[0]
	return
}

func getpcsource(pc uintptr) Source {
	frames := runtime.CallersFrames([]uintptr{pc})
	frame, _ := frames.Next()
	return Source{
		Function: frame.Function,
		File:     checkpath(frame.File),
		Line:     frame.Line,
	}
}

func stack(skip, nFrames int) string {
	pcs := make([]uintptr, nFrames+1)
	n := runtime.Callers(skip+1, pcs)
	if n == 0 {
		return "(no stack)"
	}
	frames := runtime.CallersFrames(pcs[:n])
	var b strings.Builder
	i := 0
	for {
		frame, more := frames.Next()
		fmt.Fprintf(&b, "called from %s (%s:%d)\n", frame.Function, frame.File, frame.Line)
		if !more {
			break
		}
		i++
		if i >= nFrames {
			fmt.Fprintf(&b, "(rest of stack elided)\n")
			break
		}
	}
	return b.String()
}

// Safety applies logg/slog's security policies on given file pathname.
//
// - When you have preset tiled pathes, such as `~work' pointed
// to `/Volumes/vWork/work', a safety path of `/Volumes/vWork/work/a.go'
// will be translated as `~work/a.go'.
//
// It relies on Lprivacypath is set. See AddFlags.
//
// These ways can merge flags into internal settings:
//
//	AddFlags(Lprivacypathregexp | Lprivacypath)
//	AddFlags(Lprivacypathregexp, Lprivacypath)
//
// - When Lprivacypath | Lprivacypathregexp enabled, the
// user's homedir will be translated to '~' to avoid user account
// name leaked.
//
// You may add more policies with AddKnownPathMapping and
// AddKnownPathRegexpMapping.
func Safety(file string) string {
	return checkpath(file)
}

// SafetyFiles is an overrided prototype of Safety.
func SafetyFiles(files []string) (ret []string) {
	for _, f := range files {
		ret = append(ret, checkpath(f))
	}
	return
}

func checkpath(file string) string {
	// if s.curdir == "" {
	// 	s.curdir, _ = os.Getwd()
	// }
	// if strings.HasPrefix(file, s.curdir) {
	// 	file= file[len(s.curdir)+1:]
	// }

	privfile := file
	if IsAnyBitsSet(Lprivacypath) {
		for k, v := range knownPathMap {
			if strings.HasPrefix(privfile, k) {
				privfile = strings.ReplaceAll(privfile, k, v)
			}
		}

		if IsAnyBitsSet(Lprivacypathregexp) {
			for _, rpl := range knownPathRegexpMap {
				if rpl.expr.MatchString(file) {
					privfile = rpl.expr.ReplaceAllString(privfile, rpl.repl)
				}
			}
		} else {
			if strings.HasPrefix(privfile, "/Volumes/") {
				if pos := strings.IndexRune(privfile[9:], '/'); pos >= 0 {
					privfile = "~" + privfile[9+pos:]
				}
			}
		}
	}
	cwd, _ := os.Getwd()
	relfile, _ := filepath.Rel(cwd, file)
	if l := len(relfile); l > 0 && l < len(privfile) {
		return relfile
	}
	return privfile
}

func checkedfuncname(name string) string {
	// name := s.Function
	if IsAnyBitsSet(Lcallerpackagename) {
		for k, v := range codeHostingProvidersMap {
			name = strings.ReplaceAll(name, k, v) // replace github.com with "GH", ...
		}
	} else {
		if pos := strings.LastIndex(name, "/"); pos >= 0 {
			name = name[pos+1:] // strip the leading package names, eg: "GH/hedzr/logg/slog/" will be removed
		}
	}
	return name
}
