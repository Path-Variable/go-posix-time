package p_time

import (
	"testing"
	"time"
)

const expectedCET = "CET-1CEST,M3.5.0/2,M10.5.0/3"
const expectedEST = "EST5EDT,M3.2.0/2,M11.1.0"
const zagrebZone = "Europe/Zagreb"
const newyorkZone = "America/New_York"
const errorTemplate = "Time zone string for test time %s is not correct. Expected %s but got %s"

func TestFormatTimeZoneCET(t *testing.T) {
	testFormatTimeZone(t, zagrebZone, expectedCET, 2022, 2, 12)
}

func TestFormatTimeZoneCETDSTFirst(t *testing.T) {
	testFormatTimeZone(t, zagrebZone, expectedCET, 2022, 7, 12)
}

func TestFormatTimeZoneEST(t *testing.T) {
	testFormatTimeZone(t, newyorkZone, expectedEST, 2022, 2, 12)
}

func TestFormatTimeZoneESTDSTFirst(t *testing.T) {
	testFormatTimeZone(t, newyorkZone, expectedEST, 2022, 7, 12)
}

func testFormatTimeZone(t *testing.T, location, expectedFormat string, year, month, day int) {
	loc, _ := time.LoadLocation(location)
	testTime := time.Date(year, time.Month(month), day, 0, 0, 0, 0, loc)
	formatted := FormatTimeZone(testTime)
	if formatted != expectedFormat {
		t.Fatalf(errorTemplate, testTime, expectedFormat, formatted)
	}
}