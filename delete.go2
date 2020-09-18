package rangeset

// Delete removes an element from the set
// The set is unchanged if the element is not in the set
// It has time complexity O(log r) where r is the number of ranges, O(log n) in the worst case.
func (s *Set[T]) Delete(e T) {
	var endMark = minInt[T]()  // indicates top/bottom of range of valid elements
	idx := s.bsearch(e)
	if idx == 0 || (e >= (*s)[idx-1].t && (*s)[idx-1].t != endMark) {
		return // outside any existing range
	}
	if e == (*s)[idx-1].b && e == (*s)[idx-1].t - 1 {
		// Element matches a single elt range, so delete it
		*s = (*s)[:idx-1 + copy((*s)[idx-1:], (*s)[idx:])]
		return
	}
	if e > (*s)[idx-1].b && e < (*s)[idx-1].t - 1 {
		// Inside an existing range, so split it
		*s = append(*s, Span[T]{})  // add empty element at end
		copy((*s)[idx:], (*s)[idx-1:])  // move ranges to make space to ...
		(*s)[idx-1].t, (*s)[idx].b = e, e + 1
		return
	}
	if e == (*s)[idx-1].b {
		(*s)[idx-1].b = e + 1
		return
	}
	//assert(e == (*s)[idx-1].t - 1)
	(*s)[idx-1].t = e
}

// DeleteRange removes a range of elements from the set
func (s *Set[T]) DeleteRange(b, t T) {
	var endMark = minInt[T]()  // indicates top/bottom of range of valid elements
	if t <= b && t != endMark {
		return  // nothing to delete
	}

	// Work out which "range" of Spans to delete.  Note that the range given by
	// [b,t) may overlap zero or more spans or even be entirely within a span.
	bIdx, tIdx := s.bsearch(b), s.bsearch(t)
	//assert(bIdx >= 0 && bIdx <= len(*s) && tIdx >= 0 && tIdx <= len(*s))
	if bIdx > 0 && b == (*s)[bIdx-1].b {
		bIdx-- // we don't need to keep any of the bottom Span
	}
	//if tIdx > 0 && t < (*s)[tIdx-1].t {
	//	tIdx-- // we need to keep some of the top Span
	//}
	if t == endMark {
		tIdx = len(*s)
	} else if tIdx > 0 && (t < (*s)[tIdx-1].t || (*s)[tIdx-1].t == endMark) {
		tIdx-- // we need to keep some of the top Span
	}

	// At this point the number of Spans to be deleted is given by tIdx-bIdx which can be
	// -1 (Span to be inserted), 0 (only Span ends to be adjusted), 1, ... up to len(*s).

	// Check for situation where the deleted range is entirely within the span with index tIdx
	if tIdx < bIdx {
		//assert(bIdx == tIdx + 1, "tIdx should only be one less than bIdx")
		// Split an existing span in two - insert a new span and adjust the ends
		*s = append(*s, Span[T]{})
		copy((*s)[bIdx:], (*s)[tIdx:])
		(*s)[bIdx].b = t  // bottom of above
		(*s)[tIdx].t = b  // top of below
		return
	}

	//assert(tIdx >= bIdx)
	// Delete all the spans we don't need
	copy((*s)[bIdx:], (*s)[tIdx:])
	*s = (*s)[:len(*s) - (tIdx-bIdx)]

	// Adjust the ends of kept adjacent Spans if necessary
	if bIdx < len(*s) && (t > (*s)[bIdx].b || t == endMark) {
		(*s)[bIdx].b = t
	}
	if bIdx > 0 && (b < (*s)[bIdx-1].t || (*s)[bIdx-1].t == endMark) {
		(*s)[bIdx-1].t = b
	}
}
