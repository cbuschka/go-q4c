package internal

import (
	"iter"

	"github.com/cbuschka/go-q4c/types"
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

func (u *uniSelectImpl[E1]) Stream() iter.Seq2[E1, error] {
	return u.source()
}

func (u *uniSelectImpl[E1]) ToSlice() ([]E1, error) {
	return u.source.ToSlice()
}

func (u *uniSelectImpl[E1]) First() (element E1, found bool, err error) {
	var empty E1
	slice, err := u.ToSlice()
	if err != nil {
		return empty, false, err
	}
	if len(slice) == 0 {
		return empty, false, nil
	}
	return slice[0], true, nil
}

func (u *UniSetImpl[E1]) SelectFrom(elements []E1) types.FilterableUniSelect[E1] {
	source := SourceFromSlice(elements)

	return types.FilterableUniSelect[E1](&uniSelectImpl[E1]{source})
}
