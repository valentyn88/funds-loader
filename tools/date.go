package tools

import (
	"time"
)

// IsSameDay checks whether the same date or not.
func IsSameDay(date1, date2 time.Time) bool {
	y1, m1, d1 := date1.Date()
	y2, m2, d2 := date2.Date()

	return y1 == y2 && m1 == m2 && d1 == d2
}

// IsSameWeek checks whether the same week or not.
func IsSameWeek(date1, date2 time.Time) bool {
	y1, w1 := date1.ISOWeek()
	y2, w2 := date2.ISOWeek()

	return y1 == y2 && w1 == w2
}

// ParseDate parses input date.
func ParseDate(val string) (time.Time, error) {
	return time.Parse("2006-01-02T15:04:05Z", val)
}
