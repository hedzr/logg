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
			"codeberg.org":  "CB",
			"framagit.org":  "FR",
		}

		// for hardening your privacy,
		homeDir, _ = os.UserHomeDir()
		currDir, _ = os.Getwd()
		knownPathMap = map[string]string{
			homeDir: "~", // convert $homeDir to '~'
			currDir: ".", // convert abs-path currDir prefix to relative-path (started with '.')
		}
		// A tilde directory can be represented as `~work/...` if `work`
		// was defined by:
		//    hash -d work="/Volumes/VolWork/work"
		//    ls -la ~work/
		knownPathRegexpMap = append(knownPathRegexpMap, regRepl{
			expr: regexp.MustCompile(`/Volumes/[^/]+/`),
			repl: `~`, // convert tilde folder (in bash/zsh) to abbreviate mode.
		})

		if is.Tracing() || is.TraceMode() {
			lvlCurrent = TraceLevel
		} else if is.DebuggerAttached() || inTesting || is.DebugBuild() || is.DebugMode() || strings.StringToBool(os.Getenv("DEBUG")) {
			lvlCurrent = DebugLevel
			RemoveFlags(Lprivacypathregexp) // disable tilde directory to make the logging msg clickable
		} else {
			RemoveFlags(LlocalTime)
		}
		is.Env().SetOnDevModeChanged(func(mod bool, level int) {
			// if mod {
			if lvlCurrent < InfoLevel {
				lvlCurrent = InfoLevel
				Debug("[logz][onDevModeChanged] dev-mode changed, set to InfoLevel", "mode", mod, "level", level, "log-level", lvlCurrent)
			} else {
				Debug("[logz][onDevModeChanged] dev-mode changed, no action", "mode", mod, "level", level, "log-level", lvlCurrent)
				AddFlags(LlocalTime)
			}
			// }
		})
		is.SetOnDebugChanged(func(mod bool, level int) {
			// if mod {
			if lvlCurrent < DebugLevel {
				lvlCurrent = DebugLevel
				Debug("[logz][onDebugChanged] debug mode changed", "mode", mod, "level", level, "log-level", lvlCurrent)
				AddFlags(LlocalTime)
			}
			// }
		})
		is.SetOnTraceChanged(func(mod bool, level int) {
			// if mod {
			// Print("[logz][onTraceChanged]", "current", lvlCurrent, "new-trace-mode", mod)
			if lvlCurrent < TraceLevel {
				lvlCurrent = TraceLevel
				Trace("[logz][onTraceChanged] trace mode changed", "mode", mod, "level", level, "log-level", lvlCurrent)
				AddFlags(LlocalTime)
			}
			// }
		})

		defaultWriter = newDualWriter()
		defaultLog = newDetachedLogger()

		warmupAttrs := poolAttrs.Get()
		poolAttrs.Put(warmupAttrs)

		warmupPC := poolPrintCtx.Get()
		poolPrintCtx.Put(warmupPC)
	})
}

var (
	codeHostingProvidersMap map[string]string // eg: "github.com" -> "GH"

	knownPathMap map[string]string // eg: "$HOME" -> "~"

	knownPathRegexpMap []regRepl // eg: "/Volumes/(.+)" -> "~$1"

	homeDir, currDir string
)

type regRepl struct {
	expr *regexp.Regexp
	repl string
}
