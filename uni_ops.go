package q4c

import "iter"

type UniSelect[E1 any] interface {
	Stream() iter.Seq[E1]
	ToSlice() []E1
}

type uniSelectImpl[E1 any] struct {
	source func() iter.Seq[E1]
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

func SelectFrom[E1 any](elements []E1) UniSelect[E1] {
	source := func() iter.Seq[E1] {
		return func(yield func(E1) bool) {
			for _, v := range elements {
				if !yield(v) {
					return
				}
			}
		}
	}

	return UniSelect[E1](&uniSelectImpl[E1]{source})
}
