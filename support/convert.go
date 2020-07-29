package support

import (
	"strconv"

	"github.com/volatiletech/null/v8"
)

// StringToNullstring conversion
func StringToNullstring(s string) (null.String, error) {
	return null.StringFrom(s), nil
}

// StringToInt conversion
func StringToInt(s string) (int, error) {
	return strconv.Atoi(s)
}

// StringToInt32 conversion
func StringToInt32(s string) (int32, error) {
	i, err := strconv.ParseInt(s, 10, 32)
	if err != nil {
		return 0, err
	}
	return int32(i), nil
}

// StringToInt64 conversion
func StringToInt64(s string) (int64, error) {
	i, err := strconv.ParseInt(s, 10, 64)
	if err != nil {
		return 0, err
	}
	return i, nil
}

// StringToFloat32 conversion
func StringToFloat32(s string) (float32, error) {
	i, err := strconv.ParseFloat(s, 32)
	if err != nil {
		return 0, err
	}
	return float32(i), nil
}

// StringToFloat64 conversion
func StringToFloat64(s string) (float64, error) {
	i, err := strconv.ParseFloat(s, 64)
	if err != nil {
		return 0, err
	}
	return i, nil
}

// StringToBool conversion
func StringToBool(s string) (bool, error) {
	return strconv.ParseBool(s)
}

// StringToNullbool conversion
func StringToNullbool(s string) (null.Bool, error) {
	if len(s) == 0 {
		return null.Bool{}, nil
	}
	b, err := strconv.ParseBool(s)
	if err != nil {
		return null.Bool{}, err
	}
	return null.BoolFrom(b), nil
}
