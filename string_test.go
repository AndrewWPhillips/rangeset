package rangeset_test

import (
	"github.com/andrewwphillips/rangeset"
	"testing"
)

type StringElementType int16

var stringData = map[string]struct {
	in, expected string
}{
	"StringEmpty":    {"{}", "{}"},
	"StringNegative": {"{-1}", "{-1}"},
	// The following inputs are not in the normal serialisation format but will work
	"StringMerge1":    {"{1,2,4}", "{1:2,4}"},
	"StringMerge2":    {"{1,3,4}", "{1,3:4}"},
	"StringMergeAll":  {"{1,2,3}", "{1:3}"},
	"StringReverse":   {"{5,2}", "{2,5}"},
	"StringDuplicate": {"{2,1,2}", "{1:2}"},
	"StringMixed":     {"{3,1,8,2}", "{1:3,8}"},
	"StringMixRanges": {"{1:2,5,2:3,8}", "{1:3,5,8}"},

	// Note: the following assume that we are testing with int16 element type
	"Universal":  {"{U}", "{-32768:32767}"},
	"Universal2": {"{E:E}", "{-32768:32767}"},
	"BeginMark":  {"{E:10}", "{-32768:10}"},
	"EndMark":    {"{10:E}", "{10:32767}"},
	"TwoEnd":     {"{-1,100:E}", "{-1,100:32767}"},
	"TwoBegin":   {"{E:1,100}", "{-32768:1,100}"},
}

// TestString uses stringData map to perform (table-driven) tests of String() and NewFromString()
// Note that these mainly test unusual cases for NewFromString.  There are many other tests of
// String() and NewFromString() in "round trip" tests like TestRoundTripAdd() etc.
func TestString(t *testing.T) {
	for name, data := range stringData {
		s, _ := rangeset.NewFromString[StringElementType](data.in)
		got := s.String()
		Assertf(t, got == data.expected, "TestString: %20s: expected %q got %q", name, data.expected, got)
	}
}

var stringErrorData = map[string]struct {
	in string
}{
	"EmptyString":    {""},
	"NoBraces":       {"1:2"},
	"NoBraceLeft":    {"1:2}"},
	"NoBraceRight":   {"{1:2"},
	"RangeBadLeft":   {"{1.2:3}"},
	"RangeBadRight":  {"{1:!}"},
	"RangeMissRight": {"{1:}"},
	"ValueBad":       {"{ABC}"},
	"Value2Bad":      {"{1:2,#}"},
	"RangeLess":      {"{2:1}"},
	"RangeRandom":    {"{%89i:djsa.mdaja,esreiop}"},
	"RangeComma":     {"{1;3:4}"},
	"RangeColon":     {"{1,3-4}"},
}

// TestStringError checks that malformed set strings cause an error in NewFromString
func TestStringError(t *testing.T) {
	for name, data := range stringErrorData {
		_, err := rangeset.NewFromString[StringElementType](data.in)
		Assertf(t, err != nil, "StringError: %16s: expected an error got %v", name, err)
	}
}
