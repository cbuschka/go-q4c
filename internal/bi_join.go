package internal

import (
	"fmt"
	"iter"

	"github.com/cbuschka/go-q4c/types"
)

type JoinType int

const (
	Inner = iota
	Left
)

func join[E1 any, E2 any](joinType JoinType, source1 Source[E1], key1 types.KeyFunc[E1, interface{}],
	source2 Source[E2], key2 types.KeyFunc[E2, any]) Source[types.Pair[E1, E2]] {
	index, err := createIndex(source2, key2)
	if err != nil {
		return func() iter.Seq2[types.Pair[E1, E2], error] {
			return func(yield func(types.Pair[E1, E2], error) bool) {
				var empty types.Pair[E1, E2]
				yield(empty, err)
			}
		}
	}

	switch joinType {
	case Inner:
		return generateInnerJoinSeq(source1, key1, index)
	case Left:
		return generateLeftOuterJoinSeq(source1, key1, index)
	default:
		panic(fmt.Errorf("unknown join type"))
	}
}

func generateLeftOuterJoinSeq[E1 any, E2 any](source1 Source[E1], key1 types.KeyFunc[E1, interface{}],
	index map[interface{}][]E2) Source[types.Pair[E1, E2]] {

	return func() iter.Seq2[types.Pair[E1, E2], error] {
		return func(yield func(types.Pair[E1, E2], error) bool) {
			var empty types.Pair[E1, E2]
			for element1, err := range source1() {
				if err != nil {
					if !yield(empty, err) {
						break
					}
				}
				element1Key := key1(element1)
				elements, found := index[element1Key]
				if !found {
					var emptyE2 E2
					if !yield(types.Pair[E1, E2]{Element1: element1, Element2: emptyE2}, nil) {
						break
					}
				} else {
					for _, element2 := range elements {
						if !yield(types.Pair[E1, E2]{Element1: element1, Element2: element2}, nil) {
							break
						}
					}
				}
			}
		}
	}
}

func createIndex[E any](source2 Source[E], key2 types.KeyFunc[E, any]) (map[interface{}][]E, error) {
	index := make(map[interface{}][]E)
	for element2 := range source2() {
		element2Key := key2(element2)
		elements, found := index[element2Key]
		if !found {
			elements = make([]E, 0)
		}
		elements = append(elements, element2)
		index[element2Key] = elements
	}
	return index, nil
}

func generateInnerJoinSeq[E1 any, E2 any](source1 Source[E1], key1 types.KeyFunc[E1, interface{}], index map[interface{}][]E2) Source[types.Pair[E1, E2]] {
	return func() iter.Seq2[types.Pair[E1, E2], error] {
		return func(yield func(types.Pair[E1, E2], error) bool) {
			var empty types.Pair[E1, E2]
			for element1, err := range source1() {
				if err != nil {
					if !yield(empty, err) {
						break
					}
				}
				element1Key := key1(element1)
				elements, found := index[element1Key]
				if found {
					for _, element2 := range elements {
						if !yield(types.Pair[E1, E2]{Element1: element1, Element2: element2}, nil) {
							break
						}
					}
				}
			}
		}
	}
}
