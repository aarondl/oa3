package support

import (
	"strconv"

	"golang.org/x/exp/constraints"
)

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
