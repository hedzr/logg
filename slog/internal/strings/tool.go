package strings

import (
	"strings"
)

func StringToBool(val string, defaultVal ...bool) (ret bool) {
	return toBool(val, defaultVal...)
}

func toBool(val string, defaultVal ...bool) (ret bool) {
	// ret = ToBool(val, defaultVal...)
	switch strings.ToLower(val) {
	case "1", "y", "t", "yes", "true", "ok", "on":
		ret = true
	case "":
		for _, vv := range defaultVal {
			ret = vv
		}
	}
	return
}
