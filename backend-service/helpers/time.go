package helpers

import "time"

func SecondsUntilMidnight() int {
	now := time.Now()
	midnight := time.Date(
		now.Year(),
		now.Month(),
		now.Day()+1,
		0, 0, 0, 0,
		now.Location(),
	)
	return int(midnight.Sub(now).Seconds())
}
