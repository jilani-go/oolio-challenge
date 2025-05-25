package utils

import "time"

// NormalizeToMidnightUTC normalizes a time to midnight UTC
func NormalizeToMidnightUTC(t time.Time) time.Time {
	return time.Date(t.Year(), t.Month(), t.Day(), 0, 0, 0, 0, time.UTC)
}
