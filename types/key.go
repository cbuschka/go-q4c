package types

type KeyFunc[E any, K comparable] func(element1 E) K
