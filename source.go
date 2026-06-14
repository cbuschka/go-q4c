package q4c

import "iter"

type source[E1 any] func() iter.Seq2[E1, error]

func newSourceFromSlice[E any](elements []E) source[E] {
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

func (s source[E]) filteredBy(cond func(element E) bool) source[E] {
	filteredSource := func() iter.Seq2[E, error] {
		return func(yield func(E, error) bool) {
			var empty E
			for e, err := range s() {
				if err != nil {
					yield(empty, err)
					return
				}
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

func (s source[E]) toSlice() ([]E, error) {
	elements := make([]E, 0)
	for element, err := range s() {
		if err != nil {
			return nil, err
		}
		elements = append(elements, element)
	}
	return elements, nil
}

func (s source[E]) first() (E, bool, error) {
	var empty E
	for element, err := range s() {
		if err != nil {
			return empty, false, err
		}
		return element, true, err
	}
	return empty, false, nil
}
