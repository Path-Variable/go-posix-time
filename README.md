# go-posix-time

[![Go](https://github.com/isaric/go-posix-time/actions/workflows/go.yml/badge.svg?branch=main)](https://github.com/isaric/go-posix-time/actions/workflows/go.yml)


This is a small library whose purpose is to format a Go time.Time struct into a POSIX.1 TZ string.

Since the POSIX timezone string contains information about transitions between DST and Standard time
zones, it was not possible to implement such a format before Go 1.19 when the ZoneBounds struct became
available.

The library was created because the standard time format function of Go cannot produce this format.

For specifications of the format, please refer to UNIX Standard, Base Specifications, Issue 7 2018 available on-line at the
[Open Group Library](https://publications.opengroup.org/c181)

If you encounter any bugs, please open issues [here](https://github.com/isaric/go-posix-time/issues).

## Usage

```go
import "github.com/isaric/go-posix-time/pkg/p_time"

loc, _ := time.LoadLocation("Europe/Madrid")
p_time.FormatTimeZone(time.Now().In(loc)) // "CET-1CEST,M3.5.0,M10.5.0/3"
p_time.FormatTimeZone(time.Now().In(sydney)) // "AEST-10AEDT,M10.1.0,M4.1.0/3"
p_time.FormatTimeZone(time.Now().In(kolkata)) // "IST-5:30"
```

## Behaviour

* Transition rules are written in the wall clock in force just before the
  transition, as POSIX requires, so zones that fall back at midnight
  (`EET-2EEST,M3.5.0/0,M10.5.0/0`) come out right.
* A transition time of 02:00 is omitted, matching the strings in the IANA
  tzdata footers; any other time is written as `/h[:mm[:ss]]`.
* Offsets keep minutes and seconds when present (`<+0545>-5:45`), and
  abbreviations that are not purely alphabetic are quoted in angle brackets.
* Zones that have abolished daylight saving produce a plain `std offset`.
* Rules are derived from the transitions of the year that contains the given
  time, so a zone whose rules change from year to year is only right for that
  year.
