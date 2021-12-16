Rangeset implements a set container using the Go generics as implemented in the Go 1.18. It is somewhat similar to the
example "sets" package (new in Go m1.18)) but uses a slice of ranges rather than a map internally. This has advantages,
such as space-efficiency for some large sets, elements are returned in order, complement (inverse) operation is
supported as is Universal sets, etc.

## Go Generics (Parametric Polymorphism)

Parametric polymorphism, or what the Go community now commonly call **generics**, had been proposed for "Go 2" for some
time. It was officially added to the language in early 2022 (Go 1.18).

As in other languages, generics allow you to add "type parameters" to *functions* and to *types*. Unlike normal (value)
parameters, type parameters must be known at compile time. This rangeset package implements a generic *type* - a "set"
that has a type parameter specifying the type of elements that can be added to the set.

## Sets in Go

Set in Go (without generics) are typically implemented using a map, where the key is the set element type and the value
is ignored. This is efficient but not that useful as it does not provide many operations that are commonly used with
sets. Note that the element type of a set must be comparable (ie support the == and != operations). This is enforced
when using a map as a map key must be comparable.

When map is used as a makeshift set the value type is not important (only the key) so typically `struct{}` is used. (
Note that `struct{}` is a type and `struct{}{}` is a literal of that type.)  So for a set of ints you would
use `map[int]struct{}`. It's also not unusual to see a map of `bool` (eg `map[int]bool`) but that is somewhat redundant
as storing false would be the same as the element not being present at all. Moreover, for large sets it will save a lot
of space as `struct{}` has zero size whereas a `bool` uses one byte.

Alternatively, there are open source projects that implement sets, such as the
excellent https://github.com/deskarep/golang-set which provides all manner of set operations. However, a problem with
this type of solution is that set elements must be "boxed" (stored in an `interface{}`), which affects performance. It
also impacts type safety - for example, the compiler can't prevent you from accidentally adding a `string` to a set
of `int`. Moreover, you could even add a value of a non-comparable type to a set which would cause a run-time panic.

This sort of problem is where generics shine. It is easy to create a generic set type that is as performant and as
type-safe as a map, with the convenience and safety of pre-written set operations (and also hides the confusing use
of `struct{}` as the map value type). This is exactly what the example set type in the standard library set package
provides.

## rangeset

The rangeset package similarly implements a generic set but with a twist. It uses a slice of ranges to store the set,
instead of a map, which can be advantageous for sets with large contiguous ranges of elements. However, due to the use
of ranges the element type must be orderable (sets usually only require their elements to be comparable). That is, the
type of the element must support operations like less than (<, <=, >, >=) as well as incrementation (++). Hence the only
types in Go that can be used are integer types (byte, int, uint64, rune, etc).

I first had the idea for a "range set" at least 30 years ago and implemented a simple one in C. However, it wasn't until
templates were added to C++ that I created an efficient and useful implementation after I first started using the
ground-breaking STL in C++. My C++ range_set class was compatible with std::set of the STL (apart from the fact that the
type parameter had to be of an integral type). See the article I wrote on the class for the C/C++ User Journal in June
1999.

Although my C++ implementation used a linked list of ranges, I found that in Go a slice of ranges worked equally well.
Each range in the slice simply stores the bounds of the range using asymmetric bounds (inclusive lower bound, exclusive
upper bound). All operations maintain the property that the ranges are kept in numeric order and non-overlapping.

For compatibility with the example generic set mentioned above, that the Go Authors created, I have used the same method
names, including `Contains`, `AddSet`, `SubSet`, `Iterate`, etc, so rangeset could act as a drop-in replacement for a
set (of integers). Of course, there are additional methods that take advantage of the unique properties of a rangeset,
such as the ability to return all the ranges.

## Disadvantages

Apart from the obvious fact that you can't have a rangeset of `string`s (or any non-integer type) there are a few other
things to be aware of.

First, many set operations require a search of the slice, such as checking for the presence of an element or finding
where to add an element. These require a binary search which has time complexity of O(log r), where r is the number of
ranges in the set. In the worst case, where n/r == 1 (ie each element is in its own range) then the time complexity is
O(log n).

Hence time complexity is worse than that for a set implemented using a map (hash table) which has time complexity O(1).
That said, for sets with a small number of ranges the times become similar - ie, as n/r goes to infinity the time
complexity goes towards O(1). In fact, benchmarks show that lookups are faster for rangesets (with small number of
ranges) than for sets based on maps.

Space complexity is much better, being O(r). Even in the worst case (n/r == 1) it is O(n) the same as a map. In
practice, for the worst case (all non-contiguous elements) a rangeset will be almost twice the size as the map with the
same elements, since each range stores two values, but that is not the sort of set to use with a rangeset.

Perhaps not a disadvantage, but another thing to be aware of is that *a rangeset is not safe for concurrent access*. If
you are accessing a rangeset from more than one goroutine at once, you must protect the access (unless all accesses are
reads) - for example with a `mutex`.

## Advantages

The obvious advantage is that for sets with a large number of elements, if all the elements fall into one or a small
number of contiguous ranges there can be space savings. You can store huge sets in a small amount of memory and adding
elements can even *reduce* memory requirements, if it results in two existing ranges being joined.

Of course, memory is cheap (and most Go programs run in environments with huge amounts of memory) so here are a few more
benefits.

Sometimes, it is very useful to be able to get the set elements in order. Methods, such as `Values` and `Iterate` return
the element values in order. In contrast, with a set implemented using a map (hash table) the order that values are
returned is indeterminate, and you would need to store and sort them yourself.

Another useful feature (that I found later was useful in the C++ rangeset) is to be able to get all the ranges of the
set. This is possible using the `Spans` method.  (A pair of values representing a range that is part of a rangeset is
stored in a `Span struct`.)

Finally, the `String` method and `NewFromString` function allow you to easily encode and decode sets for storage. These
are also useful for displaying/obtaining a set to/from the user.

## Uses

A rangeset may not be appropriate for every use of a set, especially sparse sets, but over many years I have found a
surprising number of uses for the C++ version. For example, it is used in several places in my open-source hex editor (
see https://github.com/AndrewWPhillips/HexEdit).

The first use I made of the C++ rangeset was in an implementation of a Windows "virtual" list control. It allowed for a
list control with up to 4 billion (virtual) items. (The list box that Windows provided, at the time, had trouble
handling 1,000.)

Using this list control, a user could select large swathes of elements (by clicking, scrolling and Shift+clicking) which
would simply be stored as a range in the range set. Selecting the whole list (with Ctrl+A) resulted in a rangeset of
just one range. Using a "normal" set (such as std:set) it would have been impossible to store the selection for a list
of many millions of elements (given the memory sizes of 20 years ago).

## Types

There are three exported types: `Set` is the range set, `Element` constrains the `Set`s type parameters to only be of
integer types, `Span` stores two values representing a range (as in the slice returned by the `Spans` method).

Normally, you would just use the `Set` type by creating one something like this:

```
    set := rangeset.Make[uint16](1, 42, 3e4)
```

## Methods

Type `Set` implements these methods:

`Add` adds an element to a set (returns true if added, false if already present)

`AddRange` adds a range of elements to the set

`Delete` removes an element from the set

`DeleteRange` removes a range of elements

`Contains` returns true if the element is in the set

`Len` returns the number of elements as an `int` (may wrap around if larger than maxInt)

`Length` returns the number of elements as `uint64` and the number of ranges

`Values` returns a slice containing all the elements (in numeric order)

`Spans` returns a slice of `Span`s containing all the ranges in the set

`String` returns a string encoding of a rangeset

`Copy` returns a copy of a set

`AddSet` adds all the elements of another set (Union)

`SubSet` deletes all the elements of another set

`Intersect` deletes all elements *not* in another set

`Complement` returns the inverse set

`Iterate` calls a function on every element of a set (in numeric order)

`Filter` deletes every element on which a boolean function fails

`Iterator` returns a <-chan on which every element in the set is placed (in order)

`ReadAll` adds all of the elements read from a <-chan (inverse of Iterator)

## Functions

`Make` creates a new set (optionally with initial elements)

`Universal` returns a set including all elements allowed by the type parameter

`NewFromString` returns a new set from a string encoded with the `String` method (above)

`NewFromRange` returns a new set given an asymmetric range of values

`Equal` compares two sets

`Union` finds the union of one or more sets

`Intersect` finds the intersection of one or more sets

## Acknowledgements

Thanks to Robert Greisemer for providing the generic `minInt` function

Thanks to Dave Cheney for many ideas such as using a map for table-driven tests

Thanks to Bill Kennedy for ideas on more readable messages in tests

Thanks to Steve Mann for the idea of using a colon (:) rather than a dash (-) in the string encoding
