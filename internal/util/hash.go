package util

import (
	"time"
)

// Hash generates code from date to hour as unix hash code
func Hash(d time.Time) int64 {
	temp := time.Date(d.Year(), d.Month(), d.Day(), d.Hour(), 0, 0, 0, time.UTC)
	return temp.Unix()
}

// HashHours returns unix days by date time range
func HashHours(start, end time.Time) []int64 {
	start = time.Date(start.Year(), start.Month(), start.Day(), start.Hour(), 0, 0, 0, time.UTC)
	end = time.Date(end.Year(), end.Month(), end.Day(), end.Hour(), 0, 0, 0, time.UTC)

	var hashes []int64
	for start.Before(end) {
		hashes = append(hashes, start.Unix())
		start = start.Add(1 * time.Hour)
	}

	return append(hashes, end.Unix())
}
