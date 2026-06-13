package internal

import (
	"iter"

	"github.com/cbuschka/go-q4c/types"
)

type BiSetImpl[E1 any, E2 any] struct {
	source Source[E1]
}

func (s *BiSetImpl[E1, E2]) SelectFrom(elements []E1) types.JoinableBiSelect[E1, E2] {
	s.source = SourceFromSlice[E1](elements)
	return types.JoinableBiSelect[E1, E2](s)
}

func (s *BiSetImpl[E1, E2]) Stream() iter.Seq2[E1, error] {
	return s.source()
}

func (s *BiSetImpl[E1, E2]) ToSlice() ([]E1, error) {
	return s.source.ToSlice()
}

func (s *BiSetImpl[E1, E2]) First() (E1, bool, error) {
	return s.source.First()
}

type joiningBiSetImpl[E1 any, E2 any] struct {
	source   Source[E1]
	source2  Source[E2]
	joinType JoinType
}

type joinedBiSetImpl[E1 any, E2 any] struct {
	source Source[types.Pair[E1, E2]]
}

func (s *BiSetImpl[E1, E2]) LeftOuterJoin(elements []E2) types.BiSelectJoin[E1, E2] {
	source2 := SourceFromSlice[E2](elements)
	return types.BiSelectJoin[E1, E2](&joiningBiSetImpl[E1, E2]{source: s.source, source2: source2, joinType: Left})
}

func (s *BiSetImpl[E1, E2]) Join(elements []E2) types.BiSelectJoin[E1, E2] {
	source2 := SourceFromSlice[E2](elements)
	return types.BiSelectJoin[E1, E2](&joiningBiSetImpl[E1, E2]{source: s.source, source2: source2, joinType: Inner})
}

func (s *joiningBiSetImpl[E1, E2]) On(key1 types.KeyFunc[E1, any], key2 types.KeyFunc[E2, any]) types.FilterableBiSelect[E1, E2] {
	var s2 Source[types.Pair[E1, E2]]
	s2 = join(s.joinType, s.source, key1, s.source2, key2)
	impl := joinedBiSetImpl[E1, E2]{source: s2}
	return types.FilterableBiSelect[E1, E2](&impl)
}

func (s *joinedBiSetImpl[E1, E2]) Where(cond types.BiFilterCondition[E1, E2]) types.BiSelect[E1, E2] {
	filteredSource := s.source.FilteredBy(func(pair types.Pair[E1, E2]) bool {
		return cond(pair.Element1, pair.Element2)
	})
	return types.BiSelect[E1, E2](&joinedBiSetImpl[E1, E2]{source: filteredSource})
}

func (s *joinedBiSetImpl[E1, E2]) Stream() iter.Seq2[types.Pair[E1, E2], error] {
	var empty types.Pair[E1, E2]
	return func(yield func(types.Pair[E1, E2], error) bool) {
		for pair, err := range s.source() {
			if err != nil {
				if !yield(empty, err) {
					break
				}
			}
			if !yield(pair, nil) {
				break
			}
		}
	}
}

func (s *joinedBiSetImpl[E1, E2]) ToSlice() ([]types.Pair[E1, E2], error) {
	return s.source.ToSlice()
}

func (s *joinedBiSetImpl[E1, E2]) First() (types.Pair[E1, E2], bool, error) {
	return s.source.First()
}
