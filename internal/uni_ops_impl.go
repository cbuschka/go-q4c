package internal

import (
	"github.com/cbuschka/go-q4c/types"
	"iter"
)

type uniSelectImpl[E1 any] struct {
	source func() iter.Seq[E1]
}

func (u *uniSelectImpl[E1]) Where(cond types.UniFilterCondition[E1]) types.UniSelect[E1] {
	filteredSource := func() iter.Seq[E1] {
		return func(yield func(E1) bool) {
			for e := range u.source() {
				if !cond(e) {
					continue
				}

				if !yield(e) {
					return
				}
			}
		}
	}

	return types.UniSelect[E1](&uniSelectImpl[E1]{filteredSource})
}

func (u *uniSelectImpl[E1]) Stream() iter.Seq[E1] {
	return u.source()
}

func (u *uniSelectImpl[E1]) ToSlice() []E1 {
	elements := make([]E1, 0)
	for e := range u.Stream() {
		elements = append(elements, e)
	}

	return elements
}

func SelectFrom[E1 any](elements []E1) types.FilterableUniSelect[E1] {
	source := func() iter.Seq[E1] {
		return func(yield func(E1) bool) {
			for _, v := range elements {
				if !yield(v) {
					return
				}
			}
		}
	}

	return types.FilterableUniSelect[E1](&uniSelectImpl[E1]{source})
}
