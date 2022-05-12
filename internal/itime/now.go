package itime

import (
	"time"

	"github.com/jinzhu/now"
)

var config = &now.Config{
	WeekStartDay: time.Monday,
	TimeLocation: time.Local,
	TimeFormats:  now.TimeFormats,
}

func T() *now.Config {
	return config
}

func Now() *now.Now {
	return config.With(time.Now())
}

func New(t time.Time) *now.Now {
	return config.With(t)
}

func With(t time.Time) *now.Now {
	return config.With(t)
}
