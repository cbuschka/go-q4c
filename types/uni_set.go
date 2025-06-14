package types

import "iter"

type UniFilterCondition[E1 any] func(e E1) bool

type UniSelect[E1 any] interface {
	Stream() iter.Seq[E1]
	ToSlice() []E1
	First() (element E1, found bool)
}

type FilterableUniSelect[E1 any] interface {
	UniSelect[E1]
	Where(cond UniFilterCondition[E1]) UniSelect[E1]
}

type UniSet[E1 any] interface {
	SelectFrom(elements []E1) FilterableUniSelect[E1]
}
