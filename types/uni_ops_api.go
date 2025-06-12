package types

import "iter"

type UniFilterCondition[E1 any] func(e E1) bool

type UniSelect[E1 any] interface {
	Stream() iter.Seq[E1]
	ToSlice() []E1
}

type FilterableUniSelect[E1 any] interface {
	UniSelect[E1]
	Where(cond UniFilterCondition[E1]) UniSelect[E1]
}
