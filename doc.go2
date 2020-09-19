// Package rangeset implement a "set" container that uses ranges to efficiently store and
// manipulate common types of large sets.  Due to it's use of ranges the type parameter (the
// set's element type) must be orderable (ie support < operation etc) and also support
// incrementation (++ operation).  Hence the only element types that can currently be used
// are the integers (int, uint, byte, int64, etc).
//
// Note that (until parametric polymorphism is added to GO) use of the package requires running the
// go2go tool to convert .go2 files into .go files for specific element types in use.  Hence the
// **repository "source" files all have extension .go2**. However, one generated file (doc.go)
// is checked in to avoid go get/build failing due to an "empty" download package.
//
// It is similar to the example go2go "sets" package (src\cmd\go2go\testdata\go2path\src\sets in
// Go source repo) but has some pros and cons.  For many common uses of sets (ie those with large
// contiguous ranges) it can be much more space efficient.  On the other hand, adding, deleting,
// and finding elements has time complexity of O(log r), where r is the number of ranges, and
// O(log n) in the worst case, but can still be faster in practice (despite the example "sets"
// package mainly having constant time complexity).
//
// Apart from performance benefits, it also has the advantage that the elements are ordered (eg
// the Values() method returns the elements in order.  It also has additional facilities such as
// the ability to take the complement of a set, and create a Universal set (a set of all elements
// of the element type). It also has methods for serializing and de-serializing sets as strings.
// A useful method is "Spans" that returns a slice of all the ranges of the set.
//
// Creating sets
//
//  s := rangeset.Make[int]()           // create an empty set of ints
//  s := rangeset.Make[byte](1, 2, 7)   // create a set of bytes with 3 elements
//  s := rangeset.Universal[uint64]()   // create a set of all uint64 elements
//  s2 := s.Copy()                      // make a copy of a set
//
// Adding Elements
//
//  s.Add(42)                          // returns true if added or false if already present
//  s.AddRange('a', 'z'+1)             // add a range (using asymmetric bounds)
//  s.AddSet(s2)                       // s becomes the union of s and s2
//  s.ReadAll(ctx, ch)                 // adds from a chan until it's closed (or ctx is cancelled)
//
// Deleting Elements
//
//  s.Delete(42)                       // remove an element - does nothing if not present
//  s.DeleteRange(1, 11)               // delete elements 1 to 10 (inclusive)
//  s.SubSet(s2)                       // remove from s all the elements of s2 (if present)
//  s.Intersect(s2)                    // remove from s all the elements *not* in s2
//
// Getting Elements
//
//  b := s.Contains(42)                // returns true if element is present
//  all := s.Values()                  // returns a slice with all elements in the set (in order)
//  all := s.Spans()                   // returns a slice of "Spans" representing all ranges in the set
//
// Operations (see also AddSet above, which performs a set union)
//
//  b := rangeset.Equal(s, s2)         // returns true if 2 sets are identical
//  s := rangeset.Union(s1, s2, ...)   // returns a new set that is the union of 1 or more sets
//  s := rangeset.Intersect(s1,s2,...) // returns a set that's the intersection of 1 or more sets
//  s2 := rangeset.Complement(s1)      // returns the inverse of a set
//
// Iterating (see also Values above, which returns a slice of all elements)
//
//  s.Iterate(f)                       // calls the function f on each element of the set
//  s.Filter(f)                        // call f on each element and deletes the element if f() returns false
//  s.Iterator(ctx)                    // returns a chan that is sent all the set's elements (in order)
//
package rangeset
