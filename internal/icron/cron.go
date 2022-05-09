package icron

import (
	"time"

	"github.com/robfig/cron/v3"
	"github.com/rs/zerolog/log"
	"github.com/spf13/viper"
)

func New() *cron.Cron {
	nlog := &logger{log.Logger}
	if viper.GetBool("debug") {
		nlog = &logger{log.Logger.With().CallerWithSkipFrameCount(3).Logger()}
	}

	c := cron.New(
		cron.WithLocation(time.Local),
		cron.WithSeconds(),
		cron.WithChain(cron.Recover(nlog)),
		cron.WithLogger(nlog),
	)

	return c
}
