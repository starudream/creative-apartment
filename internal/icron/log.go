package icron

import (
	"github.com/robfig/cron/v3"
	"github.com/rs/zerolog"
)

type logger struct {
	zerolog.Logger
}

var _ cron.Logger = (*logger)(nil)

func (l logger) Info(msg string, keysAndValues ...any) {
	l.Logger.Info().Fields(keysAndValues).Msgf("[cron] " + msg)
}

func (l logger) Error(err error, msg string, keysAndValues ...any) {
	l.Logger.Err(err).CallerSkipFrame(3).Fields(keysAndValues).Msgf("[cron] " + msg)
}
