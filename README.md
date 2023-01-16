Lazy Fingertree in Go based on [Qiao's JavaScript version](https://github.com/qiao/fingertree.js/) of [Ralf Hinze's and Ross Paterson's Haskell version](http://www.soi.city.ac.uk/~ross/papers/FingerTree.html)

The public API is generic but the implementation is currently not generic because of some type resolution problems. In the future, if it's actually possible, the parameters could be pushed down into the implementation to make the code cleaner.

[adapters.go](adapters.go) contains the public API:

You provide your own object that supports the Measurer[Value, Measurement] interface. `Values` are in the leaves of the tree and your `Measurer` computes the `Measurements` in the `Measure()` and `Sum()` methods. `Measurements` can be any go objects but they *should be immutable* or there could be trouble. Please see [Ralf Hinze's and Ross Paterson's finger tree paper](http://www.soi.city.ac.uk/~ross/papers/FingerTree.html) (and the test code) for more information.

Here's the go doc:

# lazyfingertree
--
    import "github.com/zot/lazyfingertree"

Package lazyfingertree implements lazy finger trees. See the [fingertree paper]
for details.

[fingertree paper]: http://www.soi.city.ac.uk/~ross/papers/FingerTree.html

## Usage

```go
var ErrBadMeasurer = fmt.Errorf("%w, bad measurer", ErrFingerTree)
```

```go
var ErrBadValue = fmt.Errorf("%w, bad value", ErrFingerTree)
```

```go
var ErrEmptyTree = fmt.Errorf("%w, empty tree", ErrFingerTree)
```

```go
var ErrExpectedNode = fmt.Errorf("%w, expected a node", ErrFingerTree)
```

```go
var ErrFingerTree = errors.New("finger tree")
```

```go
var ErrUnsupported = fmt.Errorf("%w, unsupported operation", ErrFingerTree)
```

#### func  Diag

```go
func Diag(v any) string
```
Return a string that represents a value. Diag calls the diagstr method if the
value implements it:

    diagstr() string

otherwise, it calls fmt.Sprintf("%v", v)

#### type FingerTree

```go
type FingerTree[MS Measurer[V, M], V, M any] struct {
}
```

FingerTree is a parameterized wrapper on a low-level finger tree.

#### func  FromArray

```go
func FromArray[MS Measurer[V, M], V, M any](measurer MS, values []V) FingerTree[MS, V, M]
```
Create a finger tree. You shouldn't need to provide the type parameters, Go
should be able to infer them from your arguments. So you should just be able to
say,

    t := FromArray(myMeasurer, []Plant{plant1, plant2})

#### func (FingerTree[MS, V, M]) AddFirst

```go
func (t FingerTree[MS, V, M]) AddFirst(value any) FingerTree[MS, V, M]
```
Add a value to the start of the tree.

#### func (FingerTree[MS, V, M]) AddLast

```go
func (t FingerTree[MS, V, M]) AddLast(value any) FingerTree[MS, V, M]
```
Add a value to the and of the tree.

#### func (FingerTree[MS, V, M]) Concat

```go
func (t FingerTree[MS, V, M]) Concat(other FingerTree[MS, V, M]) FingerTree[MS, V, M]
```
Join two finger trees together

#### func (FingerTree[MS, V, M]) DropUntil

```go
func (t FingerTree[MS, V, M]) DropUntil(pred Predicate[M]) FingerTree[MS, V, M]
```
Discard all the initial values in the tree that do not satisfy the predicate

#### func (FingerTree[MS, V, M]) Each

```go
func (t FingerTree[MS, V, M]) Each(iter IterFunc[V])
```
Iterate through the tree starting at the beginning

#### func (FingerTree[MS, V, M]) EachReverse

```go
func (t FingerTree[MS, V, M]) EachReverse(iter IterFunc[V])
```
Iterate through the tree starting at the end

#### func (FingerTree[MS, V, M]) IsEmpty

```go
func (t FingerTree[MS, V, M]) IsEmpty() bool
```
Return whether the tree is empty

#### func (FingerTree[MS, V, M]) Measure

```go
func (t FingerTree[MS, V, M]) Measure() M
```
Return the measure of all the tree's values

#### func (FingerTree[MS, V, M]) PeekFirst

```go
func (t FingerTree[MS, V, M]) PeekFirst() V
```
Return the first value in the tree. Make sure to test whether the tree is empty
because this will panic if it is.

#### func (FingerTree[MS, V, M]) PeekLast

```go
func (t FingerTree[MS, V, M]) PeekLast() V
```
Return the last value in the tree. Make sure to test whether the tree is empty
because this will panic if it is.

#### func (FingerTree[MS, V, M]) RemoveFirst

```go
func (t FingerTree[MS, V, M]) RemoveFirst() FingerTree[MS, V, M]
```
Remove the first value in the tree. Make sure to test whether the tree is empty
because this will panic if it is.

#### func (FingerTree[MS, V, M]) RemoveLast

```go
func (t FingerTree[MS, V, M]) RemoveLast() FingerTree[MS, V, M]
```
Remove the last value in the tree. Make sure to test whether the tree is empty
because this will panic if it is.

#### func (FingerTree[MS, V, M]) Split

```go
func (t FingerTree[MS, V, M]) Split(predicate Predicate[M]) (FingerTree[MS, V, M], FingerTree[MS, V, M])
```
Split the tree. The first tree is all the starting values that do not satisfy
the predicate. The second tree is the first value that satisfies the predicate,
followed by the rest of the values.

#### func (FingerTree[MS, V, M]) TakeUntil

```go
func (t FingerTree[MS, V, M]) TakeUntil(pred Predicate[M]) FingerTree[MS, V, M]
```
Return all the initial values in the tree that do not satisfy the predicate

#### func (FingerTree[MS, V, M]) ToSlice

```go
func (t FingerTree[MS, V, M]) ToSlice() []V
```
Return a slice containing all of the values in the tree

#### type IterFunc

```go
type IterFunc[V any] func(value V) bool
```

An IterFunc is a function that takes a value and returns true or false. It's
used by [Each] and [EachReverse]. Returning true means to continue iteration.
Returning false means to stop.

#### type Measurer

```go
type Measurer[Value, Measure any] interface {
	// The "zero" measure
	Identity() Measure
	// Return the measure for a value.
	// Measuring a value could technically produce an error but really should not.
	// Make sure to validate inputs or to use a panic if you need error support.
	Measure(value Value) Measure
	// Add two measures together
	Sum(a Measure, b Measure) Measure
}
```

The measurer interface

#### type Predicate

```go
type Predicate[M any] func(measure M) bool
```

A Predicate is a function that takes a measure and returns true or false. It's
used by [Split], [TakeUntil], and [DropUntil].
