package types

import "iter"

type BiFilterCondition[E1 any, E2 any] func(e E1, e2 E2) bool

type BiSet[E1 any, E2 any] interface {
	SelectFrom(elements []E1) JoinableBiSelect[E1, E2]
}

type JoinableBiSelect[E1 any, E2 any] interface {
	UniSelect[E1]
	Join(elements []E2) BiSelectJoin[E1, E2]
}

type BiSelectJoin[E1 any, E2 any] interface {
	On(key1 KeyFunc[E1, any], key2 KeyFunc[E2, any]) FilterableBiSelect[E1, E2]
}

type FilterableBiSelect[E1 any, E2 any] interface {
	BiSelect[E1, E2]
	Where(cond BiFilterCondition[E1, E2]) BiSelect[E1, E2]
}

type BiSelect[E1 any, E2 any] interface {
	Stream() iter.Seq[Pair[E1, E2]]
	ToSlice() []Pair[E1, E2]
}
