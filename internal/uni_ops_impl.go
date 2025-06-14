package internal

import (
	"github.com/cbuschka/go-q4c/types"
	"iter"
)

type UniSetImpl[E1 any] struct {
}

type uniSelectImpl[E1 any] struct {
	source Source[E1]
}

func (u *uniSelectImpl[E1]) Where(cond types.UniFilterCondition[E1]) types.UniSelect[E1] {
	filteredSource := u.source.FilteredBy(cond)

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

func (u *uniSelectImpl[E1]) First() (element E1, found bool) {
	slice := u.ToSlice()
	if len(slice) == 0 {
		return *new(E1), false
	}
	return slice[0], true
}

func (u *UniSetImpl[E1]) SelectFrom(elements []E1) types.FilterableUniSelect[E1] {
	source := SourceFromSlice(elements)

	return types.FilterableUniSelect[E1](&uniSelectImpl[E1]{source})
}
