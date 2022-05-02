package helpers

import "time"

func WithinLimit(limit int, target int) bool {
	return limit-target > 0
}

func AreDatesWithinInterval(from, to time.Time, interval int) bool {
	diffMinutes := to.Sub(from).Minutes()
	return diffMinutes < float64(interval)
}
