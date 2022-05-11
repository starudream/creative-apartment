package itask

import (
	"github.com/panjf2000/ants/v2"
	"github.com/rs/zerolog"
)

type logger struct {
	zerolog.Logger
}

var _ ants.Logger = (*logger)(nil)

func (l logger) Printf(format string, args ...interface{}) {
	l.Logger.Info().Msgf("[task] "+format, args...)
}
