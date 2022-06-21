package support

import (
	"errors"
	"regexp"
	"strconv"
	"time"

	"github.com/aarondl/chrono"
	"github.com/peterhellberg/duration"
	"golang.org/x/exp/constraints"
)

var (
	rgxDateTime = regexp.MustCompile(`^\d{4}-\d{2}-\d{2}T\d{2}:\d{2}:\d{2}(?:\.\d+)?(Z|[\+-]\d{2}:\d{2})?`)
	rgxTime     = regexp.MustCompile(`^\d{2}:\d{2}:\d{2}(?:\.\d+)?(Z|[\+-]\d{2}:\d{2})?`)
	rgxDate     = regexp.MustCompile(`^\d{4}-\d{2}-\d{2}$`)
)

// StringToChronoDateTime parses an RFC3339 "date-time" production using chrono
// library
func StringToChronoDateTime(s string) (chrono.DateTime, error) {
	if !rgxDateTime.MatchString(s) {
		return chrono.DateTime{}, errors.New("must be a valid date-time")
	}

	return chrono.DateTimeFromString(s)
}

// StringToChronoDate checks for an RFC3339 "date" production using chrono
// library
func StringToChronoDate(s string) (chrono.Date, error) {
	if !rgxDate.MatchString(s) {
		return chrono.Date{}, errors.New("must be a valid date")
	}

	return chrono.DateFromString(s)
}

// StringToChronoTime parses an RFC3339 "time" production using chrono library
func StringToChronoTime(s string) (chrono.Time, error) {
	if !rgxTime.MatchString(s) {
		return chrono.Time{}, errors.New("must be a valid time")
	}

	return chrono.TimeFromString(s)
}

// StringToDateTime parses an RFC3339 "date-time" production as time.Time
func StringToDateTime(s string) (time.Time, error) {
	if !rgxDateTime.MatchString(s) {
		return time.Time{}, errors.New("must be a valid date-time")
	}

	return time.Parse(time.RFC3339, s)
}

// StringToDate checks for an RFC3339 "date" production as time.Time
func StringToDate(s string) (time.Time, error) {
	if !rgxDate.MatchString(s) {
		return time.Time{}, errors.New("must be a valid date")
	}

	return time.Parse(`2006-01-02`, s)
}

// StringToTime parses an RFC3339 "time" production as time.Time
func StringToTime(s string) (time.Time, error) {
	if !rgxTime.MatchString(s) {
		return time.Time{}, errors.New("must be a valid time")
	}

	return time.Parse(`15:04:05Z07:00`, s)
}

// StringToDuration parses an RFC3339 "duration" production
func StringToDuration(s string) (time.Duration, error) {
	dur, err := duration.Parse(s)
	if err != nil {
		return 0, errors.New("must be a valid duration")
	}

	return dur, nil
}

// StringToInt conversion
func StringToInt[T constraints.Signed](s string, bits int) (T, error) {
	i64, err := strconv.ParseInt(s, 10, bits)
	if err != nil {
		var zero T
		return zero, err
	}

	return T(i64), nil
}

// StringToUint conversion
func StringToUint[T constraints.Unsigned](s string, bits int) (T, error) {
	u64, err := strconv.ParseUint(s, 10, bits)
	if err != nil {
		var zero T
		return zero, err
	}

	return T(u64), nil
}

// StringToFloat conversion
func StringToFloat[T constraints.Float](s string, bits int) (T, error) {
	f64, err := strconv.ParseFloat(s, bits)
	if err != nil {
		var zero T
		return zero, err
	}

	return T(f64), nil
}

// StringToBool conversion
func StringToBool(s string) (bool, error) {
	return strconv.ParseBool(s)
}
