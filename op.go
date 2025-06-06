package rangeset

// op.go implement set operations like union, etc

// Equal compares two sets
// TODO: extend to compare any number of sets
func Equal[T Element](s1, s2 Set[T]) bool {
	if len(s1) != len(s2) {
		return false
	}
	for idx := 0; idx < len(s1); idx++ {
		if s1[idx].Bot != s2[idx].Bot || s1[idx].Top != s2[idx].Top {
			return false
		}
	}
	return true
}

// Complement creates the inverse set where the "universal" set (the set of all
// elements) is all the integers of the element type. For example, for an element
// type of int8 the universal set is the integers -128 to +127 and the inverse
// (using our string serialisation notation) of {1:10} is {-128:0,11:127}
// TODO: make this a method
func Complement[T Element](s Set[T]) Set[T] {
	var endMark = minInt[T]() // indicates top/bottom of range of valid elements
	retval := make(Set[T], 0, len(s)+1)
	if len(s) == 0 {
		// Inverse of empty set is U (set of all valid elements)
		retval = append(retval, Span[T]{endMark, endMark}) // represents U
		return retval
	}
	bot := endMark
	for idx, v := range s {
		if idx > 0 || v.Bot != endMark {
			retval = append(retval, Span[T]{bot, v.Bot})
		}
		bot = v.Top
	}
	if bot != endMark {
		retval = append(retval, Span[T]{bot, endMark})
	}

	return retval
}

// Union finds the union of zero or more sets and returns a new set
func Union[T Element](sets ...Set[T]) Set[T] {
	if len(sets) == 0 {
		return Set[T]{} // return empty set
	}
	// Copy one set and add (union) all the other sets to it
	// TODO: check if copying largest set then merging others is faster for common scenarios
	retval := sets[0].Copy()
	for _, other := range sets[1:] {
		retval.AddSet(other)
	}
	return retval
}

// Intersect finds the intersection of zero or more sets, returning a new set
func Intersect[T Element](sets ...Set[T]) Set[T] {
	if len(sets) == 0 {
		return Set[T]{} // return empty set
	}

	retval := sets[0].Copy()
	for _, other := range sets[1:] {
		retval.Intersect(other)
	}

	return retval
}
