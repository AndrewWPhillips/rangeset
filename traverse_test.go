package rangeset_test

import (
	"context"
	"github.com/andrewwphillips/rangeset"
	"iter"
	"testing"
)

type traverseType int // ATM we seem to need a diff. element type in each test file

// traverseData is for table-driven tests of Iterate and Filter methods
var traverseData = map[string]struct {
	in        string
	evenCount int
}{
	"Empty":   {"{}", 0},
	"One":     {"{1}", 0},
	"Two":     {"{2}", 1},
	"Odds":    {"{1,3,5,7,9,11,13,15,17,19}", 0},
	"Evens":   {"{2,4,6,8,10,12,14,16,18,20}", 10},
	"Both":    {"{1:20}", 10},
	"Both2":   {"{2:20}", 10},
	"Range1A": {"{1,4:5}", 1},
	"Range1B": {"{1,4:6}", 2},
	"Range1C": {"{1,4:7}", 2},
	"Range1D": {"{1,4:8}", 3},
	"Range2A": {"{2,4:5}", 2},
	"Range2B": {"{2,4:6}", 3},
	"Range2C": {"{2,4:7}", 3},
	"Range2D": {"{2,4:8}", 4},
	"Range3A": {"{0:3,5}", 2},
	"Range3B": {"{0:3,5:6}", 3},
	"Range3C": {"{0:3,5:7}", 3},
	"Range3D": {"{0:3,5:8}", 4},
	// TODO: add ranges that includes U and E
}

// TestGo1_23Iterator is almost identical to TestIteratorAll but tests the Seq() method (which
// returns a Go 1.23 iterator) instead of the Iterator() method (which returns a channel).
func TestGo1_23Iterator(t *testing.T) {
	for name, data := range traverseData {
		in, _ := rangeset.NewFromString[traverseType](data.in)
		var got []traverseType
		for v := range in.Seq() { // use Go 1.23 range over function
			got = append(got, v)
		}
		values := in.Values()
		Assertf(t, len(got) == len(values), "Go1_23Iterator: %20s: expected %d calls, got %d\n",
			name, len(values), len(got))
		for idx, v := range got {
			Assertf(t, v == values[idx], "Go1_23Iterator: %20s: at position %d expected %v, got %v",
				name, idx, values[idx], v)
		}
	}
}

// TestGo1_23IteratorStop is like TestIteratorCancel but using Seq method rather than Iterator() method
func TestGo1_23IteratorStop(t *testing.T) {
	in := rangeset.Make[traverseType](7, 42, 73, 86, 99)
	next, stop := iter.Pull(in.Seq())

	v, ok := next()
	Assertf(t, ok, "TestGo1_23IteratorStop: Expected first next() status to be true, got %v\n", ok)
	Assertf(t, v == 7, "TestGo1_23IteratorStop: Expected initial next() value to be 7, got %v\n", v)

	stop()
	v, ok = next()
	Assertf(t, !ok, "TestGo1_23IteratorStop: After stop() expected status to be false, got %t\n", ok)
}

// TestGo1_23IteratorEnd is like TestGo1_23IteratorStop but does not stop the iteration
func TestGo1_23IteratorEnd(t *testing.T) {
	in := rangeset.Make[traverseType](49)
	next, stop := iter.Pull(in.Seq())
	defer stop()

	v, ok := next()
	Assertf(t, ok, "TestGo1_23IteratorEnd: Expected first next() status to be true, got %v\n", ok)
	Assertf(t, v == 49, "TestGo1_23IteratorEnd: Expected initial next() value to be 49, got %v\n", v)

	v, ok = next()
	Assertf(t, !ok, "TestGo1_23IteratorEnd: Expected last next() status to be false, got %v\n", ok)
	Assertf(t, v == 0, "TestGo1_23IteratorEnd: Expected last next() value to be zero, got %v\n", v)
}

// TestSpansSeq is a simple test of the SpansSeq() method
func TestSpansSeq(t *testing.T) {
	in := rangeset.NewFromRange(-1, 2)
	next, stop := iter.Pull(in.SpansSeq())
	defer stop()

	v, ok := next()
	Assertf(t, ok, "TestSpansSeq: Expected first next() status to be true, got %v\n", ok)
	Assertf(t, v == rangeset.Span[int]{-1, 2}, "TestSpansSeq: Expected initial next() value to be {-1,2}, got %v\n", v)

	v, ok = next()
	Assertf(t, !ok, "TestSpansSeq: Expected last next() status to be false, got %v\n", ok)
	Assertf(t, v == rangeset.Span[int]{}, "TestSpansSeq: Expected last next() value to be empty, got %v\n", v)

}

// TestIterateMethod tests that Iterate calls the function on every element using opData.union set
func TestIterateMethod(t *testing.T) {
	for name, data := range traverseData {
		in, _ := rangeset.NewFromString[traverseType](data.in)
		got := make([]traverseType, 0, in.Len())
		in.Iterate(func(v traverseType) { got = append(got, v) })
		values := in.Values()
		Assertf(t, len(got) == len(values), "Iterate: %20s: expected %d calls, got %d\n",
			name, len(values), len(got))
		for idx, v := range got {
			Assertf(t, v == values[idx], "Iterate: %20s: at position %d expected %v, got %v",
				name, idx, values[idx], v)
		}
	}
}

// TestFilterNone tests the Filter method keeping all elements
func TestFilterNone(t *testing.T) {
	for name, data := range traverseData {
		in, _ := rangeset.NewFromString[traverseType](data.in)
		origLen := in.Len()
		in.Filter(func(v traverseType) bool { return true }) // keep all elts
		Assertf(t, in.Len() == origLen, "FilterNone: %20s: expected %d elements remaining, got %d\n",
			name, origLen, in.Len())
	}
}

// TestFilterAll tests the Filter method removing all elements
func TestFilterAll(t *testing.T) {
	for name, data := range traverseData {
		in, _ := rangeset.NewFromString[traverseType](data.in)
		in.Filter(func(v traverseType) bool { return false }) // keep no elts
		Assertf(t, in.Len() == 0, "FilterAll: %20s: expected no elements remaining, got %d\n",
			name, in.Len())
	}
}

// TestFilterOdd tests the Filter method by removing all odd elements
func TestFilterOdd(t *testing.T) {
	for name, data := range traverseData {
		in, _ := rangeset.NewFromString[traverseType](data.in)
		in.Filter(func(v traverseType) bool { return v%2 == 0 }) // keep even elts
		Assertf(t, in.Len() == data.evenCount,
			"FilterOdd: %20s: expected %d even elements remaining, got %d\n",
			name, data.evenCount, in.Len())
	}
}

// TestFilterEven tests the Filter method by removing all even elements
func TestFilterEven(t *testing.T) {
	for name, data := range traverseData {
		in, _ := rangeset.NewFromString[traverseType](data.in)
		origLen := in.Len()
		in.Filter(func(v traverseType) bool { return v%2 != 0 }) // keep odd elts
		Assertf(t, in.Len() == origLen-data.evenCount,
			"FilterEven: %20s: expected %d odd elements remaining, got %d\n",
			name, origLen-data.evenCount, in.Len())
	}
}

// TestIteratorAll tests that the chan returned by the Iterator method sends all elements
func TestIteratorAll(t *testing.T) {
	for name, data := range traverseData {
		in, _ := rangeset.NewFromString[traverseType](data.in)
		var got []traverseType
		for v := range in.Iterator(context.Background()) {
			got = append(got, v)
		}
		values := in.Values()
		Assertf(t, len(got) == len(values), "IteratorAll: %20s: expected %d calls, got %d\n",
			name, len(values), len(got))
		for idx, v := range got {
			Assertf(t, v == values[idx], "IteratorAll: %20s: at position %d expected %v, got %v",
				name, idx, values[idx], v)
		}
	}
}

// TestIteratorCancel tests that cancelling an iteration stops sending elements and closes the chan
func TestIteratorCancel(t *testing.T) {
	in := rangeset.Make[traverseType](7, 42, 73, 86, 99)
	ctx, cancel := context.WithCancel(context.Background())
	ch := in.Iterator(ctx)
	v, ok := <-ch
	Assertf(t, ok, "IteratorCancel: Expected initial channel read status to be true, got %v\n", ok)
	Assertf(t, v == 7, "IteratorCancel: Expected initial channel value to be 7, got %v\n", v)
	cancel()
	// Note that after canceling the context we may or may not get a final value before the chan is closed
	_, _ = <-ch
	_, ok = <-ch
	Assertf(t, !ok, "IteratorCancel: After cancel expected read status to be false, got %t\n", ok)
}

// TestChanRoundTrip tests the ReadAll (and Iterator) methods using a "round trip" - ie write the elements of
// a set to a chan that is read into a new set, then checks that the sets are the same.
func TestChanRoundTrip(t *testing.T) {
	for name, data := range traverseData {
		in, _ := rangeset.NewFromString[traverseType](data.in)
		var out rangeset.Set[traverseType]
		out.ReadAll(context.Background(), in.Iterator(context.Background()))
		Assertf(t, rangeset.Equal(in, out), "ChanRoundTrip: %12s: expected %v, got %v", name, in, out)
	}
}
