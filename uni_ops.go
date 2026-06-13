package q4c

import (
	"iter"
)

type UniFilterCondition[E1 any] func(e E1) bool

func NewUniSet[E1 any]() *UniSet[E1] {
	return &UniSet[E1]{}
}

func SelectFrom[E1 any](elements []E1) *FilterableUniSelect[E1] {
	return NewUniSet[E1]().SelectFrom(elements)
}

type UniSet[E1 any] struct {
}

type FilterableUniSelect[E1 any] struct {
	source source[E1]
}

type FilteredUniSelect[E1 any] struct {
	source source[E1]
}

func (u *FilterableUniSelect[E1]) Where(cond UniFilterCondition[E1]) *FilteredUniSelect[E1] {
	filteredSource := u.source.filteredBy(cond)

	return &FilteredUniSelect[E1]{filteredSource}
}

func (u *FilterableUniSelect[E1]) Stream() iter.Seq2[E1, error] {
	return u.source()
}

func (u *FilterableUniSelect[E1]) ToSlice() ([]E1, error) {
	return u.source.toSlice()
}

func (u *FilterableUniSelect[E1]) First() (element E1, found bool, err error) {
	return u.source.first()
}

func (u *UniSet[E1]) SelectFrom(elements []E1) *FilterableUniSelect[E1] {
	s := newSourceFromSlice(elements)
	return &FilterableUniSelect[E1]{source: s}
}

func (u *FilteredUniSelect[E1]) Stream() iter.Seq2[E1, error] {
	return u.source()
}

func (u *FilteredUniSelect[E1]) ToSlice() ([]E1, error) {
	return u.source.toSlice()
}

func (u *FilteredUniSelect[E1]) First() (element E1, found bool, err error) {
	return u.source.first()
}
