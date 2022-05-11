package icron

import (
	"time"

	"github.com/robfig/cron/v3"
	"github.com/rs/zerolog/log"
)

func New() *cron.Cron {
	nlog := &logger{log.Logger}

	c := cron.New(
		cron.WithLocation(time.Local),
		cron.WithSeconds(),
		cron.WithChain(cron.Recover(nlog)),
		cron.WithLogger(nlog),
	)

	return c
}
