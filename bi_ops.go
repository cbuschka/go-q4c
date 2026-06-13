package q4c

import (
	"iter"
)

type BiSet[E1 any, E2 any] struct {
}

type JoinableBiSelect[E1 any, E2 any] struct {
	source source[E1]
}

type BiFilterCondition[E1 any, E2 any] func(e E1, e2 E2) bool

func NewBiSet[E1 any, E2 any]() *BiSet[E1, E2] {
	return &BiSet[E1, E2]{}
}

func (s *BiSet[E1, E2]) SelectFrom(elements []E1) *JoinableBiSelect[E1, E2] {
	source := newSourceFromSlice[E1](elements)
	return &JoinableBiSelect[E1, E2]{source: source}
}

func (s *JoinableBiSelect[E1, E2]) Stream() iter.Seq2[E1, error] {
	return s.source()
}

func (s *JoinableBiSelect[E1, E2]) ToSlice() ([]E1, error) {
	return s.source.toSlice()
}

func (s *JoinableBiSelect[E1, E2]) First() (E1, bool, error) {
	return s.source.first()
}

type JoiningBiSet[E1 any, E2 any] struct {
	source   source[E1]
	source2  source[E2]
	joinType joinType
}

type FilterableBiSelect[E1 any, E2 any] struct {
	source source[Pair[E1, E2]]
}

type FilteredBiSelect[E1 any, E2 any] struct {
	source source[Pair[E1, E2]]
}

func (s *JoinableBiSelect[E1, E2]) LeftOuterJoin(elements []E2) *JoiningBiSet[E1, E2] {
	source2 := newSourceFromSlice[E2](elements)
	return &JoiningBiSet[E1, E2]{source: s.source, source2: source2, joinType: leftJoin}
}

func (s *JoinableBiSelect[E1, E2]) Join(elements []E2) *JoiningBiSet[E1, E2] {
	source2 := newSourceFromSlice[E2](elements)
	return &JoiningBiSet[E1, E2]{source: s.source, source2: source2, joinType: innerJoin}
}

func (s *JoiningBiSet[E1, E2]) On(key1 KeyFunc[E1, any], key2 KeyFunc[E2, any]) *FilterableBiSelect[E1, E2] {
	var s2 source[Pair[E1, E2]]
	s2 = join(s.joinType, s.source, key1, s.source2, key2)
	impl := FilterableBiSelect[E1, E2]{source: s2}
	return &impl
}

func (s *FilterableBiSelect[E1, E2]) Where(cond BiFilterCondition[E1, E2]) *FilteredBiSelect[E1, E2] {
	filteredSource := s.source.filteredBy(func(pair Pair[E1, E2]) bool {
		return cond(pair.Element1, pair.Element2)
	})
	return &FilteredBiSelect[E1, E2]{source: filteredSource}
}

func (s *FilterableBiSelect[E1, E2]) Stream() iter.Seq2[Pair[E1, E2], error] {
	return s.source()
}

func (s *FilterableBiSelect[E1, E2]) ToSlice() ([]Pair[E1, E2], error) {
	return s.source.toSlice()
}

func (s *FilterableBiSelect[E1, E2]) First() (Pair[E1, E2], bool, error) {
	return s.source.first()
}

func (s *FilteredBiSelect[E1, E2]) ToSlice() ([]Pair[E1, E2], error) {
	return s.source.toSlice()
}

func (s *FilteredBiSelect[E1, E2]) First() (Pair[E1, E2], bool, error) {
	return s.source.first()
}

func (s *FilteredBiSelect[E1, E2]) Stream() iter.Seq2[Pair[E1, E2], error] {
	return s.source()
}
