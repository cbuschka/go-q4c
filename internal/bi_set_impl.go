package internal

import (
	"github.com/cbuschka/go-q4c/types"
	"iter"
)

type BiSetImpl[E1 any, E2 any] struct {
	source Source[E1]
}

func (s *BiSetImpl[E1, E2]) SelectFrom(elements []E1) types.JoinableBiSelect[E1, E2] {
	s.source = SourceFromSlice[E1](elements)
	return types.JoinableBiSelect[E1, E2](s)
}

func (s *BiSetImpl[E1, E2]) Stream() iter.Seq[E1] {
	return s.source()
}

func (s *BiSetImpl[E1, E2]) ToSlice() []E1 {
	return s.source.ToSlice()
}

type completeBiSetImpl[E1 any, E2 any] struct {
	source  Source[E1]
	source2 Source[E2]
}

func (s *BiSetImpl[E1, E2]) Join(elements []E2) types.BiSelectJoin[E1, E2] {
	source2 := SourceFromSlice[E2](elements)
	return types.BiSelectJoin[E1, E2](&completeBiSetImpl[E1, E2]{s.source, source2})
}

func (s *completeBiSetImpl[E1, E2]) On(key1 types.KeyFunc[E1, any], key2 types.KeyFunc[E2, any]) types.FilterableBiSelect[E1, E2] {
	// FIXME create a join iter
	return types.FilterableBiSelect[E1, E2](s)
}

func (s *completeBiSetImpl[E1, E2]) Where(cond types.BiFilterCondition[E1, E2]) types.BiSelect[E1, E2] {
	return types.BiSelect[E1, E2](s)
}

func (s *completeBiSetImpl[E1, E2]) Stream() iter.Seq[types.Pair[E1, E2]] {
	return NewSourceMappedBy(s.source, func(e E1) types.Pair[E1, E2] {
		return types.Pair[E1, E2]{e, *new(E2)}
	})()
}

func (s *completeBiSetImpl[E1, E2]) ToSlice() []types.Pair[E1, E2] {
	return NewSourceMappedBy(s.source, func(e E1) types.Pair[E1, E2] {
		return types.Pair[E1, E2]{e, *new(E2)}
	}).ToSlice()
}
