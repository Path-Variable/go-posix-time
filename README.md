# go-posix-time

This is a small library whose purpose is to format a Go time.Time struct into a POSIX.1 TZ string.

Since the POSIX timezone string contains information about transitions between DST and Standard time
zones, it was not possible to implement such a format before Go 1.19 when the ZoneBounds struct became
available.

The library was created because the standard time format function of Go cannot produce this format.

For specifications of the format, please refer to UNIX Standard, Base Specifications, Version 7 2018 available on-line at the
[Open Group Library](https://publications.opengroup.org/c181)

If you encounter any bugs, please open issues here.