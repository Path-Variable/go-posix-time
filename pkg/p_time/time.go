// Package p_time formats a Go time.Time into a POSIX.1 TZ string such as
// "CET-1CEST,M3.5.0,M10.5.0/3", the format expected by devices that cannot
// carry the IANA time-zone database (IP cameras, embedded systems).
//
// The format is specified in the Open Group Base Specifications, Issue 7,
// section 8.3 "Other Environment Variables": std offset[dst[offset][,rule]].
// Note that POSIX offsets are counted west of UTC, the opposite of ISO 8601.
package p_time

import (
	"fmt"
	"strings"
	"time"
)

// posixDefaultTransitionHour is the hour POSIX assumes when a rule omits it.
const posixDefaultTransitionHour = 2

// FormatTimeZone returns the POSIX.1 TZ representation of the zone that
// current belongs to, including the daylight-saving transition rules in
// effect for that year. Zones without daylight saving produce "std offset".
func FormatTimeZone(current time.Time) string {
	stdName, stdOffset, dstName, dstStart, dstEnd, hasDST := zonePeriods(current)
	result := quoteName(stdName) + formatOffset(stdOffset)
	if !hasDST {
		return result
	}
	// Rule 1 is always the start of daylight saving, rule 2 its end, in both
	// hemispheres; only the calendar order of the two differs.
	return fmt.Sprintf("%s%s,%s,%s", result, quoteName(dstName), formatRule(dstStart), formatRule(dstEnd))
}

// GetPosixOffset returns the standard-time UTC offset of current's zone in
// whole hours, counted west of UTC as POSIX does (CET is -1, EST is 5).
// Deprecated: minutes are truncated. Use FormatTimeZone for the full string.
func GetPosixOffset(current time.Time) int {
	_, offset, _, _, _, _ := zonePeriods(current)
	return -(offset / 3600)
}

// zonePeriods locates the standard and daylight periods around current.
// dstStart is the instant standard time ends; dstEnd the instant it resumes.
func zonePeriods(current time.Time) (stdName string, stdOffset int, dstName string, dstStart, dstEnd time.Time, hasDST bool) {
	name, offset := current.Zone()
	start, end := current.ZoneBounds()
	if start.IsZero() && end.IsZero() {
		return name, offset, "", time.Time{}, time.Time{}, false
	}
	if current.IsDST() {
		// Step into the following standard period and describe that instead.
		return zonePeriods(end.Add(time.Hour))
	}
	// A zone that abolished daylight saving keeps its last historical
	// transition as `start` and reports a far-future or zero `end`.
	if end.IsZero() || end.Sub(current) > 366*24*time.Hour {
		return name, offset, "", time.Time{}, time.Time{}, false
	}
	after := end.Add(time.Hour) // some instant inside the daylight period
	dstName, dstOffset := after.Zone()
	if dstName == name && dstOffset == offset {
		return name, offset, "", time.Time{}, time.Time{}, false
	}
	_, dstEnd = after.ZoneBounds()
	return name, offset, dstName, end, dstEnd, true
}

// formatRule renders one transition as Mm.w.d[/time]. POSIX expresses the
// transition in the local wall clock *before* it happens, so the instant is
// shifted by the offset that was in force one moment earlier. The time is
// omitted when it is 02:00, the POSIX default, matching tzdata's own strings.
func formatRule(transition time.Time) string {
	_, offsetBefore := transition.Add(-time.Nanosecond).Zone()
	wall := transition.UTC().Add(time.Duration(offsetBefore) * time.Second)
	month := int(wall.Month())
	weekday := int(wall.Weekday())
	week := (wall.Day()-1)/7 + 1
	if wall.Day()+7 > daysInMonth(wall) {
		week = 5 // POSIX: 5 means "the last such weekday of the month"
	}
	rule := fmt.Sprintf("M%d.%d.%d", month, week, weekday)
	if wall.Hour() == posixDefaultTransitionHour && wall.Minute() == 0 && wall.Second() == 0 {
		return rule
	}
	return rule + "/" + formatClock(wall.Hour()*3600+wall.Minute()*60+wall.Second())
}

// formatOffset renders a UTC offset (seconds east) as POSIX hours west,
// keeping minutes and seconds only when they are non-zero ("-5:30").
func formatOffset(offsetEast int) string {
	sign := ""
	if offsetEast > 0 {
		sign = "-"
	}
	return sign + formatClock(abs(offsetEast))
}

// formatClock renders seconds as h[:mm[:ss]].
func formatClock(seconds int) string {
	h, m, s := seconds/3600, seconds%3600/60, seconds%60
	switch {
	case s != 0:
		return fmt.Sprintf("%d:%02d:%02d", h, m, s)
	case m != 0:
		return fmt.Sprintf("%d:%02d", h, m)
	default:
		return fmt.Sprintf("%d", h)
	}
}

// quoteName wraps zone names that are not purely alphabetic ("+03") in
// angle brackets, as POSIX requires.
func quoteName(name string) string {
	for _, r := range name {
		if !(r >= 'A' && r <= 'Z' || r >= 'a' && r <= 'z') {
			return "<" + strings.Trim(name, "<>") + ">"
		}
	}
	return name
}

func daysInMonth(t time.Time) int {
	return time.Date(t.Year(), t.Month()+1, 0, 0, 0, 0, 0, t.Location()).Day()
}

func abs(v int) int {
	if v < 0 {
		return -v
	}
	return v
}
