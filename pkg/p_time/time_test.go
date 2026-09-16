package p_time

import (
	"testing"
	"time"
)

// Expected strings match the POSIX TZ strings shipped in the IANA tzdata
// footers for these zones (see `zdump -v` or the tzfile(5) footer).
func TestFormatTimeZone(t *testing.T) {
	cases := []struct {
		zone string
		when time.Time
		want string
	}{
		// Northern hemisphere, checked from inside standard time and inside DST.
		{"Europe/Zagreb", date(2022, 2, 12), "CET-1CEST,M3.5.0,M10.5.0/3"},
		{"Europe/Zagreb", date(2022, 7, 12), "CET-1CEST,M3.5.0,M10.5.0/3"},
		{"America/New_York", date(2022, 2, 12), "EST5EDT,M3.2.0,M11.1.0"},
		{"America/New_York", date(2022, 7, 12), "EST5EDT,M3.2.0,M11.1.0"},
		// Transition at 01:00 local, so the hour must be written.
		{"Europe/London", date(2024, 1, 10), "GMT0BST,M3.5.0/1,M10.5.0"},
		// Southern hemisphere: DST starts in October and ends in April.
		{"Australia/Sydney", date(2024, 2, 10), "AEST-10AEDT,M10.1.0,M4.1.0/3"},
		{"Australia/Sydney", date(2024, 6, 10), "AEST-10AEDT,M10.1.0,M4.1.0/3"},
		// Fall-back at midnight: the wall clock before the transition is
		// Sunday 00:00, not Saturday 23:00.
		{"Asia/Beirut", date(2024, 1, 10), "EET-2EEST,M3.5.0/0,M10.5.0/0"},
		// Half-hour offset and no daylight saving.
		{"Asia/Kolkata", date(2024, 1, 10), "IST-5:30"},
		{"Asia/Tokyo", date(2024, 1, 10), "JST-9"},
		// Numeric zone abbreviations must be quoted.
		{"Europe/Istanbul", date(2024, 1, 10), "<+03>-3"},
		{"America/Sao_Paulo", date(2024, 1, 10), "<-03>3"},
		{"UTC", date(2024, 1, 10), "UTC0"},
	}
	for _, tc := range cases {
		t.Run(tc.zone, func(t *testing.T) {
			loc, err := time.LoadLocation(tc.zone)
			if err != nil {
				t.Fatalf("load %s: %v", tc.zone, err)
			}
			got := FormatTimeZone(tc.when.In(loc))
			if got != tc.want {
				t.Errorf("FormatTimeZone(%s @ %s) = %q, want %q", tc.zone, tc.when.Format("2006-01-02"), got, tc.want)
			}
		})
	}
}

func TestGetPosixOffset(t *testing.T) {
	cases := []struct {
		zone string
		when time.Time
		want int
	}{
		{"Europe/Zagreb", date(2022, 2, 12), -1},
		{"Europe/Zagreb", date(2022, 7, 12), -1}, // standard offset even during DST
		{"America/New_York", date(2022, 7, 12), 5},
		{"Australia/Sydney", date(2024, 1, 10), -10},
		{"Asia/Tokyo", date(2024, 1, 10), -9},
	}
	for _, tc := range cases {
		loc, _ := time.LoadLocation(tc.zone)
		if got := GetPosixOffset(tc.when.In(loc)); got != tc.want {
			t.Errorf("GetPosixOffset(%s) = %d, want %d", tc.zone, got, tc.want)
		}
	}
}

func TestFormatClock(t *testing.T) {
	cases := map[int]string{0: "0", 3600: "1", 19800: "5:30", 20700: "5:45", 3661: "1:01:01", 86400: "24"}
	for in, want := range cases {
		if got := formatClock(in); got != want {
			t.Errorf("formatClock(%d) = %q, want %q", in, got, want)
		}
	}
}

func date(y, m, d int) time.Time {
	return time.Date(y, time.Month(m), d, 12, 0, 0, 0, time.UTC)
}
