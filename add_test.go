package rangeset_test

import (
	"github.com/andrewwphillips/rangeset"
	"testing"
)

type (
	AddType int // set element type used for "add" tests
)

const (
	// maxAddtype is the maximum possible value of AddType.
	// Note that we use uint because AddType (above) is int (uint and int have same number of bits)
	maxAddType = AddType(^uint(0) >> 1)
	// minAddtype is the minimum possible value of AddType.
	// Note that this calculation assumes 2's-complement arithmetic but Go integers always use that
	minAddType = -maxAddType - 1
)

var addData = map[string]struct {
	elts      []AddType                // initial elements in the set
	addElt    AddType                  // element to be added
	expResult bool                     // whether it should have been added or not (already existed)
	expected  []rangeset.Span[AddType] // expected set of spans after elt was added
}{
	"AddEmpty":        {[]AddType{}, 42, true, []rangeset.Span[AddType]{{42, 43}}},
	"AddNegative":     {[]AddType{}, -1, true, []rangeset.Span[AddType]{{-1, 0}}},
	"Add1Before":      {[]AddType{42}, 40, true, []rangeset.Span[AddType]{{40, 41}, {42, 43}}},
	"AddAtStartRange": {[]AddType{42}, 41, true, []rangeset.Span[AddType]{{41, 43}}},
	"AddExisting":     {[]AddType{42}, 42, false, []rangeset.Span[AddType]{{42, 43}}},
	"AddMin":          {[]AddType{}, minAddType, true, []rangeset.Span[AddType]{{minAddType, minAddType + 1}}},
	"AddMax": {
		[]AddType{}, maxAddType, true, []rangeset.Span[AddType]{{maxAddType, minAddType}},
	}, // wraps around
	"AddExistingMin": {
		[]AddType{minAddType}, minAddType, false,
		[]rangeset.Span[AddType]{{minAddType, minAddType + 1}},
	},
	"AddExistingMax": {
		[]AddType{maxAddType}, maxAddType, false, []rangeset.Span[AddType]{{maxAddType, minAddType}},
	}, // wraps around
	"AddMinToMax": {
		[]AddType{maxAddType}, minAddType, true,
		[]rangeset.Span[AddType]{{minAddType, minAddType + 1}, {maxAddType, minAddType}},
	}, // wraps around
	"AddMaxToMin": {
		[]AddType{minAddType}, maxAddType, true,
		[]rangeset.Span[AddType]{{minAddType, minAddType + 1}, {maxAddType, minAddType}},
	}, // wraps around
	"AddAtEndRange": {[]AddType{42}, 43, true, []rangeset.Span[AddType]{{42, 44}}},
	"AddAfter":      {[]AddType{42}, 44, true, []rangeset.Span[AddType]{{42, 43}, {44, 45}}},

	"AddTo2BeforeFirst": {[]AddType{1, 3}, -1, true, []rangeset.Span[AddType]{{-1, 0}, {1, 2}, {3, 4}}},
	"AddTo2StartFirst":  {[]AddType{1, 3}, 0, true, []rangeset.Span[AddType]{{0, 2}, {3, 4}}},
	"AddTo2InFirst":     {[]AddType{1, 3}, 1, false, []rangeset.Span[AddType]{{1, 2}, {3, 4}}},
	"AddTo2Join1And2":   {[]AddType{1, 3}, 2, true, []rangeset.Span[AddType]{{1, 4}}},
	"AddTo2In2nd":       {[]AddType{1, 3}, 3, false, []rangeset.Span[AddType]{{1, 2}, {3, 4}}},
	"AddTo2EndLast":     {[]AddType{1, 3}, 4, true, []rangeset.Span[AddType]{{1, 2}, {3, 5}}},
	"AddTo2AfterLast":   {[]AddType{1, 3}, 5, true, []rangeset.Span[AddType]{{1, 2}, {3, 4}, {5, 6}}},

	"AddToR2ExtendFirst":     {[]AddType{1, 10, 11}, 2, true, []rangeset.Span[AddType]{{1, 3}, {10, 12}}},
	"AddToR2NewRangeBetween": {[]AddType{1, 10, 11}, 3, true, []rangeset.Span[AddType]{{1, 2}, {3, 4}, {10, 12}}},
	"AddToR2ExtendBack2nd":   {[]AddType{1, 10, 11}, 9, true, []rangeset.Span[AddType]{{1, 2}, {9, 12}}},
	"AddToR2Extend2nd":       {[]AddType{1, 10, 11}, 12, true, []rangeset.Span[AddType]{{1, 2}, {10, 13}}},
	"AddToR2AfterLast": {
		[]AddType{1, 10, 11}, 1e9, true,
		[]rangeset.Span[AddType]{{1, 2}, {10, 12}, {1000000000, 1000000001}},
	},

	"AddTo3BeforeFirst": {
		[]AddType{11, 12, 101, 1001, 1002}, 1, true,
		[]rangeset.Span[AddType]{{1, 2}, {11, 13}, {101, 102}, {1001, 1003}},
	},
	"AddTo3StartFirst": {
		[]AddType{11, 12, 101, 1001, 1002}, 10, true,
		[]rangeset.Span[AddType]{{10, 13}, {101, 102}, {1001, 1003}},
	},
	"AddTo3Before2nd": {
		[]AddType{11, 12, 101, 1001, 1002}, 20, true,
		[]rangeset.Span[AddType]{{11, 13}, {20, 21}, {101, 102}, {1001, 1003}},
	},
	"AddTo3Start2nd": {
		[]AddType{11, 12, 101, 1001, 1002}, 100, true,
		[]rangeset.Span[AddType]{{11, 13}, {100, 102}, {1001, 1003}},
	},
	"AddTo3In2nd": {
		[]AddType{11, 12, 101, 1001, 1002}, 101, false,
		[]rangeset.Span[AddType]{{11, 13}, {101, 102}, {1001, 1003}},
	},
	"AddTo3End2nd": {
		[]AddType{11, 12, 101, 1001, 1002}, 102, true,
		[]rangeset.Span[AddType]{{11, 13}, {101, 103}, {1001, 1003}},
	},
	"AddTo3Before3rd": {
		[]AddType{11, 12, 101, 1001, 1002}, 103, true,
		[]rangeset.Span[AddType]{{11, 13}, {101, 102}, {103, 104}, {1001, 1003}},
	},
	"AddTo3End3rd": {
		[]AddType{11, 12, 101, 1001, 1002}, 1003, true,
		[]rangeset.Span[AddType]{{11, 13}, {101, 102}, {1001, 1004}},
	},
	"AddTo3After3rd": {
		[]AddType{11, 12, 101, 1001, 1002}, 1004, true,
		[]rangeset.Span[AddType]{{11, 13}, {101, 102}, {1001, 1003}, {1004, 1005}},
	},
}

// TestTableAdd is a table driven test that adds an element (using Add() method) to a set taking data from the above addData map.
func TestTableAdd(t *testing.T) {
	for name, data := range addData {
		s := rangeset.Make(data.elts...)
		gotResult := s.Add(data.addElt)
		Assertf(t, gotResult == data.expResult, "%20s: expected result %t, got %t", name, data.expResult, gotResult)
		Assertf(t, len(s) == len(data.expected), "%20s: expected %d ranges, got %d", name, len(data.expected), len(s))
		for i := 0; i < len(data.expected); i++ {
			Assertf(t, s[i].Bot == data.expected[i].Bot, "%20s: range %d got start=%v (expected %v)\n", name, i, s[i].Bot, data.expected[i].Bot)
			Assertf(t, s[i].Top == data.expected[i].Top, "%20s: range %d got end=%v (expected %v)\n", name, i, s[i].Top, data.expected[i].Top)
		}
	}
}

// TestRoundTrip1 uses the addData table (above) to perform round trip tests, by
// converting a set to a string then back to a set.  This provides extra tests of the
// String() method and the NewFromString() and Equal() functions.
func TestRoundTrip1(t *testing.T) {
	for name, data := range addData {
		s := rangeset.Make(data.elts...)
		got, err := rangeset.NewFromString[AddType](s.String())
		Assertf(t, err == nil, "%20s: Round trip string parsing expected no error, got %v", name, err)
		same := rangeset.Equal(s, got)
		Assertf(t, same, "%20s: After round trip comparing %q: expected true, got %t", name, s.String(), same)
	}
}

// TestAddReallocStart checks the special case where append() on the Span slice has to
// reallocate memory due to inserting lots of ranges at the beginning of the slice
func TestAddReallocStart(t *testing.T) {
	s := rangeset.NewFromRange[AddType](20, 40)
	s.Add(2)
	s.Add(4)
	s.Add(6)
	s.Add(8)
	s.Add(10)

	const expected = "{2,4,6,8,10,20:39}"
	got := s.String()
	Assertf(t, got == expected, "%24s: expected %q got %q\n", "AddReallocStart", expected, got)
}

// TestAddReallocMiddle checks the special case where append() on the Span slice has to
// reallocate memory due to inserting lots of ranges in the middle of the slice
func TestAddReallocMiddle(t *testing.T) {
	s, _ := rangeset.NewFromString[AddType]("{10:20,40:50}")
	s.Add(22)
	s.Add(24)
	s.Add(26)
	s.Add(28)
	s.Add(30)
	s.Add(32)

	const expected = "{10:20,22,24,26,28,30,32,40:50}"
	got := s.String()
	Assertf(t, got == expected, "%24s: expected %q got %q\n", "AddReallocMiddle", expected, got)
}

// TestAddReallocEnd checks the special case where append() on the Span slice has to
// reallocate memory due to inserting lots of ranges at the end of the slice
func TestAddReallocEnd(t *testing.T) {
	s := rangeset.Make[AddType](10)
	s.Add(22)
	s.Add(24)
	s.Add(26)
	s.Add(28)
	s.Add(30)
	s.Add(32)

	const expected = "{10,22,24,26,28,30,32}"
	got := s.String()
	Assertf(t, got == expected, "%24s: expected %q got %q\n", "AddReallocEnd", expected, got)
}

// TestAddUniversal tests adding to the start and end of the universal set
func TestAddUniversal(t *testing.T) {
	U := rangeset.Complement(rangeset.Make[AddType]()) // complement of empty set is universal set
	s := U.Copy()
	got := s.Add(maxAddType)
	Assertf(t, got == false, "%24s: expected false got %v\n", "AddUniversal max", got)
	got = s.Add(minAddType)
	Assertf(t, got == false, "%24s: expected false got %v\n", "AddUniversal min", got)

	Assertf(t, rangeset.Equal(s, U), "%24s: after adding to U expected %v got %v\n", "AddUniversal", U, s)
}

var rangeData = map[string]struct {
	elts     []AddType
	bot, top AddType
	expected string
}{
	"AddRangeEmptyAdd1": {[]AddType{}, 42, 43, "{42}"},
	"AddRangeEmptyAdd2": {[]AddType{}, 42, 44, "{42:43}"},

	"AddRangeR1Before1":    {[]AddType{4, 5, 6, 7}, 0, 3, "{0:2,4:7}"},
	"AddRangeR1Before2":    {[]AddType{4, 5, 6, 7}, 2, 3, "{2,4:7}"},
	"AddRangeR1After1":     {[]AddType{4, 5, 6, 7}, 9, 10, "{4:7,9}"},
	"AddRangeR1After2":     {[]AddType{4, 5, 6, 7}, 9, 12, "{4:7,9:11}"},
	"AddRangeExtBelow3":    {[]AddType{4, 5, 6, 7}, 1, 4, "{1:7}"},
	"AddRangeExtBelow1":    {[]AddType{4, 5, 6, 7}, 3, 4, "{3:7}"},
	"AddRangeExtBelow1In1": {[]AddType{4, 5, 6, 7}, 3, 5, "{3:7}"},
	"AddRangeExtAbove1":    {[]AddType{4, 5, 6, 7}, 8, 9, "{4:8}"},
	"AddRangeR1In2Above1":  {[]AddType{4, 5, 6, 7}, 6, 9, "{4:8}"},
	"AddRangeR1Above4":     {[]AddType{4, 5, 6, 7}, 8, 12, "{4:11}"},
	"AddRangeR1InStart1 ":  {[]AddType{4, 5, 6, 7}, 4, 5, "{4:7}"},
	"AddRangeR1InStart2":   {[]AddType{4, 5, 6, 7}, 4, 6, "{4:7}"},
	"AddRangeR1InStart3":   {[]AddType{4, 5, 6, 7}, 4, 7, "{4:7}"},
	"AddRangeR1InAll":      {[]AddType{4, 5, 6, 7}, 4, 8, "{4:7}"},
	"AddRangeR1In2":        {[]AddType{4, 5, 6, 7}, 5, 7, "{4:7}"},
	"AddRangeR1InEnd1":     {[]AddType{4, 5, 6, 7}, 7, 8, "{4:7}"},

	"AddRange3NewAfter1":       {[]AddType{1, 2, 9, 20, 21, 22, 23}, 5, 8, "{1:2,5:7,9,20:23}"},
	"AddRange3AddStart2":       {[]AddType{1, 2, 9, 20, 21, 22, 23}, 5, 9, "{1:2,5:9,20:23}"},
	"AddRange3MergeStart2":     {[]AddType{1, 2, 9, 20, 21, 22, 23}, 5, 10, "{1:2,5:9,20:23}"},
	"AddRange2Replace2":        {[]AddType{1, 2, 9, 20, 21, 22, 23}, 5, 11, "{1:2,5:10,20:23}"},
	"AddRange3AddEnd2":         {[]AddType{1, 2, 9, 20, 21, 22, 23}, 10, 11, "{1:2,9:10,20:23}"},
	"AddRange3Merge2And3":      {[]AddType{1, 2, 9, 20, 21, 22, 23}, 5, 20, "{1:2,5:23}"},
	"AddRange3Inside3":         {[]AddType{1, 2, 9, 20, 21, 22, 23}, 21, 22, "{1:2,9,20:23}"},
	"AddRange3MergeAllBeg1":    {[]AddType{1, 2, 9, 20, 21, 22, 23}, 2, 30, "{1:29}"},
	"AddRange3MergeAllBefore1": {[]AddType{1, 2, 9, 20, 21, 22, 23}, 0, 30, "{0:29}"},
	"AddRange3MergeAllEnd3":    {[]AddType{1, 2, 9, 20, 21, 22, 23}, 0, 21, "{0:23}"},
}

// TestTableAddRange is a table driven test (using rangeData map above) that adds a range of elements to a set (using
// the AddRange() method).  It also tests converting a set to a string (using the String() method()).
func TestTableAddRange(t *testing.T) {
	for name, data := range rangeData {
		s := rangeset.Make(data.elts...)
		s.AddRange(data.bot, data.top)
		got := s.String()
		Assertf(t, got == data.expected, "%24s: expected %q got %q\n", name, data.expected, got)
	}
}

// TestRoundTrip2 uses the rangeData table (above) to perform round trip tests, by
// converting a set to a string then back to a set.  This provides extra tests of the
// String() method and the NewFromString() and Equal() functions.
func TestRoundTrip2(t *testing.T) {
	for name, data := range rangeData {
		s := rangeset.Make(data.elts...)
		got, _ := rangeset.NewFromString[AddType](s.String())
		same := rangeset.Equal(s, got)
		Assertf(t, same, "%20s: After round trip comparing %q: expected true, got %t", name, s.String(), same)
	}
}

// TestAddRangeReallocStart checks the special case where append() on the Span slice has to
// reallocate memory due to inserting lots of ranges at the beginning of the slice
func TestAddRangeReallocStart(t *testing.T) {
	s := rangeset.NewFromRange[AddType](20, 40)
	s.AddRange(2, 3)
	s.AddRange(4, 5)
	s.AddRange(6, 7)
	s.AddRange(8, 9)
	s.AddRange(10, 11)

	const expected = "{2,4,6,8,10,20:39}"
	got := s.String()
	Assertf(t, got == expected, "%24s: expected %q got %q\n", "AddRangeReallocStart", expected, got)
}

// TestAddRangeReallocMiddle checks the special case where append() on the Span slice has to
// reallocate memory due to inserting lots of ranges in the middle of the slice
func TestAddRangeReallocMiddle(t *testing.T) {
	s, _ := rangeset.NewFromString[AddType]("{10:20,40:50}")
	s.AddRange(22, 23)
	s.AddRange(24, 25)
	s.AddRange(26, 27)
	s.AddRange(28, 29)
	s.AddRange(30, 31)
	s.AddRange(32, 33)

	const expected = "{10:20,22,24,26,28,30,32,40:50}"
	got := s.String()
	Assertf(t, got == expected, "%24s: expected %q got %q\n", "AddRangeReallocMiddle", expected, got)
}

// TestAddRangeReallocEnd checks the special case where append() on the Span slice has to
// reallocate memory due to inserting lots of ranges at the end of the slice
func TestAddRangeReallocEnd(t *testing.T) {
	s := rangeset.Make[AddType](10)
	s.AddRange(22, 23)
	s.AddRange(24, 25)
	s.AddRange(26, 27)
	s.AddRange(28, 29)
	s.AddRange(30, 31)
	s.AddRange(32, 33)

	const expected = "{10,22,24,26,28,30,32}"
	got := s.String()
	Assertf(t, got == expected, "%24s: expected %q got %q\n", "AddRangeReallocEnd", expected, got)
}
