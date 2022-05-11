package icron

import (
	"github.com/robfig/cron/v3"
)

func WrapError(_ cron.EntryID, err error) {
	if err != nil {
		panic(err)
	}
}
