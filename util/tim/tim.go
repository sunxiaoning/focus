package timutil

import (
	"time"
)

const (
	DefaultTimeFormat = "2006-01-02 15:04:05"
	MinTimStr         = "0000-01-01 00:00:00"
)

var ZERO = func() time.Time {
	zero, _ := time.ParseInLocation(DefaultTimeFormat, MinTimStr, time.Local)
	return zero
}()

func Format(t time.Time, f string) string {
	if t.Before(ZERO) || t.Equal(ZERO) {
		return ""
	} else {
		return t.Format(f)
	}
}

func Parse(s string, f string) (time.Time, error) {
	t, err := time.ParseInLocation(DefaultTimeFormat, s, time.Local)
	if err != nil {
		return time.Time{}, err
	}
	if t.Before(ZERO) || t.Equal(ZERO) {
		return ZERO, nil
	} else {
		return time.ParseInLocation(f, s, time.Local)
	}
}

func DefFormat(t time.Time) string {
	return Format(t, DefaultTimeFormat)
}

func DefParse(s string) (time.Time, error) {
	return Parse(s, DefaultTimeFormat)
}
