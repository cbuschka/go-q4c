package q4c

import (
	"github.com/cbuschka/go-q4c/internal"
	"github.com/cbuschka/go-q4c/types"
)

func NewUniSet[E1 any]() types.UniSet[E1] {
	return types.UniSet[E1](&internal.UniSetImpl[E1]{})
}

func SelectFrom[E1 any](elements []E1) types.FilterableUniSelect[E1] {
	return NewUniSet[E1]().SelectFrom(elements)
}
