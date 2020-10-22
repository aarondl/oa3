package support

import (
	"errors"
	"fmt"
	"reflect"
	"regexp"
	"strings"
)

var (
	rgxUUIDv4 = regexp.MustCompile(`(?i)^[0-9A-F]{8}-[0-9A-F]{4}-4[0-9A-F]{3}-[89AB][0-9A-F]{3}-[0-9A-F]{12}$`)
)

// ValidateUUIDv4 makes it look like it's in a UUID v4 shape
func ValidateUUIDv4(s string) error {
	if !rgxUUIDv4.MatchString(s) {
		return errors.New("must be a valid uuid v4")
	}

	return nil
}

// ValidateMaxInt checks that i <= max or if exclusive then i < max
func ValidateMaxInt(i int64, max int, exclusive bool) error {
	if exclusive {
		if i >= int64(max) {
			return fmt.Errorf("must be less than %d", max)
		}
	} else {
		if i > int64(max) {
			return fmt.Errorf("must be less than or equal to %d", max)
		}
	}

	return nil
}

// ValidateMinInt checks that i >= min or if exclusive then i > min
func ValidateMinInt(i int64, min int, exclusive bool) error {
	if exclusive {
		if i <= int64(min) {
			return fmt.Errorf("must be greater than %d", min)
		}
	} else {
		if i < int64(min) {
			return fmt.Errorf("must be greater than or equal to %d", min)
		}
	}

	return nil
}

// ValidateMultipleOfInt checks that i % factor == 0
func ValidateMultipleOfInt(i int64, factor int) error {
	if i%int64(factor) != 0 {
		return fmt.Errorf("must be a multiple of %d", factor)
	}

	return nil
}

// ValidateMaxUint checks that i <= max or if exclusive then i < max
func ValidateMaxUint(i uint64, max uint, exclusive bool) error {
	if exclusive {
		if i >= uint64(max) {
			return fmt.Errorf("must be less than %d", max)
		}
	} else {
		if i > uint64(max) {
			return fmt.Errorf("must be less than or equal to %d", max)
		}
	}

	return nil
}

// ValidateMinUint checks that i >= min or if exclusive then i > min
func ValidateMinUint(i uint64, min uint, exclusive bool) error {
	if exclusive {
		if i <= uint64(min) {
			return fmt.Errorf("must be greater than %d", min)
		}
	} else {
		if i < uint64(min) {
			return fmt.Errorf("must be greater than or equal to %d", min)
		}
	}

	return nil
}

// ValidateMultipleOfUint checks that i % factor == 0
func ValidateMultipleOfUint(i uint64, factor uint) error {
	if i%uint64(factor) != 0 {
		return fmt.Errorf("must be a multiple of %d", factor)
	}

	return nil
}

// ValidateMaxFloat64 checks that i <= max or if exclusive then i < max
func ValidateMaxFloat64(f float64, max float64, exclusive bool) error {
	if exclusive {
		if f >= float64(max) {
			return fmt.Errorf("must be less than %f", max)
		}
	} else {
		if f > float64(max) {
			return fmt.Errorf("must be less than or equal to %f", max)
		}
	}

	return nil
}

// ValidateMinFloat64 checks that i >= min or if exclusive then i > min
func ValidateMinFloat64(f float64, min float64, exclusive bool) error {
	if exclusive {
		if f <= float64(min) {
			return fmt.Errorf("must be greater than %f", min)
		}
	} else {
		if f < float64(min) {
			return fmt.Errorf("must be greater than or equal to %f", min)
		}
	}

	return nil
}

// ValidateMultipleOfFloat64 checks that i / factor == 0
func ValidateMultipleOfFloat64(i, factor float64) error {
	if i/factor != 0 {
		return fmt.Errorf("must be a multiple of %f", factor)
	}

	return nil
}

// ValidateMaxLength ensures a string's length is <= max
func ValidateMaxLength(s string, max int) error {
	if len(s) <= max {
		return nil
	}

	return fmt.Errorf("length must be less than or equal to %d", max)
}

// ValidateMinLength ensures a string's length is >= min
func ValidateMinLength(s string, min int) error {
	if len(s) >= min {
		return nil
	}

	return fmt.Errorf("length must be greater than or equal to %d", min)
}

// ValidateMaxItems ensures a array's length is <= max
func ValidateMaxItems(a interface{}, max int) error {
	val := reflect.ValueOf(a)
	if val.Len() <= max {
		return nil
	}

	return fmt.Errorf("length must be less than or equal to %d", max)
}

// ValidateMinItems ensures an array's length is >= min
func ValidateMinItems(a interface{}, min int) error {
	val := reflect.ValueOf(a)
	if val.Len() >= min {
		return nil
	}

	return fmt.Errorf("length must be greater than or equal to %d", min)
}

// ValidateUniqueItems ensures an arrays items are unique. Uses
// reflect.DeepEqual and a very naive algorithm, not very performant.
func ValidateUniqueItems(a interface{}) error {
	val := reflect.ValueOf(a)

	ln := val.Len()
	for i := 0; i < ln; i++ {
		a := val.Index(i)
		for j := i + 1; j < ln; j++ {
			b := val.Index(j)

			if reflect.DeepEqual(a.Interface(), b.Interface()) {
				return fmt.Errorf("items must all be unique")
			}
		}
	}

	return nil
}

// ValidateMaxProperties ensures a map[string]X's length is <= max
func ValidateMaxProperties(m interface{}, max int) error {
	val := reflect.ValueOf(m)
	if val.Len() <= max {
		return nil
	}

	return fmt.Errorf("number of properties must be less than or equal to %d", max)
}

// ValidateMinProperties ensures a map[string]X's length is >= min
func ValidateMinProperties(m interface{}, min int) error {
	val := reflect.ValueOf(m)
	if val.Len() >= min {
		return nil
	}

	return fmt.Errorf("number of properties must be greater than or equal to %d", min)
}

// ValidatePattern validates a string against a pattern
func ValidatePattern(s string, pattern string) error {
	matched, err := regexp.MatchString(pattern, s)
	if err != nil {
		panic(err)
	}

	if !matched {
		return fmt.Errorf("must conform to the pattern `%s`", pattern)
	}

	return nil
}

// ValidateEnum validates a string against a whitelisted set of values
func ValidateEnum(s string, whitelist []string) error {
	for _, w := range whitelist {
		if s == w {
			return nil
		}
	}

	return fmt.Errorf(`must be one of: "%s"`, strings.Join(whitelist, `", "`))
}
