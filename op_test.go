package rangeset_test

import (
	"github.com/andrewwphillips/rangeset"
	"testing"
)

// signedData contains table data specifically for testing Complement() on sets of int8
var signedData = map[string]struct {
	in, expected string
}{
	// Note: these test values assume an element type of int8
	"Empty":     {"{}", "{-128:127}"},
	"Universal": {"{-128:127}", "{}"},
	"Zero":      {"{0}", "{-128:-1,1:127}"},
	"Bottom":    {"{-128}", "{-127:127}"},
	"Top":       {"{127}", "{-128:126}"},
	"Both":      {"{-128:-127,126:127}", "{-126:125}"},
	"TwoBot":    {"{-128:0,2}", "{1,3:127}"},
	"TwoTop":    {"{-5:-3,100:127}", "{-128:-6,-2:99}"},
	"TwoInside": {"{-1,1}", "{-128:-2,0,2:127}"},
	"EndsMid":   {"{-128,0,127}", "{-127:-1,1:126}"},
}

// TestComplementSigned tests inverting a set that has a signed integer element type
func TestComplementSigned(t *testing.T) {
	for name, data := range signedData {
		s, _ := rangeset.NewFromString[int8](data.in)

		// Take the complement of the set and check it is as expected
		s2 := rangeset.Complement(s)
		expected, _ := rangeset.NewFromString[int8](data.expected)
		Assertf(t, rangeset.Equal(s2, expected), "ComplementSigned: %16s: expected %q got %q\n",
			name, expected, s2)

		// Take the complement of the complement and check it matches the original
		s3 := rangeset.Complement(s2)
		Assertf(t, rangeset.Equal(s3, s), "ComplementSigned: %16s: expected %q got %q\n",
			"reverse "+name, s, s3)
	}
}

// unsignedData contains data specifically for testing Complement() on sets of uint16
var unsignedData = map[string]struct {
	in, expected string
}{
	"Empty":     {"{}", "{0:65535}"},
	"Bottom":    {"{0}", "{1:65535}"},
	"Top":       {"{65535}", "{0:65534}"},
	"Both":      {"{0:5,65530:65535}", "{6:65529}"},
	"TwoBot":    {"{0:100,200}", "{101:199,201:65535}"},
	"TwoTop":    {"{10,100:65535}", "{0:9,11:99}"},
	"TwoInside": {"{100,1000}", "{0:99,101:999,1001:65535}"},
	"ThreeEnds": {"{0:100,200,300:65535}", "{101:199,201:299}"},
}

// TestComplementUnsigned tests inverting a set that has a unsigned integer (uint16) elements
func TestComplementUnsigned(t *testing.T) {
	for name, data := range unsignedData {
		s, _ := rangeset.NewFromString[uint16](data.in)

		// Take the complement of the set and check it is as expected
		s2 := rangeset.Complement(s)
		expected, _ := rangeset.NewFromString[uint16](data.expected)
		Assertf(t, rangeset.Equal(s2, expected), "Complement Unsigned: %16s: expected %q got %q\n",
			name, expected, s2)

		// Take the complement of the complement and check it matches the original
		s3 := rangeset.Complement(s2)
		Assertf(t, rangeset.Equal(s3, s), "Complement Unsigned: %16s: expected %q got %q\n",
			"reverse "+name, s, s3)
	}
}

// TODO: int64 and uint64 complement tests (endMark)

// opData provides table data for testing various operations: Union(), Intersect(), etc
var opData = map[string]struct {
	in               []string // zero or more sets to be combined (intersected, etc)
	union, intersect string
	sub              string
}{
	"Empty0":     {[]string{}, "{}", "{}", "{}"},
	"Empty1":     {[]string{"{}"}, "{}", "{}", "{}"},
	"Empty2":     {[]string{"{}", "{}"}, "{}", "{}", "{}"},
	"Single1":    {[]string{"{1}"}, "{1}", "{1}", "{1}"},
	"Single2":    {[]string{"{1}", "{1}"}, "{1}", "{1}", "{}"},
	"Mult1":      {[]string{"{1,3,5}"}, "{1,3,5}", "{1,3,5}", "{1,3,5}"},
	"Mult2":      {[]string{"{1,3,5}", "{7,9}"}, "{1,3,5,7,9}", "{}", "{1,3,5}"},
	"EmptyAdd":   {[]string{"{}", "{1}"}, "{1}", "{}", "{}"},
	"EmptyAddR1": {[]string{"{}", "{1:2}", "{3:4}"}, "{1:4}", "{}", "{}"},
	"EmptyAddR2": {[]string{"{}", "{1}", "{3:4}"}, "{1,3:4}", "{}", "{}"},
	"3Adj":       {[]string{"{1}", "{2}", "{3}"}, "{1:3}", "{}", "{1}"},
	"3Sep":       {[]string{"{1}", "{3}", "{5}"}, "{1,3,5}", "{}", "{1}"},
	"SepRanges":  {[]string{"{1:5}", "{7:8}"}, "{1:5,7:8}", "{}", "{1:5}"},
	"Touch":      {[]string{"{1:5}", "{6:8}"}, "{1:8}", "{}", "{1:5}"},
	"Subset":     {[]string{"{1:5}", "{2:4}"}, "{1:5}", "{2:4}", "{1,5}"},
	"Subset2":    {[]string{"{1:5}", "{2:4}", "{3}"}, "{1:5}", "{3}", "{1,5}"},
	"Overlap1":   {[]string{"{1:5}", "{4:8}"}, "{1:8}", "{4:5}", "{1:3}"},
	"Overlap2":   {[]string{"{1:5}", "{4:8}", "{8:9}"}, "{1:9}", "{}", "{1:3}"},
	"SubLap":     {[]string{"{1:5}", "{3:4}", "{2:9}"}, "{1:9}", "{3:4}", "{1}"},
	"SubLap2":    {[]string{"{1:5}", "{3:4}", "{3:9}"}, "{1:9}", "{3:4}", "{1:2}"},
	"3SubSep":    {[]string{"{1:5}", "{2:4}", "{7:8}"}, "{1:5,7:8}", "{}", "{1,5}"},
	"3SubSep2":   {[]string{"{1:5}", "{7:8}", "{2:4}"}, "{1:5,7:8}", "{}", "{1,5}"},
	"3LapSep":    {[]string{"{1:5}", "{5:6}", "{8:9}"}, "{1:6,8:9}", "{}", "{1:4}"},
	"3LapSep2":   {[]string{"{1:5}", "{8:9}", "{5:6}"}, "{1:6,8:9}", "{}", "{1:4}"},
	"3Way":       {[]string{"{4,7}", "{1,7}", "{1,4}"}, "{1,4,7}", "{}", "{}"},

	"UAndU":      {[]string{"{U}", "{U}"}, "{U}", "{U}", "{}"},
	"UAndEmpty":  {[]string{"{U}", "{}"}, "{U}", "{}", "{U}"},
	"UAndOne":    {[]string{"{U}", "{1}"}, "{U}", "{1}", "{E:0,2:E}"},
	"OneAndU":    {[]string{"{1}", "{U}"}, "{U}", "{1}", "{}"},
	"UAndTwo":    {[]string{"{U}", "{4}", "{2}"}, "{U}", "{}", "{E:1,3,5:E}"},
	"EmptyAndU":  {[]string{"{}", "{U}"}, "{U}", "{}", "{}"},
	"StartPlus2": {[]string{"{E:10}", "{5}", "{20}"}, "{E:10,20}", "{}", "{E:4,6:10}"},
}

// TestUnionFunction performs a union of opData table's "in" sets and compares to "union" set
func TestUnionFunction(t *testing.T) {
	for name, data := range opData {
		var sets []rangeset.Set[uint]
		expected, _ := rangeset.NewFromString[uint](data.union)

		for _, str := range data.in {
			s, _ := rangeset.NewFromString[uint](str)
			sets = append(sets, s)
		}
		got := rangeset.Union(sets...)
		Assertf(t, rangeset.Equal(got, expected), "        Union: %20s: expected %q got %q\n",
			name, expected, got)

		// Reverse the order and check we get the same result
		sets = nil
		for idx := len(data.in); idx > 0; idx-- {
			s, _ := rangeset.NewFromString[uint](data.in[idx-1])
			sets = append(sets, s)
		}
		got = rangeset.Union(sets...)
		Assertf(t, rangeset.Equal(got, expected), "Reverse Union: %20s: expected %q got %q\n",
			name, expected, got)
	}
}

// TestIntersectFunction performs the intersection of opData table's "in" sets and compares to "intersect" set
func TestIntersectFunction(t *testing.T) {
	for name, data := range opData {
		var sets []rangeset.Set[uint]
		expected, _ := rangeset.NewFromString[uint](data.intersect)

		for _, str := range data.in {
			s, _ := rangeset.NewFromString[uint](str)
			sets = append(sets, s)
		}
		got := rangeset.Intersect(sets...)
		Assertf(t, rangeset.Equal(got, expected), "For Intersect: %20s: expected %q got %q\n",
			name, expected, got)

		// Reverse the order and check we get the same result
		sets = nil
		for idx := len(data.in); idx > 0; idx-- {
			s, _ := rangeset.NewFromString[uint](data.in[idx-1])
			sets = append(sets, s)
		}
		got = rangeset.Intersect(sets...)
		Assertf(t, rangeset.Equal(got, expected), "Rev Intersect: %20s: expected %q got %q\n",
			name, expected, got)
	}
}

// TestCopyMethod checks copying of a set using the "union" field of the above opData table
func TestCopyMethod(t *testing.T) {
	for name, data := range opData {
		orig, _ := rangeset.NewFromString[uint](data.union)
		got := orig.Copy()
		Assertf(t, rangeset.Equal(orig, got), " Set Copy: %20s: copying %q got %q\n",
			name, orig, got)
	}
}

// TestIntersectMethod tests the Intersect method
func TestIntersectMethod(t *testing.T) {
	for name, data := range opData {
		if len(data.in) == 0 {
			continue
		}
		got, _ := rangeset.NewFromString[uint](data.in[0])
		for idx := 1; idx < len(data.in); idx++ {
			s, _ := rangeset.NewFromString[uint](data.in[idx])
			got.Intersect(s)
		}
		expected, _ := rangeset.NewFromString[uint](data.intersect)
		Assertf(t, rangeset.Equal(got, expected), " Intersect: %20s: expected %q got %q\n",
			name, expected, got)
	}
}

// TestSubSetMethod tests the SubSet method
func TestSubSetMethod(t *testing.T) {
	for name, data := range opData {
		if len(data.in) == 0 {
			continue
		}
		got, _ := rangeset.NewFromString[uint](data.in[0])
		for idx := 1; idx < len(data.in); idx++ {
			s, _ := rangeset.NewFromString[uint](data.in[idx])
			got.SubSet(s)
		}
		expected, _ := rangeset.NewFromString[uint](data.sub)
		Assertf(t, rangeset.Equal(got, expected), "    SubSet: %20s: expected %q got %q\n",
			name, expected, got)
	}
}

// TestAddSetMethod tests the AddSet method
func TestAddSetMethod(t *testing.T) {
	for name, data := range opData {
		if len(data.in) == 0 {
			continue
		}
		got, _ := rangeset.NewFromString[uint](data.in[0])
		for idx := 1; idx < len(data.in); idx++ {
			s, _ := rangeset.NewFromString[uint](data.in[idx])
			got.AddSet(s)
		}
		expected, _ := rangeset.NewFromString[uint](data.union)
		Assertf(t, rangeset.Equal(got, expected), "    AddSet: %20s: expected %q got %q\n",
			name, expected, got)
	}
}
