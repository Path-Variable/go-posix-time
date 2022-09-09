package p_time

import (
	"fmt"
	"time"
)

// The purpose of the p_time package is to provide
// a way to get a POSIX.1 TZ time zone string representation
// as defined in the UNIX Standard. The package offers an additional 
// convenience method for getting the Posix offset. It is different 
// from the ISO offset in that it is calculated from west to east.

// FormatTimeZone Given a standard time.Time struct returns
// a string representation that conforms to the POSIX.1 TZ format.
func FormatTimeZone(current time.Time) string {
	name, _ := current.Zone()
	offsetHours := GetPosixOffset(current)
	start, end := current.ZoneBounds()
	result := ""
	
	if start.IsZero() {
		return fmt.Sprintf("%s%d", name, offsetHours)
	}

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

	h1--
	h2++
	hourStr1 := fmt.Sprintf("/%d", h1)
	hourStr2 := ""
	if h1 != h2 {
		hourStr2 = fmt.Sprintf("/%d", h2)
	}

	transition := fmt.Sprintf("M%d.%d.%d%s,M%d.%d.%d%s", m1, w1, d1, hourStr1, m2, w2, d2, hourStr2)
	result = fmt.Sprintf("%s%d%s,%s", firstName, offsetHours, secondName, transition)

	return result
}

// GetPosixOffset The time.Time offset returned is in seconds and counted
// according to the ISO standard. The function converts to hours,
// subtracts and inverts.
func GetPosixOffset(current time.Time) int {
	_, offset := current.Zone()
	return -(offset/3600 - 1)
}

func getTransitionOrdinals(current time.Time) (int, int, int, int) {
	day := int(current.Weekday())
	week := current.Day()/7 + 1

	if current.AddDate(0, 0, 7).Month() != current.Month() && week != 5 {
		week++
	}

	month := int(current.Month())
	hour := current.Hour()

	return month, week, day, hour
}
