package support

import (
	"errors"
	"fmt"
	"reflect"
	"regexp"
	"strconv"
	"strings"
	"time"

	"github.com/aarondl/chrono"
	"github.com/google/uuid"
	"github.com/peterhellberg/duration"
	"github.com/shopspring/decimal"
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
func StringToInt[T constraints.Signed](s string) (T, error) {
	var t T
	typ := reflect.TypeOf(t)

	i64, err := strconv.ParseInt(s, 10, typ.Bits())
	if err != nil {
		var zero T
		return zero, err
	}

	return T(i64), nil
}

// StringToUint conversion
func StringToUint[T constraints.Unsigned](s string) (T, error) {
	var t T
	typ := reflect.TypeOf(t)

	u64, err := strconv.ParseUint(s, 10, typ.Bits())
	if err != nil {
		var zero T
		return zero, err
	}

	return T(u64), nil
}

// StringToFloat conversion
func StringToFloat[T constraints.Float](s string) (T, error) {
	var t T
	typ := reflect.TypeOf(t)

	f64, err := strconv.ParseFloat(s, typ.Bits())
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

// StringToDecimal converts a string to a decimal type
func StringToDecimal(s string) (decimal.Decimal, error) {
	return decimal.NewFromString(s)
}

// StringToUUID converts a string to a uuid type
func StringToUUID(s string) (uuid.UUID, error) {
	return uuid.Parse(s)
}

// StringToString is a somewhat useless function but handy for using in
// tandem with a Map() like function over a slice to convert a string to
// a specialized version of a string and back.
//
// The error is to make it easier to pass in to a function that expects a
// conversion may create an error.
//
// In particular this helps eliminate special cases for string where other
// types will require conversion.
func StringToString[A, B ~string](s A) (B, error) {
	return B(s), nil
}

// StringNoOp is used to prevent extra special cases where other types will
// require conversion.
func StringNoOp(s string) (string, error) {
	return s, nil
}

// ExplodedFormArrayToSlice simply takes an array gathered from color=blue&color=black
// and converts it into []T using a conversion function
func ExplodedFormArrayToSlice[T any](formArray []string, convert func(s string) (T, error)) ([]T, error) {
	var err error
	out := make([]T, len(formArray))
	for i, s := range formArray {
		out[i], err = convert(s)
		if err != nil {
			return nil, fmt.Errorf("error converting form value (%d) from []string to %T: %w", i, out, err)
		}
	}

	return out, nil
}

// FlatFormArrayToSlice takes the first value from the form color=blue,black&color=yellow
// (ie. [blue,black]) and converts it into []T using a conversion function.
func FlatFormArrayToSlice[T any](formArray []string, convert func(s string) (T, error)) ([]T, error) {
	var err error
	out := make([]T, len(formArray))
	for i, s := range strings.Split(formArray[0], ",") {
		out[i], err = convert(s)
		if err != nil {
			return nil, fmt.Errorf("error converting form value (%d) from []string to %T: %w", i, out, err)
		}
	}

	return out, nil
}
