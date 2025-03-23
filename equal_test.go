package rangeset

import (
	"slices"
	"testing"
)

type EqualElementType int // Element type for tests in this file (it seems we have to declare a different element type in each test file)

// notEqualdata provides a table of data for tests that compare that 2 sets with slight differences
// do not compare equal.  (Tests that equal sets compare equal are found in many other tests).
// It's also used for other test of the elements, some of which depend on the slice of elements
// being ordered without duplicates.
var notEqualData = map[string]struct {
	left, right []EqualElementType // elements 2 sets with some difference(s)
}{
	"NotEqualEmptyVs1":     {[]EqualElementType{}, []EqualElementType{42}},
	"NotEqualEmptyVs2":     {[]EqualElementType{}, []EqualElementType{1, 42}},
	"NotEqualEmptyVs3":     {[]EqualElementType{}, []EqualElementType{1, 2, 42}},
	"NotEqualEmptyVsMany":  {[]EqualElementType{}, []EqualElementType{1, 2, 3, 42, 43, 73, 99}},
	"NotEqualOneElement":   {[]EqualElementType{-1}, []EqualElementType{0}},
	"NotEqualOneVs2Start":  {[]EqualElementType{1}, []EqualElementType{1, 2}},
	"NotEqualOneVs2End":    {[]EqualElementType{2}, []EqualElementType{1, 2}},
	"NotEqualOneVs2Sep":    {[]EqualElementType{1}, []EqualElementType{2, 3, 4}},
	"NotEqualOneVs2Neg":    {[]EqualElementType{1}, []EqualElementType{-3, -2, -1}},
	"NotEqualTwoOverlap1":  {[]EqualElementType{1, 2}, []EqualElementType{0, 1}},
	"NotEqualTwoOverlap2":  {[]EqualElementType{1, 2}, []EqualElementType{2, 3}},
	"NotEqualTwoOverlap3":  {[]EqualElementType{1, 2}, []EqualElementType{0, 1, 2}},
	"NotEqualTwoOverlap4":  {[]EqualElementType{1, 2}, []EqualElementType{1, 2, 3}},
	"NotEqualTwoOverlap5":  {[]EqualElementType{1, 2}, []EqualElementType{0, 1, 2, 3}},
	"NotEqualTwoNoOverlap": {[]EqualElementType{1, 2}, []EqualElementType{3, 4}},
	"NotEqual2Ranges1":     {[]EqualElementType{1, 9}, []EqualElementType{1, 8}},
	"NotEqual2Ranges2":     {[]EqualElementType{1, 9}, []EqualElementType{2, 9}},
	"NotEqual2Ranges3":     {[]EqualElementType{1, 9}, []EqualElementType{1, 2, 9}},
	"NotEqual2Ranges4":     {[]EqualElementType{1, 9}, []EqualElementType{1, 8, 9}},
	"NotEqual2Ranges5":     {[]EqualElementType{1, 9}, []EqualElementType{0, 1, 9}},
	"NotEqual2Ranges6":     {[]EqualElementType{1, 9}, []EqualElementType{1, 9, 10}},
	"NotEqual2Ranges7":     {[]EqualElementType{1, 2, 9}, []EqualElementType{0, 1, 2, 9}},
	"NotEqual2Ranges8":     {[]EqualElementType{1, 2, 9}, []EqualElementType{1, 2, 3, 9}},
	"NotEqual2Ranges9":     {[]EqualElementType{1, 2, 9}, []EqualElementType{0, 1, 2, 3, 9}},
	"NotEqual2Ranges10":    {[]EqualElementType{1, 2, 9}, []EqualElementType{1, 2, 8, 9}},
	"NotEqual2Ranges11":    {[]EqualElementType{1, 2, 9}, []EqualElementType{1, 2, 9, 10}},
	"NotEqual2Ranges12":    {[]EqualElementType{1, 2, 9}, []EqualElementType{1, 2, 8, 9, 10}},
	"NotEqual2Ranges13":    {[]EqualElementType{1, 2, 9}, []EqualElementType{-1, 1, 2, 9}},
	"NotEqual2Ranges14":    {[]EqualElementType{1, 2, 9}, []EqualElementType{1, 2, 4, 9}},
	"NotEqual2Ranges15":    {[]EqualElementType{1, 2, 9}, []EqualElementType{1, 2, 9, 42}},
	"NotEqual2Ranges16":    {[]EqualElementType{1, 8, 9}, []EqualElementType{0, 1, 8, 9}},
	"NotEqual2Ranges17":    {[]EqualElementType{1, 8, 9}, []EqualElementType{1, 2, 8, 9}},
	"NotEqual2Ranges18":    {[]EqualElementType{1, 8, 9}, []EqualElementType{0, 1, 2, 8, 9}},
	"NotEqual2Ranges19":    {[]EqualElementType{1, 8, 9}, []EqualElementType{1, 7, 8, 9}},
	"NotEqual2Ranges20":    {[]EqualElementType{1, 8, 9}, []EqualElementType{1, 8, 9, 10}},
	"NotEqual2Ranges21":    {[]EqualElementType{1, 8, 9}, []EqualElementType{1, 7, 8, 9, 10}},
	"NotEqual2Ranges22":    {[]EqualElementType{1, 2, 8, 9}, []EqualElementType{0, 1, 2, 8, 9}},
	"NotEqual2Ranges23":    {[]EqualElementType{1, 2, 8, 9}, []EqualElementType{1, 2, 3, 8, 9}},
	"NotEqual2Ranges24":    {[]EqualElementType{1, 2, 8, 9}, []EqualElementType{0, 1, 2, 3, 8, 9}},
	"NotEqual2Ranges25":    {[]EqualElementType{1, 2, 8, 9}, []EqualElementType{1, 2, 7, 8, 9}},
	"NotEqual2Ranges26":    {[]EqualElementType{1, 2, 8, 9}, []EqualElementType{1, 2, 8, 9, 10}},
	"NotEqual2Ranges27":    {[]EqualElementType{1, 2, 8, 9}, []EqualElementType{1, 2, 7, 8, 9, 10}},
}

// TestNotEqual tests the Equal() function for cases where it returns false.
// Note that TestRoundTripEqual1(), TestRoundTripEqual2(), etc test cases where Equal() should return true.
func TestNotEqual(t *testing.T) {
	for name, data := range notEqualData {
		set1, set2 := Make(data.left...), Make(data.right...)
		got := Equal(set1, set2)
		Assertf(t, !got, "%20s: Expecting Equal() on different sets (%s and %s) to return false, got %t",
			name, set1.String(), set2.String(), got)
		// Also do reverse comparison
		got = Equal(set2, set1)
		Assertf(t, !got, "%20s: Expecting Equal() on different sets (%s and %s) to return false, got %t",
			name, set2.String(), set2.String(), got)
	}
}

// TestRoundTripLeft uses the first (left) set of the notEqualData table (above) to perform round trip tests, by converting
// to a string then back to a set.  This provides extra tests of Equal() (as well as String() and NewFromString()).
func TestRoundTripLeft(t *testing.T) {
	for name, data := range notEqualData {
		s := Make(data.left...)
		got, _ := NewFromString[EqualElementType](s.String())
		same := Equal(s, got)
		Assertf(t, same, "%20s: After round trip comparing (1st) %q: expected true, got %t", name, s.String(), same)
	}
}

// TestRoundTripRight is like TestRoundTripLetf (above) but uses the 2nd (right) set from the notEqualData table.
func TestRoundTripEqual2(t *testing.T) {
	for name, data := range notEqualData {
		s := Make(data.right...)
		got, _ := NewFromString[EqualElementType](s.String())
		same := Equal(s, got)
		Assertf(t, same, "%20s: After round trip comparing (2nd) %q: expected true, got %t", name, s.String(), same)
	}
}

// TestValuesLeft uses the "left" data to check that slice returned by Values() method is correct.
// Note that this relies on the values in data.left being sorted without duplicates
func TestValuesLeft(t *testing.T) {
	for name, data := range notEqualData {
		s := Make(data.left...)
		got := s.Values()
		Assertf(t, slices.Equal(data.left, got), "Values test %20s: expected result %v, got %v", name, data.left, got)
	}
}

// TestValuesRight uses the "right" data to check that slice returned by Values() method is correct.
// Note that this relies on the values in data.right being sorted without duplicates
func TestValuesRight(t *testing.T) {
	for name, data := range notEqualData {
		s := Make(data.right...)
		got := s.Values()
		Assertf(t, slices.Equal(data.right, got), "Values test %20s: expected result %v, got %v", name, data.right, got)
	}
}
