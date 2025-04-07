package rangeset_test

import (
	"github.com/andrewwphillips/rangeset"
	"testing"
)

type NewElementType int // set element type used for most of the tests

var newData = map[string]struct {
	elts     []NewElementType // set elements to be passed to Make
	length   int
	expected []rangeset.Span[NewElementType] // expected set of spans
}{
	"NewEmpty":           {[]NewElementType{}, 0, []rangeset.Span[NewElementType]{}},
	"NewOne":             {[]NewElementType{42}, 1, []rangeset.Span[NewElementType]{{42, 43}}},
	"NewTwo":             {[]NewElementType{42, 44}, 2, []rangeset.Span[NewElementType]{{42, 43}, {44, 45}}},
	"NewTwoMerge":        {[]NewElementType{43, 42}, 2, []rangeset.Span[NewElementType]{{42, 44}}},
	"NewTwoSame":         {[]NewElementType{43, 43}, 1, []rangeset.Span[NewElementType]{{43, 44}}},
	"NewThree":           {[]NewElementType{2, 0, -2}, 3, []rangeset.Span[NewElementType]{{-2, -1}, {0, 1}, {2, 3}}},
	"NewThreeMerge":      {[]NewElementType{0, 1, -1}, 3, []rangeset.Span[NewElementType]{{-1, 2}}},
	"NewThreeMergeStart": {[]NewElementType{2, 4, 1}, 3, []rangeset.Span[NewElementType]{{1, 3}, {4, 5}}},
	"NewThreeMergeEnd":   {[]NewElementType{3, 4, 1}, 3, []rangeset.Span[NewElementType]{{1, 2}, {3, 5}}},
	"NewThreeSame":       {[]NewElementType{42, 42, 42}, 1, []rangeset.Span[NewElementType]{{42, 43}}},
	"NewFourWithDupes":   {[]NewElementType{1, 3, 3, 1}, 2, []rangeset.Span[NewElementType]{{1, 2}, {3, 4}}},
	"NewFour2Ranges":     {[]NewElementType{1, 2, 4, 5}, 4, []rangeset.Span[NewElementType]{{1, 3}, {4, 6}}},
}

// TestTableNew calls rangeset.Make() with various sets of parameters from above table (newData)
func TestTableNew(t *testing.T) {
	for name, data := range newData {
		s := rangeset.Make(data.elts...)
		Assertf(t, s.Len() == data.length, "%20s: expected %d elements (got %d)\n", name, data.length, s.Len())
		Assertf(t, len(s) == len(data.expected), "%20s: expected %d ranges (got %d)\n", name, len(data.expected), len(s))
		for i := 0; i < len(data.expected); i++ {
			Assertf(t, s[i].Bot == data.expected[i].Bot, "%20s: range %d start=%v (expected %v)\n", name, i, s[i].Bot, data.expected[i].Bot)
			Assertf(t, s[i].Top == data.expected[i].Top, "%20s: range %d end=%v (expected %v)\n", name, i, s[i].Top, data.expected[i].Top)
		}
	}
}
