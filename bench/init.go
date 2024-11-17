package bench

import (
	slogg "github.com/hedzr/logg/slog"
)

func init() {
	// no source:
	slogg.RemoveFlags(slogg.Lcaller | slogg.LattrsR)
}
