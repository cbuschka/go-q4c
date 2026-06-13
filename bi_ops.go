package q4c

import (
	"github.com/cbuschka/go-q4c/internal"
	"github.com/cbuschka/go-q4c/types"
)

func NewBiSet[E1 any, E2 any]() types.BiSet[E1, E2] {
	return types.BiSet[E1, E2](&internal.BiSetImpl[E1, E2]{})
}
