package q4c

import (
	"github.com/cbuschka/go-q4c/internal"
	"github.com/cbuschka/go-q4c/types"
)

func SelectFrom[E1 any](elements []E1) types.FilterableUniSelect[E1] {
	return internal.SelectFrom(elements)
}
