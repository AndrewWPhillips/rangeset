package rangeset

// bsearch.go implements an internal method to search a set

// bsearch performs a binary search and returns the index of the Span (range) that
// is immediately above the element searched for (irrespective of whether or not
// the value is found in the set). The number of possible return values is one more
// than the number of ranges in the set. So for example if you search for a value
// that is less than any element in the set then 0 is returned; if you search for a
// value that is in or beyond the last range then the number of ranges is returned.
// Note that if the set includes the smallest element, given by minInt[T](), then zero
// (below bottom range) cannot be returned, since no element can be below the smallest.
func (s Set[T]) bsearch(value T) int {
	bot, top := 0, len(s)
	for bot < top {
		curr := bot + (top-bot)/2
		if value < s[curr].Bot {
			top = curr
		} else {
			bot = curr + 1
		}
	}
	//assert(bot == top, "bottom should not be greater than top")
	return bot
}
