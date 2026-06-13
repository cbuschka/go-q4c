package q4c

import (
	"github.com/stretchr/testify/require"

	"testing"
)

func TestUniOps_SelectAllToSlice(t *testing.T) {
	type Person struct {
		Id   int
		Name string
	}

	persons := []Person{
		{Id: 1, Name: "Jane"},
		{Id: 2, Name: "Tarzan"},
	}

	result, err := SelectFrom(persons).ToSlice()
	require.NoError(t, err)
	require.Equal(t, result, persons)
}

func TestUniOps_SelectWithFilterToSlice(t *testing.T) {
	type Person struct {
		Id   int
		Name string
	}

	persons := []Person{
		{Id: 1, Name: "Jane"},
		{Id: 2, Name: "Tarzan"},
	}

	result, err := SelectFrom(persons).Where(func(p Person) bool {
		return p.Id == 1
	}).ToSlice()
	require.NoError(t, err)
	require.Equal(t, result, []Person{persons[0]})
}
