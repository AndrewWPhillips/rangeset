package rangeset

// traverse.go has methods to process the elements of a set
//  Iterate and Filter methods - use a function to operate on the whole set
//  Iterator and ReadAll - use channels of the element type
//  Seq - returns an iterator (Go 1.23) over all elements of the set

import (
	"context"
	"iter"
)

// Seq returns a Go 1.23 iterator of the set elements in order
func (s Set[T]) Seq() iter.Seq[T] {
	return func(yield func(T) bool) {
		for _, v := range s {
			for e := v.Bot; e < v.Top; e++ {
				if !yield(e) {
					return
				}
			}
		}
	}
}

// SpansSeq returns a Go 1.23 iterator of the ranges of the set
func (s Set[T]) SpansSeq() iter.Seq[Span[T]] {
	return func(yield func(Span[T]) bool) {
		for _, v := range s {
			if !yield(v) {
				return
			}
		}
	}
}

// Iterate calls f on every element in the set.
func (s Set[T]) Iterate(f func(T)) {
	for _, v := range s {
		for e := v.Bot; e < v.Top; e++ {
			f(e)
		}
	}
}

// Filter deletes elements from s for which f returns false.
// Time complexity is O(n)
// TODO look at optimising (eg by bunching into whole ranges before adding to toDelete)
func (s *Set[T]) Filter(f func(T) bool) {
	toDelete := Make[T]()
	for _, v := range *s {
		for e := v.Bot; e < v.Top; e++ {
			if !f(e) {
				toDelete.Add(e) // keep track of elts to delete
			}
		}
	}
	s.SubSet(toDelete) // delete the elts
}

// Iterator returns a channel that receives the elements of the set in order.
// To terminate the go-routine before all elements have been seen, cancel the
// context (first parameter to the method).  Note that after the context is
// canceled another element *may* be seen on the chan before it is closed.
// To avoid a goroutine leak you need to read from the channel until it is
// closed or cancel the context.
func (s Set[T]) Iterator(ctx context.Context) <-chan T {
	r := make(chan T)
	go func(ch chan<- T) {
		defer close(ch)
		for _, v := range s {
			for e := v.Bot; e < v.Top; e++ {
				select {
				case <-ctx.Done():
					return
				case ch <- e:
				}
			}
		}
	}(r)
	return r
}

// ReadAll adds elements to the set by reading them from the channel until it
// is closed or the context is canceled.
func (s *Set[T]) ReadAll(ctx context.Context, c <-chan T) {
	for {
		select {
		case <-ctx.Done():
			return
		case v, ok := <-c:
			if !ok {
				return
			}
			_ = s.Add(v)
		}
	}
}
