package rangeset

// add.go implements methods to add to a rangeset

// Add inserts a single element into the set
// It returns true if added or false if it already existed in the set
// It has time complexity O(log r) where r is the number of ranges or O(log n) worst case.
func (s *Set[T]) Add(e T) bool {
	idx := s.bsearch(e)
	//assert(idx >= 0 && idx <= len(*s))
	var endMark = minInt[T]() // in a range it flags: bottom/top of all valid elements
	if idx == 0 || (e > (*s)[idx-1].Top && (*s)[idx-1].Top != endMark) {
		// New element is before range [idx] and after range [idx-1] (+ not just past end)
		if idx < len(*s) && e == (*s)[idx].Bot-1 {
			// Extend range [idx] down by one
			(*s)[idx].Bot = e
		} else {
			// Add new range between [idx-1] and [idx] (incl. before 1st and after last)
			*s = append(*s, Span[T]{})
			copy((*s)[idx+1:], (*s)[idx:])
			(*s)[idx] = Span[T]{e, e + 1}
		}
		return true
	}
	// assert(idx > 0)
	if e == (*s)[idx-1].Top && (*s)[idx-1].Top != endMark {
		// New element is just past the end of range [idx-1]
		if idx < len(*s) && e == (*s)[idx].Bot-1 {
			// New element joins range [idx-1] and [idx] together
			e = (*s)[idx-1].Bot            // save beginning of previous range
			copy((*s)[idx-1:], (*s)[idx:]) // move rest of the ranges down
			(*s)[idx-1].Bot = e            // restore beginning of current
			*s = (*s)[:len(*s)-1]          // remove duplicate end range
		} else {
			// Extend range [idx-1] up by one
			(*s)[idx-1].Top = e + 1
		}
		return true
	}
	return false // e is already in range [idx-1]
}

// AddRange inserts a range of values into the set
// The range is specified using asymmetric bounds: b (1st param) is the lowest element
// of the range to be added and t (2nd param) is one more than the highest element
// Like Add() above it has time complexity O(log r) - or O(log n) in the worst case.
func (s *Set[T]) AddRange(b, t T) {
	var endMark = minInt[T]() // indicates top/bottom of range of valid elements
	if t <= b && t != endMark {
		return // nothing needs to be added
	}

	// Work out where Spans will be added or deleted. Note that since the added
	// range can cover multiple Spans we may have to delete 0 or more, but if the
	// range doesn't overlap any existing Spans we have to insert one.
	var bIdx, tIdx int
	bIdx = s.bsearch(b)
	//assert(bIdx <= len(*s) && tIdx <= len(*s) // must be valid Span index or 1 past end
	if bIdx == 0 || (b > (*s)[bIdx-1].Top && (*s)[bIdx-1].Top != endMark) {
		bIdx++ // past the end of the idx-1 span
	}
	if t == endMark {
		tIdx = len(*s)
	} else {
		tIdx = s.bsearch(t)
	}

	// At this point the number of Spans to be deleted is given by tIdx-bIdx which can be
	// -1 - Span to be inserted, 0 - no Spans added/deleted (but some Span ends may need
	// to be adjusted), or greater than zero - Span(s) to be deleted.

	// Check for situation where range is outside any Span
	if tIdx < bIdx {
		//assert(bIdx == tIdx + 1, "tIdx should only be one less than bIdx")
		*s = append(*s, Span[T]{})
		copy((*s)[bIdx:], (*s)[tIdx:])
		(*s)[tIdx].Bot, (*s)[tIdx].Top = b, t
		return
	}

	//assert(tIdx > 0 && tIdx >= bIdx)
	if (t < (*s)[tIdx-1].Top && t != endMark) || (*s)[tIdx-1].Top == endMark {
		t = (*s)[tIdx-1].Top
	}
	// Delete the spans we don't need
	copy((*s)[bIdx:], (*s)[tIdx:])
	*s = (*s)[:len(*s)-(tIdx-bIdx)]

	// Adjust, as necessary, the ends of the retained span
	if bIdx > 0 && b < (*s)[bIdx-1].Bot {
		(*s)[bIdx-1].Bot = b
	}
	if bIdx > 0 && (t > (*s)[bIdx-1].Top || t == endMark) && (*s)[bIdx-1].Top != endMark {
		(*s)[bIdx-1].Top = t
	}
}
