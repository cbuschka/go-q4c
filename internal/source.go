package internal

import "iter"

type Source[E1 any] func() iter.Seq[E1]

func SourceFromSlice[E any](elements []E) Source[E] {
	source := func() iter.Seq[E] {
		return func(yield func(E) bool) {
			for _, v := range elements {
				if !yield(v) {
					return
				}
			}
		}
	}
	return source
}

func (s Source[E]) FilteredBy(cond func(element E) bool) Source[E] {
	filteredSource := func() iter.Seq[E] {
		return func(yield func(E) bool) {
			for e := range s() {
				if !cond(e) {
					continue
				}

				if !yield(e) {
					return
				}
			}
		}
	}
	return filteredSource
}

func NewSourceMappedBy[E any, R any](source Source[E], mapper func(element E) R) Source[R] {
	mappedSource := func() iter.Seq[R] {
		return func(yield func(R) bool) {
			for e := range source() {
				r := mapper(e)

				if !yield(r) {
					return
				}
			}
		}
	}
	return mappedSource
}

func (s Source[E]) ToSlice() []E {
	elements := make([]E, 0)
	for element := range s() {
		elements = append(elements, element)
	}
	return elements
}
