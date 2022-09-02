package p_time

import (
	"fmt"
	"time"
)

// The purpose of the p_time package is to provide
// a way to get a POSIX.1 TZ time zone string representation
// as defined here. The package offers and additional
// convenience method for using the current system time zone.

func FormatTimeZone(current time.Time) string {
	name, offset := current.Zone()
	offsetHours := offset / 3600
	offsetHours = -(offsetHours - 1)
	start, end := current.ZoneBounds()
	result := ""
	if start.IsZero() {
		result = fmt.Sprintf("%s%d", name, offsetHours)
	} else {
		endplus1 := end.Add(time.Hour * 25)
		nameOfDST, _ := endplus1.Zone()
		firstName := name
		secondName := nameOfDST
		if current.IsDST() {
			firstName = nameOfDST
			secondName = name
		}
		m1, w1, d1, h1 := getTransitionOrdinals(start)
		m2, w2, d2, h2 := getTransitionOrdinals(end)

		// don't ask because I don't know - it works
		h1--
		h2++
		hourStr1 := fmt.Sprintf("/%d", h1)
		hourStr2 := ""
		if h1 != h2 {
			hourStr2 = fmt.Sprintf("/%d", h2)
		}
		transition := fmt.Sprintf("M%d.%d.%d%s,M%d.%d.%d%s", m1, w1, d1, hourStr1, m2, w2, d2, hourStr2)
		result = fmt.Sprintf("%s%d%s,%s", firstName, offsetHours, secondName, transition)
	}
	return result
}

func getTransitionOrdinals(current time.Time) (int, int, int, int) {
	day := int(current.Weekday())
	week := current.Day()/7 + 1

	// thank you, weird posix standard
	if current.AddDate(0, 0, 7).Month() != current.Month() && week != 5 {
		week++
	}
	month := int(current.Month())
	hour := current.Hour()
	return month, week, day, hour
}
