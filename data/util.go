package data

import (
	"time"
)

func sameDay(t1, t2 time.Time) bool {
	if t1.Day() == t2.Day() && t1.Month() == t2.Month() && t1.Year() == t2.Year() {
		return true
	}
	return false
}
