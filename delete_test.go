package rangeset_test

import (
	"github.com/andrewwphillips/rangeset"
	"testing"
)

type DeleteElementType int

var deleteData = map[string]struct {
	elts      []DeleteElementType
	deleteElt DeleteElementType
	expected  string
}{
	"DelFromEmptySet":      {[]DeleteElementType{}, 42, "{}"},
	"DelBeforeStart":       {[]DeleteElementType{1, 5, 6, 7}, 0, "{1,5:7}"},
	"DelSingleEltRange":    {[]DeleteElementType{1, 5, 6, 7}, 1, "{5:7}"},
	"DelBetweenRanges2":    {[]DeleteElementType{1, 5, 6, 7}, 2, "{1,5:7}"},
	"DelBetweenRanges3":    {[]DeleteElementType{1, 5, 6, 7}, 3, "{1,5:7}"},
	"DelBetweenRanges4":    {[]DeleteElementType{1, 5, 6, 7}, 4, "{1,5:7}"},
	"DelStartRange":        {[]DeleteElementType{1, 5, 6, 7}, 5, "{1,6:7}"},
	"DelWithinRange":       {[]DeleteElementType{1, 5, 6, 7}, 6, "{1,5,7}"},
	"DelEndRange":          {[]DeleteElementType{1, 5, 6, 7}, 7, "{1,5:6}"},
	"DelAfterLast":         {[]DeleteElementType{1, 5, 6, 7}, 8, "{1,5:7}"},
	"DelEndMidRange":       {[]DeleteElementType{1, 5, 6, 7, 10, 11}, 7, "{1,5:6,10:11}"},
	"Del2ndGap":            {[]DeleteElementType{1, 5, 6, 7, 10, 11}, 8, "{1,5:7,10:11}"},
	"DelStartLastRange":    {[]DeleteElementType{1, 5, 6, 7, 10, 11}, 10, "{1,5:7,11}"},
	"DelEndLastRange":      {[]DeleteElementType{1, 5, 6, 7, 10, 11}, 11, "{1,5:7,10}"},
	"DelAfterLastRangeOf3": {[]DeleteElementType{1, 5, 6, 7, 10, 11}, 12, "{1,5:7,10:11}"},
	"DelSingleEltRangeEnd": {[]DeleteElementType{1, 5, 6, 7, 11}, 11, "{1,5:7}"},
}

// TestTableDelete is a table driven test (using deleteData map above) that removes a set element (using
// the Delete() method).  It also tests converting a set to a string (using the String() method()).
func TestTableDelete(t *testing.T) {
	for name, data := range deleteData {
		s := rangeset.Make(data.elts...)
		s.Delete(data.deleteElt)
		got := s.String()
		Assertf(t, got == data.expected, "TableDelete:%24s: expected %q got %q\n", name, data.expected, got)
	}
}

// TestDeleteRealloc tests a special case where append() has to reallocate memory
// We start with a single range and repeatedly split it by deleting elements in the
// (original) range - each time splitting a range in two to expand the slice.
func TestDeleteRealloc(t *testing.T) {
	//s := NewFromRange[int](1, 12)
	// The above causes an error: instantiate୦୦NewFromRange୦int redeclared in this block
	s := rangeset.NewFromRange[DeleteElementType](1, 12)
	// Cause 5 appends - will trigger 2 or 3 memory reallocations
	s.Delete(8)
	s.Delete(2)
	s.Delete(10)
	s.Delete(4)
	s.Delete(6)

	const expected = "{1,3,5,7,9,11}"
	got := s.String()
	Assertf(t, got == expected, "%24s: expected %q got %q\n", "DeleteRealloc", expected, got)
}

var deleteRangeData = map[string]struct {
	elts       []DeleteElementType
	bElt, tElt DeleteElementType
	expected   string
}{
	"DelEmptySet":       {[]DeleteElementType{}, 1, 2, "{}"},
	"DelEmptyEmpty":     {[]DeleteElementType{}, -1, -1, "{}"},
	"DelEmptyBackwards": {[]DeleteElementType{}, 2, 1, "{}"},

	"DelOneRangeExact":      {[]DeleteElementType{1, 2, 3, 4}, 1, 5, "{}"},
	"DelOneRangeOverlap":    {[]DeleteElementType{1, 2, 3, 4}, 0, 99, "{}"},
	"DelOneRangeStart":      {[]DeleteElementType{1, 2, 3, 4}, 1, 4, "{4}"},
	"DelOneRangeStartMore":  {[]DeleteElementType{1, 2, 3, 4}, 1, 6, "{}"},
	"DelOneRangeEnd":        {[]DeleteElementType{1, 2, 3, 4}, 4, 5, "{1:3}"},
	"DelOneRangeEndMore":    {[]DeleteElementType{1, 2, 3, 4}, 0, 5, "{}"},
	"DelOneRangeBothBefore": {[]DeleteElementType{1, 2, 3, 4}, 0, 1, "{1:4}"},
	"DelOneRangeBothAfter":  {[]DeleteElementType{1, 2, 3, 4}, 5, 99, "{1:4}"},
	"DelOneRangeBeforeInto": {[]DeleteElementType{1, 2, 3, 4}, 0, 3, "{3:4}"},
	"DelOneRangeAfterOutOf": {[]DeleteElementType{1, 2, 3, 4}, 3, 99, "{1:2}"},
	"DelOneRangeSplit":      {[]DeleteElementType{1, 2, 3, 4}, 2, 3, "{1,3:4}"},

	"DelTwoRangeExact":      {[]DeleteElementType{1, 2, 3, 4, 7, 8, 9}, 1, 10, "{}"},
	"DelTwoRangeOverlap":    {[]DeleteElementType{1, 2, 3, 4, 7, 8, 9}, 0, 11, "{}"},
	"DelTwoRangeStartFirst": {[]DeleteElementType{1, 2, 3, 4, 7, 8, 9}, 1, 2, "{2:4,7:9}"},
	"DelTwoRangeWholeFirst": {[]DeleteElementType{1, 2, 3, 4, 7, 8, 9}, 1, 5, "{7:9}"},
	"DelTwoRangeAfterFirst": {[]DeleteElementType{1, 2, 3, 4, 7, 8, 9}, 1, 6, "{7:9}"},
	"DelTwoRangeStart2nd":   {[]DeleteElementType{1, 2, 3, 4, 7, 8, 9}, 1, 8, "{8:9}"},
	"DelTwoRangeAfter2nd":   {[]DeleteElementType{1, 2, 3, 4, 7, 8, 9}, 1, 99, "{}"},
	"DelTwoRangeGap":        {[]DeleteElementType{1, 2, 3, 4, 7, 8, 9}, 5, 7, "{1:4,7:9}"},

	"Del3RangeExact":       {[]DeleteElementType{1, 2, 3, 4, 7, 8, 9, 20, 21, 22}, 1, 23, "{}"},
	"Del3RangeOverlap":     {[]DeleteElementType{1, 2, 3, 4, 7, 8, 9, 20, 21, 22}, 0, 99, "{}"},
	"Del3RangeStartExact":  {[]DeleteElementType{1, 2, 3, 4, 7, 8, 9, 20, 21, 22}, 1, 5, "{7:9,20:22}"},
	"Del3RangeStartEnd":    {[]DeleteElementType{1, 2, 3, 4, 7, 8, 9, 20, 21, 22}, 1, 6, "{7:9,20:22}"},
	"Del3RangeStartsStart": {[]DeleteElementType{1, 2, 3, 4, 7, 8, 9, 20, 21, 22}, 0, 5, "{7:9,20:22}"},
	"Del3RangeStartBoth":   {[]DeleteElementType{1, 2, 3, 4, 7, 8, 9, 20, 21, 22}, 0, 6, "{7:9,20:22}"},
	"Del3RangeStartSplit":  {[]DeleteElementType{1, 2, 3, 4, 7, 8, 9, 20, 21, 22}, 2, 4, "{1,4,7:9,20:22}"},
	"Del3RangeMiddleExact": {[]DeleteElementType{1, 2, 3, 4, 7, 8, 9, 20, 21, 22}, 7, 10, "{1:4,20:22}"},
	"Del3RangeMiddleEnd":   {[]DeleteElementType{1, 2, 3, 4, 7, 8, 9, 20, 21, 22}, 7, 12, "{1:4,20:22}"},
	"Del3RangeMiddleStart": {[]DeleteElementType{1, 2, 3, 4, 7, 8, 9, 20, 21, 22}, 5, 10, "{1:4,20:22}"},
	"Del3RangeMiddleBoth":  {[]DeleteElementType{1, 2, 3, 4, 7, 8, 9, 20, 21, 22}, 6, 11, "{1:4,20:22}"},
	"Del3RangeMiddleSplit": {[]DeleteElementType{1, 2, 3, 4, 7, 8, 9, 20, 21, 22}, 8, 9, "{1:4,7,9,20:22}"},
	"Del3RangeEndExact":    {[]DeleteElementType{1, 2, 3, 4, 7, 8, 9, 20, 21, 22}, 20, 23, "{1:4,7:9}"},
	"Del3RangeEndStart":    {[]DeleteElementType{1, 2, 3, 4, 7, 8, 9, 20, 21, 22}, 10, 23, "{1:4,7:9}"},
	"Del3RangeEndEnd":      {[]DeleteElementType{1, 2, 3, 4, 7, 8, 9, 20, 21, 22}, 20, 99, "{1:4,7:9}"},
	"Del3RangeEndBoth":     {[]DeleteElementType{1, 2, 3, 4, 7, 8, 9, 20, 21, 22}, 19, 24, "{1:4,7:9}"},
	"Del3RangeEndSplit":    {[]DeleteElementType{1, 2, 3, 4, 7, 8, 9, 20, 21, 22}, 21, 22, "{1:4,7:9,20,22}"},
	"Del3RangeEndDelStart": {[]DeleteElementType{1, 2, 3, 4, 7, 8, 9, 20, 21, 22}, 20, 22, "{1:4,7:9,22}"},
	"Del3RangeEndDelEnd":   {[]DeleteElementType{1, 2, 3, 4, 7, 8, 9, 20, 21, 22}, 21, 23, "{1:4,7:9,20}"},
}

// TestTableDeleteRange is a table driven test that removes a range of elements (using the
// the DeleteRange() method).
func TestTableDeleteRange(t *testing.T) {
	for name, data := range deleteRangeData {
		s := rangeset.Make(data.elts...)
		s.DeleteRange(data.bElt, data.tElt)
		got := s.String()
		Assertf(t, got == data.expected, "Delete Range:%24s: expected %q got %q\n", name, data.expected, got)
	}
}

var stringDeleteData = map[string]struct {
	in       string
	toDelete uint16
	expected string
}{
	"DeleteStartFromU":  {"{U}", 0, "{1:65535}"},
	"DeleteOneFromU":    {"{U}", 1, "{0,2:65535}"},
	"DeleteTwoFromU":    {"{U}", 2, "{0:1,3:65535}"},
	"DeletePenultimate": {"{U}", 65534, "{0:65533,65535}"},
	"DeleteEndFromU":    {"{U}", 65535, "{0:65534}"},
}

// TestTableDeleteFromU does table driven tests deleting elements from Universal set
func TestTableDeleteFromU(t *testing.T) {
	for name, data := range stringDeleteData {
		s, err := rangeset.NewFromString[uint16](data.in)
		Assertf(t, err == nil, "DeleteFromU:%24s: expected <nil> (no error) got %v\n", name, err)
		s.Delete(data.toDelete)
		got := s.String()
		Assertf(t, got == data.expected, "DeleteFromU:%24s: expected %q got %q\n", name, data.expected, got)
	}
}
