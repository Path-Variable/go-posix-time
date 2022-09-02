package p_time

import (
	"testing"
	"time"
)

const expectedCET = "CET-1CEST,M3.5.0/2,M10.5.0/3"
const expectedEST = "EST5EDT,M3.2.0/2,M11.1.0"
const zagrebZone = "Europe/Zagreb"
const newyorkZone = "America/New_York"

func TestFormatTimeZoneCET(t *testing.T) {
	loc, _ := time.LoadLocation(zagrebZone)
	testTime := time.Date(2022, 4, 12, 0, 0, 0, 0, loc)
	formatted := FormatTimeZone(testTime)
	if formatted != expectedCET {
		t.Fatalf("Time zone string for test time %s is not correct. Expected %s but got %s", testTime, expectedCET, formatted)
	}
}

func TestFormatTimeZoneCETDSTFirst(t *testing.T) {
	loc, _ := time.LoadLocation(zagrebZone)
	testTime := time.Date(2022, 7, 12, 0, 0, 0, 0, loc)
	formatted := FormatTimeZone(testTime)
	if formatted != expectedCET {
		t.Fatalf("Time zone string for test time %s is not correct. Expected %s but got %s", testTime, expectedCET, formatted)
	}
}

func TestFormatTimeZoneEST(t *testing.T) {
	loc, _ := time.LoadLocation(newyorkZone)
	testTime := time.Date(2022, 4, 12, 0, 0, 0, 0, loc)
	formatted := FormatTimeZone(testTime)
	if formatted != expectedEST {
		t.Fatalf("Time zone string for test time %s is not correct. Expected %s but got %s", testTime, expectedEST, formatted)
	}
}

func TestFormatTimeZoneESTDSTFirst(t *testing.T) {
	loc, _ := time.LoadLocation(newyorkZone)
	testTime := time.Date(2022, 7, 12, 0, 0, 0, 0, loc)
	formatted := FormatTimeZone(testTime)
	if formatted != expectedEST {
		t.Fatalf("Time zone string for test time %s is not correct. Expected %s but got %s", testTime, expectedEST, formatted)
	}
}
