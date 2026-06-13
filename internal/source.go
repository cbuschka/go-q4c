package internal

import "iter"

type Source[E1 any] func() iter.Seq2[E1, error]

func SourceFromSlice[E any](elements []E) Source[E] {
	source := func() iter.Seq2[E, error] {
		return func(yield func(E, error) bool) {
			for _, v := range elements {
				if !yield(v, nil) {
					return
				}
			}
		}
	}
	return source
}

func (s Source[E]) FilteredBy(cond func(element E) bool) Source[E] {
	filteredSource := func() iter.Seq2[E, error] {
		return func(yield func(E, error) bool) {
			for e := range s() {
				if !cond(e) {
					continue
				}

				if !yield(e, nil) {
					return
				}
			}
		}
	}
	return filteredSource
}

func NewSourceMappedBy[E any, R any](source Source[E], mapper func(element E) (R, error)) Source[R] {
	mappedSource := func() iter.Seq2[R, error] {
		return func(yield func(R, error) bool) {
			var emptyR R
			for e, err := range source() {
				if err != nil {
					if !yield(emptyR, err) {
						return
					}
				}
				r, err := mapper(e)
				if err != nil {
					if !yield(emptyR, err) {
						return
					}
				}

				if !yield(r, nil) {
					return
				}
			}
		}
	}
	return mappedSource
}

func (s Source[E]) ToSlice() ([]E, error) {
	elements := make([]E, 0)
	for element, err := range s() {
		if err != nil {
			return nil, err
		}
		elements = append(elements, element)
	}
	return elements, nil
}

func (s Source[E]) First() (E, bool, error) {
	var empty E
	for element, err := range s() {
		if err != nil {
			return empty, false, err
		}
		return element, true, err
	}
	return empty, false, nil
}
