package rangeset

// integer.go provides internal functions to deal with the element type, such as
// string conversions and getting the min/max allowed values for an element type

import (
	"math"
	"strconv"
)

// isUnsigned checks whether its type param (must be integer type) is unsigned
// TODO: make sure this function is inlined (otherwise replace it where used)
func isUnsigned[T Element]() bool {
	var zero T
	return zero-1 > 0
}

// parseInt converts a decimal string to the type of it's type parameter
// The return value is undefined if the value overflows the size of the type parameter
// TODO: pass bitSize to ParseInt/Uint depending on type instead of just using 64 (to get error on overflow)
func parseInt[T Element](s string) (T, error) {
	if isUnsigned[T]() {
		// parse unsigned int
		v, err := strconv.ParseUint(s, 10, 64)
		return T(v), err
	}
	v, err := strconv.ParseInt(s, 10, 64)
	return T(v), err
}

// intToString converts an integer to a decimal string
func intToString[T Element](i T) string {
	if isUnsigned[T]() {
		return strconv.FormatUint(uint64(i), 10)
	}
	return strconv.FormatInt(int64(i), 10)
}

// minInt returns the smallest allowed integer for its type param. (signed/unsigned integer)
// Thanks to Robert Greisemer for this (see https://github.com/golang/go/issues/40301)
func minInt[T Element]() T {
	if isUnsigned[T]() {
		return 0 // unsigned int types all start with zero
	}
	// signed int types start with all but top (sign) bit off (2's complement)
	var m int64 = math.MinInt64
	//for s := 32; T(m) == 0; s >>= 1 {
	//	m >>= s
	//}
	// Note unrolling the above loop (as below) made this 10 to 20 times faster
	if T(m) != 0 {
		return T(m)
	}
	m >>= 32
	if T(m) != 0 {
		return T(m)
	}
	m >>= 16
	if T(m) != 0 {
		return T(m)
	}
	m >>= 8
	//assert(T(m) != 0)
	return T(m)
}

// maxInt returns the largest allowed integer for Element type - a signed/unsigned integer
func maxInt[T Element]() T {
	return minInt[T]() - 1 // This works for unsigned and signed (2's complement) ints
}
