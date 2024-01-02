package slog

import (
	"os"
	"regexp"

	"github.com/hedzr/is"
	"github.com/hedzr/logg/slog/internal/strings"
)

func init() {
	onceInit.Do(func() {
		lvlCurrent = WarnLevel

		codeHostingProvidersMap = map[string]string{
			"github.com":    "GH",
			"gitlab.com":    "GL",
			"gitee.com":     "GT",
			"bitbucket.com": "BB",
		}

		homeDir, _ = os.UserHomeDir()
		currDir, _ = os.Getwd()
		knownPathMap = map[string]string{
			homeDir: "~",
			currDir: ".",
		}
		knownPathRegexpMap = append(knownPathRegexpMap, regRepl{
			expr: regexp.MustCompile(`/Volumes/[^/]+/`),
			repl: `~`,
		})

		if is.InDebugging() || is.InTesting() || strings.StringToBool(os.Getenv("DEBUG")) {
			lvlCurrent = DebugLevel
		}

		defaultWriter = newDualWriter()
		defaultLog = newDetachedLogger()
	})
}

var (
	codeHostingProvidersMap map[string]string

	knownPathMap map[string]string

	knownPathRegexpMap []regRepl

	homeDir, currDir string
)

type regRepl struct {
	expr *regexp.Regexp
	repl string
}
