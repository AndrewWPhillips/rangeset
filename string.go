package rangeset

// string.go implements functions to convert a set to/from a string.

import (
	"fmt"
	"strings"
)

// NewFromString deserializes a string (eg from format created by the String() method below).
// If the format of the string is invalid it returns an error.
// Also note that it will accept more variations in strings than generated by the String()
// method (below) (eg "{1,1}" which has a redundant element) but must always consist of zero
// or more comma-separated integers or integer ranges (where a range is two integers separated by
// a colon), without whitespace, and all enclosed in braces.  Ranges are "inclusive" (unlike eg the
// AddRange() method) so, for example, "{1:2}" creates a set with 2 elements stored as the range [1,3).
//
// As a special case it will accept the special symbol E in a range indicating the "end" of valid
// element values - ie the smallest element at the bottom of a range or the largest element at the
// top. Eg: "{E:10}" represents all possible elements up to (and including) 10; "{1:E}" means all
// elements from 1 upwards.  The universal set is represented by the string "{U}" or "{E:E}".
func NewFromString[T Element](s string) (_ Set[T], _ error) {
	s = strings.TrimSpace(s)
	if len(s) < 2 || s[0] != '{' || s[len(s)-1] != '}' {
		return nil, fmt.Errorf("rangeset: string %q for new set is not enclosed in braces", s)
	}
	s = s[1 : len(s)-1]
	if len(s) == 0 { // we need to handle this specially as strings.Split("") returns a slice with one element
		return nil, nil // return empty set
	}
	if s == "U" {
		s = "E:E" // U (universal set abbrev.) means all elts (end to end)
	}

	endMark := minInt[T]() // indicates top/bottom of range of valid elements
	ranges := strings.Split(s, ",")
	retval := make(Set[T], 0, len(ranges))

	for _, r := range ranges {
		var b, t T
		var err error
		if strings.ContainsRune(r, ':') {
			parts := strings.Split(r, ":")
			if len(parts) != 2 {
				//assert(len(parts) > 2)
				return nil, fmt.Errorf("rangeset: too many parts in range %q for set {%s}", r, s)
			}
			if parts[0] == "E" {
				b = endMark
			} else {
				b, err = parseInt[T](parts[0])
				if err != nil {
					return nil, fmt.Errorf("rangeset: %w for an integer at start of range for set {%s}", err, s)
				}
			}
			if parts[1] == "E" {
				t = endMark - 1 // wraps around to max element
			} else {
				t, err = parseInt[T](parts[1])
				if err != nil {
					return nil, fmt.Errorf("rangeset: %w for an integer at end of range for set {%s}", err, s)
				}
			}
		} else {
			b, err = parseInt[T](r)
			if err != nil {
				return nil, fmt.Errorf("rangeset: %w for an integer in set {%s}", err, s)
			}
			t = b
		}
		if b > t {
			return nil, fmt.Errorf("rangeset: invalid range %q (end < start) for set {%s}", r, s)
		}
		// Note that this relies on integer overflow wrapping around if t is
		// the maximum value for T (maxInt[T]()), whence t+1 wraps around to minInt[T]
		retval.AddRange(b, t+1)
	}
	return retval, nil
}

// String generates a string representation of (ie "serialises) a set.
// Such a string can be "deserialised" using the above NewFromString() function.
// TODO: add options to encode min/max values as "E" and universal set as U, add option to change sep. char (:)
func (s Set[T]) String() string {
	var retval strings.Builder
	retval.Grow(5 * len(s)) // est. of generated string length TODO: better estimate
	retval.WriteRune('{')
	first := true
	for _, r := range s {
		if !first {
			retval.WriteRune(',')
		} else {
			first = false
		}
		retval.WriteString(intToString(r.b))
		if r.t != r.b+1 {
			retval.WriteRune(':')
			retval.WriteString(intToString(r.t - 1))
		}
	}
	retval.WriteRune('}')

	return retval.String()
}