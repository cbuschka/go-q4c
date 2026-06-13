package q4c

import (
	"fmt"
	"iter"
)

type KeyFunc[E any, K comparable] func(element1 E) K

type Pair[E1 any, E2 any] struct {
	Element1 E1
	Element2 E2
}

type joinType int

const (
	innerJoin = iota
	leftJoin
)

func join[E1 any, E2 any](joinType joinType, source1 source[E1], key1 KeyFunc[E1, interface{}],
	source2 source[E2], key2 KeyFunc[E2, any]) source[Pair[E1, E2]] {
	index, err := createIndex(source2, key2)
	if err != nil {
		return func() iter.Seq2[Pair[E1, E2], error] {
			return func(yield func(Pair[E1, E2], error) bool) {
				var empty Pair[E1, E2]
				yield(empty, err)
			}
		}
	}

	switch joinType {
	case innerJoin:
		return generateInnerJoinSeq(source1, key1, index)
	case leftJoin:
		return generateLeftOuterJoinSeq(source1, key1, index)
	default:
		panic(fmt.Errorf("unknown join type"))
	}
}

func generateLeftOuterJoinSeq[E1 any, E2 any](source1 source[E1], key1 KeyFunc[E1, interface{}],
	index map[interface{}][]E2) source[Pair[E1, E2]] {

	return func() iter.Seq2[Pair[E1, E2], error] {
		return func(yield func(Pair[E1, E2], error) bool) {
			var empty Pair[E1, E2]
			for element1, err := range source1() {
				if err != nil {
					yield(empty, err)
					return
				}
				element1Key := key1(element1)
				elements, found := index[element1Key]
				if !found {
					var emptyE2 E2
					if !yield(Pair[E1, E2]{Element1: element1, Element2: emptyE2}, nil) {
						return
					}
				} else {
					for _, element2 := range elements {
						if !yield(Pair[E1, E2]{Element1: element1, Element2: element2}, nil) {
							return
						}
					}
				}
			}
		}
	}
}

func createIndex[E any](source2 source[E], key2 KeyFunc[E, any]) (map[interface{}][]E, error) {
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

func generateInnerJoinSeq[E1 any, E2 any](source1 source[E1], key1 KeyFunc[E1, interface{}], index map[interface{}][]E2) source[Pair[E1, E2]] {
	return func() iter.Seq2[Pair[E1, E2], error] {
		return func(yield func(Pair[E1, E2], error) bool) {
			var empty Pair[E1, E2]
			for element1, err := range source1() {
				if err != nil {
					yield(empty, err)
					return
				}
				element1Key := key1(element1)
				elements, found := index[element1Key]
				if found {
					for _, element2 := range elements {
						if !yield(Pair[E1, E2]{Element1: element1, Element2: element2}, nil) {
							return
						}
					}
				}
			}
		}
	}
}
