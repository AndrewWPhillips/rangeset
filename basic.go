package rangeset

import (
	"slices"
)

// basic.go implements basic functions such as creating range sets, checking if a set contains
// an element, number of elements, adding (union), subtraction, intersection etc.

type Element interface {
	//constraints.Integer
	~int | ~int8 | ~int16 | ~int32 | ~int64 |
		~uint | ~uint8 | ~uint16 | ~uint32 | ~uint64 | ~uintptr
}

type (
	Span[T Element] struct{ Bot, Top T } // one range [b..t) ie asymmetric bounds
	Set[T Element]  []Span[T]            // a set is just a slice of ranges
)

// Make creates a new set optionally taking initial element(s)
// TODO: check if sorting elems (in copy of it) could improve perf. for large len(elems)
func Make[T Element](elems ...T) Set[T] {
	s := make(Set[T], 0, len(elems))
	for _, v := range elems {
		_ = s.Add(v)
	}
	return s
}

// NewFromRange creates a new set by specifying an initial range of elements
// Equivalent to: Make[T]() then AddRange()
func NewFromRange[T Element](b, t T) Set[T] {
	if t <= b {
		//panic("Invalid range in NewFromRange")
		// Return an empty set (this is consistent with for loop where end < start)
		return nil
	}
	return append(make(Set[T], 0, 1), Span[T]{b, t})
}

// Universal returns the set of all elements
// Note to create an Empty set simply use Make[T]()
// Equivalent to: Complement(Make[T]())
func Universal[T Element]() Set[T] {
	var endMark = minInt[T]()
	// Integer "wrap-around" means that "one more" than max element is min element
	// assert(maxInt[T]() + 1 == minInt[T]())
	return Set[T]{{endMark, endMark}}
}

// Length returns the number of elements and number of ranges in the set.
// Note if the element type is 64-bit the size of a universal set is too large to be represented
// as uint64 - in this case 0 is returned for the number of elements and 1 for the number of spans.
// It has time complexity of O(r) where r is the number of ranges, and O(n) in the worst case.
func (s Set[T]) Length() (length uint64, spans int) {
	spans = len(s)
	for _, r := range s {
		// assert(r.t > r.b)
		length += uint64(r.Top - r.Bot)
	}
	return
}

// Len returns the number of elements, or -1 if it's more than the largest int.
// It has time complexity of O(r) where r is the number of ranges, and O(n) in the worst case.
// Note: As sets are stored using ranges it is easy to have huge sets, where the number of
// elements is too large for an int.  You should prefer using the Length() method above if
// the number of elements can possibly be >= 2^31 (32-bit systems) or 2^63 (64-bit systems).
func (s Set[T]) Len() int {
	length, spans := s.Length()
	if length > uint64(^uint(0)>>1) || length == 0 && spans > 0 {
		return -1
	}
	return int(length)
}

// Contains tests whether a set contains an element
// It has time complexity O(log r) where r is the number of ranges of the set (since it
// does a binary search over the ranges).  In the worst case it is O(log n).
func (s Set[T]) Contains(e T) bool {
	idx := s.bsearch(e)
	var endMark = minInt[T]() // in a range it flags: bottom/top of all valid elements
	return idx > 0 && (e < s[idx-1].Top || s[idx-1].Top == endMark)
}

// Values returns all the values in the set as a slice (in numeric order).
// WARNING: if your range set contains large ranges this may take a
// long time and return a slice with a large number of elements.
//
// Deprecated: from Go 12.3 use the Seq() method to obtain an iterator
func (s Set[T]) Values() []T {
	retval := make([]T, 0, s.Len())
	for _, v := range s {
		for e := v.Bot; e < v.Top; e++ {
			retval = append(retval, e)
		}
	}
	return retval
}

// Spans returns all the ranges of the set as a slice of "Span" structures.
// Note that these use asymmetric ranges where the t (top) field is one more than the
// last element in the range. The Spans are sorted within the slice and do not overlap.
//
// Deprecated: from Go 1.23 use the SpansSeq() method to obtain an iterator
func (s Set[T]) Spans() []Span[T] {
	retval := make([]Span[T], 0, len(s))
	for _, v := range s {
		retval = append(retval, v)
	}
	return retval
}

// Copy makes a copy of a set
func (s Set[T]) Copy() Set[T] {
	return Set[T](slices.Clone(s))
}

// AddSet finds the union of s with s2 (ie, adds all the elements of s2 to s)
func (s *Set[T]) AddSet(s2 Set[T]) {
	for _, v := range s2 {
		s.AddRange(v.Bot, v.Top)
	}
}

// SubSet removes all elements of s2 from s
func (s *Set[T]) SubSet(s2 Set[T]) {
	for _, v := range s2 {
		s.DeleteRange(v.Bot, v.Top)
	}
}

// Intersect finds the intersection of s with s2 (ie, deletes from s any elts not in s2)
func (s *Set[T]) Intersect(s2 Set[T]) {
	endMark := minInt[T]()
	bDel := endMark
	for _, v := range s2 {
		if bDel != endMark || v.Bot != endMark {
			s.DeleteRange(bDel, v.Bot)
		}
		bDel = v.Top
	}
	if len(s2) == 0 || bDel != endMark {
		s.DeleteRange(bDel, endMark)
	}
}
