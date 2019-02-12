package goexec

import (
	"github.com/pkg/errors"
	"github.com/robfig/cron"
)

type CronSpec int

const (
	HOURLY CronSpec = iota
	DAILY
	WEEKLY
	MONTHLY
	YEARLY
)

func (s CronSpec) String() string {
	// declare an array of strings
	// ... operator counts how many
	// items in the array (7)
	specs := [...]string{
		"HOURLY",
		"DAILY",
		"WEEKLY",
		"MONTHLY",
		"YEARLY"}
	// → `day`: It's one of the
	// values of Weekday constants.
	// If the constant is Sunday,
	// then day is 0.
	// prevent panicking in case of
	// `day` is out of range of Weekday
	if s < HOURLY || s > YEARLY {
		panic(errors.New("unknown cron spec"))
	}
	// return the name of a Weekday
	// constant from the names array
	// above.
	str := specs[s]
	switch str {
	case "HOURLY":
		str = "@hourly "
	case "DAILY":
		str = "@daily"
	case "WEEKLY":
		str = "@weekly"
	case "MONTHLY":
		str = "@monthly"
	case "YEARLY":
		str = "@yearly"
	}
	return str
}

func (c *Command) Cron(spec CronSpec, fn func()) {
	cronz := cron.New()
	c.Panic(cronz.AddFunc(spec.String(), fn), "failed to add cron")
}
