// Package lazyfingertree implements parameterized lazy finger trees.
// See the [readme](README.md) for details.
package lazyfingertree

import (
	"fmt"
	"io"
)

// A Predicate is a function that takes a measure and returns true or false.
// It's used by [Split], [TakeUntil], and [DropUntil].
type Predicate[M any] func(measure M) bool

// An IterFunc is a function that takes a value and returns true or false.
// It's used by [Each] and [EachReverse]. Returning true means to continue
// iteration. Returning false means to stop.
type IterFunc[V any] func(value V) bool

// FingerTree is a parameterized wrapper on a low-level finger tree.
type FingerTree[MS Measurer[Value, Measure], Value, Measure any] struct {
	f fingerTree
}

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

func wrapTree[MS Measurer[V, M], V, M any](tree fingerTree) FingerTree[MS, V, M] {
	return FingerTree[MS, V, M]{tree}
}

var ErrBadValue = fmt.Errorf("%w, bad value", ErrFingerTree)

func wrapPredicate[M any](pred Predicate[M]) func(any) bool {
	return func(m any) bool {
		if wm, ok := m.(M); !ok {
			panic(fmt.Errorf("%w, predicate input: %v", ErrBadValue, m))
		} else {
			return pred(wm)
		}
	}
}

func wrapIter[V any](iter IterFunc[V]) func(any) bool {
	return func(v any) bool {
		if wv, ok := v.(V); !ok {
			panic(fmt.Errorf("%w, iteration value: %v", ErrBadValue, v))
		} else {
			return iter(wv)
		}
	}
}

func null[T any]() T {
	return *new(T)
}

// Add a value to the start of the tree.
func (t FingerTree[MS, V, M]) AddFirst(value V) FingerTree[MS, V, M] {
	return wrapTree[MS, V, M](t.f.AddFirst(value))
}

// Add a value to the and of the tree.
func (t FingerTree[MS, V, M]) AddLast(value V) FingerTree[MS, V, M] {
	return wrapTree[MS, V, M](t.f.AddLast(value))
}

// Remove the first value in the tree. Make sure to test whether the tree is empty
// because this will panic if it is.
func (t FingerTree[MS, V, M]) RemoveFirst() FingerTree[MS, V, M] {
	return wrapTree[MS, V, M](t.f.RemoveFirst())
}

// Remove the last value in the tree. Make sure to test whether the tree is empty
// because this will panic if it is.
func (t FingerTree[MS, V, M]) RemoveLast() FingerTree[MS, V, M] {
	return wrapTree[MS, V, M](t.f.RemoveLast())
}

// Return the first value in the tree. Make sure to test whether the tree is empty
// because this will panic if it is.
func (t FingerTree[MS, V, M]) PeekFirst() V {
	if cv, ok := t.f.PeekFirst().(V); !ok {
		panic(fmt.Errorf("%w, first value in tree: %v", ErrBadValue, t.f.PeekFirst()))
	} else {
		return cv
	}
}

// Return the last value in the tree. Make sure to test whether the tree is empty
// because this will panic if it is.
func (t FingerTree[MS, V, M]) PeekLast() V {
	if cv, ok := t.f.PeekLast().(V); !ok {
		panic(fmt.Errorf("%w, last value in tree: %v", ErrBadValue, t.f.PeekLast()))
	} else {
		return cv
	}
}

// Join two finger trees together
func (t FingerTree[MS, V, M]) Concat(other FingerTree[MS, V, M]) FingerTree[MS, V, M] {
	return wrapTree[MS, V, M](t.f.Concat(other.f))
}

// Split the tree. The first tree is all the starting values that do not satisfy the predicate.
// The second tree is the first value that satisfies the predicate, followed by the rest of the values.
func (t FingerTree[MS, V, M]) Split(predicate Predicate[M]) (FingerTree[MS, V, M], FingerTree[MS, V, M]) {
	left, right := t.f.Split(wrapPredicate(predicate))
	return wrapTree[MS, V, M](left), wrapTree[MS, V, M](right)
}

// Return a slice containing all of the values in the tree
func (t FingerTree[MS, V, M]) ToSlice() []V {
	s := t.f.ToSlice()
	result := make([]V, len(s))
	for i := 0; i < len(s); i++ {
		result[i] = s[i].(V)
	}
	return result
}

func (t FingerTree[MS, V, M]) IsZero() bool {
	return t.f == nil
}

func (t FingerTree[MS, V, M]) String() string {
	return t.f.String()
}

func (t FingerTree[MS, V, M]) Dump(w io.Writer, level int) {
	t.f.Dump(w, level)
}

// Return whether the tree is empty
func (t FingerTree[MS, V, M]) IsEmpty() bool {
	return isEmpty(t.f)
}

// Return the measure of all the tree's values
func (t FingerTree[MS, V, M]) Measure() M {
	if cm, ok := t.f.measurement().value.(M); !ok {
		panic(fmt.Errorf("%w, measurement in tree: %v", ErrBadValue, t.f.measurement().value))
	} else {
		return cm
	}
}

// Return all the initial values in the tree that do not satisfy the predicate
func (t FingerTree[MS, V, M]) TakeUntil(pred Predicate[M]) FingerTree[MS, V, M] {
	return wrapTree[MS, V, M](takeUntil(t.f, wrapPredicate(pred)))
}

// Discard all the initial values in the tree that do not satisfy the predicate
func (t FingerTree[MS, V, M]) DropUntil(pred Predicate[M]) FingerTree[MS, V, M] {
	return wrapTree[MS, V, M](dropUntil(t.f, wrapPredicate(pred)))
}

// Iterate through the tree starting at the beginning
func (t FingerTree[MS, V, M]) Each(iter IterFunc[V]) {
	t.f.Each(wrapIter(iter))
}

// Iterate through the tree starting at the end
func (t FingerTree[MS, V, M]) EachReverse(iter IterFunc[V]) {
	t.f.EachReverse(wrapIter(iter))
}

// The measurer interface
func AsMeasurer[V, M any](m any) Measurer[V, M] {
	if meas, ok := m.(Measurer[V, M]); !ok {
		panic(fmt.Errorf("%w, measurer: %v", ErrBadValue, m))
	} else {
		return meas
	}
}

type adaptedMeasurer[MS Measurer[V, M], V, M any] struct {
	am MS
}

func (m adaptedMeasurer[MS, V, M]) Identity() any {
	return m.am.Identity()
}

func (m adaptedMeasurer[MS, V, M]) Measure(value any) any {
	if v, ok := value.(V); !ok {
		panic(fmt.Errorf("%w, measure: %v", ErrBadValue, value))
	} else {
		return m.am.Measure(v)
	}
}

func (m adaptedMeasurer[MS, V, M]) Sum(a any, b any) any {
	if va, ok := a.(M); !ok {
		panic(fmt.Errorf("%w, sum: %v", ErrBadValue, a))
	} else if vb, ok2 := b.(M); !ok2 {
		panic(fmt.Errorf("%w, sum: %v", ErrBadValue, b))
	} else {
		return m.am.Sum(va, vb)
	}
}

// Create a finger tree. You shouldn't need to provide the type parameters,
// Go should be able to infer them from your arguments.
// So you should just be able to say,
//
//	t := FromArray(myMeasurer, []Plant{plant1, plant2})
func FromArray[MS Measurer[V, M], V, M any](measurer MS, values []V) FingerTree[MS, V, M] {
	cvt := make([]any, len(values))
	for i := 0; i < len(values); i++ {
		cvt[i] = values[i]
	}
	return wrapTree[MS, V, M](fromArray(adaptedMeasurer[MS, V, M]{measurer}, cvt))
}

func Concat[MS Measurer[V, M], V, M any](trees ...FingerTree[MS, V, M]) FingerTree[MS, V, M] {
	result := newEmptyTree(trees[0].f.measurement().measurer)
	for _, t := range trees {
		result = result.Concat(t.f)
	}
	return wrapTree[MS, V, M](result)
}
