package rangeset  // Note: we can't do external tests using package rangeset_test (yet) due to deficiencies of go2go

import (
	"testing"
)

type elementType int  // Note that we have to declare a different element type in each test file due to deficiencies of go2go

// Assertf is a test helper that marks a test as having failed and displays
// more information about the test then continues the current test.
// When called multiple times you will see a list of test messages each preceded
// by a tick or cross, but the list is only displayed if at least one
// test fails or the verbose (-test.v) option is used.
func Assertf(t *testing.T, succeeded bool, format string, args ...interface{}) {
	const (
		succeed = "\u2713" // tick
		failed  = "X"      //"\u2717" // cross
	)

	t.Helper()
	if !succeeded {
		t.Errorf("%s\t" + format, append([]interface{}{failed}, args...)...)
	} else {
		t.Logf("%s\t" + format, append([]interface{}{succeed}, args...)...)
	}
}

// TestZero checks that an empty set has no elements
func TestZero(t *testing.T) {
	var s Set[elementType]
	Assertf(t, s.Len() == 0, "Empty set had size %d", s.Len())
}

// TestOne checks that a set with a single element has the right length
func TestOne(t *testing.T) {
	s := Make(42)
	Assertf(t, s.Len() == 1, "Set with one element had size %d", s.Len())
}

// TestNewRangeOne tests the NewFromRange method with a single element
func TestNewRangeOne(t *testing.T) {
	s := NewFromRange(1, 2)
	Assertf(t, s.Len() == 1, "After calling NewFromRange adding a single element {1} got length %d", s.Len())
	Assertf(t, !s.Contains(0), "After calling NewFromRange the set {1} should *not* contain 0")
	Assertf(t, s.Contains(1), "After calling NewFromRange the set {1} should contain 1")
	Assertf(t, !s.Contains(2), "After calling NewFromRange the set {1} should *not* contain 2")
}

// TestNewRangeMany tests the NewFromRange method with a many elements
func TestNewRangeMany(t *testing.T) {
	s := NewFromRange(-1, 2)
	Assertf(t, s.Len() == 3, "After calling NewFromRange adding {-1,0,1} got length %d (expected 3)", s.Len())
	Assertf(t, !s.Contains(-2), "After calling NewFromRange the set {-1,0,1} should *not* contain -2")
	Assertf(t, s.Contains(-1), "After calling NewFromRange the set {-1,0,1} should contain -1")
	Assertf(t, s.Contains(0), "After calling NewFromRange the set {-1,0,1} should contain 0")
	Assertf(t, s.Contains(1), "After calling NewFromRange the set {-1,0,1} should contain 1")
	Assertf(t, !s.Contains(2), "After calling NewFromRange the set {-1,0,1} should *not* contain 2")
}

// TestLen checks that Len returns correct values
// Note that TestTableNew() also tests Len() [using many more variations]
func TestLen(t *testing.T) {
	const b1, t1, b2, t2 = 1, 10, 11, 20
	const expected = t1 - b1 + t2 - b2
	var s Set[elementType]
	s.AddRange(b1, t1)
	s.AddRange(b2, t2)
	Assertf(t, s.Len() == expected, "TestLen: after adding %d elements to empty set got len %d", expected, s.Len())
	length, spans := s.Length()
	Assertf(t, length == expected, "TestLen: expecting length of %d elements, got length %d", expected, length)
	Assertf(t, spans == 2, "TestLen: expecting span count of 2, got %v", spans)
}

// TestUniversal checks that universal on uint64 has zero length and one range
func TestUniversal(t *testing.T) {
	u := Universal[uint64]()
	size, spans := u.Length()
	Assertf(t, size == 0, "For a universal set was expecting zero length, got %d", size)
	Assertf(t, spans == 1, "For a universal set was expecting one span, got %v", spans)
}

// TestEmptyAndUniversalComplement tests universal and empty set complements
func TestEmptyAndUniversalComplement(t *testing.T) {
	empty := Make[elementType]()
	complement := Complement(empty)
	u := Universal[elementType]()
	Assertf(t, Equal(u, complement),  "%24s: expected complement of empty set to be %v got %v\n",
		"TestEmptyUniversal", u, complement)

	complement = Complement(u)
	Assertf(t, Equal(empty, complement),  "%24s: expected complement of U to be %v got %v\n",
		"TestEmptyUniversal", empty, complement)
}
