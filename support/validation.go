package support

import (
	"errors"
	"fmt"
	"math"
	"reflect"
	"regexp"
	"strings"

	"golang.org/x/exp/constraints"
)

var (
	rgxUUIDv4 = regexp.MustCompile(`(?i)^[0-9A-F]{8}-[0-9A-F]{4}-4[0-9A-F]{3}-[89AB][0-9A-F]{3}-[0-9A-F]{12}$`)
)

// ValidateFormatUUIDv4 makes it look like it's in a UUID v4 shape
func ValidateFormatUUIDv4(s string) error {
	if !rgxUUIDv4.MatchString(s) {
		return errors.New("must be a valid uuid v4")
	}

	return nil
}

// ValidateMaxNumber checks that val <= max or if exclusive then val < max
func ValidateMaxNumber[N constraints.Integer | constraints.Float](val, max N, exclusive bool) error {
	if exclusive {
		if val >= max {
			return fmt.Errorf("must be less than %v", max)
		}
	} else {
		if val > max {
			return fmt.Errorf("must be less than or equal to %v", max)
		}
	}

	return nil
}

// ValidateMinNumber checks that val >= min or if exclusive then val > min
func ValidateMinNumber[N constraints.Integer | constraints.Float](val, min N, exclusive bool) error {
	if exclusive {
		if val <= min {
			return fmt.Errorf("must be greater than %v", min)
		}
	} else {
		if val < min {
			return fmt.Errorf("must be greater than or equal to %v", min)
		}
	}

	return nil
}

// ValidateMultipleOfInt checks that val % factor == 0
func ValidateMultipleOfInt[N constraints.Integer](val, factor N) error {
	if (val % factor) != 0 {
		return fmt.Errorf("must be a multiple of %d", factor)
	}

	return nil
}

// ValidateMultipleOfInt checks that val % factor == 0
func ValidateMultipleOfFloat[N constraints.Float](val, factor N) error {
	if quot := float64(val / factor); math.Trunc(quot) != quot {
		return fmt.Errorf("must be a multiple of %f", factor)
	}

	return nil
}

// ValidateMaxLength ensures a string's length is <= max
func ValidateMaxLength[S ~string](s S, max int) error {
	if len(s) <= max {
		return nil
	}

	return fmt.Errorf("length must be less than or equal to %d", max)
}

// ValidateMinLength ensures a string's length is >= min
func ValidateMinLength[S ~string](s S, min int) error {
	if len(s) >= min {
		return nil
	}

	return fmt.Errorf("length must be greater than or equal to %d", min)
}

// ValidateMaxItems ensures a array's length is <= max
func ValidateMaxItems[T any](a []T, max int) error {
	if len(a) <= max {
		return nil
	}

	return fmt.Errorf("length must be less than or equal to %d", max)
}

// ValidateMinItems ensures an array's length is >= min
func ValidateMinItems[T any](a []T, min int) error {
	if len(a) >= min {
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
func ValidateMaxProperties[V any, M ~map[string]V](m M, max int) error {
	if len(m) <= max {
		return nil
	}

	return fmt.Errorf("number of properties must be less than or equal to %d", max)
}

// ValidateMinProperties ensures a map[string]X's length is >= min
func ValidateMinProperties[V any, M ~map[string]V](m M, min int) error {
	if len(m) >= min {
		return nil
	}

	return fmt.Errorf("number of properties must be greater than or equal to %d", min)
}

// ValidatePattern validates a string against a pattern
func ValidatePattern[S ~string](s S, pattern string) error {
	matched, err := regexp.MatchString(pattern, string(s))
	if err != nil {
		panic(err)
	}

	if !matched {
		return fmt.Errorf("must conform to the pattern `%s`", pattern)
	}

	return nil
}

// ValidateEnum validates a string against a whitelisted set of values
func ValidateEnum[S ~string](s S, whitelist []string) error {
	for _, w := range whitelist {
		if string(s) == w {
			return nil
		}
	}

	return fmt.Errorf(`must be one of: "%s"`, strings.Join(whitelist, `", "`))
}
