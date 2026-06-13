package types

import "iter"

type UniFilterCondition[E1 any] func(e E1) bool

type UniSelect[E1 any] interface {
	Stream() iter.Seq2[E1, error]
	ToSlice() ([]E1, error)
	First() (E1, bool, error)
}

type FilterableUniSelect[E1 any] interface {
	UniSelect[E1]
	Where(cond UniFilterCondition[E1]) UniSelect[E1]
}

type UniSet[E1 any] interface {
	SelectFrom(elements []E1) FilterableUniSelect[E1]
}
